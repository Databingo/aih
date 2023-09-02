package main

type Mark struct {
	Loc        *Location
	BufferName string
}

var marks = map[rune]*Mark{}

func markCreate(markLetter rune, b *Buffer) *Mark {
	m := &Mark{Loc: b.Cursor.Clone(), BufferName: b.Name}
	marks[markLetter] = m
	return m
}

func getMark(markLetter rune) *Mark {
	return marks[markLetter]
}

func markJump(markLetter rune) {
	if m, ok := marks[markLetter]; ok {
		if b := showBuffer(m.BufferName); b != nil {
			b.MoveTo(m.Loc.Char, m.Loc.Line)
		} else {
			messageError("Can't find buffer named '" + m.BufferName + "'")
		}
	}
}
