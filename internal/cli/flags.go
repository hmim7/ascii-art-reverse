package cli

import (
	"io"
	"os"
	"strings"

	"ascii-art-reverse/internal/render"
)

// ValidFlag holds a successfully parsed flag.
type ValidFlag struct {
	Name  string // "output" | "align" | "color" | "reverse" | "stdin"
	Value string // raw value after =; empty for --stdin
}

// FlagError holds a flag token that failed validation.
type FlagError struct {
	Raw      string // original token as typed
	Category string // "color"|"align"|"output"|"reverse"|"unknown"|"malformed"|"dup-output"|"dup-align"
}

// RawColorRule holds one --color=<val> [substring] pair before ANSI resolution.
type RawColorRule struct {
	ColorValue string
	Substring  string // next positional token; empty = whole string colored
}

// ParsedArgs is the result of a single-pass classification of os.Args[1:].
type ParsedArgs struct {
	OutputValue  string
	AlignValue   string // "" means default (left); set when --align= is provided
	ReverseValue string
	ColorRules   []RawColorRule
	StdinMode    bool
	Positional   []string // [text] or [text, banner]
	UnknownFlags []FlagError
	Malformed    []FlagError
}

// classifier holds mutable state for a single ClassifyArgs pass.
type classifier struct {
	args          []string
	result        ParsedArgs
	alignExplicit bool
}

func newClassifier(args []string) *classifier {
	return &classifier{
		args: args,
		result: ParsedArgs{
			Positional:   []string{},
			ColorRules:   []RawColorRule{},
			UnknownFlags: []FlagError{},
			Malformed:    []FlagError{},
		},
	}
}

// countNonFlags returns the number of non-flag tokens in args[from:].
func (c *classifier) countNonFlags(from int) int {
	n, stopFlags := 0, false
	for j := from; j < len(c.args); j++ {
		if stopFlags {
			n++
			continue
		}
		if c.args[j] == "--" && j < len(c.args)-1 && strings.HasPrefix(c.args[j+1], "--") {
			stopFlags = true
			continue
		}
		if !strings.HasPrefix(c.args[j], "--") {
			n++
		}
	}
	return n
}

func (c *classifier) parseColor(tok string, i int) int {
	rule := RawColorRule{ColorValue: tok[8:]}
	if i+1 < len(c.args) && !strings.HasPrefix(c.args[i+1], "--") && c.countNonFlags(i+2) > 0 {
		rule.Substring = c.args[i+1]
		i++
	}
	c.result.ColorRules = append(c.result.ColorRules, rule)
	return i
}

func (c *classifier) parseOutput(tok string) {
	val := tok[9:]
	if !strings.HasSuffix(strings.ToLower(val), ".txt") || val == "" {
		c.result.Malformed = append(c.result.Malformed, FlagError{Raw: tok, Category: "output"})
		return
	}
	if c.result.OutputValue != "" {
		c.result.Malformed = append(c.result.Malformed, FlagError{
			Raw:      "--output=" + c.result.OutputValue,
			Category: "dup-output",
		})
	}
	c.result.OutputValue = val
}

func (c *classifier) parseAlign(tok string) {
	val := tok[8:]
	switch val {
	case "left", "right", "center", "justify":
		if c.alignExplicit {
			c.result.Malformed = append(c.result.Malformed, FlagError{
				Raw:      c.result.AlignValue,
				Category: "dup-align",
			})
		}
		c.result.AlignValue = val
		c.alignExplicit = true
	default:
		c.result.Malformed = append(c.result.Malformed, FlagError{Raw: tok, Category: "align"})
	}
}

func (c *classifier) parseReverse(tok string) {
	val := tok[10:]
	if c.result.ReverseValue != "" {
		c.result.Malformed = append(c.result.Malformed, FlagError{
			Raw:      "--reverse=" + c.result.ReverseValue,
			Category: "dup-reverse",
		})
	}
	c.result.ReverseValue = val
}

