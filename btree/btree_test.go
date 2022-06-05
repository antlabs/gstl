package btree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试get set
// 不分裂逻辑
func Test_SetAndGet(t *testing.T) {
	b := New[int, int](0)

	max := 10
	for i := 0; i < max; i++ {
		b.Set(i, i)
	}

	for i := 0; i < max; i++ {
		v, err := b.Get(i)
		assert.NoError(t, err)
		assert.Equal(t, v, i)
	}
}

// 测试get set
func Test_SetAndGet_Split(t *testing.T) {
	b := New[int, int](2)

	max := 10
	for i := 1; i < max; i++ {
		b.Set(i, i)
	}

	for i := 1; i < max; i++ {
		v, err := b.Get(i)
		assert.NoError(t, err, fmt.Sprintf("index:%d", i))
		assert.Equal(t, v, i, fmt.Sprintf("index:%d", i))
	}
}
