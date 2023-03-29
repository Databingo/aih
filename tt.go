package main

import (
	//     "bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	//		         "strings"

	"github.com/peterh/liner"
)

func main() {
	// Define command line flags
	model := flag.String("model", "gpt2", "ChatGPT model to use")
	length := flag.Int("length", 20, "Length of generated text")
	prompt := flag.String("prompt", "", "Prompt for generated text")

	// Parse command line arguments
	flag.Parse()

	// Create a new liner instance for command history and line editing
	liner := liner.NewLiner()
	defer liner.Close()

	// Load command history from file
	if f, err := os.Open(".history"); err == nil {
		liner.ReadHistory(f)
		f.Close()
	}

	// Loop to read and execute commands
	for {
		// Read a line of input
		if input, err := liner.Prompt("> "); err == nil {
			// Add input to command history
			liner.AppendHistory(input)

			// If input is "exit", break the loop
			if input == "exit" {
				break
			}

			// Create a new chatgpt command with the specified model and length
			cmd := exec.Command("chatgpt", "--model", *model, "--length", fmt.Sprintf("%d", *length))

			// Connect to the command's standard input and output
			stdin, _ := cmd.StdinPipe()
			//stdout, _ := cmd.StdoutPipe()

			// Start the command
			cmd.Start()

			// Write the input to the command's standard input
			stdin.Write([]byte(*prompt + input + "\n"))

		}
	}
}
