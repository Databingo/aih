package main

import (
	_ "embed"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"time"
)

// Set up client of Llama2 (Rod version)
var page_llama2 *rod.Page
var relogin_llama2 = true
var channel_llama2 chan string

func Llama2() {
	channel_llama2 = make(chan string)
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

}
