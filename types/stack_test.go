package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStack(t *testing.T) {
	stackint := NewStack(1, 2, 3)
	assert.Equal(t, stackint.data, []int{1, 2, 3})
	assert.Equal(t, stackint.size, 3)

	stackstring := NewStack("a", "b", "c")
	assert.Equal(t, stackstring.data, []string{"a", "b", "c"})
	assert.Equal(t, stackstring.size, 3)
}

func TestStackMethod(t *testing.T) {
	stack := NewStack(1, 2, 3)

	stack.Push(4)
	assert.Equal(t, stack.data, []int{1, 2, 3, 4})
	assert.Equal(t, stack.size, 4)
	assert.Equal(t, stack.Top(), 4)

	stack.Push(5)
	assert.Equal(t, stack.data, []int{1, 2, 3, 4, 5})
	assert.Equal(t, stack.size, 5)
	assert.Equal(t, stack.Top(), 5)

	stack.Clear()
	assert.Equal(t, stack.size, 0)
}
