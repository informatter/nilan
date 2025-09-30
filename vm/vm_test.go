package vm

import (
	"nilan/compiler"
	"testing"
)

func TestExecuteBytecodeVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack []int64
	}{
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(1)},
			},
			expectedStack: []int64{5, 1},
		},
	}

	for _, tt := range tests {

		vm := New()
		vm.Run(tt.bytecode)
		for i := 0; i < len(vm.stack); i++ {
			if vm.stack[i] != tt.expectedStack[i] {
				t.Errorf("vm stack at index: %d - got: %d, want: %d", i, vm.stack[i], tt.expectedStack[i])
			}
		}
	}
}
