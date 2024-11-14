package main

import (
	"bufio"
	"fmt"
	"io"
	"nilan/lexer"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nHi %s Welcome to Nilan!\n\n", user.Username)
	repl(os.Stdin, os.Stdout)
}

func repl(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, ">>> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "exit" {
			os.Exit(0)
		}
		lex := lexer.CreateLexer(line)
		tokens, err := lex.Scan()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(tokens)
	}
}
