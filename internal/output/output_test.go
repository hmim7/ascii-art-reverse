package output_test

import (
	"os"
	"testing"

	"ascii-art-reverse/internal/output"
)

func TestGetWriter_Stdout(t *testing.T) {
	w, f, err := output.GetWriter("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Error("expected nil file for stdout path")
	}
	if w != os.Stdout {
		t.Error("expected os.Stdout writer")
	}
}

func TestGetWriter_File(t *testing.T) {
	path := t.TempDir() + "/test.txt"
	w, f, err := output.GetWriter(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil file")
	}
	defer f.Close()
	if w == nil {
		t.Error("expected non-nil writer")
	}
}
