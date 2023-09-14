package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/creack/pty"
	"github.com/gdamore/tcell/v2"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"github.com/go-rod/stealth"
//	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/peterh/liner"
	"github.com/rivo/tview"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var trace = false

// var trace = true
var userInput string
var color_bard = tcell.ColorDarkCyan
var color_bing = tcell.ColorDarkMagenta
var color_chat = tcell.ColorWhite
var color_chatapi = tcell.ColorWhite
var color_claude = tcell.ColorYellow

//var color_huggingchat = tcell.ColorDarkMagenta

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

func sprint(s string) {
	if userInput != ".v" && userInput != ".vi" && userInput != ".vim" {
		fmt.Println(s)
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

	// Save miv locally
	switch runtime.GOOS {
	case "linux", "darwin":
		os.WriteFile(".mvi", vi, 0755)
	case "windows":
		os.WriteFile(".mvi.exe", vi, 0755)
	}

	// Create prompt for user input
	Liner := liner.NewLiner()
	defer Liner.Close()

	// Use RESP for record response per time
	var RESP string

	// Read Aih Configure
	aih_json, err := ioutil.ReadFile(".aih.json")
	if err != nil {
		err = ioutil.WriteFile(".aih.json", []byte(""), 0644)
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
	//var browser *rod.Browser
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
	//browser.ServeMonitor(":7777")

	//////////////////////0////////////////////////////
	// Set up client of OpenAI API
	key := gjson.Get(string(aih_json), "key")
	OpenAI_Key := key.String()
	config := openai.DefaultConfig(OpenAI_Key)
	client := openai.NewClientWithConfig(config)
	messages := make([]openai.ChatCompletionMessage, 0)

	//////////////////////c1////////////////////////////
	go Bard()

	//////////////////////c2////////////////////////////
	// Set up client of Claude (Rod version)
	go Claude2()
	//var page_claude *rod.Page
	//var relogin_claude = true
	//channel_claude := make(chan string)
	//go func() {
	//	defer func() {
	//		if err := recover(); err != nil {
	//			relogin_claude = true
	//		}
	//	}()
	//	//page_claude = browser.MustPage("https://claude.ai")
	//	page_claude = stealth.MustPage(browser)
	//	page_claude.MustNavigate("https://claude.ai/api/organizations").MustWaitLoad()
	//	org_json := page_claude.MustElementX("//pre").MustText()
	//	org_uuid := gjson.Get(string(org_json), "0.uuid").String()
	//	time.Sleep(6 * time.Second)

	//	new_uuid := uuid.New().String()
	//	new_uuid_url := "https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations"
	//	create_new_converastion_json := `{"uuid":"` + new_uuid + `","name":""}`
	//	create_new_converastion_js := `
	//	 (new_uuid_url, sdata) => {
	//	 var xhr = new XMLHttpRequest();
	//	 xhr.open("POST", new_uuid_url);
	//	 xhr.setRequestHeader('Content-Type', 'application/json');
	//	 xhr.setRequestHeader('Referer', 'https://claude.ai/chats');
	//	 xhr.setRequestHeader('Origin', 'https://claude.ai');
	//	 xhr.setRequestHeader('TE', 'trailers');
	//	 xhr.setRequestHeader('DNT', '1');
	//	 xhr.setRequestHeader('Connection', 'keep-alive');
	//	 xhr.setRequestHeader('Accept', 'text/event-stream, text/event-stream');
	//	 xhr.onreadystatechange = function() {
	//	     if (xhr.readyState == XMLHttpRequest.DONE) { 
	//	         var res_text = xhr.responseText;
	//	         console.log(res_text);
	//	        } 
	//	     }
	//         console.log(sdata);
	//	 xhr.send(sdata);
	//	}
	//	`
	//	page_claude.MustEval(create_new_converastion_js, new_uuid_url, create_new_converastion_json).Str()
	//	// posted new conversation uuid
	//	time.Sleep(3 * time.Second) // delay to simulate human being

	//	var record_chat_messages string
	//	var response_chat_messages string
	//	for i := 1; i <= 20; i++ {
	//		create_json := page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()

	//		message_uuid := gjson.Get(string(create_json), "uuid").String()
	//		if message_uuid == new_uuid {
	//			//fmt.Println("create_conversation success...")
	//			relogin_claude = false
	//			record_chat_messages = gjson.Get(string(create_json), "chat_messages").String()
	//			break
	//		}
	//		time.Sleep(2 * time.Second)
	//	}

	//	if relogin_claude == true {
	//		sprint("✘ Claude")
	//		//page_claude.MustPDF("./tmp/Claude✘.pdf")
	//	}
	//	if relogin_claude == false {
	//		sprint("✔ Claude")
	//		for {
	//			select {
	//			case question := <-channel_claude:
	//				question = strings.Replace(question, "\r", "\n", -1)
	//				question = strings.Replace(question, "\"", "\\\"", -1)
	//				question = strings.Replace(question, "\n", "\\n", -1)
	//				question = strings.TrimSuffix(question, "\n")
	//				// re-activate
	//				page_claude.MustNavigate("https://claude.ai/api/account/statsig/" + org_uuid).MustWaitLoad()
	//				record_json := page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()
	//				record_chat_messages = gjson.Get(string(record_json), "chat_messages").String()
	//				record_count := gjson.Get(string(response_chat_messages), "#").Int()
	//				page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid).MustWaitLoad()
	//				time.Sleep(2 * time.Second) // delay to simulate human being
	//				//question = strings.Replace(question, `"`, `\"`, -1) // escape " in input text when code into json

	//				d := `{"completion":{"prompt":"` + question + `","timezone":"Asia/Shanghai","model":"claude-2"},"organization_uuid":"` + org_uuid + `","conversation_uuid":"` + new_uuid + `","text":"` + question + `","attachments":[]}`
	//				//fmt.Println(d)
	//				js := `
	//	                               (sdata, new_uuid) => {
	//	                               var xhr = new XMLHttpRequest();
	//	                               xhr.open("POST", "https://claude.ai/api/append_message");
	//	                               xhr.setRequestHeader('Content-Type', 'application/json');
	//	                               xhr.setRequestHeader('Referer', 'https://claude.ai/chat/new_uuid');
	//	                               xhr.setRequestHeader('Origin', 'https://claude.ai');
	//	                               xhr.setRequestHeader('TE', 'trailers');
	//	                               xhr.setRequestHeader('Connection', 'keep-alive');
	//	                               xhr.setRequestHeader('Accept', 'text/event-stream, text/event-stream');
	//	                               xhr.onreadystatechange = function() {
	//	                                   if (xhr.readyState == XMLHttpRequest.DONE) { 
	//	                                       var res_text = xhr.responseText;
	//	                                       console.log(res_text);
	//	                                   }
	//	                               }
	//	                               console.log(sdata);
	//	                               xhr.send(sdata);
	//	                               }
	//	                              `
	//				page_claude.MustEval(js, d, new_uuid).Str()
	//				fmt.Println("Claude generating...")
	//				//if role == ".all" {
	//				//	channel_claude <- "click_claude"
	//				//}
	//				time.Sleep(3 * time.Second) // delay to simulate human being

	//				// wait answer
	//				var claude_response = false
	//				var response_json string
	//				for i := 1; i <= 20; i++ {
	//					if page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustHasX("//pre") {
	//						response_json = page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()
	//						response_chat_messages = gjson.Get(string(response_json), "chat_messages").String()
	//						count := gjson.Get(string(response_chat_messages), "#").Int()

	//						if response_chat_messages != record_chat_messages && count == record_count+2 {
	//							claude_response = true
	//							record_chat_messages = response_chat_messages
	//							break
	//						}
	//					}
	//					time.Sleep(3 * time.Second)
	//				}
	//				if claude_response == true {
	//					count := gjson.Get(string(response_chat_messages), "#").Int()
	//					answer := gjson.Get(string(response_json), "chat_messages.#(index=="+strconv.FormatInt(count-1, 10)+").text").String()
	//					channel_claude <- answer
	//				} else {
	//					channel_claude <- "✘✘  Claude, Please check the internet connection and verify login status."
	//					relogin_claude = true
	//					//page_claude.MustPDF("./tmp/Claude✘.pdf")
	//				}
	//			}
	//		}
	//	}

	//}()

	////////////////////////c3////////////////////////////
	//// Set up client of Huggingchat (Rod version)
	//var page_hc *rod.Page
	//var relogin_hc = true
	//channel_hc := make(chan string)
	//go func() {
	//	defer func() {
	//		if err := recover(); err != nil {
	//			relogin_hc = true
	//		}
	//	}()
	//	//page_hc = browser.MustPage("https://huggingface.co/chat")
	//	page_hc = stealth.MustPage(browser)
	//	page_hc.MustNavigate("https://huggingface.co/chat")
	//	for i := 1; i <= 30; i++ {
	//		if page_hc.MustHasX("//button[contains(text(), 'Sign Out')]") {
	//			relogin_hc = false
	//			break
	//		}
	//		time.Sleep(time.Second)
	//	}
	//	if relogin_hc == true {
	//		sprint("✘ HuggingChat")
	//		//page_hc.MustPDF("./tmp/HuggingChat✘.pdf")
	//	}
	//	if relogin_hc == false {
	//		sprint("✔ HuggingChat")
	//		for {
	//			select {
	//			case question := <-channel_hc:
	//				//page_hc.MustActivate()
	//				//fmt.Println("HuggingChat received question...", question)
	//				for i := 1; i <= 20; i++ {
	//					if page_hc.MustHasX("//textarea[@enterkeyhint='send']") {
	//						page_hc.MustElementX("//textarea[@enterkeyhint='send']").MustInput(question)
	//						break
	//					}
	//					time.Sleep(time.Second)
	//				}
	//				//fmt.Println("HuggingChat input typed...")
	//				for i := 1; i <= 20; i++ {
	//					if page_hc.MustHas("button svg path[d='M27.71 4.29a1 1 0 0 0-1.05-.23l-22 8a1 1 0 0 0 0 1.87l8.59 3.43L19.59 11L21 12.41l-6.37 6.37l3.44 8.59A1 1 0 0 0 19 28a1 1 0 0 0 .92-.66l8-22a1 1 0 0 0-.21-1.05Z']") {
	//						page_hc.MustElement("button svg path[d='M27.71 4.29a1 1 0 0 0-1.05-.23l-22 8a1 1 0 0 0 0 1.87l8.59 3.43L19.59 11L21 12.41l-6.37 6.37l3.44 8.59A1 1 0 0 0 19 28a1 1 0 0 0 .92-.66l8-22a1 1 0 0 0-.21-1.05Z']").MustClick()
	//						break
	//					}
	//					time.Sleep(time.Second)
	//				}
	//				fmt.Println("HuggingChat generating...")
	//				//page_hc.MustActivate() // Sometime three dot to hang
	//				//if role == ".all" {
	//				//	channel_hc <- "click_hc"
	//				//}
	//				// Check too much traffic
	//				channel_hc_check := make(chan string)
	//				go func() {
	//					for i := 1; i <= 20; i++ {
	//						if page_hc.MustHasX("//*[contains(text(), 'Too much traffic, please try again')]") {
	//							channel_hc <- "✘✘  HuggingChat, Please check the internet connection and verify login status. Traffic."
	//							fmt.Println("HuggingChat too much traffic...")
	//							relogin_hc = true
	//							close(channel_hc_check)
	//							break
	//						}
	//						time.Sleep(1 * time.Second)
	//					}
	//				}()
	//				// Check redirect to conversation
	//				for i := 1; i <= 20; i++ {
	//					info := page_hc.MustInfo()
	//					if strings.HasPrefix(info.URL, "https://huggingface.co/chat/conversation") {
	//						break
	//					}
	//					time.Sleep(1 * time.Second)
	//				}

	//				// stop_icon
	//				var stop_icon_disappear = false
	//				for i := 1; i <= 80; i++ {
	//					//for {
	//					if page_hc.MustHas("svg path[d='M24 6H8a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2Z']") {
	//						stop_icon_disappear = false
	//					} else {
	//						if page_hc.MustHasX("//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')]") {

	//							if page_hc.MustElementX("(//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')])[last()]").MustHasX("following-sibling::div[1]") {
	//								stop_icon_disappear = true
	//								break
	//							}
	//						}

	//					}
	//					if relogin_hc == true {
	//						continue
	//					}
	//					time.Sleep(time.Second)
	//				}
	//				if stop_icon_disappear == true {
	//					//page_hc.MustHasX("//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')]")
	//					page_hc.MustElementX("//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')]").MustWaitVisible()
	//					img := page_hc.MustElementX("(//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')])[last()]")
	//					content := img.MustElementX("following-sibling::div[1]")
	//					answer := content.MustText()
	//					channel_hc <- answer
	//				} else {
	//					channel_hc <- "✘✘  HuggingChat, Please check the internet connection and verify login status."
	//					relogin_hc = true
	//					//page_hc.MustPDF("./tmp/HuggingChat✘.pdf")

	//				}
	//			}
	//		}
	//	}

	//}()

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
				if page_chatgpt.MustHasX("//h2[contains(text(), 'Your session has expired')]") {
					relogin_chatgpt = true
					close(channel_chatgpt_tips)
					break
				}
				time.Sleep(time.Second)
			}
		}()

		//for i := 1; i <= 30; i++ {
		//	if page_chatgpt.MustHasX("//div[contains(text(), 'Okay, let')]") {
		//		page_chatgpt.MustElementX("//div[contains(text(), 'Okay, let')]").MustWaitVisible().MustClick()
		//	}
		//	time.Sleep(time.Second)
		//}
		//for i := 1; i <= 30; i++ {
		//	if page_chatgpt.MustHasX("//h2[contains(text(), 'Your session has expired')]") {
		//		relogin_chatgpt = true
		//		break
		//	}
		//	if page_chatgpt.MustHasX("//div[contains(text(), 'Something went wrong. If this issue persists please')]") {
		//		//fmt.Println("ChatGPT web error")
		//		//channel_chatgpt <- "✘✘  ChatGPT, Please check the internet connection and verify login status."
		//		relogin_chatgpt = true
		//		break
		//		//page_chatgpt.MustPDF("ChatGPT✘.pdf")
		//	}
		//	time.Sleep(time.Second)
		//}
		for i := 1; i <= 30; i++ {
			if page_chatgpt.MustHasX("//textarea[@id='prompt-textarea']") && !page_chatgpt.MustHasX("//h2[contains(text(), 'Your session has expired')]") {
				relogin_chatgpt = false
				break
			}
			time.Sleep(time.Second)
		}

		if relogin_chatgpt == true {
			sprint("✘ ChatGPT")
			//page_chatgpt.MustPDF("./tmp/ChatGPT✘.pdf")
		}
		if relogin_chatgpt == false {
			sprint("✔ ChatGPT")
			for {
				select {
				case question := <-channel_chatgpt:
					//page_chatgpt.MustActivate()
					for i := 1; i <= 20; i++ {
						if page_chatgpt.MustHasX("//textarea[@id='prompt-textarea']") {
							page_chatgpt.MustElementX("//textarea[@id='prompt-textarea']").MustInput(question)
							break
						}
						time.Sleep(1 * time.Second)
					}
					for i := 1; i <= 20; i++ {
						if page_chatgpt.MustHasX("//textarea[@id='prompt-textarea']/..//button") {
							page_chatgpt.MustElementX("//textarea[@id='prompt-textarea']/..//button").MustClick()
							break
						}
						time.Sleep(1 * time.Second)
					}
					fmt.Println("ChatGPT generating...")
					//page_chatgpt.MustActivate()
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
						//page_chatgpt.MustPDF("./tmp/ChatGPT✘.pdf")

					}
				}
			}
		}
	}()

	//////////////////////c5////////////////////////////
	// Set up client of falcon180-180B (Rod version)
	var page_falcon180 *rod.Page
	var relogin_falcon180 = true
	channel_falcon180 := make(chan string)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				relogin_falcon180 = true
			}
		}()
		//page_hc = browser.MustPage("https://huggingface.co/chat")
		page_falcon180 = stealth.MustPage(browser)
		//page_falcon180.MustNavigate("https://huggingface.co/spaces/tiiuae/falcon-180b-demo")
		page_falcon180.MustNavigate("https://tiiuae-falcon-180b-demo.hf.space/?__theme=light")
		for i := 1; i <= 30; i++ {
			if page_falcon180.MustHasX("//textarea[@data-testid='textbox']") {
				relogin_falcon180 = false
				break
			}
			time.Sleep(time.Second)
		}
		if relogin_falcon180 == true {
			sprint("✘ Falcon180")
			//page_hc.MustPDF("./tmp/HuggingChat✘.pdf")
		}
		if relogin_falcon180 == false {
			sprint("✔ Falcon180")
			for {
				select {
				case question := <-channel_falcon180:
					//page_hc.MustActivate()
					//fmt.Println("Falcon180 received question...", question)
					for i := 1; i <= 20; i++ {
						if page_falcon180.MustHasX("//textarea[@data-testid='textbox']") {
							page_falcon180.MustElementX("//textarea[@data-testid='textbox']").MustInput(question)
							break
						}
						time.Sleep(time.Second)
					}
					//fmt.Println("Falcon180 input typed...")
					for i := 1; i <= 20; i++ {
						//if page_falcon180.MustHasX("//button[contains(text(), 'Submit')]") {
						if page_falcon180.MustHasX("//button[@id='component-17']") {
							page_falcon180.MustElementX("//button[@id='component-17']").MustClick()
							break
						}
						time.Sleep(time.Second)
					}
					fmt.Println("Falcon180 generating...")
					//page_falcon180.MustActivate() // Sometime three dot to hang
					//if role == ".all" {
					//	channel_falcon180 <- "click_falcon180"
					//}
					//// Check Error
					//channel_falcon180_check := make(chan string)
					//go func() {
					//	for i := 1; i <= 20; i++ {
					//		if page_falcon180.MustHasX("//*[contains(text(), 'Too much traffic, please try again')]") {
					//			channel_falcon180 <- "✘✘ Falcon180, Please check the internet connection and verify login status. Traffic."
					//			fmt.Println("Falcon180 too much traffic...")
					//			relogin_falcon180 = true
					//			close(channel_falcon180_check)
					//			break
					//		}
					//		time.Sleep(1 * time.Second)
					//	}
					//}()

					// stop_icon
					var stop_icon_disappear = false
					err := rod.Try(func() {
						page_falcon180.Timeout(10 * time.Second).MustElementX("//button[@id='component-18']").MustWaitVisible().CancelTimeout()
					})
					if err == nil {
						err = rod.Try(func() {
							page_falcon180.Timeout(80 * time.Second).MustElementX("//button[@id='component-18']").MustWaitInvisible().CancelTimeout()
						})
						if err == nil {
							stop_icon_disappear = true
						} else {
							//fmt.Println("err::::", err)
						}
					} else {
						//fmt.Println("err::", err)
					}

					if stop_icon_disappear == true {
						answer := page_falcon180.MustElementX("(//div[@data-testid='bot'])[last()]").MustText()
						channel_falcon180 <- answer
					} else {
						channel_falcon180 <- "✘✘  Falcon180, Please check the internet connection and verify login status."
						relogin_falcon180 = true
						//page_hc.MustPDF("./tmp/HuggingChat✘.pdf")

					}
				}
			}
		}

	}()

	//////////////////////c6////////////////////////////
	// Set up client of Llama2 (Rod version)
	var page_llama2 *rod.Page
	var relogin_llama2 = true
	channel_llama2 := make(chan string)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				relogin_llama2 = true
			}
		}()
		//page_hc = browser.MustPage("https://huggingface.co/chat")
		page_llama2 = stealth.MustPage(browser)
		page_llama2.MustNavigate("https://ysharma-explore-llamav2-with-tgi.hf.space")
		for i := 1; i <= 30; i++ {
			if page_llama2.MustHasX("//textarea[@data-testid='textbox']") {
				relogin_llama2 = false
				break
			}
			time.Sleep(time.Second)
		}
		if relogin_llama2 == true {
			sprint("✘ Llama2")
			//page_hc.MustPDF("./tmp/HuggingChat✘.pdf")
		}
		if relogin_llama2 == false {
			sprint("✔ Llama2")
			for {
				select {
				case question := <-channel_llama2:
					//page_hc.MustActivate()
					//fmt.Println("Falcon180 received question...", question)
					for i := 1; i <= 20; i++ {
						if page_llama2.MustHasX("//textarea[@data-testid='textbox']") {
							page_llama2.MustElementX("//textarea[@data-testid='textbox']").MustInput(question)
							break
						}
						time.Sleep(time.Second)
					}
					//fmt.Println("Falcon180 input typed...")
					for i := 1; i <= 20; i++ {

						//if page_llama2.MustHasX("//button[contains(text(), 'Submit')]") {
							//page_llama2.MustElementX("//button[contains(text(), 'Submit')]").MustClick()
						if page_llama2.MustHasX("//button[@id='component-14']") {
							page_llama2.MustElementX("//button[@id='component-14']").MustClick()
							break
						}
						time.Sleep(time.Second)
					}
					fmt.Println("Llama2 generating...")
					//page_falcon180.MustActivate() // Sometime three dot to hang
					//if role == ".all" {
					//	channel_falcon180 <- "click_falcon180"
					//}
					//// Check Error
					//channel_falcon180_check := make(chan string)
					//go func() {
					//	for i := 1; i <= 20; i++ {
					//		if page_falcon180.MustHasX("//*[contains(text(), 'Too much traffic, please try again')]") {
					//			channel_falcon180 <- "✘✘ Falcon180, Please check the internet connection and verify login status. Traffic."
					//			fmt.Println("Falcon180 too much traffic...")
					//			relogin_falcon180 = true
					//			close(channel_falcon180_check)
					//			break
					//		}
					//		time.Sleep(1 * time.Second)
					//	}
					//}()

					// stop_icon
					var stop_icon_disappear = false
					err := rod.Try(func() {
						//page_llama2.Timeout(10 * time.Second).MustElementX("//button[contains(text(), 'Stop')]").MustWaitVisible().CancelTimeout()
						page_llama2.Timeout(10 * time.Second).MustElementX("//button[@id='component-15']").MustWaitVisible().CancelTimeout()
					})
					if err == nil {
						err = rod.Try(func() {
							//page_llama2.Timeout(80 * time.Second).MustElementX("//button[contains(text(), 'Stop')]").MustWaitInvisible().CancelTimeout()
							page_llama2.Timeout(80 * time.Second).MustElementX("//button[@id='component-15']").MustWaitInvisible().CancelTimeout()
						})
						if err == nil {
							stop_icon_disappear = true
						} else {
							//fmt.Println("err::::", err)
						}
					} else {
						//fmt.Println("err::", err)
					}

					if stop_icon_disappear == true {
						answer := page_llama2.MustElementX("(//div[@data-testid='bot'])[last()]").MustText()
						channel_llama2 <- answer
					} else {
						channel_llama2 <- "✘✘  Llama2, Please check the internet connection and verify login status."
						relogin_llama2 = true
						//page_hc.MustPDF("./tmp/HuggingChat✘.pdf")

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
				//close(channel_hc)
				close(channel_llama2)
				close(channel_falcon180)
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

	// Start Loop to read user input
	for {
		// Re-read user input history
		if f, err := os.Open(".history"); err == nil {
			Liner.ReadHistory(f)
			f.Close()
		}

		prompt := strconv.Itoa(left_tokens) + role + "> "
		userInput = multiln_input(Liner, prompt)

		// Check Aih commands
		switch userInput {
		case "":
			continue
		case ".proxy":
			proxy, _ := Liner.Prompt("Please input your proxy:")
			if proxy == "" {
				continue
			}
			aihj, err := ioutil.ReadFile(".aih.json")
			new_aihj, _ := sjson.Set(string(aihj), "proxy", proxy)
			err = ioutil.WriteFile(".aih.json", []byte(new_aihj), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
			fmt.Println("Please restart Aih for using proxy")
			/// exit
			browser.MustClose()
			close(channel_bard)
			close(channel_chatgpt)
			close(channel_claude)
			//close(channel_hc)
			close(channel_llama2)
			close(channel_falcon180)
			Liner.Close()
			syscall.Exit(0)
			//os.Exit(0)
		case ".help":
			fmt.Println("                           ")
			fmt.Println("                 Welcome to Aih!                             ")
			fmt.Println("--------------------------------------------------------------------------- ")
			fmt.Println(" .               Select AI mode of Bard/ChatGPT/Claude2/Llama2/Falcon180 ")
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
			fmt.Println(" .v              Mini vi to edit quest, `:ai` send, `:q` cancel")
			fmt.Println(" .c or .clear    Clear screen")
			fmt.Println(" .h or .history  Show history")
			fmt.Println(" .key            Set key of ChatGPT API")
			fmt.Println(" .proxy          Set proxy")
			fmt.Println(" .help           Help")
			fmt.Println(" .exit           Exit")
			fmt.Println(" .speak          Voice speak context (MasOS only)")
			fmt.Println(" .quiet          Not speak")
			//fmt.Println(" .new            New conversation of ChatGPT")
			fmt.Println("--------------------------------------------------------------------------- ")
			fmt.Println("                           ")
			fmt.Println("                           ")
			continue
		case ".c", ".clear":
			clear()
			continue
		case ".v", ".vi", ".vim":
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "linux", "darwin":
				cmd = exec.Command("./.mvi")
			case "windows":
				cmd = exec.Command("./.mvi.exe")
			}
			// Enter mini vi
			ptmx, err := pty.Start(cmd)

			// Handle pty size.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGWINCH)
			go func() {
				for range ch {
					if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
						log.Printf("can't resizing pty: %s", err)
					}
				}
			}()
			ch <- syscall.SIGWINCH // Initial resize.

			// Set stdin in raw mode.
			oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				log.Println(fmt.Sprintf("can't make terminal raw mode: %s", err))
				//	g.Message(err.Error(), "main", func() {})
				return
			}

			// Copy stdin to the pty and the pty to stdout.
			go io.Copy(ptmx, os.Stdin)
			io.Copy(os.Stdout, ptmx)
			ptmx.Close()
			// Reset stdin model
			err = terminal.Restore(int(os.Stdin.Fd()), oldState)
			// Read quest
			ipt, _ := ioutil.ReadFile(".quest.txt")
			// Clean quest with "LF" the ryy's tmpt value
			if qs, err := os.OpenFile(".quest.txt", os.O_WRONLY|os.O_TRUNC, 0666); err == nil {
				qs.Write([]byte{0x0a})
				qs.Close()

			}
			// When no edie or q! Empty file have "LF"(\n)
			if ipt[0] != []byte{0x0a}[0] {
				// For claude ajax
				userInput = string(ipt)
				//userInput = strings.Replace(userInput, "\r", "\n", -1)
				//userInput = strings.Replace(userInput, "\"", "\\\"", -1)
				//userInput = strings.Replace(userInput, "\n", "\\n", -1)
				//userInput = strings.TrimSuffix(userInput, "\n")
				fmt.Println(userInput)

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
			} else {
				continue
			}
			//continue
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
			//close(channel_hc)
			close(channel_llama2)
			close(channel_falcon180)
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
					//"HuggingChat",
					"Llama2",
					"Falcon180",
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
			//case "HuggingChat":
			//	role = ".huggingchat"
			//	left_tokens = 0
			//	continue
			case "Falcon180":
				role = ".falcon180"
				left_tokens = 0
				continue
			case "Llama2":
				role = ".llama2"
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
			//--------------- ----------------------
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
			//if relogin_hc == false {
			//	channel_hc <- userInput
			//	//<-channel_hc
			//}
			if relogin_llama2 == false {
				channel_llama2 <- userInput
				//<-channel_hc
			}
			if relogin_falcon180 == false {
				channel_falcon180 <- userInput
				//<-channel_hc
			}
			//--------------- ----------------------

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
			//if relogin_hc == false {
			//	answer_hc := <-channel_hc
			//	fmt.Println(">HuggingChat Done.")
			//	RESP += "\n\n---------------- huggingchat answer ----------------\n"
			//	RESP += strings.TrimSpace(answer_hc)
			//}
			if relogin_llama2 == false {
				answer_llama2 := <-channel_llama2
				fmt.Println(">Llama2 Done.")
				RESP += "\n\n---------------- llama2 answer ----------------\n"
				RESP += strings.TrimSpace(answer_llama2)
			}
			if relogin_falcon180 == false {
				answer_falcon180 := <-channel_falcon180
				fmt.Println(">Falcon180 Done.")
				RESP += "\n\n---------------- falcon180 answer ----------------\n"
				RESP += strings.TrimSpace(answer_falcon180)
			}
			speak_out(speak, RESP)
			save_conversation(role, userInput, RESP)
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
				speak_out(speak, RESP)
				save_conversation(role, userInput, RESP)
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
				speak_out(speak, RESP)
				save_conversation(role, userInput, RESP)
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
				speak_out(speak, RESP)
				save_conversation(role, userInput, RESP)
				printer(color_chatapi, RESP, false)
			}

		}

		//// HUGGINGCHAT:
		//if role == ".huggingchat" {
		//	if relogin_hc == true {
		//		fmt.Println("✘ HuggingChat")
		//	} else {
		//		channel_hc <- userInput
		//		answer := <-channel_hc

		//		// Print the response to the terminal
		//		RESP = strings.TrimSpace(answer)
		//		speak_out(speak, RESP)
		//		save_conversation(role, userInput, RESP)
		//		printer(color_huggingchat, RESP, false)
		//	}

		//}

		// FALCON:
		if role == ".falcon180" {
			if relogin_falcon180 == true {
				fmt.Println("✘ Falcon180")
			} else {
				channel_falcon180 <- userInput
				answer := <-channel_falcon180

				// Print the response to the terminal
				RESP = strings.TrimSpace(answer)
				speak_out(speak, RESP)
				save_conversation(role, userInput, RESP)
				printer(color_chat, RESP, false)
			}

		}
		// LLAMA2:
		if role == ".llama2" {
			if relogin_llama2 == true {
				fmt.Println("✘ Llama2")
			} else {
				channel_llama2 <- userInput
				answer := <-channel_llama2

				// Print the response to the terminal
				RESP = strings.TrimSpace(answer)
				speak_out(speak, RESP)
				save_conversation(role, userInput, RESP)
				printer(color_chat, RESP, false)
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
				aihj, err := ioutil.ReadFile(".aih.json")
				new_aihj, _ := sjson.Set(string(aihj), "key", OpenAI_Key)
				err = ioutil.WriteFile(".aih.json", []byte(new_aihj), 0644)
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
			speak_out(speak, RESP)
			save_conversation(role, userInput, RESP)
			printer(color_chatapi, RESP, false)

		}

		save2clip_board(RESP)
		// clean RESP
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

	fmt.Fprint(textView, tview.Escape(context))
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

// Persistent conversation uInput + response
func save_conversation(role, uInput, RESP string) {
	if fs, err := os.OpenFile("history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
		time_string := time.Now().Format("2006-01-02 15:04:05")
		_, err = fs.WriteString("--------------------\n")
		_, err = fs.WriteString(time_string + role + "\nQuestion:\n" + uInput + "\n")
		_, err = fs.WriteString("Answer:" + "\n" + RESP + "\n")
		if err != nil {
			panic(err)
		}
		fs.Close()
	}
}

// Speak all the response RESP using the "say" command
func speak_out(speak int, RESP string) {
	if speak == 1 {
		//fmt.Println("speaking")
		go func() {
			switch runtime.GOOS {
			case "linux", "darwin":
				cmd := exec.Command("say", RESP)
				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
			case "windows":
				// to do
				_ = 1 + 1

			}

		}()
	}
}
