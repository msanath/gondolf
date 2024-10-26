package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const newFileTemplate = `
package sqlstorage

import (
	coreerrors "{{.GoModuleName}}/core/errors"
	"{{.GoModuleName}}/core/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/pkg/simplesql"
	// ++ledgerbuilder:Imports

	"github.com/jmoiron/sqlx"
)

type SQLStorage struct {
	{{.RecordName}} {{.PackageName}}.Repository
	// ++ledgerbuilder:RepositoryInterface
}

func NewSQLStorage(
	db *sqlx.DB,
	is_sqlite bool,
) (*SQLStorage, error) {
	var schemaMigrations = []simplesql.Migration{}

	schemaMigrations = append(schemaMigrations, {{.AttributePrefix}}TableMigrations...)
	// ++ledgerbuilder:Migrations

	var errHandler simplesql.ErrHandler
	if is_sqlite {
		errHandler = simplesql.SQLiteErrHandler
	} else {
		errHandler = simplesql.MySQLErrHandler
	}

	simpleDB := simplesql.NewDatabase(
		db, simplesql.WithErrHandler(errHandler),
	)
	err := simpleDB.ApplyMigrations(schemaMigrations)
	if err != nil {
		return nil, err
	}

	return &SQLStorage{
		{{.RecordName}}: new{{.RecordName}}Storage(simpleDB),
		// ++ledgerbuilder:RepoInstance
	}, nil
}

func errHandler(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case simplesql.ErrRecordNotFound:
		return coreerrors.NewStorageError(coreerrors.RepositoryErrorRecordNotFound, "Record not found.")
	case simplesql.ErrInsertConflict:
		return coreerrors.NewStorageError(coreerrors.RepositoryErrorRecordInsertConflict, "Duplicate entry, record already exists.")
	case simplesql.ErrInternal:
		return coreerrors.NewStorageError(coreerrors.RepositoryErrorInternal, "Internal error.")
	default:
		return err
	}
}
`

const storageTemplate = `
package sqlstorage

import (
	"context"

	"{{.GoModuleName}}/core/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/core"
	"{{.GoModuleName}}/pkg/simplesql"
)

var {{.AttributePrefix}}TableMigrations = []simplesql.Migration{
	{
		Version: 1,  // Update the version number sequentially.
		Up: ` + "`" + `
			CREATE TABLE {{.TableName}} (
				id VARCHAR(255) NOT NULL,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				UNIQUE (name, is_deleted)
			);
		` + "`" + `,
		Down: ` + "`" + `
				DROP TABLE IF EXISTS {{.TableName}};
			` + "`" + `,
	},
}

type {{.AttributePrefix}}Row struct {
	ID        string ` + "`" + `db:"id" orm:"op=create key=primary_key filter=In"` + "`" + `
	Version   uint64 ` + "`" + `db:"version" orm:"op=create,update"` + "`" + `
	Name      string ` + "`" + `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"` + "`" + `
	IsDeleted bool   ` + "`" + `db:"is_deleted"` + "`" + `
	State     string ` + "`" + `db:"state" orm:"op=create,update filter=In,NotIn"` + "`" + `
	Message   string ` + "`" + `db:"message" orm:"op=create,update"` + "`" + `
}

type {{.AttributePrefix}}UpdateFields struct {
	State   *string ` + "`" + `db:"state"` + "`" + `
	Message *string ` + "`" + `db:"message"` + "`" + `
}

type {{.AttributePrefix}}SelectFilters struct {
	IDIn       []string ` + "`" + `db:"id:in"` + "`" + `        // IN condition
	NameIn     []string ` + "`" + `db:"name:in"` + "`" + `      // IN condition
	StateIn    []string ` + "`" + `db:"state:in"` + "`" + `     // IN condition
	StateNotIn []string ` + "`" + `db:"state:not_in"` + "`" + ` // NOT IN condition
	VersionGte *uint64  ` + "`" + `db:"version:gte"` + "`" + `  // Greater than or equal condition
	VersionLte *uint64  ` + "`" + `db:"version:lte"` + "`" + `  // Less than or equal condition
	VersionEq  *uint64  ` + "`" + `db:"version:eq"` + "`" + `   // Equal condition

	IncludeDeleted bool   ` + "`" + `db:"include_deleted"` + "`" + ` // Special boolean handling
	Limit          uint32 ` + "`" + `db:"limit"` + "`" + `
}

const {{.AttributePrefix}}TableName = "{{.TableName}}"

func {{.AttributePrefix}}ModelToRow({{.AttributePrefix}} {{.PackageName}}.{{.RecordName}}Record) {{.AttributePrefix}}Row {
	return {{.AttributePrefix}}Row{
		ID:      {{.AttributePrefix}}.Metadata.ID,
		Version: {{.AttributePrefix}}.Metadata.Version,
		Name:    {{.AttributePrefix}}.Name,
		State:   {{.AttributePrefix}}.Status.State,
		Message: {{.AttributePrefix}}.Status.Message,
	}
}

func {{.AttributePrefix}}RowToModel(row {{.AttributePrefix}}Row) {{.PackageName}}.{{.RecordName}}Record {
	return {{.PackageName}}.{{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: core.Status{
			State:   row.State,
			Message: row.Message,
		},
	}
}

// {{.AttributePrefix}}Storage is a concrete implementation of {{.RecordName}}Repository using sqlx
type {{.AttributePrefix}}Storage struct {
	simplesql.Database
	tableName  string
	modelToRow func({{.PackageName}}.{{.RecordName}}Record) {{.AttributePrefix}}Row
	rowToModel func({{.AttributePrefix}}Row) {{.PackageName}}.{{.RecordName}}Record
}

// new{{.RecordName}}Storage creates a new storage instance satisfying the {{.RecordName}}Repository interface
func new{{.RecordName}}Storage(db simplesql.Database) {{.PackageName}}.Repository {
	return &{{.AttributePrefix}}Storage{
		Database:           db,
		tableName:    {{.AttributePrefix}}TableName,
		modelToRow:   {{.AttributePrefix}}ModelToRow,
		rowToModel:   {{.AttributePrefix}}RowToModel,
	}
}

func (s *{{.AttributePrefix}}Storage) Insert(ctx context.Context, record {{.PackageName}}.{{.RecordName}}Record) error {
	row := s.modelToRow(record)
	err := s.Database.InsertRow(ctx, s.DB, s.tableName, row)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *{{.AttributePrefix}}Storage) GetByMetadata(ctx context.Context, metadata core.Metadata) ({{.PackageName}}.{{.RecordName}}Record, error) {
	var row {{.AttributePrefix}}Row
	err := s.Database.GetRowByID(ctx, metadata.ID, metadata.Version, metadata.IsDeleted, s.tableName, &row)
	if err != nil {
		return {{.PackageName}}.{{.RecordName}}Record{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *{{.AttributePrefix}}Storage) GetByName(ctx context.Context, {{.AttributePrefix}}Name string) ({{.PackageName}}.{{.RecordName}}Record, error) {
	var row {{.AttributePrefix}}Row
	err := s.Database.GetRowByName(ctx, {{.AttributePrefix}}Name, s.tableName, &row)
	if err != nil {
		return {{.PackageName}}.{{.RecordName}}Record{}, errHandler(err)
	}

	return s.rowToModel(row), nil
}

func (s *{{.AttributePrefix}}Storage) UpdateState(ctx context.Context, metadata core.Metadata, status core.Status) error {
	err := s.Database.UpdateRow(ctx, s.DB, metadata.ID, metadata.Version, s.tableName, {{.AttributePrefix}}UpdateFields{
		State:   &status.State,
		Message: &status.Message,
	})
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *{{.AttributePrefix}}Storage) Delete(ctx context.Context, metadata core.Metadata) error {
	err := s.Database.MarkRowAsDeleted(ctx, s.DB, metadata.ID, metadata.Version, s.tableName)
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *{{.AttributePrefix}}Storage) List(ctx context.Context, filters core.Filters) ([]{{.PackageName}}.{{.RecordName}}Record, error) {
	dbFilters := {{.AttributePrefix}}SelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		StateIn:        append([]string{}, filters.StateIn...),
		StateNotIn:     append([]string{}, filters.StateNotIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}

	var rows []{{.AttributePrefix}}Row
	err := s.Database.SelectRows(ctx, s.tableName, dbFilters, &rows)
	if err != nil {
		return nil, err
	}

	var records []{{.PackageName}}.{{.RecordName}}Record
	for _, row := range rows {
		records = append(records, s.rowToModel(row))
	}

	return records, nil
}
`

