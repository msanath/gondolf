package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const repoTemplate = `
package sqlstorage

import (
	"context"

	"github.com/msanath/gondolf/pkg/simplesql"
	"{{.GoModuleName}}/internal/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/internal/ledger/core"
	"{{.GoModuleName}}/internal/sqlstorage/tables"
)

// {{.AttributePrefix}}Storage is a concrete implementation of {{.RecordName}}Repository using sqlx
type {{.AttributePrefix}}Storage struct {
	simplesql.Database
	{{.AttributePrefix}}Table *tables.{{.RecordName}}Table
}

// new{{.RecordName}}Storage creates a new storage instance satisfying the {{.RecordName}}Repository interface
func new{{.RecordName}}Storage(db simplesql.Database) {{.PackageName}}.Repository {
	return &{{.AttributePrefix}}Storage{
		Database:     db,
		{{.AttributePrefix}}Table: tables.New{{.RecordName}}Table(db),
	}
}

func {{.AttributePrefix}}ModelToRow(model {{.PackageName}}.{{.RecordName}}Record) tables.{{.RecordName}}Row {
	return tables.{{.RecordName}}Row{
		ID:      model.Metadata.ID,
		Version: model.Metadata.Version,
		Name:    model.Name,
		State:   model.Status.State.ToString(),
		Message: model.Status.Message,
	}
}

func {{.AttributePrefix}}RowToModel(row tables.{{.RecordName}}Row) {{.PackageName}}.{{.RecordName}}Record {
	return {{.PackageName}}.{{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      row.ID,
			Version: row.Version,
		},
		Name: row.Name,
		Status: {{.PackageName}}.{{.RecordName}}Status{
			State:   {{.PackageName}}.{{.RecordName}}StateFromString(row.State),
			Message: row.Message,
		},
	}
}

func (s *{{.AttributePrefix}}Storage) Insert(ctx context.Context, record {{.PackageName}}.{{.RecordName}}Record) error {
	execer := s.DB
	err := s.{{.AttributePrefix}}Table.Insert(ctx, execer, {{.AttributePrefix}}ModelToRow(record))
	if err != nil {
		return errHandler(err)
	}
	return nil
}

func (s *{{.AttributePrefix}}Storage) GetByMetadata(ctx context.Context, metadata core.Metadata) ({{.PackageName}}.{{.RecordName}}Record, error) {
	row, err := s.{{.AttributePrefix}}Table.GetByIDAndVersion(ctx, metadata.ID, metadata.Version, metadata.IsDeleted)
	if err != nil {
		return {{.PackageName}}.{{.RecordName}}Record{}, errHandler(err)
	}
	return {{.AttributePrefix}}RowToModel(row), nil
}

func (s *{{.AttributePrefix}}Storage) GetByName(ctx context.Context, name string) ({{.PackageName}}.{{.RecordName}}Record, error) {
	row, err := s.{{.AttributePrefix}}Table.GetByName(ctx, name)
	if err != nil {
		return {{.PackageName}}.{{.RecordName}}Record{}, errHandler(err)
	}
	return {{.AttributePrefix}}RowToModel(row), nil
}

func (s *{{.AttributePrefix}}Storage) UpdateState(ctx context.Context, metadata core.Metadata, status {{.PackageName}}.{{.RecordName}}Status) error {
	execer := s.DB
	state := status.State.ToString()
	message := status.Message
	updateFields := tables.{{.RecordName}}TableUpdateFields{
		State:   &state,
		Message: &message,
	}
	return s.{{.AttributePrefix}}Table.Update(ctx, execer, metadata.ID, metadata.Version, updateFields)
}

func (s *{{.AttributePrefix}}Storage) Delete(ctx context.Context, metadata core.Metadata) error {
	execer := s.DB
	return s.{{.AttributePrefix}}Table.Delete(ctx, execer, metadata.ID, metadata.Version)
}

func (s *{{.AttributePrefix}}Storage) List(ctx context.Context, filters {{.PackageName}}.{{.RecordName}}ListFilters) ([]{{.PackageName}}.{{.RecordName}}Record, error) {
	dbFilters := tables.{{.RecordName}}TableSelectFilters{
		IDIn:           append([]string{}, filters.IDIn...),
		NameIn:         append([]string{}, filters.NameIn...),
		VersionGte:     filters.VersionGte,
		VersionLte:     filters.VersionLte,
		VersionEq:      filters.VersionEq,
		IncludeDeleted: filters.IncludeDeleted,
		Limit:          filters.Limit,
	}
	for _, state := range filters.StateIn {
		dbFilters.StateIn = append(dbFilters.StateIn, state.ToString())
	}
	for _, state := range filters.StateNotIn {
		dbFilters.StateNotIn = append(dbFilters.StateNotIn, state.ToString())
	}

	rows, err := s.{{.AttributePrefix}}Table.List(ctx, dbFilters)
	if err != nil {
		return nil, err
	}
	var records []{{.PackageName}}.{{.RecordName}}Record
	for _, row := range rows {
		records = append(records, {{.AttributePrefix}}RowToModel(row))
	}
	return records, nil
}
`

