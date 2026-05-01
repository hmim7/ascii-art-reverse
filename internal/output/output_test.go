package output_test

import (
	"os"
	"testing"

	"ascii-art-reverse/internal/output"
)

func TestGetWriter_Table(t *testing.T) {
	tests := []struct {
		pathFn     func() string // lazy so TempDir is per-test
		name       string
		wantErr    bool
		wantFile   bool
		wantStdout bool
	}{
		{
			name:       "empty path returns stdout",
			pathFn:     func() string { return "" },
			wantStdout: true,
		},
		{
			name:     "valid txt path creates file",
			pathFn:   func() string { return t.TempDir() + "/out.txt" },
			wantFile: true,
		},
		{
			name:    "nonexistent directory returns error",
			pathFn:  func() string { return t.TempDir() + "/nodir/out.txt" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			w, f, err := output.GetWriter(tt.pathFn())
			if tt.wantErr {
				if err == nil {
					st.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				st.Fatalf("unexpected error: %v", err)
			}
			if f != nil {
				defer func() { _ = f.Close() }()
			}
			if tt.wantStdout && w != os.Stdout {
				st.Error("expected os.Stdout")
			}
			if tt.wantFile && f == nil {
				st.Error("expected non-nil file")
			}
		})
	}
}
