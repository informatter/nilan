package main

import (
	"bufio"
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

type replCompiledCmd struct {
	diassemble   bool
	dumpBytecode bool
	dumpAST      bool
}

func (*replCompiledCmd) Name() string { return "cRepl" }
func (*replCompiledCmd) Synopsis() string {
	return "Start REPL session with the compiled version of nilan"
}
func (*replCompiledCmd) Usage() string {
	return `nilan cRepl`
}

func (cmd *replCompiledCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.diassemble, "diassemble", false, "diassemble the bytecode and dump it to a .dnic file")
	f.BoolVar(&cmd.dumpBytecode, "dumpBytecode", false, "Writes the encoded bytecode as hexadecimal to a .nic file")
	f.BoolVar(&cmd.dumpAST, "dumpAST", false, "Writes the AST as JSON to a file")
	f.BoolVar(&cmd.diassemble, "di", false, "Shorthand for diassemble.")
	f.BoolVar(&cmd.dumpBytecode, "du", false, "Shorthand for dumpBytecode")
	f.BoolVar(&cmd.dumpAST, "da", false, "Shorthand for dumpAST.")

}

func (cmd *replCompiledCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	fmt.Println("\nWelcome to the Nilan programming language!")
	fmt.Println("")

	fmt.Print(`
	███╗   ██╗██╗██╗      █████╗ ███╗   ██╗    ██████╗ ███████╗██████╗ ██╗     
	████╗  ██║██║██║     ██╔══██╗████╗  ██║    ██╔══██╗██╔════╝██╔══██╗██║     
	██╔██╗ ██║██║██║     ███████║██╔██╗ ██║    ██████╔╝█████╗  ██████╔╝██║     
	██║╚██╗██║██║██║     ██╔══██║██║╚██╗██║    ██╔══██╗██╔══╝  ██╔═══╝ ██║     
	██║ ╚████║██║███████╗██║  ██║██║ ╚████║    ██║  ██║███████╗██║     ███████╗
	╚═╝  ╚═══╝╚═╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝    ╚═╝  ╚═╝╚══════╝╚═╝     ╚══════╝
																			
`)
	scanner := bufio.NewScanner(os.Stdin)
	astCompiler := compiler.NewASTCompiler()
	vm := vm.New()

	for {
		fmt.Fprintf(os.Stdout, ">>> ")
		scanned := scanner.Scan()
		if !scanned {
			err := scanner.Err()
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 %s", err.Error())
				return subcommands.ExitFailure
			}
		}

		line := scanner.Text()
		if line == "exit" {
			os.Exit(0)
		}

		lex := lexer.New(line)
		tokens, err := lex.Scan()
		if err != nil {
			fmt.Println(err)
			continue
		}

		parser := parser.Make(tokens)
		statements, parseErrs := parser.Parse()
		if len(parseErrs) > 0 {
			fmt.Fprintf(os.Stdout, "Parse error: ")
			for _, pErr := range parseErrs {
				fmt.Fprintf(os.Stdout, "%v\n", pErr)
			}
			continue
		}

		// TODO/NOTE: Previous compiled code is going to be recompiled again in the REPL,
		// but for now its fine
		bytecode, err := astCompiler.CompileAST(statements)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		if cmd.diassemble {
			_, err := astCompiler.DiassembleBytecode(true, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 Bytecode diassemble error:\n:\t%s", err.Error())
				continue
			}

		}
		if cmd.dumpBytecode {
			err := astCompiler.DumpBytecode("")
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 Dump bytecode error:\n:\t%s", err.Error())
			}
		}
		if cmd.dumpAST {
			err := parser.PrintToFile(statements, "ast.json")
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 Dump AST error:\n:\t%s", err.Error())
				continue
			}
		}

		runtimeErr := vm.Run(bytecode)
		if runtimeErr != nil {
			fmt.Fprintln(os.Stderr, runtimeErr.Error())
			continue
		}
	}
}
