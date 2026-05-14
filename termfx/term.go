//go:build unix

package termfx

import (
	"syscall"
	"unsafe"
)

// winsize mirrors the kernel's struct winsize.
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// to be responsive
func GetSize() (cols, rows int) {
	ws := winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 || ws.Col == 0 || ws.Row == 0 {
		return 80, 24
	}
	return int(ws.Col), int(ws.Row)
}
