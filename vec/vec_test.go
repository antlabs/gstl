package vec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_Push_Pop_Int(t *testing.T) {
	v := New(1, 2, 3, 4, 5, 6)
	v.Push(7)
	v.Push(8)
	n, _ := v.Pop()
	assert.Equal(t, n, 8)
}

func Test_New_Push_Pop_String(t *testing.T) {
	v := New("1", "2", "3", "4", "5", "6")
	v.Push("7")
	v.Push("8")
	n, _ := v.Pop()
	assert.Equal(t, n, "8")
}

func Test_RotateLeft(t *testing.T) {
	v := New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(2)
	assert.Equal(t, v.ToSlice(), []byte{'c', 'd', 'e', 'f', 'a', 'b'})

	v = New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(8)
	assert.Equal(t, v.ToSlice(), []byte{'c', 'd', 'e', 'f', 'a', 'b'})

	v = New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(0)
	assert.Equal(t, v.ToSlice(), []byte{'a', 'b', 'c', 'd', 'e', 'f'})
}
