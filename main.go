package main

import (
	"bufio"
	"fmt"
	"io"
	"nilan/lexer"
	"nilan/parser"
	"os"
)

func main() {
	fmt.Println("\n\nWelcome to Nilan!")
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
		parser := parser.CreateParser(tokens)
		ast, _ := parser.Parse()
		fmt.Println(ast)
	}
}
