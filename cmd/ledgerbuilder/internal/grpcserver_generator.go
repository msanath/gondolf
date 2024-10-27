package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const grpcServerTemplate = `
package grpcservers

import (
	"context"

	"{{.GoModuleName}}/gen/api/{{.ProtoPkgNamespace}}pb"
	"{{.GoModuleName}}/internal/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/internal/ledger/core"
)

type {{.RecordName}}Service struct {
	ledger              {{.PackageName}}.Ledger
	protoToLedgerRecord func(proto *{{.ProtoPkgNamespace}}pb.{{.RecordName}}) {{.PackageName}}.{{.RecordName}}Record
	ledgerRecordToProto func(record {{.PackageName}}.{{.RecordName}}Record) *{{.ProtoPkgNamespace}}pb.{{.RecordName}}

	{{.ProtoPkgNamespace}}pb.Unimplemented{{.RecordName}}sServer
}

func {{.AttributePrefix}}ProtoToLedgerRecord(proto *{{.ProtoPkgNamespace}}pb.{{.RecordName}}) {{.PackageName}}.{{.RecordName}}Record {
	return {{.PackageName}}.{{.RecordName}}Record{
		Metadata: core.Metadata{
			ID:      proto.Metadata.Id,
			Version: proto.Metadata.Version,
		},
		Name: proto.Name,
		Status: {{.PackageName}}.{{.RecordName}}Status{
			State:   {{.PackageName}}.{{.RecordName}}State(proto.Status.State.String()),
			Message: proto.Status.Message,
		},
	}
}

func {{.AttributePrefix}}LedgerRecordToProto(record {{.PackageName}}.{{.RecordName}}Record) *{{.ProtoPkgNamespace}}pb.{{.RecordName}} {
	return &{{.ProtoPkgNamespace}}pb.{{.RecordName}}{
		Metadata: &{{.ProtoPkgNamespace}}pb.Metadata{
			Id:      record.Metadata.ID,
			Version: record.Metadata.Version,
		},
		Name: record.Name,
		Status: &{{.ProtoPkgNamespace}}pb.{{.RecordName}}Status{
			State:   {{.ProtoPkgNamespace}}pb.{{.RecordName}}State({{.ProtoPkgNamespace}}pb.{{.RecordName}}State_value[record.Status.State.ToString()]),
			Message: record.Status.Message,
		},
	}
}

func New{{.RecordName}}Service(ledger {{.PackageName}}.Ledger) *{{.RecordName}}Service {
	return &{{.RecordName}}Service{
		ledger:              ledger,
		protoToLedgerRecord: {{.AttributePrefix}}ProtoToLedgerRecord,
		ledgerRecordToProto: {{.AttributePrefix}}LedgerRecordToProto,
	}
}

// Create creates a new {{.RecordName}}
func (s *{{.RecordName}}Service) Create(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Request) (*{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Response, error) {
	createResponse, err := s.ledger.Create(ctx, &{{.PackageName}}.CreateRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Response{Record: s.ledgerRecordToProto(createResponse.Record)}, nil
}

// GetByMetadata retrieves a {{.RecordName}} by its metadata
func (s *{{.RecordName}}Service) GetByMetadata(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByMetadataRequest) (*{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}Response, error) {
	getResponse, err := s.ledger.GetByMetadata(ctx, &core.Metadata{
		ID:      req.Metadata.Id,
		Version: req.Metadata.Version,
	})
	if err != nil {
		return nil, err
	}
	return &{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}Response{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// GetByName retrieves a {{.RecordName}} by its name
func (s *{{.RecordName}}Service) GetByName(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByNameRequest) (*{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}Response, error) {
	getResponse, err := s.ledger.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}Response{Record: s.ledgerRecordToProto(getResponse.Record)}, nil
}

// UpdateStatus updates the state and message of an existing {{.RecordName}}
func (s *{{.RecordName}}Service) UpdateStatus(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Update{{.RecordName}}StatusRequest) (*{{.ProtoPkgNamespace}}pb.Update{{.RecordName}}Response, error) {
	updateResponse, err := s.ledger.UpdateStatus(ctx, &{{.PackageName}}.UpdateStatusRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
		Status: {{.PackageName}}.{{.RecordName}}Status{
			State:   {{.PackageName}}.{{.RecordName}}State(req.Status.State.String()),
			Message: req.Status.Message,
		},
	})
	if err != nil {
		return nil, err
	}
	return &{{.ProtoPkgNamespace}}pb.Update{{.RecordName}}Response{Record: s.ledgerRecordToProto(updateResponse.Record)}, nil
}

// List returns a list of {{.RecordName}}s that match the provided filters
func (s *{{.RecordName}}Service) List(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Request) (*{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Response, error) {
	if req == nil {
		req = &{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Request{}
	}
	var gte, lte, eq *uint64
	if req.VersionGte != 0 {
		gte = &req.VersionGte
	}
	if req.VersionLte != 0 {
		lte = &req.VersionLte
	}
	if req.VersionEq != 0 {
		eq = &req.VersionEq
	}

	stateIn := make([]{{.PackageName}}.{{.RecordName}}State, len(req.StateIn))
	for i, state := range req.StateIn {
		stateIn[i] = {{.PackageName}}.{{.RecordName}}StateFromString(state.String())
	}

	stateNotIn := make([]{{.PackageName}}.{{.RecordName}}State, len(req.StateNotIn))
	for i, state := range req.StateNotIn {
		stateNotIn[i] = {{.PackageName}}.{{.RecordName}}StateFromString(state.String())
	}

	listResponse, err := s.ledger.List(ctx, &{{.PackageName}}.ListRequest{
		Filters: {{.PackageName}}.{{.RecordName}}ListFilters{
			IDIn:       req.IdIn,
			NameIn:     req.NameIn,
			VersionGte: gte,
			VersionLte: lte,
			VersionEq:  eq,
			StateIn:    stateIn,
			StateNotIn: stateNotIn,
		},
	})
	if err != nil {
		return nil, err
	}

	records := make([]*{{.ProtoPkgNamespace}}pb.{{.RecordName}}, len(listResponse.Records))
	for i, record := range listResponse.Records {
		records[i] = s.ledgerRecordToProto(record)
	}

	return &{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Response{Records: records}, nil
}

// Delete deletes a {{.RecordName}}
func (s *{{.RecordName}}Service) Delete(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Delete{{.RecordName}}Request) (*{{.ProtoPkgNamespace}}pb.Delete{{.RecordName}}Response, error) {
	err := s.ledger.Delete(ctx, &{{.PackageName}}.DeleteRequest{
		Metadata: core.Metadata{
			ID:      req.Metadata.Id,
			Version: req.Metadata.Version,
		},
	})
	if err != nil {
		return nil, err
	}
	return &{{.ProtoPkgNamespace}}pb.Delete{{.RecordName}}Response{}, nil
}
`

