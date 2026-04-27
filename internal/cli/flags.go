package cli

// ValidFlag holds a successfully parsed flag.
type ValidFlag struct {
	Name  string // "output" | "align" | "color" | "reverse" | "stdin"
	Value string // raw value after =; empty for --stdin
}

// FlagError holds a flag token that failed validation.
type FlagError struct {
	Raw      string // original token as typed
	Category string // "color"|"align"|"output"|"reverse"|"unknown"|"malformed"
}

// RawColorRule holds one --color=<val> [substring] pair before ANSI resolution.
type RawColorRule struct {
	ColorValue string
	Substring  string // next positional token; empty = whole string colored
}

// ParsedArgs is the result of a single-pass classification of os.Args[1:].
type ParsedArgs struct {
	OutputValue  string
	AlignValue   string // default "left"
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
	// TODO(task02): implement single-pass classifier
	return ParsedArgs{AlignValue: "left"}
}

// BuildColorRules converts []RawColorRule into []ColorRule for the render
// package, resolving ANSI codes and warning on invalid color values.
func BuildColorRules(raw []RawColorRule) []interface{} {
	// TODO(task02): resolve ANSI via colorizer; return []render.ColorRule
	// Placeholder — replace interface{} with render.ColorRule once render/ exists.
	return nil
}

// ReadStdin reads all of os.Stdin and returns it as a string.
func ReadStdin() string {
	// TODO(task02): implement stdin reader
	return ""
}
