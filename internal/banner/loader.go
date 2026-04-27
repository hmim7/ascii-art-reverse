package banner

// Load reads the banner font file for the given name and returns a
// map[rune][]string for O(1) character lookup. Falls back to "standard"
// with a stderr warning on any error.
func Load(name string) (map[rune][]string, error) {
	// TODO(task04): implement sanitizeName, isLeadingSeparator, file reading,
	// 855-line validation, leading/trailing separator dispatch, and map build.
	return nil, nil
}

func sanitizeName(name string) string {
	// TODO(task04): filepath.Base → lowercase → strip .txt → [a-z0-9\-_] filter
	return name
}

func isLeadingSeparator(lines []string) bool {
	// TODO(task04): inspect ! block lines 9–17 to detect separator format
	return true
}
