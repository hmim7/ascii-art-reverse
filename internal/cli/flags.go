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

// ClassifyArgs performs a single pass over args, routing every token into
// exactly one bucket of ParsedArgs. No token is processed twice.
func ClassifyArgs(args []string) ParsedArgs {
	result := ParsedArgs{
		Positional:   []string{},
		ColorRules:   []RawColorRule{},
		UnknownFlags: []FlagError{},
		Malformed:    []FlagError{},
	}

	// countNonFlags returns the number of non-flag tokens in args[from:].
	countNonFlags := func(from int) int {
		n := 0
		stopFlags := false
		for j := from; j < len(args); j++ {
			if stopFlags {
				n++
				continue
			}
			if args[j] == "--" && j < len(args)-1 && strings.HasPrefix(args[j+1], "--") {
				stopFlags = true
				continue
			}
			if !strings.HasPrefix(args[j], "--") {
				n++
			}
		}
		return n
	}

	alignExplicit := false

	i := 0
	for i < len(args) {
		tok := args[i]
		switch {
		case tok == "--stdin":
			result.StdinMode = true

		case strings.HasPrefix(tok, "--color="):
			val := tok[8:]
			rule := RawColorRule{ColorValue: val}
			// Peek: consume the next token as a substring only when at least
			// one more non-flag token follows it (so it's not the sole string).
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") && countNonFlags(i+2) > 0 {
				rule.Substring = args[i+1]
				i++
			}
			result.ColorRules = append(result.ColorRules, rule)

		case strings.HasPrefix(tok, "--output="):
			val := tok[9:]
			if !strings.HasSuffix(strings.ToLower(val), ".txt") || val == "" {
				// Invalid extension or empty value — treat as malformed.
				result.Malformed = append(result.Malformed, FlagError{Raw: tok, Category: "output"})
			} else if result.OutputValue != "" {
				// Duplicate valid output: last wins, record old value for warning.
				result.Malformed = append(result.Malformed, FlagError{
					Raw:      "--output=" + result.OutputValue,
					Category: "dup-output",
				})
				result.OutputValue = val
			} else {
				result.OutputValue = val
			}

		case strings.HasPrefix(tok, "--align="):
			val := tok[8:]
			switch val {
			case "left", "right", "center", "justify":
				if alignExplicit {
					// Duplicate valid align: last wins, record old value for warning.
					result.Malformed = append(result.Malformed, FlagError{
						Raw:      result.AlignValue,
						Category: "dup-align",
					})
				}
				result.AlignValue = val
				alignExplicit = true
			default:
				// Unknown align type.
				result.Malformed = append(result.Malformed, FlagError{Raw: tok, Category: "align"})
			}

		case strings.HasPrefix(tok, "--reverse="):
			val := tok[10:]
			if result.ReverseValue != "" {
				result.Malformed = append(result.Malformed, FlagError{
					Raw:      "--reverse=" + result.ReverseValue,
					Category: "dup-reverse",
				})
			}
			result.ReverseValue = val

		case tok == "--":
			// Delimiter only if it precedes something that looks like a flag.
			if i+1 < len(args) && strings.HasPrefix(args[i+1], "--") {
				result.Positional = append(result.Positional, args[i+1:]...)
				i = len(args) // Skip remaining arguments
				continue
			}
			result.UnknownFlags = append(result.UnknownFlags, FlagError{Raw: tok, Category: "unknown"})

		// Known flag names without proper =value syntax.
		case strings.HasPrefix(tok, "--color"):
			result.Malformed = append(result.Malformed, FlagError{Raw: tok, Category: "color"})
		case strings.HasPrefix(tok, "--output"):
			result.Malformed = append(result.Malformed, FlagError{Raw: tok, Category: "output"})
		case strings.HasPrefix(tok, "--align"):
			result.Malformed = append(result.Malformed, FlagError{Raw: tok, Category: "align"})
		case strings.HasPrefix(tok, "--reverse"):
			result.Malformed = append(result.Malformed, FlagError{Raw: tok, Category: "reverse"})

		case strings.HasPrefix(tok, "--"):
			result.UnknownFlags = append(result.UnknownFlags, FlagError{Raw: tok, Category: "unknown"})

		default:
			result.Positional = append(result.Positional, tok)
		}
		i++
	}

	return result
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
