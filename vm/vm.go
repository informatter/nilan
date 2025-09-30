package vm

import (
	"encoding/binary"
	"fmt"
	"nilan/compiler"
)

// Represents a stack based virtual-machine (VM).
// It is the runtime environment where Nilan bytecode
// gets executed.
type VM struct {
	stack Stack
	ip    int
	debug bool
}

// Creates a new VM instance
func New() *VM {
	return &VM{debug: true}
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
func (vm *VM) Run(bytecode compiler.Bytecode) error {

	var instructionLength int
	for {
		opCode := compiler.Opcode(bytecode.Instructions[vm.ip])

		switch opCode {
		case compiler.OP_END:
			return nil
		case compiler.OP_CONSTANT:

			index := vm.ip + compiler.OPCODE_TOTAL_BYTES
			s := bytecode.Instructions[index : vm.ip+compiler.OP_CONSTANT_TOTAL_BYTES]
			operand := binary.BigEndian.Uint16(s)
			value := bytecode.ConstantsPool[operand]
			vm.stack.Push(value)
			instructionLength = compiler.OP_CONSTANT_TOTAL_BYTES
		default:
			// NOTE: This should only happen in development mode.
			return fmt.Errorf("unknown opcode %v at ip %d", opCode, vm.ip)
		}

		vm.ip += instructionLength

	}
}
