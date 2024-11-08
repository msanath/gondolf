package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const newFileTemplate = `
package sqlstorage

import (
	"{{.GoModuleName}}/internal/ledger/{{.PackageName}}"
	ledgererrors "{{.GoModuleName}}/internal/ledger/errors"
	"{{.GoModuleName}}/internal/sqlstorage/tables"

	"github.com/msanath/gondolf/pkg/simplesql"

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

	var errHandler simplesql.ErrHandler
	if is_sqlite {
		errHandler = simplesql.SQLiteErrHandler
	} else {
		errHandler = simplesql.MySQLErrHandler
	}

	simpleDB := simplesql.NewDatabase(
		db, simplesql.WithErrHandler(errHandler),
	)
	err := tables.Initialize(simpleDB)
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
		return ledgererrors.NewLedgerError(ledgererrors.ErrRecordNotFound, "Record not found.")
	case simplesql.ErrInsertConflict:
		return ledgererrors.NewLedgerError(ledgererrors.ErrRecordInsertConflict, "Duplicate entry, record already exists.")
	case simplesql.ErrInternal:
		return ledgererrors.NewLedgerError(ledgererrors.ErrRepositoryInternal, "Internal error.")
	default:
		return err
	}
}
`

const testSQLStorageTemplate = `
// Code generated by ledger-builder. DO NOT EDIT.

package test

import (
	"testing"

	"{{.GoModuleName}}/internal/sqlstorage"

	simplesqltest "github.com/msanath/gondolf/pkg/simplesql/test"
	"github.com/stretchr/testify/require"
)

func TestSQLStorage(t *testing.T) *sqlstorage.SQLStorage {
	db, err := simplesqltest.NewTestSQLiteDB()
	require.NoError(t, err)
	storage, err := sqlstorage.NewSQLStorage(db, false)
	require.NoError(t, err)
	return storage
}
`

func (o GenerateOptions) generateStorageCommon() error {
	fmt.Println("Generating storage components")

	sqlStoragePath := filepath.Join(o.DestinationPath, "internal", "sqlstorage")
	// Create destination path if it doesn't exist
	err := os.MkdirAll(sqlStoragePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	newFilePath := filepath.Join(sqlStoragePath, "new.go")
	// If file is present, update at ++ledgerbuilder markers
	if _, err := os.Stat(newFilePath); err == nil {
		fmt.Println("... existing new.go found. Updating new.go")
		err = o.updateNewStorageFile(newFilePath)
		if err != nil {
			return fmt.Errorf("failed to update new storage file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	} else {
		fmt.Println("... creating new.go")
		err = executeTemplate("newFileTemplate", newFileTemplate, sqlStoragePath, "new.go", o)
		if err != nil {
			return fmt.Errorf("failed to generate new storage file: %w", err)
		}
	}

	// Generate the test file
	testPath := filepath.Join(sqlStoragePath, "test", "new.go")
	if _, err := os.Stat(newFilePath); err == nil {
		fmt.Println("... existing test/new.go found. Skipping")
		return nil
	}
	err = os.MkdirAll(testPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create test path: %w", err)
	}
	err = executeTemplate("testSQLStorageTemplate", testSQLStorageTemplate, testPath, "new.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate test file: %w", err)
	}

	return nil
}

// updateNewStorageFile updates the existing 'new.go' file to add new repositories and migrations.
func (o GenerateOptions) updateNewStorageFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Prepare the new lines to be added
	newRepo := fmt.Sprintf("%s %s.Repository", o.RecordName, o.PackageName)
	newMigration := fmt.Sprintf("schemaMigrations = append(schemaMigrations, %sTableMigrations...)", o.AttributePrefix)
	newRepoInstance := fmt.Sprintf("%s: new%sStorage(simpleDB),", o.RecordName, o.RecordName)
	newImport := fmt.Sprintf("\"%s/internal/ledger/%s\"", o.GoModuleName, o.PackageName)

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
