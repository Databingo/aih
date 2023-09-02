package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Location struct {
	Line int
	Char int
}

func NewLocation(l, c int) *Location {
	return &Location{Line: l, Char: c}
}

func orderLocations(l1, l2 *Location) (*Location, *Location) {
	if l1.Line < l2.Line {
		return l1, l2
	} else if l1.Line > l2.Line {
		return l2, l1
	} else {
		if l1.Char > l2.Char {
			return l2, l1
		} else {
			return l1, l2
		}
	}
}

func (l1 *Location) Equal(l2 *Location) bool {
	return l1.Line == l2.Line && l1.Char == l2.Char
}

func (l1 *Location) Before(l2 *Location) bool {
	ol1, _ := orderLocations(l1, l2)
	return l1 == ol1
}

func (l1 *Location) After(l2 *Location) bool {
	_, ol2 := orderLocations(l1, l2)
	return l1 == ol2
}

func (loc *Location) Clone() *Location {
	return &Location{Line: loc.Line, Char: loc.Char}
}

type CharRange struct {
	Beg int
	End int
}

func NewCharRange(b, e int) *CharRange {
	return &CharRange{b, e}
}

type Buffer struct {
	Data             [][]rune
	History          []*Action
	HistoryIndex     int
	Name             string
	Path             string
	Modified         bool
	Cursor           *Location
	Modes            []string
	LastRenderWidth  int
	LastRenderHeight int
}

func NewBuffer(name string, path string) *Buffer {
	b := &Buffer{
		Data:         [][]rune{{}},
		History:      []*Action{},
		HistoryIndex: -1,
		Modified:     false,
		Cursor:       NewLocation(0, 0),
		Modes:        []string{},
	}

	if path == "" {
		// TODO ensure uniqueness
		b.Name = name
	} else {
		b.SetPath(path)
		// TODO ensure uniqueness
		b.Name = name
	}

	return b
}

func (b *Buffer) IsInMode(name string) bool {
	for _, n := range b.Modes {
		if n == name {
			return true
		}
	}
	return false
}

func (b *Buffer) CharAt(l, c int) rune {
	line := b.Data[l]
	if c < 0 {
		return rune(0)
	} else if c < len(line) {
		return line[c]
	} else {
		return '\n'
	}
}

func (b *Buffer) GetLine(l int) []rune {
	return b.Data[l]
}

func (b *Buffer) CharAtLeft() rune {
	return b.CharAt(b.Cursor.Line, b.Cursor.Char-1)
}

func (b *Buffer) CharUnderCursor() rune {
	return b.CharAt(b.Cursor.Line, b.Cursor.Char)
}

func (b *Buffer) WordUnderCursor() []rune {
	ch := b.CharUnderCursor()
	if !isWord(ch) {
		return nil
	}
	line := b.GetLine(b.Cursor.Line)
	cstart := b.Cursor.Char
	for cstart > 0 && isWord(b.CharAt(b.Cursor.Line, cstart-1)) {
		cstart--
	}
	cend := cstart
	for cend < len(line) && isWord(b.CharAt(b.Cursor.Line, cend+1)) {
		cend++
	}
	return line[cstart : cend+1]
}

func (b *Buffer) FirstLine() bool {
	return b.Cursor.Line == 0
}

func (b *Buffer) LastLine() bool {
	return b.Cursor.Line == len(b.Data)-1
}

func (b *Buffer) MoveTo(c, l int) {
	b.Cursor.Line = max(min(l, len(b.Data)-1), 0)
	b.Cursor.Char = max(min(c, len(b.Data[b.Cursor.Line])), 0)
	hook_trigger_buffer("moved", b)
}

func (b *Buffer) Move(c, l int) {
	b.MoveTo(b.Cursor.Char+c, b.Cursor.Line+l)
}

func (b *Buffer) MoveWordForward() bool {
	for {
		c := b.CharUnderCursor()
		if c == '\n' {
			if b.LastLine() {
				return false
			} else {
				b.MoveTo(0, b.Cursor.Line+1)
				break
			}
		}

		for isWord(c) && c != '\n' {
			b.Move(1, 0)
			c = b.CharUnderCursor()
		}

		if c == '\n' {
			continue
		}
		break
	}

	c := b.CharUnderCursor()
	for !isWord(c) && c != '\n' {
		b.Move(1, 0)
		c = b.CharUnderCursor()
	}

	return true
}

