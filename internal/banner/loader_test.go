package banner_test

import (
	"testing"

	"ascii-art-reverse/internal/banner"
)

func TestLoad_AllBanners(t *testing.T) {
	banners := []string{
		"standard",
		"shadow",
		"doom",
		"block",
		"thinkertoy",
		"dancing",
		"greek",
	}

	for _, name := range banners {
		t.Run(name, func(t *testing.T) {
			m, fb, err := banner.Load(name)
			if err != nil {
				t.Fatalf("Load(%q) error: %v", name, err)
			}
			if fb != nil {
				t.Errorf("Load(%q): expected no fallback, got %+v", name, fb)
			}
			if len(m) != 95 {
				t.Errorf("Load(%q): expected 95 characters, got %d", name, len(m))
			}
		})
	}
}

func TestLoad_NotFound(t *testing.T) {
	m, fb, err := banner.Load("nonexistent")
	if err != nil {
		t.Fatalf("Load(nonexistent) should return FallbackInfo, not error: %v", err)
	}
	if fb == nil {
		t.Fatal("expected FallbackInfo for nonexistent banner, got nil")
	}
	if fb.Kind != "not-found" {
		t.Errorf("FallbackInfo.Kind = %q, want %q", fb.Kind, "not-found")
	}
	if m != nil {
		t.Errorf("expected nil map when fallback occurs, got map with %d entries", len(m))
	}
}

func TestLoad_CoversPrintableASCII(t *testing.T) {
	m, _, err := banner.Load("standard")
	if err != nil {
		t.Fatalf("Load(standard): %v", err)
	}
	for r := rune(32); r <= 126; r++ {
		art, ok := m[r]
		if !ok {
			t.Errorf("rune %d (%q) missing from standard map", r, r)
			continue
		}
		if len(art) != 8 {
			t.Errorf("rune %d (%q): expected 8 art lines, got %d", r, r, len(art))
		}
	}
}
