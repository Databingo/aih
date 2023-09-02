package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
)

type CommandFn func(*ViewTree, *Buffer, *KeyList)

type ModeBinding struct {
	k *KeyList
	f CommandFn
}

type Mode struct {
	name     string
	bindings []*ModeBinding
}

var modes = map[string]*Mode{}

func modeHandle(m *Mode, kl *KeyList) *KeyList {
	var match *KeyList = nil
	var matchBinding *ModeBinding = nil
	for _, binding := range m.bindings {
		if matched := kl.HasSuffix(binding.k); matched != nil {
			if match == nil || len(matched.keys) > len(match.keys) {
				matchBinding = binding
				match = matched
			}
		}
	}
	if match != nil {
		matchBinding.f(currentViewTree, currentViewTree.Leaf.Buf, match)
		return match
	}
	return nil
}

func findMode(name string) *Mode {
	if m, ok := modes[name]; ok {
		return m
	} else {
		return nil
	}
}

func mustFindMode(name string) *Mode {
	mode := findMode(name)
	if mode == nil {
		panic(fmt.Sprintf("no mode named '%s'", name))
	}
	return mode
}

// Adds a new empty mode to the mode list, if not already present
func addMode(name string) {
	if _, ok := modes[name]; !ok {
		modes[name] = &Mode{name: name, bindings: []*ModeBinding{}}
	}
}

func bind(mode_name string, k *KeyList, f CommandFn) {
	mode := mustFindMode(mode_name)

	// If this key is bound, update bound function
	for _, binding := range mode.bindings {
		if k.String() == binding.k.String() {
			binding.f = f
		}
	}
	// Else, it's a new binding, add it
	mode.bindings = append(mode.bindings, &ModeBinding{k: k, f: f})
}

func initModes() {
	addMode("normal")
	bind("normal", k("m $alpha"), commandMark)
	bind("normal", k("' $alpha"), commandMoveToMark)
	bind("normal", k(":"), promptCommand)
	bind("normal", k("h"), moveLeft)
	bind("normal", k("j"), moveDown)
	bind("normal", k("k"), moveUp)
	bind("normal", k("l"), moveRight)
	bind("normal", k("0"), moveLineBeg)
	bind("normal", k("$"), moveLineEnd)
	bind("normal", k("g g"), moveTop)
	bind("normal", k("G"), moveBottom)
	bind("normal", k("C-u"), moveJumpUp)
	bind("normal", k("C-d"), moveJumpDown)
	bind("normal", k("z z"), moveCenterLine)
	bind("normal", k("w"), moveWordForward)
	bind("normal", k("e"), moveWordEndForward)
	bind("normal", k("b"), moveWordBackward)
	bind("normal", k("C-c"), cancelKeysEntered)
	bind("normal", k("C-g"), cancelKeysEntered)
	bind("normal", k("ESC ESC"), cancelKeysEntered)
	bind("normal", k("i"), enterInsertMode)
	bind("normal", k("a"), enterInsertModeAppend)
	bind("normal", k("A"), enterInsertModeEol)
	bind("normal", k("o"), enterInsertModeNl)
	bind("normal", k("O"), enterInsertModeNlUp)
	bind("normal", k("x"), removeChar)
	bind("normal", k("d d"), removeLine)
	bind("normal", k("y y"), commandCopyLine)
	bind("normal", k("p"), commandPaste)
	bind("normal", k("u"), commandUndo)
	bind("normal", k("C-r"), commandRedo)
	bind("normal", k("v"), enterVisualMode)
	bind("normal", k("V"), enterVisualBlockMode)

	addMode("insert")
	bind("insert", k("ESC"), enterNormalMode)
	bind("insert", k("C-c"), enterNormalMode)
	bind("insert", k("RET"), insertEnter)
	bind("insert", k("BAK"), insertBackspace)
	bind("insert", k("$any"), insert)

	addMode("prompt")
	bind("prompt", k("C-c"), promptCancel)
	bind("prompt", k("C-g"), promptCancel)
	bind("prompt", k("ESC"), promptCancel)
	bind("prompt", k("RET"), promptFinish)
	bind("prompt", k("BAK"), promptBackspace)
	bind("prompt", k("$any"), promptInsert)

	addMode("buffers")
	bind("buffers", k("q"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		closeCurrentBuffer(true)
	})
	bind("buffers", k("RET"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		closeCurrentBuffer(true)
		showBuffer(string(b.Data[b.Cursor.Line]))
	})

	addMode("directory")
	bind("directory", k("q"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		closeCurrentBuffer(true)
	})
	bind("directory", k("RET"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		closeBuffer(currentViewTree.Leaf.Buf)
		selectAvailableBuffer(false)
		file_path := filepath.Join(b.Path, string(b.Data[b.Cursor.Line]))
		runCommand([]string{"edit", file_path})
	})

	// TODO Remove once we have user configurable bindings
	bind("normal", k("SPC b"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		runCommand([]string{"buffers"})
	})
	editCurrentFolder := func(vt *ViewTree, b *Buffer, kl *KeyList) {
		if b.Path == "" {
			runCommand([]string{"edit", "."})
		} else {
			runCommand([]string{"edit", filepath.Dir(b.Path)})
		}
	}
	bind("normal", k("SPC f"), editCurrentFolder)
	bind("normal", k("-"), editCurrentFolder)
	bind("normal", k("SPC n"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		runCommand([]string{"clearsearch"})
	})
}

