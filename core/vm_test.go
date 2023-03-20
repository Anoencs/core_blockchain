package core

import (
	"fmt"
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
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0d, 0x4f, 0x0b, 0x4f, 0x0b, 0x46, 0x0b, 0x03, 0x0a, 0x0e, 0x0f}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	valueBytes, err := vm.contractState.Get([]byte("FOO"))
	value := deserializeInt64(valueBytes)
	assert.Nil(t, err)

	assert.Equal(t, value, int64(1))
}

func TestStoreAndGet(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0d, 0x4f, 0x0b, 0x4f, 0x0b, 0x46, 0x0b, 0x03, 0x0a, 0x0e, 0x0f}
	pushFoo := []byte{0x4f, 0x0b, 0x4f, 0x0b, 0x46, 0x0b, 0x03, 0x0a, 0x0e, 0x001}

	data = append(data, pushFoo...)

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	//	fmt.Printf("%+v", vm.stack.Get()...)
	value := vm.stack.Pop().([]byte)
	valueSerialize := deserializeInt64(value)
	assert.Equal(t, valueSerialize, int64(1))
}

func TestStoreAndGet2(t *testing.T) {
	data := []byte{0x09, 0x0a, 0x04e, 0x0b, 0x041, 0x0b, 0x02, 0x0a, 0x0e, 0x0f}
	pushAN := []byte{0x04e, 0x0b, 0x041, 0x0b, 0x02, 0x0a, 0x0e, 0x001}
	data = append(data, pushAN...)
	contracState := NewState()
	vm := NewVM(data, contracState)
	assert.Nil(t, vm.Run())

	value := vm.stack.Pop().([]byte)
	fmt.Printf("%+v", value)
	valueSerialize := deserializeInt64(value)
	assert.Equal(t, valueSerialize, int64(9))
}
