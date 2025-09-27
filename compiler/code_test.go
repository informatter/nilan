package compiler

import (
	"testing"
)

func TestMakeInstruction(t *testing.T) {

	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		//
		{OP_CONSTANT, []int{65000}, []byte{byte(OP_CONSTANT), 253, 232}},
	}

	for _, tt := range tests {

		instruction := MakeInstruction(tt.op, tt.operands...)
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
