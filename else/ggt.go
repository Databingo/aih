package main

import (
//	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/go-rod/rod"
	//"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/stealth"
	//"github.com/go-rod/rod/lib/launcher"
	"github.com/manifoldco/promptui"
	"github.com/peterh/liner"
	"github.com/rivo/tview"
"jaytaylor.com/html2text"
	//openai "github.com/sashabaranov/go-openai"
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
	os.Setenv("http_proxy", Proxy)
	os.Setenv("https_proxy", Proxy)

	//proxy_url := launcher.NewUserMode().
	//	Proxy(Proxy).
	//	//Leakless(true).// indepent tab | work with UserDataDir()
	//	//UserDataDir("data").// indepent tab + data
	//	//Set("disable-default-apps").
	//	//Headless(true).
	//	MustLaunch()

	role := ".bard"

	// Open rod browser
	var browser *rod.Browser
	browser = rod.New().
		Trace(true).
		//ControlURL(proxy_url).
		Timeout(60 * 24 * time.Minute).
		//Headless(false).
		MustConnect()

//	//////////////////////0////////////////////////////
//	// Set up client of OpenAI API
//	key := gjson.Get(string(aih_json), "key")
//	OpenAI_Key := key.String()
//	config := openai.DefaultConfig(OpenAI_Key)
//	client := openai.NewClientWithConfig(config)
//	messages := make([]openai.ChatCompletionMessage, 0)

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
		//page_bard = browser.MustPage("https://google.com")
		page_bard = stealth.MustPage(browser)
		page_bard.MustNavigate("https://google.com")
		relogin_bard = false
		//for {
		//	if page_bard.MustHasX("//textarea[@id='APjFqb']") {
		//		relogin_bard = false
		//		break
		//	}
		//	time.Sleep(time.Second)
		//}
		if relogin_bard == true {
			fmt.Println("✘ Bard")
		}
		if relogin_bard == false {
			fmt.Println("✔ Bard")
			for {
				select {
				case question := <-channel_bard:
					//page_bard.MustActivate()
					//page_bard.MustElementX("//textarea[@id='APjFqb']").MustWaitVisible().MustSelectAllText().MustInput(question).MustType(input.Enter)
		                        page_bard.MustNavigate("https://google.com/search?q="+question)
					if role == ".all" {
						channel_bard <- "click_bard"
					} else { page_bard.MustActivate() }
					page_bard.MustElementX("//div[@id='hdtb-tls']").MustWaitVisible()
					search_dive := page_bard.MustElementX("//div[@id='search']").MustWaitVisible()
					//response := search_dive.MustElementsX("//div[@class='kvH3mc BToiNc UK95Uc']")
					//response := search_dive.MustElementsX("//div[@jscontroller='SC71Yd']")
					//response := img.MustElementX("parent::div/parent::div").MustWaitVisible()
					answerH, _ := search_dive.HTML()
					answer, _ :=  html2text.FromString(answerH, html2text.Options{PrettyTables: true})

					//answer := search_dive.MustText()
				//	var answer = ""
				//	for _, i := range response {
                                //         //link := i.MustElementX("//div[@class='Z26q7c UK95Uc jGGQ5e']//a/@href")
                                //         link := i.MustElementX("//a/@href")
				//	 //describe := i.MustElementX("//div[@class='Z26q7c UK95Uc' and @data-sncf='1']").MustWaitVisible()
				//	 //describe := i.MustElementX("//div[@data-sncf='1']").MustWaitVisible()
				//	 describe := i.MustElementX("//div[@class='VwiC3b yXK7lf MUxGbd yDYNvb lyLwlc']").MustWaitVisible()
				//	answer += link.MustText()
				//	answer += describe.MustText()
				//	answer += "\n\n"
				//       }
					channel_bard <- answer
					//channel_bard <- "answer"
				}
			}
		}

	}()

