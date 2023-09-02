package main

import (
	"strings"

	"github.com/gdamore/tcell"
)

var (
	editorIsPromptActive                           = false
	editorPrompt                                   = ""
	editorPromptValue                              = ""
	editorPromptCallbackFn   func([]string)        = nil
	editorPromptCompletionFn func(string) []string = nil
)

func prompt(prompt string, compFn func(string) []string, cbFn func([]string)) {
	editorPrompt = prompt
	editorPromptValue = ""
	editorPromptCallbackFn = cbFn
	editorPromptCompletionFn = compFn
	enterMode("prompt")
}

func noopComplete(prefix string) []string {
	return []string{}
}

func promptUpdateCompletion() {
	// TODO
}

func promptCancel(vt *ViewTree, b *Buffer, kl *KeyList) {
	enterMode("normal")
}

func promptFinish(vt *ViewTree, b *Buffer, kl *KeyList) {
	enterMode("normal")
	// TODO better args parsing
	editorPromptCallbackFn(strings.Split(editorPromptValue, " "))
}

func promptBackspace(vt *ViewTree, b *Buffer, kl *KeyList) {
	if len(editorPromptValue) > 0 {
		editorPromptValue = editorPromptValue[:len(editorPromptValue)-1]
		promptUpdateCompletion()
	}
}

func promptInsert(vt *ViewTree, b *Buffer, kl *KeyList) {
	k := kl.keys[len(kl.keys)-1]
	if k.Key == tcell.KeyRune && k.Mod == 0 {
		editorPromptValue += string(k.Chr)
		promptUpdateCompletion()
	}
}

func promptCommand(vt *ViewTree, b *Buffer, kl *KeyList) {
	prompt(":", func(prefix string) []string {
		// TODO provide command suggestions
		return []string{}
	}, runCommand)
}
