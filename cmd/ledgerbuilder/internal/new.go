package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type GenerateOptions struct {
	PackageName       string
	GoModuleName      string
	RecordName        string
	DestinationPath   string
	TableName         string
	AttributePrefix   string
	ProtoPkgNamespace string

	TablesOnly bool
}

func NewGenerator(opts GenerateOptions) *GenerateOptions {
	return &opts
}

func (o GenerateOptions) Generate(ctx context.Context) error {
	fmt.Printf("Generating record %s at %s using Go module %s\n", o.RecordName, o.DestinationPath, o.GoModuleName)

	if o.PackageName == "" {
		o.PackageName = strings.ToLower(o.RecordName)
	}

	// TableName makes the name lowercase and adds an underscore between words
	if o.TableName == "" {
		recordName := o.RecordName

		// Find the first uppercase letter after the first letter
		for i, r := range recordName[1:] {
			if r >= 'A' && r <= 'Z' {
				o.TableName = strings.ToLower(o.AttributePrefix[:i+1]) + "_" + strings.ToLower(o.AttributePrefix[i+1:])
				break
			}
		}
		o.TableName = strings.ToLower(recordName)
	}

	// protoPkgNamespace makes the last part of the Go module name lowercase
	o.ProtoPkgNamespace = filepath.Base(o.GoModuleName)

	// recordName makes the first letter of the record name lowercase
	o.AttributePrefix = strings.ToLower(o.RecordName[:1]) + o.RecordName[1:]
	fmt.Println("----------------------------------------")
	fmt.Println("Package Name: ", o.PackageName)
	fmt.Println("Go Module Name: ", o.GoModuleName)
	fmt.Println("Record Name: ", o.RecordName)
	fmt.Println("Destination Path: ", o.DestinationPath)
	fmt.Println("Table Name: ", o.TableName)
	fmt.Println("----------------------------------------")

	if o.TablesOnly {
		if err := o.generateTables(); err != nil {
			return err
		}
		return nil
	}

	if err := o.generateCoreComponents(); err != nil {
		return err
	}

	if err := o.generateLedgerRecord(); err != nil {
		return err
	}

	if err := o.generateTables(); err != nil {
		return err
	}

	if err := o.generateStorageCommon(); err != nil {
		return err
	}

	if err := o.generateRepo(); err != nil {
		return err
	}

	if err := o.generateProto(); err != nil {
		return err
	}

	if err := o.generateGRPCServerCommon(); err != nil {
		return err
	}

	if err := o.generateGRPCServer(); err != nil {
		return err
	}

	if err := o.generateTemporalActivity(); err != nil {
		return err
	}

	fmt.Println("----------------------------------------")

	fmt.Println("Tidying up generated files")
	if err := o.tidyGeneratedFile(); err != nil {
		return err
	}
	return nil
}

func executeTemplate(templateName, templateStr, path, fileName string, data any) error {
	// Generate the records file
	filePath := filepath.Join(path, fileName)
	// Check if the file already exists. If it does, return an error.
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("record file already exists at %s", filePath)
	}

	// Create the record file
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create record file: %w", err)
	}
	defer f.Close()

	// Execute the template
	t := template.Must(template.New(templateName).Parse(templateStr))
	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute record template: %w", err)
	}
	return nil
}

func (g *GenerateOptions) tidyGeneratedFile() error {
	// // Run "goimports" on the generated files in the DestinationPath
	// importsCmd := exec.Command("goimports", "-w", g.DestinationPath)
	// importsCmd.Stdout = os.Stdout
	// importsCmd.Stderr = os.Stderr
	// if err := importsCmd.Run(); err != nil {
	// 	fmt.Println("Error running 'goimports':", err)
	// 	return err
	// }
	mkProto := exec.Command("make", "proto")
	mkProto.Dir = g.DestinationPath // Set the working directory to the destination path
	mkProto.Stdout = os.Stdout
	mkProto.Stderr = os.Stderr
	if err := mkProto.Run(); err != nil {
		fmt.Println("Error running 'make proto':", err)
		return err
	}

	// Run "go mod tidy" in the DestinationPath
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = g.DestinationPath // Set the working directory to the destination path
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		fmt.Println("Error running 'go mod tidy':", err)
		return err
	}

	// Run "go fmt" on the generated files in the DestinationPath
	fmtCmd := exec.Command("go", "fmt", "./...")
	fmtCmd.Dir = g.DestinationPath // Set the working directory to the destination path
	fmtCmd.Stdout = os.Stdout
	fmtCmd.Stderr = os.Stderr
	if err := fmtCmd.Run(); err != nil {
		fmt.Println("Error running 'go fmt':", err)
		return err
	}

	return nil
}
