// This package contains the parser and compiler for Nilan. A Pratt parser is used to parse expressions,
// Each token maps to a particular infix and prefix parsing rule with its presedence level.
package compiler

import (
	"encoding/binary"
	"fmt"
	"nilan/token"
	"os"
	"strings"
)

// Precedence levels for the grammar's rules, ordered from lowest to highest.
// Highest rules will be parsed and compiled before lower presedence rules.
const (
	PREC_NONE = iota
	PREC_ASSIGNMENT
	PREC_TERM   // +,-
	PREC_FACTOR // /,*
	PREC_UNARY  // !, -,
)

var precedence = map[token.TokenType]int{
	token.ADD:  PREC_TERM,
	token.SUB:  PREC_TERM,
	token.DIV:  PREC_FACTOR,
	token.MULT: PREC_FACTOR,
	token.BANG: PREC_UNARY,
}

type ParseFunc func(*Compiler)

// Defines the parsing behavior for a specific token type.
// It contains optional prefix and infix parsing functions, and the precedence level of the token.
type parseRule struct {
	prefix     ParseFunc
	infix      ParseFunc
	precedence int
}

// Represents the compiler which will compile
// a stream of `Token`'s to `Bytecode` to be executed
// by the VM
type Compiler struct {
	bytecode     Bytecode
	readPosition int32

	totalTokens  int32
	tokens       []token.Token
	previousTok  token.Token
	currentTok   token.Token
	parsingRules map[token.TokenType]parseRule
}

// Creates a `Compiler` instance and returns
// a pointer to it.
func NewCompiler(tokens []token.Token) *Compiler {
	c := &Compiler{
		bytecode: Bytecode{
			Instructions:  Instructions{},
			ConstantsPool: []any{},
		},
		totalTokens: int32(len(tokens)),
		tokens:      tokens,

		parsingRules: map[token.TokenType]parseRule{
			token.ADD:   {prefix: nil, infix: (*Compiler).binary, precedence: PREC_TERM},
			token.SUB:   {prefix: (*Compiler).unary, infix: (*Compiler).binary, precedence: PREC_TERM},
			token.DIV:   {prefix: nil, infix: (*Compiler).binary, precedence: PREC_FACTOR},
			token.MULT:  {prefix: nil, infix: (*Compiler).binary, precedence: PREC_FACTOR},
			token.INT:   {prefix: (*Compiler).number, infix: nil, precedence: PREC_NONE},
			token.FLOAT: {prefix: (*Compiler).number, infix: nil, precedence: PREC_NONE},
		},
	}
	return c
}

// Compiles a stream of `Token`s into `Bytecode`
func (c *Compiler) Compile() (Bytecode, error) {

	c.expression()
	c.emit(OP_END)
	return c.bytecode, nil
}

// Diassembles the compiled bytecode to a human readable format
// and optionally saves it to disk.
// It returns the diassembled bytecode as a string or an error if
// the file could not be created.
func (c *Compiler) DiassembleBytecode(saveToDisk bool) (string, error) {

	var diassembledBytecode string
	var builder strings.Builder
	var instructionLength int
	ip := 0

	// NOTE: Slicing in go includes the first element, but excludes the last one.
	// for example, [0:4] will include index 0 to index 3 of the array.
	for ip <= len(c.bytecode.Instructions) {
		opCode := Opcode(c.bytecode.Instructions[ip])
		switch opCode {
		case OP_ADD, OP_SUBTRACT, OP_DIVIDE, OP_MULTIPLY, OP_NEGATE, OP_END:

			result, err := DiassembleInstruction([]byte{c.bytecode.Instructions[ip]})
			if err != nil {
				panic(err.Error())
			}
			builder.WriteString(result)

		case OP_CONSTANT:
			offset := ip + OP_CONSTANT_TOTAL_BYTES
			instruction := c.bytecode.Instructions[ip:offset]
			index := binary.BigEndian.Uint16(instruction[OPCODE_TOTAL_BYTES:])
			value := c.bytecode.ConstantsPool[index]

			diassembledInstr, err := DiassembleInstruction(instruction)
			if err != nil {
				panic(err.Error())
			}
			result := diassembledInstr + fmt.Sprintf(", value: %d", value)
			builder.WriteString(result)
			builder.WriteString("\n")
			instructionLength = OP_CONSTANT_TOTAL_BYTES

		}

		ip += instructionLength
	}
	diassembledBytecode = builder.String()
	if saveToDisk {
		fDescriptor, err := os.Create("bytecode.txt")
		if err != nil {
			return "", fmt.Errorf("error creating diassembled bytecode file: %s", err.Error())
		}
		fDescriptor.WriteString(diassembledBytecode)
		defer fDescriptor.Close()
	}
	return diassembledBytecode, nil
}

