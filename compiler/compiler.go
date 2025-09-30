package compiler

import (
	"nilan/token"
)

// Represents the compiler which will compile
// a stream of `Token`'s to `Bytecode` to be executed
// by the VM
type Compiler struct {
	bytecode Bytecode
}

// Creates `Compiler` instance and returns
// a pointer to it.
func NewCompiler() *Compiler {
	return &Compiler{
		bytecode: Bytecode{
			Instructions:  Instructions{},
			ConstantsPool: []any{},
		},
	}
}

// Compiles an array of `Token`'s into `Bytecode`
func (c *Compiler) Compile(tokens []token.Token) (Bytecode, error) {

	for _, tok := range tokens {

		switch tok.TokenType {
		// TODO: Handle other tokens.
		case token.EOF:
			c.emit(OP_END)
		default:
			c.handleNumber(tok)
		}
	}
	return c.bytecode, nil
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
