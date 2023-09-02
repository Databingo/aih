# Todo

- Autocomplete (based on open buffers tokens)
- Buffer related commands/keybindings (language modes?)
- ~Command mode (q,w,o,e,wq,!,wqall)~
- ~Marks~
- ~Visual mode (+ visual line)~
- ~Incremental Search~
- ~Buffer list/switch mode~
- ~File tree mode~

**Post 1.0**

- Syntax highlithing
- Settings
- Scripting / user settings
- Fuzzy file search mode
- Recursive grep mode
- Line wrapping
- Windows
- Color schemes
- Shell mode
- Ensure UTF-8 works
- Help/Tutorial
- Fringe for error reporting / plugins

# The big list

- [ ] Auto indent
- [ ] Custom bindings
- [ ] Line numbers
- [ ] Search and replace
  - [ ] Search
  - [ ] Replace
- [ ] Tests
- [ ] Error handling
  - [ ] Fatal
  - [ ] Script/Runtime
  - [ ] User
- [ ] Unicode support
- [ ] Command execution
- [ ] Movement keys
- [ ] Visual mode
  - [ ] Normal
  - [ ] Line
- [ ] Help
  - [ ] General
  - [ ] Per command
  - [ ] Tutorial
- [ ] Options/Configuration
  - [ ] Save/Load
  - [ ] Tabs to space
  - [ ] Tab size
  - [ ] Color scheme
- [ ] Undo/Redo
- [ ] Clipboard
  - [ ] Copy
  - [ ] Paste
  - [ ] Cut
  - [ ] Registers
- [ ] Macros
- [ ] Syntax highlighting
- [ ] Color schemes

# Vim's Perl Interface for inspiration:

```
VIM::Msg({msg}, {group}?)
VIM::SetOption({arg})             Sets a vim option.
VIM::Buffers([{bn}...])           With no arguments, returns a list of all the buffers.
VIM::Windows([{wn}...])           With no arguments, returns a list of all the windows.
VIM::DoCommand({cmd})             Executes Ex command {cmd}.
VIM::Eval({expr})                 Evaluates {expr} and returns (success, val).
Window->SetHeight({height})
Window->Cursor({row}?, {col}?)
Window->Buffer()
Buffer->Name()                    Returns the filename for the Buffer.
Buffer->Number()                  Returns the number of the Buffer.
Buffer->Count()                   Returns the number of lines in the Buffer.
Buffer->Get({lnum}, {lnum}?, ...)
Buffer->Delete({lnum}, {lnum}?)
Buffer->Append({lnum}, {line}, {line}?, ...)
Buffer->Set({lnum}, {line}, {line}?, ...)
$main::curwin
$main::curbuf
```

# Some Emacs interface functions

...to consider implementing?

- `(redraw-display)`
- `(redisplay force?)`
- `(message format &rest args)`
- `(with-temp-message message &rest body)` Show message, executes body, removes message, returns body result
- `(current-message)`
- `(make-progress-reporter message &optional min-value max-value current-value min-change min-time)`
- `(progress-reporter-update reporter &optional value)`
- `(progress-reporter-done reporter)`
- `(messages-buffer)`
- `(yes-or-no-p)`
- `cursor-in-echo-area`
- `(format)`
- `(display-warning type message &optional level)` Levels being: emergency, error, warning, debug
- Handle hiding/make buffer sections invisible (for folding)
- `before-init-time`
- `after-init-time`
- `(save-excursion)`

