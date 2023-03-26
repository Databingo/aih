package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"runtime"
//	"strconv"

//	"github.com/eiannone/keyboard"
//	"github.com/nsf/termbox-go"
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

	key := gjson.Get(string(data), "key")
	OpenAI_Key := key.String()
	//fmt.Println(OpenAI_Key)

	proxy := gjson.Get(string(data), "proxy")
	Proxy := proxy.String()
	//fmt.Println(Proxy)

	// Set up the OpenAI API client
	config := openai.DefaultConfig(OpenAI_Key)

	if Proxy != "" {
		proxyUrl, err := url.Parse(Proxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		config.HTTPClient = &http.Client{Transport: transport}
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

	//err = keyboard.Open()
	//defer keyboard.Close()

	//hist := make([]string, 0)
	////cur := ""
	//pos := -1

	for {
		fmt.Print(left_tokens, role, "> ")

	//	if keyboard.Wait() {
	//		_, key, _ := keyboard.GetKey()
	//		switch key {
	//		case keyboard.KeyArrowUp:
	//			if pos > -1 {
	//				//cur = hist[pos]
	//				pos--
	//			}

	//		case keyboard.KeyArrowDown:
	//			if pos < len(hist)-1 {
	//				//cur = hist[pos]
	//				pos++
	//			}
	//		}
	//	} else if scanner.Scan() {
                        scanner.Scan() 
			userInput := scanner.Text()

			// Parse the command line arguments to get the prompt
			//prompt := strings.Join(os.Args[1:], " ")
			//	prompt := userInput

			switch userInput {
			case "":
				continue
			case ".exit":
				fmt.Println("Byebye")
				return
				//	case ".history":
				//		for _, cmd := range history {
				//			fmt.Println(cmd)
				//		}
				//		continue
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
			case ".help":
				//fmt.Println(".info        Print the information")
				fmt.Println(".help        Show help")
				//fmt.Println(".key         Set key")
				fmt.Println(".proxy       Set proxy")
				fmt.Println(".new         New conversation")
				fmt.Println(".prompt      Role of Assistant for create precise prompt")
				fmt.Println(".speak       Voice speak context")
				fmt.Println(".quiet       Quiet not speak")
				fmt.Println(".clear       Clear screen")
				fmt.Println(".exit        Exit")
				fmt.Println("                 ")
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
			 	if runtime.GOOS != "windows"{
				 cmd = "clear"
				 clear := exec.Command(cmd)
				 clear.Stdout = os.Stdout
				 clear.Run()
				}
				 continue
			case ".promptor":
				messages = make([]openai.ChatCompletionMessage, 0)
				max_tokens = 4097
				used_tokens = 0
				left_tokens = max_tokens - used_tokens
				userInput = prompt_prove 
				role = ".promptor"
			// char, key, err := keyboard.GetKey()
			// if err != nil { panic(err) }
			// if key == keyboard.KeyArrowUp {
			//  if pos >0 {
			//   pos--
			//   userInput = hist[pos]
			//   fmt.Print(userInput)
			//  }} else if key == keyboard.KeyArrowDown{
			//   if pos < len(hist) -1 {
			//    pos++
			//    userInput = hist[pos]
			//    fmt.Print(userInput)
			//  }} else if key == keyboard.KeyEnter{
			//   fmt.Println()
			//   break
			//  } else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			//   if len(userInput) > 0{
			//    userInput = userInput[:len(userInput)-1]
			//    fmt.Print("\b \b")
			//   }
			//  }else {
			//   userInput += string(char)
			//   fmt.Print(string(char))
			//  }//}
			//// }()

		//	default:
		//		hist = append(hist, userInput)
		//		pos = len(hist) - 1
			}

			// Add input to hist
			//        hist = append(hist, userInput)
			//pos = len(hist) -1

			// char, key, err := keyboard.GetKey()
			// if err != nil { panic(err) }
			// if key == keyboard.KeyArrowUp {
			//  if pos >0 {
			//   pos--
			//   userInput = hist[pos]
			//   fmt.Print(userInput)
			//  }} else if key == keyboard.KeyArrowDown{
			//   if pos < len(hist) -1 {
			//    pos++
			//    userInput = hist[pos]
			//    fmt.Print(userInput)
			//  }} else if key == keyboard.KeyEnter{
			//   fmt.Println()
			//   break
			//  } else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			//   if len(userInput) > 0{
			//    userInput = userInput[:len(userInput)-1]
			//    fmt.Print("\b \b")
			//   }
			//  }else {
			//   userInput += string(char)
			//   fmt.Print(string(char))
			//  }//}
			//// }()

			// Porcess input
			//prompt := userInput
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			})

			// Generate a response from ChatGPT
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					//	Messages: []openai.ChatCompletionMessage{
					//		{
					//			Role:    openai.ChatMessageRoleUser,
					//			Content: prompt,
					//		}},
					Messages: messages,
					//MaxTokens: 4096,
				},
			)

			if err != nil {
				fmt.Println(err)
				//return
				continue
			}

			// Print the response to the terminal
			//c := color.New(color.FgWhite, color.Bold)
			c := color.New(color.FgWhite)
			cnt := strings.TrimSpace(resp.Choices[0].Message.Content)
			used_tokens = resp.Usage.TotalTokens
			left_tokens = max_tokens - used_tokens
			c.Println(cnt)
			//fmt.Printf("%+v\n", left_tokens)

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
