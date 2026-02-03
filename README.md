
# Nilan

Nilan is a programming language I am currently developing for fun 🚀, implemented in Go.
My goal is to learn more about how programming languages work under the hood and to explore the different pipelines involved — from taking source code as input to making the CPU execute instructions 🤖.

ℹ️ The project has transitioned from a tree-walk interpreter to a compiler-based architecture. The tree-walk interpreter is now **deprecated**. The current development focuses on the **ASTCompiler** which compiles AST nodes to bytecode executed by a stack-based **Virtual Machine (VM)**.

## Architecture

The language now follows a traditional compiler pipeline:

```
Source Code → Lexer → Tokens → Parser → AST → ASTCompiler → Bytecode → VM
```

## Features

### ASTCompiler + VM (Current) ✅

✅ Arithmetic expressions: `+`, `-`, `*`, `/`

✅ Comparison operators: `>`, `>=`, `<`, `<=`, `==`, `!=`

✅ Boolean literals: `true`, `false`

✅ Unary negation: `-10`, `!false`

✅ Assignment statements (e.g., `var a = 2`)

✅ Literal values: integers, floats, boleans, strings

✅ Grouped expressions: `(a + b) * c`

✅ REPL (Read-Eval-Print Loop) for interactive testing

✅ Execute source code from a file (via `emit` command)

### Tree-Walk Interpreter (Deprecated) ⚠️

The following features were implemented in the tree-walk interpreter but are **not yet** supported in the ASTCompiler + VM:

✅ Lexical scope

✅ Block scope: `{}`

✅ Comparison operators: `>`, `>=`, `<`, `<=`, `==`, `!=`

✅ Boolean literals: `true`, `false`

✅ String literals: `"hellow world"`

✅ Control flow: `if`, `else`

✅ Logical operators: `and`, `or`

✅ `null` literal

✅ Parenthesized expressions

✅ Variable identifiers and names

✅ Assignment statements (e.g., `var a = 2`)

✅ Unary operations: logical not `!`

## Limitations

### ASTCompiler + VM (Current) 🔴

The following features are **not yet supported** in the compiled version:

🔴 Lexical and block scope

🔴 Boolean literals and operations: `and`, `or`

🔴 String literals and string operations

🔴 Control flow: `if`, `else`, `while` loops

🔴 Functions and function calls

🔴 Classes, structs, interfaces

🔴 Exponentiation or other advanced operators

🔴 Complex features such as Module/package imports, data structures, etc ...

🔴 Static typing

### Tree-Walk Interpreter (Deprecated) 🔴

The following are **not supported** in the tree-walk interpreter (and are not planned):

🔴 Functions and function calls

🔴 Classes, structs, interfaces

🔴 Inheritance

🔴 Arrays or other complex data structures

🔴 Control flow: loops, `else if`, `break`

🔴 Exponentiation or other advanced operators

🔴 Complex features such as Module/package imports, etc ...




## Current Syntactic Grammar (ISO EBNF)

