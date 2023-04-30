Here is the corrected version of the README file:

AIH: Use GoogleBard, BingChat, and ChatGPT in the terminal.

Screenshot usage:

./aih Command list
commandoperation
.bardBard
.bingBing Chat
.chatChatGPT Web (free)
.chatapiChatGPT API (pay)
.proxySet proxy, for example socks5://127.0.0.1:7890
<<Start multiple lines input model>>
End multiple lines input model
↑Previous input value
↓Next input value
.newNew conversation of ChatGPT
.speakVoice speak context (macOS)
.quietNo speak
.bardkeySet GoogleBard cookie
.bingkeySet BingChat cookie
.chatkeySet ChatGPT Web accessToken
.chatapikeySet ChatGPT API key
.clearClear screen
.helpHelp
.exitExit

Prerequisites:

For ChatGPT Web (free), you should have an account and a logged-in access token from OpenAI. For ChatGPT API (pay), you should have a paid API on billing. For GoogleBard, you should join the waitlist and have a cookie value of __Secure-lPSID. For BingChat, you should apply for the waitlist and have a cookie.

How to get Cookies:

For GoogleBard cookie, you can log in, then add the Cookie-Editor extension. Click on the right-top corner to copy the __Secure-lPSID value. For BingChat cookie, you can log in, use Cookie-Editor, click the Cookie-Editor icon, click "Export," then click "Export as JSON" (this saves your cookies to the clipboard). Type .bingkey in AIH, and you will see a prompt "Please type <<, then paste Bing cookie, then type >>, then press Enter." By doing so, you can set the BingChat cookie via multiple lines input model.

Tips:

The returned text will be automatically saved in your system clipboard, so you can paste it anywhere directly.

Supported Operating Systems:

Mac, Linux, Windows

Installation:

Bash
$ git clone https://github.com/Databingo/aih
$ go clean -cache && go clean -modcache
$ cd aih && go mod tidy && go build

Or, download executable Binary file according to your operating system.

About:

Suggestions: This is an open plan based on the idea of "Correlation's enhancement of AI and human beings." If you have any suggestions, please write them in Issues.

Acknowledgements:

github.com/rocketlaunchr/google-search
github.com/sashabaranov/go-openai
github.com/CNZeroY/googleBard
github.com/pavel-one/EdgeGPT-Go
github.com/pengzhile/pandora
