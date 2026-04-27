//go:build windows

package render

import "errors"

// getWinSize is a stub on Windows; DetectTerminalWidth falls back to 80.
func getWinSize(fd uintptr) (int, error) {
	return 0, errors.New("terminal size detection not supported on Windows")
}
