package compiler

// This file implements the ASTCompiler, which compiles the abstract syntax tree (AST) directly to bytecode.

import (
	"encoding/binary"
	"fmt"
	"nilan/ast"
	"nilan/token"
	"os"
	"strings"
)

// Local represents a local variable in the compiler.
// NOTE/TODO: The struct layout can probably be optimised by packing the fields differently.
// So the struct has better cache locality and takes up less memory.
type Local struct {

	// The variable's name
	name string
	// The variable's depth in the scope stack. Used to determine when variables go out of scope.
	depth uint16
	// Whether the variable has been initialized. Used to prevent accessing uninitialized variables.
	initialized bool
	// The slot index where the variable is stored. Used for local variable access in the VM.
	slot uint16
}

// ASTCompiler is a visitor that compiles AST nodes directly to bytecode.
// It implements both ast.ExpressionVisitor and ast.StmtVisitor interfaces
// to traverse and compile the abstract syntax tree to bytecode.
type ASTCompiler struct {

	// The resulting compiled bytecode.
	bytecode Bytecode
	// Tracks initialized global variables
	initialized map[string]bool
	// A stack of local variables in the current scope. Used for local variable management and access.
	// Locals are orderd by by their declaration order that appears in the code. The most recently declared variable
	// will always be at the top of the stack.
	// TODO: We can re-factor the `Stack` implementation in the VM package so it can be used here. We should move that implementation
	// to a new package.
	locals []Local
	// The current depth of nested scopes. Used to determine when local variables go out of scope.
	scopeDepth uint16
}

// NewASTCompiler creates a new AST-to-bytecode compiler.
func NewASTCompiler() *ASTCompiler {
	return &ASTCompiler{
		bytecode: Bytecode{
			Instructions:  Instructions{},
			ConstantsPool: []any{},
			NameConstants: []string{},
		},
		initialized: make(map[string]bool),
		locals:      []Local{},
		scopeDepth:  0,
	}
}

// DumpBytecode writes the compiled bytecode to a file with a `.nic` extension.
// The bytecode is encoded as hexadecimal so it can be viewed in a text editor.
func (ac *ASTCompiler) DumpBytecode(filePath string) error {
	if filePath == "" {
		filePath = "bytecode.nic"
	} else {
		filePath = filePath + ".nic"
	}
	fDescriptor, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating nilan bytecode file: %s", err.Error())
	}

	encoded := fmt.Sprintf("%x", ac.bytecode.Instructions)
	fDescriptor.Write([]byte(encoded))
	defer fDescriptor.Close()
	return nil
}

