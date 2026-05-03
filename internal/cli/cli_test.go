package cli_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"ascii-art-reverse/internal/cli"
)

func TestClassifyArgs_Fundamentals(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantOut string
		args    []string
		wantPos []string
		wantCol int // count of color rules
		wantUnk int // count of unknown flags
	}{
		{
			name:    "BasicPositional",
			args:    []string{"Hello", "standard"},
			wantPos: []string{"Hello", "standard"},
			wantUnk: 0,
		},
		{
			name:    "UnknownFlag",
			args:    []string{"--unknown", "Hello"},
			wantPos: []string{"Hello"},
			wantUnk: 1,
		},
		{
			name:    "double dash as literal is unknown (triggers hint)",
			args:    []string{"--"},
			wantPos: []string{},
			wantUnk: 1,
		},
		{
			name:    "double dash as delimiter with flags",
			args:    []string{"--color=blue", "--", "--"},
			wantPos: []string{"--"}, // Protected literal "--"
			wantCol: 1,
		},
		{
			name:    "double dash stops flag processing",
			args:    []string{"--", "--output=test.txt"},
			wantPos: []string{"--output=test.txt"},
			wantOut: "",
		},
		{
			name:    "double dash ignored as delimiter if next arg is not a flag",
			args:    []string{"--", "Hello", "standard"},
			wantPos: []string{"Hello", "standard"},
			wantUnk: 1, // "--" is unknown because "Hello" is not a flag
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()
			got := cli.ClassifyArgs(tt.args)

			if !reflect.DeepEqual(got.Positional, tt.wantPos) {
				st.Errorf("Positional = %v, want %v", got.Positional, tt.wantPos)
			}
			if got.OutputValue != tt.wantOut {
				st.Errorf("OutputValue = %q, want %q", got.OutputValue, tt.wantOut)
			}
			if len(got.ColorRules) != tt.wantCol {
				st.Errorf("len(ColorRules) = %d, want %d", len(got.ColorRules), tt.wantCol)
			}
			if len(got.UnknownFlags) != tt.wantUnk {
				st.Errorf("len(UnknownFlags) = %d, want %d", len(got.UnknownFlags), tt.wantUnk)
			}
		})
	}
}

func ExampleClassifyArgs() {
	args := []string{"--color=red", "Hello", "standard"}
	parsed := cli.ClassifyArgs(args)
	fmt.Println(parsed.Positional)
	fmt.Println(parsed.ColorRules[0].ColorValue)
	// Output:
	// [Hello standard]
	// red
}

