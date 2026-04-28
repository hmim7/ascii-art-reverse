package reverse_test

import (
	"path/filepath"
	"testing"

	"ascii-art-reverse/internal/banner"
	"ascii-art-reverse/internal/reverse"
)

func TestRun_Examples(t *testing.T) {
	bannerMap, err := banner.Load("standard")
	if err != nil {
		t.Fatalf("load standard banner: %v", err)
	}

	tests := []struct {
		file string
		want string
	}{
		{"example00.txt", "Hello World"},
		{"example01.txt", "123"},
		{"example02.txt", "#=\\["},
		{"example03.txt", "something&234"},
		{"example04.txt", "abcdefghijklmnopqrstuvwxyz"},
		{"example05.txt", "\\!\" #$%&'()*+,-./"},
		{"example06.txt", ":;{=}?@"},
		{"example07.txt", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got, err := reverse.Run(filepath.Join("testdata", tt.file), bannerMap)
			if err != nil {
				t.Fatalf("reverse run failed: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
