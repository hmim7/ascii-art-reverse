//go:build !windows

package render

import (
	"syscall"
	"unsafe"
)

// getWinSize retrieves the terminal column count via ioctl TIOCGWINSZ.
func getWinSize(fd uintptr) (int, error) {
	var ws struct{ Row, Col, Xpixel, Ypixel uint16 }
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		return 0, errno
	}
	return int(ws.Col), nil
}
