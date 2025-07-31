
# Nilan

Nilan is a programming language I am currently developing for fun 🚀, implemented in Go.
My goal is to learn more about how programming languages work under the hood and to explore the different pipelines involved — from taking source code as input to making the CPU execute instructions 🤖.

## Features

✅ Write the press release

✅ Arithmetic expressions: `+`, `-`, `*`, `/`

✅ Comparison operators: `>`, `>=`, `<`, `<=`, `==`, `!=`

✅ Boolean literals: `true`, `false`

✅ `null` literal

✅ Parenthesized expressions

✅ Unary operations: logical not `!`, negation `-`

✅ REPL (Read-Eval-Print Loop) for interactive testing



## Limitations

Currently, Nilan supports only very primitive expressions and literals. The following are **not supported** yet:

🔴 Variable identifiers and names

🔴 Assignment statements (e.g., `a = 2`)

🔴 String literals and operations

🔴 Functions and function calls

🔴 Arrays or other complex data structures

🔴 Control flow constructs (e.g., `if`, loops)

🔴 Exponentiation or other advanced operators


## Quick Start

To start a REPL and try out Nilan expressions:

```bash
go run .
```

Start typing Nilan expressions in the interactive prompt.

## Current Grammar (ISO EBNF)

Nilan’s grammar is defined using **ISO Extended Backus–Naur Form (ISO EBNF)**, conforming to [ISO/IEC 14977](https://www.iso.org/standard/26153.html).

```ebnf
equality   = comparison , { ( '!=' | '==' ) , comparison } ;
comparison = term , { ( '>' | '>=' | '<' | '<=' ) , term } ;
term       = factor , { ( '+' | '-' ) , factor } ;
factor     = unary , { ( '*' | '/' ) , unary } ;
unary      = ( '!' | '-' ) , unary | primary ;
primary    = ( 'FLOAT' | 'INT' | 'true' | 'false' | 'null' ) | '(' , expression , ')' ;
```

> 💡 Currently, the grammar and parser support only basic constructs such as logical, arithmetic, unary operations, literals, and parenthesized expressions.

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

### Example Rule – Breakdown

Example rule:

```ebnf
term = factor , { ( '+' | '-' ) , factor } ;
```

Means:

- A `term` consists of:
    - A `factor`, followed by
    - Zero or more repetitions of:
        - Either `'+'` or `'-'`, and
        - Another `factor`


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

#### `term`

Handles addition and subtraction:

```text
+ , -
```

Example:

```text
3 + 5 - 2
```


#### `factor`

Handles multiplication and division:

```text
* , /
```

Example:

```text
4 * 2 / 8
```


#### `unary`

Handles unary operations like logical not and negation, with recursive chaining:

```text
! , -
```

Examples:

```text
--5
!(-3)
```


#### `primary`

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
foo + 3       # Variables/identifiers not supported yet
"abc" + "def" # String literals not supported
2 ** 3        # Exponentiation operator not supported
1 +           # Trailing operator causes parse error
```


## Extending Nilan

For example, too add new features like variables or function calls:

1. **Update the Lexer**
Add recognition of new token types such as identifiers.
2. **Extend the Grammar**
For variables, extend the grammar to include an `IDENTIFIER` token in `primary`:

```ebnf
primary = 'IDENTIFIER' | ( 'FLOAT' | 'INT' | 'true' | 'false' | 'null' ) | '(' , expression , ')' ;
```

3. **Implement Semantics**
TODO: Add section when compiler or interpreter is implemented

## Development

### Running the REPL

```bash
go run .
```


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