// parseDelimiter handles "--". Returns the updated i (before the outer i++).
// Setting i = len(args)-1 causes the outer i++ to land on len(args), ending the loop.
func (c *classifier) parseDelimiter(i int) int {
	if i+1 < len(c.args) && strings.HasPrefix(c.args[i+1], "--") {
		c.result.Positional = append(c.result.Positional, c.args[i+1:]...)
		return len(c.args) - 1
	}
	c.result.UnknownFlags = append(c.result.UnknownFlags, FlagError{Raw: "--", Category: "unknown"})
	return i
}

// classifyBareFlag handles tokens starting with "--" that lack a valid "=value".
// Checks known prefixes first; falls back to unknown.
func (c *classifier) classifyBareFlag(tok string) {
	known := []struct{ prefix, cat string }{
		{"--color", "color"}, {"--output", "output"},
		{"--align", "align"}, {"--reverse", "reverse"},
	}
	for _, k := range known {
		if strings.HasPrefix(tok, k.prefix) {
			c.result.Malformed = append(c.result.Malformed, FlagError{Raw: tok, Category: k.cat})
			return
		}
	}
	c.result.UnknownFlags = append(c.result.UnknownFlags, FlagError{Raw: tok, Category: "unknown"})
}

func (c *classifier) classify(tok string, i int) int {
	switch {
	case tok == "--stdin":
		c.result.StdinMode = true
	case strings.HasPrefix(tok, "--color="):
		i = c.parseColor(tok, i)
	case strings.HasPrefix(tok, "--output="):
		c.parseOutput(tok)
	case strings.HasPrefix(tok, "--align="):
		c.parseAlign(tok)
	case strings.HasPrefix(tok, "--reverse="):
		c.parseReverse(tok)
	case tok == "--":
		i = c.parseDelimiter(i)
	case strings.HasPrefix(tok, "--"):
		c.classifyBareFlag(tok)
	default:
		c.result.Positional = append(c.result.Positional, tok)
	}
	return i
}

// ClassifyArgs performs a single pass over args, routing every token into
// exactly one bucket of ParsedArgs. No token is processed twice.
func ClassifyArgs(args []string) ParsedArgs {
	c := newClassifier(args)
	for i := 0; i < len(args); i++ {
		i = c.classify(args[i], i)
	}
	return c.result
}

// BuildColorRules converts []RawColorRule into []render.ColorRule, resolving
// ANSI codes and warning on invalid color values.
func BuildColorRules(raw []RawColorRule) []render.ColorRule {
	result := make([]render.ColorRule, 0, len(raw))
	for _, r := range raw {
		ansi, ok := render.ColorToANSI(r.ColorValue)
		if !ok {
			WarnInvalidFlags([]FlagError{{Category: "color", Raw: "--color=" + r.ColorValue}})
			continue
		}
		result = append(result, render.ColorRule{
			ANSIStart: ansi,
			Substring: r.Substring,
		})
	}
	return result
}

// ReadStdin reads all of os.Stdin and returns it as a string.
func ReadStdin() string {
	data, _ := io.ReadAll(os.Stdin)
	return strings.TrimRight(string(data), "\n")
}

// IsKnownFlagOption reports whether s starts with a recognized flag prefix.
func IsKnownFlagOption(s string) bool {
	return strings.HasPrefix(s, "--color") ||
		strings.HasPrefix(s, "--output") ||
		strings.HasPrefix(s, "--align") ||
		strings.HasPrefix(s, "--reverse")
}

// ResolveBannerName returns the banner name from positional args, defaulting to "standard".
func ResolveBannerName(args ParsedArgs) string {
	if len(args.Positional) == 2 {
		return args.Positional[1]
	}
	return "standard"
}

// ResolveInput returns the text to render: stdin content when StdinMode is set,
// the first positional argument otherwise, or "" when neither is present.
func ResolveInput(args ParsedArgs) string {
	if args.StdinMode {
		return ReadStdin()
	}
	if len(args.Positional) > 0 {
		return args.Positional[0]
	}
	return ""
}
