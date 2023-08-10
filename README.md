# Aih: Talk with Bard/Bing/ChatGPT/Claude in the terminal.

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
|.new        | Start a new conversation of ChatGPT|
|.help       | Show help|
|.exit       | Exit Aih|


## Prerequisites
- [Chrome Browser](https://google.com/chrome), `python3`, `undected_chromedriver`.
- For Bard, you need a free logged-in cookie of [Bard](https://bard.google.com).
- For Claude, you need a free logged-in cookie of [Claude](https://claude.ai).
- For ChatGPT, you need a free logged-in cookie of [OpenAI](https://chat.openai.com).
- For ChatGPT API (paid) you need a paid API on [Billing](https://platform.openai.com/account/billing/overview). 
- For Bing Chat, you need a free logged-in cookie of [Bing](https://account.microsoft.com) 

## How to get Cookies
- For Bing Chat cookie you can log in and then use [Cookie-Editor](https://cookie-editor.cgagnier.ca) -> click Cookie-Editor icon -> click "Export" -> click "Export as JSON" (This saves your cookies to the clipboard), then type `.key` to choose `Set Bing Chat Cookie` in Aih, you will see a prompt that says **"Please type << then paste Bing cookie then type >> then press Enter"**, by doing so you can set Bing Chat cookie via multiple lines input mode.
- For Google Bard cookie, same as Bing.
- For Chatgpt cookie, same as Bing.
- For Cloude cookie, same as Bing.

## Tips
- The returned text will be auotmatically saved in your system clipboard, so you can paste it anywhere directly.
- You can see more usages of command line operation from [here](https://github.com/peterh/liner#Line-editing).
- All conversation history was persisted locally in `history.txt`, in the same directory as the aih binary .

## Co-relation's Enhancement Function
| Command    | Operation|
|------------|----------|
|.speak      | Voice speak context(macOS only)|
|.quiet      | Disable voice output |

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
- For Bard & Chatgpt & Claude
```
$ pip3 install undetected_chromedriver 
```
- Or, download the executable [binary file](https://github.com/Databingo/aih/releases) according to your operating system.

## About Suggestions
This is an open plan based on the idea of "Co-relation's enhancement of AI and human beings". If you have any suggestions, please write them in the Issues section.

## Acknowledgements
- github.com/sashabaranov/go-openai 
- github.com/pavel-one/EdgeGPT-Go
- github.com/ultrafunkamsterdam/undetected-chromedriver
