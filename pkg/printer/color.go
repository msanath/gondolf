package printer

import "github.com/fatih/color"

type colorType int

// Following is an enum of colors
const (
	systemDefault colorType = iota
	cyan
	red
	yellow
	green
	magenta
	blue
)

// withColor provides a function which renders a string with the mentioned color.
func withColor(ct colorType) func(a ...interface{}) string {
	switch ct {
	case cyan:
		return color.New(color.FgCyan).SprintFunc()
	case red:
		return color.New(color.FgRed).SprintFunc()
	case yellow:
		return color.New(color.FgYellow).SprintFunc()
	case blue:
		return color.New(color.FgBlue).SprintFunc()
	case green:
		return color.New(color.FgGreen).SprintFunc()
	case magenta:
		return color.New(color.FgMagenta).SprintFunc()
	case systemDefault:
		fallthrough
	default:
		return nil
	}
}

func GreenText(text string) string {
	return withColor(green)(text)
}

func RedText(text string) string {
	return withColor(red)(text)
}

func YellowText(text string) string {
	return withColor(yellow)(text)
}

func CyanText(text string) string {
	return withColor(cyan)(text)
}

func BlueText(text string) string {
	return withColor(blue)(text)
}

func MagentaText(text string) string {
	return withColor(magenta)(text)
}
