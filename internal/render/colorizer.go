package render

import (
	"fmt"
	"strconv"
	"strings"
)

var namedColors = map[string]string{
	"black":          "\x1b[30m",
	"red":            "\x1b[31m",
	"green":          "\x1b[32m",
	"yellow":         "\x1b[33m",
	"blue":           "\x1b[34m",
	"magenta":        "\x1b[35m",
	"cyan":           "\x1b[36m",
	"white":          "\x1b[37m",
	"orange":         "\x1b[38;5;208m",
	"bright-black":   "\x1b[90m",
	"bright-red":     "\x1b[91m",
	"bright-green":   "\x1b[92m",
	"bright-yellow":  "\x1b[93m",
	"bright-blue":    "\x1b[94m",
	"bright-magenta": "\x1b[95m",
	"bright-cyan":    "\x1b[96m",
	"bright-white":   "\x1b[97m",
}

// ColorToANSI converts any supported color notation to an ANSI CSI start
// sequence. ok=false means the input was unknown or unparsable; the caller
// should warn and skip coloring.
func ColorToANSI(color string) (start string, ok bool) {
	lower := strings.ToLower(color)

	if ansi, found := namedColors[lower]; found {
		return ansi, true
	}

	// Hex: #RRGGBB
	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r, err1 := strconv.ParseUint(color[1:3], 16, 8)
		g, err2 := strconv.ParseUint(color[3:5], 16, 8)
		b, err3 := strconv.ParseUint(color[5:7], 16, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b), true
		}
	}

	// Strip spaces for RGB/HSL parsing.
	compact := strings.ReplaceAll(lower, " ", "")

	// RGB: rgb(r,g,b)
	if strings.HasPrefix(compact, "rgb(") && strings.HasSuffix(compact, ")") {
		inner := compact[4 : len(compact)-1]
		parts := strings.Split(inner, ",")
		if len(parts) == 3 {
			r, err1 := strconv.ParseUint(parts[0], 10, 8)
			g, err2 := strconv.ParseUint(parts[1], 10, 8)
			b, err3 := strconv.ParseUint(parts[2], 10, 8)
			if err1 == nil && err2 == nil && err3 == nil && r <= 255 && g <= 255 && b <= 255 {
				return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b), true
			}
		}
	}

	// HSL: hsl(h,s%,l%)
	if strings.HasPrefix(compact, "hsl(") && strings.HasSuffix(compact, ")") {
		inner := compact[4 : len(compact)-1]
		parts := strings.Split(inner, ",")
		if len(parts) == 3 {
			h, err1 := strconv.ParseFloat(parts[0], 64)
			sStr := strings.TrimSuffix(parts[1], "%")
			lStr := strings.TrimSuffix(parts[2], "%")
			s, err2 := strconv.ParseFloat(sStr, 64)
			l, err3 := strconv.ParseFloat(lStr, 64)
			if err1 == nil && err2 == nil && err3 == nil {
				r, g, b := hslToRGB(h, s/100, l/100)
				return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b), true
			}
		}
	}

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
	mask := make([]bool, len(s))
	if !enableColor {
		return mask
	}
	if len(sub) == 0 {
		for i := range mask {
			mask[i] = true
		}
		return mask
	}
	// Case-sensitive substring scan.
	for i := 0; i <= len(s)-len(sub); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if s[i+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			for j := 0; j < len(sub); j++ {
				mask[i+j] = true
			}
		}
	}
	return mask
}

func buildPerRuneANSI(s []rune, rules []ColorRule) []string {
	result := make([]string, len(s))
	for _, rule := range rules {
		mask := buildColorMask(s, []rune(rule.Substring), true)
		for i, on := range mask {
			if on {
				result[i] = rule.ANSIStart // last matching rule wins
			}
		}
	}
	return result
}

func hslToRGB(h, s, l float64) (r, g, b uint8) {
	if s == 0 {
		v := uint8(l * 255)
		return v, v, v
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	h /= 360
	hue := func(t float64) uint8 {
		if t < 0 {
			t++
		}
		if t > 1 {
			t--
		}
		switch {
		case t < 1.0/6:
			return uint8((p + (q-p)*6*t) * 255)
		case t < 0.5:
			return uint8(q * 255)
		case t < 2.0/3:
			return uint8((p + (q-p)*(2.0/3-t)*6) * 255)
		default:
			return uint8(p * 255)
		}
	}
	return hue(h + 1.0/3), hue(h), hue(h - 1.0/3)
}
