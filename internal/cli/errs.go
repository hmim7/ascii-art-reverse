package cli

import (
	"fmt"
	"os"
	"strings"

	"ascii-art-reverse/internal/banner"
)

// Usage constants — exact strings required by sources/ audit.
const (
	UsageBasic   = "Usage: go run . [STRING] [BANNER]\n\nEX: go run . something standard"
	UsageColor   = "Usage: go run . [OPTION] [STRING]\n\nEX: go run . --color=<color> <substring to be colored> \"something\""
	UsageOutput  = "Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard"
	UsageAlign   = "Usage: go run . [OPTION] [STRING] [BANNER]\n\nExample: go run . --align=right something standard"
	UsageReverse = "Usage: go run . [OPTION]\n\nEX: go run . --reverse=<fileName>"
)

// Fatal prints the usage string to stderr and exits with code 1.
func Fatal(usage string) {
	fmt.Fprintln(os.Stderr, usage)
	os.Exit(1)
}

// SelectUsage picks the correct UsageXxx constant based on which valid flag
// was present in args. Falls back to UsageBasic when no flag is recognized.
func SelectUsage(args ParsedArgs) string {
	for _, arg := range os.Args {
		if arg == "--" {
			break
		}
		switch {
		case strings.HasPrefix(arg, "--reverse"):
			return UsageReverse
		case strings.HasPrefix(arg, "--output"):
			return UsageOutput
		case strings.HasPrefix(arg, "--align"):
			return UsageAlign
		case strings.HasPrefix(arg, "--color"):
			return UsageColor
		}
	}
	// Unknown flags that look like misspelled known flags (e.g. --colours, --ouput, --revrese, --allign).
	for _, f := range args.UnknownFlags {
		raw := strings.ToLower(f.Raw)
		switch {
		case strings.HasPrefix(raw, "--rev"):
			return UsageReverse
		case strings.HasPrefix(raw, "--out"):
			return UsageOutput
		case strings.HasPrefix(raw, "--al"):
			return UsageAlign
		case strings.HasPrefix(raw, "--col"):
			return UsageColor
		}
	}
	return UsageBasic
}

// EmitWarnings iterates args and fires the appropriate Warn* helper for
// each invalid or duplicate flag. All warnings are non-fatal.
func EmitWarnings(args ParsedArgs) {
	// Collect all validation-failure categories into one line (input order).
	var invalid []FlagError
	for _, f := range args.Malformed {
		switch f.Category {
		case "color", "output", "align", "reverse":
			invalid = append(invalid, f)
		}
	}
	if len(invalid) > 0 {
		WarnInvalidFlags(invalid)
	}

	// Duplicate-flag warnings carry semantic context; keep them separate.
	for _, f := range args.Malformed {
		switch f.Category {
		case "dup-output":
			WarnOutputRedirected(args.OutputValue, f.Raw[9:])
		case "dup-align":
			WarnAlignOverridden(f.Raw, args.AlignValue)
		case "dup-reverse":
			warnf("warning: reverse redirected to %q; previous flag %q ignored", args.ReverseValue, f.Raw)
		}
	}

	// Provide a hint for unknown flags that look like they might be intended as strings.
	var dashFlags []string
	for _, f := range args.UnknownFlags {
		if strings.HasPrefix(f.Raw, "--") && len(f.Raw) >= 2 {
			dashFlags = append(dashFlags, f.Raw)
		}
	}
	if len(dashFlags) > 0 {
		WarnDoubleDashHint(dashFlags...)
	}
}

// warnf centralizes the grey ANSI color wrapping for all non-fatal feedback.
func warnf(format string, a ...interface{}) {
	fmt.Fprint(os.Stderr, "\x1b[30m")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr, "\x1b[0m")
}

