# Aih: Talk with Bard/ChatGPT/Claude/HuggingChat in the terminal.

![screenshot2](aih.gif) 

## Usage
Download [binary file](https://github.com/Databingo/aih/releases) then type:
```bash
./aih
```
## Command list
| Command    | Operation|
|------------|----------|
|.           | Select AI mode of Bard/Bing/ChatGPT/Claude|
|.key        | Set cookie of Bard/Bing/ChatGPT/Claude|
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
|.speak      | Voice speak context(macOS only)|
|.quiet      | Disable voice output |
|.new        | Start a new conversation of ChatGPT|
|.help       | Show help|
|.exit       | Exit Aih|

## Prerequisites
- [Chrome Browser](https://google.com/chrome), `python3`, `pip3 install undetected_chromedriver`
- For Bard, you need a free logged-in cookie of [Bard](https://bard.google.com).
- For Claude, you need a free logged-in cookie of [Claude](https://claude.ai).
- For ChatGPT, you need a free logged-in cookie of [OpenAI](https://chat.openai.com).
- For HuggingChat, you need a free logged-in cookie of [HuggingChat](https://huggingface.co/chat).
- For ChatGPT API (paid) you need a paid API on [Billing](https://platform.openai.com/account/billing/overview). 

## How to get Cookies
- For Bard cookie, you can log in and then use [Cookie-Editor](https://cookie-editor.cgagnier.ca) -> click Cookie-Editor icon -> click "Export" -> click "Export as JSON" (This saves your cookies to the clipboard), then type `.key` to choose `Set Bard Cookie` in Aih, you will see a prompt that says **"Please type << then paste Bard cookie then type >> then press Enter"**, by doing so you can set Bard cookie via multiple lines input mode.
- For Chatgpt cookie, same.
- For Cloude cookie, same.
- For HuggingChat cookie, same.

## Tips
- The returned text will be auotmatically saved in your system clipboard, so you can paste it anywhere directly.
- You can see more usages of command line operation from [here](https://github.com/peterh/liner#Line-editing).
- All conversation history was persisted locally in `history.txt`, in the same directory as the Aih binary .

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
- github.com/sashabaranov/go-openai 
- github.com/ultrafunkamsterdam/undetected-chromedriver