func moveLeft(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Move(-1, 0)
}
func moveRight(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Move(1, 0)
}
func moveUp(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Move(0, -1)
}
func moveDown(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Move(0, 1)
}
func moveLineBeg(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveTo(0, b.Cursor.Line)
}
func moveLineEnd(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveTo(len(b.Data[b.Cursor.Line]), b.Cursor.Line)
}
func moveTop(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveTo(0, 0)
}
func moveBottom(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveTo(0, len(b.Data)-1)
}
func moveJumpUp(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Move(0, -15)
}
func moveJumpDown(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Move(0, 15)
}
func moveCenterLine(vt *ViewTree, b *Buffer, kl *KeyList) {
	vt.Leaf.CenterPending = true
}
func moveWordBackward(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveWordBackward()
}
func moveWordForward(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveWordForward()
}
func moveWordEndForward(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.MoveWordEndForward()
}

func cancelKeysEntered(vt *ViewTree, b *Buffer, kl *KeyList) {
	keysEntered = k("")
}

// Enter in a new mode
func enterMode(mode string) {
	editorMode = mode
	// TODO maybe not the best place to clear this
	message("")
}

func enterNormalMode(vt *ViewTree, b *Buffer, kl *KeyList) {
	moveLeft(vt, b, kl)
	enterMode("normal")
}
func enterInsertMode(vt *ViewTree, b *Buffer, kl *KeyList) {
	enterMode("insert")
}
func enterInsertModeAppend(vt *ViewTree, b *Buffer, kl *KeyList) {
	moveRight(vt, b, kl)
	enterMode("insert")
}
func enterInsertModeEol(vt *ViewTree, b *Buffer, kl *KeyList) {
	moveLineEnd(vt, b, kl)
	enterMode("insert")
}
func enterInsertModeNl(vt *ViewTree, b *Buffer, kl *KeyList) {
	moveLineEnd(vt, b, kl)
	b.Insert([]rune("\n"))
	b.Move(0, 1) // ensure a valid position
	enterMode("insert")
}
func enterInsertModeNlUp(vt *ViewTree, b *Buffer, kl *KeyList) {
	moveLineBeg(vt, b, kl)
	b.Insert([]rune("\n"))
	b.Move(0, 0) // ensure a valid position
	enterMode("insert")
}

func insertEnter(vt *ViewTree, b *Buffer, kl *KeyList) {
	i := 0
	for ; i < len(b.Data[b.Cursor.Line]) && isSpace(b.Data[b.Cursor.Line][i]); i++ {
	}
	b.Insert([]rune("\n" + strings.Repeat(" ", i)))
	b.MoveTo(i, b.Cursor.Line+1)
}
func insertBackspace(vt *ViewTree, b *Buffer, kl *KeyList) {
	if b.Cursor.Char == 0 {
		if b.Cursor.Line != 0 {
			moveUp(vt, b, kl)
			moveLineEnd(vt, b, kl)
			b.Remove(1)
		}
	} else {
		if b.CharAtLeft() != ' ' {
			b.Move(-1, 0)
			b.Remove(1)
			return
		}
		// handle spaces
		delete_n := 1
		if configGetBool("tab_to_spaces", b) {
			delete_n = int(configGetNumber("tab_width", b))
		}
		for i := 0; i < delete_n && b.CharAtLeft() == ' '; i++ {
			b.Move(-1, 0)
			b.Remove(1)
		}
	}
}
func insert(vt *ViewTree, b *Buffer, kl *KeyList) {
	k := kl.keys[len(kl.keys)-1]
	if k.Key == tcell.KeyTab {
		if configGetBool("tab_to_spaces", b) {
			tabWidth := int(configGetNumber("tab_width", b))
			b.Insert([]rune(strings.Repeat(" ", tabWidth)))
			b.Move(tabWidth, 0)
		} else {
			b.Insert([]rune{'\t'})
			b.Move(1, 0)
		}
	} else if k.Key == tcell.KeyRune && k.Mod == 0 {
		b.Insert([]rune{k.Chr})
		moveRight(vt, b, kl)
	} else {
		message("Can't insert '" + kl.String() + "'")
	}
}

func removeChar(vt *ViewTree, b *Buffer, kl *KeyList) {
	removed := b.Remove(1)
	clipboardSet(defaultClipboard, removed)
}
func removeLine(vt *ViewTree, b *Buffer, kl *KeyList) {
	moveLineBeg(vt, b, kl)
	removed := b.Remove(len(b.Data[b.Cursor.Line]) + 1)
	clipboardSet(defaultClipboard, removed)
}

func commandUndo(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Undo()
}
func commandRedo(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.Redo()
}
func commandCopyLine(vt *ViewTree, b *Buffer, kl *KeyList) {
	value := make([]rune, len(b.Data[b.Cursor.Line])+1)
	copy(value, b.Data[b.Cursor.Line])
	value[len(value)-1] = '\n'
	clipboardSet(defaultClipboard, value)
}
func commandPaste(vt *ViewTree, b *Buffer, kl *KeyList) {
	value := clipboardGet(defaultClipboard)
	if len(value) == 0 {
		message("Nothing to paste!")
		return
	}
	b.Move(1, 0)
	b.Insert(value)
}

func commandMark(vt *ViewTree, b *Buffer, kl *KeyList) {
	mark_letter := kl.keys[len(kl.keys)-1].Chr
	markCreate(mark_letter, b)
}
func commandMoveToMark(vt *ViewTree, b *Buffer, kl *KeyList) {
	mark_letter := kl.keys[len(kl.keys)-1].Chr
	markJump(mark_letter)
}
