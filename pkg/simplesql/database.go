package simplesql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB         *sqlx.DB
	errHandler ErrHandler
}

type Option func(*Database)

func WithErrHandler(errHandler ErrHandler) Option {
	return func(d *Database) {
		d.errHandler = errHandler
	}
}

func NewDatabase(db *sqlx.DB, opts ...Option) Database {
	d := Database{
		DB:         db,
		errHandler: defaultErrHandler,
	}
	for _, opt := range opts {
		opt(&d)
	}

	return d
}

func (d *Database) Insert(
	ctx context.Context, execer sqlx.ExecerContext, tableName string, row interface{},
) error {
	// Deduce the column names and placeholders from the struct tags
	columnNames, placeholders := getColumnNamesAndPlaceholders(row)

	// Build the final query string
	query := fmt.Sprintf(`
		INSERT INTO %s
		(%s)
		VALUES (%s)
	`, tableName, columnNames, placeholders)

	// d.logger.Debug("InsertRow", "query", strings.ReplaceAll(query, "\n\t\t", " "))
	// Execute the query
	_, err := d.bindAndExec(ctx, execer, query, row)
	return d.errHandler(err)
}

func (d *Database) Get(
	ctx context.Context, tableName string, key interface{}, row interface{},
) error {
	// Deduce the column names for the SELECT statement
	columnNames, _ := getColumnNamesAndPlaceholders(row)
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE 1=1`, columnNames, tableName)

	// Prepare the parameters for the WHERE clause
	params := []interface{}{}
	keyValue := reflect.ValueOf(key)
	if keyValue.Kind() == reflect.Ptr {
		keyValue = keyValue.Elem()
	}
	keyType := keyValue.Type()

	// Iterate over the struct fields to build the WHERE clause
	for i := 0; i < keyType.NumField(); i++ {
		field := keyType.Field(i)
		fieldValue := keyValue.Field(i)

		// Check if the field is a pointer and skip if nil
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		// Use the "db" struct tag if present, otherwise default to field name
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = field.Name
		}

		// Add condition to query and append field value to params
		query += fmt.Sprintf(" AND %s = ?", columnName)
		params = append(params, fieldValue.Interface())
	}

	// Execute the query
	err := d.DB.GetContext(ctx, row, query, params...)
	return d.errHandler(err)
}

func (d *Database) Update(
	ctx context.Context, execer sqlx.ExecerContext, tableName string, key interface{}, fields interface{},
) error {
	version := uint64(0)
	versionSet := true
	versionField := reflect.ValueOf(key).FieldByName("Version")
	if !versionField.IsValid() {
		versionSet = false
	} else {
		version = versionField.Uint()
	}

	query := ""
	var updates []string
	params := map[string]interface{}{}

	if versionSet {
		query = fmt.Sprintf(`
			UPDATE %s
			SET version = :new_version
		`, tableName)

		params = map[string]interface{}{
			"new_version": version + 1,
		}
	} else {
		query = fmt.Sprintf(`UPDATE %s SET `, tableName)
	}

	// Use reflection to iterate over the fields and extract db tags and values
	v := reflect.ValueOf(fields)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		dbTag := fieldType.Tag.Get("db")
		attributeTag := fmt.Sprintf("update_%s", dbTag)

		// Check if the field is nil (for pointers), if not, add to updates
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			updates = append(updates, fmt.Sprintf("%s = :%s", dbTag, attributeTag))
			params[attributeTag] = field.Elem().Interface() // Dereference pointer and add value
		}
	}

	if len(updates) > 0 {
		if versionSet {
			query += ", " + strings.Join(updates, ", ")
		} else {
			query += strings.Join(updates, ", ")
		}
	}

	query += " WHERE 1=1"
	// Add the key fields to the WHERE clause
	keyValue := reflect.ValueOf(key)
	if keyValue.Kind() == reflect.Ptr {
		keyValue = keyValue.Elem()
	}
	keyType := keyValue.Type()

	for i := 0; i < keyType.NumField(); i++ {
		field := keyType.Field(i)
		fieldValue := keyValue.Field(i)

		// Check if the field is a pointer and skip if nil
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		// Use the "db" struct tag if present, otherwise default to field name
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = field.Name
		}

		// Add condition to query and append field value to params
		query += fmt.Sprintf(" AND %s = :%s", columnName, columnName)
		params[columnName] = fieldValue.Interface()
	}

	res, err := d.bindAndExec(ctx, execer, query, params)
	if err != nil {
		return d.errHandler(err)
	}
	return d.checkOptimisticLock(res)
}

func (d *Database) Delete(
	ctx context.Context, tableName string, key interface{},
) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE 1=1`, tableName)

	// Prepare the parameters for the WHERE clause
	params := []interface{}{}
	keyValue := reflect.ValueOf(key)
	if keyValue.Kind() == reflect.Ptr {
		keyValue = keyValue.Elem()
	}
	keyType := keyValue.Type()

	// Iterate over the struct fields to build the WHERE clause
	for i := 0; i < keyType.NumField(); i++ {
		field := keyType.Field(i)
		fieldValue := keyValue.Field(i)

		// Check if the field is a pointer and skip if nil
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		// Use the "db" struct tag if present, otherwise default to field name
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = field.Name
		}

		// Add condition to query and append field value to params
		query += fmt.Sprintf(" AND %s = ?", columnName)
		params = append(params, fieldValue.Interface())
	}

	// Execute the query
	_, err := d.DB.ExecContext(ctx, query, params...)
	return d.errHandler(err)
}

