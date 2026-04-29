package render_test

import (
	"testing"

	"ascii-art-reverse/internal/render"
)

// TODO(task05): expand ParseInput and Mapper tests.
// TODO(task06): add ColorToANSI and WrapWithColor tests.
// TODO(task08): add alignment and justify tests.

func TestParseInput_NewlineSplit(t *testing.T) {
	got := render.ParseInput("Hello\nWorld")
	if len(got) != 2 {
		t.Errorf("expected 2 segments, got %d", len(got))
	}
}

func TestParseInput_Empty(t *testing.T) {
	got := render.ParseInput("")
	if len(got) != 0 {
		t.Errorf("expected 0 segments, got %d", len(got))
	}
}

func TestParseInput_OnlyNewline(t *testing.T) {
	got := render.ParseInput("\\n")
	if len(got) != 1 {
		t.Errorf("expected 1 segment for single newline, got %d", len(got))
	}
}

func TestWrapWithColor_ShortCircuit(t *testing.T) {
	if got := render.WrapWithColor("hi", ""); got != "hi" {
		t.Errorf("expected short-circuit on empty ansi, got %q", got)
	}
	if got := render.WrapWithColor("", "\x1b[31m"); got != "" {
		t.Errorf("expected short-circuit on empty s, got %q", got)
	}
}

func TestShouldRenderGopher_EmptyText(t *testing.T) {
	if render.ShouldRenderGopher("", []string{""}) {
		t.Error("empty text should not trigger gopher")
	}
}
