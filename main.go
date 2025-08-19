package main

import (
	"bufio"
	"fmt"
	"io"
	"nilan/interpreter"
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
		parser := parser.Create(tokens)
		ast, errors := parser.Parse()
		if len(errors) > 0 {
			continue
		}
		parser.Print(ast)
		interpreter := interpreter.Interpreter{}
		interpreter.Interpret(ast)

	}
}
