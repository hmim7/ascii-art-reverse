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

// RenderGopher writes the hardcoded Peeking Gopher ANSI art to w.
func RenderGopher(w io.Writer) {
	// TODO(task11): replace with full hardcoded ANSI gopher art string
	fmt.Fprintln(w, "  (^_^)  <- Peeking Gopher (placeholder)")
}
