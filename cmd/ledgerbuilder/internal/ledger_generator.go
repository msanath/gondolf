package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const interfaceTemplate = `
package {{.PackageName}}

import (
	"context"

	"{{.GoModuleName}}/internal/ledger/core"
)

// {{.RecordName}} is a representation of the {{.RecordName}} of an application.
type {{.RecordName}}Record struct {
	Metadata core.Metadata // Metadata is the metadata that identifies the {{.RecordName}}. It is a combination of the {{.RecordName}}'s name and version.
	Name     string        // Name is the name of the {{.RecordName}}.
	Status   {{.RecordName}}Status // Status is the status of the {{.RecordName}}.
}

// {{.RecordName}}State is the state of a {{.RecordName}}.
type {{.RecordName}}State string

const (
	{{.RecordName}}StateUnknown  {{.RecordName}}State = "{{.RecordName}}State_UNKNOWN"
	{{.RecordName}}StatePending  {{.RecordName}}State = "{{.RecordName}}State_PENDING"
	{{.RecordName}}StateActive   {{.RecordName}}State = "{{.RecordName}}State_ACTIVE"
	{{.RecordName}}StateInActive {{.RecordName}}State = "{{.RecordName}}State_INACTIVE"
)

// ToString returns the string representation of the {{.RecordName}}State.
func (s {{.RecordName}}State) ToString() string {
	return string(s)
}

// FromString converts a string into a {{.RecordName}}State if valid, otherwise returns an error.
func {{.RecordName}}StateFromString(s string) {{.RecordName}}State {
	switch s {
	case string({{.RecordName}}StatePending):
		return {{.RecordName}}StatePending
	case string({{.RecordName}}StateActive):
		return {{.RecordName}}StateActive
	case string({{.RecordName}}StateInActive):
		return {{.RecordName}}StateInActive
	default:
		return {{.RecordName}}State(s) // Unknown state. Return as is.
	}
}

type {{.RecordName}}Status struct {
	State   {{.RecordName}}State // State is the discrete condition of the resource.
	Message string       // Message is a human-readable description of the resource's state.
}

// Ledger provides the methods for managing {{.RecordName}} records.
type Ledger interface {
	// Create creates a new {{.RecordName}}.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// GetByMetadata retrieves a {{.RecordName}} by its metadata.
	GetByMetadata(context.Context, *core.Metadata) (*GetResponse, error)
	// GetByName retrieves a {{.RecordName}} by its name.
	GetByName(context.Context, string) (*GetResponse, error)
	// UpdateStatus updates the state and message of an existing {{.RecordName}}.
	UpdateStatus(context.Context, *UpdateStateRequest) (*UpdateResponse, error)
	// List returns a list of {{.RecordName}} that match the provided filters.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes a {{.RecordName}}.
	Delete(context.Context, *DeleteRequest) error
}

// CreateRequest represents the {{.RecordName}} creation request.
type CreateRequest struct {
	Name string
}

// CreateResponse represents the response after creating a new {{.RecordName}}.
type CreateResponse struct {
	Record {{.RecordName}}Record
}

// UpdateStateRequest represents the request to update the state and message of a {{.RecordName}}.
type UpdateStateRequest struct {
	Metadata core.Metadata
	Status   {{.RecordName}}Status
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
	Filters {{.RecordName}}ListFilters
}

// {{.RecordName}}Filters contains filters for querying the {{.RecordName}} table.
type {{.RecordName}}ListFilters struct {
	IDIn       []string // IN condition
	NameIn     []string // IN condition
	VersionGte *uint64  // Greater than or equal condition
	VersionLte *uint64  // Less than or equal condition
	VersionEq  *uint64  // Equal condition

	IncludeDeleted bool   // IncludeDeleted indicates whether to include soft-deleted records in the result.
	Limit          uint32 // Limit is the maximum number of records to return.

	StateIn    []{{.RecordName}}State
	StateNotIn []{{.RecordName}}State
}

// ListResponse represents the response to a list request.
type ListResponse struct {
	Records []{{.RecordName}}Record
}

type DeleteRequest struct {
	Metadata core.Metadata
}
`