// DiassembleBytecode disassembles the compiled bytecode to a human readable format
// and optionally saves it to disk.
// It returns the disassembled bytecode as a string or an error if the file could not be created.
func (ac *ASTCompiler) DiassembleBytecode(saveToDisk bool, filePath string) (string, error) {
	var diassembledBytecode string
	var builder strings.Builder
	var instructionLength int
	totalInstructions := len(ac.bytecode.Instructions) - 1
	ip := 0

	// NOTE: Slicing in go includes the first element, but excludes the last one.
	// for example, [0:4] will include index 0 to index 3 of the array.

	for ip <= totalInstructions {
		opCode := Opcode(ac.bytecode.Instructions[ip])
		switch opCode {
		case OP_ADD, OP_LESS, OP_LARGER, OP_PRINT, OP_SUBTRACT, OP_DIVIDE,
			OP_MULTIPLY, OP_NEGATE, OP_NOT, OP_AND, OP_OR,
			OP_EQUALITY, OP_NOT_EQUAL, OP_LARGER_EQUAL, OP_LESS_EQUAL,
			OP_END, OP_POP:

			result, err := DiassembleInstruction([]byte{ac.bytecode.Instructions[ip]})
			if err != nil {
				panic(err.Error())
			}
			builder.WriteString(result)
			if opCode == OP_END {
				break
			}
			builder.WriteString("\n")
			instructionLength = OPCODE_TOTAL_BYTES

		case OP_GET_LOCAL, OP_SET_LOCAL:
			// The  operand is the index where the local variable is stored in the VM's stack.
			operand, dia := ac.diassemble3ByteInstruction(ip)
			result := dia + fmt.Sprintf(", vm stack index: %d", operand)
			builder.WriteString(result)
			builder.WriteString("\n")
			instructionLength = THREE_BYTE_INSTRUCTION_LENGTH

		case OP_SCOPE_EXIT:
			operand, dia := ac.diassemble3ByteInstruction(ip)
			result := dia + fmt.Sprintf(", total local variables to pop from the VM's stack: %d", operand)
			builder.WriteString(result)
			builder.WriteString("\n")
			instructionLength = THREE_BYTE_INSTRUCTION_LENGTH

		// Handles all opcodes which store data in the constants pool.
		// all these opcodes have an operand (index into constants pool) with a width of 2 bytes.
		case OP_CONSTANT, OP_SET_GLOBAL, OP_GET_GLOBAL:

			// The operand is the index into the constants pool where the actual value is stored.
			operand, dia := ac.diassemble3ByteInstruction(ip)
			value := ac.bytecode.ConstantsPool[operand]
			result := dia + fmt.Sprintf(", value: %d", value)
			builder.WriteString(result)
			builder.WriteString("\n")
			instructionLength = THREE_BYTE_INSTRUCTION_LENGTH

		case OP_JUMP, OP_JUMP_IF_FALSE:

			operand, dia := ac.diassemble3ByteInstruction(ip)
			result := dia + fmt.Sprintf(", byte index in instruction array: %d", operand)
			builder.WriteString(result)
			builder.WriteString("\n")
			instructionLength = THREE_BYTE_INSTRUCTION_LENGTH

		}

		ip += instructionLength
	}
	diassembledBytecode = builder.String()
	if saveToDisk {
		if filePath == "" {
			filePath = "bytecode.dnic"
		} else {
			filePath = filePath + ".dnic"
		}
		fDescriptor, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("error creating diassembled bytecode file: %s", err.Error())
		}
		fDescriptor.WriteString(diassembledBytecode)
		defer fDescriptor.Close()
	}
	return diassembledBytecode, nil
}

func (ac *ASTCompiler) CompileAST(statements []ast.Stmt) (b Bytecode, err error) {
	// Recover from any panic that may occur during compilation
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case SemanticError:
				err = v
			case DeveloperError:
				err = v
			}
		}
	}()

	// If previous compilation left an OP_END at the end, drop it
	if len(ac.bytecode.Instructions) > 0 {
		if ac.bytecode.Instructions[len(ac.bytecode.Instructions)-1] == byte(OP_END) {
			ac.bytecode.Instructions = ac.bytecode.Instructions[:len(ac.bytecode.Instructions)-1]
		}
	}

	for _, stmt := range statements {
		func() {
			//NOTE: Catch panics per statement to avoid aborting the whole loop
			defer func() {
				if r := recover(); r != nil {
					panic(r)
				}
			}()
			stmt.Accept(ac)
		}()
	}

	ac.emit(OP_END)
	return ac.bytecode, nil
}

// VisitBinary handles binary expressions (arithmetic operators: +, -, *, /)
func (ac *ASTCompiler) VisitBinary(binary ast.Binary) any {

	// NOTE: Left expression is compiled first to ensure correct evaluation order
	binary.Left.Accept(ac)
	binary.Right.Accept(ac)

	switch binary.Operator.TokenType {
	case token.ADD:
		ac.emit(OP_ADD)
	case token.SUB:
		ac.emit(OP_SUBTRACT)
	case token.MULT:
		ac.emit(OP_MULTIPLY)
	case token.DIV:
		ac.emit(OP_DIVIDE)

	case token.EQUAL_EQUAL:
		ac.emit(OP_EQUALITY)
	case token.LARGER:
		ac.emit(OP_LARGER)
	case token.LESS:
		ac.emit(OP_LESS)
	case token.LESS_EQUAL:
		ac.emit(OP_LESS_EQUAL)
	case token.LARGER_EQUAL:
		ac.emit(OP_LARGER_EQUAL)
	case token.NOT_EQUAL:
		ac.emit(OP_NOT_EQUAL)
	}

	return nil
}

