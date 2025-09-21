package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"nilan/interpreter"
	"nilan/lexer"
	"nilan/parser"
)

// replCmd implements the REPL command
type runCmd struct{}

func (*runCmd) Name() string     { return "run" }
func (*runCmd) Synopsis() string { return "Execute Nilan code from a source file" }
func (*runCmd) Usage() string {
	return `run:
  Execute Nilan code.
`
}
func (r *runCmd) SetFlags(f *flag.FlagSet) {}

func (r *runCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "💥 File not provided\n")
		return subcommands.ExitUsageError
	}
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "💥 Failed to read file: %v\n", err)
		return subcommands.ExitFailure
	}

	interpreter := interpreter.Make()
	lex := lexer.CreateLexer(string(data))
	tokens, err := lex.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lexing error: %v\n", err)
		return subcommands.ExitFailure
	}
	parser := parser.Make(tokens)
	ast, errors := parser.Parse()
	if len(errors) > 0 {
		for _, error := range errors {
			fmt.Fprintln(os.Stderr, error)
		}
		return subcommands.ExitFailure
	}
	interpreter.Interpret(ast)
	return subcommands.ExitSuccess
}
