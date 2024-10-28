package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const tableInitializerTemplate = `
package tables

import "github.com/msanath/gondolf/pkg/simplesql"

func Initialize(simpleDB simplesql.Database) error {
	var schemaMigrations = []simplesql.Migration{}

	schemaMigrations = append(schemaMigrations, {{.AttributePrefix}}TableMigrations...)
	// ++ledgerbuilder:Migrations

	err := simpleDB.ApplyMigrations(schemaMigrations)
	if err != nil {
		return err
	}
	return nil
}
`

const tableTemplate = `
package tables

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/msanath/gondolf/pkg/simplesql"
)

var {{.AttributePrefix}}TableMigrations = []simplesql.Migration{
	{
		Version: 1, // Update the version number sequentially.
		Up: ` + "`" + `
			CREATE TABLE {{.TableName}} (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
				UNIQUE (id, name, is_deleted)
			);
		` + "`" + `,
		Down: ` + "`" + `
				DROP TABLE IF EXISTS {{.TableName}};
			` + "`" + `,
	},
}

type {{.RecordName}}Row struct {
	ID        string ` + "`" + `db:"id" orm:"op=create key=primary_key filter=In"` + "`" + `
	Version   uint64 ` + "`" + `db:"version" orm:"op=create,update"` + "`" + `
	Name      string ` + "`" + `db:"name" orm:"op=create composite_unique_key:Name,isDeleted filter=In"` + "`" + `
	IsDeleted bool   ` + "`" + `db:"is_deleted"` + "`" + `
	State     string ` + "`" + `db:"state" orm:"op=create,update filter=In,NotIn"` + "`" + `
	Message   string ` + "`" + `db:"message" orm:"op=create,update"` + "`" + `
}

type {{.RecordName}}TableUpdateFields struct {
	State   *string ` + "`" + `db:"state"` + "`" + `
	Message *string ` + "`" + `db:"message"` + "`" + `
}

type {{.RecordName}}TableSelectFilters struct {
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

type {{.RecordName}}Table struct {
	simplesql.Database
	tableName string
}

func New{{.RecordName}}Table(db simplesql.Database) *{{.RecordName}}Table {
	return &{{.RecordName}}Table{
		Database:  db,
		tableName: {{.AttributePrefix}}TableName,
	}
}

func (s *{{.RecordName}}Table) Insert(ctx context.Context, execer sqlx.ExecerContext, row {{.RecordName}}Row) error {
	return s.Database.InsertRow(ctx, execer, s.tableName, row)
}

func (s *{{.RecordName}}Table) GetByIDAndVersion(ctx context.Context, id string, version uint64, isDeleted bool) ({{.RecordName}}Row, error) {
	var row {{.RecordName}}Row
	err := s.Database.GetRowByID(ctx, id, version, isDeleted, s.tableName, &row)
	if err != nil {
		return {{.RecordName}}Row{}, err
	}
	return row, nil
}

func (s *{{.RecordName}}Table) GetByName(ctx context.Context, name string) ({{.RecordName}}Row, error) {
	var row {{.RecordName}}Row
	err := s.Database.GetRowByName(ctx, name, s.tableName, &row)
	if err != nil {
		return {{.RecordName}}Row{}, err
	}
	return row, nil
}

func (s *{{.RecordName}}Table) Update(
	ctx context.Context, execer sqlx.ExecerContext, id string, version uint64, updateFields {{.RecordName}}TableUpdateFields,
) error {
	return s.Database.UpdateRow(ctx, execer, id, version, s.tableName, updateFields)
}

func (s *{{.RecordName}}Table) Delete(ctx context.Context, execer sqlx.ExecerContext, id string, version uint64) error {
	return s.Database.MarkRowAsDeleted(ctx, execer, id, version, s.tableName)
}

func (s *{{.RecordName}}Table) List(ctx context.Context, filters {{.RecordName}}TableSelectFilters) ([]{{.RecordName}}Row, error) {
	var rows []{{.RecordName}}Row
	err := s.Database.SelectRows(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
`

func (o GenerateOptions) generateTables() error {
	fmt.Println("Generating tables")

	tablesPath := filepath.Join(o.DestinationPath, "internal", "sqlstorage", "tables")
	// Create destination path if it doesn't exist
	err := os.MkdirAll(tablesPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	initalizeFilePath := filepath.Join(tablesPath, "initialize.go")
	// If file is present, update at ++ledgerbuilder markers
	if _, err := os.Stat(initalizeFilePath); err == nil {
		fmt.Println("... existing new.go found. Updating new.go")
		err = o.addToMigrations(initalizeFilePath)
		if err != nil {
			return fmt.Errorf("failed to update new initialize file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	} else {
		fmt.Println("... creating initialize.go")
		err = executeTemplate("tableInitializerTemplate", tableInitializerTemplate, tablesPath, "initialize.go", o)
		if err != nil {
			return fmt.Errorf("failed to generate initalize file: %w", err)
		}
	}

	fileName := fmt.Sprintf("%s_table.go", o.TableName)
	fmt.Println("... creating ", fileName)
	err = executeTemplate("tablesTemplate", tableTemplate, tablesPath, fileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	return nil
}

// addToMigrations updates the existing 'initialize.go' file to add new migrations.
func (o GenerateOptions) addToMigrations(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Prepare the new lines to be added
	newMigration := fmt.Sprintf("schemaMigrations = append(schemaMigrations, %sTableMigrations...)", o.AttributePrefix)

	// Insert into the appropriate places
	updatedContent := insertAtPlaceholder(string(content), "// ++ledgerbuilder:Migrations", newMigration)

	// Write the updated content back to the file
	return os.WriteFile(filePath, []byte(updatedContent), 0644)
}
