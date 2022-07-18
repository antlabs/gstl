package rbtree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 从小到大, 插入
func Test_SetAndGet(t *testing.T) {
	b := New[int, int]()
	max := 1000
	for i := 0; i < max; i++ {
		b.SetWithPrev(i, i)
	}

	b.Range(func(k, v int) bool {
		fmt.Println("k", k, "v", v)
		return true
	})

	for i := 0; i < max; i++ {
		v, err := b.GetWithErr(i)
		assert.NoError(t, err)
		assert.Equal(t, v, i)
	}
}
