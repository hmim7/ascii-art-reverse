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
	// Return the usage for the first malformed flag (input order) that maps
	// to a known category.
	for _, f := range args.Malformed {
		switch f.Category {
		case "reverse", "dup-reverse":
			return UsageReverse
		case "output", "dup-output":
			return UsageOutput
		case "align", "dup-align":
			return UsageAlign
		case "color":
			return UsageColor
		}
	}

	// No malformed flags — fall back to the highest-priority valid flag.
	if args.ReverseValue != "" {
		return UsageReverse
	}
	if args.OutputValue != "" {
		return UsageOutput
	}
	if args.AlignValue != "" {
		return UsageAlign
	}
	if len(args.ColorRules) > 0 {
		return UsageColor
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
	for _, f := range args.UnknownFlags {
		if strings.HasPrefix(f.Raw, "--") && len(f.Raw) >= 2 {
			WarnDoubleDashHint(f.Raw)
		}
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
		cat := f.Category
		if strings.HasPrefix(cat, "dup-") {
			cat = cat[4:]
		}
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

func WarnDoubleDashHint(val string) {
	warnf("hint: to render %q as ascii-art, use the \"--\" delimiter before [STRING] (e.g., go run . -- %q)", val, val)
}

func WarnBannerNotFound(name string) {
	warnf("warning: banner %q not found", name)
}

func WarnBannerInvalid(name string) {
	warnf("warning: banner %q invalid (expected 855 lines)", name)
}

func WarnFlagsAfterString(flag string) {
	warnf("hint: %q looks like a flag option; when using \"--\", all flag options must be placed before the delimiter (e.g., go run . %s -- [STRING])", flag, flag)
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
	if len(args.UnknownFlags) > 0 || fatalMalformed {
		Fatal(SelectUsage(args))
	}
}

// CheckPositionalsOrFatal exits with UsageBasic when more than two positional
// arguments are given, emitting a hint for any that look like flag options.
func CheckPositionalsOrFatal(args ParsedArgs) {
	if len(args.Positional) > 2 {
		for _, pos := range args.Positional[1:] {
			if IsKnownFlagOption(pos) {
				WarnFlagsAfterString(pos)
			}
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
