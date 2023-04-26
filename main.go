package main

import (
	//	"bufio"
	"context"
	"fmt"
	"github.com/CNZeroY/googleBard/bard"
	"github.com/PuerkitoBio/goquery"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/headzoo/surf"
	"github.com/pavel-one/EdgeGPT-Go"
	"github.com/peterh/liner"
	"github.com/rocketlaunchr/google-search"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sohaha/cursor"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	// Read json configure
	data, err := ioutil.ReadFile("aih.json")
	////

	//Read Proxy
	proxy := gjson.Get(string(data), "proxy")
	Proxy := proxy.String()

	if Proxy != "" {
		os.Setenv("https_proxy", "http://127.0.0.1:7890")
		os.Setenv("http_proxy", "http://127.0.0.1:7890")
		//proxyUrl, err := url.Parse(Proxy)
		//if err != nil {
		//	panic(err)
		//}
		//transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		//// for openai api
		//config.HTTPClient = &http.Client{Transport: transport}
		//// for normal page
		//client_n.Transport = transport
		//bow.SetTransport(transport)
	}

	liner := liner.NewLiner()
	defer liner.Close()
	/////
	//	if err != nil {
	//		//if err == nil {
	//		//var okey string
	//		//fmt.Println("Please input your OpenAI Key: ")
	//		okey, _ := liner.Prompt("Please input your OpenAI Key: ")
	//		conf := `{"key":"` + okey + `"}`
	//		//fmt.Println(conf)
	//		err := ioutil.WriteFile("aih.json", []byte(conf), 0644)
	//		if err != nil {
	//			fmt.Println("Save failed.")
	//		}
	//	}
	//

	//Read OpenAI_Key
	key := gjson.Get(string(data), "key")
	OpenAI_Key := key.String()

	//Read Google Cookie of __Secure-lPSID
	__Secure_lPSID := gjson.Get(string(data), "__Secure-lPSID")
	bard_session_id := __Secure_lPSID.String()

	// Set up client config of OpenAI API
	config := openai.DefaultConfig(OpenAI_Key)

	// Set up client for normal page
	client_n := &http.Client{}
	client_n.Timeout = time.Second * 10
	bow := surf.NewBrowser()

	// Set up client for OpenAI API
	client := openai.NewClientWithConfig(config)
	messages := make([]openai.ChatCompletionMessage, 0)

	// Set up client for google_bard
	bard_client := bard.NewBard(bard_session_id, Proxy)
	bardOptions := bard.Options{
		ConversationID: "",
		ResponseID:     "",
		ChoiceID:       "",
	}
	printer_bard := color.New(color.FgGreen).Add(color.Bold)

	// Set up client for bing chat
	s := EdgeGPT.NewStorage()
	gpt, err := s.GetOrSet("any-key")
	if err != nil {
		panic(err)
	}
	printer_bing := color.New(color.FgBlue).Add(color.Bold)

	fmt.Println("Welcome to aih v0.1.0\nType \".help\" for more information.")
	max_tokens := 4097
	used_tokens := 0
	left_tokens := max_tokens - used_tokens
	speak := 0
	role := ""

	////
	if f, err := os.Open(".history"); err == nil {
		liner.ReadHistory(f)
		f.Close()
	}
	/////

	// Start loop to read user input and setn API requests
	for {
		/////
		promp := strconv.Itoa(left_tokens) + role + "> "
		userInput, _ := liner.Prompt(promp)
		liner.AppendHistory(userInput)

		userInput = strings.Trim(userInput, " ") // remove space after .xxx
		clipb := ""                              // for save to system clipboard

		switch userInput {
		case "":
			continue
		case ".exit":
			fmt.Println("Byebye")
			return
		case ".proxy":
			proxy, _ := liner.Prompt("Please input your proxy: ")
			data, err := ioutil.ReadFile("aih.json")
			sdata := string(data)
			njs, _ := sjson.Set(sdata, "proxy", proxy)
			err = ioutil.WriteFile("aih.json", []byte(njs), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
			fmt.Println("Please restart aih")
			continue
		case ".key":
			k, _ := liner.Prompt("Please input your OpenAI key: ")
			data, err := ioutil.ReadFile("aih.json")
			sdata := string(data)
			nnjs, _ := sjson.Set(sdata, "key", k)
			err = ioutil.WriteFile("aih.json", []byte(nnjs), 0644)
			if err != nil {
				fmt.Println("Save failed.")
			}
			fmt.Println("Please restart aih")
			continue
		case ".help":
			//fmt.Println(".info      Print the information")
			fmt.Println(".help        Show help")
			fmt.Println(".key         Set key")
			fmt.Println(".proxy       Set proxy")
			fmt.Println(".new         New conversation")
			fmt.Println(".speak       Voice speak context")
			fmt.Println(".quiet       Quiet not speak")
			fmt.Println(".clear       Clear screen")
			fmt.Println(".update      Inquery up-to-date question")
			fmt.Println(".code        Code creation by Cursor")
			fmt.Println(".exit        Exit")
			fmt.Println(" ------roles------")
			fmt.Println(".prompt      Role of Assistant for create precise prompt")
			fmt.Println(".writer      Role of Checker for create well sentences")
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
			continue
		case ".prompt":
			messages = make([]openai.ChatCompletionMessage, 0)
			max_tokens = 4097
			used_tokens = 0
			left_tokens = max_tokens - used_tokens
			userInput = prompt_prove
			role = ".prompt"
		case ".writer":
			messages = make([]openai.ChatCompletionMessage, 0)
			max_tokens = 4097
			used_tokens = 0
			left_tokens = max_tokens - used_tokens
			userInput = write_prove
			role = ".writer"
		case ".update":
			role = ".update"
			continue
		case ".code":
			role = ".code"
			continue
		case ".bard":
			role = ".bard"

			data, err := ioutil.ReadFile("aih.json")
			__Secure_lPSID := gjson.Get(string(data), "__Secure-lPSID")
			bard_session_id := __Secure_lPSID.String()
			if bard_session_id == "" {
				bard_session_id, _ := liner.Prompt("Please input your cookie value of __Secure-lPSID: ")
				sdata := string(data)
				njs, _ := sjson.Set(sdata, "__Secure-lPSID", bard_session_id)
				err = ioutil.WriteFile("aih.json", []byte(njs), 0644)
				if err != nil {
					fmt.Println("Save failed.")
				}
				// renew bard client with session id
				bard_client = bard.NewBard(bard_session_id, Proxy)
				left_tokens = 0
				role = ".bard"
				continue
			}
			left_tokens = 0
			continue
		case ".bing":
			role = ".bing"
			left_tokens = 0
			continue
		}

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
					//http
					//	req, err := http.NewRequest("GET", durl, nil)
					//	if err != nil {
					//		fmt.Println(err)
					//		return
					//	}
					//	req.Header = headers
					//	resp_p, err := client_n.Do(req)
					//	if err != nil {
					//		fmt.Println(err)
					//		return
					//	}

					//	defer resp_p.Body.Close()
					//	cnt_p, err := ioutil.ReadAll(resp_p.Body)
					//	if err != nil {
					//		fmt.Println(err)
					//		return
					//	}
					//---------------

					//surf
					err := bow.Open(durl)
					if err != nil {
						//fmt.Println(err)
						return
					}
					cnt_p := bow.Body()
					//---------------

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
			printer_bing.Println(as.Answer.GetAnswer())
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

var prompt_prove = `I want you to become my Prompt Creator. Your goal is to help me craft the best possible prompt for my needs. The prompt will be used by you, ChatGPT. You will follow the following process: 1. Your first response will be to ask me what the prompt should be about. I will provide my answer, but we will need to improve it through continual iterations by going through the next steps. 2. Based on my input, you will generate 3 sections. a) Revised prompt (provide your rewritten prompt. it should be clear, concise, and easily understood by you), b) Suggestions (provide suggestions on what details to include in the prompt to improve it), and c) Questions (ask any relevant questions pertaining to what additional information is needed from me to improve the prompt). 3. We will continue this iterative process with me providing additional information to you and you updating the prompt in the Revised prompt section until it's complete.`

var write_prove = `I want you to become my Write Checker. Your goal is to help me craft the best possible sentences for my needs. You will follow the following process: 1. Your first response will be to ask me what is the sentences. I will provide my answer, but we will need to improve it through continual iterations by going through the next steps. 2. Based on my input, you will generate 3 sections. a) Revised sentences (provide your rewritten sentences. it should be clear, concise, and easily understood), b) Suggestions (provide suggestions on what details to include in the sentence to improve it), and c) Questions (ask any relevant questions pertaining to what additional information is needed from me to improve the sentence). 3. We will continue this iterative process with me providing additional information to you and you updating the sentence in the Revised sentences section until it's complete.`
