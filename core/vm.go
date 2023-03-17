package core

import "projectx/types"

type Instruction byte

var (
	InstrPush  Instruction = 0x0a
	InstrAdd   Instruction = 0x0b
	InstrMinus Instruction = 0x0c
	InstrMul   Instruction = 0x0d
	InstrDiv   Instruction = 0x0e
)

type VM struct {
	data  []byte
	ip    int //instruction pointer
	stack *types.Stack[byte]
	sp    int //stack pointer
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: types.NewStack[byte](),
		sp:    -1,
	}
}

func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])

		if err := vm.Exec(instr); err != nil {
			return err
		}

		vm.ip++

		if vm.ip > len(vm.data)-1 {
			break
		}

	}
	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrPush:
		vm.stack.Push(vm.data[vm.ip-1])
	case InstrAdd:
		n1 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		n2 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		vm.stack.Push(n1 + n2)
		vm.sp++
	case InstrMinus:
		n1 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		n2 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		vm.stack.Push(n2 - n1)
		vm.sp++
	case InstrMul:
		n1 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		n2 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		vm.stack.Push(n1 * n2)
		vm.sp++
	case InstrDiv:
		n1 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		n2 := vm.stack.Top()
		vm.stack.Pop()
		vm.sp--
		vm.stack.Push(n2 / n1)
		vm.sp++
	}

	return nil
}
