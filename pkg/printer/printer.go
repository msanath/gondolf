package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// PlainText is a writer to display compute broker service entities in plain text format.
type PlainText interface {
	PrintTable(headers []string, rows [][]string, opts ...TablePrinterOption)
	PrintDisplayField(field DisplayField)
	PrintDisplayFieldWithIndent(field DisplayField)
	PrintKeyValue(key, value string)
	PrintKeyValueWithIndent(key, value string)
	PrintLineSeparator()
	PrintEmptyLine()
	PrintHeader(value string)
	PrintError(message string)
	PrintWarning(value string)
	PrintSuccess(value string)
	SeekConfirmation(message string) bool
	PrintInJSONFormat(obj interface{}) error
}
type plainText struct{}

func NewPlainTextPrinter() PlainText {
	return &plainText{}
}

type tablePrinterOptions struct {
	withRowSeparator bool
}

type TablePrinterOption func(*tablePrinterOptions)

func WithRowSeparator() TablePrinterOption {
	return func(o *tablePrinterOptions) {
		o.withRowSeparator = true
	}
}

func (p *plainText) PrintTable(headers []string, rows [][]string, opts ...TablePrinterOption) {
	options := tablePrinterOptions{
		withRowSeparator: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	table := tablewriter.NewWriter(os.Stdout)

	if options.withRowSeparator {
		table.SetRowLine(true)
	}
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(headers)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(rows)
	table.Render()
	p.PrintSuccess(fmt.Sprintf("Total Records: %d", len(rows)))
}

// PrintDisplayField prints the key and value in a key-value format. Ex:
func (p *plainText) PrintDisplayField(field DisplayField) {
	key := field.DisplayName
	value := field.Value()
	p.PrintKeyValue(key, value)
}

func (p *plainText) PrintKeyValue(key, value string) {
	if value == "" {
		value = "<unset>"
	}
	// if the value contains newlines, print it in a new line with an indent
	if strings.Contains(value, "\n") {
		// Split the value by newlines and print each line with an indent
		lines := strings.Split(value, "\n")
		fmt.Printf("%-35s%s\n", CyanText(key+":"), lines[0])
		for _, line := range lines[1:] {
			fmt.Printf("%-25s%s\n", "", line)
		}
		return
	}

	fmt.Printf("%-35s%s\n", CyanText(key+":"), value)
}

// PrintDisplayFieldWithIndent prints the key and value in a key-value format with an indent. Ex:
func (p *plainText) PrintDisplayFieldWithIndent(field DisplayField) {
	key := field.DisplayName
	value := field.Value()
	p.PrintKeyValueWithIndent(key, value)
}

func (p *plainText) PrintKeyValueWithIndent(key, value string) {
	if value == "" {
		value = "<unset>"
	}
	// if the value contains newlines, print it in a new line with an indent
	if strings.Contains(value, "\n") {
		// Split the value by newlines and print each line with an indent
		lines := strings.Split(value, "\n")
		fmt.Printf("\t%-35s%s\n", CyanText(key+":"), lines[0])
		for _, line := range lines[1:] {
			fmt.Printf("\t%-26s%s\n", "", line)
		}
		return
	}
	fmt.Printf("\t%-35s%s\n", CyanText(key+":"), value)
}

// prints the value in a header format. Ex:
// PrintHeader("Node ID")
// Output - Node ID
func (p *plainText) PrintHeader(value string) {
	fmt.Printf("%s\n", MagentaText(value))
}

// prints a warning with YELLOW color. Ex:
// PrintWarning("With great power comes great responsibility!")
// Output - With great power comes great responsibility!
func (p *plainText) PrintWarning(value string) {
	fmt.Printf("%s\n", withColor(yellow)(value))
}

// Error displays an error string with RED color.
func (p *plainText) PrintError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", RedText(message))
}

func (p *plainText) PrintLineSeparator() {
	fmt.Println(strings.Repeat("-", 80))
}

func (p *plainText) PrintEmptyLine() {
	fmt.Println()
}

// SeekConfirmation asks for user confirmation and returns true if the user confirms.
func (p *plainText) SeekConfirmation(message string) bool {
	fmt.Printf("%s [y/N]:: ", withColor(yellow)(message))
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "yes" || strings.ToLower(response) == "y"
}

// PrintSuccess prints the value in GREEN color. Ex:
// PrintSuccess("Completed successfully")
// Output - Completed successfully
func (p *plainText) PrintSuccess(value string) {
	fmt.Printf("%s\n", GreenText(value))
}

func (p *plainText) PrintInJSONFormat(obj interface{}) error {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}
