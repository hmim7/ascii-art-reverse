package banner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FallbackInfo is non-nil when Load fell back to the standard banner.
type FallbackInfo struct {
	Kind string // "not-found" | "invalid"
	Name string // sanitized banner name that was attempted
}

// Load reads the banner font file for the given name and returns a
// map[rune][]string for O(1) character lookup. Falls back to "standard"
// on any error; when a fallback occurs, FallbackInfo is returned so the
// caller can decide how to report it.
func Load(name string) (map[rune][]string, *FallbackInfo, error) {
	return load(name, true)
}

// bannerDir resolves the banner/ directory by walking up from the current
// working directory until it finds a banner/ folder containing standard.txt.
// This allows the loader to work correctly from any package directory.
func bannerDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "banner"
	}
	for {
		// Confirm by checking for standard.txt, not just the directory name.
		probe := filepath.Join(dir, "banner", "standard.txt")
		if _, err := os.Stat(probe); err == nil {
			return filepath.Join(dir, "banner")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "banner"
}

func load(name string, allowFallback bool) (map[rune][]string, *FallbackInfo, error) {
	clean := sanitizeName(name)
	path := filepath.Join(bannerDir(), clean+".txt")

	data, err := os.ReadFile(path) // #nosec G304 — path built from sanitizeName (alphanumeric/-/_ only) + fixed bannerDir
	if err != nil {
		if allowFallback && clean != "standard" {
			return nil, &FallbackInfo{Kind: "not-found", Name: clean}, nil
		}
		return nil, nil, fmt.Errorf("banner %q: %w", clean, err)
	}

	// Normalize line endings to handle CRLF (\r\n) correctly
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) != 855 {
		if allowFallback && clean != "standard" {
			return nil, &FallbackInfo{Kind: "invalid", Name: clean}, nil
		}
		return nil, nil, fmt.Errorf("banner %q: expected 855 lines, got %d", clean, len(lines))
	}

	leading := isLeadingSeparator(lines)
	m := make(map[rune][]string, 95)
	for r := 32; r <= 126; r++ {
		start := (r - 32) * 9
		var src []string
		if leading {
			src = lines[start+1 : start+9]
		} else {
			src = lines[start : start+8]
		}
		art := make([]string, 8)
		copy(art, src)
		m[rune(r)] = art
	}
	return m, nil, nil
}

func sanitizeName(name string) string {
	base := strings.ToLower(filepath.Base(name))
	base = strings.TrimSuffix(base, ".txt")
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return -1
	}, base)
}

func isLeadingSeparator(lines []string) bool {
	// lines[9] is the first line of the '!' (char 33) block.
	// Leading separator format: lines[9] is the separator (empty/blank).
	// Trailing separator format (doom): lines[9] is the first art line of '!'.
	return strings.TrimSpace(lines[9]) == ""
}
