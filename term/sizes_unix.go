// Copyright 2013 Apcera Inc. All rights reserved.

// +build cgo,!windows

package term

/*
#include <termios.h>

// provides struct winsize *, with:
//    ws_row, ws_col, ws_xpixel, ws_ypixel
// all short
*/
import "C"

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

var ErrGetWinsizeFailed = errors.New("term: syscall.TIOCGWINSZ failed")

// The Window Size is the terminal size maintained by the kernel for the TTY.
// This is distinct from any overrides that might exist.
func GetTerminalWindowSize(file *os.File) (*Size, error) {
	fd := uintptr(file.Fd())
	winsize := C.struct_winsize{}
	winp := uintptr(unsafe.Pointer(&winsize))
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGWINSZ, winp)
	if ep != 0 {
		return nil, ErrGetWinsizeFailed
	}
	return &Size{
		Lines:   int(winsize.ws_row),
		Columns: int(winsize.ws_col),
	}, nil
}
