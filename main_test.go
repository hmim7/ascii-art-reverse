package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"ascii-art-reverse/internal/banner"
	"ascii-art-reverse/internal/cli"
	"ascii-art-reverse/internal/render"
)

// runRender calls the render pipeline directly, bypassing main() and its os.Exit calls.
func runRender(t *testing.T, args []string) string {
	t.Helper()
	parsed := cli.ClassifyArgs(args)
	bannerName := cli.ResolveBannerName(parsed)
	bannerMap, fallback, err := banner.Load(bannerName)
	if err != nil {
		t.Fatalf("banner.Load(%q): %v", bannerName, err)
	}
	if fallback != nil {
		t.Fatalf("unexpected banner fallback: %+v", fallback)
	}
	segments := render.ParseInput(cli.ResolveInput(parsed))
	rules := cli.BuildColorRules(parsed.ColorRules)
	var buf bytes.Buffer
	render.RenderAlignedWithColorRules(&buf, segments, bannerMap, rules, parsed.AlignValue, 80)
	return buf.String()
}

// countLines returns the number of output lines, treating empty string as 0.
func countLines(s string) int {
	if s == "" {
		return 0
	}
	return len(strings.Split(strings.TrimRight(s, "\n"), "\n"))
}

func TestIntegration_RenderPipeline(t *testing.T) {
	tests := []struct {
		args             []string
		name             string
		wantContains     string // substring that must appear in output
		wantLines        int    // line count when > 0
		wantEmpty        bool   // output must be exactly ""
		wantLeadingSpace bool   // every line must start with a space
	}{
		{
			name:      "empty input produces no output",
			args:      []string{""},
			wantEmpty: true,
		},
		{
			name:      `\n escape produces one blank line`,
			args:      []string{`\n`},
			wantLines: 1,
		},
		{
			name:      "single word renders 8 art lines",
			args:      []string{"Hi"},
			wantLines: 8,
		},
		{
			name:      "two words on one line render 8 art lines",
			args:      []string{"Hi World"},
			wantLines: 8,
		},
		{
			name:      `single \n between words: no blank separator`,
			args:      []string{`Hello\nWorld`},
			wantLines: 16,
		},
		{
			name:      `double \n between words: blank separator`,
			args:      []string{`Hello\n\nWorld`},
			wantLines: 17,
		},
		{
			name:         "--color=red injects ANSI red escape",
			args:         []string{"--color=red", "Hi"},
			wantContains: "\x1b[31m",
		},
		{
			name:             "--align=center adds leading spaces",
			args:             []string{"--align=center", "A"},
			wantLines:        8,
			wantLeadingSpace: true,
		},
		{
			name:             "--align=right adds leading spaces",
			args:             []string{"--align=right", "A"},
			wantLines:        8,
			wantLeadingSpace: true,
		},
		{
			name:      "shadow banner renders non-empty output",
			args:      []string{"Hi", "shadow"},
			wantLines: 8,
		},
		{
			name:      "doom banner renders non-empty output",
			args:      []string{"Hi", "doom"},
			wantLines: 8,
		},
		{
			name:      "thinkertoy banner renders non-empty output",
			args:      []string{"Hi", "thinkertoy"},
			wantLines: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			out := runRender(st, tt.args)

			if tt.wantEmpty && out != "" {
				st.Errorf("want empty output, got %q", out)
			}
			if tt.wantLines > 0 && countLines(out) != tt.wantLines {
				st.Errorf("want %d lines, got %d", tt.wantLines, countLines(out))
			}
			if tt.wantContains != "" && !strings.Contains(out, tt.wantContains) {
				st.Errorf("want output to contain %q\nfull output: %q", tt.wantContains, out)
			}
			if tt.wantLeadingSpace {
				for _, l := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
					if !strings.HasPrefix(l, " ") {
						st.Errorf("expected leading space, got line %q", l)
					}
				}
			}
		})
	}
}

// TestMain_Happy exercises the main() orchestration path for inputs that never
// call os.Exit, covering all non-error branches of the top-level dispatcher.
func TestMain_ExecutionBranches(t *testing.T) {
	outPath := t.TempDir() + "/art.txt"
	tests := []struct {
		args  []string
		check func(st *testing.T, stdout string)
		name  string
	}{
		{name: "no args", args: []string{"ascii-art-reverse", " "}},
		{name: "empty string", args: []string{"ascii-art-reverse", ""}},
		{name: "single word", args: []string{"ascii-art-reverse", "Hi"}},
		{name: "word and banner", args: []string{"ascii-art-reverse", "Hi", "shadow"}},
		{
			name: "reverse mode",
			args: []string{"ascii-art-reverse", "--reverse=internal/reverse/testdata/example00.txt"},
			check: func(st *testing.T, stdout string) {
				if !strings.Contains(stdout, "Hello World") {
					st.Errorf("output missing 'Hello World', got: %q", stdout)
				}
			},
		},
		{
			name: "output file mode",
			args: []string{"ascii-art-reverse", "--output=" + outPath, "Hi"},
			check: func(st *testing.T, stdout string) {
				data, err := os.ReadFile(outPath)
				if err != nil {
					st.Fatalf("output file not created: %v", err)
				}
				if len(data) == 0 {
					st.Error("output file is empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			oldArgs := os.Args
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			st.Cleanup(func() {
				os.Args = oldArgs
				os.Stdout = oldStdout
			})

			os.Args = tt.args
			os.Stdout = w

			main() // main() must not panic or call os.Exit
			if err := w.Close(); err != nil {
				st.Fatalf("Failed to close pipe writer: %v", err)
			}

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)

			if tt.check != nil {
				tt.check(st, buf.String())
			}
		})
	}
}

// FuzzRenderPipeline checks the render pipeline never panics on arbitrary input.
func FuzzRenderPipeline(f *testing.F) {
	bannerMap, fallback, err := banner.Load("standard")
	if err != nil || fallback != nil {
		f.Fatal("failed to load standard banner for fuzz corpus")
	}

	f.Add("Hello World")
	f.Add("")
	f.Add("\n")
	f.Add("Hi\nThere")
	f.Add("!@#$%^&*()")
	f.Add("   spaces   ")
	f.Add("A B C D E F")
	f.Add("Hello\nWorld\n!")

	f.Fuzz(func(t *testing.T, input string) {
		segments := render.ParseInput(input)
		var buf bytes.Buffer
		render.RenderAlignedWithColorRules(&buf, segments, bannerMap, nil, "", 80)
	})
}

// FuzzClassifyArgs checks the classifier never panics and always returns
// initialized slices on arbitrary token streams.
func FuzzClassifyArgs(f *testing.F) {
	f.Add("hello standard")
	f.Add("--color=red hello")
	f.Add("--align=center hello standard")
	f.Add("--output=out.txt hello")
	f.Add("--reverse=file.txt")
	f.Add("")
	f.Add("-- --color=red")
	f.Add("--color=red --align=right hello shadow")

	f.Fuzz(func(t *testing.T, input string) {
		result := cli.ClassifyArgs(strings.Fields(input))
		if result.Positional == nil {
			t.Error("Positional must never be nil")
		}
		if result.ColorRules == nil {
			t.Error("ColorRules must never be nil")
		}
		if result.UnknownFlags == nil {
			t.Error("UnknownFlags must never be nil")
		}
		if result.Malformed == nil {
			t.Error("Malformed must never be nil")
		}
	})
}
