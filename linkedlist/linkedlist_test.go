package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Push(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").ToSlice(), []string{"1", "2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).ToSlice(), []int{1, 2, 3, 4})
}

func Test_RPush(t *testing.T) {
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4").ToSlice(), []string{"1", "2", "3", "4"})
	assert.Equal(t, New[int]().RPush(1, 2, 3, 4).ToSlice(), []int{1, 2, 3, 4})
}

func Test_Len(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").Len(), 4)
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).Len(), 4)
}

func Test_LPop(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(3), []string{"2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).RPop(3), []int{2, 3, 4})
}
