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
	for _, art := range bannerMap {
		if len(art) > 0 {
			widthSet[len(art[0])] = struct{}{}
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
		seg, err := matchSegment(block, rev, widths)
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
	for r, lines := range bannerMap {
		rev[strings.Join(lines, "|")] = r
	}
	return rev
}

// matchSegment greedy-scans one 8-line block left-to-right.
// At each column offset it tries all known glyph widths; first match consumed.
func matchSegment(lines []string, rev map[string]rune, widths []int) (string, error) {
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
	col := 0
	// Scan until we hit the visible end of the art.
	// Special case: if maxCol is 0 (all 8 lines were trimmed), try to match at least one glyph.
	for col < maxCol || (col == 0 && maxCol == 0) {
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
			if col >= maxCol && col > 0 {
				break // No match found beyond the visible art; we are done.
			}
			return "", fmt.Errorf("unrecognized glyph at column %d", col)
		}
	}
	return result.String(), nil
}
