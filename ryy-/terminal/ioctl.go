package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

func ioctl(f *os.File, cmd, p uintptr) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		syscall.TIOCSWINSZ,
		p)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (t *VT) ptyResize() error {
	if t.pty == nil {
		return nil
	}
	var w struct{ row, col, xpix, ypix uint16 }
	w.row = uint16(t.dest.rows)
	w.col = uint16(t.dest.cols)
	w.xpix = 16 * uint16(t.dest.cols)
	w.ypix = 16 * uint16(t.dest.rows)
	return ioctl(t.pty, syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&w)))
}
