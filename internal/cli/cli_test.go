package cli_test

import (
	"testing"

	"ascii-art-reverse/internal/cli"
)

// TODO(task02): expand with full ClassifyArgs test table.
// TODO(task03): add SelectUsage and EmitWarnings tests.

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