Nilan’s syntactic grammar is defined using **ISO Extended Backus–Naur Form (ISO EBNF)**, conforming to [ISO/IEC 14977](https://www.iso.org/standard/26153.html). It represents the rules used to parse a sequence of tokens into an Abstract Syntax Tree (AST)

```ebnf
program = { declaration }, EOF ;

declaration = variable-declaration | statement ;

variable-declaration = identifier , "=" , expression ;

statement = expression
          | if-statement
          | print-statement
          | while-statement
          | block-statement ;

if-statement = "if" , expression , statement , [ "else" , statement ] ;

print-statement = "print" , expression ;

while-statement = "while", expression, statement ;

block-statement  = "{" , { declaration } , "}" ;

expression = assignment-expression ;

assignment-expression = IDENTIFIER, "=", assignment-expression
           | or-expression ;

or-expression  = and-expression , { "or" , and-expression } ;

and-expression = equality-expression, { "and", equality-expression } ;

equality-expression = comparison-expression, { ("!=", "=="), comparison-expression } ;

comparison-expression = term-expression, { (">" | ">=" | "<" | "<="), term-expression } ;

term-expression = factor-expression, { ("+" | "-"), factor-expression } ;

factor-expression = unary-expression, { ("*" | "/"), unary-expression } ;

unary-expression = ("!" | "-"), unary-expression
      | primary-expression ;

primary-expression = FLOAT 
        | INT 
        | IDENTIFIER 
        | "true" 
        | "false" 
        | "null" 
        | "(", expression, ")" ;

```

This grammar is not left-recursive because none of the non-terminals start their production with themselves on the left side. Each rule begins with a different non-terminal or terminal before any recursion happens. For example, `equality` starts with `comparison`,`comparison` starts with `term`, etc...


### How to Read the Grammar

Each line defines a **production rule** in the form:

```ebnf
nonterminal = definition ;
```

- A **nonterminal** (e.g., `term`, `factor`) is a named syntactic category made of other rules.
- A **definition** consists of terminal symbols (token literals), other nonterminals, and notation operators.


### Terminals and Nonterminals

| **Type** | **Example** | **Description** |
| :-- | :-- | :-- |
| Nonterminal | `term` | Named construct that expands into other rules |
| Terminal | `'+'`, `'true'` | Fixed token literals enclosed in single quotes |

> 💡 **Note:** Tokens like `'INT'` and `'FLOAT'` are token types returned by the lexer, not literal characters.

### Grammar Notation Symbols

| Symbol | Meaning | Example |
| :-- | :-- | :-- |
| `=` | Rule definition | `term = factor , { ('+' | '-') , factor } ;` |
| `;` | End of rule | Every rule ends in a semicolon |
| `,` | Sequence | `a , b` means `a` followed by `b` |
| `|` | Alternatives | `a | b` means either `a` or `b` |
| `{ ... }` | Zero or more repetitions | `{ a }` means repeat `a` zero or more times |
| `( ... )` | Grouping | Used to group alternatives or sequences |
| `[ ... ]` | Optional | Used to speficy optional implementation, for example an `else` clause |

### Example Rule – Breakdown

Example rule:

```ebnf
term-expression = factor-expression , { ( '+' | '-' ) , factor-expression } ;
```

Means:

- A `term-expression` consists of:
    - A `factor-expression`, followed by
    - Zero or more repetitions of:
        - Either `'+'` or `'-'`, and
        - Another `factor-expression`


#### Example Matches:

- `3`
- `3 + 5`
- `3 - 4 + 2`


### Operator Precedence (Implicitly Encoded)

Precedence from lowest to highest is encoded in the grammar structure itself:


| Precedence Level | Operators | Grammar Rule |
| :-- | :-- | :-- |
| Lowest | Equality: `==`, `!=` | `equality` |
|  | Comparison: `>`, `<`, etc | `comparison` |
|  | Additive: `+`, `-` | `term` |
|  | Multiplicative: `*`, `/` | `factor` |
|  | Unary: `-`, `!` | `unary` |
| Highest | Parentheses, literals | `primary` |

> 💡 Lower-precedence rules contain (as components) higher-precedence expressions. This structure ensures operators like `*` bind more tightly than `+`. For example, the expression `5 * 5 + 10 + 2` is parsed as `(5 * 5) + 10 + 2`.


### Some Examples

#### `term-expression`

Handles addition and subtraction:

```text
+ , -
```

Example:

```text
3 + 5 - 2
```


#### `factor-expression`

Handles multiplication and division:

```text
* , /
```

Example:

```text
4 * 2 / 8
```


#### `unary-expression`

Handles unary operations like logical not and negation, with recursive chaining:

```text
! , -
```

Examples:

```text
--5
!(-3)
```


#### `primary-expression`

Handles literals and parenthesized expressions:

```text
(FLOAT | INT | true | false | null | '(' expression ')')
```

Examples:

```text
(5 + 3)
true
```


### Example: Parsing `1 + 2 * 3`

Parsing order according to precedence:

1. Multiplication `*` by `factor` rule
2. Addition `+` by `term` rule

Result:

- Multiply `2 * 3` **first**
- Add `1 + (2 * 3)`


### AST Structure

```
   +
  / \
 1   *
    / \
   2   3
```

Expressed as:

```python
Binary(
  Left=Literal(1),
  Op='+',
  Right=Binary(
    Left=Literal(2),
    Op='*',
    Right=Literal(3)
  )
)
```


### Grammar Rule Involvement

| Expression | Grammar Rule |
| :-- | :-- |
| `1` | `primary` → `INT` |
| `2 * 3` | `factor` (multiplication) |
| `1 + (...)` | `term` (addition) |

### Invalid or Unsupported Examples

These examples will **not parse** correctly with the current grammar:

```text
1 +
2++
2--
2+=
2-=
2**2
```


## Extending Nilan


1. **Update the Lexer (Optional)**

    If new token types need to be introduced, `token.go` and `lexer.go` need to be modified. 

2. **Extend the `ExpressionVisitor` or `StmtVisitor` interfaces**

    Depending on the type of new syntax introduced, make sure to add the corresponding visit method to one of the interfaces.

3. **Add a new AST node to `expressions.go` or `statements.go`**

    Depending on the type of new syntax introduced, make sure to add the corresponding AST node `struct` to `expressions.go` or `statements.go` depending if its an expression or statement node.

4. **Extend the Parser**

    Extend the Parser to handle the new syntax grammar by adding a method which creates an AST node.

6. **Extend the compiler**

    a. Add a new opcode (`code.go`) and handle diassembling and assembling the opcode

    b. Extend the compiler to compile the new AST node(s) to bytecode

7. **Extend the VM**

    Extend the vm to execute the new bytecode instruction(s)

8. **Extend the AST Printer (Optional)**

    Extend the AST Printer by adding the corresponding visitor method to handle the new AST node.


## Installation

```bash
git clone https://github.com/informatter/nilan.git
cd nilan
go install .
```

> 💡 If changes are made to the code, run `go install .` again to create a new binary with the updates.

## Usage

### Compiled Version (ASTCompiler + VM)

**1. Emit Bytecode**

Generates and optionally disassembles bytecode from a Nilan source file. Useful for debugging the compiler.

```bash
nilan emit arithmetic.ni
```

```bash
nilan emit --help
```

**2. REPL**

Interactive REPL session for testing arithmetic expressions:

```bash
nilan cRepl
```

To see all available commands:

```bash
nilan cRepl --help
```

> 💡 For iterative development, use: `go run . -- cRepl` or `go run . -- emit <file-name>` ... etc so any CLI tool can be used without needing to build a binary.


### Tree-Walk Interpreter (Deprecated)

The following commands are still available but use the deprecated tree-walk interpreter:

**1. Run**

Executes a Nilan source file:
```bash
nilan run hellow_world.ni
```

**2. REPL**

Interactive REPL session:
```bash
nilan repl
```


### Testing

Run tests for a specific package, e.g., lexer:

```bash
go test ./compiler
```

Run all unit tests recursively:

```bash
go test ./...
```

Run a particular integration test:

```bash
go test ./compiler -v -run TestFullPipeline
```


### Linting and Formatting

Format a particular package:

```bash
go fmt ./compiler
```

Format all Go files:

```bash
go fmt ./...
```

There are also handy aliases in `.aliases.sh` for running all the test suite, formatting and building the binary

```bash
source .aliases
format
build
test
repl ## NOTE: Uses the compiled version of Nilan by default
```

### Using the debugger

💡The below instructions are specifically for VSCode

A debugger can be used if desired by installing the Delve debugger for the golang programming language:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Then create a `.vscode` folder and add the following `launch.json` file:

```json

{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Nilan Code Execution (dlv-dap)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "debugAdapter": "dlv-dap",
            "console": "integratedTerminal",
            // NOTE: Change file if needed.
            "args": ["run","hellow_world.ni"]
        },
        {
            "name": "Launch Nilan REPL (dlv-dap)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "debugAdapter": "dlv-dap",
            "console": "integratedTerminal",
            "args": ["repl"]
        },
        {
            "name": "Launch Nilan Compile Emit (dlv-dap)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "debugAdapter": "dlv-dap",
            "console": "integratedTerminal",
            "args": ["emit","arithmetic.ni"]
        },
        {
            "name": "Launch Nilan Compile REPL (dlv-dap)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "debugAdapter": "dlv-dap",
            "console": "integratedTerminal",
            "args": ["cRepl"]
        }
    ]
}

```
The launch configuration can be modified as needed, this is just my current setup.

After this the debugger will be ready to be used and breakpoints can be set.




## References

- [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom
- [Extended Backus-Naur Form (EBNF)](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
- [Recursive Descent Parser](https://en.wikipedia.org/wiki/Recursive_descent_parser)
- [Pratt Parser](https://www.chidiwilliams.com/posts/on-recursive-descent-and-pratt-parsing)
- [Visitor Pattern (Go example)](https://refactoring.guru/design-patterns/visitor/go/example)
- [Operator Precedence Explained](https://dear-computer.twodee.org/expressions/precedence.html)