const repoTestTemplate = `
package sqlstorage_test

import (
	"context"
	"fmt"
	"testing"

	"{{.GoModuleName}}/internal/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/internal/ledger/core"
	ledgererrors "{{.GoModuleName}}/internal/ledger/errors"
	"{{.GoModuleName}}/internal/sqlstorage"

	"github.com/msanath/gondolf/pkg/simplesql/test"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const {{.AttributePrefix}}idPrefix = "{{.PackageName}}"

func Test{{.RecordName}}RecordLifecycle(t *testing.T) {
	db, err := test.NewTestSQLiteDB()
	require.NoError(t, err)

	storage, err := sqlstorage.NewSQLStorage(db, true)
	require.NoError(t, err)

	testRecord := {{.PackageName}}.{{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s1", {{.AttributePrefix}}idPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s1", {{.AttributePrefix}}idPrefix),
		Status: {{.PackageName}}.{{.RecordName}}Status{
			State:   {{.PackageName}}.{{.RecordName}}StateActive,
			Message: "",
		},
	}
	repo := storage.{{.RecordName}}

	test{{.RecordName}}CRUD(t, repo, testRecord)
}

func test{{.RecordName}}CRUD(t *testing.T, repo {{.PackageName}}.Repository, testRecord {{.PackageName}}.{{.RecordName}}Record) {
	ctx := context.Background()
	var err error

	t.Run("Insert Success", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.NoError(t, err)
	})

	t.Run("Insert Duplicate Failure", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
	})

	t.Run("Get By Metadata Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByMetadata(ctx, testRecord.Metadata)
		require.NoError(t, err)
		require.Equal(t, testRecord, receivedRecord)
	})

	t.Run("Get By Metadata failure", func(t *testing.T) {
		metadata := core.Metadata{
			ID:      testRecord.Metadata.ID,
			Version: 2, // Different version
		}
		_, err := repo.GetByMetadata(ctx, metadata)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code, err.Error())
	})

	t.Run("Get By Name Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(testRecord, receivedRecord))
	})

	t.Run("Get By Name Failure", func(t *testing.T) {
		_, err := repo.GetByName(ctx, "unknown")
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code, err.Error())
	})

	t.Run("Update State Success", func(t *testing.T) {
		status := {{.PackageName}}.{{.RecordName}}Status{
			State:   "error",
			Message: "Needs attention",
		}

		err = repo.UpdateState(ctx, testRecord.Metadata, status)
		require.NoError(t, err)

		updatedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Equal(t, status, updatedRecord.Status)
		require.Equal(t, testRecord.Metadata.Version+1, updatedRecord.Metadata.Version)
		testRecord = updatedRecord
	})

	t.Run("Delete Success", func(t *testing.T) {
		err = repo.Delete(ctx, testRecord.Metadata)
		require.NoError(t, err)

		_, err = repo.GetByName(ctx, testRecord.Name)
		require.Error(t, err)
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
	})

	t.Run("Create More Resources", func(t *testing.T) {
		// Create 10 records.
		for i := range 10 {
			newRecord := testRecord
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", {{.AttributePrefix}}idPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", {{.AttributePrefix}}idPrefix, i+1)
			newRecord.Status.State = {{.PackageName}}.{{.RecordName}}StateActive
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", {{.AttributePrefix}}idPrefix, i+1)

			if (i+1)%2 == 0 {
				newRecord.Status.State = {{.PackageName}}.{{.RecordName}}StateInActive
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", {{.AttributePrefix}}idPrefix, i+1)
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

	t.Run("List", func(t *testing.T) {
		records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{})
		require.NoError(t, err)
		require.Len(t, records, 10)

		receivedIDs := []string{}
		for _, record := range records {
			receivedIDs = append(receivedIDs, record.Metadata.ID)

		}
		expectedIDs := []string{}
		for i := range 10 {
			expectedIDs = append(expectedIDs, fmt.Sprintf("%s-%d", {{.AttributePrefix}}idPrefix, i+1))

		}
		require.ElementsMatch(t, expectedIDs, receivedIDs)
		allRecords := records

		t.Run("List Success With Filter", func(t *testing.T) {
			records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				StateIn: []{{.PackageName}}.{{.RecordName}}State{{"{"}}{{.PackageName}}.{{.RecordName}}StateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, {{.PackageName}}.{{.RecordName}}StateActive, record.Status.State)
			}
		})

		t.Run("List with Names Filter", func(t *testing.T) {
			records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				NameIn: []string{allRecords[0].Name, allRecords[1].Name, allRecords[2].Name},
			})
			require.NoError(t, err)
			require.Len(t, records, 3)

			// Check if the returned records are the same as the first 3 computeCapabilitys.
			for i, record := range records {
				require.Equal(t, allRecords[i], record)
			}
		})

		t.Run("List with Limit", func(t *testing.T) {
			records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				Limit: 3,
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})

		t.Run("List with IncludeDeleted", func(t *testing.T) {
			err = repo.Delete(ctx, allRecords[0].Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				StateNotIn: []{{.PackageName}}.{{.RecordName}}State{{"{"}}{{.PackageName}}.{{.RecordName}}StateActive},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, {{.PackageName}}.{{.RecordName}}StateInActive, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := {{.PackageName}}.{{.RecordName}}Status{
				State:   {{.PackageName}}.{{.RecordName}}StatePending,
				Message: "Needs attention",
			}

			err = repo.UpdateState(ctx, allRecords[1].Metadata, status)
			require.NoError(t, err)
			ve := uint64(1)
			records, err := repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 1)

			ve += 1
			records, err = repo.List(ctx, {{.PackageName}}.{{.RecordName}}ListFilters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})
	})
}
`

func (o GenerateOptions) generateRepo() error {
	fmt.Println("Generating storage repo components")

	sqlStoragePath := filepath.Join(o.DestinationPath, "internal", "sqlstorage")
	// Create destination path if it doesn't exist
	err := os.MkdirAll(sqlStoragePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	fileName := fmt.Sprintf("%s_repo.go", o.TableName)
	fmt.Println("... creating ", fileName)
	err = executeTemplate("repoTemplate", repoTemplate, sqlStoragePath, fileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	// Generate the test file
	testFileName := fmt.Sprintf("%s_repo_test.go", o.TableName)
	fmt.Println("... creating ", testFileName)
	err = executeTemplate("storageTestTemplate", repoTestTemplate, sqlStoragePath, testFileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	return nil
}
