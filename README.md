# Talk to AI models in terminal.

<img src="aih.gif" alt="screenshot" style="width:80%;">

## Usage
Download [binary file](https://github.com/Databingo/aih/releases) then type:
```bash
./aih
```
## Command list
| Command    | Operation|
|------------|----------|
|.           | Select AI mode of Bard/ChatGPT/Claude(Claude2)/HuggingChat(Llama2)|
|.proxy      | Set proxy, for example: socks5://127.0.0.1:7890|
|<<          | Start multiple lines input mode|
|>>          | End multiple lines input mode|
|↑           | Previous input|
|↓           | Next input|
|.c or .clear| Clear the screen|
|.h or .history | Show history of conversations|
|j           | Scroll down|
|k           | Scroll up|
|g           | Scroll to top|
|G           | Scroll to bottom|
|q or Enter  | Back to conversation|
|.help       | Show help|
|.exit       | Exit Aih|

## Prerequisites
- [Chrome Browser](https://google.com/chrome)
- Free account of [Bard](https://bard.google.com), [Claude](https://claude.ai), [OpenAI](https://chat.openai.com), [HuggingChat](https://huggingface.co/chat) logged-in manually on your Chrome browser.
- (Optional) Paid ChatGPT API on [Billing](https://platform.openai.com/account/billing/overview). 

## Tips
- Close Chrome browser before run Aih at the first time.
- Usage of command [Liner](https://github.com/peterh/liner).
- Answer will be auotmatically saved in system clipboard for pasting.
- Conversations were persisted in `history.txt` beside Aih binary.
- `./aih -rod=show` to display pages in real browser.
- All-In-One mode will display answers from all the AI modes.
<img src="allin1.png" alt="screenshot" style="width:80%;">

## Supported OS
- Mac/Linux/Windows

## Installation
```
$ git clone https://github.com/Databingo/aih
$ go clean -cache && go clean -modcache 
$ cd aih && go mod tidy && go build 
```
## Acknowledgements
- github.com/go-rod/rod
- github.com/sashabaranov/go-openai 

## License
MIT and "You can choose not to include this MIT license into your copies".
