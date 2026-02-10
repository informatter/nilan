package compiler

import (
	"testing"
)

func TestAssembleInstruction(t *testing.T) {
	operand := 65000
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OP_CONSTANT, []int{operand}, []byte{byte(OP_CONSTANT), 253, 232}},
		{OP_END, []int{}, []byte{byte(OP_END)}},
		{OP_ADD, []int{}, []byte{byte(OP_ADD)}},
		{OP_MULTIPLY, []int{}, []byte{byte(OP_MULTIPLY)}},
		{OP_DIVIDE, []int{}, []byte{byte(OP_DIVIDE)}},
		{OP_SUBTRACT, []int{}, []byte{byte(OP_SUBTRACT)}},
		{OP_NEGATE, []int{}, []byte{byte(OP_NEGATE)}},
		{OP_NOT, []int{}, []byte{byte(OP_NOT)}},
		{OP_PRINT, []int{}, []byte{byte(OP_PRINT)}},
		{OP_AND, []int{}, []byte{byte(OP_AND)}},
		{OP_OR, []int{}, []byte{byte(OP_OR)}},
		{OP_EQUALITY, []int{}, []byte{byte(OP_EQUALITY)}},
		{OP_NOT_EQUAL, []int{}, []byte{byte(OP_NOT_EQUAL)}},
		{OP_LARGER, []int{}, []byte{byte(OP_LARGER)}},
		{OP_LESS, []int{}, []byte{byte(OP_LESS)}},
		{OP_LARGER_EQUAL, []int{}, []byte{byte(OP_LARGER_EQUAL)}},
		{OP_LESS_EQUAL, []int{}, []byte{byte(OP_LESS_EQUAL)}},
		{OP_DEFINE_GLOBAL, []int{operand}, []byte{byte(OP_DEFINE_GLOBAL), 253, 232}},
		{OP_SET_GLOBAL, []int{operand}, []byte{byte(OP_SET_GLOBAL), 253, 232}},
		{OP_GET_GLOBAL, []int{operand}, []byte{byte(OP_GET_GLOBAL), 253, 232}},
		{OP_DEFINE_LOCAL, []int{operand}, []byte{byte(OP_DEFINE_LOCAL), 253, 232}},
		{OP_SET_LOCAL, []int{operand}, []byte{byte(OP_SET_LOCAL), 253, 232}},
		{OP_GET_LOCAL, []int{operand}, []byte{byte(OP_GET_LOCAL), 253, 232}},
		{OP_JUMP, []int{operand}, []byte{byte(OP_JUMP), 253, 232}},
		{OP_JUMP_IF_FALSE, []int{operand}, []byte{byte(OP_JUMP_IF_FALSE), 253, 232}},
		{OP_POP, []int{}, []byte{byte(OP_POP)}},
	}

	for _, tt := range tests {

		instruction, err := AssembleInstruction(tt.op, tt.operands...)

		if err != nil {
			t.Error("error assembling instruction")
		}
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
		{[]byte{byte(OP_CONSTANT), 253, 232}, "opcode: OP_CONSTANT, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_END)}, "opcode: OP_END, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_ADD)}, "opcode: OP_ADD, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_MULTIPLY)}, "opcode: OP_MULTIPLY, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_DIVIDE)}, "opcode: OP_DIVIDE, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_SUBTRACT)}, "opcode: OP_SUBTRACT, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_NEGATE)}, "opcode: OP_NEGATE, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_NOT)}, "opcode: OP_NOT, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_PRINT)}, "opcode: OP_PRINT, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_AND)}, "opcode: OP_AND, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_OR)}, "opcode: OP_OR, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_EQUALITY)}, "opcode: OP_EQUALITY, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_NOT_EQUAL)}, "opcode: OP_NOT_EQUAL, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_LARGER)}, "opcode: OP_LARGER, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_LESS)}, "opcode: OP_LESS, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_LARGER_EQUAL)}, "opcode: OP_LARGER_EQUAL, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_LESS_EQUAL)}, "opcode: OP_LESS_EQUAL, operand: None, operand widths: 0 bytes"},
		{[]byte{byte(OP_DEFINE_GLOBAL), 253, 232}, "opcode: OP_DEFINE_GLOBAL, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_SET_GLOBAL), 253, 232}, "opcode: OP_SET_GLOBAL, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_GET_GLOBAL), 253, 232}, "opcode: OP_GET_GLOBAL, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_DEFINE_LOCAL), 253, 232}, "opcode: OP_DEFINE_LOCAL, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_SET_LOCAL), 253, 232}, "opcode: OP_SET_LOCAL, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_GET_LOCAL), 253, 232}, "opcode: OP_GET_LOCAL, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_JUMP), 253, 232}, "opcode: OP_JUMP, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_JUMP_IF_FALSE), 253, 232}, "opcode: OP_JUMP_IF_FALSE, operand: 65000, operand widths: 2 bytes"},
		{[]byte{byte(OP_POP)}, "opcode: OP_POP, operand: None, operand widths: 0 bytes"},
	}

	for _, tt := range tests {
		result, err := DiassembleInstruction(tt.instruction)
		if err != nil {
			t.Errorf(err.Error())
		}

		if tt.expected != result {
			t.Errorf("wrong diassembled instruction - got: %s, want: %s", result, tt.expected)
		}
	}
}
