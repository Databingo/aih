package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

const (
	specialChars = "[]{}()/\\"
)

var (
	keysEntered                    = NewKeyList("")
	lastKey                        = NewKeyList("")
	termEvents                     = make(chan tcell.Event, 500)
	defaultClipboard               = '_'
	clipboards                     = map[rune][]rune{'_': []rune{}}
	editorMode                     = "normal"
	editorMessage                  = ""
	editorMessageType              = "info"
	editorWidth                    = 0
	editorHeight                   = 0
	buffers                        = []*Buffer{}
	screen            tcell.Screen = nil
	rootViewTree      *ViewTree    = nil
	currentViewTree   *ViewTree    = nil
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		fmt.Println("ry v0.0.0")
		os.Exit(0)
	}

	defer handlePanics()

	initModes()
	initCommands()

	initConfig()
	init_hooks()
	init_highlighting()
	init_search()
	initVisual()
	initTerm()

	initScreen()
	initTermEvents()
	initBuffers()
	initViews()

	render()

top:
	for {
		select {
		case ev := <-termEvents:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlQ {
					screen.Fini()
					screen = nil
					break top
					/*
						} else if ev.Key() == tcell.KeyEscape {
							kl := k("ESC")
							enter_normal_mode(current_view_tree, current_view_tree.leaf.buf, kl)
							last_key = kl
							keys_entered = k("")
					*/
				} else {
					keysEntered.AddKey(NewKeyFromEvent(ev))

					buf := currentViewTree.Leaf.Buf
					for _, mode_name := range buf.Modes {
						if matched := modeHandle(mustFindMode(mode_name), keysEntered); matched != nil {
							keysEntered = k("")
							lastKey = matched
							continue top
						}
					}
					if matched := modeHandle(mustFindMode(editorMode), keysEntered); matched != nil {
						keysEntered = k("")
						lastKey = matched
						continue top
					}
				}
			case *tcell.EventResize:
				editorWidth, editorHeight = screen.Size()
			}
		default:
			render()
		}
	}
}
