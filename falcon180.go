package main

import (
	_ "embed"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"time"
)

// Set up client of falcon180-180B (Rod version)
var page_falcon180 *rod.Page
var relogin_falcon180 = true
var channel_falcon180 chan string

func Falcon180() {
	channel_falcon180 = make(chan string)
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

}
