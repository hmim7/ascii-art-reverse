package render

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

// visibleWidth counts the printable byte width of s, ignoring ANSI CSI sequences.
func visibleWidth(s string) int {
	width := 0
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			i += 2
			for i < len(s) && s[i] != 'm' {
				i++
			}
			if i < len(s) {
				i++ // consume 'm'
			}
			continue
		}
		width++
		i++
	}
	return width
}

func padLine(line, align string, width int) string {
	if align == "left" || align == "justify" || align == "" {
		return line
	}
	vw := visibleWidth(line)
	pad := width - vw
	if pad <= 0 {
		return line
	}
	switch align {
	case "right":
		return strings.Repeat(" ", pad) + line
	case "center":
		return strings.Repeat(" ", pad/2) + line
	}
	return line
}

func alignRenderedLines(lines []string, align string, width int) []string {
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = padLine(l, align, width)
	}
	return result
}

// extractWordIndices returns [{start,end}, ...] index pairs for each word in runes,
// splitting on space characters only.
func extractWordIndices(runes []rune) [][]int {
	var words [][]int
	i := 0
	for i < len(runes) {
		for i < len(runes) && runes[i] == ' ' {
			i++
		}
		if i >= len(runes) {
			break
		}
		start := i
		for i < len(runes) && runes[i] != ' ' {
			i++
		}
		words = append(words, []int{start, i})
	}
	return words
}

func justifyRenderedSegment(runes []rune, get func(rune) []string, perRuneANSI []string, width int) []string {
	wordIndices := extractWordIndices(runes)
	if len(wordIndices) <= 1 {
		// Single word: fall back to left alignment.
		parts := buildRuneParts(runes, get, perRuneANSI)
		return joinRuneParts(parts)
	}

	// Sum the visual width of all non-space glyphs.
	contentWidth := 0
	for _, r := range runes {
		if r != ' ' {
			lines := get(r)
			if len(lines) > 0 {
				contentWidth += len(lines[0])
			}
		}
	}

	gaps := len(wordIndices) - 1
	totalSpace := width - contentWidth
	if totalSpace < 0 {
		totalSpace = 0
	}
	baseGap := totalSpace / gaps
	remainder := totalSpace % gaps

	output := make([]string, 8)
	for wi, wr := range wordIndices {
		// Append each glyph in the word.
		for ri := wr[0]; ri < wr[1]; ri++ {
			glyphs := get(runes[ri])
			for row := 0; row < 8 && row < len(glyphs); row++ {
				output[row] += WrapWithColor(glyphs[row], perRuneANSI[ri])
			}
		}
		// Append gap spaces between words (first remainder gaps get +1).
		if wi < gaps {
			spaces := baseGap
			if wi < remainder {
				spaces++
			}
			gap := strings.Repeat(" ", spaces)
			for row := 0; row < 8; row++ {
				output[row] += gap
			}
		}
	}
	return output
}

func buildRuneParts(runes []rune, get func(rune) []string, perRuneANSI []string) [][]string {
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
	lines := make([]string, 8)
	for _, cols := range runeParts {
		for row := 0; row < 8 && row < len(cols); row++ {
			lines[row] += cols[row]
		}
	}
	return lines
}

func renderSegment(seg string, get func(rune) []string, rules []ColorRule, align string, width int) []string {
	runes := []rune(seg)
	ansi := buildPerRuneANSI(runes, rules)
	if align == "justify" {
		return justifyRenderedSegment(runes, get, ansi, width)
	}
	parts := buildRuneParts(runes, get, ansi)
	lines := joinRuneParts(parts)
	return alignRenderedLines(lines, align, width)
}

// RenderAlignedWithColorRules is the core render entry point.
// Empty segment → one blank line; non-empty → 8 lines.
// Special case: when every segment is empty (blank-only input), write exactly
// len(segs)-1 blank lines — one per \n character in the original input —
// so that "" produces no output and "\n" produces exactly one blank line.
func RenderAlignedWithColorRules(w io.Writer, segs []string, m map[rune][]string, rules []ColorRule, align string, width int) {
	get := Mapper(m)

	allEmpty := true
	for _, s := range segs {
		if s != "" {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		for i := 0; i < len(segs)-1; i++ {
			fmt.Fprintln(w)
		}
		return
	}

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
