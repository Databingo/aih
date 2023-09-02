package main

func initTerm() {
	addMode("term")
	bind("term", k("ESC"), termExitMode)
	bind("term", k("$any"), termInput)
}

func termExitMode(vt *ViewTree, b *Buffer, kl *KeyList) {
	b.RemoveMode("term")
}

func termInput(vt *ViewTree, b *Buffer, kl *KeyList) {
}
