package main

import (
	_ "embed"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"strings"
	"time"
)

// Set up client of chatgpt (rod version)
var page_chatgpt *rod.Page
var relogin_chatgpt = true
var channel_chatgpt chan string

func Chatgpt() {
	channel_chatgpt = make(chan string)
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
				page_chatgpt.MustActivate()
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
				fmt.Println("ChatGPT generating...", role)
				//page_chatgpt.MustActivate()
				if role == ".all" {
					channel_chatgpt <- "click_chatgpt"
				}


			        //page_chatgpt.MustElementX("//div[contains(text(), 'Stop generating')]")
				//page_chatgpt.MustElement("svg path[d='M0 2a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V2z']") 
				//fmt.Println("found creating icon")
				var regenerate_icon = false

				for i := 1; i <= 60; i++ {
					//if page_chatgpt.MustHas("svg path[d='M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15']") {
					if page_chatgpt.MustHas("svg path[d='M0 2a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V2z']") {

					 time.Sleep(1 * time.Second)
					 continue
					}

					if page_chatgpt.MustHas("svg path[d='M4.5 2.5C5.05228 2.5 5.5 2.94772 5.5 3.5V5.07196C7.19872 3.47759 9.48483 2.5 12 2.5C17.2467 2.5 21.5 6.75329 21.5 12C21.5 17.2467 17.2467 21.5 12 21.5C7.1307 21.5 3.11828 17.8375 2.565 13.1164C2.50071 12.5679 2.89327 12.0711 3.4418 12.0068C3.99033 11.9425 4.48712 12.3351 4.5514 12.8836C4.98798 16.6089 8.15708 19.5 12 19.5C16.1421 19.5 19.5 16.1421 19.5 12C19.5 7.85786 16.1421 4.5 12 4.5C9.7796 4.5 7.7836 5.46469 6.40954 7H9C9.55228 7 10 7.44772 10 8C10 8.55228 9.55228 9 9 9H4.5C3.96064 9 3.52101 8.57299 3.50073 8.03859C3.49983 8.01771 3.49958 7.99677 3.5 7.9758V3.5C3.5 2.94772 3.94771 2.5 4.5 2.5Z']") {
					   //page_chatgpt.MustElement("svg path[d='M7 11L12 6L17 11M12 18V7']").MustWaitVisible().MustWaitStable() 
					if page_chatgpt.MustHas("svg path[d='M7 11L12 6L17 11M12 18V7']") {
						regenerate_icon = true
						//fmt.Println("wait...")
						break
					       }
					}
					time.Sleep(1 * time.Second)
				}
				if regenerate_icon == true { 
				        page_chatgpt.MustActivate()
					time.Sleep(3 * time.Second)

					answer := page_chatgpt.MustElementX("(//div[contains(@class, 'w-full text-token-text-primary')])[last()]").MustText()[15:]
					//answer := page_chatgpt.MustElementX("(//div[contains(@class, 'group w-full')])[last()]").MustText()[7:]
					//answer := page_chatgpt.MustElementX("(//div[contains(@class, 'group final-completion w-full')])[last()]").MustText()[7:]
					//answer_div := page_chatgpt.MustElementX("(//div[contains(@class, 'w-full text-token-text-primary')])[last()]").MustWaitStable()
					//answer := answer_div.MustText()[15:]

					//response := img.MustElementX("ancestor::bard-avatar[1]/parent::div/parent::div").MustWaitVisible()
					//answer := response.MustText()

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
}
