package banner_test

import (
	"testing"

	"ascii-art-reverse/internal/banner"
)

// TODO(task04): expand with full banner loading test table once banner/ files exist.

func TestLoad_Standard(t *testing.T) {
	m, err := banner.Load("standard")
	if err != nil {
		t.Fatalf("Load(standard) error: %v", err)
	}
	if len(m) != 95 {
		t.Errorf("expected 95 characters, got %d", len(m))
	}
}

func TestLoad_InvalidFallback(t *testing.T) {
	m, err := banner.Load("nonexistent")
	if err != nil {
		t.Fatalf("Load(nonexistent) should fall back, not error: %v", err)
	}
	if len(m) == 0 {
		t.Error("expected standard fallback map, got empty")
	}
}