func WarnInvalidFlags(flags []FlagError) {
	seenCat := make(map[string]bool)
	var cats, raws []string
	for _, f := range flags {
		cat := strings.TrimPrefix(f.Category, "dup-")
		if !seenCat[cat] {
			seenCat[cat] = true
			cats = append(cats, cat)
		}
		raws = append(raws, fmt.Sprintf("%q", f.Raw))
	}
	noun := "flag"
	if len(raws) > 1 {
		noun = "flags"
	}
	warnf("warning: invalid %s %s %s",
		strings.Join(cats, ", "), noun, strings.Join(raws, ", "))
}

func WarnOutputRedirected(newVal, oldVal string) {
	warnf("warning: output redirected to %q; previous flag %q ignored", newVal, oldVal)
}

func WarnAlignOverridden(oldVal, newVal string) {
	warnf("warning: previous align flag %q overridden by %q", oldVal, newVal)
}

func WarnDoubleDashHint(vals ...string) {
	if len(vals) == 0 {
		return
	}
	quoted := make([]string, len(vals))
	for i, v := range vals {
		quoted[i] = fmt.Sprintf("%q", v)
	}

	warnf("hint: to render %s as ascii-art, use the \"--\" delimiter before [STRING] (e.g., go run . -- %s)",
		strings.Join(quoted, ", "), strings.Join(vals, " "))
}

func WarnBannerNotFound(name string) {
	warnf("warning: banner %q not found", name)
}

func WarnBannerInvalid(name string) {
	warnf("warning: banner %q invalid (expected 855 lines)", name)
}

func WarnReverseMismatch(fileName string) {
	warnf("warning: the provided banner does not match the art in %q; please provide the correct [BANNER] argument", fileName)
}

func WarnFlagsAfterString(flags ...string) {
	if len(flags) == 0 {
		return
	}
	quoted := make([]string, len(flags))
	for i, f := range flags {
		quoted[i] = fmt.Sprintf("%q", f)
	}
	noun := "option"
	verb := "looks like a"
	if len(flags) > 1 {
		noun = "options"
		verb = "look like"
	}
	warnf("hint: %s %s flag %s; when using \"--\", all flag options must be placed before the delimiter (e.g., go run . %s -- [STRING])",
		strings.Join(quoted, ", "), verb, noun, strings.Join(flags, " "))
}

// ValidateOrFatal exits with the appropriate usage message if any malformed
// or unknown flags are present.
func ValidateOrFatal(args ParsedArgs) {
	fatalMalformed := false
	for _, f := range args.Malformed {
		switch f.Category {
		case "color", "output", "align", "reverse":
			fatalMalformed = true
		}
	}
	reverseConflict := args.ReverseValue != "" && (args.OutputValue != "" || args.AlignValue != "" || len(args.ColorRules) > 0 || args.StdinMode)
	reverseOverflow := args.ReverseValue != "" && len(args.Positional) > 1
	noInput := args.ReverseValue == "" && !args.StdinMode && len(args.Positional) == 0

	isMixed := reverseConflict || reverseOverflow || noInput
	if len(args.UnknownFlags) > 0 || fatalMalformed || isMixed {
		Fatal(SelectUsage(args))
	}
}

// CheckPositionalsOrFatal exits with UsageBasic when more than two positional
// arguments are given, emitting a hint for any that look like flag options.
func CheckPositionalsOrFatal(args ParsedArgs) {
	if len(args.Positional) > 2 {
		var misplaced []string
		for _, pos := range args.Positional[1:] {
			if IsKnownFlagOption(pos) {
				misplaced = append(misplaced, pos)
			}
		}
		if len(misplaced) > 0 {
			WarnFlagsAfterString(misplaced...)
		}
		Fatal(UsageBasic)
	}
}

// HandleBannerFallbackOrFatal emits the appropriate warning and exits when
// banner.Load returned a non-nil FallbackInfo.
func HandleBannerFallbackOrFatal(fallback *banner.FallbackInfo, bannerName string) {
	if fallback == nil {
		return
	}
	if fallback.Kind == "not-found" {
		WarnBannerNotFound(fallback.Name)
	} else {
		WarnBannerInvalid(fallback.Name)
	}
	if IsKnownFlagOption(bannerName) {
		WarnFlagsAfterString(bannerName)
	}
	Fatal(UsageBasic)
}