func (b *Buffer) MoveWordEndForward() bool {
	b.Move(1, 0)
	for {
		c := b.CharUnderCursor()
		if c == '\n' {
			if b.LastLine() {
				return false
			} else {
				b.MoveTo(0, b.Cursor.Line+1)
				break
			}
		}

		for !isWord(c) && c != '\n' {
			b.Move(1, 0)
			c = b.CharUnderCursor()
		}

		if c == '\n' {
			continue
		}
		break
	}

	c := b.CharUnderCursor()
	for isWord(c) && c != '\n' {
		b.Move(1, 0)
		c = b.CharUnderCursor()
	}
	b.Move(-1, 0)

	return true
}

func (b *Buffer) MoveWordBackward() bool {
	b.Move(-1, 0)
	for {
		c := b.CharUnderCursor()
		if b.Cursor.Char == 0 {
			if b.FirstLine() {
				return false
			} else {
				b.MoveTo(len(b.Data[b.Cursor.Line-1]), b.Cursor.Line-1)
				continue
			}
		}

		for !isWord(c) && b.Cursor.Char != 0 {
			b.Move(-1, 0)
			c = b.CharUnderCursor()
		}

		if b.Cursor.Char == 0 {
			continue
		}
		break
	}

	c := b.CharUnderCursor()
	for isWord(c) && b.Cursor.Char != 0 {
		b.Move(-1, 0)
		c = b.CharUnderCursor()
	}
	b.Move(1, 0)

	return true
}

func (b *Buffer) Insert(data []rune) {
	a := NewAction(ActionTypeInsert, b.Cursor.Clone(), data)
	b.HistoryIndex++
	b.History = tryMergeHistory(b.History[:b.HistoryIndex], a)
	a.Apply(b)
}

func (b *Buffer) RemoveAt(loc *Location, n int) []rune {
	a := NewAction(ActionTypeRemove, loc.Clone(), make([]rune, n))
	b.HistoryIndex++
	b.History = tryMergeHistory(b.History[:b.HistoryIndex], a)
	a.Apply(b)
	return a.Data
}

func (b *Buffer) Remove(n int) []rune {
	return b.RemoveAt(b.Cursor, n)
}

func (b *Buffer) Undo() {
	if b.HistoryIndex >= 0 && b.HistoryIndex != -1 {
		b.History[b.HistoryIndex].Revert(b)
		b.HistoryIndex--
	} else {
		message("Noting to undo!")
	}
}

func (b *Buffer) Redo() {
	if b.HistoryIndex+1 < len(b.History) && len(b.History) > b.HistoryIndex+1 {
		b.HistoryIndex++
		b.History[b.HistoryIndex].Apply(b)
	} else {
		message("Noting to redo!")
	}
}

func (b *Buffer) SetPath(path string) {
	var err error
	b.Path, err = filepath.Abs(path)
	if err != nil {
		b.Path = filepath.Clean(path)
	}
	b.Name = ""
	name := filepath.Base(b.Path)

	i := 1
checkName:
	for _, b2 := range buffers {
		if b2.Name == name {
			b.Name = name + " " + strconv.Itoa(i)
			i++
			goto checkName
		}
	}
	if b.Name == "" {
		b.Name = name
	}
}

func (b *Buffer) AddMode(name string) {
	if b.IsInMode(name) {
		return
	}
	b.Modes = append(b.Modes, name)
}

func (b *Buffer) RemoveMode(name string) {
	for i, n := range b.Modes {
		if n == name {
			b.Modes = append(b.Modes[:i], b.Modes[i+1:]...)
			return
		}
	}
}

func (b *Buffer) Contents() string {
	ret := ""
	for _, line := range b.Data {
		ret += string(line) + "\n"
	}
	return ret
}

func (b *Buffer) NicePath() string {
	return strings.Replace(b.Path, os.Getenv("HOME"), "~", -1)
}

func (b *Buffer) Save() {
	if b.Path == "" {
		messageError("Can't save a buffer without a path.")
		return
	}
	err := ioutil.WriteFile(b.Path, []byte(b.Contents()), 0666)
	if err != nil {
		messageError("Error saving buffer: " + err.Error())
	} else {
		b.Modified = false
		message("Buffer written to '" + b.NicePath() + "'")
	}
}

func tryMergeHistory(al []*Action, a *Action) []*Action {
	// TODO save end location on actions so that we can merge them here
	return append(al, a)
}
