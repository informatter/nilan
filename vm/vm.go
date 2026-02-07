package vm

import (
	"encoding/binary"
	"fmt"
	"nilan/compiler"
)

type arithmeticFuncFloat func(a float64, b float64) float64
type arithmeticFuncInt func(a int64, b int64) int64

func addFloat(a float64, b float64) float64 {
	return a + b
}
func addInt(a int64, b int64) int64 {
	return a + b
}
func subFloat(a float64, b float64) float64 {
	return a - b
}
func subInt(a int64, b int64) int64 {
	return a - b
}
func multFloat(a float64, b float64) float64 {
	return a * b
}
func multInt(a int64, b int64) int64 {
	return a * b
}
func divFloat(a float64, b float64) float64 {
	// TODO: Add runtime error division by zero
	return a / b
}
func divInt(a int64, b int64) int64 {
	// TODO: Add runtime error division by zero

	return a / b
}

type equalityFuncFloat func(a float64, b float64) bool
type equalityFuncInt func(a int64, b int64) bool

func largerThanInt(a int64, b int64) bool {
	return a > b
}

func largerThanFloat(a float64, b float64) bool {
	return a > b
}

func smallerThanInt(a int64, b int64) bool {
	return a < b
}

func smallerThanFloat(a float64, b float64) bool {
	return a < b
}

func largerEqualInt(a int64, b int64) bool {
	return a >= b
}

func largerEqualFloat(a float64, b float64) bool {
	return a >= b
}

func smallerEqualInt(a int64, b int64) bool {
	return a <= b
}

func smallerEqualFloat(a float64, b float64) bool {
	return a <= b
}

