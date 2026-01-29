package vm

import (
	"nilan/compiler"
	"testing"
)

func assertResults(tests []struct {
	bytecode      compiler.Bytecode
	expectedStack []float64
}, t *testing.T) {

	t.Helper()
	for _, tt := range tests {

		vm := New()
		err := vm.Run(tt.bytecode)
		if err != nil {
			t.Error(err.Error())
		}
		if len(vm.stack) == 0 {
			t.Errorf("vm stack should not be empty")
		}
		for i := 0; i < len(vm.stack); i++ {
			if vm.stack[i] != tt.expectedStack[i] {
				t.Errorf("vm stack at index: %d - got: %f, want: %f", i, vm.stack[i], tt.expectedStack[i])
			}
		}
	}
}

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

func TestExecuteBytecodeBinaryOpVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack []int64
	}{
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_ADD),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(1)},
			},
			expectedStack: []int64{6},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_CONSTANT), 0, 2,
					byte(compiler.OP_CONSTANT), 0, 3,
					byte(compiler.OP_ADD),
					byte(compiler.OP_ADD),
					byte(compiler.OP_ADD),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(1), int64(3), int64(10)},
			},
			expectedStack: []int64{19},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_CONSTANT), 0, 2,
					byte(compiler.OP_MULTIPLY),
					byte(compiler.OP_MULTIPLY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3), int64(2)},
			},
			expectedStack: []int64{30},
		},

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_SUBTRACT),
					byte(compiler.OP_CONSTANT), 0, 2,
					byte(compiler.OP_SUBTRACT),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3), int64(2)},
			},
			expectedStack: []int64{0},
		},
	}

	for _, tt := range tests {

		vm := New()
		err := vm.Run(tt.bytecode)
		if err != nil {
			t.Error(err.Error())
		}
		if len(vm.stack) == 0 {
			t.Errorf("vm stack should not be empty")
		}
		for i := 0; i < len(vm.stack); i++ {
			if vm.stack[i] != tt.expectedStack[i] {
				t.Errorf("vm stack at index: %d - got: %d, want: %d", i, vm.stack[i], tt.expectedStack[i])
			}
		}
	}
}

func TestExecuteBytecodeBinaryOpFloatVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack []float64
	}{
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_ADD),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.3), int64(3)},
			},
			expectedStack: []float64{8.3},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_ADD),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.3), float64(3.65)},
			},
			expectedStack: []float64{8.95},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_DIVIDE),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(9), int64(2)},
			},
			expectedStack: []float64{4.5},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_DIVIDE),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(4), int64(2)},
			},
			expectedStack: []float64{2.0},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_DIVIDE),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(10.55), float64(3.04)},
			},
			expectedStack: []float64{3.4703947368421053},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_DIVIDE),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.544), float64(21.943)},
			},
			expectedStack: []float64{0.25265460511324794},
		},
	}

	assertResults(tests, t)
}
