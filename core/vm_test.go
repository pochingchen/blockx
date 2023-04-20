package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value := s.Pop()
	assert.Equal(t, value, 1)

	value = s.Pop()
	assert.Equal(t, value, 2)
}

func TestVM(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	state := NewState()
	vm := NewVM(data, state)
	assert.Nil(t, vm.Run())

	valueBytes, err := state.Get([]byte("FOO"))
	assert.Nil(t, err)
	value := DeserializeInt64(valueBytes)
	assert.Equal(t, value, int64(5))
}
