package terminal

import (
	"io"
	"strings"
	"testing"
)

func extractStr(t *State, x0, x1, row int) string {
	var s []rune
	for i := x0; i <= x1; i++ {
		c, _, _ := t.Cell(i, row)
		s = append(s, c)
	}
	return string(s)
}

func TestPlainChars(t *testing.T) {
	var st State
	term, err := Create(&st, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := "Hello world!"
	_, err = term.Write([]byte(expected))
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	actual := extractStr(&st, 0, len(expected)-1, 0)
	if expected != actual {
		t.Fatal(actual)
	}
}

func TestNewline(t *testing.T) {
	var st State
	term, err := Create(&st, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := "Hello world!\n...and more."
	_, err = term.Write([]byte("\033[20h")) // set CRLF mode
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, err = term.Write([]byte(expected))
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	split := strings.Split(expected, "\n")
	actual := extractStr(&st, 0, len(split[0])-1, 0)
	actual += "\n"
	actual += extractStr(&st, 0, len(split[1])-1, 1)
	if expected != actual {
		t.Fatal(actual)
	}

	// A newline with a color set should not make the next line that color,
	// which used to happen if it caused a scroll event.
	st.moveTo(0, st.rows-1)
	_, err = term.Write([]byte("\033[1;37m\n$ \033[m"))
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	_, fg, bg := st.Cell(st.Cursor())
	if fg != DefaultFG {
		t.Fatal(st.cur.x, st.cur.y, fg, bg)
	}
}
