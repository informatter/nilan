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
	}

	for _, tt := range tests {

		instruction := AssmebleInstruction(tt.op, tt.operands...)
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

func TestDiassembleInstruction(t *testing.T   ) {
	tests := []struct {
		instruction []byte
		expected    string
	}{
		// TODO: add more test cases
		{[]byte{byte(OP_CONSTANT), 253, 232}, "opcode: OP_CONSTANT, operand: 65000, operand widths: 2 bytes"},
	}

	for _, tt := range tests {
		err := DiassembleInstruction(tt.instruction)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
