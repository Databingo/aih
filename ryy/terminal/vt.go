package terminal

import (
	"bufio"
	"bytes"
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
	"unicode"
	"unicode/utf8"
)

// VT represents the virtual terminal emulator.
type VT struct {
	dest *State
	rc   io.ReadCloser
	br   *bufio.Reader
	pty  *os.File
}

// Start initializes a virtual terminal emulator with the target state
// and a new pty file by starting the *exec.Command. The returned
// *os.File is the pty file.
func Start(state *State, cmd *exec.Cmd) (*VT, *os.File, error) {
	var err error
	t := &VT{
		dest: state,
	}
	t.pty, err = pty.Start(cmd)
	if err != nil {
		return nil, nil, err
	}
	t.rc = t.pty
	t.init()
	return t, t.pty, nil
}

// Create initializes a virtual terminal emulator with the target state
// and io.ReadCloser input.
func Create(state *State, rc io.ReadCloser) (*VT, error) {
	t := &VT{
		dest: state,
		rc:   rc,
	}
	t.init()
	return t, nil
}

func (t *VT) init() {
	t.br = bufio.NewReader(t.rc)
	t.dest.numlock = true
	t.dest.state = t.dest.parse
	t.dest.cur.attr.fg = DefaultFG
	t.dest.cur.attr.bg = DefaultBG
	t.Resize(80, 24)
	t.dest.reset()
}

// File returns the pty file.
func (t *VT) File() *os.File {
	return t.pty
}

// Write parses input and writes terminal changes to state.
func (t *VT) Write(p []byte) (int, error) {
	var written int
	r := bytes.NewReader(p)
	t.dest.lock()
	defer t.dest.unlock()
	for {
		c, sz, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return written, err
		}
		written += sz
		if c == unicode.ReplacementChar && sz == 1 {
			if r.Len() == 0 {
				// not enough bytes for a full rune
				return written - 1, nil
			}
			t.dest.logln("invalid utf8 sequence")
			continue
		}
		t.dest.put(c)
	}
	return written, nil
}

// Close closes the pty or io.ReadCloser.
func (t *VT) Close() error {
	return t.rc.Close()
}

// Parse blocks on read on pty or io.ReadCloser, then parses sequences until
// buffer empties. State is locked as soon as first rune is read, and unlocked
// when buffer is empty.
// TODO: add tests for expected blocking behavior
func (t *VT) Parse() error {
	var locked bool
	defer func() {
		if locked {
			t.dest.unlock()
		}
	}()
	for {
		c, sz, err := t.br.ReadRune()
		if err != nil {
			return err
		}
		if c == unicode.ReplacementChar && sz == 1 {
			t.dest.logln("invalid utf8 sequence")
			break
		}
		if !locked {
			t.dest.lock()
			locked = true
		}

		// put rune for parsing and update state
		t.dest.put(c)

		// break if our buffer is empty, or if buffer contains an
		// incomplete rune.
		n := t.br.Buffered()
		if n == 0 || (n < 4 && !fullRuneBuffered(t.br)) {
			break
		}
	}
	return nil
}

func fullRuneBuffered(br *bufio.Reader) bool {
	n := br.Buffered()
	buf, err := br.Peek(n)
	if err != nil {
		return false
	}
	return utf8.FullRune(buf)
}

// Resize reports new size to pty and updates state.
func (t *VT) Resize(cols, rows int) {
	t.dest.lock()
	defer t.dest.unlock()
	_ = t.dest.resize(cols, rows)
	t.ptyResize()
}