//	//////////////////////c2////////////////////////////
//	// Set up client of Claude (Rod version)
//	var page_claude *rod.Page
//	var relogin_claude = true
//	channel_claude := make(chan string)
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				relogin_claude = true
//			}
//		}()
//		//page_claude = browser.MustPage("https://bing.com")
//		page_claude = stealth.MustPage(browser)
//		page_claude.MustNavigate("https://bing.com")
//		relogin_claude = false
//		//for {
//		//	if page_claude.MustHasX("//input[@id='sb_form_q']") {
//		//		relogin_claude = false
//		//		//fmt.Println("1✔ Claude")
//		//		break
//		//	}
//		//	time.Sleep(time.Second)
//		//}
//
//		if relogin_claude == true {
//			fmt.Println("✘ Claude")
//		}
//		if relogin_claude == false {
//			fmt.Println("✔ Claude")
//			for {
//				select {
//				case question := <-channel_claude:
//					//page_claude.MustElementX("//input[@id='sb_form_q']").MustWaitVisible().MustSelectAllText().MustInput(question).MustType(input.Enter)
//		                        page_claude.MustNavigate("https://bing.com/search?q="+question)
//					if role == ".all" {
//						channel_claude <- "click_claude"
//					} else {
//					page_claude.MustActivate()
//				       }
//					page_claude.MustElement("svg path[d='M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z']").MustWaitVisible()
//					channel_claude <- "answer"
//				}
//			}
//		}
//
//	}()
//
//	//////////////////////c3////////////////////////////
//	// Set up client of Huggingchat (Rod version)
//	var page_hc *rod.Page
//	var relogin_hc = true
//	channel_hc := make(chan string)
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				relogin_hc = true
//			}
//		}()
//		//page_hc = browser.MustPage("https://baidu.com")
//		page_hc = stealth.MustPage(browser)
//		page_hc.MustNavigate("https://baidu.com")
//		relogin_hc = false
//		//for {
//		//	if page_hc.MustHasX("//input[@id='kw']") {
//		//		relogin_hc = false
//		//		break
//		//	}
//		//	time.Sleep(time.Second)
//		//}
//		if relogin_hc == true {
//			fmt.Println("✘ HuggingChat")
//		}
//		if relogin_hc == false {
//			fmt.Println("✔ HuggingChat")
//			for {
//				select {
//				case question := <-channel_hc:
//					//page_hc.MustElementX("//input[@id='kw']").MustWaitVisible().MustSelectAllText().MustInput(question).MustType(input.Enter)
//		                        page_hc.MustNavigate("https://baidu.com/s?wd="+question)
//					if role == ".all" {
//						channel_hc <- "click_hc"
//					} else {
//					page_hc.MustActivate()
//				       }
//
//					page_hc.MustElementX("//a[contains(text(), '更多')]")
//					channel_hc <- "baidu"
//				}
//			}
//		}
//
//	}()
//
//	//////////////////////c4////////////////////////////
//	// Set up client of chatgpt (rod version)
//	var page_chatgpt *rod.Page
//	var relogin_chatgpt = true
//	channel_chatgpt := make(chan string)
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				relogin_chatgpt = true
//			}
//		}()
//		//page_chatgpt = browser.MustPage("https://zhihu.com/topic/19585187/top-answers")
//		page_chatgpt = stealth.MustPage(browser)
//		page_chatgpt.MustNavigate("https://zhihu.com/topic/19585187/top-answers")
//		relogin_chatgpt = false
//		//channel_chatgpt_verify := make(chan string)
//		//go func() {
//		//	for {
//		//		if page_chatgpt.MustHasX("//div[contains(text(), '其他方式登录')]") {
//		//	fmt.Println("Found Click ✘ ")
//		//			//page_chatgpt.MustElement("button svg path[d='M18.22 19.28a.75.75 0 1 0 1.06-1.06L13.06 12l6.22-6.22a.75.75 0 0 0-1.06-1.06L12 10.94 5.78 4.72a.75.75 0 0 0-1.06 1.06L10.94 12l-6.22 6.22a.75.75 0 1 0 1.06 1.06L12 13.06l6.22 6.22Z']").MustWaitVisible().MustClick()
//		//			page_chatgpt.MustElementX("button[@aria-label='关闭']").MustWaitVisible().MustClick()
//		//	fmt.Println("Click ✘ ")
//		//			break
//		//		}
//		//		time.Sleep(time.Second)
//		//	}
//		//	close(channel_chatgpt_verify)
//		//}()
//		//for {
//		//	//if page_chatgpt.MustHasX("//div[contains(text(), '其他方式登录')]") {
//		//	//   page_chatgpt.MustElement("svg path[d='M18.22 19.28a.75.75 0 1 0 1.06-1.06L13.06 12l6.22-6.22a.75.75 0 0 0-1.06-1.06L12 10.94 5.78 4.72a.75.75 0 0 0-1.06 1.06L10.94 12l-6.22 6.22a.75.75 0 1 0 1.06 1.06L12 13.06l6.22 6.22Z']").MustClick()
//		//	//	relogin_chatgpt = false
//
//		//	//	break
//		//	//}
//		//	if page_chatgpt.MustHasX("//input[@id='Popover1-toggle']") {
//		//		relogin_chatgpt = false
//		//		break
//		//	}
//		//	time.Sleep(time.Second)
//		//}
//
//		if relogin_chatgpt == true {
//			fmt.Println("✘ ChatGPT")
//			// Automatic login
//			//page_chatgpt.MustElementX("//div[contains(text(), 'Welcome to ChatGPT')] | //h2[contains(text(), 'Get started')]").MustWaitVisible()
//			//page_chatgpt.MustElementX("//div[not(contains(@class, 'mb-4')) and contains(text(), 'Log in')]").MustClick()
//			//utils.Sleep(1.5)
//			//page_chatgpt.MustElementX("//input[@id='username']").MustWaitVisible().MustInput(chatgpt_user)
//			//utils.Sleep(1.5)
//			//page_chatgpt.MustElementX("//button[contains(text(), 'Continue')]").MustClick()
//			//utils.Sleep(1.5)
//			//page_chatgpt.MustElementX("//input[@id='password']").MustWaitVisible().MustInput(chatgpt_password)
//			//utils.Sleep(1.5)
//			//page_chatgpt.MustElementX("//button[not(contains(@aria-hidden, 'true')) and contains(text(), 'Continue')]").MustClick()
//			////page_chatgpt.MustElementX("//h4[contains(text(), 'This is a free research preview.')]").MustWaitVisible()
//			////utils.Sleep(1.5)
//			////page_chatgpt.MustElementX("//button/div[contains(text(), 'Next')]").MustClick()
//			////page_chatgpt.MustElementX("//h4[contains(text(), 'How we collect data')]").MustWaitVisible()
//			////utils.Sleep(1.5)
//			////page_chatgpt.MustElementX("//button/div[contains(text(), 'Next')]").MustClick()
//			////page_chatgpt.MustElementX("//h4[contains(text(), 'love your feedback!')]").MustWaitVisible()
//			////utils.Sleep(1.5)
//			////page_chatgpt.MustElementX("//button/div[contains(text(), 'Done')]").MustClick()
//			////utils.Sleep(1.5)
//			//page_chatgpt.MustElementX("//a[contains(text(), 'New chat')]").MustWaitVisible().MustClick()
//			//page_chatgpt.MustElementX("//textarea[@id='prompt-textarea']").MustWaitVisible()
//			//utils.Sleep(1.5)
//			//page_chatgpt.MustElementX("//textarea[@id='prompt-textarea']").MustInput("hello")
//			//utils.Sleep(1.5)
//			//page_chatgpt.MustElement("svg:last-of-type path[d='M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15']").MustWaitVisible()
//			//fmt.Println("Retry icon show")
//			//page_chatgpt.MustElementX("(//div[contains(@class, 'group w-full')])[last()]").MustText()
//			//fmt.Println("✔ ChatGPT Ready")
//		}
//		if relogin_chatgpt == false {
//			fmt.Println("✔ ChatGPT")
//			for {
//				select {
//				case question := <-channel_chatgpt:
//		                        page_chatgpt.MustNavigate("https://zhihu.com/search?type=content&q="+question)
//					if role == ".all" {
//						channel_chatgpt <- "click_chatgpt"
//					} else {
//					page_chatgpt.MustActivate()
//				       }
//					page_chatgpt.MustElementX("//div[contains(text(), '筛选')]").MustWaitVisible()
//					channel_chatgpt <- "zhihu"
//				}
//			}
//		}
//	}()

	// Clean screen
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
	//chat_mode := openai.GPT3Dot5Turbo
	//chat_completion := true

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
			Liner.Close()
			/// exit_safe()
			if relogin_bard == false {
				page_bard.MustClose()
			}
			//if relogin_chatgpt == false {
			//	page_chatgpt.MustClose()
			//}
			//if relogin_claude == false {
			//	page_claude.MustClose()
			//}
			//if relogin_hc == false {
			//	page_hc.MustClose()
			//}
			close(channel_bard)
			//close(channel_chatgpt)
			//close(channel_claude)
			//close(channel_hc)
			Liner.Close()
			syscall.Exit(0)
			//os.Exit(0)
		case ".help":
			fmt.Println("                           ")
			fmt.Println("                 Welcome to Aih!                             ")
			fmt.Println("------------------------------------------------------------ ")
			fmt.Println(" .               Select AI mode of Bard/Bing/ChatGPT/Claude")
			fmt.Println(" .proxy          Set proxy")
			fmt.Println(" .key            Set key of ChatGPT API")
			fmt.Println(" <<              Start multiple lines input")
			fmt.Println(" >>              End multiple lines input")
			fmt.Println(" ↑               Previous input")
			fmt.Println(" ↓               Next input")
			fmt.Println(" .c or .clear    Clear screen")
			fmt.Println(" .h or .history  Show history")
			fmt.Println(" j               Scroll down")
			fmt.Println(" k               Scroll up")
			fmt.Println(" gg              Scroll top")
			fmt.Println(" G               Scroll bottom")
			fmt.Println(" q or Enter      Back to conversation")
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
			//	exit_safe()
			if relogin_bard == false {
				page_bard.MustClose()
			}
			//if relogin_chatgpt == false {
			//	page_chatgpt.MustClose()
			//}
			//if relogin_claude == false {
			//	page_claude.MustClose()
			//}
			//if relogin_hc == false {
			//	page_hc.MustClose()
			//}
			close(channel_bard)
			//close(channel_chatgpt)
			//close(channel_claude)
			//close(channel_hc)
			Liner.Close()
			syscall.Exit(0)
			//os.Exit(0)
		case ".new":
			// For role .chat
			//conversation_id = ""
			//parent_id = ""
			// For role .chatapi
			//messages = make([]openai.ChatCompletionMessage, 0)
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
				//chat_mode = openai.GPT3Dot5Turbo
				max_tokens = 4097
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				//chat_completion = true
				continue
			case "ChatGPT API gpt-4 8K Prompt, $0.03/1K tokens":
				role = ".chatapi"
				//chat_mode = openai.GPT4
				max_tokens = 8192
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				//chat_completion = false
				continue
			case "ChatGPT API gpt-4 8K Completion, $0.06/1K tokens":
				role = ".chatapi"
				//chat_mode = openai.GPT4
				max_tokens = 8192
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				//chat_completion = true
				continue
			case "ChatGPT API gpt-4 32K Prompt, $0.06/1K tokens":
				role = ".chatapi"
				//chat_mode = openai.GPT432K
				max_tokens = 32768
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				//chat_completion = false
				continue
			case "ChatGPT API gpt-4 32K Completion, $0.12/1K tokens":
				role = ".chatapi"
				//chat_mode = openai.GPT432K
				max_tokens = 32768
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				//chat_completion = true
				continue
			case "Exit":
				continue
			}
		case ".key":
			prom := promptui.Select{
				Label: "Select:",
				Size:  6,
				Items: []string{
					//"Set Bard Cookie",
					//"Set ChatGPT Cookie",
					//"Set Claude Cookie",
					//"Set HuggingChat Cookie",
					"Set ChatGPT API Key",
					"Exit",
				},
			}

			_, keyy, err := prom.Run()
			if err != nil {
				panic(err)
			}

			switch keyy {
			//case "Set Bard Cookie":
			//	role = ".bard"
			//	goto BARD
			//case "Set ChatGPT Cookie":
			//	role = ".chat"
			//	goto CHAT
			//case "Set Claude Cookie":
			//	role = ".claude"
			//	goto CLAUDE
			//case "Set HuggingChat Cookie":
			//	role = ".huggingchat"
			//	goto HUGGINGCHAT
			case "Set ChatGPT API Key":
				//OpenAI_Key = ""
				role = ".chatapi"
				//goto CHATAPI
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

		//ALL-IN-ONE:
		if role == ".all" {
			if relogin_bard == false {
				channel_bard <- userInput
				<-channel_bard
			}
			//if relogin_chatgpt == false {
			//	channel_chatgpt <- userInput
			//	<-channel_chatgpt
			//}
			//if relogin_claude == false {
			//	channel_claude <- userInput
			//	<-channel_claude
			//}
			//if relogin_hc == false {
			//	channel_hc <- userInput
			//	<-channel_hc
			//}

			if relogin_bard == false {
				answer_bard := <-channel_bard
				RESP += "\n\n---------------- bard answer ----------------\n"
				RESP += strings.TrimSpace(answer_bard)
			}
			//if relogin_chatgpt == false {
			//	answer_chatgpt := <-channel_chatgpt
			//	RESP += "\n\n---------------- chatgpt answer ----------------\n"
			//	RESP += strings.TrimSpace(answer_chatgpt)
			//}
			//if relogin_claude == false {
			//	answer_claude := <-channel_claude
			//	RESP += "\n\n---------------- claude answer ----------------\n"
			//	RESP += strings.TrimSpace(answer_claude)
			//}
			//if relogin_hc == false {
			//	answer_hc := <-channel_hc
			//	RESP += "\n\n---------------- huggingchat answer ----------------\n"
			//	RESP += strings.TrimSpace(answer_hc)
			//}
			printer(color_chat, RESP, false)

		}
		if role == ".bard" {
			if relogin_bard == true {
				fmt.Println("✘ Bard")
			} else {
				channel_bard <- userInput
				answer := <-channel_bard

				// Print the response to the terminal
				RESP = strings.TrimSpace(answer)
				printer(color_bard, RESP, false)
			}

		}

		//	CLAUDE:
//		// Check role for correct actions
//		if role == ".claude" {
//			if relogin_claude == true {
//				fmt.Println("✘ Claude")
//			} else {
//				channel_claude <- userInput
//				answer := <-channel_claude
//
//				RESP = strings.TrimSpace(answer)
//				printer(color_claude, RESP, false)
//			}
//
//		}
//		//	CHAT:
//		if role == ".chat" {
//			if relogin_chatgpt == true {
//				fmt.Println("✘ ChatGPT")
//			} else {
//				channel_chatgpt <- userInput
//				answer := <-channel_chatgpt
//
//				// Print the response to the terminal
//				RESP = strings.TrimSpace(answer)
//				printer(color_chatapi, RESP, false)
//			}
//
//		}
//
//		//	HUGGINGCHAT:
//		if role == ".huggingchat" {
//			if relogin_hc == true {
//				fmt.Println("✘ HuggingChat")
//			} else {
//				channel_hc <- userInput
//				answer := <-channel_hc
//
//				// Print the response to the terminal
//				RESP = strings.TrimSpace(answer)
//				printer(color_huggingchat, RESP, false)
//			}
//
//		}
//	CHATAPI:
//		if role == ".chatapi" {
//			// Check ChatGPT API Key
//			if OpenAI_Key == "" {
//				OpenAI_Key, _ = Liner.Prompt("Please input your OpenAI Key: ")
//				if OpenAI_Key == "" {
//					continue
//				}
//				aihj, err := ioutil.ReadFile("aih.json")
//				new_aihj, _ := sjson.Set(string(aihj), "key", OpenAI_Key)
//				err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
//				if err != nil {
//					fmt.Println("Save failed.")
//				}
//				// Renew ChatGPT client with key
//				config = openai.DefaultConfig(OpenAI_Key)
//				client = openai.NewClientWithConfig(config)
//				messages = make([]openai.ChatCompletionMessage, 0)
//				left_tokens = 0
//				continue
//			}
//			// Porcess input
//			messages = append(messages, openai.ChatCompletionMessage{
//				Role:    openai.ChatMessageRoleUser,
//				Content: userInput,
//			})
//
//			// Generate a response from ChatGPT
//			resp, err := client.CreateChatCompletion(
//				context.Background(),
//				openai.ChatCompletionRequest{
//					Model:    chat_mode,
//					Messages: messages,
//				},
//			)
//
//			if err != nil {
//				fmt.Println(">>>", err)
//				continue
//			}
//
//			// Record in coversation context
//			if chat_completion {
//				messages = append(messages, openai.ChatCompletionMessage{
//					Role:    openai.ChatMessageRoleUser,
//					Content: RESP,
//				})
//			}
//
//			// Print the response to the terminal
//			RESP = strings.TrimSpace(resp.Choices[0].Message.Content)
//			used_tokens = resp.Usage.TotalTokens
//			left_tokens = max_tokens - used_tokens
//			printer(color_chatapi, RESP, false)
//
//		}

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

func scrollDown(textView *tview.TextView) {
	row, _ := textView.GetScrollOffset()
	textView.ScrollTo(row+1, 0)
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