// VisitUnary handles unary expressions (operators: -, !)
func (ac *ASTCompiler) VisitUnary(unary ast.Unary) any {

	unary.Right.Accept(ac)

	switch unary.Operator.TokenType {
	case token.SUB:
		ac.emit(OP_NEGATE)
	case token.BANG:
		ac.emit(OP_NOT)
	}
	return nil
}

// VisitLiteral handles literal values (numbers, strings, booleans, null)
// Adds the literal value to the constants pool.
func (ac *ASTCompiler) VisitLiteral(literal ast.Literal) any {
	ac.addConstant(literal.Value)
	return nil
}

// VisitGrouping handles parenthesized expressions
func (ac *ASTCompiler) VisitGrouping(grouping ast.Grouping) any {
	// Recursively compile the inner expression
	grouping.Expression.Accept(ac)
	return nil
}

// VisitVariableExpression compiles variable access by emitting bytecode to load the variable's
// value onto the VM's stack.
//
// For local variabables, it emites an OP_GET_LOCAL instruction with the variable's slot index as the operand.
//
// For global variables, it emits an OP_GET_GLOBAL instruction with the variable's index in the NameConstants pool as the operand.
//
// For example, this compiles code such as `x` or `y` by emitting the appropriate instruction to get
// the variable's value from the VM's stack.
func (ac *ASTCompiler) VisitVariableExpression(variable ast.Variable) any {

	identifier := variable.Name.Lexeme

	slotIndex := ac.resolveLocal(identifier)
	if slotIndex != -1 {
		if !ac.locals[slotIndex].initialized {
			panic(SemanticError{
				Message: fmt.Sprintf("Cant access uninitialised variable '%s'", identifier),
			})
		}
		ac.emit(OP_GET_LOCAL, slotIndex)
		return nil
	}

	globalIndex := ac.resolveGlobal(identifier)
	if globalIndex == -1 {
		panic(SemanticError{
			Message: fmt.Sprintf("name '%s' is not defined", identifier),
		})
	}
	if !ac.initialized[identifier] {
		panic(SemanticError{
			Message: fmt.Sprintf("Cant access uninitialised variable '%s'", identifier),
		})
	}

	ac.emit(OP_GET_GLOBAL, globalIndex)
	return nil
}

// VisitAssignExpression compiles an assignment expression by first compiling the right-hand side expression,
// and then attempting to resolve the variable name as local or global.
//
// For local variables, it emits an OP_SET_LOCAL instruction with the variable's slot index as the operand.
//
// For global variables, it emits an OP_SET_GLOBAL instruction with the variable's index in the NameConstants pool as the operand.
//
// For exmaple, this compiles code such as `x = 5` or `y = x + 2` by first compiling the right hand side expression
// (`5` or `x + 2`), then emitting the appropriate instruction to store the value in the corresponding variable.
func (ac *ASTCompiler) VisitAssignExpression(assign ast.Assign) any {

	name := assign.Name.Lexeme

	// compile the right hand side expression first.
	// This ensures that the correct value is on top of the stack when the OP_SET_LOCAL
	// or OP_SET_GLOBAL instruction is emitted.
	assign.Value.Accept(ac)

	slotIndex := ac.resolveLocal(name)
	if slotIndex != -1 {
		ac.locals[slotIndex].initialized = true
		ac.emit(OP_SET_LOCAL, slotIndex)
		return nil
	}

	globalIndex := ac.resolveGlobal(name)
	if globalIndex == -1 {
		panic(SemanticError{
			Message: fmt.Sprintf("name '%s' is not defined", name),
		})
	}

	ac.initialized[name] = true
	ac.emit(OP_SET_GLOBAL, globalIndex)
	return nil
}