// Retrieves the parsing rule associated with the given token type.
// It returns the parseRule and true if found, otherwise returns an empty parseRule and false.
func (c *Compiler) getParseRule(tokenType token.TokenType) (parseRule, bool) {
	rule, ok := c.parsingRules[tokenType]
	if !ok {
		return parseRule{}, false
	}

	return rule, true
}

// begins parsing an expression from the assignment presedence level
func (c *Compiler) expression() {
	c.parsePresedence(PREC_ASSIGNMENT)
}

// Parses expressions with the provided precedence level.
// It advances the token stream, applies the parse rule, and continues while
// the next token precedence is higher or equal.
func (c *Compiler) parsePresedence(presedence int) {
	c.advance()

	rule, success := c.getParseRule(c.previousTok.TokenType)
	if !success {
		panic("Expected expression")
	}

	rule.prefix(c)

	for presedence <= c.getPresedence(c.currentTok.TokenType) && !c.isFinished() {
		c.advance()
		rule, success := c.getParseRule(c.previousTok.TokenType)
		if !success {
			// Any token sequence without a valid infix or separator rule between them is invalid.
			// for example, two identifiers like x y or two numbers like 5 5 would be considered
			// invalid in the grammar. An infix rule is expected after a valid left-hand expression
			panic("SyntaxError: invalid syntax")
		}
		rule.infix(c)
	}
}

// Parses and emits code for binary operators (+, -, *, /).
// It parses the right-hand operand with higher precedence and
// emits the corresponding bytecode for the operator.
func (c *Compiler) binary() {
	tokenType := c.previousTok.TokenType
	prec := c.getPresedence(tokenType)
	// +1 because each binary operator's right-hand presedence is one
	// level higher than its own
	c.parsePresedence(prec + 1) // compile right hand expression (operand) first
	switch tokenType {
	case token.SUB:
		c.emit(OP_SUBTRACT)
	case token.ADD:
		c.emit(OP_ADD)
	case token.MULT:
		c.emit(OP_MULTIPLY)
	case token.DIV:
		c.emit(OP_DIVIDE)
	}
}

// Parses and emits code for unary operators (!,-).
// It parses the operand and emits the appropriate bytecode for the unary operation.
func (c *Compiler) unary() {
	tokenType := c.previousTok.TokenType
	c.parsePresedence(PREC_UNARY) // // compile right hand expression (oparand) first
	switch tokenType {
	case token.SUB:
		c.emit(OP_NEGATE)
	case token.BANG:
		c.emit(OP_NEGATE)
	default:
		return

	}
}

// parses integer and floating-point literals and emits their
// bytecode representation
func (c *Compiler) number() {
	tokenType := c.previousTok.TokenType
	switch tokenType {
	case token.INT:
		c.handleNumber(c.previousTok)
	case token.FLOAT:
		c.handleNumber(c.previousTok)
	}
}

// Gets the precedence of the given token type.
// If the token type has no defined precedence, PREC_NONE is returned.
func (c *Compiler) getPresedence(tokenType token.TokenType) int {
	prec, ok := precedence[tokenType]
	if !ok {
		return PREC_NONE
	}

	return prec
}

// isFinished returns true if the parser has reached the end of token stream (EOF).
func (c *Compiler) isFinished() bool {
	return c.currentTok.TokenType == token.EOF
}

// advance moves the parser to the next token in the input stream.
// It updates previousTok and currentTok accordingly.
func (c *Compiler) advance() {

	if c.isFinished() {
		return
	}
	c.previousTok = c.tokens[c.readPosition]
	c.readPosition++
	c.currentTok = c.tokens[c.readPosition]
}

// Processes a numeric token into a bytecode instruction.
func (c *Compiler) handleNumber(token token.Token) {
	switch value := token.Literal.(type) {
	case float64:
		c.addConstant(value)
	case int64:
		c.addConstant(value)
	}
}

// Appends a value to the compiler's constant pool and emits an
// `OP_CONSTANT` instruction that references the index of the newly added constant.
// This allows the constant to be used during runtime.
func (c *Compiler) addConstant(value any) {
	c.bytecode.ConstantsPool = append(c.bytecode.ConstantsPool, value)
	index := len(c.bytecode.ConstantsPool) - 1
	c.emit(OP_CONSTANT, index)
}

// Constructs a bytecode instruction from the given opcode and operands,
// then appends the resulting instruction bytes to the compiler's instruction
// stream. This is the low-level mechanism for building the VM instructions.
func (c *Compiler) emit(opcode Opcode, operands ...int) {
	instruction := AssembleInstruction(opcode, operands...)
	c.bytecode.Instructions = append(c.bytecode.Instructions, instruction...)
}
