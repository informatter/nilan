package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(&runCmd{}, "tree-walk-interpreter")
	subcommands.Register(&replCmd{}, "tree-walk-interpreter")
	subcommands.Register(&emitBytecodeCmd{}, "compiler")
	subcommands.Register(&replCompiledCmd{}, "compiler")
	subcommands.Register(&runCompiledCmd{}, "compiler")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
