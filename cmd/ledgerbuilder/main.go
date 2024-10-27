package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/msanath/gondolf/cmd/ledgerbuilder/internal"
	"github.com/spf13/cobra"
)

func main() {
	o := internal.GenerateOptions{}

	cleanAndStart := false
	cmd := cobra.Command{
		Use: "ledger-builder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cleanAndStart {
				fmt.Println("Cleaning destination/core path")
				path := filepath.Join(o.DestinationPath, "internal")
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to clean destination path: %w", err)
				}

				path = filepath.Join(o.DestinationPath, "api")
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to clean destination path: %w", err)
				}

				path = filepath.Join(o.DestinationPath, "pkg", "grpcservers")
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to clean destination path: %w", err)
				}

				path = filepath.Join(o.DestinationPath, "pkg", "controlplane", "temporal", "activities")
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to clean destination path: %w", err)
				}
			}

			return o.Generate(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.RecordName, "record-name", "", "Name of the record to generate")
	cmd.MarkFlagRequired("record-name")

	cmd.Flags().StringVar(&o.DestinationPath, "destination-path", "", "Path to save the generated record")
	cmd.MarkFlagRequired("destination-path")

	cmd.Flags().StringVar(&o.GoModuleName, "go-module-name", "", "Name of the Go module to use")
	cmd.MarkFlagRequired("go-module-name")

	cmd.Flags().StringVar(&o.PackageName, "package-name", "", "Name of the package to use")
	cmd.Flags().StringVar(&o.TableName, "table-name", "", "Name of the table to use")

	cmd.Flags().BoolVar(&cleanAndStart, "clean", false, "Clean the core path before generating the record")
	cmd.Execute()
}
