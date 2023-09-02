package main

type ActionType int

const (
	ActionTypeInsert ActionType = 1
	ActionTypeRemove            = -1
)

type Action struct {
	Typ  ActionType
	Loc  *Location
	Data []rune
}

func NewAction(typ ActionType, loc *Location, data []rune) *Action {
	return &Action{Typ: typ, Loc: loc, Data: data}
}

func (a *Action) Apply(b *Buffer) {
	a.Do(b, a.Typ)
	b.Modified = true
	hook_trigger_buffer("modified", b)
}

func (a *Action) Revert(b *Buffer) {
	a.Do(b, -a.Typ)
	b.Modified = true
	hook_trigger_buffer("modified", b)
}

func (a *Action) Do(b *Buffer, typ ActionType) {
	if typ == ActionTypeInsert {
		a.Insert(b)
	} else {
		a.Remove(b)
	}
}

func (a *Action) Insert(b *Buffer) {
	c, l := a.Loc.Char, a.Loc.Line
	for i := len(a.Data) - 1; i >= 0; i-- {
		ch := a.Data[i]
		if ch == '\n' {
			rest := append([]rune(nil), b.Data[l][c:]...)
			b.Data[l] = b.Data[l][:c]
			b.Data = append(b.Data[:l+1],
				append([][]rune{rest}, b.Data[l+1:]...)...)
		} else {
			b.Data[l] = append(b.Data[l][:c],
				append([]rune{ch}, b.Data[l][c:]...)...)
		}
	}
}

func (a *Action) Remove(b *Buffer) {
	n := len(a.Data)
	c, l := a.Loc.Char, a.Loc.Line
	removed := []rune{}
	for i := 0; i < n; i++ {
		removed = append(removed, b.CharAt(l, c))
		if b.CharAt(l, c) == '\n' {
			if len(b.Data)-1 == l {
				a.Data = removed
				return
			}
			b.Data[l] = append(b.Data[l], b.Data[l+1]...)
			b.Data = append(b.Data[:l+1], b.Data[l+2:]...)
		} else {
			b.Data[l] = append(b.Data[l][:c], b.Data[l][c+1:]...)
		}
	}
	a.Data = removed
}
