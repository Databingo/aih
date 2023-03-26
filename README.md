# aih

## Introduce
This is a terminal AI assistant tool based on the idea of "Co-relation's enhancement of AI and human beings". 
Current functions:
1. Assist you chat with ChatGPT-3.5-turbo from terminal.
2. Automatic voice reading AI returns text. 
3. Assist you create precise prompt sequentially by type .prompt
![screenshot](aih.gif)

## Support OS
- Mac
- Linux

## Installation
```bash
$ git clone https://github.com/Databingo/aih
$ cd aih && go mod tidy && go build 
```
## Command list
|command   | operation|
|----------|----------|
|.help      | Show help|
|.key       | Set key|
|.proxy     | Set proxy|
|.new       | New conversation|
|.speak     | Voice speak context|
|.quiet     | Quiet not speak |
|.clear     | Clear screen|
|.exit      | Exit|
|-----------|----------|
|.prompt    | Role of Assistant for create precise prompt|

## Usage
```bash
./aih
```
1. Paste you OpenAI key from the terminal the first time you run aih;
2. Type .proxy to set proxy if you need;
3. It automatic start with conversation mode.
4. The number in prompt is the left tokens in this conversation.

## Todo
1. Tidy code.

## About Suggestions
This is just a test in concept of "Co-relation's enhancement of AI and human beings". 
If you have any suggestions please write in Issues.