func TestSelectUsage(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		osArgs []string
	}{
		{
			name:   "Basic Usage",
			osArgs: []string{"."},
			want:   cli.UsageBasic,
		},
		{
			name:   "Reverse Priority over All",
			osArgs: []string{".", "--reverse=file.txt", "--output=out.txt"},
			want:   cli.UsageReverse,
		},
		{
			name:   "Output Priority over Align",
			osArgs: []string{".", "--output=out.txt", "--align=center"},
			want:   cli.UsageOutput,
		},
		{
			name:   "Align Priority over Color",
			osArgs: []string{".", "--align=center", "--color=blue"},
			want:   cli.UsageAlign,
		},
		{
			name:   "Color Usage",
			osArgs: []string{".", "--color=green"},
			want:   cli.UsageColor,
		},
		{
			name:   "Malformed Output triggers Output Usage",
			osArgs: []string{".", "--output="},
			want:   cli.UsageOutput,
		},
		{
			name:   "Unknown Flag defaults to Basic Usage",
			osArgs: []string{".", "--unknown"},
			want:   cli.UsageBasic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			oldArgs := os.Args
			os.Args = tt.osArgs
			defer func() { os.Args = oldArgs }()

			if got := cli.SelectUsage(cli.ParsedArgs{}); got != tt.want {
				st.Errorf("SelectUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmitWarnings(t *testing.T) {
	tests := []struct {
		name           string
		wantSubstrings []string
		args           cli.ParsedArgs
	}{
		{
			name: "MixedWarnings",
			args: cli.ParsedArgs{
				OutputValue: "new.txt",
				Malformed: []cli.FlagError{
					{Category: "dup-output", Raw: "--output=old.txt"},
					{Category: "color", Raw: "invalid"},
				},
				UnknownFlags: []cli.FlagError{{Category: "unknown", Raw: "---"}},
			},
			wantSubstrings: []string{
				`warning: --output overridden: old.txt → new.txt`,
				`warning: invalid color flag "invalid"`,
				`hint: use "--" before [STRING] that look like flags (e.g., go run . -- ---)`,
			},
		},
		{
			name: "output then align produces one line in input order",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "output", Raw: "--output=bad.pdf"},
					{Category: "align", Raw: "--align=middle"},
				},
			},
			wantSubstrings: []string{`warning: invalid output, align flags "--output=bad.pdf", "--align=middle"`},
		},
		{
			name: "align then color preserves input order",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "align", Raw: "--align=middle"},
					{Category: "color", Raw: "--color"},
				},
			},
			wantSubstrings: []string{`warning: invalid align, color flags "--align=middle", "--color"`},
		},
		{
			name: "reverse then output then align in one line",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "reverse", Raw: "--reverse"},
					{Category: "output", Raw: "--output=bad.pdf"},
					{Category: "align", Raw: "--align=x"},
				},
			},
			wantSubstrings: []string{`warning: invalid reverse, output, align flags "--reverse", "--output=bad.pdf", "--align=x"`},
		},
		{
			name: "duplicate align flag",
			args: cli.ParsedArgs{
				AlignValue: "right",
				Malformed:  []cli.FlagError{{Category: "dup-align", Raw: "left"}},
			},
			wantSubstrings: []string{`--align overridden: left → right`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			cli.EmitWarnings(tt.args)

			_ = w.Close()
			os.Stderr = old

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			for _, sub := range tt.wantSubstrings {
				if !strings.Contains(output, sub) {
					st.Errorf("EmitWarnings() missing expected substring: %q\nGot: %q", sub, output)
				}
			}
		})
	}
}

