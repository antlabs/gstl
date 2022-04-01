package vec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Int(t *testing.T) {
	v := New(1, 2, 3, 4, 5, 6)
	v.Push(7)
	v.Push(8)
	n, _ := v.Pop()
	assert.Equal(t, n, 8)
}

func Test_String(t *testing.T) {
	v := New("1", "2", "3", "4", "5", "6")
	v.Push("7")
	v.Push("8")
	n, _ := v.Pop()
	assert.Equal(t, n, "8")
}
