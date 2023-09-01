package main

import (
	zclip "github.com/zyedidia/clipboard"
)

func clipboardGet(register rune) []rune {
	if register == defaultClipboard {
		if value, err := zclip.ReadAll("clipboard"); err == nil {
			return []rune(value)
		}
	}
	if value, ok := clipboards[register]; ok {
		return value
	}
	return []rune{}
}

func clipboardSet(register rune, value []rune) {
	if register == defaultClipboard {
		if err := zclip.WriteAll(string(value), "clipboard"); err != nil {
			messageError("Error clipboard_get: " + err.Error())
			clipboards[register] = value
		}
	} else {
		clipboards[register] = value
	}
}
