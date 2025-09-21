package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/google/subcommands"
	"nilan/interpreter"
	"nilan/lexer"
	"nilan/parser"
)

// replCmd implements the REPL command
type replCmd struct{}

func (*replCmd) Name() string     { return "repl" }
func (*replCmd) Synopsis() string { return "Start REPL session" }
func (*replCmd) Usage() string {
	return `repl:
  Start interactive REPL session.
`
}
func (r *replCmd) SetFlags(f *flag.FlagSet) {}

func repl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	interpreter := interpreter.Make()

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
		parser := parser.Make(tokens)
		ast, errors := parser.Parse()
		if len(errors) > 0 {
			for _, error := range errors {
				fmt.Fprintln(os.Stderr, error)
			}
			continue
		}
		parser.Print(ast)
		interpreter.Interpret(ast)

	}
}

func (r *replCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Println("\n\nWelcome to Nilan!")
	repl(os.Stdin, os.Stdout)
	return subcommands.ExitSuccess
}
