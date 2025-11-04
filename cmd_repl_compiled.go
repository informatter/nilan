package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"nilan/compiler"
	"nilan/lexer"
	"nilan/vm"

	"github.com/google/subcommands"
)

type replCompiledCmd struct {
	diassemble   bool
	dumpBytecode bool
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
	f.BoolVar(&cmd.diassemble, "di", false, "Shorthand for diassemble.")
	f.BoolVar(&cmd.dumpBytecode, "du", false, "Shorthand for dumpBytecode")
}

func (cmd *replCompiledCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	fmt.Println("\n\nWelcome to the compiled version of Nilan!")

	scanner := bufio.NewScanner(os.Stdin)

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
		compiler := compiler.New(tokens)
		bytecode, err := compiler.Compile()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if cmd.diassemble {
			_, err := compiler.DiassembleBytecode(true, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 Bytecode diassemble error:\n:\t%s", err.Error())
				continue
			}

		}
		if cmd.dumpBytecode {
			err := compiler.DumpBytecode("")
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 Dump bytecode error:\n:\t%s", err.Error())
			}
		}
		vm := vm.New()
		runtimeErr := vm.Run(bytecode)
		if runtimeErr != nil {
			fmt.Fprintln(os.Stderr, runtimeErr.Error())
			continue
		}
	}
}