func TestClassifyArgs_MixedFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args           []string
		wantPos        []string
		name           string
		wantOutput     string
		wantReverse    string
		wantAlign      string
		wantColorSub   string
		wantColorCount int
		wantMalformed  int
		wantUnknown    int
	}{
		{
			name:           "valid output and malformed color",
			args:           []string{"--output=test.txt", "--color red", "hello"},
			wantOutput:     "test.txt",
			wantReverse:    "",
			wantPos:        []string{"hello"},
			wantColorCount: 0,
			wantMalformed:  1,
		},
		{
			name:           "invalid align type and valid color",
			args:           []string{"--align=middle", "--color=blue", "hello"},
			wantAlign:      "",
			wantReverse:    "",
			wantPos:        []string{"hello"},
			wantColorCount: 1,
			wantMalformed:  1,
		},
		{
			name:          "multiple malformed flags with positional",
			args:          []string{"--output=test.pdf", "--align top", "standard"},
			wantOutput:    "", // Malformed, so not set
			wantReverse:   "",
			wantAlign:     "", // Malformed, so not set
			wantPos:       []string{"standard"},
			wantMalformed: 2, // Two malformed flags
		},
		{
			name:          "malformed output extension",
			args:          []string{"--output=image.jpg", "text"},
			wantOutput:    "",
			wantReverse:   "",
			wantAlign:     "",
			wantPos:       []string{"text"},
			wantMalformed: 1,
		},
		{
			name:          "empty output value",
			args:          []string{"--output=", "text"},
			wantOutput:    "",
			wantReverse:   "",
			wantAlign:     "",
			wantPos:       []string{"text"},
			wantMalformed: 1,
		},
		{
			name:          "malformed reverse flag",
			args:          []string{"--reverse", "file.txt"},
			wantOutput:    "",
			wantReverse:   "", // Malformed, so not set
			wantAlign:     "",
			wantPos:       []string{"file.txt"}, // "file.txt" is positional after malformed flag
			wantMalformed: 1,
		},
		{
			name:          "malformed reverse extension",
			args:          []string{"--reverse=example07", "text"},
			wantOutput:    "",
			wantReverse:   "",
			wantAlign:     "",
			wantPos:       []string{"text"},
			wantMalformed: 1,
			wantUnknown:   0,
		},
		{
			name:           "duplicate output flag (last wins)",
			args:           []string{"--output=first.txt", "--output=second.txt", "hello"},
			wantOutput:     "second.txt",
			wantReverse:    "",
			wantAlign:      "",
			wantPos:        []string{"hello"},
			wantColorCount: 0,
			wantMalformed:  1, // For the "dup-output" warning
		},
		{
			name:           "duplicate align flag (last wins)",
			args:           []string{"--align=left", "--align=right", "hello"},
			wantOutput:     "",
			wantReverse:    "",
			wantAlign:      "right",
			wantPos:        []string{"hello"},
			wantColorCount: 0,
			wantMalformed:  1, // For the "dup-align" warning
		},
		{
			name:          "duplicate reverse flag",
			args:          []string{"--reverse=first.txt", "--reverse=second.txt"},
			wantReverse:   "second.txt",
			wantPos:       []string{},
			wantMalformed: 1,
		},
		{
			name:        "reverse with specific banner",
			args:        []string{"--reverse=test.txt", "thinkertoy"},
			wantReverse: "test.txt",
			wantPos:     []string{"thinkertoy"},
			wantUnknown: 0,
		},
		{
			name:           "stop flags delimiter in color parse",
			args:           []string{"--color=red", "sub", "--", "--output=x", "text"},
			wantPos:        []string{"sub", "--output=x", "text"},
			wantColorCount: 1,
			wantColorSub:   "",
			wantUnknown:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()
			got := cli.ClassifyArgs(tt.args)

			if got.OutputValue != tt.wantOutput {
				st.Errorf("%s: OutputValue = %q, want %q", tt.name, got.OutputValue, tt.wantOutput)
			}
			if got.AlignValue != tt.wantAlign {
				st.Errorf("%s: AlignValue = %q, want %q", tt.name, got.AlignValue, tt.wantAlign)
			}
			if got.ReverseValue != tt.wantReverse {
				st.Errorf("%s: ReverseValue = %q, want %q", tt.name, got.ReverseValue, tt.wantReverse)
			}
			if !reflect.DeepEqual(got.Positional, tt.wantPos) {
				st.Errorf("%s: Positional = %v, want %v", tt.name, got.Positional, tt.wantPos)
			}
			if len(got.ColorRules) != tt.wantColorCount {
				st.Errorf("%s: ColorRules count = %d, want %d", tt.name, len(got.ColorRules), tt.wantColorCount)
			}
			if tt.wantColorSub != "" && len(got.ColorRules) > 0 && got.ColorRules[0].Substring != tt.wantColorSub {
				st.Errorf("%s: ColorSub = %q, want %q", tt.name, got.ColorRules[0].Substring, tt.wantColorSub)
			}
			if len(got.Malformed) != tt.wantMalformed {
				st.Errorf("%s: Malformed count = %d, want %d", tt.name, len(got.Malformed), tt.wantMalformed)
			}
			if len(got.UnknownFlags) != tt.wantUnknown {
				st.Errorf("%s: Unknown count = %d, want %d", tt.name, len(got.UnknownFlags), tt.wantUnknown)
			}
		})
	}
}

// TestWarnFunctions covers every Warn* helper that writes to stderr.
func TestWarnFunctions(t *testing.T) {
	tests := []struct {
		name    string
		fn      func()
		wantSub string
	}{
		{
			name:    "WarnAlignOverridden",
			fn:      func() { cli.WarnAlignOverridden("left", "right") },
			wantSub: `--align overridden: left → right`,
		},
		{
			name:    "WarnBannerNotFound",
			fn:      func() { cli.WarnBannerNotFound("mytheme") },
			wantSub: `banner "mytheme" not found`,
		},
		{
			name:    "WarnBannerInvalid",
			fn:      func() { cli.WarnBannerInvalid("broken") },
			wantSub: `invalid banner "broken"`,
		},
		{
			name:    "WarnFlagsAfterString",
			fn:      func() { cli.WarnFlagsAfterString("--color=red") },
			wantSub: `--color=red`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			r, w, _ := os.Pipe()
			old := os.Stderr
			os.Stderr = w
			tt.fn()
			_ = w.Close()
			os.Stderr = old

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			if !strings.Contains(buf.String(), tt.wantSub) {
				st.Errorf("output missing %q\nGot: %q", tt.wantSub, buf.String())
			}
		})
	}
}

