## Development

To start a REPL:
```bash
go run .
```

Start typing nilan code

### Testing

To test a particular package:

```bash
go fmt ./lexer
```

To run all unit tests:

```bash
go test ./...
```

**Note:** Go's testing framework automatically runs files ending with `_test.go` and functions starting with `Test`, `Benchmark`, or `Example`

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
