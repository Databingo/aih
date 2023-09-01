package main

import (
	"fmt"
	"strconv"
)

func render() {
	width, height := editorWidth, editorHeight

	screen.Clear()

	renderViewTree(rootViewTree, 0, 0, width, height-1)

	renderMessageBar(width, height)

	screen.Show()
}

func renderMessageBar(width, height int) {
	s := style("default")

	if editorMode == "prompt" {
		p := editorPrompt + editorPromptValue
		write(s, 0, height-1, p)
		write(s.Reverse(true), len(p), height-1, " ")
		return
	}

	smb := s
	if editorMessageType == "error" {
		smb = style("message.error")
	}
	if editorMessage != "" {
		write(smb, 0, height-1, editorMessage)
	} else {
		write(smb, 0, height-1, keysEntered.String())
	}
	lastKeyText := lastKey.String()
	write(s, width-len(lastKeyText)-1, height-1, lastKeyText)
}

func renderViewTree(vt *ViewTree, x, y, w, h int) {
	if vt.Leaf != nil {
		renderView(vt.Leaf, x, y, w, h)
		return
	}
	panic("unreachable")
}

func renderView(v *View, x, y, w, h int) {
	sc := style("cursor")
	sln := style("linenumber")
	ssb := style("statusbar")
	ssbh := style("statusbar.highlight")
	b := v.Buf

	b.LastRenderWidth = w
	b.LastRenderHeight = h

	styleMap := highlighting_styles(b)

	gutterw := len(strconv.Itoa(len(b.Data))) + 1
	sy := y
	line := v.LineOffset
	for line < len(b.Data) && sy < y+h-1 {
		write(sln, x, sy, padl(strconv.Itoa(line+1), gutterw-1, ' '))

		sx := x + gutterw
		for c, char := range b.Data[line] {
			if v == currentViewTree.Leaf && line == b.Cursor.Line && c == b.Cursor.Char {
				sx += write(sc, sx, sy, string(char))
			} else {
				sx += write(styleMap[line][c], sx, sy, string(char))
			}
			if sx >= x+w {
				break
			}
		}
		if v == currentViewTree.Leaf &&
			line == b.Cursor.Line &&
			b.Cursor.Char == len(b.Data[b.Cursor.Line]) {
			write(sc, sx, sy, " ")
		}

		line++
		sy++
	}

	// Current mode
	modeStatus := editorMode
	for _, modeName := range b.Modes {
		modeStatus += "+" + modeName
	}
	modeStatus = " " + modeStatus + " "
	write(ssbh, x, y+h-1, modeStatus)

	// Position
	statusRight := fmt.Sprintf("(%d,%d) %d ", b.Cursor.Char+1, b.Cursor.Line+1, len(b.Data))
	write(ssb, x+w-len(statusRight), y+h-1, statusRight)
	// File name
	statusLeft := " " + b.Name
	if b.Modified {
		statusLeft += " [+]"
	}
	write(ssb, x+len(modeStatus), y+h-1, padr(statusLeft, w-len(statusRight)-len(modeStatus), ' '))
}
