# Aih: Talk with Bard/ChatGPT/Claude2/Llama2(HuggingChat) in the terminal.

![screenshot1](aih.gif) 

## Usage
Download [binary file](https://github.com/Databingo/aih/releases) then type:
```bash
./aih
```
## Command list
| Command    | Operation|
|------------|----------|
|.           | Select AI mode of Bard/ChatGPT/Claude/HuggingChat|
|.proxy      | Set proxy, for example: socks5://127.0.0.1:7890|
|<<          | Start multiple lines input mode|
|>>          | End multiple lines input mode|
|↑           | Previous input|
|↓           | Next input|
|.c or .clear| Clear the screen|
|.h or .history | Show history of conversations|
|j           | Scroll down|
|k           | Scroll up|
|g           | Scroll top|
|G           | Scroll bottom|
|q or Enter  | Back to conversation|
|.help       | Show help|
|.exit       | Exit Aih|

## Prerequisites
- [Chrome Browser](https://google.com/chrome)
- Free account of [Bard](https://bard.google.com), [Claude](https://claude.ai), [OpenAI](https://chat.openai.com), [HuggingChat](https://huggingface.co/chat) logged-in manually on your Chrome browser.
- Paid ChatGPT API on [Billing](https://platform.openai.com/account/billing/overview)(optional). 

## Tips
- The returned text will be auotmatically saved in your system clipboard for pasting.
- You can see more operations of command line [here](https://github.com/peterh/liner#Line-editing).
- All conversation history was persisted in `history.txt` locally along with Aih binary.
- All-In-One mode will display answers of all the AI modes.
<img src="allinone.png" alt="screenshot2" style="width:80%;">

## Supported Operating Systems:
- Mac
- Linux
- Windows

## Installation
- Bash
```
$ git clone https://github.com/Databingo/aih
$ go clean -cache && go clean -modcache 
$ cd aih && go mod tidy && go build 
```
- Or, download the executable [binary file](https://github.com/Databingo/aih/releases) according to your operating system.

## About Suggestions
This is an open plan based on the idea of "Co-relation's enhancement of AI and human beings". If you have any suggestions, please write them in the Issues section.

## Acknowledgements
- github.com/go-rod/rod
- github.com/sashabaranov/go-openai 
