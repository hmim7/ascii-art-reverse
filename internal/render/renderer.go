package render

import "io"

// ColorRule specifies an ANSI color to apply to an optional substring.
// Empty Substring means the whole rendered string is colored.
type ColorRule struct {
	ANSIStart string
	Substring string
}

// Render renders segments with no color and left alignment.
func Render(w io.Writer, segs []string, m map[rune][]string) {
	RenderAlignedWithColorRules(w, segs, m, nil, "left", 0)
}

// RenderWithColor renders with a single color rule applied to an optional substring.
func RenderWithColor(w io.Writer, segs []string, m map[rune][]string, ansi, sub string) {
	RenderAlignedWithColorRules(w, segs, m, []ColorRule{{ANSIStart: ansi, Substring: sub}}, "left", 0)
}

// RenderWithColorRules renders with multiple color rules, left alignment.
func RenderWithColorRules(w io.Writer, segs []string, m map[rune][]string, rules []ColorRule) {
	RenderAlignedWithColorRules(w, segs, m, rules, "left", 0)
}

// RenderAligned renders with no color and the specified alignment.
func RenderAligned(w io.Writer, segs []string, m map[rune][]string, align string, width int) {
	RenderAlignedWithColorRules(w, segs, m, nil, align, width)
}

// RenderAlignedWithColor renders with a single color rule and alignment.
func RenderAlignedWithColor(w io.Writer, segs []string, m map[rune][]string, ansi, sub, align string, width int) {
	RenderAlignedWithColorRules(w, segs, m, []ColorRule{{ANSIStart: ansi, Substring: sub}}, align, width)
}