// VisitVarStmt handles variable declaration statements.
//
// For global variables, it adds the variable name to the NameConstants pool and
// emits an OP_SET_GLOBAL instruction.
//
// For local variables it declares the variable in the current scope and emits an OP_SET_LOCAL instruction.
//
// For example, this compiles code such as `var x = 5`,  `var y`, var z = 10+2` ... etc
func (ac *ASTCompiler) VisitVarStmt(varStmt ast.VarStmt) any {

	variableName := varStmt.Name.Lexeme
	if ac.scopeDepth == 0 {
		// Handles global variable declaration.
		index := ac.addNameConstant(variableName)
		if varStmt.Initializer != nil {
			varStmt.Initializer.Accept(ac)
			ac.emit(OP_SET_GLOBAL, index)
		}
		ac.initialized[variableName] = varStmt.Initializer != nil
	} else {
		// Handles local variable declaration.
		ac.declareLocal(variableName)
		if varStmt.Initializer != nil {
			varStmt.Initializer.Accept(ac)
		} else {
			ac.addConstant(nil)
		}
		slot := ac.locals[len(ac.locals)-1].slot
		ac.emit(OP_SET_LOCAL, int(slot))
		ac.locals[len(ac.locals)-1].initialized = varStmt.Initializer != nil
	}

	return nil
}

// VisitLogicalExpression compiles logical expressions (and, or) by emitting bytecode that implements short-circuiting behaviour.
func (ac *ASTCompiler) VisitLogicalExpression(logical ast.Logical) any {

	// left expression is compiled first to ensure correct evaluation order and short-circuiting behaviour.
	logical.Left.Accept(ac)

	switch logical.Operator.TokenType {
	case token.OR:
		// For an "or" expression, if the left operand is truthy, we want to short-circuit and skip
		// evaluating the right operand.

		jumpIfFalsePos := ac.emitPlaceholderJump(OP_JUMP_IF_FALSE)
		jumpEndPos := ac.emitPlaceholderJump(OP_JUMP)

		rightStart := len(ac.bytecode.Instructions)
		ac.patchJump(jumpIfFalsePos, rightStart)

		ac.emit(OP_POP)

		// The right expression is compiled after emitting the jump instruction. If the left operand is truthy,
		// the VM will jump over the right expression. This is achieved by the below patchJump call.
		logical.Right.Accept(ac)

		ac.patchJump(jumpEndPos, len(ac.bytecode.Instructions))
	case token.AND:
		// For an "and" expression, if the left operand is falsy, we want to short-circuit and skip evaluating the right operand.
		jumpIfFalsePos := ac.emitPlaceholderJump(OP_JUMP_IF_FALSE)

		ac.emit(OP_POP)
		logical.Right.Accept(ac)

		ac.patchJump(jumpIfFalsePos, len(ac.bytecode.Instructions))
	}
	return nil
}

// VisitExpressionStmt is not directly called; handled by CompileAST
func (ac *ASTCompiler) VisitExpressionStmt(exprStmt ast.ExpressionStmt) any {
	exprStmt.Expression.Accept(ac)
	return nil
}

func (ac *ASTCompiler) VisitPrintStmt(printStmt ast.PrintStmt) any {
	printStmt.Expression.Accept(ac)
	ac.emit(OP_PRINT)
	return nil
}

// VisitBlockStmt compiles a block statement by sequentially compiling each statement
// in the block.
func (ac *ASTCompiler) VisitBlockStmt(blockStmt ast.BlockStmt) any {

	ac.beginScope()
	for _, stmt := range blockStmt.Statements {
		func() {
			//NOTE: Catch panics per statement to avoid aborting the whole loop
			defer func() {
				if r := recover(); r != nil {
					panic(r)
				}
			}()
			stmt.Accept(ac)
		}()
	}

	popped := ac.endScope()
	if popped > 0 {
		ac.emit(OP_SCOPE_EXIT, popped)
	}
	return nil
}