func (d *Database) List(
	ctx context.Context, tableName string, filters interface{}, result interface{},
) error {
	// Deduce the column names and placeholders from the struct tags
	columnNames, _ := getColumnNamesAndPlaceholders(result)
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE 1=1`, columnNames, tableName)
	params := map[string]interface{}{}

	// Use reflection to iterate over the filters struct and build query conditions
	v := reflect.ValueOf(filters)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		dbTag := fieldType.Tag.Get("db")

		if dbTag == "" {
			continue // Skip fields with no db tag
		}

		// Split the tag to handle operations (e.g., eq, lt, gt)
		tagParts := strings.Split(dbTag, ":")
		columnName := tagParts[0]
		operation := "none"

		if len(tagParts) > 1 {
			operation = tagParts[1] // Extract the operation from the tag
		}

		// Handle slice types (IN and NOT IN clauses)
		if field.Kind() == reflect.Slice && field.Len() > 0 {
			if strings.Contains(operation, "not_in") {
				query += fmt.Sprintf(" AND %s NOT IN (:%s)", columnName, fieldType.Name)
			} else {
				query += fmt.Sprintf(" AND %s IN (:%s)", columnName, fieldType.Name)
			}
			params[fieldType.Name] = field.Interface()

		} else if field.IsValid() && !isEmptyValue(field) {
			// Handle different operations
			switch operation {
			case "lt":
				query += fmt.Sprintf(" AND %s < :%s", columnName, fieldType.Name)
			case "gt":
				query += fmt.Sprintf(" AND %s > :%s", columnName, fieldType.Name)
			case "lte":
				query += fmt.Sprintf(" AND %s <= :%s", columnName, fieldType.Name)
			case "gte":
				query += fmt.Sprintf(" AND %s >= :%s", columnName, fieldType.Name)
			case "eq":
				query += fmt.Sprintf(" AND %s = :%s", columnName, fieldType.Name)

			}
			params[fieldType.Name] = field.Interface()
		}
	}

	// Handle limit if it's provided
	_, ok := t.FieldByName("Limit")
	if ok && v.FieldByName("Limit").Uint() > 0 {
		query += " LIMIT :limit"
		params["limit"] = v.FieldByName("Limit").Uint()
	}

	// Prepare the final query with expanded parameters
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return fmt.Errorf("failed to expand IN clause: %s, %w", err.Error(), ErrInternal)
	}

	// Expand IN clause
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return fmt.Errorf("failed to expand IN clause: %s, %w", err.Error(), ErrInternal)
	}

	// Rebind for the current SQL driver
	query = d.DB.Rebind(query)

	// Execute the query
	err = d.DB.SelectContext(ctx, result, query, args...)
	if err != nil {
		return d.errHandler(err)
	}

	return nil
}

// Helper function to check if a field is empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (d *Database) bindAndExec(
	ctx context.Context, execer sqlx.ExecerContext, query string, row interface{}) (sql.Result, error) {
	query, args, err := sqlx.Named(query, row)
	if err != nil {
		return nil, err
	}
	query = d.DB.Rebind(query)
	// d.logger.Debug("bindAndExec", "query", strings.ReplaceAll(query, "\n\t\t", " "))
	return execer.ExecContext(ctx, query, args...)
}

func (d *Database) checkOptimisticLock(res sql.Result) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return d.errHandler(err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected: %w", ErrInsertConflict)
	}
	return nil
}

// Helper function to get column names and placeholders from struct tags
func getColumnNamesAndPlaceholders(row interface{}) (string, string) {
	v := reflect.ValueOf(row)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if we received a slice
	if v.Kind() == reflect.Slice {
		// Check if the slice is empty
		if v.Len() == 0 {
			// Get the element type of the slice to deduce the columns, even if empty
			v = reflect.New(v.Type().Elem()).Elem() // Create a new instance of the element type
		} else {
			// If not empty, we can use the first element
			v = v.Index(0)
		}
	}

	t := v.Type()
	var columnNames []string
	var placeholders []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnNames = append(columnNames, dbTag)
			placeholders = append(placeholders, ":"+dbTag)
		}
	}

	return strings.Join(columnNames, ", "), strings.Join(placeholders, ", ")
}
