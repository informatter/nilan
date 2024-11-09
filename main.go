package main

import (
	"bufio"
	"fmt"
	"nilan/lexer"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\nNilan\n\n")
	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		cleanedInput := strings.TrimSpace(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if cleanedInput == "exit" {
			os.Exit(0)
		}

		scanner := lexer.Lexer{
			Input: cleanedInput,
		}

		tokens, err := scanner.Scan()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(tokens)
	}
}
