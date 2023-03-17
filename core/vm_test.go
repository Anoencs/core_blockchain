package core

import (
	"projectx/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x9, 0x0a, 0x2, 0x0a}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	assert.Equal(t, vm.stack, types.NewStack[uint8](0x1, 0x2, 0x9, 0x2))
}

func TestAdd(t *testing.T) {
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	assert.Equal(t, vm.stack.Top(), uint8(0x3))
}

func TestMinus(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x01, 0x0a, 0x0c}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	assert.Equal(t, vm.stack.Top(), uint8(0x1))
}
