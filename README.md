# Nilan

Nilan is a programming language I am currently developing for fun 🚀 and I am implementing it in Go 

My goal is to learn more about how a programming language works under the hood and further explore the different pipelines involved  to take source code as input and make the CPU execute the instructions 🤖

### References

- [Crafting Interpreters](https://craftinginterpreters.com/), by Robert Nystrom
- [Extended Backus-Naur Form (EBNF)](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
- [Recursive Descent Parser](https://en.wikipedia.org/wiki/Recursive_descent_parser)

## Development

To start a REPL:
```bash
go run .
```

Start typing nilan code

### Testing

To test a particular package:

```bash
go test ./lexer
```

To run all unit tests:

```bash
go test ./...
```

### Linting

**Formatting**

To format a particular package:

```bash
go fmt ./lexer
```

To format all files:

```bash
go fmt ./...
```
