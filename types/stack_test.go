package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(128)
	s.Push(1)
	s.Push(2)
	value := s.Pop()
	assert.Equal(t, value, 2)
	fmt.Print(s)
	value = s.Pop()
	assert.Equal(t, value, 1)
	fmt.Print(s)

}
