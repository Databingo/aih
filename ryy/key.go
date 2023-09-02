package main

import (
	"strings"
	"unicode"

	"github.com/gdamore/tcell"
)

type Key struct {
	Mod tcell.ModMask
	Key tcell.Key
	Chr rune
}

const (
	KeyTypeCatchall tcell.Key = iota + 5000
	KeyTypeAlpha
	KeyTypeNum
	KeyTypeAlphaNum
)

func NewKeyFromEvent(ev *tcell.EventKey) *Key {
	k, r, m := ev.Key(), ev.Rune(), ev.Modifiers()

	keyName := ev.Name()
	if strings.HasPrefix(keyName, "Ctrl+") {
		k = tcell.KeyRune
		r = unicode.ToLower([]rune(keyName[5:6])[0])
	}

	// Handle Ctrl-h
	if k == tcell.KeyBackspace {
		m |= tcell.ModCtrl
		k = tcell.KeyRune
		r = 'h'
	}

	if k != tcell.KeyRune {
		r = 0
	}

	return &Key{Mod: ev.Modifiers(), Key: k, Chr: r}
}

func NewKey(rep string) *Key {
	if rep == "$any" {
		return &Key{Key: KeyTypeCatchall}
	} else if rep == "$num" {
		return &Key{Key: KeyTypeNum}
	} else if rep == "$alpha" {
		return &Key{Key: KeyTypeAlpha}
	} else if rep == "$alphanum" {
		return &Key{Key: KeyTypeAlphaNum}
	}

	parts := strings.Split(rep, "-")
	if rep == "-" {
		parts = []string{"-"}
	}

	// Modifiers
	modMask := tcell.ModNone
	for _, part := range parts[:len(parts)-1] {
		switch part {
		case "C":
			modMask |= tcell.ModCtrl
		case "S":
			modMask |= tcell.ModShift
		case "A":
			modMask |= tcell.ModAlt
		case "M":
			modMask |= tcell.ModMeta
		}
	}

	// Key
	var r rune = 0
	var k tcell.Key
	lastPart := parts[len(parts)-1]
	switch lastPart {
	case "DEL":
		k = tcell.KeyDelete
	case "BAK":
		k = tcell.KeyBackspace2
	case "RET":
		k = tcell.KeyEnter
	case "SPC":
		k = tcell.KeyRune
		r = ' '
	case "ESC":
		k = tcell.KeyEscape
	case "TAB":
		k = tcell.KeyTab
	default:
		k = tcell.KeyRune
		r = []rune(lastPart)[0]
	}

	return &Key{Mod: modMask, Key: k, Chr: r}
}

func (k *Key) String() string {
	if k.Key == KeyTypeCatchall {
		return "$any"
	} else if k.Key == KeyTypeNum {
		return "$num"
	} else if k.Key == KeyTypeAlpha {
		return "$alpha"
	} else if k.Key == KeyTypeAlphaNum {
		return "$alphanum"
	}

	mods := []string{}
	if k.Mod&tcell.ModCtrl != 0 {
		mods = append(mods, "C")
	}
	if k.Mod&tcell.ModShift != 0 {
		mods = append(mods, "S")
	}
	if k.Mod&tcell.ModAlt != 0 {
		mods = append(mods, "A")
	}
	if k.Mod&tcell.ModMeta != 0 {
		mods = append(mods, "M")
	}

	name := string(k.Chr)
	switch k.Key {
	case tcell.KeyDelete:
		name = "DEL"
	case tcell.KeyBackspace2:
		name = "BAK"
	case tcell.KeyEnter:
		name = "RET"
	case tcell.KeyEscape:
		name = "ESC"
	case tcell.KeyTab:
		name = "TAB"
	}
	if k.Key == tcell.KeyRune && k.Chr == ' ' {
		name = "SPC"
	}

	return strings.Join(append(mods, name), "-")
}

func (k *Key) IsRune() bool {
	return k.Mod == 0 && k.Key == tcell.KeyRune
}

// TODO implement alphanum match
func (k1 *Key) Matches(k2 *Key) bool {
	if k1.Key == KeyTypeCatchall || k2.Key == KeyTypeCatchall {
		return true
	}
	if k1.Key == KeyTypeAlpha && k2.IsRune() && isAlpha(k2.Chr) {
		return true
	}
	if k2.Key == KeyTypeAlpha && k1.IsRune() && isAlpha(k1.Chr) {
		return true
	}
	if k1.Key == KeyTypeNum && k2.IsRune() && isNum(k2.Chr) {
		return true
	}
	if k2.Key == KeyTypeNum && k1.IsRune() && isNum(k1.Chr) {
		return true
	}
	return k1.Mod == k2.Mod && k1.Key == k2.Key && k1.Chr == k2.Chr
}

type KeyList struct {
	keys []*Key
}

func NewKeyList(rep string) *KeyList {
	kl := &KeyList{[]*Key{}}
	parts := strings.Split(rep, " ")
	for _, part := range parts {
		if part != "" {
			kl.keys = append(kl.keys, NewKey(part))
		}
	}
	return kl
}

var k = NewKeyList

func (kl *KeyList) String() string {
	rep := []string{}
	for _, k := range kl.keys {
		rep = append(rep, k.String())
	}
	return strings.Join(rep, " ")
}

func (kl *KeyList) AddKey(k *Key) {
	kl.keys = append(kl.keys, k)
}

func (kl1 *KeyList) Matches(kl2 *KeyList) bool {
	if len(kl1.keys) != len(kl2.keys) {
		return false
	}
	for i := range kl1.keys {
		if !kl1.keys[i].Matches(kl2.keys[i]) {
			return false
		}
	}
	return true
}

func (kl1 *KeyList) HasSuffix(kl2 *KeyList) *KeyList {
	for i := len(kl1.keys) - 1; i >= 0; i-- {
		tmp_kl := KeyList{kl1.keys[i:]}
		if tmp_kl.Matches(kl2) {
			return &tmp_kl
		}
	}
	return nil
}
