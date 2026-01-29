package main

import (
	"context"
	"flag"
	"fmt"
	"nilan/compiler"
	"nilan/lexer"
	"nilan/parser"
	"os"
	"strings"

	"github.com/google/subcommands"
)

type emitBytecodeCmd struct {
	diassemble   bool
	dumpBytecode bool
	filePath     string
}

func (*emitBytecodeCmd) Name() string { return "emit" }
func (*emitBytecodeCmd) Synopsis() string {
	return "Emit the bytecode representation from a source file"
}
func (*emitBytecodeCmd) Usage() string {
	return `nilan emit <file>`
}

func (cmd *emitBytecodeCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.diassemble, "diassemble", true, "diassemble the bytecode and dump it to a text file.")
	f.BoolVar(&cmd.dumpBytecode, "dumpBytecode", true, "Writes the encoded bytecode as hexadecimal to a .nic file")
	f.StringVar(&cmd.filePath, "file path", "/", "The file path to write the diassembled bytecode to. If no file path is provided the file will be saved under the same directory where this command is executed from.")
}

func (r *emitBytecodeCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "💥 File not provided\n")
		return subcommands.ExitUsageError
	}
	nilanFile := args[0]
	data, err := os.ReadFile(nilanFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "💥 Failed to read file:\n\t%v", err.Error())
		return subcommands.ExitFailure
	}

	lex := lexer.New(string(data))
	tokens, err := lex.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lexing error: %v\n", err)
		return subcommands.ExitFailure
	}

	parser := parser.Make(tokens)
	statements, parseErrs := parser.Parse()
	if len(parseErrs) > 0 {
		fmt.Fprintf(os.Stderr, "💥 Parsing error:\n")
		for _, pErr := range parseErrs {
			fmt.Fprintf(os.Stderr, "\t%v\n", pErr)
		}
		return subcommands.ExitFailure
	}

	astCompiler := compiler.NewASTCompiler()
	_, cErr := astCompiler.CompileAST(statements)

	if cErr != nil {
		fmt.Fprintf(os.Stderr, "💥 Compilation error:\n\t%v\n", cErr)
		return subcommands.ExitFailure
	}

	if r.diassemble {
		parts := strings.Split(nilanFile, ".")
		fileName := parts[0]
		astCompiler.DumpBytecode(fileName)

		_, dErr := astCompiler.DiassembleBytecode(true, fileName)
		if dErr != nil {
			fmt.Fprintf(os.Stderr, "💥 Bytecode diassemble error:\n:\t%s", dErr.Error())
			return subcommands.ExitFailure
		}
	}

	if r.dumpBytecode {
		parts := strings.Split(nilanFile, ".")
		fileName := parts[0]
		err := astCompiler.DumpBytecode(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "💥 Dump bytecode error:\n:\t%s", err.Error())
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess

}
