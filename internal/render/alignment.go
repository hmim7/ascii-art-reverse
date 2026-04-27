package render

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// DetectTerminalWidth returns the terminal column width using three sources
// in priority order: COLUMNS env var → ioctl → fallback 80.
func DetectTerminalWidth() int {
	if v := os.Getenv("COLUMNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if n, err := getWinSize(os.Stdout.Fd()); err == nil && n > 0 {
		return n
	}
	return 80
}

// ResolveWidth returns 80 when writing to a file, otherwise DetectTerminalWidth.
func ResolveWidth(usingFileOutput bool) int {
	if usingFileOutput {
		return 80
	}
	return DetectTerminalWidth()
}

func visibleWidth(s string) int {
	// TODO(task08): strip \x1b[…m CSI sequences; count remaining bytes
	return len(s)
}

func padLine(line, align string, width int) string {
	// TODO(task08): implement right / center padding; left returns unchanged
	return line
}

func alignRenderedLines(lines []string, align string, width int) []string {
	// TODO(task08): apply padLine to each line
	return lines
}

func extractWordIndices(runes []rune) [][]int {
	// TODO(task08): split on ' ' only; return slice of index-slices per word
	return nil
}

func justifyRenderedSegment(runes []rune, get func(rune) []string, perRuneANSI []string, width int) []string {
	// TODO(task08): gap distribution: baseGap=total/gaps; remainder=total%gaps
	return nil
}

func buildRuneParts(runes []rune, get func(rune) []string, perRuneANSI []string) [][]string {
	// TODO(task08): for each rune fetch 8-line glyph; wrap with WrapWithColor
	parts := make([][]string, len(runes))
	for i, r := range runes {
		lines := get(r)
		cols := make([]string, 8)
		for j := range cols {
			if j < len(lines) {
				cols[j] = WrapWithColor(lines[j], perRuneANSI[i])
			}
		}
		parts[i] = cols
	}
	return parts
}

func joinRuneParts(runeParts [][]string) []string {
	// TODO(task08): concatenate rune columns per row → 8 complete lines
	lines := make([]string, 8)
	for _, cols := range runeParts {
		for row := 0; row < 8 && row < len(cols); row++ {
			lines[row] += cols[row]
		}
	}
	return lines
}

func renderSegment(seg string, get func(rune) []string, rules []ColorRule, align string, width int) []string {
	// TODO(task08): full per-segment render → color → align pipeline
	runes := []rune(seg)
	ansi := buildPerRuneANSI(runes, rules)
	parts := buildRuneParts(runes, get, ansi)
	lines := joinRuneParts(parts)
	return alignRenderedLines(lines, align, width)
}

// RenderAlignedWithColorRules is the core render entry point.
// Empty segment → one blank line; non-empty → 8 lines.
func RenderAlignedWithColorRules(w io.Writer, segs []string, m map[rune][]string, rules []ColorRule, align string, width int) {
	get := Mapper(m)
	for _, seg := range segs {
		if seg == "" {
			fmt.Fprintln(w)
			continue
		}
		for _, line := range renderSegment(seg, get, rules, align, width) {
			fmt.Fprintln(w, line)
		}
	}
}
