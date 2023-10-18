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
					//answer := page_chatgpt.MustElementX("(//div[contains(@class, 'group w-full')])[last()]").MustText()[7:]
					answer := page_chatgpt.MustElementX("(//div[contains(@class, 'group final-completion w-full')])[last()]").MustText()[7:]
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
