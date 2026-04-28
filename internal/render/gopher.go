package render

import (
	"fmt"
	"io"
)

// ShouldRenderGopher returns true when input is non-empty, contains only
// non-ASCII runes (after allowing \n \t \\), and every parsed segment is blank.
func ShouldRenderGopher(text string, segments []string) bool {
	if text == "" {
		return false
	}
	if !containsNonASCII(text) {
		return false
	}
	for _, seg := range segments {
		if seg != "" {
			return false
		}
	}
	return true
}

func containsNonASCII(text string) bool {
	// TODO(task11): skip \n \t \\; flag rune <32 or >126 as non-ASCII
	for _, r := range text {
		if r > 126 {
			return true
		}
	}
	return false
}

const (
	gopherRed    = "\033[31m"
	gopherGreen  = "\033[32m"
	gopherCyan   = "\033[36m"
	gopherOrange = "\033[38;5;208m"
	gopherReset  = "\033[0m"
)

// RenderGopher writes the hardcoded Peeking Gopher ANSI art to w.
func RenderGopher(writer io.Writer) {
	fmt.Fprint(writer, gopherGreen+"                   \n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"                   \n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"                   \n"+gopherReset)
	fmt.Fprint(writer, gopherRed+"🐸 Ooops! Peeking Gopher appeared...\n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"          ___   ___ \n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"         ( o ) ( o )\n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"        /           \\\n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"       |  (  ---  )  |\n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"       ^^^         ^^^\n"+gopherReset)
	fmt.Fprint(writer, gopherRed+"     [⛔ UNSUPPORTED ⛔]\n\n"+gopherReset)
	fmt.Fprint(writer, gopherOrange+"     -- Gopher module --\n"+gopherReset)
	fmt.Fprint(writer, gopherRed+"Input is currently unsupported.\n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"                   \n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"                   \n"+gopherReset)
	fmt.Fprint(writer, gopherGreen+"                   \n"+gopherReset)
}
