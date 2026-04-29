package cli_test

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"ascii-art-reverse/internal/cli"
)

func TestClassifyArgs_BasicPositional(t *testing.T) {
	got := cli.ClassifyArgs([]string{"Hello", "standard"})
	if len(got.Positional) != 2 {
		t.Errorf("expected 2 positional args, got %d", len(got.Positional))
	}
}

func TestClassifyArgs_UnknownFlag(t *testing.T) {
	got := cli.ClassifyArgs([]string{"--unknown", "Hello"})
	if len(got.UnknownFlags) == 0 {
		t.Error("expected unknown flag to be captured")
	}
}

func TestClassifyArgs_DoubleDashDelimiter(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantPos []string
		wantOut string
		wantCol int // count of color rules
		wantUnk int // count of unknown flags
	}{
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
		t.Run(tt.name, func(t *testing.T) {
			got := cli.ClassifyArgs(tt.args)

			if !reflect.DeepEqual(got.Positional, tt.wantPos) {
				t.Errorf("Positional = %v, want %v", got.Positional, tt.wantPos)
			}
			if got.OutputValue != tt.wantOut {
				t.Errorf("OutputValue = %q, want %q", got.OutputValue, tt.wantOut)
			}
			if len(got.ColorRules) != tt.wantCol {
				t.Errorf("len(ColorRules) = %d, want %d", len(got.ColorRules), tt.wantCol)
			}
			if len(got.UnknownFlags) != tt.wantUnk {
				t.Errorf("len(UnknownFlags) = %d, want %d", len(got.UnknownFlags), tt.wantUnk)
			}
		})
	}
}

func TestSelectUsage(t *testing.T) {
	tests := []struct {
		name string
		args cli.ParsedArgs
		want string
	}{
		{
			name: "Basic Usage",
			args: cli.ParsedArgs{},
			want: cli.UsageBasic,
		},
		{
			name: "Reverse Priority over All",
			args: cli.ParsedArgs{
				ReverseValue: "file.txt",
				OutputValue:  "out.txt",
				AlignValue:   "right",
				ColorRules:   []cli.RawColorRule{{ColorValue: "red"}},
			},
			want: cli.UsageReverse,
		},
		{
			name: "Output Priority over Align",
			args: cli.ParsedArgs{
				OutputValue: "out.txt",
				AlignValue:  "center",
			},
			want: cli.UsageOutput,
		},
		{
			name: "Align Priority over Color",
			args: cli.ParsedArgs{
				AlignValue: "justify",
				ColorRules: []cli.RawColorRule{{ColorValue: "blue"}},
			},
			want: cli.UsageAlign,
		},
		{
			name: "Color Usage",
			args: cli.ParsedArgs{
				ColorRules: []cli.RawColorRule{{ColorValue: "green"}},
			},
			want: cli.UsageColor,
		},
		{
			name: "Malformed Output triggers Output Usage",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{{Category: "output", Raw: "--output="}},
			},
			want: cli.UsageOutput,
		},
		{
			name: "Unknown Flag defaults to Basic Usage",
			args: cli.ParsedArgs{
				UnknownFlags: []cli.FlagError{{Category: "unknown", Raw: "--unknown"}},
			},
			want: cli.UsageBasic,
		},
		{
			name: "Unknown Flag alongside malformed Align shows Align Usage",
			args: cli.ParsedArgs{
				UnknownFlags: []cli.FlagError{{Category: "unknown", Raw: "--"}},
				Malformed:    []cli.FlagError{{Category: "align", Raw: "--align"}},
			},
			want: cli.UsageAlign,
		},
		{
			name: "Output first then Align returns Output Usage",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "output", Raw: "--output=bad.pdf"},
					{Category: "align", Raw: "--align=middle"},
				},
			},
			want: cli.UsageOutput,
		},
		{
			name: "Align first then Color returns Align Usage",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "align", Raw: "--align=middle"},
					{Category: "color", Raw: "--color red"},
				},
			},
			want: cli.UsageAlign,
		},
		{
			name: "Color first then Output returns Color Usage",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "color", Raw: "--color red"},
					{Category: "output", Raw: "--output=bad.pdf"},
				},
			},
			want: cli.UsageColor,
		},
		{
			name: "Reverse first then Output and Align returns Reverse Usage",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "reverse", Raw: "--reverse"},
					{Category: "output", Raw: "--output=bad.pdf"},
					{Category: "align", Raw: "--align=middle"},
				},
			},
			want: cli.UsageReverse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cli.SelectUsage(tt.args); got != tt.want {
				t.Errorf("SelectUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmitWarnings(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	args := cli.ParsedArgs{
		OutputValue: "new.txt",
		Malformed: []cli.FlagError{
			{Category: "dup-output", Raw: "--output=old.txt"},
			{Category: "color", Raw: "invalid"},
		},
		UnknownFlags: []cli.FlagError{
			{Category: "unknown", Raw: "---"},
		},
	}

	cli.EmitWarnings(args)

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedSubstrings := []string{
		`warning: output redirected to "new.txt"; previous flag "old.txt" ignored`,
		`warning: invalid color flag "invalid"`,
		`hint: to render "---" as ascii-art, use the "--" delimiter before [STRING] (e.g., go run . -- "---")`,
	}

	for _, exp := range expectedSubstrings {
		if !strings.Contains(output, exp) {
			t.Errorf("EmitWarnings() output missing expected substring: %q\nFull output: %q", exp, output)
		}
	}
}

