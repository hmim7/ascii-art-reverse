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
		// Handle real tab character.
		if s[i] == '\t' {
			b.WriteString("   ")
			i++
			continue
		}
		// Normalize actual \r\n and \r to \n.
		if s[i] == '\r' {
			if i+1 < len(s) && s[i+1] == '\n' {
				i++ // skip \r; the \n will be written on next iteration
			}
			b.WriteByte('\n')
			i++
			continue
		}
		// Handle backslash escape sequences.
		if s[i] == '\\' && i+1 < len(s) {
			if s[i+1] == '\\' {
				// Check for \\n or \\t sequences.
				if i+2 < len(s) {
					if s[i+2] == 'n' {
						b.WriteString("\\\n")
						i += 3
						continue
					}
					if s[i+2] == 't' {
						b.WriteString("\\   ")
						i += 3
						continue
					}
				}
				b.WriteByte('\\')
				i += 2
				continue
			}
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
				i += 2
				continue
			case 't':
				b.WriteString("   ")
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
