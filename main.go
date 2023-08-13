package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/go-rod/rod/lib/utils"
	"github.com/go-rod/rod/lib/launcher"

	"bufio"
	//"bytes"
	"context"
	"encoding/json"
	"fmt"
	//"github.com/Databingo/EdgeGPT-Go"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	//"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/peterh/liner"
	"github.com/rivo/tview"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"io/ioutil"
	//"net/http"
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
	
	// Set proxy for rod
	//proxy_url := launcher.New().Proxy(Proxy).Delete("use-mock-keychain").MustLaunch()
	proxy_url := launcher.New().Proxy(Proxy).MustLaunch()

        // open rod browser
        var browser *rod.Browser
	if Proxy != "" {
	browser = rod.New().
	           Trace(true).
		   ControlURL(proxy_url).
	           Timeout(3 * time.Minute).
		   MustConnect()
		  } else {
	browser = rod.New().
	           Trace(true).
	           Timeout(3 * time.Minute).
		   MustConnect()
		  }

	// Test Proxy
	//TEST_PROXY:
	//	fmt.Println("Checking network accessing...")
	//	ops1 := googlesearch.SearchOptions{Limit: 12}
	//	_, err = googlesearch.Search(nil, "BTC", ops1)
	//	if err != nil {
	//		fmt.Println("Need proxy to access GoogleBard, BingChat, ChatGPT")
	//		proxy, _ := Liner.Prompt("Please input proxy: ")
	//		if proxy == "" {
	//			goto TEST_PROXY
	//		}
	//		aihj, err := ioutil.ReadFile("aih.json")
	//		new_aihj, _ := sjson.Set(string(aihj), "proxy", proxy)
	//		err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
	//		if err != nil {
	//			fmt.Println("Save failed.")
	//		}
	//		fmt.Println("Please restart Aih for using proxy...")
	//		Liner.Close()
	//		syscall.Exit(0)
	//
	//	}

	// Set up client of OpenAI API
	key := gjson.Get(string(aih_json), "key")
	OpenAI_Key := key.String()
	config := openai.DefaultConfig(OpenAI_Key)
	client := openai.NewClientWithConfig(config)
	messages := make([]openai.ChatCompletionMessage, 0)

	//////////////////////0////////////////////////////
	// Set up client of ChatGPT (chromedriver version)
	pf_chatgpt, _ := os.CreateTemp("", "pf_chatgpt.py")
	_, _ = pf_chatgpt.WriteString(ps_chatgpt)
	_ = pf_chatgpt.Close()
	defer os.Remove(pf_chatgpt.Name())

	// Read cookie
	chatgpt_json, err := ioutil.ReadFile("./4.json")
	if err != nil {
		err = ioutil.WriteFile("./4.json", []byte(""), 0644)
	}
	var chatgptjs string
	chatgptjs = gjson.Parse(string(chatgpt_json)).String()
	var cmd_chatgpt *exec.Cmd
	var stdout_chatgpt io.ReadCloser
	var stdin_chatgpt io.WriteCloser
	var login_chatgpt bool
	var relogin_chatgpt bool
	var scanner_chatgpt *bufio.Scanner
	channel_chatgpt_answer := make(chan string)
	if chatgptjs != "" {
		//cmd_bard = exec.Command("python3", "-u", "./bard.py", "load")
		cmd_chatgpt = exec.Command("python3", "-u", pf_chatgpt.Name(), "load")
		stdout_chatgpt, _ = cmd_chatgpt.StdoutPipe()
		stdin_chatgpt, _ = cmd_chatgpt.StdinPipe()

		go func(cmd *exec.Cmd) {
			time.Sleep(3 * 0 * time.Second)
			if err := cmd.Start(); err != nil {
				panic(err)
			}
		}(cmd_chatgpt)

		login_chatgpt = false
		relogin_chatgpt = false
		go func(login_chatgpt, relogin_chatgpt *bool) {
			time.Sleep(3 * 0 * time.Second)
			scanner_chatgpt = bufio.NewScanner(stdout_chatgpt)
			for scanner_chatgpt.Scan() {
				RESP = scanner_chatgpt.Text()
				if RESP == "login work" {
					*login_chatgpt = true
				} else if RESP == "relogin" {
					*relogin_chatgpt = true
				} else {
					channel_chatgpt_answer <- RESP
				}
			}
		}(&login_chatgpt, &relogin_chatgpt)
	}

	//////////////////////1////////////////////////////
	// Set up client of Bard (chromedriver version)
	pf_bard, _ := os.CreateTemp("", "pf_bard.py")
	_, _ = pf_bard.WriteString(ps_bard)
	_ = pf_bard.Close()
	defer os.Remove(pf_bard.Name())

	// Read cookie
	bard_json, err := ioutil.ReadFile("./2.json")
	if err != nil {
		err = ioutil.WriteFile("./2.json", []byte(""), 0644)
	}
	var bjs string
	bjs = gjson.Parse(string(bard_json)).String()
	var cmd_bard *exec.Cmd
	var stdout_bard io.ReadCloser
	var stdin_bard io.WriteCloser
	var login_bard bool
	var relogin_bard bool
	var scanner_bard *bufio.Scanner
	channel_bard_answer := make(chan string)
	if bjs != "" {
		//cmd_bard = exec.Command("python3", "-u", "./bard.py", "load")
		cmd_bard = exec.Command("python3", "-u", pf_bard.Name(), "load")
		stdout_bard, _ = cmd_bard.StdoutPipe()
		stdin_bard, _ = cmd_bard.StdinPipe()

		go func(cmd *exec.Cmd) {
			time.Sleep(3 * 1 * time.Second)
			if err := cmd.Start(); err != nil {
				panic(err)
			}
		}(cmd_bard)

		login_bard = false
		relogin_bard = false
		go func(login_bard, relogin_bard *bool) {
			time.Sleep(3 * 1 * time.Second)
			scanner_bard = bufio.NewScanner(stdout_bard)
			for scanner_bard.Scan() {
				RESP = scanner_bard.Text()
				if RESP == "login work" {
					*login_bard = true
				} else if RESP == "relogin" {
					*relogin_bard = true
				} else {
					channel_bard_answer <- RESP
				}
			}
		}(&login_bard, &relogin_bard)
	}

	//////////////////////2////////////////////////////
	// Set up client of Claude2 (chromedriver version)
	pf_claude, _ := os.CreateTemp("", "pf_claude.py")
	_, _ = pf_claude.WriteString(ps_claude)
	_ = pf_claude.Close()
	defer os.Remove(pf_claude.Name())

	// Read cookie
	claude2_json, err := ioutil.ReadFile("./3.json")
	if err != nil {
		err = ioutil.WriteFile("./3.json", []byte(""), 0644)
	}
	var c2js string
	c2js = gjson.Parse(string(claude2_json)).String()
	var cmd_claude2 *exec.Cmd
	var stdout_claude2 io.ReadCloser
	var stdin_claude2 io.WriteCloser
	var login_claude2 bool
	var relogin_claude2 bool
	var scanner_claude2 *bufio.Scanner
	channel_claude2_answer := make(chan string)
	if c2js != "" {
		//cmd_claude2 = exec.Command("python3", "-u", "./claude2.py", "load")
		cmd_claude2 = exec.Command("python3", "-u", pf_claude.Name(), "load")
		stdout_claude2, _ = cmd_claude2.StdoutPipe()
		stdin_claude2, _ = cmd_claude2.StdinPipe()

		go func(cmd *exec.Cmd) {
			time.Sleep(3 * 2 * time.Second)
			if err := cmd.Start(); err != nil {
				panic(err)
			}
		}(cmd_claude2)

		login_claude2 = false
		relogin_claude2 = false
		go func(login_claude2, relogin_claude2 *bool) {
			time.Sleep(3 * 2 * time.Second)
			scanner_claude2 = bufio.NewScanner(stdout_claude2)
			for scanner_claude2.Scan() {
				RESP = scanner_claude2.Text()
				if RESP == "login work" {
					*login_claude2 = true
				} else if RESP == "relogin" {
					*relogin_claude2 = true
				} else {
					channel_claude2_answer <- RESP
				}
			}
		}(&login_claude2, &relogin_claude2)
	}

	//////////////////////3////////////////////////////
	// Set up client of huggingchat (chromedriver version)
	pf_hc, _ := os.CreateTemp("", "pf_hc.py")
	_, _ = pf_hc.WriteString(ps_huggingchat)
	_ = pf_hc.Close()
	defer os.Remove(pf_hc.Name())

	// Read cookie
	hc_json, err := ioutil.ReadFile("./5.json")
	if err != nil {
		err = ioutil.WriteFile("./5.json", []byte(""), 0644)
	}
	var hcjs string
	hcjs = gjson.Parse(string(hc_json)).String()
	var cmd_hc *exec.Cmd
	var stdout_hc io.ReadCloser
	var stdin_hc io.WriteCloser
	var login_hc bool
	var relogin_hc bool
	var scanner_hc *bufio.Scanner
	channel_hc_answer := make(chan string)
	if hcjs != "" {
		//cmd_hc = exec.Command("python3", "-u", "./huggingchat.py", "load")
		cmd_hc = exec.Command("python3", "-u", pf_hc.Name(), "load")
		stdout_hc, _ = cmd_hc.StdoutPipe()
		stdin_hc, _ = cmd_hc.StdinPipe()

		//time.Sleep(time.Duration(3*3) * time.Second)
		go func(cmd *exec.Cmd) {
			time.Sleep(3 * 3 * time.Second)
			if err := cmd.Start(); err != nil {
				panic(err)
			}
		}(cmd_hc)

		login_hc = false
		relogin_hc = false
		go func(login_hc, relogin_hc *bool) {
			time.Sleep(3 * 3 * time.Second)
			scanner_hc = bufio.NewScanner(stdout_hc)
			for scanner_hc.Scan() {
				RESP = scanner_hc.Text()
				if RESP == "login work" {
					*login_hc = true
				} else if RESP == "relogin" {
					*relogin_hc = true
				} else {
					channel_hc_answer <- RESP
				}
			}
		}(&login_hc, &relogin_hc)
	}

	//////////////////////4////////////////////////////
	// Set up client of Rod_chatgpt (rod version)

	// Read user/password
	u_json, _ := ioutil.ReadFile("user.json")
	if err != nil {
		err = ioutil.WriteFile("user.json", []byte(""), 0644)
	}
	var chatgpt_user string
	var chatgpt_password string
	chatgpt_user = gjson.Get(string(u_json), "chatgpt.user").String()
	chatgpt_password = gjson.Get(string(u_json), "chatgpt.password").String()
	if chatgpt_user == "Name" { chatgpt_user = "" }
	if chatgpt_password == "Password" { chatgpt_password = "" }
	fmt.Println(chatgpt_user)
	fmt.Println(chatgpt_password)

	//// Read cookie
	//chatgpt_json, err := ioutil.ReadFile("cookies/chatgpt.json")
	//if err != nil {
	//	err = ioutil.WriteFile("cookies/chatgpt.json", []byte(""), 0644)
	//}
	//var chatgptjs string
	//chatgptjs = gjson.Parse(string(chatgpt_json)).String()

        //// Open page with cookie 
	//if chatgptjs != "" {
	//    //cookie
	//    page := stealth.MustPage(browser)
	//    page.MustNavigate("https://chat.openai.com")
	//}

        // Open page with password
        var page_chatgpt  *rod.Page
	if (chatgpt_user != "" && chatgpt_password != "") {
	    //cookie
	    page_chatgpt = stealth.MustPage(browser)
	    page_chatgpt.MustNavigate("https://chat.openai.com")
	    utils.Pause()
	}



	// Set up client of Bing Chat
	//var gpt *EdgeGPT.GPT
	//_, err = ioutil.ReadFile("./cookies/1.json")
	//if err == nil {
	//	s := EdgeGPT.NewStorage()
	//	ch := make(chan bool)
	//	go func() {
	//		defer func() {
	//			if r := recover(); r != nil {
	//				_ = os.Remove("./cookies/1.json")
	//				ch <- true
	//				return
	//			}
	//		}()
	//		gpt, err = s.GetOrSet("any-key")
	//		ch <- true
	//	}()
	//	<-ch
	//}

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
	role := ".bard"
	uInput := ""
	//price := ""
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
			Liner.Close()
			syscall.Exit(0)
		case ".help":
			fmt.Println("                           ")
			fmt.Println(" .               Select AI mode of Bard/Bing/ChatGPT/Claude")
			fmt.Println(" .key            Set cookie of Bard/Bing/ChatGPT/Claude")
			fmt.Println(" .proxy          Set proxy")
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
			fmt.Println(" .new            New conversation of ChatGPT")
			fmt.Println(" .speak          Voice speak context (MasOS only)")
			fmt.Println(" .quiet          Not speak")
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
			cmd_bard.Process.Kill()
			cmd_claude2.Process.Kill()
			cmd_chatgpt.Process.Kill()
			cmd_hc.Process.Kill()
			switch runtime.GOOS {
			case "linux", "darwin":
				cmd := exec.Command("pkill", "-f", "undetected_chromedriver")
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
			case "windows":
				cmd := exec.Command("taskkill", "/IM", "undetected_chromedriver", "/F")
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
			}

			return
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
					"Bard",
					//"Bing",
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
					"Set Bard Cookie",
					"Set ChatGPT Cookie",
					"Set Claude Cookie",
					"Set HuggingChat Cookie",
					"Set ChatGPT API Key",
					"Exit",
				},
			}

			_, keyy, err := prom.Run()
			if err != nil {
				panic(err)
			}

			switch keyy {
			case "Set Bard Cookie":
				if bjs != "" {
					cmd_bard.Process.Kill()
				}
				bjs = ""
				role = ".bard"
				goto BARD
			case "Set ChatGPT Cookie":
				if chatgptjs != "" {
					cmd_chatgpt.Process.Kill()
				}
				chatgptjs = ""
				role = ".chat"
				goto CHAT
			case "Set ChatGPT API Key":
				OpenAI_Key = ""
				role = ".chatapi"
				goto CHATAPI
			//case "Set Bing Chat Cookie":
			//	_ = os.Remove("./cookies/1.json")
			//	role = ".bing"
			//	goto BING
			case "Set Claude Cookie":
				if c2js != "" {
					cmd_claude2.Process.Kill()
				}
				c2js = ""
				role = ".claude"
				goto CLAUDE
			case "Set HuggingChat Cookie":
				if hcjs != "" {
					cmd_hc.Process.Kill()
				}
				hcjs = ""
				role = ".huggingchat"
				goto HUGGINGCHAT

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

	BARD:
		// Check role for correct actions
		if role == ".bard" {

			if bjs == "" {
				prom := "Please type << then paste Bard cookie then type >> then press Enter: "
				cook := multiln_input(Liner, prom)

				// Clear screen of input cookie string
				clear()

				// Check cookie
				cook = strings.Replace(cook, "\r", "", -1)
				cook = strings.Replace(cook, "\n", "", -1)
				if len(cook) < 100 {
					fmt.Println("Invalid cookie")
					continue
				}
				if !json.Valid([]byte(cook)) {
					fmt.Println("Invalid JSON format")
					continue
				}
				if !strings.Contains(cook, ".google.com") {
					fmt.Println("Invalid cookie, please make sure the tab is bard.google.com")
					continue

				}

				// Save cookie
				err = ioutil.WriteFile("./2.json", []byte(cook), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}

				// Reload bard cookie
				bard_json, err = ioutil.ReadFile("./2.json")
				bjs = gjson.Parse(string(bard_json)).String()
				if bjs == "" {
					continue
				}
				if bjs != "" {
					//cmd_bard = exec.Command("python3", "-u", "./bard.py", "load")
					cmd_bard = exec.Command("python3", "-u", pf_bard.Name(), "load")
					stdout_bard, _ = cmd_bard.StdoutPipe()
					stdin_bard, _ = cmd_bard.StdinPipe()
					go func(cmd *exec.Cmd) {
						if err := cmd.Start(); err != nil {
							panic(err)
						}
					}(cmd_bard)

					scanner_bard = bufio.NewScanner(stdout_bard)
					login_bard = false
					relogin_bard = false
					go func(login_bard, relogin_bard *bool) {
						for scanner_bard.Scan() {
							RESP = scanner_bard.Text()
							if RESP == "login work" {
								*login_bard = true
							} else if RESP == "relogin" {
								*relogin_bard = true
							} else {
								channel_bard_answer <- RESP
							}
						}
					}(&login_bard, &relogin_bard)
				}
			}
			if relogin_bard == true {
				fmt.Println("Cookie failed, please renew bard cookie...")
				bjs = ""
				continue

			}
			if login_bard != true {
				fmt.Println("Bard initializing...")
				continue
			}

			spc := strings.Replace(userInput, "\n", "(-:]", -1)
			_, err = io.WriteString(stdin_bard, spc+"\n")
			if err != nil {
				panic(err)
			}

			RESP = <-channel_bard_answer
			RESP = strings.Replace(RESP, "(-:]", "\n", -1)
			printer(color_bard, RESP, false)
			save2clip_board(RESP)

		}

	CLAUDE:
		// Check role for correct actions
		if role == ".claude" {

			if c2js == "" {
				prom := "Please type << then paste Claude2 cookie then type >> then press Enter: "
				cook := multiln_input(Liner, prom)

				// Clear screen of input cookie string
				clear()

				// Check cookie
				cook = strings.Replace(cook, "\r", "", -1)
				cook = strings.Replace(cook, "\n", "", -1)
				if len(cook) < 100 {
					fmt.Println("Invalid cookie")
					continue
				}
				if !json.Valid([]byte(cook)) {
					fmt.Println("Invalid JSON format")
					continue
				}
				if !strings.Contains(cook, ".claude") {
					fmt.Println("Invalid cookie, please make sure the tab is claude.ai")
					continue

				}

				// Save cookie
				err = ioutil.WriteFile("./3.json", []byte(cook), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}

				// Reload claude2 cookie
				claude2_json, err = ioutil.ReadFile("./3.json")
				c2js = gjson.Parse(string(claude2_json)).String()
				if c2js == "" {
					continue
				}
				if c2js != "" {
					//cmd_claude2 = exec.Command("python3", "-u", "./claude2.py", "load")
					cmd_claude2 = exec.Command("python3", "-u", pf_claude.Name(), "load")
					stdout_claude2, _ = cmd_claude2.StdoutPipe()
					stdin_claude2, _ = cmd_claude2.StdinPipe()

					go func(cmd *exec.Cmd) {
						if err := cmd.Start(); err != nil {
							panic(err)
						}
					}(cmd_claude2)
					scanner_claude2 = bufio.NewScanner(stdout_claude2)
					login_claude2 = false
					relogin_claude2 = false
					go func(login_claude2, relogin_claude2 *bool) {
						for scanner_claude2.Scan() {
							RESP = scanner_claude2.Text()
							if RESP == "login work" {
								*login_claude2 = true
							} else if RESP == "relogin" {
								*relogin_claude2 = true
							} else {
								channel_claude2_answer <- RESP
							}
						}
					}(&login_claude2, &relogin_claude2)
				}
			}
			if relogin_claude2 == true {
				fmt.Println("Cookie failed, please renew claude2 cookie...")
				c2js = ""
				continue

			}
			if login_claude2 != true {
				fmt.Println("Claude2 initializing...")
				continue
			}

			spc := strings.Replace(userInput, "\n", "(-:]", -1)
			_, err = io.WriteString(stdin_claude2, spc+"\n")
			if err != nil {
				panic(err)
			}

			RESP = <-channel_claude2_answer
			RESP = strings.Replace(RESP, "(-:]", "\n", -1)
			printer(color_claude, RESP, false)
			save2clip_board(RESP)

		}
		//	BING:
		//		if role == ".bing" {
		//			// Check BingChat cookie
		//			_, err := ioutil.ReadFile("./cookies/1.json")
		//			if err != nil {
		//				prom := "Please type << then paste Bing cookie then type >> then press Enter: "
		//				cook := multiln_input(Liner, prom)
		//
		//				// Clear screen of input cookie string
		//				clear()
		//
		//				// Check cookie
		//				cook = strings.Replace(cook, "\r", "", -1)
		//				cook = strings.Replace(cook, "\n", "", -1)
		//				if len(cook) < 100 {
		//					fmt.Println("Invalid cookie")
		//					continue
		//				}
		//				if !json.Valid([]byte(cook)) {
		//					fmt.Println("Invalid JSON format")
		//					continue
		//				}
		//				if !strings.Contains(cook, ".bing.com") {
		//					fmt.Println("Invalid cookie, please make sure the tab is bing.com")
		//					continue
		//
		//				}
		//
		//				// Save cookie
		//				_ = os.MkdirAll("./cookies", 0755)
		//				err = ioutil.WriteFile("./cookies/1.json", []byte(cook), 0644)
		//				if err != nil {
		//					fmt.Println("Save failed.")
		//				}
		//
		//				// Renew BingChat client with cookie
		//				s := EdgeGPT.NewStorage()
		//				// Test gpt with cookie in gorountine
		//				ch := make(chan bool)
		//				go func() {
		//					// If invalid, remove cookie
		//					defer func() {
		//						if r := recover(); r != nil {
		//							_ = os.Remove("./cookies/1.json")
		//							fmt.Println("Invalid cookie value")
		//							ch <- true
		//							return
		//						}
		//					}()
		//					gpt, err = s.GetOrSet("any-key")
		//					ch <- true
		//				}()
		//				<-ch
		//				continue
		//			}
		//
		//			// Send message
		//			as, err := gpt.AskSync("creative", userInput)
		//			if err != nil {
		//				fmt.Println(err)
		//				continue
		//			}
		//			RESP = strings.TrimSpace(as.Answer.GetAnswer())
		//			save2clip_board(RESP)
		//			printer(color_bing, RESP, false)
		//		}
		//
	CHAT:
		if role == ".chat" {
			if chatgptjs == "" {
				prom := "Please type << then paste ChatGPT cookie then type >> then press Enter: "
				cook := multiln_input(Liner, prom)

				// Clear screen of input cookie string
				clear()

				// Check cookie
				cook = strings.Replace(cook, "\r", "", -1)
				cook = strings.Replace(cook, "\n", "", -1)
				if len(cook) < 100 {
					fmt.Println("Invalid cookie")
					continue
				}
				if !json.Valid([]byte(cook)) {
					fmt.Println("Invalid JSON format")
					continue
				}
				if !strings.Contains(cook, "chat.openai.com") {
					fmt.Println("Invalid cookie, please make sure the tab is chat.openai.com")
					continue

				}

				// Save cookie
				err = ioutil.WriteFile("./4.json", []byte(cook), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}

				// Reload claude2 cookie
				chatgpt_json, err = ioutil.ReadFile("./4.json")
				chatgptjs = gjson.Parse(string(chatgpt_json)).String()
				if chatgptjs == "" {
					continue
				}
				if chatgptjs != "" {
					//cmd_bard = exec.Command("python3", "-u", "./bard.py", "load")
					cmd_chatgpt = exec.Command("python3", "-u", pf_chatgpt.Name(), "load")
					stdout_chatgpt, _ = cmd_chatgpt.StdoutPipe()
					stdin_chatgpt, _ = cmd_chatgpt.StdinPipe()

					go func(cmd *exec.Cmd) {
						if err := cmd.Start(); err != nil {
							panic(err)
						}
					}(cmd_chatgpt)

					login_chatgpt = false
					relogin_chatgpt = false
					go func(login_chatgpt, relogin_chatgpt *bool) {
						scanner_chatgpt = bufio.NewScanner(stdout_chatgpt)
						for scanner_chatgpt.Scan() {
							RESP = scanner_chatgpt.Text()
							if RESP == "login work" {
								*login_chatgpt = true
							} else if RESP == "relogin" {
								*relogin_chatgpt = true
							} else {
								channel_chatgpt_answer <- RESP
							}
						}
					}(&login_chatgpt, &relogin_chatgpt)
				}
			}
			if relogin_chatgpt == true {
				fmt.Println("Cookie failed, please renew chatgpt cookie...")
				chatgptjs = ""
				continue

			}
			if login_chatgpt != true {
				fmt.Println("chatgpt initializing...")
				continue
			}

			spc := strings.Replace(userInput, "\n", "(-:]", -1)
			_, err = io.WriteString(stdin_chatgpt, spc+"\n")
			if err != nil {
				panic(err)
			}

			RESP = <-channel_chatgpt_answer
			RESP = strings.Replace(RESP, "(-:]", "\n", -1)
			printer(color_chat, RESP, false)
			save2clip_board(RESP)

		}

	HUGGINGCHAT:
		if role == ".huggingchat" {
			if hcjs == "" {
				prom := "Please type << then paste HuggingChat cookie then type >> then press Enter: "
				cook := multiln_input(Liner, prom)

				// Clear screen of input cookie string
				clear()

				// Check cookie
				cook = strings.Replace(cook, "\r", "", -1)
				cook = strings.Replace(cook, "\n", "", -1)
				if len(cook) < 100 {
					fmt.Println("Invalid cookie")
					continue
				}
				if !json.Valid([]byte(cook)) {
					fmt.Println("Invalid JSON format")
					continue
				}
				if !strings.Contains(cook, "huggingface.co") {
					fmt.Println("Invalid cookie, please make sure the tab is huggingface.co")
					continue

				}

				// Save cookie
				err = ioutil.WriteFile("./5.json", []byte(cook), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}

				// Reload huggingchat cookie
				hc_json, err = ioutil.ReadFile("./5.json")
				hcjs = gjson.Parse(string(hc_json)).String()
				if hcjs == "" {
					continue
				}
				if hcjs != "" {
					//cmd_hc = exec.Command("python3", "-u", "./huggingchat.py", "load")
					cmd_hc = exec.Command("python3", "-u", pf_hc.Name(), "load")
					stdout_hc, _ = cmd_hc.StdoutPipe()
					stdin_hc, _ = cmd_hc.StdinPipe()

					go func(cmd *exec.Cmd) {
						if err := cmd.Start(); err != nil {
							panic(err)
						}
					}(cmd_hc)

					login_hc = false
					relogin_hc = false
					go func(login_hc, relogin_hc *bool) {
						scanner_hc = bufio.NewScanner(stdout_hc)
						for scanner_hc.Scan() {
							RESP = scanner_hc.Text()
							if RESP == "login work" {
								*login_hc = true
							} else if RESP == "relogin" {
								*relogin_hc = true
							} else {
								channel_hc_answer <- RESP
							}
						}
					}(&login_hc, &relogin_hc)
				}
			}
			if relogin_hc == true {
				fmt.Println("Cookie failed, please renew huggingchat cookie...")
				hcjs = ""
				continue

			}
			if login_hc != true {
				fmt.Println("huggingchat initializing...")
				continue
			}

			spc := strings.Replace(userInput, "\n", "(-:]", -1)
			_, err = io.WriteString(stdin_hc, spc+"\n")
			if err != nil {
				panic(err)
			}

			RESP = <-channel_hc_answer
			RESP = strings.Replace(RESP, "(-:]", "\n", -1)
			printer(color_huggingchat, RESP, false)
			save2clip_board(RESP)

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
					Model:    chat_mode, //openai.GPT3Dot5Turbo,
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
			//printer_chat.Println(RESP)
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

var ps_bard = `
import undetected_chromedriver as uc
#from selenium import webdriver as uc
import random,time,os,sys
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support    import expected_conditions as EC
import json
import sys

# Restart session
#########################
#driver = uc.Chrome(options=chrome_options, headless=True)
chrome_options = uc.ChromeOptions()
chrome_options.add_argument("--disable-extensions")
chrome_options.add_argument("--disable-popup-blocking")
chrome_options.add_argument("--profile-directory=Default")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-plugins-discovery")
chrome_options.add_argument("--incognito")
chrome_options.add_argument("--headless")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

# Load cookie
driver.get("https://bard.google.com")
with open("./2.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://bard.google.com")
wait = WebDriverWait(driver, 20)
try:
    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='mat-input-0']")))
    print("login work")
except:
    print("relogin")
   #open("./2.json", "w").close()
    driver.quit()
    os.exit()

wait = WebDriverWait(driver, 30000)
while 1:
   #ori = input(":")
   #if ori:
    for line in sys.stdin:
        message = line.strip()
        ori = message.replace("(-:]", " ")
        work.send_keys(ori)
        driver.find_element(By.XPATH, "//button[@mattooltip='Submit']").click()
       #ini_source = driver.page_source
        if ori:
            try:
                img_thinking = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_thinking_v2_e272afd4f8d4bbd25efe.gif')]")))
               #print("get img_thinking")
                img = wait.until(EC.presence_of_element_located((By.XPATH,  "//img[contains(@src, 'https://www.gstatic.com/lamda/images/sparkle_resting_v2_1ff6f6a71f2d298b1a31.gif')]")))
               #print("get img")
                response = img.find_element(By.XPATH,  "ancestor::model-response")
               #print("get response content img")
                google  = response.find_element(By.XPATH,  ".//button[@aria-label='Google it']")
                
                contents = response.find_elements(By.XPATH, ".//message-content")
                texts= "\n".join(content.text for content in contents)
                text = "(-:]".join(line for line in texts.splitlines() if line)

                text = response.text
                text = text.replace("\n","(-:]")
                text = text.replace("View other drafts","")
                text = text.replace("Regenerate draft","")
                text = text.replace("thumb_up","")
                text = text.replace("thumb_down","")
                text = text.replace("upload","")
                text = text.replace("Google it","")
                text = text.replace("more_vert","")
                text = text.replace("volume_up","")
                text = "(-:]".join(line for line in text.splitlines() if line)
                print(text)
                sys.stdout.flush()

                cookies = driver.get_cookies()
                with open("./2.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass

`

var ps_claude = `
import undetected_chromedriver as uc
import random,time,os,sys
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support    import expected_conditions as EC
import json
import sys

# Restart session
#########################
#driver = uc.Chrome(options=chrome_options, headless=True)
chrome_options = uc.ChromeOptions()
chrome_options.add_argument("--disable-extensions")
chrome_options.add_argument("--disable-popup-blocking")
chrome_options.add_argument("--profile-directory=Default")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-plugins-discovery")
chrome_options.add_argument("--incognito")
chrome_options.add_argument("--headless")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

driver.get("https://claude.ai")

# Load cookie
with open("./3.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://claude.ai")
wait = WebDriverWait(driver, 200)
try:
    work = wait.until(EC.visibility_of_element_located((By.XPATH,  "//p[@data-placeholder='Message Claude or search past chats...']")))
    driver.find_element(By.XPATH, "//div[contains(text(), 'Start a new chat')]").click()
    input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//p[@data-placeholder='Message Claude...']")))
    print("login work")                                                
   #driver.find_element(By.XPATH, "//button[@class='sc-dAOort']").click()
except:
    print("relogin")
   #open("./3.json", "w").close()
    driver.quit()
    os.exit()

while 1:
   #ori = input(":")
   #if ori:
    for line in sys.stdin:
        message = line.strip()
        ori = message.replace("(-:]", " ")
        input_space.send_keys(ori)
        driver.find_element(By.XPATH, "//button[@aria-label='Send Message']").click()
        if ori:
            try:
                retry_icon = wait.until(EC.presence_of_element_located((By.XPATH,  "//svg:path[@d= 'M224,128a96,96,0,0,1-94.71,96H128A95.38,95.38,0,0,1,62.1,197.8a8,8,0,0,1,11-11.63A80,80,0,1,0,71.43,71.39a3.07,3.07,0,0,1-.26.25L44.59,96H72a8,8,0,0,1,0,16H24a8,8,0,0,1-8-8V56a8,8,0,0,1,16,0V85.8L60.25,60A96,96,0,0,1,224,128Z']")))
               #print("get last retry_icon")
                content = retry_icon.find_element(By.XPATH,  "preceding::div[2]")
                text = content.get_attribute("textContent")
                text = text.replace("\n","(-:]")
                print(text)
                sys.stdout.flush()
                
                # Save cookie
                cookies = driver.get_cookies()
                with open("./3.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass
`

var ps_chatgpt = `
import undetected_chromedriver as uc
import random,time,os,sys
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support    import expected_conditions as EC
import json
import sys

# Restart session
#########################
#driver = uc.Chrome(options=chrome_options, headless=True)
chrome_options = uc.ChromeOptions()
chrome_options.add_argument("--disable-extensions")
chrome_options.add_argument("--disable-popup-blocking")
chrome_options.add_argument("--profile-directory=Default")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-plugins-discovery")
chrome_options.add_argument("--incognito")
chrome_options.add_argument("--headless")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

# Load cookie
driver.get("https://chat.openai.com")
with open("./4.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://chat.openai.com")
wait = WebDriverWait(driver, 200)
try:
    notice1 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//h4[contains(text(), 'This is a free research preview.')]")))
   #print("notice1")
    next1 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//button/div[contains(text(), 'Next')]")))
   #print("next1")
    next1.click()
   #print("next1.click")
    notice2 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//h4[contains(text(), 'How we collect data')]")))
   #print("notice2")
    next2 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//button/div[contains(text(), 'Next')]")))
   #print("next2")
    next2.click()
   #print("next2.click")
    notice3 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//h4[contains(text(), 'love your feedback!')]")))
   #print("notice3")
    next3 = wait.until(EC.visibility_of_element_located((By.XPATH,  "//button/div[contains(text(), 'Done')]")))
   #print("next3")
    next3.click()
   #print("next3.click")
    driver.find_element(By.XPATH, "//a[contains(text(), 'New chat')]").click()
    input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='prompt-textarea']")))
    input_space.send_keys("hello")
    driver.find_element(By.XPATH, "//button//svg:path[@d='M.5 1.163A1 1 0 0 1 1.97.28l12.868 6.837a1 1 0 0 1 0 1.766L1.969 15.72A1 1 0 0 1 .5 14.836V10.33a1 1 0 0 1 .816-.983L8.5 8 1.316 6.653A1 1 0 0 1 .5 5.67V1.163Z']").click()
    print("login work")
except:
    print("relogin")
   #open("./2.json", "w").close()
    driver.quit()
    os.exit()

while 1:
   #ori = input(":")
   #if ori:
    for line in sys.stdin:
        message = line.strip()
        ori = message.replace("(-:]", " ")
        input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@id='prompt-textarea']")))
        input_space.send_keys(ori)
        driver.find_element(By.XPATH, "//button//svg:path[@d='M.5 1.163A1 1 0 0 1 1.97.28l12.868 6.837a1 1 0 0 1 0 1.766L1.969 15.72A1 1 0 0 1 .5 14.836V10.33a1 1 0 0 1 .816-.983L8.5 8 1.316 6.653A1 1 0 0 1 .5 5.67V1.163Z']").click()
       #ini_source = driver.page_source
        if ori:
            try:
                retry_icon = wait.until(EC.presence_of_element_located((By.XPATH,  "//svg:path[@d='M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15']")))
               #print("get retry_icon")
                content = retry_icon.find_element(By.XPATH,  "(//div[contains(@class, 'group w-full')])[last()]")
                text = content.get_attribute("textContent")
                text = text.replace("ChatGPTChatGPT","")
                text = text.replace("1 / 1","")
                text = text.replace("\n","(-:]")
                print(text)
                sys.stdout.flush()
                cookies = driver.get_cookies()
                with open("./4.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass
`

var ps_huggingchat = `
import undetected_chromedriver as uc
import random,time,os,sys
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support    import expected_conditions as EC
import json
import sys

#########################
chrome_options = uc.ChromeOptions()
chrome_options.add_argument("--disable-extensions")
chrome_options.add_argument("--disable-popup-blocking")
chrome_options.add_argument("--profile-directory=Default")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-plugins-discovery")
chrome_options.add_argument("--incognito")
chrome_options.add_argument("--headless")
chrome_options.add_argument("user_agent=DN")
driver = uc.Chrome(options=chrome_options)

driver.get("https://huggingface.co/chat")
with open("./5.json", "r", newline='') as inputdata:
    ck = json.load(inputdata)
for c in ck:
    driver.add_cookie({k:c[k] for k in {'name', 'value'}})

# Renew with cookie
driver.get("https://huggingface.co/chat")
wait = WebDriverWait(driver, 200)
try:
    input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@enterkeyhint='send']")))
    print("login work")
    input_space.send_keys("hello")
    driver.find_element(By.XPATH, "//button//svg:path[@d='M27.71 4.29a1 1 0 0 0-1.05-.23l-22 8a1 1 0 0 0 0 1.87l8.59 3.43L19.59 11L21 12.41l-6.37 6.37l3.44 8.59A1 1 0 0 0 19 28a1 1 0 0 0 .92-.66l8-22a1 1 0 0 0-.21-1.05Z']").click()
except:
    print("relogin")
    driver.quit()
    os.exit()

while 1:
   #ori = input(":")
   #if ori:
    for line in sys.stdin:
        message = line.strip()
        ori = message.replace("(-:]", " ")
        input_space = wait.until(EC.visibility_of_element_located((By.XPATH,  "//textarea[@enterkeyhint='send']")))
        input_space.send_keys(ori)
        driver.find_element(By.XPATH, "//button//svg:path[@d='M27.71 4.29a1 1 0 0 0-1.05-.23l-22 8a1 1 0 0 0 0 1.87l8.59 3.43L19.59 11L21 12.41l-6.37 6.37l3.44 8.59A1 1 0 0 0 19 28a1 1 0 0 0 .92-.66l8-22a1 1 0 0 0-.21-1.05Z']").click()
        if ori:
            try:
                stop_icon = wait.until(EC.presence_of_element_located((By.XPATH,  "//svg:path[@d='M24 6H8a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2Z']")))
               #print("get stop_icon")
                wait.until(EC.staleness_of(stop_icon))
               #print("disappear stop_icon")
                img = driver.find_element(By.XPATH,  "(//img[contains(@src, 'https://huggingface.co/avatars/2edb18bd0206c16b433841a47f53fa8e.svg')])[last()]")
               #print("img")
                content = img.find_element(By.XPATH,  "following-sibling::div[1]")
                text = content.get_attribute("textContent")
                text = text.replace("\n","(-:]")
                print(text)
                sys.stdout.flush()
                cookies = driver.get_cookies()
                with open("./5.json", "w", newline='') as outputdata:
                    json.dump(cookies, outputdata)

            except Exception as e:
                pass
`
