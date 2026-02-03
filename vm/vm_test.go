package vm

import (
	"nilan/compiler"
	"testing"
)

func assertResults(tests []struct {
	bytecode      compiler.Bytecode
	expectedStack any
}, t *testing.T) {
	t.Helper()
	for _, tt := range tests {
		vm := New()
		err := vm.Run(tt.bytecode)
		if err != nil {
			t.Error(err.Error())
		}

		expStack, ok := tt.expectedStack.([]any)
		if !ok {
			t.Errorf("expectedStack must be []any containing float64 or int64")
			continue
		}
		if len(vm.stack) != len(expStack) {
			t.Errorf("stack length mismatch: got %d, want %d", len(vm.stack), len(expStack))
			continue
		}

		for i := 0; i < len(vm.stack); i++ {
			expected := expStack[i]
			actual := vm.stack[i]
			if actual != expected {
				t.Errorf("vm stack at index: %d - got: %d, want: %d", i, actual, expected)
			}
			switch exp := expected.(type) {
			case float64:
				if actual != exp {
					t.Errorf("stack[%d]: got %f, want %f", i, actual, exp)
				}
			case int64:
				if actual != exp {
					t.Errorf("stack[%d]: got %d, want %d", i, actual, exp)
				}
			case bool:
				if actual != exp {
					t.Errorf("stack[%d]: got %t, want %t", i, actual, exp)
				}
			default:
				t.Errorf("stack[%d]: unsupported expected type %T", i, expected)
			}
		}
	}
}

// Tests logical expressions in the bytecode without using variables
func TestExecuteBytecodeLogicalOpVMStack(t *testing.T) {
	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
	}{

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_AND),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{true, true},
			},
			expectedStack: []any{true},
		},

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_AND),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{true, false},
			},
			expectedStack: []any{false},
		},

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_OR),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{true, false},
			},
			expectedStack: []any{true},
		},

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_OR),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{false, false},
			},
			expectedStack: []any{false},
		},

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_AND),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{false, false},
			},
			expectedStack: []any{false},
		},

		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_OR),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{false, true},
			},
			expectedStack: []any{true},
		},
	}
	assertResults(tests, t)
}

// Tests logical expressions with variables in the bytecode
func TestExecuteBytecodeLogicalOpWithVariables(t *testing.T) {
	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
	}{
		// Define global var a = true, b = false, then a && b, a || b
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					// Define a = true
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_DEFINE_GLOBAL), 0, 0,
					// Define b = false
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_DEFINE_GLOBAL), 0, 1,
					// Get a
					byte(compiler.OP_GET_GLOBAL), 0, 0,
					// Get b
					byte(compiler.OP_GET_GLOBAL), 0, 1,
					// AND
					byte(compiler.OP_AND),
					// Get a
					byte(compiler.OP_GET_GLOBAL), 0, 0,
					// Get b
					byte(compiler.OP_GET_GLOBAL), 0, 1,
					// OR
					byte(compiler.OP_OR),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{true, false},
				NameConstants: []string{"a", "b"},
			},
			expectedStack: []any{false, true}, // a && b, a || b
		},
	}
	assertResults(tests, t)
}

func TestExecuteBytecodeVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
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
			expectedStack: []any{int64(5), int64(1)},
		},
	}

	assertResults(tests, t)
}

func TestExecuteBytecodeBinaryOpVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
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
			expectedStack: []any{int64(6)},
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
			expectedStack: []any{int64(19)},
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
			expectedStack: []any{int64(30)},
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
			expectedStack: []any{int64(0)},
		},
	}

	assertResults(tests, t)
}

func TestExecuteBytecodeBinaryOpFloatVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
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
			expectedStack: []any{float64(8.3)},
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
			expectedStack: []any{float64(8.95)},
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
			expectedStack: []any{float64(4.5)},
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
			expectedStack: []any{float64(2.0)},
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
			expectedStack: []any{float64(3.4703947368421053)},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_DIVIDE),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{5.544, 21.943},
			},
			expectedStack: []any{0.25265460511324794},
		},
	}

	assertResults(tests, t)
}

func TestExecuteBytecodeNegateOpVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
	}{
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_NEGATE),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(10)},
			},
			expectedStack: []any{int64(-10)},
		},
	}

	assertResults(tests, t)
}

func TestExecuteBytecodeNotOpVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
	}{
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_NOT),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(true)},
			},
			expectedStack: []any{bool(false)},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_NOT),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(false)},
			},
			expectedStack: []any{bool(true)},
		},
	}

	assertResults(tests, t)
}

func TestExecuteBytecodePrintStatement(t *testing.T) {
	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
	}{
		{
			bytecode: compiler.Bytecode{
				Instructions:  []byte{byte(compiler.OP_CONSTANT), 0, 0, byte(compiler.OP_CONSTANT), 0, 1, byte(compiler.OP_ADD), byte(compiler.OP_PRINT), byte(compiler.OP_END)},
				ConstantsPool: []any{int64(10), int64(3)},
			},
			expectedStack: []any{},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_NOT),
					byte(compiler.OP_PRINT),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(false)},
			},
			expectedStack: []any{},
		},
	}

	assertResults(tests, t)
}

func TestExecuteBytecodeComparisonOpVMStack(t *testing.T) {

	tests := []struct {
		bytecode      compiler.Bytecode
		expectedStack any
	}{
		// Equality tests
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_EQUALITY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(5)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_EQUALITY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_EQUALITY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), float64(5.0)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_EQUALITY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), float64(3.0)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_EQUALITY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(true), bool(true)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_EQUALITY),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(true), bool(false)},
			},
			expectedStack: []any{false},
		},
		// Not Equal tests
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_NOT_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(5)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_NOT_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_NOT_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), float64(5.0)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_NOT_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), float64(3.0)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_NOT_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(true), bool(true)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_NOT_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{bool(true), bool(false)},
			},
			expectedStack: []any{true},
		},
		// Larger tests
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(3), int64(5)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), float64(3.0)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), float64(3.0)},
			},
			expectedStack: []any{true},
		},

		// Less tests
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(3), int64(5)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(3.0), float64(5.0)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), int64(3)},
			},
			expectedStack: []any{false},
		},
		// Larger Equal tests
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(5)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(3), int64(5)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LARGER_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(5.0), float64(5.0)},
			},
			expectedStack: []any{true},
		},
		// Less Equal tests
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(5)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(3), int64(5)},
			},
			expectedStack: []any{true},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{int64(5), int64(3)},
			},
			expectedStack: []any{false},
		},
		{
			bytecode: compiler.Bytecode{
				Instructions: []byte{
					byte(compiler.OP_CONSTANT), 0, 0,
					byte(compiler.OP_CONSTANT), 0, 1,
					byte(compiler.OP_LESS_EQUAL),
					byte(compiler.OP_END),
				},
				ConstantsPool: []any{float64(3.0), float64(5.0)},
			},
			expectedStack: []any{true},
		},
	}

	assertResults(tests, t)
}
