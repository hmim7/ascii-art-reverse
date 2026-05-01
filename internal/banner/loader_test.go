package banner_test

import (
	"testing"

	"ascii-art-reverse/internal/banner"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name         string
		banner       string
		fallbackKind string
		wantErr      bool
		wantFallback bool
	}{
		{name: "Standard", banner: "standard"},
		{name: "Shadow", banner: "shadow"},
		{name: "Doom", banner: "doom"},
		{name: "Block", banner: "block"},
		{name: "Thinkertoy", banner: "thinkertoy"},
		{name: "Dancing", banner: "dancing"},
		{name: "Greek", banner: "greek"},
		{name: "NotFound", banner: "nonexistent", wantFallback: true, fallbackKind: "not-found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			m, fb, err := banner.Load(tt.banner)
			if (err != nil) != tt.wantErr {
				st.Fatalf("Load(%q) error = %v, wantErr %v", tt.banner, err, tt.wantErr)
			}
			if tt.wantFallback {
				if fb == nil {
					st.Fatal("expected FallbackInfo, got nil")
				}
				if m != nil {
					st.Errorf("expected nil map when fallback occurs, got map with %d entries", len(m))
				}
				if fb.Kind != tt.fallbackKind {
					st.Errorf("FallbackInfo.Kind = %q, want %q", fb.Kind, tt.fallbackKind)
				}
				return
			}
			if fb != nil {
				st.Errorf("unexpected FallbackInfo: %+v", fb)
			}
			if m == nil {
				st.Fatal("expected map, got nil")
			}
			if len(m) != 95 {
				st.Errorf("expected 95 characters, got %d", len(m))
			}
			for r := rune(32); r <= 126; r++ {
				art, ok := m[r]
				if !ok {
					st.Errorf("rune %d (%q) missing from map", r, r)
					continue
				}
				if len(art) != 8 {
					st.Errorf("rune %d (%q): expected 8 art lines, got %d", r, r, len(art))
				}
			}
		})
	}
}

// FuzzLoad checks that Load never panics on arbitrary banner name input.
// Unknown names fall back to "standard"; the function must always return
// either a valid map or a FallbackInfo — never a panic.
func FuzzLoad(f *testing.F) {
	f.Add("standard")
	f.Add("shadow")
	f.Add("nonexistent")
	f.Add("")
	f.Add("../../etc/passwd")
	f.Add("STANDARD")
	f.Add("stan dard")
	f.Add("../banner/standard")
	f.Add("standard.txt")

	f.Fuzz(func(st *testing.T, name string) {
		m, fb, err := banner.Load(name)
		if err != nil {
			// Only "standard" should ever return an error (if the file is missing).
			// For all other names Load must return a FallbackInfo instead.
			st.Errorf("Load(%q) returned unexpected error: %v", name, err)
			return
		}
		if m == nil && fb == nil {
			st.Errorf("Load(%q): both map and FallbackInfo are nil", name)
		}
	})
}

func BenchmarkLoad(b *testing.B) {
	for _, name := range []string{"standard", "shadow", "doom"} {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = banner.Load(name)
			}
		})
	}
}
