package compiler

import (
	"testing"
)

func TestAssembleInstruction(t *testing.T) {

	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		// TODO: add more test cases
		{OP_CONSTANT, []int{65000}, []byte{byte(OP_CONSTANT), 253, 232}},
		{OP_NEGATE, []int{}, []byte{byte(OP_NEGATE)}},
		{OP_SUBTRACT, []int{}, []byte{byte(OP_SUBTRACT)}},
		{OP_ADD, []int{}, []byte{byte(OP_ADD)}},
		{OP_MULTIPLY, []int{}, []byte{byte(OP_MULTIPLY)}},
		{OP_DIVIDE, []int{}, []byte{byte(OP_DIVIDE)}},
	}

	for _, tt := range tests {

		instruction := AssembleInstruction(tt.op, tt.operands...)
		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length - got: %d, want: %d", len(instruction), len(tt.expected))
		}

		for i, byte := range tt.expected {
			instructionByte := instruction[i]
			if instruction[i] != byte {
				t.Errorf("instruction has wrong byte - got: %v, want: %v", instructionByte, byte)
			}
		}

	}

}

func TestDiassembleInstruction(t *testing.T) {
	tests := []struct {
		instruction []byte
		expected    string
	}{
		// TODO: add more test cases
		{[]byte{byte(OP_CONSTANT), 253, 232}, "opcode: OP_CONSTANT, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_SUBTRACT)}, "opcode: OP_SUBTRACT, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_MULTIPLY)}, "opcode: OP_MULTIPLY, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_DIVIDE)}, "opcode: OP_DIVIDE, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_ADD)}, "opcode: OP_ADD, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_NEGATE)}, "opcode: OP_NEGATE, operand: None, operand widths: 0 bytes"},
	}

	for _, tt := range tests {
		err := DiassembleInstruction(tt.instruction)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
