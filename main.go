package main

import (
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"github.com/go-rod/stealth"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/peterh/liner"
	"github.com/rivo/tview"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// var trace = true
var trace = false

var color_bard = tcell.ColorDarkCyan
var color_bing = tcell.ColorDarkMagenta
var color_chat = tcell.ColorWhite
var color_chatapi = tcell.ColorWhite
var color_claude = tcell.ColorYellow
var color_huggingchat = tcell.ColorDarkMagenta

func clear() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func multiln_input(Liner *liner.State, prompt string) string {
	// For recognize multipile lines input module
	// |--------------------------|------
	// |recording && input        | action
	// |--------------------------|------
	// |false && == "" or x       | record; break;
	// |false && != "<<"          | record; break;
	// |false && == "<<" + ">>"   | record; break; rm << >>;
	// |false && == "<<"          | record; true; rm <<;
	// |true  && == "" or x       | record;
	// |true  && != ">>"          | record;
	// |true  && == ">>"          | record; break; rm >>;
	// |--------------------------|------

	var ln string
	var lns []string
	recording := false
	for {
		if ln == "" && !recording {
			ln, _ = Liner.Prompt(prompt)
		} else {
			ln, _ = Liner.Prompt("")
		}
		ln = strings.Trim(ln, " ")
		if !recording && (ln == "" || len(ln) == 1) {
			lns = append(lns, ln)
			break
		} else if !recording && ln[:2] != "<<" {
			lns = append(lns, ln)
			break
		} else if !recording && ln[:2] == "<<" && len(ln) >= 4 && ln[len(ln)-2:] == ">>" {
			lns = append(lns, ln[2:len(ln)-2])
			break
		} else if !recording && ln[:2] == "<<" {
			recording = true
			lns = append(lns, ln[2:])
		} else if recording && (ln == "" || len(ln) == 1) {
			lns = append(lns, ln)
		} else if recording == true && ln[len(ln)-2:] != ">>" {
			lns = append(lns, ln)
		} else if recording == true && ln[len(ln)-2:] == ">>" {
			recording = false
			lns = append(lns, ln[:len(ln)-2])
			break
		}
	}

	long_str := strings.Join(lns, "\n")
	return long_str
}

// Write response RESP to clipboard
func save2clip_board(rs string) {
	err := clipboard.WriteAll(rs)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Create prompt for user input
	Liner := liner.NewLiner()
	defer Liner.Close()

	// Use RESP for record response per time
	var RESP string

	// Read Aih Configure
	aih_json, err := ioutil.ReadFile("aih.json")
	if err != nil {
		err = ioutil.WriteFile("aih.json", []byte(""), 0644)
	}

	// Read Proxy
	Proxy := gjson.Get(string(aih_json), "proxy").String()

	// Set proxy for system of current program
	//os.Setenv("http_proxy", Proxy)
	//os.Setenv("https_proxy", Proxy)

	role := ".bard"

	// Set proxy
	proxy_u := launcher.NewUserMode().
		//Proxy(Proxy).
		//Leakless(true).// indepent tab | work with UserDataDir()
		//UserDataDir("data").// indepent tab + data
		//Set("disable-default-apps").
		//Headless(true).
		MustLaunch()

	// Open rod browser
	var browser *rod.Browser
	browser = rod.New().
		Trace(trace).
		ControlURL(proxy_u).
		Timeout(60 * 24 * time.Minute).
		MustConnect()

	// Get cookies (for login AI accounts)
	cookies := browser.MustGetCookies()

	// Set proxy for daemon browser_
	proxy_url := launcher.New().
		Proxy(Proxy).
		MustLaunch()

	// Open rod daemon browser
	var browser_ *rod.Browser
	if Proxy != "" {
		browser_ = rod.New().
			Trace(trace).
			ControlURL(proxy_url).
			Timeout(60 * 24 * time.Minute).
			MustConnect()
	} else {
		browser_ = rod.New().
			Trace(trace).
			Timeout(60 * 24 * time.Minute).
			MustConnect()

	}

	// Share cookies
	for _, i := range cookies {
		browser_.MustSetCookies(i)
	}

	// Renew browser to daemon
	browser = browser_
	browser.ServeMonitor(":7777")

	//////////////////////0////////////////////////////
	// Set up client of OpenAI API
	key := gjson.Get(string(aih_json), "key")
	OpenAI_Key := key.String()
	config := openai.DefaultConfig(OpenAI_Key)
	client := openai.NewClientWithConfig(config)
	messages := make([]openai.ChatCompletionMessage, 0)

	//////////////////////c1////////////////////////////
	// Set up client of Bard (Rod version)
	var page_bard *rod.Page
	var relogin_bard = true
	channel_bard := make(chan string)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				relogin_bard = true
			}
		}()
		//page_bard := browser.MustPage("https://bard.google.com")
		page_bard = stealth.MustPage(browser)
		page_bard.MustNavigate("https://bard.google.com")

		for i := 1; i <= 30; i++ {
			if page_bard.MustHasX("//textarea[@id='mat-input-0']") {
				relogin_bard = false
				break
			}
			time.Sleep(time.Second)
		}
		if relogin_bard == true {
			fmt.Println("✘ Bard")
		}
		if relogin_bard == false {
			fmt.Println("✔ Bard")
			for {
				select {
				case question := <-channel_bard:
					//page_bard.MustActivate()
					page_bard.MustElementX("//textarea[@id='mat-input-0']").MustWaitVisible().MustInput(question)
					page_bard.MustElementX("//button[@mattooltip='Submit']").MustClick()
					fmt.Println("Bard generating...")
					//if role == ".all" {
					//	channel_bard <- "click_bard"
					//}
					// wait generated icon
					var generated_icon_appear = false
					for i := 1; i <= 60; i++ {
						if page_bard.MustHasX("//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')]") {
							generated_icon_appear = true
							break
						}
						time.Sleep(1 * time.Second)
					}
					if generated_icon_appear == true {
						img := page_bard.MustElementX("//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')][last()]").MustWaitVisible()
						response := img.MustElementX("parent::div/parent::div").MustWaitVisible()
						answer := response.MustText()
						channel_bard <- answer
					} else {
						channel_bard <- "✘✘  Bard, Please check the internet connection and verify login status."
						relogin_bard = true

					}
				}
			}
		}

	}()

	//////////////////////c2////////////////////////////
	// Set up client of Claude (Rod version)
	var page_claude *rod.Page
	var relogin_claude = true
	channel_claude := make(chan string)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				relogin_claude = true
			}
		}()
		//page_claude = browser.MustPage("https://claude.ai")
		page_claude = stealth.MustPage(browser)
		page_claude.MustNavigate("https://claude.ai/api/organizations").MustWaitLoad()
		org_json := page_claude.MustElementX("//pre").MustText()
		org_uuid := gjson.Get(string(org_json), "0.uuid").String()
		time.Sleep(6 * time.Second)

		new_uuid := uuid.New().String()
		new_uuid_url := "https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations"
		create_new_converastion_json := `{"uuid":"` + new_uuid + `","name":""}`
		create_new_converastion_js := `
		 (new_uuid_url, sdata) => {
		 var xhr = new XMLHttpRequest();
		 xhr.open("POST", new_uuid_url);
		 xhr.setRequestHeader('Content-Type', 'application/json');
		 xhr.setRequestHeader('Referer', 'https://claude.ai/chats');
		 xhr.setRequestHeader('Origin', 'https://claude.ai');
		 xhr.setRequestHeader('TE', 'trailers');
		 xhr.setRequestHeader('DNT', '1');
		 xhr.setRequestHeader('Connection', 'keep-alive');
		 xhr.setRequestHeader('Accept', 'text/event-stream, text/event-stream');
		 xhr.onreadystatechange = function() {
		     if (xhr.readyState == XMLHttpRequest.DONE) { 
		         var res_text = xhr.responseText;
		         console.log(res_text);
		        } 
		     }
	         console.log(sdata);
		 xhr.send(sdata);
		}
		`
		page_claude.MustEval(create_new_converastion_js, new_uuid_url, create_new_converastion_json).Str()
		// posted new conversation uuid
		time.Sleep(3 * time.Second) // delay to simulate human being

		var record_chat_messages string
		var response_chat_messages string
		for i := 1; i <= 20; i++ {
			create_json := page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()

			message_uuid := gjson.Get(string(create_json), "uuid").String()
			if message_uuid == new_uuid {
				//fmt.Println("create_conversation success...")
				relogin_claude = false
				record_chat_messages = gjson.Get(string(create_json), "chat_messages").String()
				break
			}
			time.Sleep(2 * time.Second)
		}

		if relogin_claude == true {
			fmt.Println("✘ Claude")
		}
		if relogin_claude == false {
			fmt.Println("✔ Claude")
			for {
				select {
				case question := <-channel_claude:
					// re-activate
					page_claude.MustNavigate("https://claude.ai/api/account/statsig/" + org_uuid).MustWaitLoad()
					record_json := page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()
					record_chat_messages = gjson.Get(string(record_json), "chat_messages").String()
					record_count := gjson.Get(string(response_chat_messages), "#").Int()
					page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid).MustWaitLoad()
					time.Sleep(1 * time.Second)                         // delay to simulate human being
					question = strings.Replace(question, `"`, `\"`, -1) // escape " in input text when code into json

					d := `{"completion":{"prompt":"` + question + `","timezone":"Asia/Shanghai","model":"claude-2"},"organization_uuid":"` + org_uuid + `","conversation_uuid":"` + new_uuid + `","text":"` + question + `","attachments":[]}`
					js := `
		                               (sdata) => {
		                               var xhr = new XMLHttpRequest();
		                               xhr.open("POST", "https://claude.ai/api/append_message");
		                               xhr.setRequestHeader('Content-Type', 'application/json');
		                               xhr.setRequestHeader('Referer', 'https://claude.ai/chats');
		                               xhr.setRequestHeader('Origin', 'https://claude.ai');
		                               xhr.setRequestHeader('TE', 'trailers');
		                               xhr.setRequestHeader('Connection', 'keep-alive');
		                               xhr.setRequestHeader('Accept', 'text/event-stream, text/event-stream');
		                               xhr.onreadystatechange = function() {
		                                   if (xhr.readyState == XMLHttpRequest.DONE) { 
		                                       var res_text = xhr.responseText;
		                                       console.log(res_text);
		                                   }
		                               }
		                               console.log(sdata);
		                               xhr.send(sdata);
		                               }
		                              `
					page_claude.MustEval(js, d).Str()
					fmt.Println("Claude generating...")
					//if role == ".all" {
					//	channel_claude <- "click_claude"
					//}
					time.Sleep(3 * time.Second) // delay to simulate human being

					// wait answer
					var claude_response = false
					var response_json string
					for i := 1; i <= 20; i++ {
						response_json = page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()
						response_chat_messages = gjson.Get(string(response_json), "chat_messages").String()
						count := gjson.Get(string(response_chat_messages), "#").Int()

						if response_chat_messages != record_chat_messages && count == record_count+2 {
							claude_response = true
							record_chat_messages = response_chat_messages
							break
						}
						time.Sleep(3 * time.Second)
					}
					if claude_response == true {
						count := gjson.Get(string(response_chat_messages), "#").Int()
						answer := gjson.Get(string(response_json), "chat_messages.#(index=="+strconv.FormatInt(count-1, 10)+").text").String()
						channel_claude <- answer
					} else {
						channel_claude <- "✘✘  Claude, Please check the internet connection and verify login status."
						relogin_claude = true
					}
				}
			}
		}

	}()

	//////////////////////c3////////////////////////////
	// Set up client of Huggingchat (Rod version)
	var page_hc *rod.Page
	var relogin_hc = true
	channel_hc := make(chan string)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				relogin_hc = true
			}
		}()
		//page_hc = browser.MustPage("https://huggingface.co/chat")
		page_hc = stealth.MustPage(browser)
		page_hc.MustNavigate("https://huggingface.co/chat")
		for i := 1; i <= 30; i++ {
			if page_hc.MustHasX("//button[contains(text(), 'Sign Out')]") {
				relogin_hc = false
				break
			}
			time.Sleep(time.Second)
		}
		if relogin_hc == true {
			fmt.Println("✘ HuggingChat")
		}
		if relogin_hc == false {
			fmt.Println("✔ HuggingChat")
			for {
				select {
				case question := <-channel_hc:
					//page_hc.MustActivate()
					page_hc.Timeout(20 * time.Second).MustElementX("//textarea[@enterkeyhint='send']").MustInput(question)
					page_hc.Timeout(20 * time.Second).MustElement("button svg path[d='M27.71 4.29a1 1 0 0 0-1.05-.23l-22 8a1 1 0 0 0 0 1.87l8.59 3.43L19.59 11L21 12.41l-6.37 6.37l3.44 8.59A1 1 0 0 0 19 28a1 1 0 0 0 .92-.66l8-22a1 1 0 0 0-.21-1.05Z']").MustClick()
					fmt.Println("HuggingChat generating...")
					//if role == ".all" {
					//	channel_hc <- "click_hc"
					//}
					for {
						info := page_hc.MustInfo()
						if strings.HasPrefix(info.URL, "https://huggingface.co/chat/conversation") {
							break
						}
						time.Sleep(1 * time.Second)
					}

					// stop_icon
					var stop_icon_disappear = false
					for i := 1; i <= 60; i++ {
						if page_hc.MustHas("svg path[d='M24 6H8a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2Z']") {
							stop_icon_disappear = false
						} else {
							stop_icon_disappear = true
							break

						}
						time.Sleep(time.Second)
					}
					if stop_icon_disappear == true {
						page_hc.MustHasX("//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')]")
						page_hc.MustElementX("//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')]").MustWaitVisible()
						img := page_hc.MustElementX("(//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')])[last()]")
						content := img.MustElementX("following-sibling::div[1]")
						answer := content.MustText()
						channel_hc <- answer
					} else {
						channel_hc <- "✘✘  HuggingChat, Please check the internet connection and verify login status."
						relogin_hc = true

					}
				}
			}
		}

	}()

	//////////////////////c4////////////////////////////
	// Set up client of chatgpt (rod version)
	var page_chatgpt *rod.Page
	var relogin_chatgpt = true
	channel_chatgpt := make(chan string)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				relogin_chatgpt = true
			}
		}()
		//page_chatgpt = browser.MustPage("https://chat.openai.com")
		page_chatgpt = stealth.MustPage(browser)
		page_chatgpt.MustNavigate("https://chat.openai.com")

		channel_chatgpt_tips := make(chan string)
		go func() {
			for i := 1; i <= 30; i++ {
				if page_chatgpt.MustHasX("//div[contains(text(), 'Okay, let')]") {
					page_chatgpt.MustElementX("//div[contains(text(), 'Okay, let')]").MustWaitVisible().MustClick()
					close(channel_chatgpt_tips)
					break
				}
				time.Sleep(time.Second)
			}
		}()

		for i := 1; i <= 30; i++ {
			if page_chatgpt.MustHasX("//textarea[@id='prompt-textarea']") && !page_chatgpt.MustHasX("//h2[contains(text(), 'Your session has expired')]") {
				relogin_chatgpt = false
				break
			}
			time.Sleep(time.Second)
		}

		if relogin_chatgpt == true {
			fmt.Println("✘ ChatGPT")
		}
		if relogin_chatgpt == false {
			fmt.Println("✔ ChatGPT")
			for {
				select {
				case question := <-channel_chatgpt:
					//page_chatgpt.MustActivate()
					if page_chatgpt.MustHasX("//div[contains(text(), 'Something went wrong. If this issue persists please')]") {

						fmt.Println("ChatGPT web error")
						channel_chatgpt <- "✘✘  ChatGPT, Please check the internet connection and verify login status."
						relogin_chatgpt = true

					}
					page_chatgpt.MustElementX("//textarea[@id='prompt-textarea']").MustWaitVisible().MustInput(question)
					page_chatgpt.MustElementX("//textarea[@id='prompt-textarea']/..//button").MustClick()
					fmt.Println("ChatGPT generating...")
					//if role == ".all" {
					//	channel_chatgpt <- "click_chatgpt"
					//}

					var regenerate_icon = false
					for i := 1; i <= 60; i++ {
						if page_chatgpt.MustHas("svg path[d='M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15']") {
							regenerate_icon = true
							break
						}
						time.Sleep(1 * time.Second)
					}
					if regenerate_icon == true {
						answer := page_chatgpt.MustElementX("(//div[contains(@class, 'group w-full')])[last()]").MustText()[7:]
						if strings.Contains(answer,
							"An error occurred. Either the engine you requested does not exist or there was another issue processing your request. If this issue persists please contact us through our help center at help.openai.com.") {
							relogin_chatgpt = true
						}
						channel_chatgpt <- answer
					} else {
						channel_chatgpt <- "✘✘  ChatGPT, Please check the internet connection and verify login status."
						relogin_chatgpt = true

					}
				}
			}
		}
	}()

	// Exit when wake up for the disconnecting with daemon browser
	go func() {
		//fmt.Println("wake monitor...")
		for {
			utils.Sleep(3)
			if _, err := browser.Version(); err != nil {
				browser = rod.New().MustConnect()
				browser.MustClose()
				fmt.Println("Please restart Aih because the daemon process has been disconnected.")
				close(channel_bard)
				close(channel_chatgpt)
				close(channel_claude)
				close(channel_hc)
				Liner.Close()
				syscall.Exit(0)
			}

		}

	}()

	clear()

	// Welcome to Aih
	welcome := `
╭ ────────────────────────────── ╮
│    Welcome to Aih              │ 
│    Type .help for help         │ 
╰ ────────────────────────────── ╯ `
	fmt.Println(welcome)

	max_tokens := 4097
	used_tokens := 0
	left_tokens := 0
	speak := 0
	uInput := ""
	chat_mode := openai.GPT3Dot5Turbo
	chat_completion := true

	// Start loop to read user input
	for {
		// Re-read user input history
		if f, err := os.Open(".history"); err == nil {
			Liner.ReadHistory(f)
			f.Close()
		}

		prompt := strconv.Itoa(left_tokens) + role + "> "
		userInput := multiln_input(Liner, prompt)
		//fmt.Println("userInput:", userInput)

		// Check Aih commands
		switch userInput {
		case "":
			continue
		case ".proxy":
			proxy, _ := Liner.Prompt("Please input your proxy:")
			if proxy == "" {
				continue
			}
			aihj, err := ioutil.ReadFile("aih.json")
			new_aihj, _ := sjson.Set(string(aihj), "proxy", proxy)
			err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
			fmt.Println("Please restart Aih for using proxy")
			/// exit
			browser.MustClose()
			close(channel_bard)
			close(channel_chatgpt)
			close(channel_claude)
			close(channel_hc)
			Liner.Close()
			syscall.Exit(0)
			//os.Exit(0)
		case ".help":
			fmt.Println("                           ")
			fmt.Println("                 Welcome to Aih!                             ")
			fmt.Println("------------------------------------------------------------ ")
			fmt.Println(" .               Select AI mode of Bard/ChatGPT/Claude2/HuggingChat(Llama2)")
			fmt.Println(" ↑               Previous input")
			fmt.Println(" ↓               Next input")
			fmt.Println(" <<              Start multiple lines input")
			fmt.Println(" >>              End multiple lines input")
			fmt.Println(" j               Scroll down")
			fmt.Println(" k               Scroll up")
			fmt.Println(" f               Page down")
			fmt.Println(" p               Page up")
			fmt.Println(" g               Scroll top")
			fmt.Println(" G               Scroll bottom")
			fmt.Println(" q or Enter      Back to conversation")
			fmt.Println(" .c or .clear    Clear screen")
			fmt.Println(" .h or .history  Show history")
			fmt.Println(" .key            Set key of ChatGPT API")
			fmt.Println(" .proxy          Set proxy")
			fmt.Println(" .help           Help")
			fmt.Println(" .exit           Exit")
			fmt.Println(" .speak          Voice speak context (MasOS only)")
			fmt.Println(" .quiet          Not speak")
			//fmt.Println(" .new            New conversation of ChatGPT")
			fmt.Println("------------------------------------------------------------ ")
			fmt.Println("                           ")
			fmt.Println("                           ")
			continue
		case ".c", ".clear":
			clear()
			continue
		case ".h", ".history":
			cnt, _ := ioutil.ReadFile("history.txt")
			printer(color_chat, string(cnt), true)
			continue
		case ".exit":
			//exit
			browser.MustClose()
			close(channel_bard)
			close(channel_chatgpt)
			close(channel_claude)
			close(channel_hc)
			Liner.Close()
			syscall.Exit(0)
		case ".new":
			// For role .chat
			//conversation_id = ""
			//parent_id = ""
			// For role .chatapi
			messages = make([]openai.ChatCompletionMessage, 0)
			//max_tokens = 4097
			used_tokens = 0
			left_tokens = max_tokens - used_tokens
			continue
		case ".", "/":
			proms := promptui.Select{
				Label: "Select AI mode to chat",
				Size:  10,
				Items: []string{
					"All-In-One",
					"Bard",
					"ChatGPT",
					"Claude",
					"HuggingChat",
					"ChatGPT API gpt-3.5-turbo, $0.002/1K tokens",
					"ChatGPT API gpt-4 8K Prompt, $0.03/1K tokens",
					"ChatGPT API gpt-4 8K Completion, $0.06/1K tokens",
					"ChatGPT API gpt-4 32K Prompt, $0.06/1K tokens",
					"ChatGPT API gpt-4 32K Completion, $0.12/1K tokens",
					"Exit",
				},
			}

			_, ai, err := proms.Run()
			if err != nil {
				panic(err)
			}

			switch ai {
			case "Bard":
				role = ".bard"
				left_tokens = 0
				continue
			case "Bing":
				role = ".bing"
				left_tokens = 0
				continue
			case "ChatGPT":
				role = ".chat"
				left_tokens = 0
				continue
			case "Claude":
				role = ".claude"
				left_tokens = 0
				continue
			case "HuggingChat":
				role = ".huggingchat"
				left_tokens = 0
				continue
			case "All-In-One":
				role = ".all"
				left_tokens = 0
				continue
			case "ChatGPT API gpt-3.5-turbo, $0.002/1K tokens":
				role = ".chatapi"
				chat_mode = openai.GPT3Dot5Turbo
				max_tokens = 4097
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				chat_completion = true
				continue
			case "ChatGPT API gpt-4 8K Prompt, $0.03/1K tokens":
				role = ".chatapi"
				chat_mode = openai.GPT4
				max_tokens = 8192
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				chat_completion = false
				continue
			case "ChatGPT API gpt-4 8K Completion, $0.06/1K tokens":
				role = ".chatapi"
				chat_mode = openai.GPT4
				max_tokens = 8192
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				chat_completion = true
				continue
			case "ChatGPT API gpt-4 32K Prompt, $0.06/1K tokens":
				role = ".chatapi"
				chat_mode = openai.GPT432K
				max_tokens = 32768
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				chat_completion = false
				continue
			case "ChatGPT API gpt-4 32K Completion, $0.12/1K tokens":
				role = ".chatapi"
				chat_mode = openai.GPT432K
				max_tokens = 32768
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				chat_completion = true
				continue
			case "Exit":
				continue
			}
		case ".key":
			prom := promptui.Select{
				Label: "Select:",
				Size:  6,
				Items: []string{
					"Set ChatGPT API Key",
					"Exit",
				},
			}

			_, keyy, err := prom.Run()
			if err != nil {
				panic(err)
			}

			switch keyy {
			case "Set ChatGPT API Key":
				OpenAI_Key = ""
				role = ".chatapi"
				goto CHATAPI
			case "Exit":
				continue
			}

		case ".speak":
			speak = 1
			continue
		case ".quiet":
			speak = 0
			continue
		default:
			// Re-read user input history in case other process alternated
			if f, err := os.Open(".history"); err == nil {
				Liner.ReadHistory(f)
				f.Close()
			}
			// Record user input without Aih commands
			uInput = strings.Replace(userInput, "\r", "\n", -1)
			uInput = strings.Replace(uInput, "\n", " ", -1)
			Liner.AppendHistory(uInput)
			// Persistent user input
			if f, err := os.Create(".history"); err == nil {
				Liner.WriteHistory(f)
				f.Close()
			}

		}

		// ALL-IN-ONE:
		if role == ".all" {
			if relogin_bard == false {
				channel_bard <- userInput
				//<-channel_bard
			}
			if relogin_chatgpt == false {
				channel_chatgpt <- userInput
				//<-channel_chatgpt
			}
			if relogin_claude == false {
				channel_claude <- userInput
				//<-channel_claude
			}
			if relogin_hc == false {
				channel_hc <- userInput
				//<-channel_hc
			}

			if relogin_bard == false {
				answer_bard := <-channel_bard
				fmt.Println(">Bard Done.")
				RESP += "\n\n---------------- bard answer ----------------\n"
				RESP += strings.TrimSpace(answer_bard)
			}
			if relogin_chatgpt == false {
				answer_chatgpt := <-channel_chatgpt
				fmt.Println(">ChatGPT Done.")
				RESP += "\n\n---------------- chatgpt answer ----------------\n"
				RESP += strings.TrimSpace(answer_chatgpt)
			}
			if relogin_claude == false {
				answer_claude := <-channel_claude
				fmt.Println(">Claude Done.")
				RESP += "\n\n---------------- claude answer ----------------\n"
				RESP += strings.TrimSpace(answer_claude)
			}
			if relogin_hc == false {
				answer_hc := <-channel_hc
				fmt.Println(">HuggingChat Done.")
				RESP += "\n\n---------------- huggingchat answer ----------------\n"
				RESP += strings.TrimSpace(answer_hc)
			}
			printer(color_chat, RESP, false)

		}
		if role == ".bard" {
			if relogin_bard == true {
				fmt.Println("✘ Bard")
			} else {
				//fmt.Println("main thread get", userInput)
				channel_bard <- userInput
				//fmt.Println("put", userInput, "into channel_bard")
				answer := <-channel_bard

				// Print the response to the terminal
				RESP = strings.TrimSpace(answer)
				printer(color_bard, RESP, false)
			}

		}

		// CLAUDE:
		if role == ".claude" {
			if relogin_claude == true {
				fmt.Println("✘ Claude")
			} else {
				channel_claude <- userInput
				answer := <-channel_claude

				RESP = strings.TrimSpace(answer)
				printer(color_claude, RESP, false)
			}

		}
		// CHATGPT:
		if role == ".chat" {
			if relogin_chatgpt == true {
				fmt.Println("✘ ChatGPT")
			} else {
				channel_chatgpt <- userInput
				answer := <-channel_chatgpt

				// Print the response to the terminal
				RESP = strings.TrimSpace(answer)
				printer(color_chatapi, RESP, false)
			}

		}

		// HUGGINGCHAT:
		if role == ".huggingchat" {
			if relogin_hc == true {
				fmt.Println("✘ HuggingChat")
			} else {
				channel_hc <- userInput
				answer := <-channel_hc

				// Print the response to the terminal
				RESP = strings.TrimSpace(answer)
				printer(color_huggingchat, RESP, false)
			}

		}
	CHATAPI:
		if role == ".chatapi" {
			// Check ChatGPT API Key
			if OpenAI_Key == "" {
				OpenAI_Key, _ = Liner.Prompt("Please input your OpenAI Key: ")
				if OpenAI_Key == "" {
					continue
				}
				aihj, err := ioutil.ReadFile("aih.json")
				new_aihj, _ := sjson.Set(string(aihj), "key", OpenAI_Key)
				err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}
				// Renew ChatGPT client with key
				config = openai.DefaultConfig(OpenAI_Key)
				client = openai.NewClientWithConfig(config)
				messages = make([]openai.ChatCompletionMessage, 0)
				left_tokens = 0
				continue
			}
			// Porcess input
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			})

			// Generate a response from ChatGPT
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:    chat_mode,
					Messages: messages,
				},
			)

			if err != nil {
				fmt.Println(">>>", err)
				continue
			}

			// Record in coversation context
			if chat_completion {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: RESP,
				})
			}

			// Print the response to the terminal
			RESP = strings.TrimSpace(resp.Choices[0].Message.Content)
			used_tokens = resp.Usage.TotalTokens
			left_tokens = max_tokens - used_tokens
			printer(color_chatapi, RESP, false)

		}

		// -------------for all AI's RESP---------------

		// Persistent conversation uInput + response
		if fs, err := os.OpenFile("history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
			time_string := time.Now().Format("2006-01-02 15:04:05")
			_, err = fs.WriteString("--------------------\n")
			_, err = fs.WriteString(time_string + role + "\n\nQuestion:\n" + uInput + "\n\n")
			_, err = fs.WriteString("Answer:" + "\n" + RESP + "\n")
			if err != nil {
				panic(err)
			}
			fs.Close()
		}

		// Speak all the response RESP using the "say" command
		if speak == 1 {

			fmt.Println("speaking")
			go func() {
				switch runtime.GOOS {
				case "linux", "darwin":
					cmd := exec.Command("say", RESP)
					err = cmd.Run()
					if err != nil {
						fmt.Println(err)
					}
				case "windows":
					// to do
					_ = 1 + 1

				}

			}()
		}

		RESP = ""

	}
}

