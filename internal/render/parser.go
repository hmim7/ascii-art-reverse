package render

import "strings"

// ParseInput preprocesses the raw input string, splits on newlines, and
// filters each part to printable ASCII [32, 126].
func ParseInput(input string) []string {
	// TODO(task05): implement preprocess → Split → filterASCII pipeline
	parts := strings.Split(input, "\n")
	return parts
}

func preprocess(s string) string {
	// TODO(task05): normalize \r\n/\r; collapse escape sequences \n \t \\ etc.
	return s
}

func filterASCII(s string) string {
	// TODO(task05): keep only runes in [32, 126]
	return s
}
