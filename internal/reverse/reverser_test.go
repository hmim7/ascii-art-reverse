package reverse_test

import (
	"testing"

	"ascii-art-reverse/internal/banner"
	"ascii-art-reverse/internal/reverse"
)

func loadBanner(t *testing.T, name string) map[rune][]string {
	t.Helper()
	m, fb, err := banner.Load(name)
	if err != nil {
		t.Fatalf("banner.Load(%q): %v", name, err)
	}
	if fb != nil {
		t.Fatalf("banner.Load(%q): unexpected fallback %+v", name, fb)
	}
	return m
}

func TestRun_StandardExamples(t *testing.T) {
	std := loadBanner(t, "standard")

	tests := []struct {
		file string
		want string
	}{
		{"testdata/example00.txt", "Hello World"},
		{"testdata/example01.txt", "123"},
		{"testdata/example02.txt", `#=\[`},
		{"testdata/example03.txt", "something&234"},
		{"testdata/example04.txt", "abcdefghijklmnopqrstuvwxyz"},
		{"testdata/example05.txt", `\!" #$%&'()*+,-./`},
		{"testdata/example06.txt", ":;{=}?@"},
		{"testdata/example07.txt", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got, err := reverse.Run(tt.file, std)
			if err != nil {
				t.Fatalf("Run(%q): unexpected error: %v", tt.file, err)
			}
			if got != tt.want {
				t.Errorf("Run(%q) = %q, want %q", tt.file, got, tt.want)
			}
		})
	}
}

func TestRun_ShadowBanner(t *testing.T) {
	shadow := loadBanner(t, "shadow")
	got, err := reverse.Run("testdata/shadow_hello_world.txt", shadow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("got %q, want %q", got, "Hello World")
	}
}

func TestRun_MultiLine(t *testing.T) {
	std := loadBanner(t, "standard")
	got, err := reverse.Run("testdata/multi_line_art.txt", std)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Hello\nWorld" {
		t.Errorf("got %q, want %q", got, "Hello\nWorld")
	}
}

func TestRun_InvalidArt(t *testing.T) {
	std := loadBanner(t, "standard")
	_, err := reverse.Run("testdata/invalid_art.txt", std)
	if err == nil {
		t.Error("expected error for invalid art, got nil")
	}
}

func TestRun_FileNotFound(t *testing.T) {
	std := loadBanner(t, "standard")
	_, err := reverse.Run("testdata/nonexistent.txt", std)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