func TestEmitWarnings_MultipleInvalid(t *testing.T) {
	tests := []struct {
		name          string
		args          cli.ParsedArgs
		wantSubstring string
	}{
		{
			name: "output then align produces one line in input order",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "output", Raw: "--output=bad.pdf"},
					{Category: "align", Raw: "--align=middle"},
				},
			},
			wantSubstring: `warning: invalid output, align flags "--output=bad.pdf", "--align=middle"`,
		},
		{
			name: "align then color preserves input order",
			args: cli.ParsedArgs{
				Malformed: []cli.FlagError{
					{Category: "align", Raw: "--align=middle"},
					{Category: "color", Raw: "--color"},
				},
			},
			wantSubstring: `warning: invalid align, color flags "--align=middle", "--color"`,
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
			wantSubstring: `warning: invalid reverse, output, align flags "--reverse", "--output=bad.pdf", "--align=x"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			cli.EmitWarnings(tt.args)

			w.Close()
			os.Stderr = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if !strings.Contains(output, tt.wantSubstring) {
				t.Errorf("EmitWarnings() missing %q\nGot: %q", tt.wantSubstring, output)
			}
		})
	}
}

func TestClassifyArgs_MixedFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantOutput     string
		wantReverse    string
		wantAlign      string
		wantPos        []string
		wantColorCount int
		wantMalformed  int
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cli.ClassifyArgs(tt.args)

			if got.OutputValue != tt.wantOutput {
				t.Errorf("%s: OutputValue = %q, want %q", tt.name, got.OutputValue, tt.wantOutput)
			}
			if got.AlignValue != tt.wantAlign {
				t.Errorf("%s: AlignValue = %q, want %q", tt.name, got.AlignValue, tt.wantAlign)
			}
			if got.ReverseValue != tt.wantReverse {
				t.Errorf("%s: ReverseValue = %q, want %q", tt.name, got.ReverseValue, tt.wantReverse)
			}
			if !reflect.DeepEqual(got.Positional, tt.wantPos) {
				t.Errorf("%s: Positional = %v, want %v", tt.name, got.Positional, tt.wantPos)
			}
			if len(got.ColorRules) != tt.wantColorCount {
				t.Errorf("%s: ColorRules count = %d, want %d", tt.name, len(got.ColorRules), tt.wantColorCount)
			}
			if len(got.Malformed) != tt.wantMalformed {
				t.Errorf("%s: Malformed count = %d, want %d", tt.name, len(got.Malformed), tt.wantMalformed)
			}
		})
	}
}

func TestClassifyArgs_StdinMode(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantStdin   bool
		wantPos     []string
		wantOutput  string
		wantAlign   string
		wantReverse string
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
		t.Run(tt.name, func(t *testing.T) {
			got := cli.ClassifyArgs(tt.args)
			if got.StdinMode != tt.wantStdin {
				t.Errorf("StdinMode = %v, want %v", got.StdinMode, tt.wantStdin)
			}
			if !reflect.DeepEqual(got.Positional, tt.wantPos) {
				t.Errorf("Positional = %v, want %v", got.Positional, tt.wantPos)
			}
			if got.OutputValue != tt.wantOutput {
				t.Errorf("OutputValue = %q, want %q", got.OutputValue, tt.wantOutput)
			}
			if got.AlignValue != tt.wantAlign {
				t.Errorf("AlignValue = %q, want %q", got.AlignValue, tt.wantAlign)
			}
			if got.ReverseValue != tt.wantReverse {
				t.Errorf("ReverseValue = %q, want %q", got.ReverseValue, tt.wantReverse)
			}
		})
	}
}
