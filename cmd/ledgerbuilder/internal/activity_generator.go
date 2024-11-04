package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const temporalActivityTemplate = `
package {{.ProtoPkgNamespace}}activities

import (
	"context"
	"fmt"

	"{{.GoModuleName}}/gen/api/{{.ProtoPkgNamespace}}pb"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type {{.RecordName}}Activities struct {
	client {{.ProtoPkgNamespace}}pb.{{.RecordName}}sClient
}

// New{{.RecordName}}Activities creates a new instance of {{.RecordName}}Activities.
func New{{.RecordName}}Activities(client {{.ProtoPkgNamespace}}pb.{{.RecordName}}sClient, registry worker.Registry) *{{.RecordName}}Activities {
	a := &{{.RecordName}}Activities{client: client}
	registry.RegisterActivity(a.Create{{.RecordName}})
	registry.RegisterActivity(a.Get{{.RecordName}}ByMetadata)
	registry.RegisterActivity(a.Get{{.RecordName}}ByName)
	registry.RegisterActivity(a.Update{{.RecordName}}Status)
	registry.RegisterActivity(a.List{{.RecordName}})
	registry.RegisterActivity(a.Delete{{.RecordName}})
	return a
}

// Create{{.RecordName}} is an activity that interacts with the gRPC service to create a {{.RecordName}}.
func (c *{{.RecordName}}Activities) Create{{.RecordName}}(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Request) (*{{.ProtoPkgNamespace}}pb.Create{{.RecordName}}Response, error) {
	activity.GetLogger(ctx).Info("Creating {{.RecordName}}", "request", req)

	// Check if the context has a deadline to handle timeout.
	if deadline, ok := ctx.Deadline(); ok {
		activity.GetLogger(ctx).Info("Context has a deadline: %v", deadline)
	}

	// Call gRPC method with context for timeout
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to create {{.RecordName}}", "error", err)
		return nil, fmt.Errorf("failed to create {{.RecordName}}: %w", err)
	}

	return resp, nil
}

// Get{{.RecordName}}ByMetadata fetches {{.RecordName}} details based on metadata.
func (c *{{.RecordName}}Activities) Get{{.RecordName}}ByMetadata(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByMetadataRequest) (*{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}Response, error) {
	activity.GetLogger(ctx).Info("Fetching {{.RecordName}} by metadata", "request", req)

	resp, err := c.client.GetByMetadata(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get {{.RecordName}} by metadata", "error", err)
		return nil, fmt.Errorf("failed to get {{.RecordName}} by metadata: %w", err)
	}

	return resp, nil
}

// Get{{.RecordName}}ByName fetches {{.RecordName}} details based on name.
func (c *{{.RecordName}}Activities) Get{{.RecordName}}ByName(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}ByNameRequest) (*{{.ProtoPkgNamespace}}pb.Get{{.RecordName}}Response, error) {
	activity.GetLogger(ctx).Info("Fetching {{.RecordName}} by name", "request", req)

	resp, err := c.client.GetByName(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to get {{.RecordName}} by name", "error", err)
		return nil, fmt.Errorf("failed to get {{.RecordName}} by name: %w", err)
	}

	return resp, nil
}

// Update{{.RecordName}}State updates the state of a {{.RecordName}}.
func (c *{{.RecordName}}Activities) Update{{.RecordName}}Status(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Update{{.RecordName}}StatusRequest) (*{{.ProtoPkgNamespace}}pb.Update{{.RecordName}}Response, error) {
	activity.GetLogger(ctx).Info("Updating {{.RecordName}} state", "request", req)

	resp, err := c.client.UpdateStatus(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to update {{.RecordName}} state", "error", err)
		return nil, fmt.Errorf("failed to update {{.RecordName}} state: %w", err)
	}

	return resp, nil
}

// List{{.RecordName}} lists all {{.RecordName}}s based on the request.
func (c *{{.RecordName}}Activities) List{{.RecordName}}(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Request) (*{{.ProtoPkgNamespace}}pb.List{{.RecordName}}Response, error) {
	activity.GetLogger(ctx).Info("Listing {{.RecordName}}s", "request", req)

	resp, err := c.client.List(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to list {{.RecordName}}s", "error", err)
		return nil, fmt.Errorf("failed to list {{.RecordName}}s: %w", err)
	}

	return resp, nil
}

func (c *{{.RecordName}}Activities) Delete{{.RecordName}}(ctx context.Context, req *{{.ProtoPkgNamespace}}pb.Delete{{.RecordName}}Request) (*{{.ProtoPkgNamespace}}pb.Delete{{.RecordName}}Response, error) {
	activity.GetLogger(ctx).Info("Deleting {{.RecordName}}", "request", req)

	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		activity.GetLogger(ctx).Error("Failed to delete {{.RecordName}}", "error", err)
		return nil, fmt.Errorf("failed to delete {{.RecordName}}: %w", err)
	}

	return resp, nil
}
`

func (o GenerateOptions) generateTemporalActivity() error {
	fmt.Println("Generating Temporal activities")

	activityPath := filepath.Join(o.DestinationPath, "controlplane", "temporal", "activities", fmt.Sprintf("%sactivities", o.ProtoPkgNamespace))
	err := os.MkdirAll(activityPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create activities path: %w", err)
	}

	activityFileName := fmt.Sprintf("%s.go", o.PackageName)
	fmt.Println("... creating ", activityFileName)
	err = executeTemplate("temporalActivityTemplate", temporalActivityTemplate, activityPath, activityFileName, o)
	if err != nil {
		return fmt.Errorf("failed to generate record file: %w", err)
	}

	return nil
}
