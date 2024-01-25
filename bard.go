package main

import (
	_ "embed"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"strings"
	"time"
)


// Set up client of Bard (Rod version)
var page_bard *rod.Page
var relogin_bard = true
var channel_bard chan string

func Bard() {
	channel_bard = make(chan string)
	defer func() {
		if err := recover(); err != nil {
			relogin_bard = true
		}
	}()
	page_bard = stealth.MustPage(browser)
	page_bard.MustNavigate("https://bard.google.com")

	for i := 1; i <= 30; i++ {
		//if page_bard.MustHasX("//textarea[@id='mat-input-0']") {
		//if page_bard.MustHasX("//rich-textarea[@aria-label='Input for prompt text']") {
		if page_bard.MustHasX("//rich-textarea[@enterkeyhint='send']") {
			relogin_bard = false
			break
		}
		// Check "I'm not a robot"
		info := page_bard.MustInfo()
		if strings.HasPrefix(info.URL, "https://google.com/sorry") {
			relogin_bard = true
			break
		}
		// Check "Sign in"
		if page_bard.MustHasX("//a[contains(text(), 'Sign in')]") {
			relogin_bard = true
			break
		}
		// Check "You've been signed out"
		if page_bard.MustHasX("//h1[contains(text(), 've been signed out')]") {
			relogin_bard = true
			break
		}

		time.Sleep(time.Second)
	}
	if relogin_bard == true {
		sprint("✘ Bard")
	}
	if relogin_bard == false {
		sprint("✔ Bard")
		for {
			select {
			case question := <-channel_bard:
				page_bard.MustActivate()       
				//page_bard.MustElementX("//rich-textarea[@aria-label='Input for prompt text']").MustWaitVisible().MustInput(question)
				page_bard.MustElementX("//rich-textarea[@enterkeyhint='send']").MustWaitVisible().MustInput(question)
				page_bard.MustElementX("//button[@mattooltip='Submit']").MustClick()
				fmt.Println("Bard generating...", role)
				//page_bard.MustActivate()
				if role == ".all" {
				        //fmt.Println("Bard role", role)
					channel_bard <- "click_bard"
				}
				// wait generated icon
				var generated_icon_appear = false
				var c = 0
				for i := 1; i <= 60; i++ {
					if page_bard.MustHasX("//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')]") {
						generated_icon_appear = true
						break
					}
					c = c + 1 
			                if c == 5 { 
					     page_bard.MustActivate()
					     c = 0
					    }
					time.Sleep(1 * time.Second)
				}
				if generated_icon_appear == true {
					img := page_bard.MustElementX("//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')][last()]").MustWaitVisible()
					//response := img.MustElementX("parent::div/parent::bard-logo/parent::div/parent::div").MustWaitVisible()
					//response := img.MustElementX("parent::div/parent::div/parent::bard-avatar/parent::div/parent::div").MustWaitVisible()
					response := img.MustElementX("ancestor::bard-avatar[1]/parent::div/parent::div").MustWaitVisible()
					answer := response.MustText()
					channel_bard <- answer
				} else {
					channel_bard <- "✘✘  Bard, Please check the internet connection and verify login status."
					relogin_bard = true

				}
			}
		}
	}

}
