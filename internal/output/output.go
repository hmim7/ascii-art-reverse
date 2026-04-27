package output

import (
	"io"
	"os"
	"path/filepath"
)

// GetWriter resolves the io.Writer for the render pipeline.
//
//   - outputPath == "" → returns (os.Stdout, nil, nil)
//   - outputPath != "" → opens the file with O_WRONLY|O_CREATE|O_TRUNC, perm 0644
//
// The *os.File is returned separately so the caller can defer file.Close().
// Always overwrites (O_TRUNC); appending is not supported.
// ANSI escape codes are written unchanged into the file.
func GetWriter(outputPath string) (io.Writer, *os.File, error) {
	if outputPath == "" {
		return os.Stdout, nil, nil
	}
	f, err := os.OpenFile(filepath.Clean(outputPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}
