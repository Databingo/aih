package main

func initVisual() {
	addMode("visual")
	bind("visual", k("ESC"), exitVisualMode)
	bind("visual", k("y"), visualModeYank)
	bind("visual", k("d"), visualModeDelete)
	bind("visual", k("p"), visualModePaste)
	bind("visual", k("c"), visualModeChange)

	addMode("visual-line")
	bind("visual-line", k("ESC"), exitVisualMode)
	bind("visual-line", k("y"), visualModeYank)
	bind("visual-line", k("d"), visualModeDelete)
	bind("visual-line", k("p"), visualModePaste)
	bind("visual-line", k("c"), visualModeChange)

	hook_buffer("moved", visualRehighlight)
}

// Run highlight when in visual mode and cursor mode as normal we
// only recompute highlights when the buffer changes
func visualRehighlight(b *Buffer) {
	if b.IsInMode("visual") || b.IsInMode("visual-line") {
		highlight_buffer(b)
	}
}

func visualHighlight(b *Buffer, l, c int) bool {
	in_visual_line := b.IsInMode("visual-line")
	if !b.IsInMode("visual") && !in_visual_line {
		return false
	}

	l1, l2 := orderLocations(b.Cursor, getMark('∫').Loc)

	if in_visual_line {
		// compare using line numbers
		return l1.Line <= l && l <= l2.Line
	} else {
		// compare using line numbers + char position
		loc := NewLocation(l, c)
		return loc.After(l1) && loc.Before(l2) || loc.Equal(l1) || loc.Equal(l2)
	}
}

func exitVisualMode(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.RemoveMode("visual")
	b.RemoveMode("visual-line")
	highlight_buffer(b)
}

func enterVisualMode(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.AddMode("visual")
	markCreate('∫', b)
}

func enterVisualBlockMode(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.AddMode("visual-line")
	markCreate('∫', b)
	highlight_buffer(b)
}

func visualModeSelection(b *Buffer) ([]rune, *Location, *Location) {
	in_visual_line := b.IsInMode("visual-line")
	l1, l2 := orderLocations(b.Cursor, getMark('∫').Loc)
	data := []rune{}
	if in_visual_line {
		l1.Char = 0
		l2.Char = len(b.GetLine(l2.Line))
	}
	for l := l1.Line; l <= l2.Line; l++ {
		start_char := 0
		if l == l1.Line {
			start_char = l1.Char
		}
		line_data := b.GetLine(l)
		end_char := len(line_data)
		if l == l2.Line {
			end_char = min(l2.Char+1, len(line_data))
		}
		data = append(data, line_data[start_char:end_char]...)
		if l != l2.Line || end_char == len(line_data) {
			data = append(data, '\n')
		}
	}
	return data, l1, l2
}

func visualModeYank(vt *ViewTree, b *Buffer, kl *KeyList) {
	text, l1, _ := visualModeSelection(b)
	clipboardSet(defaultClipboard, text)

	b.MoveTo(l1.Char, l1.Line)
	exitVisualMode(vt, b, kl)
}
func visualModeDelete(vt *ViewTree, b *Buffer, kl *KeyList) {
	text, l1, _ := visualModeSelection(b)
	b.MoveTo(l1.Char, l1.Line)
	b.Remove(len(text))

	exitVisualMode(vt, b, kl)
}
func visualModePaste(vt *ViewTree, b *Buffer, kl *KeyList) {
	text, l1, _ := visualModeSelection(b)
	clipboard_text := clipboardGet(defaultClipboard)
	b.MoveTo(l1.Char, l1.Line)
	b.Remove(len(text))
	b.Insert(clipboard_text)
	clipboardSet(defaultClipboard, text)

	b.MoveTo(l1.Char, l1.Line)
	exitVisualMode(vt, b, kl)
}
func visualModeChange(vt *ViewTree, b *Buffer, kl *KeyList) {
	text, l1, _ := visualModeSelection(b)
	b.MoveTo(l1.Char, l1.Line)
	b.Remove(len(text))

	exitVisualMode(vt, b, kl)
	enterInsertMode(vt, b, kl)
}
