package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"nilan/compiler"
	"nilan/lexer"
	"nilan/parser"
	"nilan/vm"

	"github.com/google/subcommands"
)

// replCmd implements the REPL command
type runCompiledCmd struct{}

func (*runCompiledCmd) Name() string     { return "runC" }
func (*runCompiledCmd) Synopsis() string { return "Execute Nilan code from a source file" }
func (*runCompiledCmd) Usage() string {
	return `run:
  Execute Nilan code.
`
}
func (r *runCompiledCmd) SetFlags(f *flag.FlagSet) {}

func (r *runCompiledCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	compiler := compiler.NewASTCompiler()
	vm := vm.New()
	lex := lexer.New(string(data))
	tokens, err := lex.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
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
	bytecode, err := compiler.CompileAST(ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	err = vm.Run(bytecode)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
