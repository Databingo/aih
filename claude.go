package main

import (
	//	"context"
	_ "embed"
	"fmt"
	//	"github.com/atotto/clipboard"
	//	"github.com/creack/pty"
	//	"github.com/gdamore/tcell/v2"
	"github.com/go-rod/rod"
	//	"github.com/go-rod/rod/lib/launcher"
	//	"github.com/go-rod/rod/lib/utils"
	"github.com/go-rod/stealth"
	"github.com/google/uuid"
	//	"github.com/manifoldco/promptui"
	//	"github.com/peterh/liner"
	//	"github.com/rivo/tview"
	//	openai "github.com/sashabaranov/go-openai"
	"github.com/tidwall/gjson"
	//	"github.com/tidwall/sjson"
	//	"golang.org/x/crypto/ssh/terminal"
	//	"io"
	//	"io/ioutil"
	//	"log"
	//	"os"
	//	"os/exec"
	//	"os/signal"
	//	"runtime"
	"strconv"
	"strings"
	//	"syscall"
	"time"
)

var page_claude *rod.Page
var relogin_claude = true
var channel_claude chan string

//	channel_claude := make(chan string)
func Claude2() {
	channel_claude = make(chan string)
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
		sprint("✘ Claude")
		//page_claude.MustPDF("./tmp/Claude✘.pdf")
	}
	if relogin_claude == false {
		sprint("✔ Claude")
		for {
			select {
			case question := <-channel_claude:
				question = strings.Replace(question, "\r", "\n", -1)
				question = strings.Replace(question, "\"", "\\\"", -1)
				question = strings.Replace(question, "\n", "\\n", -1)
				question = strings.TrimSuffix(question, "\n")
				// re-activate
				page_claude.MustNavigate("https://claude.ai/api/account/statsig/" + org_uuid).MustWaitLoad()
				record_json := page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()
				record_chat_messages = gjson.Get(string(record_json), "chat_messages").String()
				record_count := gjson.Get(string(response_chat_messages), "#").Int()
				page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid).MustWaitLoad()
				time.Sleep(2 * time.Second) // delay to simulate human being
				//question = strings.Replace(question, `"`, `\"`, -1) // escape " in input text when code into json

				d := `{"completion":{"prompt":"` + question + `","timezone":"Asia/Shanghai","model":"claude-2"},"organization_uuid":"` + org_uuid + `","conversation_uuid":"` + new_uuid + `","text":"` + question + `","attachments":[]}`
				//fmt.Println(d)
				js := `
		                               (sdata, new_uuid) => {
		                               var xhr = new XMLHttpRequest();
		                               xhr.open("POST", "https://claude.ai/api/append_message");
		                               xhr.setRequestHeader('Content-Type', 'application/json');
		                               xhr.setRequestHeader('Referer', 'https://claude.ai/chat/new_uuid');
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
				page_claude.MustEval(js, d, new_uuid).Str()
				fmt.Println("Claude generating...")
				//if role == ".all" {
				//	channel_claude <- "click_claude"
				//}
				time.Sleep(3 * time.Second) // delay to simulate human being

				// wait answer
				var claude_response = false
				var response_json string
				for i := 1; i <= 20; i++ {
					if page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustHasX("//pre") {
						response_json = page_claude.MustNavigate("https://claude.ai/api/organizations/" + org_uuid + "/chat_conversations/" + new_uuid).MustElementX("//pre").MustText()
						response_chat_messages = gjson.Get(string(response_json), "chat_messages").String()
						count := gjson.Get(string(response_chat_messages), "#").Int()

						if response_chat_messages != record_chat_messages && count == record_count+2 {
							claude_response = true
							record_chat_messages = response_chat_messages
							break
						}
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
					//page_claude.MustPDF("./tmp/Claude✘.pdf")
				}
			}
		}
	}

}
