package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const recordTemplate = `
package {{.PackageName}}

import "{{.GoModuleName}}/core"

// {{.RecordName}} is a representation of the {{.RecordName}} of an application.
type {{.RecordName}}Record struct {
	Metadata core.Metadata // Metadata is the metadata that identifies the {{.RecordName}}. It is a combination of the {{.RecordName}}'s name and version.
	Name     string        // Name is the name of the {{.RecordName}}.
	Status   core.Status   // Status is the status of the {{.RecordName}}.
}

// {{.RecordName}}State is the state of a {{.RecordName}}.
type {{.RecordName}}State string

const (
	{{.RecordName}}StateActive   {{.RecordName}}State = "{{.RecordName}}State_ACTIVE"
	{{.RecordName}}StateInActive {{.RecordName}}State = "{{.RecordName}}State_INACTIVE"
)

// ToString returns the string representation of the {{.RecordName}}State.
func (cs {{.RecordName}}State) ToString() string {
	return string(cs)
}

// FromString converts a string into a {{.RecordName}}State if valid, otherwise returns an error.
func {{.RecordName}}StateFromString(s string) {{.RecordName}}State {
	switch s {
	case string({{.RecordName}}StateActive):
		return {{.RecordName}}StateActive
	case string({{.RecordName}}StateInActive):
		return {{.RecordName}}StateInActive
	default:
		return "unknown"
	}
}
`

const ledgerTemplate = `
package {{.PackageName}}

import (
	"context"
	"{{.GoModuleName}}/core"
)

// Ledger provides the methods for managing {{.RecordName}} records.
type Ledger interface {
	// Create creates a new {{.RecordName}}.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a {{.RecordName}} by its metadata.
	GetByMetadata(context.Context, *core.Metadata) (*GetResponse, error)
	// GetByName retrieves a {{.RecordName}} by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateState updates the state and message of an existing {{.RecordName}}.
	UpdateState(context.Context, *UpdateStateRequest) (*UpdateResponse, error)
	// List returns a list of {{.RecordName}} that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
}

// CreateRequest represents the {{.RecordName}} creation request.
type CreateRequest struct {
	Name  string
	Type  string
	Score int
}

// CreateResponse represents the response after creating a new {{.RecordName}}.
type CreateResponse struct {
	Record {{.RecordName}}Record
}

// UpdateStateRequest represents the request to update the state and message of a {{.RecordName}}.
type UpdateStateRequest struct {
	Metadata core.Metadata
	Status   core.Status
}

// GetResponse represents the response for fetching a {{.RecordName}}.
type GetResponse struct {
	Record {{.RecordName}}Record
}

// UpdateResponse represents the response after updating the state of a {{.RecordName}}.
type UpdateResponse struct {
	Record {{.RecordName}}Record
}

// ListRequest represents the request to list {{.RecordName}}s with filters.
type ListRequest struct {
	core.Filters
}

// {{.RecordName}}SelectFilters contains filters for querying the {{.RecordName}} table.
type {{.RecordName}}ListFilters struct {
	IDIn       []string
	NameIn     []string
	StateIn    []{{.RecordName}}State
	StateNotIn []{{.RecordName}}State

	IncludeDeleted bool
	Limit          uint32
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []{{.RecordName}}Record
}

// Repository provides the methods that the storage layer must implement to support the ledger.
type Repository interface {
	Insert(ctx context.Context, cluster {{.RecordName}}Record) error
	GetByMetadata(ctx context.Context, metadata core.Metadata) ({{.RecordName}}Record, error)
	GetByName(ctx context.Context, name string) ({{.RecordName}}Record, error)
	UpdateState(ctx context.Context, metadata core.Metadata, state core.Status) error
	Delete(ctx context.Context, metadata core.Metadata) error
	List(ctx context.Context, filters core.Filters) ([]{{.RecordName}}Record, error)
}

// ledger implements the Ledger interface.
type ledger struct {
	repo Repository
}

// NewLedger creates a new Ledger instance.
func NewLedger(repo Repository) Ledger {
	return &ledger{repo: repo}
}

// Create creates a new {{.RecordName}}.
func (l *ledger) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	rec := {{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      req.Name,
			Version: 0,
		},
		Name: req.Name,
		Status: core.Status{
			State:   {{.RecordName}}StateActive.ToString(),
			Message: "",
		},
	}

	err := l.repo.Insert(ctx, rec)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		Record: rec,
	}, nil
}

// GetByMetadata retrieves a {{.RecordName}} by its metadata.
func (l *ledger) GetByMetadata(ctx context.Context, metadata *core.Metadata) (*GetResponse, error) {
	record, err := l.repo.GetByMetadata(ctx, *metadata)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

// GetByName retrieves a {{.RecordName}} by its name.
func (l *ledger) GetByName(ctx context.Context, name string) (*GetResponse, error) {
	record, err := l.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

// UpdateState updates the state and message of an existing {{.RecordName}}.
func (l *ledger) UpdateState(ctx context.Context, req *UpdateStateRequest) (*UpdateResponse, error) {
	err := l.repo.UpdateState(ctx, req.Metadata, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByMetadata(ctx, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Record: record,
	}, nil
}

// List returns a list of {{.RecordName}}s that match the provided filters.
func (l *ledger) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	records, err := l.repo.List(ctx, req.Filters)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Records: records,
	}, nil
}
`

