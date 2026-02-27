package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

// func foo(){}
// foo()
// NOTE 💡: If the above code is uncommented, the following Go syntax error will be raised when attempting to
// run this file: `syntax error: non-declaration statement outside function body`
// Even though the error that is actually being raised in the Go parses is `expected declaration, found foo`
// It will be good to have the same type of error handling in Nilan.
func main() {
	subcommands.Register(&emitBytecodeCmd{}, "compiler")
	subcommands.Register(&replCompiledCmd{}, "compiler")
	subcommands.Register(&runCompiledCmd{}, "compiler")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

}
