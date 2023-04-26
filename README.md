# Aih: use chatGPT, google-bard, bing-chat in terminal. 

![screenshot](aih2.png)

## Usage
```bash
./aih
```
## Command list
|command   | operation|
|----------|----------|
|.bard      | Bard|
|.bing      | Bing Chat|
|.chat      | ChatGPT|
|.proxy     | Set proxy for example: socks5://127.0.0.1:7890|
|.key       | Set ChatGPT key|
|.new       | New conversation of ChatGPT|
|.speak     | Voice speak context(macos)|
|.quiet     | No speak |
|.clear     | Clear screen|
|.help      | Show help|
|.exit      | Exit|

## Pre-requests
- For ChatGPT you should have an account with payed API on [Billing](https://platform.openai.com/account/billing/overview). 
- For GoogleBard you should in the waitlist on https://bard.google.com "Join Waitlist". 
- For GoogleBard you should login and add [Cookie-Editor](https://cookie-editor.cgagnier.ca) extension then Click it to copy __Secure-lPSID value.
- For BingChat you should in the waitlist on https://bing.com/new "Join Waitlist".
- For BingChat you should login and add [Cookie-Editor](https://cookie-editor.cgagnier.ca) extension then click the Icon on the right-top corner then click "Export" then click "Export as JSON"(This saves your cookies to clipboard), then type .bing in aih, paste in terminal and hit Enter, then hit Ctrl+D.

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
- github.com/pavel-one/EdgeGPT-Go