const grpcServerTestTemplate = `
package grpcservers_test

import (
	"context"
	"testing"

	"{{.GoModuleName}}/gen/api/{{.ProtoPkgNamespace}}pb"
	servertest "{{.GoModuleName}}/pkg/grpcservers/test"

	"github.com/stretchr/testify/require"
)

func Test{{.RecordName}}Server(t *testing.T) {
	ts, err := servertest.NewTestServer()
	require.NoError(t, err)
	defer ts.Close()

	client := {{.ProtoPkgNamespace}}pb.New{{.RecordName}}sClient(ts.Conn())
	ctx := context.Background()

	// create
	resp, err := client.Create(ctx, &{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Request{Name: "test-{{.RecordName}}"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "test-{{.RecordName}}", resp.Record.Name)

	// get by metadata
	getResp, err := client.GetByMetadata(ctx, &{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByMetadataRequest{
		Metadata: resp.Record.Metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.Equal(t, "test-{{.RecordName}}", getResp.Record.Name)

	// get by name
	getByNameResp, err := client.GetByName(ctx, &{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByNameRequest{Name: "test-{{.RecordName}}"})
	require.NoError(t, err)
	require.NotNil(t, getByNameResp)
	require.Equal(t, "test-{{.RecordName}}", getByNameResp.Record.Name)

	// update
	updateResp, err := client.UpdateStatus(ctx, &{{.ProtoPkgNamespace}}pb.Update{{.RecordName}}StatusRequest{
		Metadata: resp.Record.Metadata,
		Status: &{{.ProtoPkgNamespace}}pb.{{.RecordName}}Status{
			State:   {{.ProtoPkgNamespace}}pb.{{.RecordName}}State_{{.RecordName}}State_ACTIVE,
			Message: "test-message",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.Equal(t, "test-{{.RecordName}}", updateResp.Record.Name)
	require.Equal(t, {{.ProtoPkgNamespace}}pb.{{.RecordName}}State_{{.RecordName}}State_ACTIVE, updateResp.Record.Status.State)

	// Create another
	resp2, err := client.Create(ctx, &{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Request{Name: "test-{{.RecordName}}-2"})
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// list
	listResp, err := client.List(ctx, &{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Request{
		StateIn: []{{.ProtoPkgNamespace}}pb.{{.RecordName}}State{{"{"}}{{.ProtoPkgNamespace}}pb.{{.RecordName}}State_{{.RecordName}}State_ACTIVE},
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Records, 1)

	// Delete
	_, err = client.Delete(ctx, &{{.ProtoPkgNamespace}}pb.Delete{{.RecordName}}Request{Metadata: resp2.Record.Metadata})
	require.NoError(t, err)

	// Get deleted by name
	_, err = client.GetByName(ctx, &{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByNameRequest{Name: "test-{{.RecordName}}-2"})
	require.Error(t, err)
}
`

func (o GenerateOptions) generateGRPCServer() error {
	grpcServersPath := filepath.Join(o.DestinationPath, "pkg", "grpcservers")
	// Create destination path if it doesn't exist
	err := os.MkdirAll(grpcServersPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create grpcservers path: %w", err)
	}

	fileName := fmt.Sprintf("%s.go", o.PackageName)
	fmt.Println("... creating ", fileName)
	err = executeTemplate("grpcServerTemplate", grpcServerTemplate, grpcServersPath, fileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	testFileName := fmt.Sprintf("%s_test.go", o.PackageName)
	fmt.Println("... creating ", testFileName)
	err = executeTemplate("grpcServerTestTemplate", grpcServerTestTemplate, grpcServersPath, testFileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	return nil
}