// VisitIfStmt compiles an if or if-else statement by emitting bytecode.
// It uses backpatching to resolve jump offsets for branching.
func (ac *ASTCompiler) VisitIfStmt(ifStmt ast.IfStmt) any {

	// compile the condition expression first
	ifStmt.Condition.Accept(ac)

	jumpIfFalsePatch := ac.emitPlaceholderJump(OP_JUMP_IF_FALSE)
	// For example, the intructions would now be something like: [..., OP_JUMP_IF_FALSE,  0x00, 0x00]
	// where `0x00, 0x0` are the placeholder operand bytes.

	ifStmt.Then.Accept(ac)

	if ifStmt.Else != nil {
		// If there is an "else" branch, emit a jump instruction to skip over it after executing the "then" branch.
		jumpPatch := ac.emitPlaceholderJump(OP_JUMP)

		// Patch the operand of the OP_JUMP_IF_FALSE instruction defined at the beginning.
		// This allows the VM to correctly jump to the start of the "else" branch, if the "then"
		// branch condition evaluates false.
		elsePos := len(ac.bytecode.Instructions)
		ac.patchJump(jumpIfFalsePatch, elsePos)

		ifStmt.Else.Accept(ac)

		endPos := len(ac.bytecode.Instructions)
		// Patch the operand of `OP_JUMP` so the VM can jump to the end of the "else" branch.
		ac.patchJump(jumpPatch, endPos)
	} else {
		// If there is no "else" branch, patch the OP_JUMP_IF_FALSE so that
		// control jumps to the instruction after the "then" branch when
		// the condition is false.
		afterPos := len(ac.bytecode.Instructions)
		ac.patchJump(jumpIfFalsePatch, afterPos)
	}
	// Emits `OP_POP` so the VM can pop the condition expression's value from the stack.
	ac.emit(OP_POP)
	return nil
}

func (ac *ASTCompiler) VisitWhileStmt(whileStmt ast.WhileStmt) any {

	loopstartPos := len(ac.bytecode.Instructions)

	// compile the condition expression first
	whileStmt.Condition.Accept(ac)

	jumpIfFalsePatch := ac.emitPlaceholderJump(OP_JUMP_IF_FALSE)

	// compile the loop body
	whileStmt.Body.Accept(ac)

	// After compiling the loop body, we need to emit a jump instruction
	// so the VM can jump back to the start of the loop condition.
	ac.emit(OP_POP)
	ac.emit(OP_JUMP, loopstartPos)

	// if the while condition is false, the VM needs to jump to the end of the loop body,
	// which is the current position in the instruction array.
	loopEndPos := len(ac.bytecode.Instructions)
	ac.patchJump(jumpIfFalsePatch, loopEndPos)
	ac.emit(OP_POP)

	return nil
}

// patchjump overwrites a jump instruction's operand with the actual correct byte offset.
// When compiling if statements, its not possible to know the else branch (or the statement after
// the if) will be until the then-branch is compiled. Jump instructions are emmited with placeholder operands,
// then later call patchJump to fix those operands.

// The jumpPos is the byte index where the jump instruction's OPCODE is located.
//
//	This is the position BEFORE the jump was emitted
//
// The targetPos is the byte index where the jump instruction should jump to.
// Example:
// jumpPos = 10, targetPos = 20
// Before patching: [..., OP_JUMP_IF_FALSE, 0x00, 0x00, ...] (jump instruction starts at index 10)
// After patching: [..., OP_JUMP_IF_FALSE, 0x00, 0x0A, ...] (jump instruction now correctly jumps to index 20)
func (ac *ASTCompiler) patchJump(jumpPos int, targetPos int) {

	operandPos := jumpPos + OPCODE_TOTAL_BYTES

	instruction := make([]byte, 2)
	binary.BigEndian.PutUint16(instruction, uint16(targetPos))

	// override the 2-byte placeholder operand in the instruction array with
	// the correct operand bytes that will make the jump instruction jump to the target position.
	ac.bytecode.Instructions[operandPos] = instruction[0]
	ac.bytecode.Instructions[operandPos+1] = instruction[1]

}

// addConstant appends a value to the constant pool and emits an OP_CONSTANT instruction.
// The operand of the instruction will be its index in the constants pool.
func (ac *ASTCompiler) addConstant(value any) {
	ac.bytecode.ConstantsPool = append(ac.bytecode.ConstantsPool, value)
	index := len(ac.bytecode.ConstantsPool) - 1
	ac.emit(OP_CONSTANT, index)
}

// addNameConstant adds a variable name to the NameConstants pool
// and returns its index.
func (ac *ASTCompiler) addNameConstant(value string) int {

	for _, name := range ac.bytecode.NameConstants {
		if name == value {
			panic(SemanticError{
				Message: fmt.Sprintf("Redefinition of variable '%s'", value),
			})
		}
	}
	ac.bytecode.NameConstants = append(ac.bytecode.NameConstants, value)
	return len(ac.bytecode.NameConstants) - 1
}

