package compiler

import (
	"encoding/binary"
	"fmt"
)

// Represents the definition of the `Bytecode`
// which will be created by the compiler and passed to
// the Virtual Machine (VM) to execute
//
// Fields:
//   - Instructions: An array of instructions defined by opcodes and
// 	   					their operands
//   - ConstantsPool: An array containing all the constant values from the source code.
type Bytecode struct {
	Instructions  Instructions
	ConstantsPool []any
}

type Opcode byte

type Instructions []byte

// opcodes
//iota generates a distinct byte for each bytecode
const (
	// represents a opcode constant with a single operand with a size of
	// 2 bytes, which represents a `uint16`.
	// `uint16` -> set of all unsigned 16-bit integers (0 to 65535)
	// this will restrict a nilan program to have a total of 65535 constants.
	// NOTE: This is not a hard constraint, could be changed to uint32 if needed
	OP_CONSTANT Opcode = iota
)

// Represents a definition of an opcode.
// Fields:
//   - Name: The human-readable name for the opcode e.g "OP_CONSTANT"
//   - OperandBytes: The number of bytes each operand takes up.
type OpCodeDefinition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*OpCodeDefinition{
	// has a single operand which takes two bytes of memory.
	OP_CONSTANT: {Name: "OP_CONSTANT", OperandWidths: []int{2}},
}

func Get(op Opcode) (*OpCodeDefinition, error) {
	def, ok := definitions[op]
	if !ok {
		return nil, fmt.Errorf("opcode: '%c' undefined", op)
	}
	return def, nil
}

// Constructs a bytecode instruction from an opcode and its operands.
// The bytecode operands are encoded in BigEndian order
//
// The resulting byte slice always begins with the opcode, followed by each
// operand encoded according to its defined width in Big-Endian order. This
// means that each `uint16` operand will be encoded with the two bytes stored with the most significant
// byte first (the largest byte), followed by the least significant byte (the smallest byte).
// For example, the instruction for OP_CONSTANT could be defined as:
// [0,253,232] , if its operand is 65000. 65000 in Big Endian format is defined as
// 255 and 232.
//
// Parameters:
//   - op: The opcode representing the instruction to encode.
//   - operands: A variadic list of integers providing the operand values
//     corresponding to the opcode's expected operand widths.
//
// Returns:
//   - A byte slice containing the encoded instruction. If the opcode is not
//     recognized, an empty slice is returned.
//
// Example:
//   // Suppose OP_CONSTANT expects a 2-byte operand (index into constants table).
//   instr := MakeBytecode(OP_CONSTANT, 42)
//   // instr now contains: [<opcode for OP_CONSTANT>, 0x00, 0x2A]
//
func MakeInstruction(op Opcode, operands ...int) []byte {
	def, err := Get(op)
	if err != nil {
		return []byte{}
	}

	byteOffset := 1
	instructionLength := byteOffset // starts at one for the opcode
	for _, i := range def.OperandWidths {
		instructionLength += i
	}

	instruction := make([]byte, instructionLength)

	// The firt byte of the instruction will be the opcode
	instruction[0] = byte(op)

	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[byteOffset:], uint16(o))
		}
		byteOffset += width
	}
	return instruction
}
