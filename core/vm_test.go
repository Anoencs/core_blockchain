package core

import (
	"projectx/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x9, 0x0a, 0x2, 0x0a}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Equal(t, vm.stack, types.NewStack(128))
}
func TestVM2(t *testing.T) {
	data := []byte{0x4f, 0x0b, 0x4f, 0x0b, 0x46, 0x0b, 0x03, 0x0a, 0x0e, 0x02, 0x0a, 0x03, 0x0a, 0x0d, 0x0f}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	// result := vm.stack.Pop().([]byte)
	//fmt.Printf("%+v", vm.stack.Top())
	//fmt.Printf("%+v\n", vm.stack.Get()...)
	//fmt.Println(string(result))
	//assert.Equal(t, "FOO", string(result))
}
func TestAdd(t *testing.T) {
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0c}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	assert.Equal(t, 0x3, vm.stack.Top())
}

func TestMinus(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0d}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	assert.Equal(t, 0x1, vm.stack.Top())
}
