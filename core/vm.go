package core

import (
	"encoding/binary"
)

type Instruction byte

const (
	InstrPushInt  Instruction = 0x0a // 10
	InstrAdd      Instruction = 0x0b // 11
	InstrPushByte Instruction = 0x0c // 12
	InstrPack     Instruction = 0x0d // 13
	InstrSub      Instruction = 0x0e // 14
	InstrStore    Instruction = 0x0f // 15
)

type Stack struct {
	data []any
	sp   int
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

func (s *Stack) Push(v any) {
	s.data[s.sp] = v
	s.sp++
}

func (s *Stack) Pop() any {
	value := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)
	s.sp--

	return value
}

type VM struct {
	data          []byte
	ip            int // instruction pointer
	stack         *Stack
	contractState *State
}

func NewVM(data []byte, s *State) *VM {
	return &VM{
		contractState: s,
		data:          data,
		ip:            0,
		stack:         NewStack(128),
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
	case InstrStore:
		var (
			key             = vm.stack.Pop().([]byte)
			value           = vm.stack.Pop()
			serializedValue []byte
		)

		switch v := value.(type) {
		case int:
			serializedValue = SerializeInt64(int64(v))
		default:
			panic("todo: unknown type")
		}

		vm.contractState.Put(key, serializedValue)
	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1]))

	case InstrPushByte:
		vm.stack.Push(vm.data[vm.ip-1])

	case InstrPack:
		n := vm.stack.Pop().(int)
		b := make([]byte, n)

		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)

	case InstrSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a - b
		vm.stack.Push(c)

	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a + b
		vm.stack.Push(c)
	}

	return nil
}

func SerializeInt64(value int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))

	return buf
}

func DeserializeInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
