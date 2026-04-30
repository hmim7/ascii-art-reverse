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
	for _, r := range text {
		if r == '\n' || r == '\t' || r == '\\' {
			continue
		}
		if r < 32 || r > 126 {
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
	lines := []string{
		gopherGreen + "                   \n" + gopherReset,
		gopherGreen + "                   \n" + gopherReset,
		gopherGreen + "                   \n" + gopherReset,
		gopherRed + "🐸 Ooops! Peeking Gopher appeared...\n" + gopherReset,
		gopherGreen + "          ___   ___ \n" + gopherReset,
		gopherGreen + "         ( o ) ( o )\n" + gopherReset,
		gopherGreen + "        /           \\\n" + gopherReset,
		gopherGreen + "       |  (  ---  )  |\n" + gopherReset,
		gopherGreen + "       ^^^         ^^^\n" + gopherReset,
		gopherRed + "     [⛔ UNSUPPORTED ⛔]\n\n" + gopherReset,
		gopherOrange + "     -- Gopher module --\n" + gopherReset,
		gopherRed + "Input is currently unsupported.\n" + gopherReset,
		gopherGreen + "                   \n" + gopherReset,
		gopherGreen + "                   \n" + gopherReset,
		gopherGreen + "                   \n" + gopherReset,
	}
	for _, l := range lines {
		_, _ = fmt.Fprint(writer, l)
	}
}
