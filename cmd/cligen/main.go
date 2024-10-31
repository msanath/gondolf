package main

import (
	"fmt"

	"os"

	"github.com/msanath/gondolf/cmd/cligen/internal"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 4 { // nolint:gomnd
		fmt.Println("Usage: cligen <GoStructType> <pkgName> <outputFile>")
		os.Exit(1)
	}
	structType := os.Args[1]
	pkgName := os.Args[2]
	outputFile := os.Args[3]

	// Get a writer to write to the output file.
	generator, err := internal.NewGenerator(structType, pkgName, outputFile)
	if err != nil {
		fmt.Println("Unable to initialize generator:", err)
		os.Exit(1)
	}
	err = generator.Generate()
	if err != nil {
		fmt.Println("Error generating statements:", err)
		os.Exit(1)
	}
}
