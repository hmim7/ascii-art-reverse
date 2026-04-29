package banner_test

import (
	"testing"

	"ascii-art-reverse/internal/banner"
)

// TODO(task04): expand with full banner loading test table once banner/ files exist.

func TestLoad_Standard(t *testing.T) {
	m, fb, err := banner.Load("standard")
	if err != nil {
		t.Fatalf("Load(standard) error: %v", err)
	}
	if fb != nil {
		t.Errorf("expected no fallback for standard, got %+v", fb)
	}
	if len(m) != 95 {
		t.Errorf("expected 95 characters, got %d", len(m))
	}
}

func TestLoad_InvalidFallback(t *testing.T) {
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
