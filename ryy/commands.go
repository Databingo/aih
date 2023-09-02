package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
//	"fmt"
)

func openBufferFromFile(path string) *Buffer {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) || err != nil {
		return openBufferNamed(filepath.Base(path))
	} else {
		if fileInfo.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				messageError("Error opening directory: " + err.Error())
				return nil
			}
			files, err := file.Readdir(0)
			if err != nil {
				messageError("Error opening directory: " + err.Error())
				return nil
			}
			buf := NewBuffer(filepath.Base(path), path)
			buf.Data = [][]rune{}
			file_names := []string{}
			for _, file_info := range files {
				if file_info.IsDir() {
					file_names = append(file_names, " "+file_info.Name())
				} else {
					file_names = append(file_names, file_info.Name())
				}
			}
			file_names = append(file_names, " ..")
			sort.Strings(file_names)
			for _, file_name := range file_names {
				if file_name[0] == ' ' { // is dir
					buf.Data = append(buf.Data, []rune(file_name[1:]+"/"))
				} else {
					buf.Data = append(buf.Data, []rune(file_name))
				}
			}
			buf.AddMode("directory")
			buffers = append(buffers, buf)
			hook_trigger_buffer("modified", buf)
			return buf
		}
	}

	buf := NewBuffer(filepath.Base(path), path)
	if buf.Path != "" {
		contents, err := ioutil.ReadFile(buf.Path)
		if err != nil {
			messageError("Error reading file '" + buf.NicePath() + "'")
			return nil
		}
		buf.Data = [][]rune{}
		for _, line := range strings.Split(string(contents), "\n") {
			buf.Data = append(buf.Data, []rune(line))
		}
		if len(buf.Data) > 1 {
			buf.Data = buf.Data[:len(buf.Data)-1]
		}
	}
	buffers = append(buffers, buf)
	hook_trigger_buffer("modified", buf)
	return buf
}

func openBufferNamed(name string) *Buffer {
	buf := NewBuffer(name, "")
	buffers = append(buffers, buf)
	hook_trigger_buffer("modified", buf)
	return buf
}

// TODO check if is shown first, then if not create split not replace
func showBuffer(buffer_name string) *Buffer {
	for _, b := range buffers {
		if b.Name == buffer_name {
			if currentViewTree.Leaf.Buf == b {
				return b // already shown
			}
			currentViewTree = NewViewTreeLeaf(nil, NewView(b))
			rootViewTree = currentViewTree
			return b
		}
	}
	return nil
}

func closeBuffer(b *Buffer) {
	for i, b2 := range buffers {
		if b == b2 {
			buffers = append(buffers[:i], buffers[i+1:]...)
			break
		}
	}
}

func selectAvailableBuffer(closeIfNone bool) {
	if len(buffers) == 0 {
		if closeIfNone {
			// TODO Use method here (don't handcode screen.Fini())
			screen.Fini()
			os.Exit(0)
		}
	} else {
		currentViewTree = NewViewTreeLeaf(nil, NewView(buffers[0]))
		rootViewTree = currentViewTree
	}
}

func closeCurrentBuffer(force bool) {
	b := currentViewTree.Leaf.Buf
	if b.Modified && !force {
		messageError("Save buffer before closing it.")
		return
	}
	closeBuffer(b)
	selectAvailableBuffer(true)
}

func findBuffer(name string) *Buffer {
	for _, b := range buffers {
		if b.Name == name {
			return b
		}
	}
	return nil
}

var commands = map[string]func([]string){}
var commandAliases = map[string]string{}

func runCommand(args []string) {
	if len(args) == 0 {
		messageError("No command given!")
		return
	}
	command_name := args[0]
	if full_command_name, ok := commandAliases[command_name]; ok {
		command_name = full_command_name
	}
	if c, ok := commands[command_name]; ok {
		c(args)
	} else {
		messageError("No command named '" + command_name + "'")
	}
}

func addCommand(name string, fn func([]string)) {
	commands[name] = fn
}
func addAlias(alias, name string) {
	commandAliases[alias] = name
}

func initCommands() {
	addCommand("quit", func(args []string) {
		/////
		screen.Fini()
		os.Exit(0)
		/////
		//closeCurrentBuffer(false)
		closeCurrentBuffer(true)
	})
	addAlias("q", "quit")
	addCommand("quit!", func(args []string) {
		closeCurrentBuffer(false)
	})
	addAlias("q!", "quit!")
	addCommand("write", func(args []string) {
		b := currentViewTree.Leaf.Buf
		if len(args) > 1 {
			b.SetPath(args[1])
		}
		b.Save()
	})
	addAlias("w", "write")
	////////
	addCommand("wai", func(args []string) {
		b := currentViewTree.Leaf.Buf
		b.SetPath("./.quest.txt")
		b.Save()
		//fmt.Println("close")
		screen.Fini()
		os.Exit(0)
		//closeCurrentBuffer(false)
	})
	addAlias("ai", "wai")
	////////
	addCommand("edit", func(args []string) {
		if len(args) < 2 {
			messageError("Can't open buffer without a name or file path.")
		} else {
			path, err := filepath.Abs(args[1])
			if err != nil {
				path = args[1]
			}
			if b := openBufferFromFile(path); b != nil {
				showBuffer(b.Name)
			}
		}

	})
	addAlias("e", "edit")
	addAlias("o", "edit")
	addCommand("writequit", func(args []string) {
		runCommand([]string{"write"})
		runCommand([]string{"quit"})
	})
	addAlias("wq", "writequit")
	addCommand("buffers", func(args []string) {
		var b *Buffer
		if b = findBuffer("*buffers*"); b == nil {
			b = openBufferNamed("*buffers*")
			b.AddMode("buffers")
		}
		b.Data = [][]rune{}
		for _, buf := range buffers {
			if buf.Name != "*buffers*" {
				b.Data = append(b.Data, []rune(buf.Name))
			}
		}
		hook_trigger_buffer("modified", b)
		showBuffer(b.Name)
	})
	addAlias("b", "buffers")
}
