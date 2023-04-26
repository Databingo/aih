# Aih: use chatGPT, bard in terminal. 

![screenshot](aih.gif)

## Usage
```bash
./aih
```
## Command list
|command   | operation|
|----------|----------|
|.bard      | Bard|
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
