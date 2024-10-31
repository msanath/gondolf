package internal

import (
	"errors"
	"fmt"
	"go/types"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type Generator struct {
	structType *types.Struct
	structName string
	pkgName    string
	outputFile string
}

func NewGenerator(structName string, pkgName string, outputFile string) (*Generator, error) {
	pkg, err := getCurrentPackage()
	if err != nil {
		return nil, err
	}
	obj := pkg.Types.Scope().Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("type '%s' not found", structName)
	}

	s, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("type '%s' must be a struct", structName)
	}
	err = isValidStructForDisplayGen(s)
	if err != nil {
		return nil, err
	}

	return &Generator{
		pkgName:    pkgName,
		structName: structName,
		outputFile: outputFile,
		structType: s,
	}, nil
}

// isValidStructForDisplayGen checks if the struct is valid for ORM generation.
func isValidStructForDisplayGen(s *types.Struct) error {
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		if f.Embedded() {
			// TODO: handle embedded struct fields
			continue
		}

		if !f.Exported() {
			// Skip private fields
			continue
		}

		name := f.Name()
		tags := reflect.StructTag(s.Tag(i))
		_, ok := tags.Lookup("json")
		if !ok {
			return fmt.Errorf("all fields must have json tag. field %s has no `db` tag name", name)
		}
	}
	return nil
}

func (g *Generator) Generate() error {
	f, err := os.Create(g.outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write the struct body
	bodyData, err := g.getBody()
	if err != nil {
		return err
	}
	bt := template.Must(template.New("body").Parse(bodyTemplate))
	err = bt.Execute(f, bodyData)
	if err != nil {
		return err
	}

	return g.tidyGeneratedFile()
}

func (g *Generator) tidyGeneratedFile() error {
	// Run "goimports"
	importsCmd := exec.Command("goimports", "-w", g.outputFile)
	importsCmd.Stdout = os.Stdout
	importsCmd.Stderr = os.Stderr
	if err := importsCmd.Run(); err != nil {
		fmt.Println("Error running 'goimports':", err)
		return err
	}

	// Run "go mod tidy"
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		fmt.Println("Error running 'go mod tidy':", err)
		return err
	}

	// Run "go fmt"
	fmtCmd := exec.Command("go", "fmt")
	fmtCmd.Stdout = os.Stdout
	fmtCmd.Stderr = os.Stderr
	if err := fmtCmd.Run(); err != nil {
		fmt.Println("Error running 'go fmt':", err)
		return err
	}
	return nil
}

func getCurrentPackage() (*packages.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedName | packages.NeedImports | packages.NeedDeps}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, err
	}

	// A directory can have both a regular package and "_test" package in it, and we want the regular one.
	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
		return pkg, nil
	}

	return nil, errors.New("no package found")
}
