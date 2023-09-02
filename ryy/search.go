package main

import (
	"regexp"
	"strings"
)

var (
	last_search_buffer    *Buffer = nil
	last_search                   = ""
	last_search_index             = 0
	last_search_highlight         = false
	last_search_results           = []*Location{}
)

func init_search() {
	bind("normal", k("/"), handle_search_start)
	bind("normal", k("N"), handle_search_prev)
	bind("normal", k("n"), handle_search_next)
	bind("normal", k("*"), handle_search_search_work_under_cursor)
	bind("normal", k("SPC n"), func(vt *ViewTree, b *Buffer, kl *KeyList) {
		search_clear()
	})

	addCommand("clearsearch", func(args []string) {
		search_clear()
	})
	addAlias("cs", "clearsearch")

	hook_buffer("modified", func(b *Buffer) {
		// Update search result indexes on buffer changes
		if last_search != "" && b == last_search_buffer {
			search_find_matches(b, last_search)
			last_search_index = len(last_search_results) - 1
		}
	})
}

func search_clear() {
	last_search_highlight = false
	highlight_buffer(currentViewTree.Leaf.Buf)
}

func search_find_matches(b *Buffer, search string) {
	last_search = search
	re := regexp.MustCompile(regexp.QuoteMeta(search))

	last_search_buffer = b
	last_search_results = []*Location{}
	for i, line := range b.Data {
		idxs := re.FindAllStringIndex(string(line), -1)
		for _, idx := range idxs {
			last_search_results = append(last_search_results, NewLocation(i, idx[0]))
		}
	}
}

func search_start(b *Buffer, search string) {
	if len(search) > 0 {
		search_find_matches(b, search)

		// TODO start index at first match after cursor
		last_search_index = len(last_search_results) - 1
		search_next(b)
	}
}

func handle_search_start(vt *ViewTree, b *Buffer, kl *KeyList) {
	prompt("/", noopComplete, func(args []string) {
		search_start(b, strings.Join(args, " "))
	})
}

func search_prev(b *Buffer) {
	if len(last_search_results) == 0 {
		message("No search result.")
		return
	}
	if last_search_index == 0 {
		last_search_index = len(last_search_results) - 1
	} else {
		last_search_index--
	}
	last_search_highlight = true
	highlight_buffer(currentViewTree.Leaf.Buf)
	loc := last_search_results[last_search_index]
	b.MoveTo(loc.Char, loc.Line)

}

func handle_search_prev(vt *ViewTree, b *Buffer, kl *KeyList) {
	search_prev(b)
}

func search_next(b *Buffer) {
	if len(last_search_results) == 0 {
		message("No search result.")
		return
	}
	if last_search_index+1 == len(last_search_results) {
		last_search_index = 0
	} else {
		last_search_index++
	}
	last_search_highlight = true
	highlight_buffer(currentViewTree.Leaf.Buf)
	loc := last_search_results[last_search_index]
	b.MoveTo(loc.Char, loc.Line)

}

func handle_search_next(vt *ViewTree, b *Buffer, kl *KeyList) {
	search_next(b)
}

func search_highlight(b *Buffer, l, c int) int {
	if !last_search_highlight {
		return 0
	}
	if last_search_buffer != b {
		return 0
	}
	for _, loc := range last_search_results {
		if l == loc.Line && c == loc.Char {
			return len(last_search)
		}
	}
	return 0
}

func handle_search_search_work_under_cursor(vt *ViewTree, b *Buffer, kl *KeyList) {
	search_start(b, string(b.WordUnderCursor()))
}
