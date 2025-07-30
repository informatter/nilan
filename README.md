# Nilan

Nilan is a programming language I am currently developing for fun 🚀 and I am implementing it in Go 

My goal is to learn more about how a programming language works under the hood and further explore the different pipelines involved  to take source code as input and make the CPU execute the instructions 🤖

### References

- [Crafting Interpreters](https://craftinginterpreters.com/), by Robert Nystrom
- [Extended Backus-Naur Form (EBNF)](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
- [Recursive Descent Parser](https://en.wikipedia.org/wiki/Recursive_descent_parser)


## Current Grammar (ISO EBNF)


### Grammar Definition

Uses **ISO EBNF syntax**, conforming to [ISO/IEC 14977](https://www.iso.org/standard/26153.html)

```ebnf
equality    = comparison , { ( '!=' | '==' ) , comparison } ;
comparison  = term , { ( '>' | '>=' | '<' | '<=' ) , term } ;
term        = factor , { ( '+' | '-' ) , factor } ;
factor      = unary , { ( '*' | '/' ) , unary } ;
unary       = ( '!' | '-' ) , unary | primary ;
primary     = ( 'FLOAT' | 'INT' | 'true' | 'false' | 'null' ) | '(' , expression , ')' ;
```

💡 Currently, the grammar and thus the parser itself only supports very primitive constructs such as logical, arithmetic, unary, literals and parenthesized expressions 




### How to Read the Grammar

Each line in the grammar defines a **production rule** of the following form:

```ebnf
nonterminal = definition ;
```

- A **nonterminal** (e.g. `term`, `factor`) is a named structure that refers to patterns in the language.
- A **definition** is made up of terminal symbols, other nonterminals, and structural operators.

### Terminals and Nonterminals

| **Type**       | **Example**       | **Description**                                  |
|----------------|-------------------|--------------------------------------------------|
| Nonterminal    | `term`            | A named part of the grammar that expands into other rules. |
| Terminal       | `'+'`, `'true'`   | A fixed literal token, enclosed in single quotes `'...'`. |

💡 **Note:** Tokens like `'INT'`, `'FLOAT'`, `'true'` refer to token types returned by the **lexer**, not literal character strings.

###  Grammar Notation Symbols

<table>
  <thead>
    <tr>
      <th><strong>Symbol</strong></th>
      <th><strong>Meaning</strong></th>
      <th><strong>Example</strong></th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>=</code></td>
      <td>Rule definition</td>
      <td><code>term = factor , { ('+' | '-') , factor } ;</code></td>
    </tr>
    <tr>
      <td><code>;</code></td>
      <td>End of rule</td>
      <td>Every rule ends in a semicolon</td>
    </tr>
    <tr>
      <td><code>,</code></td>
      <td>Sequence (then)</td>
      <td><code>a , b</code> means <code>a</code> followed by <code>b</code></td>
    </tr>
    <tr>
      <td><code>|</code></td>
      <td>Alternatives</td>
      <td><code>a | b</code> means either <code>a</code> or <code>b</code></td>
    </tr>
    <tr>
      <td><code>{ ... }</code></td>
      <td>Zero or more repetitions</td>
      <td><code>{ a }</code> means repeat <code>a</code> zero or more times</td>
    </tr>
    <tr>
      <td><code>( ... )</code></td>
      <td>Grouping</td>
      <td>Used to group alternatives or sequences</td>
    </tr>
  </tbody>
</table>



### Example Rule — Breakdown

Here’s an example rule from the grammar:

```ebnf
term = factor , { ( '+' | '-' ) , factor } ;
```

#### This means:

- A `term` consists of:
  - A `factor`, 
  - Followed by **zero or more** repetitions of:
    - A `'+'` or `'-'`,
    - Followed by another `factor`.

#### Example Matches:

- `3`
- `3 + 5`
- `3 - 4 + 2`

### Operator Precedence (Implicit)

The grammar structure also defines **operator precedence** without needing extra tables:



| **Precedence** (Low → High) | **Operators**            | **Rule**      |
|-----------------------------|--------------------------|---------------|
| Lowest                      | Equality: `==`, `!=`      | `equality`    |
|                             | Comparison: `>`, `<`, etc | `comparison`  |
|                             | Additive: `+`, `-`        | `term`        |
|                             | Multiplicative: `*`, `/`  | `factor`      |
|                             | Unary: `-`, `!`           | `unary`       |
| Highest                     | `()`, literals            | `primary`     |


💡Each level of the grammar (from `equality` up to `primary`) ensures that lower-precedence operations are parsed before the highest-precedence ones. lower-precedence expressions are parsed first because they might contain sub-expressions of higher-presedence.




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

Handles unary operations (negation, logical not):
```text
! , -
```
Recursive form allows chaining:
```text
--5
!(-3)
```

#### `primary`

Handles literal values and parenthesized expressions:
```text
(FLOAT | INT | true | false | null | '(' expression ')')
```
Example:
```text
(5 + 3)
true
```

### Example Expression

```text
1 + 2 * 3
```

#### Parsing According to Grammar

Using the grammar's precedence (from high to low):

1. **Multiplication (`*`)** → handled by `factor`
2. **Addition (`+`)** → handled by `term`

So:

- `2 * 3` is evaluated **first**, then
- `1 + (result of 2 * 3)`

#### AST Structure

the expression:

```
1 + 2 * 3
```

is parsed as:

```plaintext
   +
  / \
 1   *
    / \
   2   3
```

Expressed as a structured AST:

```python
Binary(
  op='+',
  left=Literal(1),
  right=Binary(
    op='*',
    left=Literal(2),
    right=Literal(3)
  )
)
```

#### Grammar Rule Involvement

| Expression | Grammar Rule      |
|------------|-------------------|
| `1`        | `primary` → `INT` |
| `2 * 3`    | `factor` (multiplication) |
| `1 + (...)`| `term` (addition) |



#### Summary

- The **AST mirrors operator precedence** encoded in the grammar.
- `2 * 3` is deeper in the tree than `1 + (...)`, reflecting that `*` has higher precedence than `+`.
- This method works for any expression the grammar can parse.


---

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
