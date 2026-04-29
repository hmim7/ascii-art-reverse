package render

// Mapper wraps the banner map in an O(1) closure for use by the render pipeline.
// Out-of-range rune (<32 or >126) or rune missing from the map → 8 empty strings.
func Mapper(bannerMap map[rune][]string) func(rune) []string {
	return func(r rune) []string {
		if r < 32 || r > 126 {
			return make([]string, 8)
		}
		if lines, ok := bannerMap[r]; ok {
			return lines
		}
		return make([]string, 8)
	}
}
