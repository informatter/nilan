
# Nilan

Nilan is a programming language I am currently developing for fun 🚀, implemented in Go.
My goal is to learn more about how programming languages work under the hood and to explore the different pipelines involved — from taking source code as input to making the CPU execute instructions 🤖.


## Features

✅ Arithmetic expressions: `+`, `-`, `*`, `/`

✅ Lexical scope

✅ Block scope: `{}`

✅ Comparison operators: `>`, `>=`, `<`, `<=`, `==`, `!=`

✅ Boolean literals: `true`, `false`

✅ String literals: `"hellow world"`

✅ Control flow: `if`, `else`

✅ Logical operators: `and`, `or`

✅ Boolean literals: `true`, `false`

✅ `null` literal

✅ Parenthesized expressions

✅ Variable identifiers and names

✅ Assignment statements (e.g., `var a = 2`)

✅ Unary operations: logical not `!`, negation `-`

✅ REPL (Read-Eval-Print Loop) for interactive testing

✅ Execute source code from a file.

## Limitations

The following are **not supported** yet:

🔴 Functions and function calls

🔴 Classes, structs, interfaces

🔴 Inheritance

🔴 Arrays or other complex data structures

🔴 Control flow: loops, `else if`, `break`

🔴 Logical operators: `not`

🔴 Exponentiation or other advanced operators

🔴 Tree-Walk interpreter

🔴 Complex features such as Module/package imports, etc ...

## Short term TODOs

- While loop
- For loop
- Functions

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


1. **Update the Lexer (optional)**
If new token types need to be introduced, `token.go` and `lexer.go` need to be modified. 

2. **Extend the `ExpressionVisitor` or `StmtVisitor` interfaces
Depending on the type of new syntax introduced, make sure to add the corresponding visit method to one of the interfaces.

3. **Add a new AST node to `expressions.go` or `statements.go`
Depending on the type of new syntax introduced, make sure to add the corresponding AST node `struct` to `expressions.go` or `statements.go` depending if its an expression or statement node.

4. **Extend the Parser**
Extend the Parser to handle the new syntax grammar by adding a method which creates an AST node.

5. **Extend the Interpreter**
Extend the interpreter to execute the the new AST node returned by the parser. This will involve implementing the method added to the `ExpressionVisitor` or `StmtVisitor` interfaces


## Installation

```bash
git clone https://github.com/informatter/nilan.git
cd nilan
go install .
```

## Usage

Once installed there are two main commands than can be used:

**1. REPL**

Start a REPL session
```bash
nilan repl
```

**2. Run**

Compiles the specified file and executes it directly
```bash
nilan run hellow_world.ni
```

💡If changes are made to the code, run `go install .` once again so a new binary is created with the new changes.

For iterative development is recommended to simply run:

`go run . -- repl` **or** `go run . -- run <file-name>`

for a more efficient workflow.



### Testing

Run tests for a specific package, e.g., lexer:

```bash
go test ./lexer
```

Run all unit tests recursively:

```bash
go test ./...
```


### Linting and Formatting

Format a particular package:

```bash
go fmt ./lexer
```

Format all Go files:

```bash
go fmt ./...
```


## References

- [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom
- [Extended Backus-Naur Form (EBNF)](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
- [Recursive Descent Parser](https://en.wikipedia.org/wiki/Recursive_descent_parser)
- [Visitor Pattern (Go example)](https://refactoring.guru/design-patterns/visitor/go/example)
- [Operator Precedence Explained](https://dear-computer.twodee.org/expressions/precedence.html)

