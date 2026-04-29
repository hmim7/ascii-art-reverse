package render_test

import (
	"bytes"
	"strings"
	"testing"

	"ascii-art-reverse/internal/render"
)

// --- ParseInput (task05) ---

func TestParseInput_NewlineSplit(t *testing.T) {
	got := render.ParseInput("Hello\nWorld")
	if len(got) != 2 {
		t.Errorf("expected 2 segments, got %d", len(got))
	}
}

func TestParseInput_Empty(t *testing.T) {
	// strings.Split("", "\n") = [""] → 1 segment (empty string).
	// RenderAlignedWithColorRules handles all-empty via special case.
	got := render.ParseInput("")
	if len(got) != 1 {
		t.Errorf("expected 1 segment for empty input, got %d", len(got))
	}
	if got[0] != "" {
		t.Errorf("expected empty segment, got %q", got[0])
	}
}

func TestParseInput_OnlyNewline(t *testing.T) {
	// "\\n" is preprocessed to a real newline → split → ["", ""] (2 segments).
	got := render.ParseInput("\\n")
	if len(got) != 2 {
		t.Errorf("expected 2 segments for single-newline input, got %d", len(got))
	}
}

func TestParseInput_MultipleNewlines(t *testing.T) {
	got := render.ParseInput("A\\nB\\nC")
	if len(got) != 3 {
		t.Errorf("expected 3 segments, got %d: %v", len(got), got)
	}
	if got[0] != "A" || got[1] != "B" || got[2] != "C" {
		t.Errorf("unexpected segments: %v", got)
	}
}

func TestParseInput_FiltersNonASCII(t *testing.T) {
	got := render.ParseInput("héllo")
	if len(got) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(got))
	}
	if strings.ContainsRune(got[0], 'é') {
		t.Errorf("non-ASCII rune 'é' was not filtered: %q", got[0])
	}
}

// --- Mapper (task05) ---

func TestMapper_ValidRune(t *testing.T) {
	m := map[rune][]string{
		'A': {"line1", "line2", "line3", "line4", "line5", "line6", "line7", "line8"},
	}
	get := render.Mapper(m)
	lines := get('A')
	if len(lines) != 8 || lines[0] != "line1" {
		t.Errorf("Mapper returned unexpected lines for 'A': %v", lines)
	}
}

func TestMapper_MissingRune(t *testing.T) {
	get := render.Mapper(map[rune][]string{})
	lines := get('Z')
	if len(lines) != 8 {
		t.Errorf("expected 8 empty lines for missing rune, got %d", len(lines))
	}
	for i, l := range lines {
		if l != "" {
			t.Errorf("line %d should be empty, got %q", i, l)
		}
	}
}

func TestMapper_OutOfRangeBelow(t *testing.T) {
	get := render.Mapper(map[rune][]string{'\x01': {"x"}})
	lines := get('\x01')
	if len(lines) != 8 {
		t.Errorf("expected 8 empty lines for rune <32, got %d", len(lines))
	}
}

func TestMapper_OutOfRangeAbove(t *testing.T) {
	get := render.Mapper(map[rune][]string{'é': {"x"}})
	lines := get('é')
	if len(lines) != 8 {
		t.Errorf("expected 8 empty lines for rune >126, got %d", len(lines))
	}
}

// --- ColorToANSI (task06) ---

