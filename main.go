package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"io/ioutil"
	"github.com/fatih/color"
	"github.com/peterh/liner"
	"github.com/sohaha/cursor"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/atotto/clipboard"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/CNZeroY/googleBard/bard"
	"github.com/rocketlaunchr/google-search"
	openai "github.com/sashabaranov/go-openai"
)

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

func main() {

	// Create prompt for user input
	Liner := liner.NewLiner()
	defer Liner.Close()
	if f, err := os.Open(".history"); err == nil {
		Liner.ReadHistory(f)
		f.Close()
	}
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

	// Test Proxy
	fmt.Println("Checking network accessing...")
	ops1 := googlesearch.SearchOptions{Limit: 12}
	_, err = googlesearch.Search(nil, "BTC", ops1)
	if err != nil {
		fmt.Println("Need proxy to access GoogleBard, BingChat, ChatGPT")
		proxy, _ := Liner.Prompt("Please input proxy: ")
		aihj, err := ioutil.ReadFile("aih.json")
		new_aihj, _ := sjson.Set(string(aihj), "proxy", proxy)
		err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
		if err != nil {
			fmt.Println("Save failed.")
		}
		fmt.Println("Please restart aih for using proxy...")
		Liner.Close()
		syscall.Exit(0)

	}

	// Set up client for normal_page
	client_n := &http.Client{}
	client_n.Timeout = time.Second * 10
	//bow := surf.NewBrowser()

	// Set up client for OpenAI_API
	key := gjson.Get(string(aih_json), "key")
	OpenAI_Key := key.String()
	config := openai.DefaultConfig(OpenAI_Key)
	client := openai.NewClientWithConfig(config)
	messages := make([]openai.ChatCompletionMessage, 0)

	// Set up client for GoogleGard
	bard_session_id := gjson.Get(string(aih_json), "__Secure-lPSID").String()
	bard_client := bard.NewBard(bard_session_id, "")
	bardOptions := bard.Options{
		ConversationID: "",
		ResponseID:     "",
		ChoiceID:       "",
	}
	printer_bard := color.New(color.FgRed).Add(color.Bold)

	// Set up client for BingChat
	var gpt *EdgeGPT.GPT
	_, err = ioutil.ReadFile("./cookies/1.json")
	if err == nil {
		s := EdgeGPT.NewStorage()
		gpt, err = s.GetOrSet("any-key")
	}
	printer_bing := color.New(color.FgCyan).Add(color.Bold)

	// Clean screen
	clear()

	// Welcome to Aih
	fmt.Println("---------------------")
	fmt.Println("Welcome to Aih v0.1.0")
	fmt.Println("Type .help for help")
	fmt.Println("---------------------")
	max_tokens := 4097
	used_tokens := 0
	//left_tokens := max_tokens - used_tokens
	left_tokens := 0
	speak := 0
	role := ".bard"

	// Start loop to read user input
	for {
		promp := strconv.Itoa(left_tokens) + role + "> "
		userInput, _ := Liner.Prompt(promp)
		userInput = strings.Trim(userInput, " ") // remove side space

		// for save to system clipboard
		clipb := ""

		// Check Aih commands
		switch userInput {
		case "":
			continue
		case ".proxy":
			proxy, _ := Liner.Prompt("Please input your proxy:")
			aihj, err := ioutil.ReadFile("aih.json")
			new_aihj, _ := sjson.Set(string(aihj), "proxy", proxy)
			err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
			fmt.Println("Please restart Aih for using proxy")
			Liner.Close()
			syscall.Exit(0)
		case ".chatkey":
			//k, _ := Liner.Prompt("Please input your OpenAI key: ")
			//aihj, err := ioutil.ReadFile("aih.json")
			//new_aihj, _ := sjson.Set(string(aihj), "key", k)
			//err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
			//if err != nil {
			//	fmt.Println("Save failed.")
			//}
			//fmt.Println("Please restart Aih for using key")
                        OpenAI_Key = ""
			role = ".chat"
			continue
	        case ".bardkey":		
                        bard_session_id = "" 
			role = ".bard"
			continue
	        case ".bingkey":		
			err := os.Remove("./cookies/1.json")
			if err != nil {panic(err)}
			role = ".bing"
			continue
		case ".help":
			fmt.Println(".bard        Bard")
			fmt.Println(".bing        Bing")
			fmt.Println(".chat        ChatGPT")
			fmt.Println(".help        Help")
			fmt.Println(".proxy       Set proxy")
			fmt.Println(".chatkey     Set ChatGPT key")
			fmt.Println(".new         New conversation of ChatGPT")
			fmt.Println(".speak       Voice speak context")
			fmt.Println(".quiet       Not speak")
			fmt.Println(".clear       Clear screen")
			//fmt.Println(".code        Code creation by Cursor")
			fmt.Println(".exit        Exit")
			continue
		case ".speak":
			speak = 1
			continue
		case ".quiet":
			speak = 0
			continue
		case ".clear":
			clear()
			continue
		case ".exit":
			return
		case ".new":
			role = "chat"
			messages = make([]openai.ChatCompletionMessage, 0)
			max_tokens = 4097
			used_tokens = 0
			left_tokens = max_tokens - used_tokens
			continue
		//case ".code":
		//	role = ".code"
		//	continue
		case ".bard":
			role = ".bard"
			left_tokens = 0
			continue
		case ".bing":
			role = ".bing"
			left_tokens = 0
			continue
		case ".chat":
			role = ".chat"
			left_tokens = max_tokens - used_tokens
			continue
		}

		// Record user input without Aih commands
		Liner.AppendHistory(userInput)

		var RESP string
		// Check role for currect actions
		if role == ".code" {
			res_code, err := cursor.Conv(userInput)
			if err != nil {
				panic(err)
				return
			}
			cg := color.New(color.FgGreen)
			cg.Println(res_code)

			// write to clipboard
			err = clipboard.WriteAll(res_code)
			if err != nil {
				panic(err)
				return
			}
			continue
		}

		if role == ".bard" {
			// Check GoogleBard session
			if bard_session_id == "" {
				bard_session_id, _ = Liner.Prompt("Please input your cookie value of __Secure-lPSID: ")
				aihj, err := ioutil.ReadFile("aih.json")
				nj, _ := sjson.Set(string(aihj), "__Secure-lPSID", bard_session_id)
				err = ioutil.WriteFile("aih.json", []byte(nj), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}
				// Renew GoogleBard client with __Secure-lPSID
				bard_client = bard.NewBard(bard_session_id, "")
				continue
			}

			// Send message
			response, err := bard_client.SendMessage(userInput, bardOptions)
			if err != nil {
				panic(err)
			}

			all_resp := response
			if all_resp != nil {
				RESP = response.Choices[0].Answer
				printer_bard.Println(RESP)
				// Write to clipboard
				err = clipboard.WriteAll(RESP)
				if err != nil {
					panic(err)
					return
				}
			} else {
				break
			}
			bardOptions.ConversationID = response.ConversationID
			bardOptions.ResponseID = response.ResponseID
			bardOptions.ChoiceID = response.Choices[0].ChoiceID
			//continue

		}

		if role == ".bing" {
			// Check BingChat cookie
			_, err := ioutil.ReadFile("./cookies/1.json")
			if err != nil {
				var lines []string
				fmt.Println("Please paste bing cookie here then press Enter then Ctrl+D:")
				for {
					line, err := Liner.Prompt("")
					if err == io.EOF {
						break
					}
					lines = append(lines, line)
				}
				longString := strings.Join(lines, "\n")
				_ = os.MkdirAll("./cookies", 0755)
				err = ioutil.WriteFile("./cookies/1.json", []byte(longString), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}
				// Renew BingChat client with cookie
				s := EdgeGPT.NewStorage()
				gpt, err = s.GetOrSet("any-key")
				// Clear screen
				clear()
				continue
			}
			// Send message
			as, err := gpt.AskSync("creative", userInput)
			if err != nil {
				panic(err)
			}
			RESP = strings.TrimSpace(as.Answer.GetAnswer())
			printer_bing.Println(RESP)
			//continue

		}
		if role == ".chat" {
		        // Check ChatGPT Key
			if OpenAI_Key == "" {
				OpenAI_Key, _ = Liner.Prompt("Please input your OpenAI Key: ")
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
					Model:    openai.GPT3Dot5Turbo,
					Messages: messages,
				},
			)

			if err != nil {
				fmt.Println(err)
				continue
			}

			// Print the response to the terminal
			c := color.New(color.FgWhite)
			RESP = strings.TrimSpace(resp.Choices[0].Message.Content)
			used_tokens = resp.Usage.TotalTokens
			left_tokens = max_tokens - used_tokens
			c.Println(RESP)

			// Write to clipboard
			clipb += fmt.Sprintln(RESP)
			err = clipboard.WriteAll(clipb)
			if err != nil {
				panic(err)
				return
			}

			// Record in coversation context
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: RESP,
			})

		}

			// Speak the response using the "say" command
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
						_ = 1 + 1

					}

				}()
			}
	}
}
