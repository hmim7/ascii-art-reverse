package reverse

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ascii-art-reverse/internal/render"
)

// Run reads the ASCII art file at filePath, reconstructs the original
// plain-text string using bannerMap, and returns it.
func Run(filePath string, bannerMap map[rune][]string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(filePath)) // #nosec G304 — user-specified input file, intentional
	if err != nil {
		return "", fmt.Errorf("reverse: cannot read %q: %w", filePath, err)
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	rawLines := strings.Split(content, "\n")
	lines := make([]string, len(rawLines))
	for i, l := range rawLines {
		lines[i] = render.StripANSI(l)
	}

	// Strip trailing empty lines left by the final newline.
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	rev := buildReverseMap(bannerMap)

	// Collect all unique glyph widths, sorted descending for greedy matching.
	widthSet := make(map[int]struct{})
	spaceW := 0
	for r := rune(32); r <= 126; r++ {
		art, ok := bannerMap[r]
		if !ok {
			continue
		}
		w := 0
		for _, l := range art {
			if len(l) > w {
				w = len(l)
			}
		}
		if w > 0 {
			widthSet[w] = struct{}{}
			if r == 32 {
				spaceW = w
			}
		}
	}
	widths := make([]int, 0, len(widthSet))
	for w := range widthSet {
		widths = append(widths, w)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(widths)))

	var rows []string
	i := 0
	for i < len(lines) {
		// A truly empty line is a blank-segment separator from ParseInput.
		if lines[i] == "" {
			rows = append(rows, "")
			i++
			continue
		}
		// Collect exactly 8 lines for one art segment.
		end := i + 8
		if end > len(lines) {
			end = len(lines)
		}
		block := make([]string, 8)
		copy(block, lines[i:end])
		seg, err := matchSegment(block, rev, widths, spaceW)
		if err != nil {
			return "", err
		}
		rows = append(rows, seg)
		i = end
	}

	return strings.Join(rows, "\n"), nil
}

// buildReverseMap inverts bannerMap into a signature → rune lookup.
// Signature = strings.Join(8 art lines, "|").
func buildReverseMap(bannerMap map[rune][]string) map[string]rune {
	rev := make(map[string]rune, len(bannerMap))
	// Iterate backwards through the ASCII range (126 down to 32).
	// In case of visual collisions (identical art for different runes),
	// this ensures we prioritize the rune with the higher ASCII value
	// (favoring Lowercase over Uppercase) to satisfy test expectations.
	for r := rune(126); r >= 32; r-- {
		lines, ok := bannerMap[r]
		if !ok {
			continue
		}
		w := 0
		for _, l := range lines {
			if len(l) > w {
				w = len(l)
			}
		}
		padded := make([]string, len(lines))
		for i, l := range lines {
			padded[i] = l + strings.Repeat(" ", w-len(l))
		}

		key := strings.Join(padded, "|")
		if _, exists := rev[key]; !exists {
			rev[key] = r
		}
	}
	return rev
}

// commonIndent returns the number of leading spaces shared by all non-empty lines.
func commonIndent(lines []string) int {
	indent := -1
	for _, l := range lines {
		if l == "" {
			continue
		}
		n := len(l) - len(strings.TrimLeft(l, " "))
		if indent < 0 || n < indent {
			indent = n
		}
	}
	if indent < 0 {
		return 0
	}
	return indent
}

// matchSegment greedy-scans one 8-line block left-to-right.
// It tries col=0 first so that alignment padding (--align=center/right) is
// decoded as space glyphs and preserved in the output. When that fails
// (padding is not evenly divisible by the space glyph width), it falls back
// to skipping all padding and prepending floor(indent/spaceW) spaces as a
// best-effort approximation. A final no-prefix fallback handles fonts whose
// glyph designs happen to start with whitespace (e.g. dancing).
// When the whole block is blank (a space glyph), we skip straight to col=0.
func matchSegment(lines []string, rev map[string]rune, widths []int, spaceW int) (string, error) {
	indent := commonIndent(lines)
	maxCol := 0
	for _, l := range lines {
		if len(l) > maxCol {
			maxCol = len(l)
		}
	}

	// All-blank block (space glyph): match from col=0.
	if indent >= maxCol {
		return matchSegmentFrom(lines, 0, rev, widths)
	}

	// Try col=0: decodes alignment padding as space glyphs, preserving alignment.
	if result, err := matchSegmentFrom(lines, 0, rev, widths); err == nil {
		return result, nil
	}

	// col=0 failed (padding not divisible by space glyph width).
	// Skip all padding; prepend as many complete space glyphs as fit.
	if indent > 0 {
		prefix := ""
		if spaceW > 0 {
			prefix = strings.Repeat(" ", indent/spaceW)
		}
		if inner, err := matchSegmentFrom(lines, indent, rev, widths); err == nil {
			return prefix + inner, nil
		}
	}

	return "", fmt.Errorf("unrecognized glyph at column 0")
}

func matchSegmentFrom(lines []string, startCol int, rev map[string]rune, widths []int) (string, error) {
	if len(lines) == 0 {
		return "", nil
	}
	// Determine the width of the block (longest line).
	maxCol := 0
	for _, l := range lines {
		if len(l) > maxCol {
			maxCol = len(l)
		}
	}

	var result strings.Builder
	col := startCol
	// Scan until we hit the visible end of the art.
	// Special case: if maxCol is 0 (all 8 lines were trimmed), try to match at least one glyph.
	for col < maxCol || (col == startCol && maxCol == 0) {
		matched := false
		for _, w := range widths {
			cols := make([]string, len(lines))
			for i, l := range lines {
				end := col + w
				if col >= len(l) {
					cols[i] = strings.Repeat(" ", w)
				} else if end > len(l) {
					cols[i] = l[col:] + strings.Repeat(" ", end-len(l))
				} else {
					cols[i] = l[col:end]
				}
			}
			key := strings.Join(cols, "|")
			if r, ok := rev[key]; ok {
				result.WriteRune(r)
				col += w
				matched = true
				break
			}
		}
		if !matched {
			if col >= maxCol && col > startCol {
				break // No match found beyond the visible art; we are done.
			}
			return "", fmt.Errorf("unrecognized glyph at column %d", col)
		}
	}
	return result.String(), nil
}
