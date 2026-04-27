package reverse

import "strings"

// Run reads the ASCII art file at filePath, reconstructs the original
// plain-text string using bannerMap, and returns it.
func Run(filePath string, bannerMap map[rune][]string) (string, error) {
	// TODO(task12): read file → strip trailing blanks → group 8-line blocks
	// → matchSegment each block → join rows with \n
	return "", nil
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
	// TODO(task12): implement greedy column scanner
	return "", nil
}
