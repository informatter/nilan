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
	PREC_NONE = iota // LOWEST PRESEDENCE
	PREC_ASSIGNMENT
	PREC_TERM   // +,-
	PREC_FACTOR // /,*
	PREC_UNARY  // !, -, // HIGHEST PRESEDENCE
)

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
	currentTok   token.Token
	nextTok      token.Token
	parsingRules map[token.TokenType]parseRule
}

// Creates a `Compiler` instance and returns
// a pointer to it.
func New(tokens []token.Token) *Compiler {
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

	err := c.expression()
	if err != nil {
		return c.bytecode, err
	}
	c.emit(OP_END)
	return c.bytecode, nil
}

// Writes the compiled bytecode to a file with a `.nic` extension.
// The bytecode  is encoded as hexadecimal so it can be viewed in a
// text editor
func (c *Compiler) DumpBytecode(filePath string) error {

	if filePath == "" {
		filePath = "bytecode.nic"
	} else {
		filePath = filePath + ".nic"
	}
	fDescriptor, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating nilan bytecode file: %s", err.Error())
	}

	encoded := fmt.Sprintf("%x", c.bytecode.Instructions)
	fDescriptor.Write([]byte(encoded))
	defer fDescriptor.Close()
	return nil
}

// Diassembles the compiled bytecode to a human readable format
// and optionally saves it to disk.
// It returns the diassembled bytecode as a string or an error if
// the file could not be created.
func (c *Compiler) DiassembleBytecode(saveToDisk bool, filePath string) (string, error) {

	var diassembledBytecode string
	var builder strings.Builder
	var instructionLength int
	totalInstructions := len(c.bytecode.Instructions) - 1
	ip := 0

	// NOTE: Slicing in go includes the first element, but excludes the last one.
	// for example, [0:4] will include index 0 to index 3 of the array.
	for ip <= totalInstructions {
		opCode := Opcode(c.bytecode.Instructions[ip])
		switch opCode {
		case OP_ADD, OP_SUBTRACT, OP_DIVIDE, OP_MULTIPLY, OP_NEGATE, OP_END:

			result, err := DiassembleInstruction([]byte{c.bytecode.Instructions[ip]})
			if err != nil {
				panic(err.Error())
			}
			builder.WriteString(result)
			if opCode == OP_END {
				break
			}
			builder.WriteString("\n")
			instructionLength = OPCODE_TOTAL_BYTES

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
		if filePath == "" {
			filePath = "bytecode.dnic"
		} else {
			filePath = filePath + ".dnic"
		}
		fDescriptor, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("error creating diassembled bytecode file: %s", err.Error())
		}
		fDescriptor.WriteString(diassembledBytecode)
		defer fDescriptor.Close()
	}
	return diassembledBytecode, nil
}

// Retrieves the parsing rule associated with the given token type.
// It returns a valid `parseRule“, or an invalid `parseRule` if a `parseRule`
// was not found for the `TokenType`.
func (c *Compiler) getParseRule(tokenType token.TokenType) parseRule {
	rule, ok := c.parsingRules[tokenType]
	if !ok {
		return parseRule{prefix: nil, infix: nil}
	}

	return rule
}

// begins parsing an expression from the assignment presedence level
// A `SyntaxError` is returned if an error occurs.
func (c *Compiler) expression() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				err = SyntaxError{
					Message: v,
				}
			}
		}
	}()
	c.parsePresedence(PREC_ASSIGNMENT)
	return nil
}

// Parses expressions with the provided precedence level.
// It advances the token stream, applies the parse rule, and continues while
// the next token precedence is higher or equal.
func (c *Compiler) parsePresedence(presedence int) {
	c.advance()

	rule := c.getParseRule(c.currentTok.TokenType)
	if rule.prefix == nil {
		panic("Expected expression")
	}

	rule.prefix(c)

	for c.getParseRule(c.nextTok.TokenType).precedence >= presedence && !c.isFinished() {
		c.advance()
		rule := c.getParseRule(c.currentTok.TokenType)
		if rule.infix == nil {
			// Any token sequence without a valid infix or separator rule between them is invalid.
			// for example, two identifiers like x y or two numbers like 5 5 would be considered
			// invalid in the grammar. An infix rule is expected after a valid left-hand expression
			panic("Invalid syntax")
		}
		rule.infix(c)
	}
}

// Parses and emits code for binary operators (+, -, *, /).
// It parses the right-hand operand with higher precedence and
// emits the corresponding bytecode for the operator.
func (c *Compiler) binary() {
	operator := c.currentTok
	rule := c.getParseRule(operator.TokenType)
	// +1 because each binary operator's right-hand presedence is one
	// level higher than its own
	c.parsePresedence(rule.precedence + 1) // compile right hand expression (operand) first
	switch operator.TokenType {
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
	tokenType := c.currentTok.TokenType
	c.parsePresedence(PREC_UNARY) // compile right hand expression (oparand) first
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
	tokenType := c.currentTok.TokenType
	switch tokenType {
	case token.INT:
		c.handleNumber(c.currentTok)
	case token.FLOAT:
		c.handleNumber(c.currentTok)
	}
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
	c.currentTok = c.tokens[c.readPosition]
	c.readPosition++
	c.nextTok = c.tokens[c.readPosition]
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