// Determines if a value is a float.
//
// Parameters:
//   - value: the literal value (various possible types).
//
// Returns:
//   - bool: true if the value is an float, false otherwise.
func isFloat(value any) bool {
	switch value.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// Determines if a value is an integer.
//
// Parameters:
//   - value: the literal value (various possible types).
//
// Returns:
//   - bool: true if the value is an integer, false otherwise.
func isInt(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

// isNumeric determines if a value is either an integer or a float.
func isNumeric(value any) bool {
	return isFloat(value) || isInt(value)
}

func isBool(val any) bool {
	_, ok := val.(bool)
	return ok
}

// Attempts to convert a literal value into a int64.
//
// Parameters:
//   - value: the literal value (various possible types).
//
// Returns:
//   - int64: the converted numeric value.
//   - error: on failure to convert value to float64.
func literalToInt64(value any) (int64, error) {

	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// Attempts to convert a literal value into a float64.
//
// Parameters:
//   - value: the literal value (various possible types).
//
// Returns:
//   - float64: the converted numeric value.
//   - error: on failure to convert value to float64.
func literalToFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// comparisonOpHandler defines a function type for handling comparison operations.
type comparisonOpHandler func(*VirtualMachine, equalityFuncFloat, equalityFuncInt) error

// makeComparisonHandler creates a comparison operation handler
// using the provided float and int equality functions.
func makeComparisonHandler(f equalityFuncFloat, i equalityFuncInt) comparisonOpHandler {
	return func(vm *VirtualMachine, _ equalityFuncFloat, _ equalityFuncInt) error {
		return vm.handleNumericEqualityOps(f, i)
	}

}

func isFalsey(value any) bool {
	if value == nil {
		return true
	}
	if isBool(value) && value.(bool) == false {
		return true
	}

	// Everything else is considered to be true.
	return false
}

// Represents a stack based virtual-machine (VirtualMachine).
// It is the runtime environment where Nilan bytecode
// gets executed.
type VirtualMachine struct {

	// stack stores intermediate values during bytecode execution,
	// which are poped and pushed as instructions are executed.
	stack Stack
	// instruction pointer stores the address of the current bytecode instruction.
	// It determines where the VM is in the program.
	ip    int
	debug bool
	// globalVars stores the mapping of global variable names to their corresponding values.
	globalVars map[string]any
	// comparisonOpHandlers maps comparison opcodes to their corresponding handler functions.
	comparisonOpHandlers map[compiler.Opcode]comparisonOpHandler
}

// Creates a new VM instance
func New() *VirtualMachine {
	return &VirtualMachine{
		debug:      true,
		globalVars: make(map[string]any),
		comparisonOpHandlers: map[compiler.Opcode]comparisonOpHandler{
			compiler.OP_LARGER:       makeComparisonHandler(largerThanFloat, largerThanInt),
			compiler.OP_LESS:         makeComparisonHandler(smallerThanFloat, smallerThanInt),
			compiler.OP_LARGER_EQUAL: makeComparisonHandler(largerEqualFloat, largerEqualInt),
			compiler.OP_LESS_EQUAL:   makeComparisonHandler(smallerEqualFloat, smallerEqualInt),
		},
	}
}

// handleNumericEqualityOps applies numeric comparison functions to the two topmost
// values on the VM stack.
func (vm *VirtualMachine) handleNumericEqualityOps(floatFunc equalityFuncFloat, intFunc equalityFuncInt) error {
	b := vm.stack.Pop()
	a := vm.stack.Pop()
	if !isNumeric(a) || !isNumeric(b) {
		return RuntimeError{Message: fmt.Sprintf("operands must be numeric values: %v,%v", a, b)}
	}
	isAFloat := isFloat(a)
	isBFloat := isFloat(b)
	if isAFloat && isBFloat {
		a, error := literalToFloat64(a)
		b, err := literalToFloat64(b)
		if error != nil || err != nil {
			// This currently handles a very rage edge condition, where either number is not
			// a float32 or float64, despite isFloat returning true.
			message := fmt.Sprintf("operands must be valid floating point values: %v,%v", a, b)
			return RuntimeError{Message: message}
		}
		vm.stack.Push(floatFunc(a, b))
	} else if isAFloat && !isBFloat {
		a, error := literalToFloat64(a)

		if error != nil {
			// This currently handles a very rage edge condition, where either number is not
			// a float32/float64 or int types, despite isFloat/isInt returning true.
			message := fmt.Sprintf("operands must be valid integer and floating point values: %v,%v", a, b)
			return RuntimeError{Message: message}
		}
		if bv, ok := b.(int64); ok {
			vm.stack.Push(floatFunc(a, float64(bv)))
		}
	} else if !isAFloat && isBFloat {

		b, error := literalToFloat64(b)
		if error != nil {
			// This currently handles a very rage edge condition, where either number is not
			// a float32/float64 or int types, despite isFloat/isInt returning true.
			message := fmt.Sprintf("operands must be valid integer and floating point values: %v,%v", a, b)
			return RuntimeError{Message: message}
		}
		if av, ok := a.(int64); ok {
			vm.stack.Push(floatFunc(float64(av), b))
		}
	} else {
		a, error := literalToInt64(a)
		b, err := literalToInt64(b)
		if error != nil || err != nil {
			// This currently handles a very rage edge condition, where either number is not
			// an int type, despite isInt returning true.
			message := fmt.Sprintf("operands must be valid integer values: %v,%v", a, b)
			return RuntimeError{Message: message}
		}
		vm.stack.Push(intFunc(a, b))
	}
	return nil
}

// Executes the provided bytecode on the virtual machine (VM).
//
// It fetches and decodes each instruction starting at the VM's current
// instruction pointer (ip), processes the instruction based on its opcode,
// and modifies the VM's state accordingly (e.g. pushing constants onto the stack).
//
// The instruction pointer (ip) is incremented by the size of the current
// instruction after its execution.
//
// Execution terminates normally when an OP_END opcode is encountered,
// or returns an error if an unknown opcode is found.
//
// Parameters:
//   - bytecode: The compiled instructions to execute.
//
// Returns:
//   - error: Any error encountered during execution, including unknown opcodes.
func (vm *VirtualMachine) Run(bytecode compiler.Bytecode) error {

	var instructionLength int
	for {
		opCode := compiler.Opcode(bytecode.Instructions[vm.ip])
		intOpCode := int(opCode)

		switch opCode {
		case compiler.OP_END:
			if vm.stack.Peek() != nil {
				// NOTE: temp code to handle operations such as 2+2 to be printed in the REPL
				// Can there be a more suitable place to handle this other than in the VM?
				// for now it does not hurt to leave it here...
				fmt.Println(vm.stack.Peek())
			}
			return nil

		case compiler.OP_POP:
			vm.stack.Pop()
			instructionLength = compiler.OPCODE_TOTAL_BYTES

		case compiler.OP_PRINT:
			l := vm.execPrintInstruction()
			instructionLength = l
		case compiler.OP_CONSTANT:
			instructionLength = vm.execConstantInstruction(bytecode)

		case compiler.OP_ADD:
			l, err := vm.execArithmeticInstruction(addFloat, addInt, intOpCode)
			if err != nil {
				return err
			}
			instructionLength = l
		case compiler.OP_SUBTRACT:
			l, err := vm.execArithmeticInstruction(subFloat, subInt, intOpCode)
			if err != nil {
				return err
			}
			instructionLength = l
		case compiler.OP_MULTIPLY:
			l, err := vm.execArithmeticInstruction(multFloat, multInt, intOpCode)
			if err != nil {
				return err
			}
			instructionLength = l
		case compiler.OP_DIVIDE:
			l, err := vm.execArithmeticInstruction(divFloat, divInt, intOpCode)
			if err != nil {
				return err
			}
			instructionLength = l

		case compiler.OP_NEGATE, compiler.OP_NOT:
			l, err := vm.execUnaryInstruction(opCode)
			if err != nil {
				return err
			}
			instructionLength = l

		case compiler.OP_LARGER, compiler.OP_LESS, compiler.OP_LARGER_EQUAL, compiler.OP_LESS_EQUAL:
			l, err := vm.execComparisonInstruction(opCode)
			if err != nil {
				return err
			}
			instructionLength = l
		case compiler.OP_EQUALITY, compiler.OP_NOT_EQUAL:
			l, err := vm.execEqualityInstruction(opCode)
			if err != nil {
				return err
			}
			instructionLength = l

		case compiler.OP_AND, compiler.OP_OR:
			l, err := vm.execLogicalInstruction(opCode)
			if err != nil {
				return err
			}
			instructionLength = l

		// NOTE: `continue` is needed as the VM needs to jump to a target instruction
		// instead of incrementing the instruction pointer to the next instruction in sequence.
		case compiler.OP_JUMP:
			vm.ip = vm.execJumpInstruction(bytecode)
			continue
		case compiler.OP_JUMP_IF_FALSE:
			vm.ip = vm.execJumpIfFalseInstruction(bytecode)
			continue

		case compiler.OP_DEFINE_GLOBAL, compiler.OP_SET_GLOBAL:
			instructionLength = vm.execDefineGlobalInstruction(bytecode)
		case compiler.OP_GET_GLOBAL:
			instructionLength = vm.execGetGlobalInstruction(bytecode)
		default:
			// NOTE: This should only happen in development mode.
			return fmt.Errorf("unknown opcode %v at ip %d", opCode, vm.ip)
		}

		vm.ip += instructionLength
	}
}

func (vm *VirtualMachine) execPrintInstruction() int {
	value := vm.stack.Pop()
	if value == nil {
		fmt.Println("null")
		return compiler.OPCODE_TOTAL_BYTES
	}

	fmt.Println(value)
	return compiler.OPCODE_TOTAL_BYTES
}

// execJumpInstruction executes a `OP_JUMP` instruction by reading the target byte
// offset from the instruction's operand and returning it.
func (vm *VirtualMachine) execJumpInstruction(bytecode compiler.Bytecode) int {

	operandIndex := vm.ip + compiler.OPCODE_TOTAL_BYTES
	// skips the opcode byte and reads the next 2 bytes, to retrieve the
	// operand which represents the target byte offset to jump to
	instruction := bytecode.Instructions[operandIndex : operandIndex+2]
	targetByteOffset := binary.BigEndian.Uint16(instruction)

	return int(targetByteOffset)
}

// execJumpIfFalseInstruction executes a `OP_JUMP_IF_FALSE` instruction by evaluating the condition
// on top of the stack and determining whether to jump to the target byte offset
// or continue to the next instruction.
func (vm *VirtualMachine) execJumpIfFalseInstruction(bytecode compiler.Bytecode) int {

	condition := vm.stack.Peek()
	if isFalsey(condition) {
		operandIndex := vm.ip + compiler.OPCODE_TOTAL_BYTES
		instruction := bytecode.Instructions[operandIndex : operandIndex+2]
		targetByteOffset := binary.BigEndian.Uint16(instruction)
		// If the condition is falsey, the VM should jump to the beginning of the
		// else block (or the end of the if statement if there is no else block),
		// which is located at the byte offset specified by the instruction's operand.
		return int(targetByteOffset)
	}

	// if the condition is truthy, the VM should continue executing the next
	// instruction, which is located immediatelty after the currrent instructions
	// opcode and operand bytes.

	return vm.ip + compiler.OP_JUMP_TOTAL_BYTES
}

// execEqualityInstruction executes equality or inequality operations based on the provided opcode.
func (vm *VirtualMachine) execEqualityInstruction(opCode compiler.Opcode) (int, error) {

	if opCode == compiler.OP_EQUALITY {
		b := vm.stack.Pop()
		a := vm.stack.Pop()
		vm.stack.Push(a == b)
		return compiler.OPCODE_TOTAL_BYTES, nil
	}
	if opCode == compiler.OP_NOT_EQUAL {
		b := vm.stack.Pop()
		a := vm.stack.Pop()
		vm.stack.Push(a != b)
		return compiler.OPCODE_TOTAL_BYTES, nil
	}
	return 0, RuntimeError{Message: fmt.Sprintf("unknown equality opcode %v", opCode)}
}

// execComparisonInstruction executes a comparison or equality operation based on the provided opcode.
func (vm *VirtualMachine) execComparisonInstruction(opCode compiler.Opcode) (int, error) {

	handler, ok := vm.comparisonOpHandlers[opCode]
	if !ok {
		return 0, RuntimeError{Message: fmt.Sprintf("unknown comparison opcode %v", opCode)}
	}
	if err := handler(vm, nil, nil); err != nil {
		return 0, err
	}
	return compiler.OPCODE_TOTAL_BYTES, nil
}

// execLogicalInstruction executes logical operations (AND, OR) on the VM's stack.
func (vm *VirtualMachine) execLogicalInstruction(opCode compiler.Opcode) (int, error) {
	if opCode == compiler.OP_AND {
		b := vm.stack.Pop()
		a := vm.stack.Pop()

		if isBool(a) && isBool(b) {
			vm.stack.Push(a.(bool) && b.(bool))
		} else {
			return 0, RuntimeError{Message: "operands must be boolean values"}
		}
	}
	if opCode == compiler.OP_OR {
		b := vm.stack.Pop()
		a := vm.stack.Pop()
		if isBool(a) && isBool(b) {
			vm.stack.Push(a.(bool) || b.(bool))
		} else {
			return 0, RuntimeError{Message: "operands must be boolean values"}
		}
	}
	return compiler.OPCODE_TOTAL_BYTES, nil

}

// execDefineGlobalInstruction defines a global variable, and assigns the corresponding
// value from the top of the stack to it.
func (vm *VirtualMachine) execDefineGlobalInstruction(bytecode compiler.Bytecode) int {
	index := vm.ip + compiler.OPCODE_TOTAL_BYTES
	instruction := bytecode.Instructions[index : vm.ip+compiler.OP_CONSTANT_TOTAL_BYTES]
	operand := binary.BigEndian.Uint16(instruction)
	name := bytecode.NameConstants[operand]

	vm.globalVars[name] = vm.stack.Pop()
	return compiler.OP_CONSTANT_TOTAL_BYTES
}

// execSetGlobalInstruction sets the value of an existing global variable
func (vm *VirtualMachine) execGetGlobalInstruction(bytecode compiler.Bytecode) int {
	index := vm.ip + compiler.OPCODE_TOTAL_BYTES
	instruction := bytecode.Instructions[index : vm.ip+compiler.OP_CONSTANT_TOTAL_BYTES]
	operand := binary.BigEndian.Uint16(instruction)
	name := bytecode.NameConstants[operand]
	vm.stack.Push(vm.globalVars[name])
	return compiler.OP_CONSTANT_TOTAL_BYTES
}

// Executes a unary operations and pushes the result onto the VM's stack.
func (vm *VirtualMachine) execUnaryInstruction(opCode compiler.Opcode) (int, error) {
	value := vm.stack.Pop()
	if value == nil {
		return 0, RuntimeError{Message: "stack underflow on unary operation"}
	}

	if opCode == compiler.OP_NEGATE {
		if isFloat(value) {
			val, err := literalToFloat64(value)
			if err != nil {
				return 0, RuntimeError{Message: err.Error()}
			}
			vm.stack.Push(-val)
			return compiler.OPCODE_TOTAL_BYTES, nil
		}

		val, err := literalToInt64(value)
		if err != nil {
			return 0, RuntimeError{Message: err.Error()}
		}
		vm.stack.Push(-val)

	}

	if opCode == compiler.OP_NOT {
		switch v := value.(type) {
		case bool:
			vm.stack.Push(!v)
		default:
			// non-nil, non-falsey values are truthy -> !truthy == false
			vm.stack.Push(false)
		}
	}

	return compiler.OPCODE_TOTAL_BYTES, nil
}

// Fetches and pushes a constant value from the bytecode
// onto the VM's stack.
//
// It reads the operand following the OP_CONSTANT opcode to locate the
// appropriate entry in the constants pool, retrieves that value, and pushes it
// onto the stack for subsequent instruction execution.
//
// Parameters:
//   - bytecode: The compiled sequence of instructions containing both opcodes
//     and constant pool references.
//
// Returns:
//   - int: The total number of bytes consumed by this instruction, used to
//     increment the VM's instruction pointer.
func (vm *VirtualMachine) execConstantInstruction(bytecode compiler.Bytecode) int {
	index := vm.ip + compiler.OPCODE_TOTAL_BYTES
	instruction := bytecode.Instructions[index : vm.ip+compiler.OP_CONSTANT_TOTAL_BYTES]
	operand := binary.BigEndian.Uint16(instruction)
	value := bytecode.ConstantsPool[operand]
	vm.stack.Push(value)
	return compiler.OP_CONSTANT_TOTAL_BYTES
}

// Executes an arithmetic operation on the VM's stack
// based on the operand types and provided arithmetic functions.
//
// It pops two operands from the stack, determines whether they are integers
// or floats, and applies the corresponding arithmetic function.
// The result is then pushed back onto the stack.
//
// Parameters:
//   - operationFloat: Function handling arithmetic between floating-point values.
//   - operationInt:   Function handling arithmetic between integer values.
//   - opCode:         The opcode representing the arithmetic operation.
//
// Returns:
//   - int: The number of bytes consumed by the instruction, used to advance the
//     instruction pointer.
//   - error: A RuntimeError if operand types are invalid, otherwise nil.
func (vm *VirtualMachine) execArithmeticInstruction(operationFloat arithmeticFuncFloat, operationInt arithmeticFuncInt, opCode int) (int, error) {
	b := vm.stack.Pop()
	a := vm.stack.Pop()

	if a != nil && b != nil {
		var aFloatVal float64
		var aIntVal int64
		var bFloatVal float64
		var bIntVal int64
		isAFloat := isFloat(a)
		isBFloat := isFloat(b)
		isAInt := isInt(a)
		isBInt := isInt(b)

		if isAFloat {
			val, _ := literalToFloat64(a)
			aFloatVal = val
		}
		if isBFloat {
			val, _ := literalToFloat64(b)
			bFloatVal = val
		}
		if isAInt {
			val, _ := literalToInt64(a)
			aIntVal = val
		}
		if isBInt {
			val, _ := literalToInt64(b)
			bIntVal = val
		}

		if !isAFloat && !isBFloat && !isAInt && !isBInt {
			message := fmt.Sprintf("operands must be numeric values: %v,%v", a, b)
			return 0, RuntimeError{Message: message}
		}

		if isAFloat && isBFloat {
			result := operationFloat(aFloatVal, bFloatVal)
			vm.stack.Push(result)
		}
		if isAFloat && isBInt {
			result := operationFloat(aFloatVal, float64(bIntVal))
			vm.stack.Push(result)
		}
		if isAInt && isBFloat {
			result := operationFloat(float64(aIntVal), bFloatVal)
			vm.stack.Push(result)
		}
		if isAInt && isBInt {
			if opCode == int(compiler.OP_DIVIDE) {
				result := operationFloat(float64(aIntVal), float64(bIntVal))
				vm.stack.Push(result)
			} else {
				result := operationInt(aIntVal, bIntVal)
				vm.stack.Push(result)
			}
		}
	}

	return compiler.OPCODE_TOTAL_BYTES, nil
}
