package main

import (
	"math"

	"github.com/gdamore/tcell"
)

type ViewHighlight struct {
	Beg   *Location
	End   *Location
	Style tcell.Style
}

type View struct {
	Buf           *Buffer
	LineOffset    int
	CenterPending bool

	Highlights []*ViewHighlight
}

func NewView(buf *Buffer) *View {
	return &View{
		Buf:           buf,
		LineOffset:    0,
		CenterPending: false,
		Highlights:    []*ViewHighlight{},
	}
}

func (v *View) AdjustScroll(w, h int) {
	l := v.Buf.Cursor.Line
	if v.CenterPending {
		v.LineOffset = max(l-int(math.Floor(float64(h-1)/2)), 1)
		v.CenterPending = false
		return
	}
	// too low
	// (h-2) as height includes status bar and moving to 0 based
	if l > h-2+v.LineOffset {
		v.LineOffset = max(l-h+2, 0)
	}
	// too high
	if l < v.LineOffset {
		v.LineOffset = l
	}
}

// }}}

// {{{ ViewTree
type ViewTree struct {
	Parent *ViewTree
	Left   *ViewTree
	Right  *ViewTree
	Top    *ViewTree
	Bottom *ViewTree
	Leaf   *View
	Size   int
}

func NewViewTreeLeaf(parent *ViewTree, v *View) *ViewTree {
	return &ViewTree{Parent: parent, Leaf: v, Size: 50}
}

// }}}

// {{{ message
func message(m string) {
	editorMessage = m
	editorMessageType = "info"
}

func messageError(m string) {
	editorMessage = m
	editorMessageType = "error"
}
