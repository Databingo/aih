# Aih: use chatGPT, bard in terminal. 

![screenshot](aih.gif)

## Usage
```bash
./aih
```
## Command list
|command   | operation|
|----------|----------|
|.help      | Show help|
|.proxy     | Set proxy for example: socks5://127.0.0.1:7890|
|.bard      | Bard|
|.chat      | ChatGPT|
|.key       | Set ChatGPT key|
|.new       | New conversation of ChatGPT|
|.speak     | Voice speak context(macos)|
|.quiet     | No speak |
|.clear     | Clear screen|
|.exit      | Exit|


## Support OS
- Mac
- Linux
- Windows

## Installation
```bash
$ git clone https://github.com/Databingo/aih
$ cd aih && go mod tidy && go build 
```
## About Suggestions
This is an open plan based on the idea of "Co-relation's enhancement of AI and human beings".
If you have any suggestions please write in Issues.

## Acknowledge
- github.com/sashabaranov/go-openai 
- github.com/CNZeroY/googleBard
