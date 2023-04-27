package main

import (
	"os"
	"io"
	"fmt"
	"sync"
	"time"
	"syscall"
	"context"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/fatih/color"
	"github.com/peterh/liner"
	"github.com/headzoo/surf"
	"github.com/sohaha/cursor"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/atotto/clipboard"
	"github.com/PuerkitoBio/goquery"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/CNZeroY/googleBard/bard"
	"github.com/rocketlaunchr/google-search"
	openai "github.com/sashabaranov/go-openai"

       )

func clear(){
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
		fmt.Println("Need setting proxy to access bard, bing, chatGPT")
		proxy, _ := Liner.Prompt("Please input your proxy: ")
		data, err := ioutil.ReadFile("aih.json")
		sdata := string(data)
		njs, _ := sjson.Set(sdata, "proxy", proxy)
		err = ioutil.WriteFile("aih.json", []byte(njs), 0644)
		if err != nil {
			fmt.Println("Save failed.")
		}
		fmt.Println("Please restart aih for using proxy")
		Liner.Close()
		syscall.Exit(0)

	}

	// Set up client for normal_page
	client_n := &http.Client{}
	client_n.Timeout = time.Second * 10
	bow := surf.NewBrowser()

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
	fmt.Println("Welcome to Aih v0.1.0\nType \".help\" for help")
	fmt.Println("---------------------")
	max_tokens := 4097
	used_tokens := 0
	left_tokens := max_tokens - used_tokens
	speak := 0
	role := ".bard"

	// Start loop to read user input
	for {
		promp := strconv.Itoa(left_tokens) + role + "> "
		userInput, _ := Liner.Prompt(promp)
		userInput = strings.Trim(userInput, " ") // remove side space

		// for save to system clipboard
		clipb := ""

		switch userInput {
		case "":
			continue
		case ".exit":
			return
		case ".proxy":
			proxy, _ := Liner.Prompt("Please input your proxy:")
			aihj, err := ioutil.ReadFile("aih.json")
			str_aihj := string(aihj)
			new_aihj, _ := sjson.Set(str_aihj, "proxy", proxy)
			err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
		        fmt.Println("Please restart aih for using proxy")
		        Liner.Close()
		        syscall.Exit(0)
		case ".key":
			k, _ := Liner.Prompt("Please input your OpenAI key: ")
			aihj, err := ioutil.ReadFile("aih.json")
			str_aihj := string(aihj)
			new_aihj, _ := sjson.Set(str_aihj, "key", k)
			err = ioutil.WriteFile("aih.json", []byte(new_aihj), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
			fmt.Println("Please restart aih")
			continue
		case ".help":
			//fmt.Println(".info      Print the information")
			fmt.Println(".bard        Bard")
			fmt.Println(".bing        Bing")
			fmt.Println(".chat        ChatGPT")
			fmt.Println(".help        Help")
			fmt.Println(".key         Set key")
			fmt.Println(".proxy       Set proxy")
			fmt.Println(".new         New conversation of ChatGPT")
			fmt.Println(".speak       Voice speak context")
			fmt.Println(".quiet       Not speak")
			fmt.Println(".clear       Clear screen")
			//fmt.Println(".update      Inquery up-to-date question")
			//fmt.Println(".code        Code creation by Cursor")
			fmt.Println(".exit        Exit")
			continue
		case ".speak":
			speak = 1
			continue
		case ".new":
			messages = make([]openai.ChatCompletionMessage, 0)
			max_tokens = 4097
			used_tokens = 0
			left_tokens = max_tokens - used_tokens
			role = ""
			continue
		case ".quiet":
			speak = 0
			continue
		case ".clear":
		        clear()
			//switch runtime.GOOS {
			//case "linux", "darwin":
			//	cmd := exec.Command("clear")
			//	cmd.Stdout = os.Stdout
			//	cmd.Run()
			//case "windows":
			//	cmd := exec.Command("cmd", "/c", "cls")
			//	cmd.Stdout = os.Stdout
			//	cmd.Run()
			//}
			continue
		//case ".prompt":
		//	messages = make([]openai.ChatCompletionMessage, 0)
		//	max_tokens = 4097
		//	used_tokens = 0
		//	left_tokens = max_tokens - used_tokens
		//	userInput = prompt_prove
		//	role = ".prompt"
		//case ".writer":
		//	messages = make([]openai.ChatCompletionMessage, 0)
		//	max_tokens = 4097
		//	used_tokens = 0
		//	left_tokens = max_tokens - used_tokens
		//	userInput = write_prove
		//	role = ".writer"
		case ".update":
			role = ".update"
			continue
		case ".code":
			role = ".code"
			continue
		case ".bard":
			role = ".bard"

			if bard_session_id == "" {
				bard_session_id, _ = Liner.Prompt("Please input your cookie value of __Secure-lPSID: ")
				data, err := ioutil.ReadFile("aih.json")
				sdata := string(data)
				njs, _ := sjson.Set(sdata, "__Secure-lPSID", bard_session_id)
				err = ioutil.WriteFile("aih.json", []byte(njs), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}
				// renew bard client with session id
				bard_client = bard.NewBard(bard_session_id, "")
				fmt.Println("Renew bard client with session id ready")
				left_tokens = 0
				continue
			}

			left_tokens = 0
			continue

		case ".bing":
			role = ".bing"

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

				// renew bing client with cookie
				s := EdgeGPT.NewStorage()
				gpt, err = s.GetOrSet("any-key")
				// Clean screen
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

			left_tokens = 0
			continue

		case ".chat":
			role = ".chat"

			key := gjson.Get(string(aih_json), "key")
			OpenAI_Key := key.String()
			if OpenAI_Key == "" {
				okey, _ := Liner.Prompt("Please input your OpenAI Key: ")
				conf := `{"key":"` + okey + `"}`
				err := ioutil.WriteFile("aih.json", []byte(conf), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}
				// renew chatgpt client with key
				config = openai.DefaultConfig(OpenAI_Key)
				client = openai.NewClientWithConfig(config)
				messages = make([]openai.ChatCompletionMessage, 0)
				left_tokens = 0
				continue
			}
			continue
		}

		Liner.AppendHistory(userInput)

		if role == ".update" {
			// Generate a abstract response from ChatGPT
			prompt := "Please abstract or extent keywords from this message for precise information on search engine in one line separate by ',' : " + userInput
			resp_, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt,
						}},
				},
			)

			if err != nil {
				fmt.Println(err)
				continue
			}
			cc := color.New(color.FgYellow)
			key_ := strings.TrimSpace(resp_.Choices[0].Message.Content)
			datetime := time.Now().Format("2006-01-02")
			key_ = key_ + ", " + datetime
			cc.Println("Key:", key_)
			clipb += fmt.Sprintln("Key:", key_)
			// Search in google
			results := make([]googlesearch.Result, 0)
			ops1 := googlesearch.SearchOptions{Limit: 12}
			if Proxy != "" {
				ops := googlesearch.SearchOptions{ProxyAddr: Proxy}
				results, _ = googlesearch.Search(nil, key_, ops, ops1)
			} else {
				results, _ = googlesearch.Search(nil, key_, ops1)
			}
			cc.Println("------up-to-date------")
			clipb += fmt.Sprintln("------up-to-date------")
			for index, i := range results {
				cc.Print("[", index, "] ")
				cc.Println(i.URL)
				cc.Println(i.Title)
				clipb += fmt.Sprintln("[", index, "]")
				clipb += fmt.Sprintln(i.URL)
				clipb += fmt.Sprintln(i.Title)
			}

			var wg sync.WaitGroup
			pages := ""
			headers := http.Header{}
			headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/538.36")
			headers.Set("Accept-Language", "")
			for index, i := range results {
				durl := i.URL
				wg.Add(1)
				go func(index int) {
					defer wg.Done()
					//surf
					err := bow.Open(durl)
					if err != nil {
						//fmt.Println(err)
						return
					}
					cnt_p := bow.Body()

					raw_page := string(cnt_p)
					r_page, err := goquery.NewDocumentFromReader(strings.NewReader(raw_page))
					page := r_page.Text()
					page = strings.ReplaceAll(page, "\n", " ")
					page = strings.ReplaceAll(page, "  ", " ")
					//fmt.Println(page)
					if len(page) > 6000 {
						page = page[:6000]
					}
					// Generate a summary response from ChatGPT
					prompt_p := "Please abstract usefull information about `" + userInput + "` from message below, today is" + datetime + ", if no usefull information, return 0: " + page
					resps, err := client.CreateChatCompletion(
						context.Background(),
						openai.ChatCompletionRequest{
							Model: openai.GPT3Dot5Turbo,
							Messages: []openai.ChatCompletionMessage{
								{
									Role:    openai.ChatMessageRoleUser,
									Content: prompt_p,
								}},
						},
					)

					if err != nil {
						//fmt.Println(err)
						return
					}
					summary := strings.TrimSpace(resps.Choices[0].Message.Content)
					//fmt.Println("[", index, "] ", summary)
					pages += summary
				}(index)
			}
			wg.Wait()
			// Generate a summary_total response from ChatGPT
			if len(pages) > 7000 {
				pages = pages[:7000]
			}

			prompt_ps := "Please well manage information from message below, for precise conscise, ignore useless answer, only useful answer: " + pages
			resps, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt_ps,
						}},
				},
			)

			if err != nil {
				fmt.Println(err)
				return
				//continue
			}
			summary_total := strings.TrimSpace(resps.Choices[0].Message.Content)
			cc.Println("------summary------")
			cc.Println(summary_total)
			cc.Println("-------------------")

			clipb += fmt.Sprintln("------summary------")
			clipb += fmt.Sprintln(summary_total)
			clipb += fmt.Sprintln("-------------------")
		}

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
			response, err := bard_client.SendMessage(userInput, bardOptions)
			if err != nil {
				panic(err)
			}

			all_resp := response
			if all_resp != nil {
				resp := response.Choices[0].Answer
				printer_bard.Println(resp)
				// write to clipboard
				err = clipboard.WriteAll(resp)
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

			continue

		}

		if role == ".bing" {
			as, err := gpt.AskSync("creative", userInput)
			if err != nil {
				panic(err)
			}
		        res := strings.TrimSpace(as.Answer.GetAnswer())
			printer_bing.Println(res)
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
		cnt := strings.TrimSpace(resp.Choices[0].Message.Content)
		used_tokens = resp.Usage.TotalTokens
		left_tokens = max_tokens - used_tokens
		c.Println(cnt)

		// write to clipboard
		clipb += fmt.Sprintln(cnt)
		err = clipboard.WriteAll(clipb)
		if err != nil {
			panic(err)
			return
		}

		// Speak the response using the "say" command
		if speak == 1 {
			go func() {
				switch runtime.GOOS {
				case "linux", "darwin":
					cmd := exec.Command("say", cnt)
					err = cmd.Run()
					if err != nil {
						fmt.Println(err)
					}
				case "windows":
					_ = 1 + 1

				}

			}()
		}

		// record in coversation context
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: cnt,
		})

	}
}

