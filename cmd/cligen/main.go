package main

import (
	"context"
	"fmt"

	"os"

	"github.com/msanath/gondolf/cmd/cligen/internal"
	"github.com/spf13/cobra"
)

type cliGenOptions struct {
	StructName string
	PkgName    string
	OutputFile string
}

func main() {
	fmt.Println(os.Args)

	o := cliGenOptions{}

	cmd := cobra.Command{
		Use:   "cligen",
		Short: "Generate display helpers for a given struct",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(context.Background())
		},
	}

	cmd.Flags().StringVar(&o.StructName, "struct-name", "", "Name of the struct to generate ORM code for")
	cmd.MarkFlagRequired("struct-name")

	cmd.Flags().StringVar(&o.PkgName, "pkg-name", "", "Name of the package")
	cmd.MarkFlagRequired("pkg-name")

	cmd.Flags().StringVar(&o.OutputFile, "output-file", "", "Name of the output file")
	cmd.MarkFlagRequired("output-file")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func (o cliGenOptions) Run(ctx context.Context) error {
	// Get a writer to write to the output file.
	generator, err := internal.NewGenerator(o.StructName, o.PkgName, o.OutputFile)
	if err != nil {
		fmt.Println("Unable to initialize generator:", err)
		os.Exit(1)
	}
	err = generator.Generate()
	if err != nil {
		fmt.Println("Error generating statements:", err)
		os.Exit(1)
	}
	return nil
}
