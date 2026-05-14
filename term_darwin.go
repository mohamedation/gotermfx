//go:build darwin

package main

import (
	"syscall"
	"unsafe"
)

func setRawMode() (*syscall.Termios, error) {
	old := &syscall.Termios{}
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		syscall.TIOCGETA,
		uintptr(unsafe.Pointer(old)),
	); errno != 0 {
		return nil, errno
	}
	raw := *old
	// Disable echo, line buffering, signals, and extended processing.
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		syscall.TIOCSETA,
		uintptr(unsafe.Pointer(&raw)),
	); errno != 0 {
		return nil, errno
	}
	return old, nil
}

func restoreMode(old *syscall.Termios) {
	syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		syscall.TIOCSETA,
		uintptr(unsafe.Pointer(old)),
	)
}