const repoTestTemplate = `
package {{.PackageName}}_test

import (
	"context"
	"fmt"
	"{{.GoModuleName}}/core"
	coreerrors "{{.GoModuleName}}/core/errors"
	"{{.GoModuleName}}/internal/sqlstorage"
	"{{.GoModuleName}}/core/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/pkg/simplesql/test"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const idPrefix = "{{.PackageName}}"

func Test{{.RecordName}}RecordLifecycle(t *testing.T) {
	db, err := test.NewTestSQLiteDB()
	require.NoError(t, err)

	storage, err := sqlstorage.NewSQLStorage(db, true)
	require.NoError(t, err)

	testRecord := {{.PackageName}}.{{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      fmt.Sprintf("%s1", idPrefix),
			Version: 1,
		},
		Name: fmt.Sprintf("%s1", idPrefix),
		Status: core.Status{
			State:   "active-status",
			Message: "",
		},
	}
	repo := storage.{{.RecordName}}

	testCRUD(t, repo, testRecord)
}

func testCRUD(t *testing.T, repo {{.PackageName}}.Repository, testRecord {{.PackageName}}.{{.RecordName}}Record) {
	ctx := context.Background()
	var err error

	t.Run("Insert Success", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.NoError(t, err)
	})

	t.Run("Insert Duplicate Failure", func(t *testing.T) {
		err = repo.Insert(ctx, testRecord)
		require.Error(t, err)
		require.Equal(t, coreerrors.RepositoryErrorRecordInsertConflict, err.(coreerrors.StorageError).Code)
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
		require.Equal(t, coreerrors.RepositoryErrorRecordNotFound, err.(coreerrors.StorageError).Code, err.Error())
	})

	t.Run("Get By Name Success", func(t *testing.T) {
		receivedRecord, err := repo.GetByName(ctx, testRecord.Name)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(testRecord, receivedRecord))
	})

	t.Run("Get By Name Failure", func(t *testing.T) {
		_, err := repo.GetByName(ctx, "unknown")
		require.Error(t, err)
		require.Equal(t, coreerrors.RepositoryErrorRecordNotFound, err.(coreerrors.StorageError).Code, err.Error())
	})

	t.Run("Update State Success", func(t *testing.T) {
		status := core.Status{
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
		require.Equal(t, coreerrors.RepositoryErrorRecordNotFound, err.(coreerrors.StorageError).Code)
	})

	activeState := "active"
	inactiveState := "inactive"
	t.Run("Create More Resources", func(t *testing.T) {
		// Create 10 records.
		for i := range 10 {
			newRecord := testRecord
			newRecord.Metadata.ID = fmt.Sprintf("%s-%d", idPrefix, i+1)
			newRecord.Metadata.Version = 0
			newRecord.Name = fmt.Sprintf("%s-%d", idPrefix, i+1)
			newRecord.Status.State = activeState
			newRecord.Status.Message = fmt.Sprintf("%s-%d is active", idPrefix, i+1)

			if (i+1)%2 == 0 {
				newRecord.Status.State = inactiveState
				newRecord.Status.Message = fmt.Sprintf("%s-%d is inactive", idPrefix, i+1)
			}

			err = repo.Insert(ctx, newRecord)
			require.NoError(t, err)
		}
	})

	t.Run("List", func(t *testing.T) {
		records, err := repo.List(ctx, core.Filters{})
		require.NoError(t, err)
		require.Len(t, records, 10)

		receivedIDs := []string{}
		for _, record := range records {
			receivedIDs = append(receivedIDs, record.Metadata.ID)

		}
		expectedIDs := []string{}
		for i := range 10 {
			expectedIDs = append(expectedIDs, fmt.Sprintf("%s-%d", idPrefix, i+1))

		}
		require.ElementsMatch(t, expectedIDs, receivedIDs)
		allRecords := records

		t.Run("List Success With Filter", func(t *testing.T) {
			records, err := repo.List(ctx, core.Filters{
				StateIn: []string{activeState},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, activeState, record.Status.State)
			}
		})

		t.Run("List with Names Filter", func(t *testing.T) {
			records, err := repo.List(ctx, core.Filters{
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
			records, err := repo.List(ctx, core.Filters{
				Limit: 3,
			})
			require.NoError(t, err)
			require.Len(t, records, 3)
		})

		t.Run("List with IncludeDeleted", func(t *testing.T) {
			err = repo.Delete(ctx, allRecords[0].Metadata)
			require.NoError(t, err)

			records, err := repo.List(ctx, core.Filters{
				IncludeDeleted: true,
			})
			require.NoError(t, err)
			require.Len(t, records, 11)
		})

		t.Run("List with StateNotIn", func(t *testing.T) {
			records, err := repo.List(ctx, core.Filters{
				StateIn: []string{inactiveState},
			})
			require.NoError(t, err)
			require.Len(t, records, 5)
			for _, record := range records {
				require.Equal(t, inactiveState, record.Status.State)
			}
		})

		t.Run("Update State and check version", func(t *testing.T) {
			status := core.Status{
				State:   "error",
				Message: "Needs attention",
			}

			err = repo.UpdateState(ctx, allRecords[1].Metadata, status)
			require.NoError(t, err)
			ve := uint64(1)
			records, err := repo.List(ctx, core.Filters{
				VersionEq: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 1)

			ve += 1
			records, err = repo.List(ctx, core.Filters{
				VersionGte: &ve,
			})
			require.NoError(t, err)
			require.Len(t, records, 0)
		})
	})
}
`

func (o generateOptions) generateLedgerRecord() error {
	ledgerSubPath := filepath.Join(o.DestinationPath, "core", "ledger", strings.ToLower(o.RecordName))
	fmt.Println("Generating Ledger at ", ledgerSubPath)
	// Create destination path if it doesn't exist
	err := os.MkdirAll(ledgerSubPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	fmt.Println("... creating record.go")
	err = executeTemplate("recordTemplate", recordTemplate, ledgerSubPath, "record.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	fmt.Println("... creating ledger.go")
	err = executeTemplate("ledgerTemplate", ledgerTemplate, ledgerSubPath, "ledger.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate ledger file: %w", err)
	}

	fmt.Println("... creating repo_test.go")
	err = executeTemplate("repoTest", repoTestTemplate, ledgerSubPath, "repo_test.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate ledger file: %w", err)
	}

	return nil
}
