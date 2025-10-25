package vm

import (
	"encoding/binary"
	"fmt"
	"nilan/compiler"
)

type arithmeticFuncFloat func(a float64, b float64) float64
type arithmeticFuncInt func(a int64, b int64) int64

func addFloat(a float64, b float64) float64{
	return a+b
}
func addInt(a int64, b int64) int64{
	return a+b
}
func subFloat(a float64, b float64) float64{
	return a-b
}
func subInt(a int64, b int64) int64{
	return a-b
}
func multFloat(a float64, b float64) float64{
	return a*b
}
func multInt(a int64, b int64) int64{
	return a*b
}
func divFloat(a float64, b float64) float64{
	// TODO: Add runtime error division by zero
	return a/b
}
func divInt(a int64, b int64) int64{
	// TODO: Add runtime error division by zero
	return a/b
}

func literalToInt64(value any) (int64,error){

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

// literalToFloat64 attempts to convert a literal value into a float64.
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

// Represents a stack based virtual-machine (VirtualMachine).
// It is the runtime environment where Nilan bytecode
// gets executed.
type VirtualMachine struct {
	stack Stack
	ip    int
	debug bool
}

// Creates a new VM instance
func New() *VirtualMachine {
	return &VirtualMachine{debug: true}
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

		switch opCode {
		case compiler.OP_END:
			fmt.Println(vm.stack.Peek()) // temporary code just for viz
			return nil
		case compiler.OP_CONSTANT:
			instructionLength = vm.execConstantInstruction(bytecode)
		case compiler.OP_ADD:
			l,err := vm.execArithmeticInstruction(addFloat,addInt)
			if err!=nil{
				return err
			}
			instructionLength = l
		case compiler.OP_SUBTRACT:
			l,err := vm.execArithmeticInstruction(subFloat,subInt)
			if err!=nil{
				return err
			}
			instructionLength = l			
		case compiler.OP_MULTIPLY:
			l,err := vm.execArithmeticInstruction(multFloat,multInt)
			if err!=nil{
				return err
			}
			instructionLength = l
		case compiler.OP_DIVIDE:
			l,err := vm.execArithmeticInstruction(divFloat,divInt)
			if err!=nil{
				return err
			}
			instructionLength = l
		default:
			// NOTE: This should only happen in development mode.
			return fmt.Errorf("unknown opcode %v at ip %d", opCode, vm.ip)
		}

		vm.ip += instructionLength
	}
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
func (vm *VirtualMachine) execConstantInstruction(bytecode compiler.Bytecode) int{
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
//
// Returns:
//   - int: The number of bytes consumed by the instruction, used to advance the
//     instruction pointer.
//   - error: A RuntimeError if operand types are invalid, otherwise nil.
func (vm *VirtualMachine) execArithmeticInstruction(operationFloat arithmeticFuncFloat,operationInt arithmeticFuncInt) (int, error) {
	b:=vm.stack.Pop()
	a:= vm.stack.Pop()

	if a!=nil && b!=nil{
		bIntVal,aErr := literalToInt64(b)
		aIntVal, bErr := literalToInt64(a)
		
		if aErr!=nil && bErr != nil{
			bFloatVal,aFErr := literalToFloat64(b)
			aFloatVal, bFErr := literalToFloat64(a)

			if aFErr == nil && bFErr == nil{

				result:= operationFloat(aFloatVal,bFloatVal)
				vm.stack.Push(result)
			}else{
				message := fmt.Sprintf("operands must be numeric values: %v,%v",a,b)
				return 0, RuntimeError{Message: message}
			}

		}
		
		result :=operationInt(aIntVal,bIntVal)
		vm.stack.Push(result)
	}

	return compiler.OPCODE_TOTAL_BYTES,nil
}
