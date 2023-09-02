package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/go-errors/errors"
	runewidth "github.com/mattn/go-runewidth"
)

func fatalError(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func fatal(message string) {
	if screen != nil {
		screen.Fini()
		screen = nil
	}
	fmt.Printf("%v\n", message)
	os.Exit(1)
}

func handlePanics() {
	if err := recover(); err != nil {
		fatal(fmt.Sprintf("ry fatal error:\n%v\n%s", err, errors.Wrap(err, 2).ErrorStack()))
	}
}

func write(style tcell.Style, x, y int, str string) int {
	s := screen
	i := 0
	var deferred []rune
	dwidth := 0
	for _, r := range str {
		// Handle tabs
		if r == '\t' {
			// TODO setting
			tabWidth := int(configGetNumber("tab_width", nil))

			// Print first tab char
			s.SetContent(x+i, y, '>', nil, style.Foreground(tcell.ColorAqua))
			i++

			// Add space till we reach tab column or tabWidth
			for j := 0; j < tabWidth-1 || i%tabWidth == tabWidth-1; j++ {
				s.SetContent(x+i, y, ' ', nil, style)
				i++
			}

			deferred = nil
			continue
		}

		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}

	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}

	// i is the real width of what we just outputed
	return i
}

func listContainsString(list []string, search string) bool {
	for _, item := range list {
		if item == search {
			return true
		}
	}
	return false
}

func padr(str string, length int, padding rune) string {
	for utf8.RuneCountInString(str) < length {
		str = str + string(padding)
	}
	return str
}

func padl(str string, length int, padding rune) string {
	for utf8.RuneCountInString(str) < length {
		str = string(padding) + str
	}
	return str
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isWord(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || strings.ContainsRune("_", r)
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNum(r rune) bool {
	return r >= '0' && r <= '9'
}
