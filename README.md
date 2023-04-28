# Aih: use ChatGPT, GoogleBard, BingChat in terminal. 

![screenshot](aih.png)

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
|.bardkey   | Set GoogleBard cookie|
|.bingkey   | Set BingChat cookie|
|.chatkey   | Set ChatGPT key|
|<          | Start multiple lines input|
|>          | End multiple lines input|
|.new       | New conversation of ChatGPT|
|.speak     | Voice speak context(macos)|
|.quiet     | No speak |
|.clear     | Clear screen|
|.help      | Help|
|.exit      | Exit|

## Pre-requests
- For ChatGPT you should have an account with payed API on [Billing](https://platform.openai.com/account/billing/overview). 
- For GoogleBard & BingChat you should in the waitlist on https://bard.google.com "Join Waitlist" or https://bing.com/new "Chat now".
- For GoogleBard you should login and add [Cookie-Editor](https://cookie-editor.cgagnier.ca) extension then Click it on the right-top corner to copy __Secure-lPSID value.
- For BingChat you should login then click Cookie-Editor icon then click "Export" then click "Export as JSON"(This saves your cookies to clipboard), then type .bing in Aih, when you type in something then hit Enter, you will see a prompt "Please type > then paste Bing cookie then type > then press Enter", by doing these you could set BingChat cookie via multiple lines input module.

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
- github.com/rocketlaunchr/google-search
- github.com/sashabaranov/go-openai 
- github.com/CNZeroY/googleBard
- github.com/pavel-one/EdgeGPT-Go
