package main

import (
	"context"

	"github.com/spf13/cobra"
)

type ORMGenOptions struct {
	StructName string
	PkgName    string
	// OutputFile string
	TableName string
}

func (o ORMGenOptions) Run(ctx context.Context) error {
	generator, err := newGenerator(o)
	if err != nil {
		return err
	}
	return generator.Generate()
}

func main() {
	o := ORMGenOptions{}

	cmd := cobra.Command{
		Use:   "simplesqlorm-gen",
		Short: "Generate ORM code for a given struct",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(context.Background())
		},
	}

	cmd.Flags().StringVar(&o.StructName, "struct-name", "", "Name of the struct to generate ORM code for")
	cmd.MarkFlagRequired("struct-name")

	cmd.Flags().StringVar(&o.PkgName, "pkg-name", "", "Name of the package the struct belongs to")
	cmd.MarkFlagRequired("pkg-name")

	// cmd.Flags().StringVar(&o.OutputFile, "output-file", "", "Output file to write the generated code to")
	// cmd.MarkFlagRequired("output-file")

	cmd.Flags().StringVar(&o.TableName, "table-name", "", "Name of the table in the database")
	cmd.MarkFlagRequired("table-name")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
