package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"nilan/compiler"
	"nilan/lexer"
	"nilan/parser"
	"nilan/token"
	"nilan/vm"

	"github.com/chzyer/readline"
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
	rl, err := readline.New(">>> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "💥 Failed to initialize REPL: %s\n", err.Error())
		return subcommands.ExitFailure
	}
	defer rl.Close()

	astCompiler := compiler.NewASTCompiler()
	vm := vm.New()
	var buffer strings.Builder
	// persistent src for entire REPL session
	var sessionSource strings.Builder
	for {
		var prompt string
		if buffer.Len() == 0 {
			prompt = ">>> "
		} else {
			prompt = "... "
		}
		rl.SetPrompt(prompt)

		line, err := rl.Readline()
		if err != nil {
			if err.Error() == "EOF" {
				return subcommands.ExitSuccess
			}
			fmt.Fprintf(os.Stderr, "💥 %s\n", err.Error())
			return subcommands.ExitFailure
		}
		if strings.TrimSpace(line) == "exit" && buffer.Len() == 0 {
			os.Exit(0)
		}

		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)
		currentChunk := buffer.String()
		fullSource := currentChunk
		if sessionSource.Len() > 0 {
			fullSource = sessionSource.String() + "\n" + currentChunk
		}

		lex := lexer.New(fullSource)
		tokens, err := lex.Scan()
		if err != nil {
			fmt.Println(err)
			buffer.Reset()
			continue
		}

		if !isInputReady(tokens) {
			continue
		}

		parser := parser.Make(tokens)
		// `true` indicates that the input is from a REPL session, allowing expressions at the top level.
		statements, parseErrs := parser.Parse(true)
		if len(parseErrs) > 0 {
			// If all parse errors are syntax errors that occur at the position of the EOF token,
			// it means that the user has not finished typing their input yet.
			// We should wait for more input instead of showing an error.
			if allParseErrorsAtEOF(parseErrs, tokens[len(tokens)-1]) {
				continue
			}
			fmt.Fprintf(os.Stdout, "Parse error: ")
			for _, pErr := range parseErrs {
				fmt.Fprintf(os.Stdout, "%v\n", pErr)
			}
			buffer.Reset()
			continue
		}

		// Commit current chunk into session source after successful parse.
		if sessionSource.Len() > 0 {
			sessionSource.WriteString("\n")
		}
		sessionSource.WriteString(currentChunk)

		if cmd.dumpAST {
			err := parser.PrintToFile(statements, "ast.json")
			if err != nil {
				fmt.Fprintf(os.Stderr, "💥 Dump AST error:\n:\t%s", err.Error())
				continue
			}
		}
		// TODO/NOTE: Previous compiled code is going to be recompiled again in the REPL,
		// but for now its fine
		bytecode, err := astCompiler.CompileAST(statements)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			buffer.Reset()
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

		runtimeErr := vm.Run(bytecode)
		if runtimeErr != nil {
			fmt.Fprintln(os.Stderr, runtimeErr.Error())
			buffer.Reset()
			continue
		}
		buffer.Reset()
	}
}

// isInputReady checks if the input is ready to be parsed and executed. It checks for balanced parentheses and braces,
// and also checks if the last non-EOF token is an operator or a keyword that expects more input.
//
// For example, if the user types `if (x > 5) {`, the REPL should wait for more input until the
// user finishes the block with a `}`.
func isInputReady(tokens []token.Token) bool {

	braceBalance := 0
	for _, tok := range tokens {
		switch tok.TokenType {
		case token.LCUR:
			braceBalance++
		case token.RCUR:
			braceBalance--
		}
	}

	if braceBalance > 0 {
		return false
	}

	last := lastNonEOF(tokens)
	if last == nil {
		return true
	}

	switch last.TokenType {
	case token.ASSIGN,
		token.ADD,
		token.SUB,
		token.MULT,
		token.DIV,
		token.BANG,
		token.EQUAL_EQUAL,
		token.NOT_EQUAL,
		token.LESS,
		token.LESS_EQUAL,
		token.LARGER,
		token.LARGER_EQUAL,
		token.COMMA,
		token.LPA,
		token.LCUR,
		token.IF,
		token.ELSE,
		token.ELIF,
		token.WHILE,
		token.FOR,
		token.FUNC,
		token.RETURN,
		token.VAR,
		token.CONST,
		token.AND,
		token.OR,
		token.PRINT:
		return false
	}

	return true
}

// lastNonEOF returns the last non-EOF token from the list of tokens. If all tokens are EOF, it returns nil.
func lastNonEOF(tokens []token.Token) *token.Token {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].TokenType != token.EOF {
			return &tokens[i]
		}
	}
	return nil
}

// allParseErrorsAtEOF checks if all parse errors are syntax errors that occur at the position of the EOF token.
func allParseErrorsAtEOF(parseErrs []error, eof token.Token) bool {
	for _, parseErr := range parseErrs {
		syntaxErr, ok := parseErr.(parser.SyntaxError)
		if !ok {
			return false
		}
		if syntaxErr.Line != eof.Line || syntaxErr.Column != eof.Column {
			return false
		}
	}
	return len(parseErrs) > 0
}