func TestColorToANSI_NamedColors(t *testing.T) {
	tests := []struct {
		color string
		want  string
	}{
		{"red", "\x1b[31m"},
		{"green", "\x1b[32m"},
		{"blue", "\x1b[34m"},
		{"cyan", "\x1b[36m"},
		{"white", "\x1b[37m"},
		{"black", "\x1b[30m"},
	}
	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			got, ok := render.ColorToANSI(tt.color)
			if !ok {
				t.Fatalf("ColorToANSI(%q): ok=false", tt.color)
			}
			if got != tt.want {
				t.Errorf("ColorToANSI(%q) = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}

func TestColorToANSI_Hex(t *testing.T) {
	got, ok := render.ColorToANSI("#FF0000")
	if !ok {
		t.Fatal("ColorToANSI(#FF0000): ok=false")
	}
	if got != "\x1b[38;2;255;0;0m" {
		t.Errorf("got %q, want %q", got, "\x1b[38;2;255;0;0m")
	}
}

func TestColorToANSI_RGB(t *testing.T) {
	got, ok := render.ColorToANSI("rgb(0,255,0)")
	if !ok {
		t.Fatal("ColorToANSI(rgb(0,255,0)): ok=false")
	}
	if got != "\x1b[38;2;0;255;0m" {
		t.Errorf("got %q, want %q", got, "\x1b[38;2;0;255;0m")
	}
}

func TestColorToANSI_Invalid(t *testing.T) {
	_, ok := render.ColorToANSI("notacolor")
	if ok {
		t.Error("expected ok=false for invalid color")
	}
}

// --- WrapWithColor (task06) ---

func TestWrapWithColor_ShortCircuit(t *testing.T) {
	if got := render.WrapWithColor("hi", ""); got != "hi" {
		t.Errorf("expected short-circuit on empty ansi, got %q", got)
	}
	if got := render.WrapWithColor("", "\x1b[31m"); got != "" {
		t.Errorf("expected short-circuit on empty s, got %q", got)
	}
}

func TestWrapWithColor_Wrap(t *testing.T) {
	got := render.WrapWithColor("hello", "\x1b[31m")
	want := "\x1b[31m" + "hello" + "\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- ShouldRenderGopher ---

func TestShouldRenderGopher_EmptyText(t *testing.T) {
	if render.ShouldRenderGopher("", []string{""}) {
		t.Error("empty text should not trigger gopher")
	}
}

func TestShouldRenderGopher_NonASCII(t *testing.T) {
	if !render.ShouldRenderGopher("こんにちは", []string{""}) {
		t.Error("non-ASCII text should trigger gopher")
	}
}

func TestShouldRenderGopher_ASCII(t *testing.T) {
	if render.ShouldRenderGopher("Hello", []string{"Hello"}) {
		t.Error("ASCII text should not trigger gopher")
	}
}

// --- Alignment (task08) ---

func buildSimpleMap() map[rune][]string {
	m := make(map[rune][]string)
	for r := rune(32); r <= 126; r++ {
		art := make([]string, 8)
		// Use single-char wide glyphs for predictable width.
		art[4] = string(r)
		for i := range art {
			if art[i] == "" {
				art[i] = " "
			}
		}
		m[r] = art
	}
	return m
}

func TestRenderAligned_Left(t *testing.T) {
	var buf bytes.Buffer
	m := buildSimpleMap()
	segs := []string{"AB"}
	render.RenderAlignedWithColorRules(&buf, segs, m, nil, "left", 80)
	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[4], "AB") {
		t.Errorf("left-aligned line[4] should start with \"AB\", got %q", lines[4])
	}
}

func TestRenderAligned_Right(t *testing.T) {
	var buf bytes.Buffer
	m := buildSimpleMap()
	segs := []string{"A"}
	render.RenderAlignedWithColorRules(&buf, segs, m, nil, "right", 20)
	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 lines, got %d", len(lines))
	}
	if !strings.HasSuffix(strings.TrimRight(lines[4], " "), "A") {
		t.Errorf("right-aligned line[4] should end with \"A\", got %q", lines[4])
	}
	if !strings.HasPrefix(lines[4], " ") {
		t.Errorf("right-aligned line[4] should have leading spaces, got %q", lines[4])
	}
}

func TestRenderAligned_Center(t *testing.T) {
	var buf bytes.Buffer
	m := buildSimpleMap()
	segs := []string{"A"}
	render.RenderAlignedWithColorRules(&buf, segs, m, nil, "center", 21)
	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[4], " ") {
		t.Errorf("center-aligned line[4] should have leading spaces, got %q", lines[4])
	}
}

func TestRenderAligned_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	m := buildSimpleMap()
	segs := []string{""}
	render.RenderAlignedWithColorRules(&buf, segs, m, nil, "left", 80)
	if buf.Len() != 0 {
		t.Errorf("empty input should produce no output, got %q", buf.String())
	}
}

func TestRenderAligned_NewlineOnly(t *testing.T) {
	var buf bytes.Buffer
	m := buildSimpleMap()
	segs := []string{"", ""}
	render.RenderAlignedWithColorRules(&buf, segs, m, nil, "left", 80)
	// Two all-empty segments → 1 blank line (len(segs)-1).
	if buf.String() != "\n" {
		t.Errorf("two empty segments should yield one blank line, got %q", buf.String())
	}
}