const ledgerTemplate = `
package {{.PackageName}}

import (
	"context"
	"fmt"

	"{{.GoModuleName}}/internal/ledger/core"
	ledgererrors "{{.GoModuleName}}/internal/ledger/errors"

	"github.com/google/uuid"
)

// ledger implements the Ledger interface.
type ledger struct {
	repo Repository
}

// Repository provides the methods that the storage layer must implement to support the ledger.
type Repository interface {
	Insert(context.Context, {{.RecordName}}Record) error
	GetByMetadata(context.Context, core.Metadata) ({{.RecordName}}Record, error)
	GetByName(context.Context, string) ({{.RecordName}}Record, error)
	UpdateState(context.Context, core.Metadata, {{.RecordName}}Status) error
	Delete(context.Context, core.Metadata) error
	List(context.Context, {{.RecordName}}ListFilters) ([]{{.RecordName}}Record, error)
}

// NewLedger creates a new Ledger instance.
func NewLedger(repo Repository) Ledger {
	return &ledger{repo: repo}
}

// Create creates a new {{.RecordName}}.
func (l *ledger) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	// validate the request
	if req.Name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"{{.RecordName}} name is required",
		)
	}

	rec := {{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      uuid.New().String(),
			Version: 0,
		},
		Name: req.Name,
		Status: {{.RecordName}}Status{
			State:   {{.RecordName}}StatePending,
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
	// validate the request
	if metadata == nil || metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to fetch by metadata",
		)
	}

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
	// validate the request
	if name == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"Name missing. Name is required to fetch by name",
		)
	}

	record, err := l.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Record: record,
	}, nil
}

var validStateTransitions = map[{{.RecordName}}State][]{{.RecordName}}State{
	{{.RecordName}}StatePending:  {{"{"}}{{.RecordName}}StateActive, {{.RecordName}}StateInActive},
	{{.RecordName}}StateActive:   {{"{"}}{{.RecordName}}StateInActive},
	{{.RecordName}}StateInActive: {{"{"}}{{.RecordName}}StateActive},
}

// UpdateStatus updates the state and message of an existing {{.RecordName}}.
func (l *ledger) UpdateStatus(ctx context.Context, req *UpdateStateRequest) (*UpdateResponse, error) {
	// validate the request
	if req.Metadata.ID == "" {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			"ID missing. ID is required to update state",
		)
	}

	existingRecord, err := l.repo.GetByMetadata(ctx, req.Metadata)
	if err != nil {
		if ledgererrors.IsLedgerError(err) && ledgererrors.AsLedgerError(err).Code == ledgererrors.ErrRecordNotFound {
			return nil, ledgererrors.NewLedgerError(
				ledgererrors.ErrRecordInsertConflict,
				"Either record does not exist or version mismatch resulted in conflict. Check and retry.",
			)
		}
	}

	// validate the state transition
	validStates, ok := validStateTransitions[existingRecord.Status.State]
	if !ok {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			fmt.Sprintf("Invalid state transition from %s to %s", existingRecord.Status.State, req.Status.State),
		)
	}
	var valid bool
	for _, state := range validStates {
		if state == req.Status.State {
			valid = true
			break
		}
	}
	if !valid {
		return nil, ledgererrors.NewLedgerError(
			ledgererrors.ErrRequestInvalid,
			fmt.Sprintf("Invalid state transition from %s to %s", existingRecord.Status.State, req.Status.State),
		)
	}

	err = l.repo.UpdateState(ctx, req.Metadata, req.Status)
	if err != nil {
		return nil, err
	}

	record, err := l.repo.GetByMetadata(ctx, core.Metadata{
		ID:      req.Metadata.ID,
		Version: req.Metadata.Version + 1,
	})
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

func (l *ledger) Delete(ctx context.Context, req *DeleteRequest) error {
	return l.repo.Delete(ctx, req.Metadata)
}
`

