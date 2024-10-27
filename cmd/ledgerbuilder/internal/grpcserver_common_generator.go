package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const testServerFileTemplate = `
package test

import (
	"context"
	"fmt"
	"net"

	"{{.GoModuleName}}/gen/api/mrdspb"
	"{{.GoModuleName}}/internal/ledger/{{.PackageName}}"
	"{{.GoModuleName}}/internal/sqlstorage"
	"{{.GoModuleName}}/pkg/grpcservers"

	// ++ledgerbuilder:Imports

	"github.com/msanath/gondolf/pkg/simplesql/test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type TestServer struct {
	server *grpc.Server
	conn   *grpc.ClientConn
}

var testDb = test.NewTestSQLiteDB

// var testDb = test.NewTestMySQLDB

func NewTestServer() (*TestServer, error) {
	gServer := grpc.NewServer()

	db, err := testDb()
	if err != nil {
		return nil, fmt.Errorf("failed to create test sqlite db: %w", err)
	}
	storage, err := sqlstorage.NewSQLStorage(db, false)
	if err != nil {
		return nil, err
	}

	{{.AttributePrefix}}Ledger := {{.PackageName}}.NewLedger(storage.{{.RecordName}})
	{{.ProtoPkgNamespace}}pb.Register{{.RecordName}}sServer(
		gServer,
		grpcservers.New{{.RecordName}}Service({{.AttributePrefix}}Ledger),
	)
	// ++ledgerbuilder:TestServerRegister

	listener := bufconn.Listen(1024 * 1024)
	go func() {
		if err := gServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	//nolint:staticcheck
	conn, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client connection: %w", err)
	}

	return &TestServer{
		server: gServer,
		conn:   conn,
	}, nil
}

func (s *TestServer) Conn() *grpc.ClientConn {
	return s.conn
}

func (s *TestServer) Close() {
	s.conn.Close()
	s.server.Stop()
}
`

func (o GenerateOptions) generateGRPCServerCommon() error {
	fmt.Println("Generating storage components")

	testServerPath := filepath.Join(o.DestinationPath, "pkg", "grpcservers", "test")
	err := os.MkdirAll(testServerPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create test Server path: %w", err)
	}

	testFilePath := filepath.Join(testServerPath, "new.go")
	// If file is present, update at ++ledgerbuilder markers
	if _, err := os.Stat(testFilePath); err == nil {
		fmt.Println("... existing new.go found. Updating new.go")
		err = o.addToTestServerPath(testFilePath)
		if err != nil {
			return fmt.Errorf("failed to update new storage file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	} else {
		fmt.Println("... creating new.go")
		err = executeTemplate("testServerFileTemplate", testServerFileTemplate, testServerPath, "new.go", o)
		if err != nil {
			return fmt.Errorf("failed to generate new test server file: %w", err)
		}
	}

	return nil
}

// updateNewFile updates the existing 'new.go' file to add new repositories and migrations.
func (o GenerateOptions) addToTestServerPath(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Prepare the new lines to be added
	newServer := fmt.Sprintf(`
	%sLedger := %s.NewLedger(storage.%s)
	%spb.Register%ssServer(
		gServer,
		grpcservers.New%sService(%sLedger),
	)`, o.AttributePrefix, o.PackageName, o.RecordName, o.ProtoPkgNamespace, o.RecordName, o.RecordName, o.AttributePrefix)

	newImport := fmt.Sprintf("\"%s/internal/ledger/%s\"", o.GoModuleName, o.PackageName)

	// Insert into the appropriate places
	updatedContent := insertAtPlaceholder(string(content), "// ++ledgerbuilder:TestServerRegister", newServer)
	updatedContent = insertAtPlaceholder(updatedContent, "// ++ledgerbuilder:Imports", newImport)

	// Write the updated content back to the file
	return os.WriteFile(filePath, []byte(updatedContent), 0644)
}
