package render

import "strings"

// ParseInput preprocesses the raw input string, splits on newlines, and
// filters each part to printable ASCII [32, 126].
func ParseInput(input string) []string {
	p := preprocess(input)
	parts := strings.Split(p, "\n")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = filterASCII(part)
	}
	return result
}

func preprocess(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		// Normalize actual \r\n and \r to \n.
		if s[i] == '\r' {
			if i+1 < len(s) && s[i+1] == '\n' {
				i++ // skip \r; the \n will be written on next iteration
			}
			b.WriteByte('\n')
			i++
			continue
		}
		// Expand backslash escape sequences from CLI input.
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
				i += 2
				continue
			case 't':
				b.WriteByte('\t')
				i += 2
				continue
			case '\\':
				b.WriteByte('\\')
				i += 2
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func filterASCII(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= 32 && r <= 126 {
			b.WriteRune(r)
		}
	}
	return b.String()
}
