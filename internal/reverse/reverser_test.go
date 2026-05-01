package reverse_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"ascii-art-reverse/internal/banner"
	"ascii-art-reverse/internal/render"
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

func TestRun(t *testing.T) {
	// Cache for loaded banners to avoid redundant I/O
	bannerCache := make(map[string]map[rune][]string)
	getBanner := func(name string) map[rune][]string {
		if m, ok := bannerCache[name]; ok {
			return m
		}
		m := loadBanner(t, name)
		bannerCache[name] = m
		return m
	}

	tests := []struct {
		name       string
		file       string
		bannerName string
		want       string
		wantErr    bool
	}{
		{name: "Example00", file: "testdata/example00.txt", bannerName: "standard", want: "Hello World"},
		{name: "Example01", file: "testdata/example01.txt", bannerName: "standard", want: "123"},
		{name: "Example02", file: "testdata/example02.txt", bannerName: "standard", want: `#=\[`},
		{name: "Example03", file: "testdata/example03.txt", bannerName: "standard", want: "something&234"},
		{name: "Example04", file: "testdata/example04.txt", bannerName: "standard", want: "abcdefghijklmnopqrstuvwxyz"},
		{name: "Example05", file: "testdata/example05.txt", bannerName: "standard", want: `\!" #$%&'()*+,-./`},
		{name: "Example06", file: "testdata/example06.txt", bannerName: "standard", want: ":;{=}?@"},
		{name: "Example07", file: "testdata/example07.txt", bannerName: "standard", want: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{name: "ShadowBanner", file: "testdata/shadow_hello_world.txt", bannerName: "shadow", want: "Hello World"},
		{name: "MultiLine", file: "testdata/multi_line_art.txt", bannerName: "standard", want: "Hello\nWorld"},
		{name: "A47MixedCase", file: "testdata/example_a47.txt", bannerName: "standard", want: "rEvErSe"},
		{name: "A48LowerNumSpace", file: "testdata/example_a48.txt", bannerName: "standard", want: "abc 123 def"},
		{name: "A49SpecialChars", file: "testdata/example_a49.txt", bannerName: "standard", want: "(^_^)"},
		{name: "InvalidArt", file: "testdata/invalid_art.txt", bannerName: "standard", wantErr: true},
		{name: "FileNotFound", file: "testdata/nonexistent.txt", bannerName: "standard", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			bm := getBanner(tt.bannerName)
			got, err := reverse.Run(tt.file, bm)
			if (err != nil) != tt.wantErr {
				st.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				st.Errorf("Run() = %q, want %q", got, tt.want)
			}
		})
	}
}

// FuzzReverseRoundTrip renders a printable-ASCII string to ASCII art, writes it
// to a temp file, then reverses it and checks the result matches the original.
// It verifies the full render→reverse pipeline never panics and stays consistent.
func FuzzReverseRoundTrip(f *testing.F) {
	f.Add("Hello World")
	f.Add("A")
	f.Add("Hi")
	f.Add("abc")
	f.Add("ABC")
	f.Add("Hello!")
	f.Add(" ")
	f.Add("abcdefghijklmnopqrstuvwxyz")
	f.Add("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	f.Add("0123456789")

	bannerMap, fb, err := banner.Load("standard")
	if err != nil || fb != nil {
		f.Fatal("failed to load standard banner for fuzz corpus")
	}

	f.Fuzz(func(st *testing.T, input string) {
		// Filter to printable ASCII only (same as the render pipeline).
		var clean strings.Builder
		for _, r := range input {
			if r >= 32 && r <= 126 {
				clean.WriteRune(r)
			}
		}
		text := clean.String()
		if text == "" {
			return
		}

		// Render to ASCII art in a temp file.
		// ParseInput preprocesses escape sequences (e.g. \\ → \, \n → newline),
		// so compare the reverse output against the processed form, not raw input.
		var buf bytes.Buffer
		segs := render.ParseInput(text)
		processed := strings.Join(segs, "\n")
		render.RenderAlignedWithColorRules(&buf, segs, bannerMap, nil, "left", 0)
		if buf.Len() == 0 {
			return
		}

		tmp, err := os.CreateTemp(st.TempDir(), "fuzz-art-*.txt")
		if err != nil {
			st.Fatal(err)
		}
		if _, werr := tmp.Write(buf.Bytes()); werr != nil {
			_ = tmp.Close()
			st.Fatal(werr)
		}
		_ = tmp.Close()

		// Reverse the art file back to text and compare against the processed form.
		got, err := reverse.Run(tmp.Name(), bannerMap)
		if err != nil {
			st.Fatalf("Run() error for input %q: %v", text, err)
		}
		if got != processed {
			st.Fatalf("round-trip mismatch: processed %q, reversed to %q", processed, got)
		}
	})
}

func BenchmarkRun(b *testing.B) {
	bannerMap, _, _ := banner.Load("standard")
	cases := []struct {
		name string
		file string
	}{
		{"hello_world", "testdata/example00.txt"},
		{"alpha_26", "testdata/example07.txt"},
		{"multi_line", "testdata/multi_line_art.txt"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, err := reverse.Run(tc.file, bannerMap)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