// emit constructs a bytecode instruction and appends it to the instruction stream
func (ac *ASTCompiler) emit(opcode Opcode, operands ...int) {
	instruction, err := AssembleInstruction(opcode, operands...)
	if err != nil {
		// TODO: Improve error handling in compiler.
		// Although in this case its can be OK as the error returned is of type `DeveloperError`
		// which would only be raised during development.
		panic(err.Error())
	}
	ac.bytecode.Instructions = append(ac.bytecode.Instructions, instruction...)
}

// emitPlaceholderJump emits a jump instruction with the specified opcode and a placeholder operand (0).
// It returns the position in the bytecode where the jump instruction was emitted,
// which can later be passed to `patchJump` to update the operand with
// the correct jump target.
func (ac *ASTCompiler) emitPlaceholderJump(opcode Opcode) int {
	position := len(ac.bytecode.Instructions)
	ac.emit(opcode, 0)
	return position
}

// beginScope increments the scope depth, when compiling a block statement.
func (ac *ASTCompiler) beginScope() {
	ac.scopeDepth++
}

// endScope decrements the scope depth and removes any local variables that go out of scope.
// It returns the number of local variables that went out of scope,
// which is used by the VM to pop them from the stack.
func (ac *ASTCompiler) endScope() int {
	ac.scopeDepth--

	count := 0
	for len(ac.locals) > 0 && ac.locals[len(ac.locals)-1].depth > ac.scopeDepth {
		ac.locals = ac.locals[:len(ac.locals)-1]
		count++
	}

	return count
}

// declareLocal adds a local variable name, checking for same-scope duplicates
// and assigns it a slot index for the VM to access it.
// It panics if there is a duplicate variable declaration in the same scope.
func (ac *ASTCompiler) declareLocal(name string) {

	for i := len(ac.locals) - 1; i >= 0; i-- {

		// By virtue of iterating backwards through the local stack,
		// we can stop checking
		if ac.locals[i].depth < ac.scopeDepth {
			break
		}
		if ac.locals[i].name == name {
			panic(SemanticError{
				Message: fmt.Sprintf("Redefinition of variable '%s'", name),
			})
		}
	}

	slot := uint16(len(ac.locals))
	local := Local{
		name:        name,
		depth:       ac.scopeDepth,
		initialized: false,
		slot:        slot,
	}
	ac.locals = append(ac.locals, local)

}

// defineLocal marks the most recently declared local variable as initialized.
func (ac *ASTCompiler) defineLocal() {
	if len(ac.locals) > 0 {
		ac.locals[len(ac.locals)-1].initialized = true
	}
}

// resolveLocal checks if a variable name exists in the current local scope and returns its slot index.
// It returns -1 if the variable is not found in the local scope.
func (ac *ASTCompiler) resolveLocal(name string) int {
	for i := len(ac.locals) - 1; i >= 0; i-- {
		if ac.locals[i].name == name {
			return int(ac.locals[i].slot)
		}
	}
	return -1
}

// resolveGlobal checks if a variable name exists in the global scope and returns its index in the NameConstants pool.
// It returns -1 if the variable is not found in the global scope.
func (ac ASTCompiler) resolveGlobal(name string) int {
	for i, n := range ac.bytecode.NameConstants {
		if n == name {
			return i
		}
	}
	return -1
}

// diassemble3ByteInstruction reads a 3-byte instruction starting at the instruction pointer(ip),
// in the bytecodes instruction array. IT interprets the final two bytes as a big-endian uint16 operand,
// and returns it along with the textual disassembly produced by DiassembleInstruction.
// A panic is raised if DiassembleInstruction returns an error.
func (ac *ASTCompiler) diassemble3ByteInstruction(ip int) (uint16, string) {
	offset := ip + 3
	instruction := ac.bytecode.Instructions[ip:offset]
	operand := binary.BigEndian.Uint16(instruction[OPCODE_TOTAL_BYTES:])
	dia, err := DiassembleInstruction(instruction)
	if err != nil {
		panic(err.Error())
	}

	return operand, dia
}