func (o generateOptions) generateStorage() error {
	fmt.Println("Generating storage components")

	subPath := filepath.Join(o.DestinationPath, "internal", "sqlstorage")
	// Create destination path if it doesn't exist
	err := os.MkdirAll(subPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	newFilePath := filepath.Join(subPath, "new.go")
	// If file is present, update at ++ledgerbuilder markers
	if _, err := os.Stat(newFilePath); err == nil {
		fmt.Println("... existing new.go found. Updating new.go")
		err = o.updateNewFile(newFilePath)
		if err != nil {
			return fmt.Errorf("failed to update new storage file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	} else {
		fmt.Println("... creating new.go")
		err = executeTemplate("newFileTemplate", newFileTemplate, subPath, "new.go", o)
		if err != nil {
			return fmt.Errorf("failed to generate new storage file: %w", err)
		}
	}

	fileName := fmt.Sprintf("%s_table.go", strings.ToLower(o.RecordName))
	fmt.Println("... creating ", fileName)
	err = executeTemplate("storageTemplate", storageTemplate, subPath, fileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	return nil
}

// updateNewFile updates the existing 'new.go' file to add new repositories and migrations.
func (o generateOptions) updateNewFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Prepare the new lines to be added
	newRepo := fmt.Sprintf("%s %s.Repository", o.RecordName, o.PackageName)
	newMigration := fmt.Sprintf("schemaMigrations = append(schemaMigrations, %sTableMigrations...)", o.AttributePrefix)
	newRepoInstance := fmt.Sprintf("%s: new%sStorage(simpleDB),", o.RecordName, o.RecordName)
	newImport := fmt.Sprintf("\"%s/core/ledger/%s\"", o.GoModuleName, o.PackageName)

	// Insert into the appropriate places
	updatedContent := insertAtPlaceholder(string(content), "// ++ledgerbuilder:RepositoryInterface", newRepo)
	updatedContent = insertAtPlaceholder(updatedContent, "// ++ledgerbuilder:Migrations", newMigration)
	updatedContent = insertAtPlaceholder(updatedContent, "// ++ledgerbuilder:RepoInstance", newRepoInstance)
	updatedContent = insertAtPlaceholder(updatedContent, "// ++ledgerbuilder:Imports", newImport)

	// Write the updated content back to the file
	return os.WriteFile(filePath, []byte(updatedContent), 0644)
}

// insertAtPlaceholder inserts a new line of code at the first occurrence of the placeholder
func insertAtPlaceholder(content, placeholder, addition string) string {
	if strings.Contains(content, placeholder) {
		parts := strings.SplitAfter(content, placeholder)
		// Remove the placeholder from the first part and insert the new line
		newContent := strings.Replace(parts[0], placeholder, "", 1)
		return newContent + addition + "\n" + placeholder + parts[1]
	}
	return content
}
