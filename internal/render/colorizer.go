package render

// ColorToANSI converts any supported color notation to an ANSI CSI start
// sequence. ok=false means the input was unknown or unparsable; the caller
// should warn and skip coloring.
func ColorToANSI(color string) (start string, ok bool) {
	// TODO(task06): implement named, bright-*, #RRGGBB, rgb(), hsl() resolution
	return "", false
}

// WrapWithColor wraps s with ansiStart and the reset code \x1b[0m.
// Short-circuits (returns s unchanged) when either argument is empty.
func WrapWithColor(s, ansiStart string) string {
	if s == "" || ansiStart == "" {
		return s
	}
	return ansiStart + s + "\x1b[0m"
}

func buildColorMask(s, sub []rune, enableColor bool) []bool {
	// TODO(task06): case-sensitive substring scan; all-true when sub is empty
	mask := make([]bool, len(s))
	if !enableColor {
		return mask
	}
	for i := range mask {
		mask[i] = len(sub) == 0
	}
	return mask
}

func buildPerRuneANSI(s []rune, rules []ColorRule) []string {
	// TODO(task06): iterate rules; last matching rule wins on overlapping positions
	return make([]string, len(s))
}

func hslToRGB(h, s, l float64) (r, g, b uint8) {
	// TODO(task06): standard HSL → RGB conversion
	return 0, 0, 0
}
