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

func Test_RPop(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(3), []string{"2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).RPop(3), []int{2, 3, 4})
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(10), []string{"1", "2", "3", "4"})
}

func Test_LPop(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(3), []string{"1", "2", "3"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).LPop(3), []int{1, 2, 3})
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(10), []string{"1", "2", "3", "4"})
}

func Test_RangeSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangeSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	assert.Equal(t, all, []string{"1", "2", "3", "4"})
}

func Test_RangePrevSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangePrevSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	assert.Equal(t, all, []string{"4", "3", "2", "1"})
}
