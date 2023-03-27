package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
//	"github.com/headzoo/surf"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	//	"strconv"
	//	"github.com/eiannone/keyboard"
	//	"github.com/nsf/termbox-go"
	"github.com/rocketlaunchr/google-search"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func main() {

	// Read json configure
	data, err := ioutil.ReadFile("aih.json")
	if err != nil {
		//if err == nil {
		var okey string
		fmt.Println("Please input your OpenAI Key: ")
		fmt.Scanln(&okey)
		conf := `{"key":"` + okey + `"}`
		//fmt.Println(conf)
		err := ioutil.WriteFile("aih.json", []byte(conf), 0644)
		if err != nil {
			fmt.Println("Save failed.")
		}
	}

	//Read OpenAI_Key
	key := gjson.Get(string(data), "key")
	OpenAI_Key := key.String()

	//Read Proxy
	proxy := gjson.Get(string(data), "proxy")
	Proxy := proxy.String()

	// Set up the OpenAI API client
	config := openai.DefaultConfig(OpenAI_Key)

	// for normal page
	client_n := &http.Client{}
	client_n.Timeout = time.Second * 10
	//bow := surf.NewBrowser()

	if Proxy != "" {
		proxyUrl, err := url.Parse(Proxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		// for openai api
		config.HTTPClient = &http.Client{Transport: transport}
		// for normal page
		//client_n := &http.Client{Transport: transport}
		client_n.Transport = transport
	}

	client := openai.NewClientWithConfig(config)
	messages := make([]openai.ChatCompletionMessage, 0)

	fmt.Println("Welcome to aih v0.1.0\nType \".help\" for more information.")
	// Start loop to read user input and setn API requests
	scanner := bufio.NewScanner(os.Stdin)
	max_tokens := 4097
	used_tokens := 0
	left_tokens := max_tokens - used_tokens
	speak := 0
	role := ""

	for {
		fmt.Print(left_tokens, role, "> ")

		// Parse the command line arguments to get the prompt
		scanner.Scan()
		userInput := scanner.Text()

		switch userInput {
		case "":
			continue
		case ".exit":
			fmt.Println("Byebye")
			return
		case ".proxy":
			var proxy string
			fmt.Println("Please input your proxy: ")
			fmt.Scanln(&proxy)
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
			var proxy string
			fmt.Println("Please input your OpenAI key: ")
			fmt.Scanln(&proxy)
			data, err := ioutil.ReadFile("aih.json")
			sdata := string(data)
			nnjs, _ := sjson.Set(sdata, "key", proxy)
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
			cmd := ""
			if runtime.GOOS != "windows" {
				cmd = "clear"
				clear := exec.Command(cmd)
				clear.Stdout = os.Stdout
				clear.Run()
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
			cc.Println("Key:", key_)
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
			for index, i := range results {
				cc.Print("[", index, "] ")
				cc.Println(i.URL)
				cc.Println(i.Title)
			}
			cc.Println("------summary------")
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
						req, err := http.NewRequest("GET", durl, nil)
						if err != nil {
							fmt.Println(err)
							return
						}
						req.Header = headers
						resp_p, err := client_n.Do(req)
						if err != nil {
							fmt.Println(err)
							return
						}

						defer resp_p.Body.Close()
						cnt_p, err := ioutil.ReadAll(resp_p.Body)
						if err != nil {
							fmt.Println(err)
							return
						}
					//---------------

					//surf
				//	err := bow.Open(durl)
				//	if err != nil {
				//		fmt.Println(err)
				//		return
				//	}
				//	cnt_p := bow.Body()
					//---------------

					page := string(cnt_p)
					//fmt.Println(page)
					if len(page) > 6000 {
						page = page[:6000]
					}
					// Generate a summary response from ChatGPT
					prompt_p := "Please abstract usefull information about `" + userInput + "` from message below, if no usefull information, return 0: " + page
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
						fmt.Println(err)
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
			//fmt.Println(">>", pages)
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
				continue
			}
			summary_total := strings.TrimSpace(resps.Choices[0].Message.Content)
			cc.Println(summary_total)
			cc.Println("-------------------")

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

		// Speak the response using the "say" command
		if speak == 1 {
			go func() {
				cmd := exec.Command("say", cnt)
				err = cmd.Run()
				if err != nil {
					fmt.Println(err)
				}
			}()
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: cnt,
		})

	}
}

var prompt_prove = `I want you to become my Prompt Creator. Your goal is to help me craft the best possible prompt for my needs. The prompt will be used by you, ChatGPT. You will follow the following process: 1. Your first response will be to ask me what the prompt should be about. I will provide my answer, but we will need to improve it through continual iterations by going through the next steps. 2. Based on my input, you will generate 3 sections. a) Revised prompt (provide your rewritten prompt. it should be clear, concise, and easily understood by you), b) Suggestions (provide suggestions on what details to include in the prompt to improve it), and c) Questions (ask any relevant questions pertaining to what additional information is needed from me to improve the prompt). 3. We will continue this iterative process with me providing additional information to you and you updating the prompt in the Revised prompt section until it's complete.`

var write_prove = `I want you to become my Write Checker. Your goal is to help me craft the best possible sentences for my needs. You will follow the following process: 1. Your first response will be to ask me what is the sentences. I will provide my answer, but we will need to improve it through continual iterations by going through the next steps. 2. Based on my input, you will generate 3 sections. a) Revised sentences (provide your rewritten sentences. it should be clear, concise, and easily understood), b) Suggestions (provide suggestions on what details to include in the sentence to improve it), and c) Questions (ask any relevant questions pertaining to what additional information is needed from me to improve the sentence). 3. We will continue this iterative process with me providing additional information to you and you updating the sentence in the Revised sentences section until it's complete.`