const ledgerTestTemplate = `
package {{.PackageName}}_test

import (
	"context"
	"testing"

	"{{.GoModuleName}}/internal/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/internal/ledger/core"
	ledgererrors "{{.GoModuleName}}/internal/ledger/errors"
	"{{.GoModuleName}}/internal/sqlstorage/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerCreate(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

		req := &{{.PackageName}}.CreateRequest{
			Name: "test-{{.PackageName}}",
		}
		resp, err := l.Create(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-{{.PackageName}}", resp.Record.Name)
		require.NotEmpty(t, resp.Record.Metadata.ID)
		require.Equal(t, uint64(0), resp.Record.Metadata.Version)
		require.Equal(t, {{.PackageName}}.{{.RecordName}}StatePending, resp.Record.Status.State)
	})

	t.Run("Create EmptyName Failure", func(t *testing.T) {
		storage := test.TestSQLStorage(t)
		l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

		req := &{{.PackageName}}.CreateRequest{
			Name: "", // Empty name
		}
		resp, err := l.Create(context.Background(), req)

		require.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		require.Nil(t, resp)
	})
}

func TestLedgerGetByMetadata(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

	req := &{{.PackageName}}.CreateRequest{
		Name: "test-{{.PackageName}}",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByMetadata Success", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &createResp.Record.Metadata)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-{{.PackageName}}", resp.Record.Name)
	})

	t.Run("GetByMetadata InvalidID Failure", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &core.Metadata{ID: ""})

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})

	t.Run("GetByMetadata NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByMetadata(context.Background(), &core.Metadata{ID: "unknown"})

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})
}

func TestLedgerGetByName(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

	req := &{{.PackageName}}.CreateRequest{
		Name: "test-{{.PackageName}}",
	}
	_, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	t.Run("GetByName Success", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "test-{{.PackageName}}")

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-{{.PackageName}}", resp.Record.Name)
	})

	t.Run("GetByName InvalidName Failure", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "")

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})

	t.Run("GetByName NotFound Failure", func(t *testing.T) {
		resp, err := l.GetByName(context.Background(), "unknown")

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})
}

func TestLedgerUpdateStatus(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

	req := &{{.PackageName}}.CreateRequest{
		Name: "test-{{.PackageName}}",
	}
	createResp, err := l.Create(context.Background(), req)
	require.NoError(t, err)

	lastUpdatedRecord := createResp.Record
	t.Run("UpdateStatus Success", func(t *testing.T) {
		updateReq := &{{.PackageName}}.UpdateStateRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: {{.PackageName}}.{{.RecordName}}Status{
				State:   {{.PackageName}}.{{.RecordName}}StateInActive,
				Message: "{{.RecordName}} is inactive now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, {{.PackageName}}.{{.RecordName}}StateInActive, resp.Record.Status.State)
		lastUpdatedRecord = resp.Record
	})

	t.Run("UpdateStatus InvalidTransition Failure", func(t *testing.T) {
		updateReq := &{{.PackageName}}.UpdateStateRequest{
			Metadata: lastUpdatedRecord.Metadata,
			Status: {{.PackageName}}.{{.RecordName}}Status{
				State:   {{.PackageName}}.{{.RecordName}}StatePending, // Invalid transition
				Message: "Invalid state transition",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRequestInvalid, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})

	t.Run("Update conflict Failure", func(t *testing.T) {
		updateReq := &{{.PackageName}}.UpdateStateRequest{
			Metadata: createResp.Record.Metadata, // This is the old metadata. Should cause a conflict.
			Status: {{.PackageName}}.{{.RecordName}}Status{
				State:   {{.PackageName}}.{{.RecordName}}StateActive,
				Message: "{{.RecordName}} is active now",
			},
		}

		resp, err := l.UpdateStatus(context.Background(), updateReq)

		assert.Error(t, err)
		require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
		require.Equal(t, ledgererrors.ErrRecordInsertConflict, err.(ledgererrors.LedgerError).Code)
		assert.Nil(t, resp)
	})
}

func TestLedgerList(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

	// Create two {{.RecordName}}s
	_, err := l.Create(context.Background(), &{{.PackageName}}.CreateRequest{Name: "{{.RecordName}}1"})
	assert.NoError(t, err)

	_, err = l.Create(context.Background(), &{{.PackageName}}.CreateRequest{Name: "{{.RecordName}}2"})
	assert.NoError(t, err)

	// List the {{.RecordName}}s
	listReq := &{{.PackageName}}.ListRequest{}
	resp, err := l.List(context.Background(), listReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Records, 2)
}

func TestLedgerDelete(t *testing.T) {
	storage := test.TestSQLStorage(t)
	l := {{.PackageName}}.NewLedger(storage.{{.RecordName}})

	// First, create the {{.RecordName}}
	createResp, err := l.Create(context.Background(), &{{.PackageName}}.CreateRequest{Name: "test-{{.PackageName}}"})
	assert.NoError(t, err)

	// Now, delete the {{.RecordName}}
	err = l.Delete(context.Background(), &{{.PackageName}}.DeleteRequest{Metadata: createResp.Record.Metadata})
	assert.NoError(t, err)

	// Try to get the {{.RecordName}} again
	_, err = l.GetByMetadata(context.Background(), &createResp.Record.Metadata)
	assert.Error(t, err)
	require.ErrorAs(t, err, &ledgererrors.LedgerError{}, "error should be of type LedgerError")
	require.Equal(t, ledgererrors.ErrRecordNotFound, err.(ledgererrors.LedgerError).Code)
}
`

func (o GenerateOptions) generateLedgerRecord() error {
	ledgerSubPath := filepath.Join(o.DestinationPath, "internal", "ledger", o.PackageName)

	fmt.Println("Generating Ledger at ", ledgerSubPath)
	// Create destination path if it doesn't exist
	err := os.MkdirAll(ledgerSubPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination path: %w", err)
	}

	fmt.Println("... creating interface.go")
	err = executeTemplate("interfaceTemplate", interfaceTemplate, ledgerSubPath, "interface.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate interface file: %w", err)
	}

	fmt.Println("... creating ledger.go")
	err = executeTemplate("ledgerTemplate", ledgerTemplate, ledgerSubPath, "ledger.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate ledger file: %w", err)
	}

	fmt.Println("... creating ledger_test.go")
	err = executeTemplate("ledgerTest", ledgerTestTemplate, ledgerSubPath, "ledger_test.go", o)
	if err != nil {
		return fmt.Errorf("failed to generate ledger file: %w", err)
	}

	return nil
}