// TestIsKnownFlagOption covers every recognized prefix and unknown inputs.
func TestIsKnownFlagOption(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  bool
	}{
		{"--color=red", true},
		{"--color", true},
		{"--output=file.txt", true},
		{"--align=center", true},
		{"--reverse=art.txt", true},
		{"--stdin", false},
		{"hello", false},
		{"", false},
		{"--unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(st *testing.T) {
			st.Parallel()
			if got := cli.IsKnownFlagOption(tt.input); got != tt.want {
				st.Errorf("IsKnownFlagOption(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStdinHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		fn    func(args cli.ParsedArgs) string
		args  cli.ParsedArgs
	}{
		{
			name:  "ReadStdin",
			input: "hello from stdin\n",
			fn:    func(_ cli.ParsedArgs) string { return cli.ReadStdin() },
		},
		{
			name:  "ResolveInput_Stdin",
			input: "piped input\n",
			fn:    cli.ResolveInput,
			args:  cli.ParsedArgs{StdinMode: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			r, w, _ := os.Pipe()
			old := os.Stdin
			os.Stdin = r
			st.Cleanup(func() { os.Stdin = old })

			_, _ = w.WriteString(tt.input)
			_ = w.Close()

			got := tt.fn(tt.args)
			want := strings.TrimRight(tt.input, "\n")
			if got != want {
				st.Errorf("%s = %q, want %q", tt.name, got, want)
			}
		})
	}
}

func TestBuildColorRules(t *testing.T) {
	tests := []struct {
		name      string
		raw       []cli.RawColorRule
		wantCount int
	}{
		{
			name:      "valid named color",
			raw:       []cli.RawColorRule{{ColorValue: "red"}},
			wantCount: 1,
		},
		{
			name:      "valid hex color with substring",
			raw:       []cli.RawColorRule{{ColorValue: "#00FF00", Substring: "Hello"}},
			wantCount: 1,
		},
		{
			name:      "invalid color is dropped",
			raw:       []cli.RawColorRule{{ColorValue: "notacolor"}},
			wantCount: 0,
		},
		{
			name: "mix of valid and invalid",
			raw: []cli.RawColorRule{
				{ColorValue: "blue"},
				{ColorValue: "bad"},
				{ColorValue: "green"},
			},
			wantCount: 2,
		},
		{
			name:      "empty list",
			raw:       []cli.RawColorRule{},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			// Suppress any warning output to stderr.
			r, w, _ := os.Pipe()
			old := os.Stderr
			os.Stderr = w
			got := cli.BuildColorRules(tt.raw)
			_ = w.Close()
			os.Stderr = old
			_, _ = io.Copy(io.Discard, r)

			if len(got) != tt.wantCount {
				st.Errorf("BuildColorRules() count = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}

// TestResolveInput covers positional, empty, and stdin paths.
func TestResolveInput(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args cli.ParsedArgs
		want string
	}{
		{
			name: "first positional arg",
			args: cli.ParsedArgs{Positional: []string{"Hello", "standard"}},
			want: "Hello",
		},
		{
			name: "no positionals",
			args: cli.ParsedArgs{Positional: []string{}},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()
			if got := cli.ResolveInput(tt.args); got != tt.want {
				st.Errorf("ResolveInput() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestClassifyArgs_ColorWithSubstring covers parseColor consuming the next token.
func TestClassifyArgs_ColorWithSubstring(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args         []string
		name         string
		wantSub      string
		wantColorLen int
	}{
		{
			name:         "color grabs next positional as substring",
			args:         []string{"--color=red", "Hello", "Hello World"},
			wantSub:      "Hello",
			wantColorLen: 1,
		},
		{
			name:         "color with no following positional",
			args:         []string{"--color=blue", "Text"},
			wantSub:      "",
			wantColorLen: 1,
		},
		{
			name:         "color flag last arg has no substring",
			args:         []string{"Text", "--color=green"},
			wantSub:      "",
			wantColorLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()
			got := cli.ClassifyArgs(tt.args)
			if len(got.ColorRules) != tt.wantColorLen {
				st.Fatalf("ColorRules count = %d, want %d", len(got.ColorRules), tt.wantColorLen)
			}
			if len(got.ColorRules) > 0 && got.ColorRules[0].Substring != tt.wantSub {
				st.Errorf("Substring = %q, want %q", got.ColorRules[0].Substring, tt.wantSub)
			}
		})
	}
}

// TestGuardFunctions covers the non-os.Exit branches of the validation guards.
func TestGuardFunctions(t *testing.T) {
	tests := []struct {
		fn   func()
		name string
	}{
		{
			name: "ValidateOrFatal_Clean",
			fn: func() {
				cli.ValidateOrFatal(cli.ParsedArgs{
					Positional: []string{"Hi"}, ColorRules: []cli.RawColorRule{},
					UnknownFlags: []cli.FlagError{}, Malformed: []cli.FlagError{},
				})
			},
		},
		{
			name: "ValidateOrFatal_DupOutput",
			fn: func() {
				cli.ValidateOrFatal(cli.ParsedArgs{
					Positional: []string{"hello"}, OutputValue: "new.txt",
					Malformed: []cli.FlagError{{Category: "dup-output", Raw: "--output=old.txt"}},
				})
			},
		},
		{
			name: "CheckPositionalsOrFatal_Two",
			fn: func() {
				cli.CheckPositionalsOrFatal(cli.ParsedArgs{
					Positional: []string{"Hello", "standard"},
				})
			},
		},
		{
			name: "HandleBannerFallbackOrFatal_Nil",
			fn: func() {
				cli.HandleBannerFallbackOrFatal(nil, "standard")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.fn() // must not call os.Exit
		})
	}
}

func TestClassifyArgs_StdinMode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args        []string
		wantPos     []string
		name        string
		wantOutput  string
		wantAlign   string
		wantReverse string
		wantStdin   bool
	}{
		{
			name:      "stdin only",
			args:      []string{"--stdin"},
			wantStdin: true,
			wantPos:   []string{},
		},
		{
			name:      "stdin with banner",
			args:      []string{"--stdin", "shadow"},
			wantStdin: true,
			wantPos:   []string{"shadow"},
		},
		{
			name:        "stdin with other flags",
			args:        []string{"--stdin", "--color=red", "hello", "--output=file.txt"},
			wantStdin:   true,
			wantPos:     []string{"hello"},
			wantOutput:  "file.txt",
			wantAlign:   "",
			wantReverse: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()
			got := cli.ClassifyArgs(tt.args)
			if got.StdinMode != tt.wantStdin {
				st.Errorf("StdinMode = %v, want %v", got.StdinMode, tt.wantStdin)
			}
			if !reflect.DeepEqual(got.Positional, tt.wantPos) {
				st.Errorf("Positional = %v, want %v", got.Positional, tt.wantPos)
			}
			if got.OutputValue != tt.wantOutput {
				st.Errorf("OutputValue = %q, want %q", got.OutputValue, tt.wantOutput)
			}
			if got.AlignValue != tt.wantAlign {
				st.Errorf("AlignValue = %q, want %q", got.AlignValue, tt.wantAlign)
			}
			if got.ReverseValue != tt.wantReverse {
				st.Errorf("ReverseValue = %q, want %q", got.ReverseValue, tt.wantReverse)
			}
		})
	}
}

func BenchmarkClassifyArgs(b *testing.B) {
	cases := []struct {
		name string
		args []string
	}{
		{"minimal", []string{"Hello"}},
		{"with_flags", []string{"--color=red", "Hello", "standard"}},
		{"complex", []string{"--color=red", "sub", "--align=center", "--output=out.txt", "Hello World", "shadow"}},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = cli.ClassifyArgs(tc.args)
			}
		})
	}
}

func BenchmarkBuildColorRules(b *testing.B) {
	old := os.Stderr
	os.Stderr, _ = os.Open(os.DevNull)
	b.Cleanup(func() { os.Stderr = old })

	cases := []struct {
		name string
		raw  []cli.RawColorRule
	}{
		{"named", []cli.RawColorRule{{ColorValue: "red"}}},
		{"hex", []cli.RawColorRule{{ColorValue: "#FF0000", Substring: "Hello"}}},
		{"multi", []cli.RawColorRule{
			{ColorValue: "red"},
			{ColorValue: "#00FF00", Substring: "World"},
			{ColorValue: "hsl(240,100%,50%)"},
		}},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = cli.BuildColorRules(tc.raw)
			}
		})
	}
}
