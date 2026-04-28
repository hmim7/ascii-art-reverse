package cli

import (
	"fmt"
	"os"
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
// Priority order: Reverse > Output > Align > Color > Basic.
func SelectUsage(args ParsedArgs) string {
	hasMal := func(cat string) bool {
		for _, f := range args.Malformed {
			if f.Category == cat {
				return true
			}
		}
		return false
	}
	if args.ReverseValue != "" || hasMal("reverse") {
		return UsageReverse
	}
	if args.OutputValue != "" || hasMal("output") || hasMal("dup-output") {
		return UsageOutput
	}
	if args.AlignValue != "" || hasMal("align") || hasMal("dup-align") {
		return UsageAlign
	}
	if len(args.ColorRules) > 0 || hasMal("color") {
		return UsageColor
	}
	return UsageBasic
}

// EmitWarnings iterates args and fires the appropriate Warn* helper for
// each invalid or duplicate flag. All warnings are non-fatal.
func EmitWarnings(args ParsedArgs) {
	for _, f := range args.Malformed {
		switch f.Category {
		case "dup-output":
			// Raw = "--output=<old_value>"; strip prefix to get old value.
			WarnOutputRedirected(args.OutputValue, f.Raw[9:])
		case "dup-align":
			// Raw = old align value string.
			WarnAlignOverridden(f.Raw, args.AlignValue)
		case "output":
			WarnInvalidOutputFlag(f.Raw)
		case "align":
			WarnInvalidAlignFlag(f.Raw)
		case "color":
			WarnInvalidColor(f.Raw)
		}
	}
}

func WarnInvalidColor(val string) {
	fmt.Fprintf(os.Stderr, "warning: invalid color %q, rendering without color\n", val)
}

func WarnInvalidOutputFlag(flag string) {
	fmt.Fprintf(os.Stderr, "warning: invalid output flag %q ignored\n", flag)
}

func WarnInvalidAlignFlag(flag string) {
	fmt.Fprintf(os.Stderr, "warning: invalid align flag %q ignored\n", flag)
}

func WarnOutputRedirected(newVal, oldVal string) {
	fmt.Fprintf(os.Stderr, "warning: output redirected to %q; previous flag %q ignored\n", newVal, oldVal)
}

func WarnAlignOverridden(oldVal, newVal string) {
	fmt.Fprintf(os.Stderr, "warning: previous align flag %q overridden by %q\n", oldVal, newVal)
}

func WarnBannerNotFound(name string) {
	fmt.Fprintf(os.Stderr, "warning: banner %q not found, default banner \"standard\" applied\n", name)
}

func WarnBannerInvalid(name string) {
	fmt.Fprintf(os.Stderr, "warning: banner %q invalid (expected 855 lines), default banner \"standard\" applied\n", name)
}
