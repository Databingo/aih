package main

import (
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	fatalError(err)
	err = screen.Init()
	fatalError(err)

	encoding.Register()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	editorWidth, editorHeight = screen.Size()
}

func initTermEvents() {
	go func() {
		for {
			if screen == nil {
				break
			}
			termEvents <- screen.PollEvent()
		}
	}()
}

func initBuffers() {
	for _, arg := range os.Args[1:] {
		openBufferFromFile(arg)
	}
	if len(buffers) == 0 {
		openBufferNamed("*scratch*")
	}
}

func initViews() {
	view := NewView(buffers[0])
	rootViewTree = &ViewTree{Leaf: view}
	currentViewTree = rootViewTree
}