func scrollUp(textView *tview.TextView) {
	row, _ := textView.GetScrollOffset()
	if row > 0 {
		textView.ScrollTo(row-1, 0)
	}
}

func scrollPageUp(textView *tview.TextView) {
	row, _ := textView.GetScrollOffset()
	if row > 0 {
		textView.ScrollTo(row-30, 0)
	}
}

func scrollDown(textView *tview.TextView) {
	row, _ := textView.GetScrollOffset()
	textView.ScrollTo(row+1, 0)
}

func scrollPageDown(textView *tview.TextView) {
	row, _ := textView.GetScrollOffset()
	textView.ScrollTo(row+30, 0)
}

func printer(colour tcell.Color, context string, history bool) {
	app := tview.NewApplication()
	flex := tview.NewFlex()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetTextColor(colour)

	flex.AddItem(tview.NewTextView(), 0, 1, false).AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(textView, 0, 6, true).
		AddItem(tview.NewTextView(), 0, 1, false), 0, 8, false).
		AddItem(tview.NewTextView(), 0, 1, false)

	fmt.Fprintf(textView, context)
	if history {
		textView.ScrollToEnd()
	}

	// Handle 'jkgGq'
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			app.Stop()
			//	case tcell.KeyUp: // maybe use for last response
			//		scrollUp(textView)
			//	case tcell.KeyDown:
		case tcell.KeyRune:
			switch event.Rune() {
			case 'k':
				scrollUp(textView)
			case 'j':
				scrollDown(textView)
			case 'p':
				scrollPageUp(textView)
			case 'f':
				scrollPageDown(textView)
			case 'g':
				textView.ScrollToBeginning()
			case 'G':
				textView.ScrollToEnd()
			case 'q':
				app.Stop()
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}

}
