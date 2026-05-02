package render_test

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"ascii-art-reverse/internal/render"
)

// --- ParseInput (task05) ---

func TestParseInput(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "NewlineSplit",
			input: "Hello\nWorld",
			want:  []string{"Hello", "World"},
		},
		{
			name:  "Empty",
			input: "",
			want:  []string{""},
		},
		{
			name:  "OnlyNewline",
			input: "\\n",
			want:  []string{"", ""},
		},
		{
			name:  "MultipleNewlines",
			input: "A\\nB\\nC",
			want:  []string{"A", "B", "C"},
		},
		{
			name:  "FiltersNonASCII",
			input: "héllo",
			want:  []string{"hllo"},
		},
		{
			name:  "TabExpanded",
			input: "A\\tB",
			want:  []string{"A   B"},
		},
		{
			name:  "BareCRSplits",
			input: "A\rB",
			want:  []string{"A", "B"},
		},
		{
			name:  "CRLFNormalized",
			input: "A\r\nB",
			want:  []string{"A", "B"},
		},
		{
			name:  "MultipleCRLF",
			input: "A\r\nB\r\nC",
			want:  []string{"A", "B", "C"},
		},
		{
			name:  "TrailingBackslash",
			input: `\`,
			want:  []string{`\`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			got := render.ParseInput(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				st.Errorf("ParseInput() = %v, want %v", got, tt.want)
			}
			if len(got) != len(tt.want) {
				st.Fatalf("got %d segments, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					st.Errorf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// --- Mapper (task05) ---

func TestMapper(t *testing.T) {
	t.Parallel()
	bannerMap := map[rune][]string{
		'A': {"l1", "l2", "l3", "l4", "l5", "l6", "l7", "l8"},
	}
	get := render.Mapper(bannerMap)

	tests := []struct {
		wantVal string
		name    string
		input   rune
		wantIdx int
	}{
		{name: "ValidRune", input: 'A', wantIdx: 0, wantVal: "l1"},
		{name: "MissingRune", input: 'Z', wantIdx: 0, wantVal: ""},
		{name: "OutOfRangeBelow", input: '\x1f', wantIdx: 0, wantVal: ""},
		{name: "OutOfRangeAbove", input: '\x7f', wantIdx: 0, wantVal: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			got := get(tt.input)
			if len(got) != 8 {
				st.Fatalf("expected 8 lines, got %d", len(got))
			}
			if got[tt.wantIdx] != tt.wantVal {
				st.Errorf("line[%d] = %q, want %q", tt.wantIdx, got[tt.wantIdx], tt.wantVal)
			}
		})
	}
}

// --- ColorToANSI (task06) ---

func TestColorToANSI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  string
		want   string
		wantOK bool
	}{
		{"NamedRed", "red", "\x1b[31m", true},
		{"NamedGreen", "green", "\x1b[32m", true},
		{"NamedBlue", "blue", "\x1b[34m", true},
		{"NamedCyan", "cyan", "\x1b[36m", true},
		{"NamedWhite", "white", "\x1b[37m", true},
		{"NamedBlack", "black", "\x1b[30m", true},
		{"HexRed", "#FF0000", "\x1b[38;2;255;0;0m", true},
		{"RGBGreen", "rgb(0,255,0)", "\x1b[38;2;0;255;0m", true},
		{"Invalid", "notacolor", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			got, ok := render.ColorToANSI(tt.input)
			if ok != tt.wantOK {
				st.Fatalf("ColorToANSI(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
			if tt.wantOK && got != tt.want {
				st.Errorf("ColorToANSI(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// --- StripANSI ---

func TestStripANSI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "NoANSI", input: "hello", want: "hello"},
		{name: "SimpleRed", input: "\x1b[31mhello\x1b[0m", want: "hello"},
		{name: "MultipleSequences", input: "\x1b[31mH\x1b[0m\x1b[32me\x1b[0mll\x1b[34mo\x1b[0m", want: "Hello"},
		{name: "Empty", input: "", want: ""},
		{name: "OnlyANSI", input: "\x1b[31m\x1b[0m", want: ""},
		{name: "Incomplete", input: "hi\x1b[31", want: "hi\x1b[31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			if got := render.StripANSI(tt.input); got != tt.want {
				st.Errorf("StripANSI() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- WrapWithColor (task06) ---

func TestWrapWithColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		ansiStart string
		want      string
	}{
		{
			name:      "EmptyAnsiShortCircuit",
			input:     "hi",
			ansiStart: "",
			want:      "hi",
		},
		{
			name:      "EmptyStringShortCircuit",
			input:     "",
			ansiStart: "\x1b[31m",
			want:      "",
		},
		{
			name:      "WrapRed",
			input:     "hello",
			ansiStart: "\x1b[31m",
			want:      "\x1b[31mhello\x1b[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			got := render.WrapWithColor(tt.input, tt.ansiStart)
			if got != tt.want {
				st.Errorf("WrapWithColor() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- ShouldRenderGopher ---

func TestShouldRenderGopher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		text     string
		segments []string
		want     bool
	}{
		{"EmptyText", "", []string{""}, false},
		{"NonASCII", "こんにちは", []string{""}, true},
		{"ASCII", "Hello", []string{"Hello"}, false},
		{
			name:     "NonASCII_But_NotBlankSegments",
			text:     "こんにちはA",
			segments: []string{"A"},
			want:     false,
		},
		{
			name:     "AllowedEscapes_Only",
			text:     "\\n",
			segments: []string{"", ""},
			want:     false,
		},
		{
			name:     "TabEscape_Only",
			text:     "\\t",
			segments: []string{""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			if got := render.ShouldRenderGopher(tt.text, tt.segments); got != tt.want {
				st.Errorf("ShouldRenderGopher(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
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

func TestRenderAligned(t *testing.T) {
	m := buildSimpleMap()
	tests := []struct {
		name  string
		align string
		check func(st *testing.T, out string)
		segs  []string
		width int
	}{
		{
			name:  "Left",
			segs:  []string{"AB"},
			align: "left",
			width: 80,
			check: func(st *testing.T, out string) {
				lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
				if len(lines) != 8 {
					st.Fatalf("expected 8 lines, got %d", len(lines))
				}
				if !strings.HasPrefix(lines[4], "AB") {
					st.Errorf("expected AB prefix, got %q", lines[4])
				}
			},
		},
		{
			name:  "Right",
			segs:  []string{"A"},
			align: "right",
			width: 20,
			check: func(st *testing.T, out string) {
				lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
				if !strings.HasSuffix(strings.TrimRight(lines[4], " "), "A") {
					st.Errorf("expected A suffix, got %q", lines[4])
				}
			},
		},
		{
			name:  "Center",
			segs:  []string{"A"},
			align: "center",
			width: 21,
			check: func(st *testing.T, out string) {
				lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
				if !strings.HasPrefix(lines[4], " ") {
					st.Errorf("expected leading space, got %q", lines[4])
				}
			},
		},
		{
			name:  "EmptyInput",
			segs:  []string{""},
			align: "left",
			width: 80,
			check: func(st *testing.T, out string) {
				if out != "" {
					st.Errorf("expected empty output, got %q", out)
				}
			},
		},
		{
			name:  "NewlineOnly",
			segs:  []string{"", ""},
			align: "left",
			width: 80,
			check: func(st *testing.T, out string) {
				if out != "\n" {
					st.Errorf("expected one newline, got %q", out)
				}
			},
		},
		{
			name:  "JustifyWords", // This was the problematic one, now fixed with named fields
			segs:  []string{"A B"},
			align: "justify",
			width: 22,
			check: func(st *testing.T, out string) {
				lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
				if len(lines[4]) != 22 {
					st.Errorf("Justify width mismatch: got %d, want 22", len(lines[4]))
				}
			},
		},
		{
			name:  "JustifyNarrowClamp",
			segs:  []string{"A B"},
			align: "justify",
			width: 1,
			check: func(st *testing.T, out string) {
				lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
				if len(lines) != 8 {
					st.Errorf("expected 8 lines for narrow justify, got %d", len(lines))
				}
			},
		},
		{
			name:  "JustifySpacesOnly",
			segs:  []string{"  "},
			align: "justify",
			width: 20,
			check: func(st *testing.T, out string) {
				if out == "" {
					st.Error("expected non-empty output for space-only justify")
				}
			},
		},
		{
			name:  "UnknownAlignmentFallback",
			segs:  []string{"A"},
			align: "unknown",
			width: 20,
			check: func(st *testing.T, out string) {
				if out == "" {
					st.Error("expected output for unknown alignment")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			var buf bytes.Buffer
			render.RenderAlignedWithColorRules(&buf, tt.segs, m, nil, tt.align, tt.width)
			tt.check(st, buf.String())
		})
	}
}

// TestRenderConvenienceWrappers verifies every delegating entry point produces output.
func TestRenderConvenienceWrappers(t *testing.T) {
	m := buildSimpleMap()
	segs := []string{"A"}

	tests := []struct {
		name string
		fn   func() int
	}{
		{"Render", func() int {
			var buf bytes.Buffer
			render.Render(&buf, segs, m)
			return buf.Len()
		}},
		{"RenderWithColor", func() int {
			var buf bytes.Buffer
			render.RenderWithColor(&buf, segs, m, "\x1b[31m", "")
			return buf.Len()
		}},
		{"RenderWithColorRules", func() int {
			var buf bytes.Buffer
			render.RenderWithColorRules(&buf, segs, m, []render.ColorRule{{ANSIStart: "\x1b[31m"}})
			return buf.Len()
		}},
		{"RenderAligned", func() int {
			var buf bytes.Buffer
			render.RenderAligned(&buf, segs, m, "left", 80)
			return buf.Len()
		}},
		{"RenderAlignedWithColor", func() int {
			var buf bytes.Buffer
			render.RenderAlignedWithColor(&buf, segs, m, "\x1b[31m", "", "left", 80)
			return buf.Len()
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			if n := tt.fn(); n == 0 {
				st.Errorf("%s: expected non-empty output", tt.name)
			}
		})
	}
}

// TestRenderGopher verifies RenderGopher writes non-empty ANSI art.
func TestRenderGopher(t *testing.T) {
	tests := []struct {
		name          string
		renderFn      func(io.Writer)
		wantSubstring string
	}{
		{
			name:          "PeekingGopher",
			renderFn:      render.RenderGopher,
			wantSubstring: "Gopher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			var buf bytes.Buffer
			tt.renderFn(&buf)
			if buf.Len() == 0 {
				st.Errorf("%s: expected non-empty output", tt.name)
			}
			if !strings.Contains(buf.String(), tt.wantSubstring) {
				st.Errorf("%s: expected %q in output", tt.name, tt.wantSubstring)
			}
		})
	}
}

// TestDetectTerminalWidth covers the COLUMNS env path and fallback.
func TestDetectTerminalWidth(t *testing.T) {
	tests := []struct {
		name    string
		columns string
		wantMin int // result must be >= wantMin (fallback is ≥80)
		wantMax int // result must be <= wantMax (0 = no upper bound check)
		wantVal int // exact value when > 0
	}{
		{
			name:    "valid COLUMNS env",
			columns: "120",
			wantVal: 120,
		},
		{
			name:    "invalid COLUMNS falls back to positive width",
			columns: "bad",
			wantMin: 1,
		},
		{
			name:    "zero COLUMNS falls back to positive width",
			columns: "0",
			wantMin: 1,
		},
		{
			name:    "negative COLUMNS falls back to positive width",
			columns: "-5",
			wantMin: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Setenv("COLUMNS", tt.columns)
			got := render.DetectTerminalWidth()
			if tt.wantVal > 0 && got != tt.wantVal {
				st.Errorf("DetectTerminalWidth() = %d, want %d", got, tt.wantVal)
			}
			if tt.wantMin > 0 && got < tt.wantMin {
				st.Errorf("DetectTerminalWidth() = %d, want >= %d", got, tt.wantMin)
			}
		})
	}
}

// TestResolveWidth covers file-output (always 80) and terminal paths.
func TestResolveWidth(t *testing.T) {
	tests := []struct {
		name     string
		fileMode bool
		want     int
	}{
		{
			name:     "file output always 80",
			fileMode: true,
			want:     80,
		},
		{
			name:     "terminal mode uses env",
			fileMode: false,
			want:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			if !tt.fileMode {
				st.Setenv("COLUMNS", "100")
			}
			if got := render.ResolveWidth(tt.fileMode); got != tt.want {
				st.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

// --- ANSI and Styling (task06) ---

// TestColorSubstringHighlight exercises the substring-matching path in buildColorMask.
func TestColorSubstringHighlight(t *testing.T) {
	m := buildSimpleMap()
	segs := []string{"A"}
	rulesAll := []render.ColorRule{{ANSIStart: "\x1b[31m", Substring: ""}}

	tests := []struct {
		rules    []render.ColorRule
		align    string
		name     string
		segs     []string
		width    int
		wantANSI bool
	}{
		{
			name: "substring match injects ANSI",
			segs: []string{"Hello World"},
			rules: []render.ColorRule{
				{ANSIStart: "\x1b[31m", Substring: "Hello"},
			},
			wantANSI: true,
			align:    "left",
			width:    0,
		},
		{
			name: "non-matching substring produces no ANSI",
			segs: []string{"Hi"},
			rules: []render.ColorRule{
				{ANSIStart: "\x1b[31m", Substring: "ZZZZ"},
			},
			wantANSI: false,
			align:    "left",
			width:    0,
		},
		{
			name: "empty substring colors whole string",
			segs: []string{"Hi"},
			rules: []render.ColorRule{
				{ANSIStart: "\x1b[31m", Substring: ""},
			},
			wantANSI: true,
			align:    "left",
			width:    0,
		},
		{
			name:     "VisibleWidth_ANSI_Stripping_Right",
			segs:     segs,
			rules:    rulesAll,
			wantANSI: true,
			align:    "right",
			width:    20,
		},
		{
			name:     "VisibleWidth_ANSI_Stripping_Center",
			segs:     segs,
			rules:    rulesAll,
			wantANSI: true,
			align:    "center",
			width:    20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			var buf bytes.Buffer
			render.RenderAlignedWithColorRules(&buf, tt.segs, m, tt.rules, tt.align, tt.width)
			hasANSI := strings.Contains(buf.String(), "\x1b[")
			if hasANSI != tt.wantANSI {
				st.Errorf("wantANSI=%v but output ANSI presence=%v", tt.wantANSI, hasANSI)
			}
			if tt.width > 0 && tt.wantANSI {
				lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
				if !strings.HasPrefix(lines[4], " ") {
					st.Errorf("%s: expected leading space, got %q", tt.name, lines[4])
				}
			}
		})
	}
}

// FuzzParseInput checks that ParseInput never panics on arbitrary string input.
func FuzzParseInput(f *testing.F) {
	f.Add("Hello World")
	f.Add("")
	f.Add("\\n")
	f.Add("A\\nB\\nC")
	f.Add("héllo")
	f.Add("!@#$%^&*()")
	f.Add("   spaces   ")
	f.Add("Hello\nWorld")
	f.Add("\x00\x01\x7f")

	f.Fuzz(func(t *testing.T, input string) {
		segs := render.ParseInput(input)
		if segs == nil {
			t.Error("ParseInput must never return nil")
		}
	})
}

// FuzzColorToANSI checks that ColorToANSI never panics on arbitrary input.
func FuzzColorToANSI(f *testing.F) {
	f.Add("red")
	f.Add("green")
	f.Add("#FF0000")
	f.Add("#abc")
	f.Add("rgb(0,255,0)")
	f.Add("rgb(999,0,0)")
	f.Add("hsl(120,100%,50%)")
	f.Add("")
	f.Add("notacolor")
	f.Add("#GGGGGG")

	f.Fuzz(func(t *testing.T, color string) {
		// Must not panic; ok=false is a valid result for unknown colors.
		_, _ = render.ColorToANSI(color)
	})
}

func BenchmarkParseInput(b *testing.B) {
	cases := []struct{ name, input string }{
		{"short", "Hello"},
		{"with_escape", `Hello\nWorld`},
		{"long", "Hello World How Are You Doing Today Fine Thanks"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = render.ParseInput(tc.input)
			}
		})
	}
}

func BenchmarkColorToANSI(b *testing.B) {
	cases := []struct{ name, input string }{
		{"named", "red"},
		{"hex", "#FF0000"},
		{"rgb", "rgb(0,255,0)"},
		{"hsl", "hsl(120,100%,50%)"},
		{"invalid", "notacolor"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = render.ColorToANSI(tc.input)
			}
		})
	}
}

func BenchmarkRenderAlignedWithColorRules(b *testing.B) {
	m := buildSimpleMap()
	red := []render.ColorRule{{ANSIStart: "\x1b[31m", Substring: "Hello"}}
	cases := []struct {
		segs  []string
		rules []render.ColorRule
		name  string
		align string
		width int
	}{
		{name: "short_left", segs: []string{"Hello"}, rules: nil, align: "left", width: 80},
		{name: "long_left", segs: []string{"Hello World How Are You"}, rules: nil, align: "left", width: 80},
		{name: "with_color", segs: []string{"Hello World"}, rules: red, align: "left", width: 80},
		{name: "right_align", segs: []string{"Hello"}, rules: nil, align: "right", width: 80},
		{name: "center_align", segs: []string{"Hello"}, rules: nil, align: "center", width: 80},
		{name: "justify", segs: []string{"Hello World"}, rules: nil, align: "justify", width: 80},
		{name: "multi_segment", segs: []string{"Hello", "World"}, rules: nil, align: "left", width: 80},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				render.RenderAlignedWithColorRules(io.Discard, tc.segs, m, tc.rules, tc.align, tc.width)
			}
		})
	}
}
