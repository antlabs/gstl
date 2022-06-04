package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试get set
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
