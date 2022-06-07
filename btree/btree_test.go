package btree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试get set
// 不分裂逻辑
func Test_Btree_SetAndGet(t *testing.T) {
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
// 分裂逻辑
func Test_Btree_SetAndGet_Split(t *testing.T) {
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

// 测试get set, 大数据量下面的测试
// 分裂逻辑
func Test_Btree_SetAndGet_Split_Big(t *testing.T) {
	max := 10000
	b := New[int, int](max)

	for i := 0; i < max; i++ {
		b.Set(i, i)
	}

	for i := 0; i < max; i++ {
		v, err := b.Get(i)
		assert.NoError(t, err, fmt.Sprintf("index:%d", i))
		assert.Equal(t, v, i, fmt.Sprintf("index:%d", i))
	}
}

func Test_Btree_Range(t *testing.T) {
	b := New[int, int](2)
	max := 10
	for i := 0; i < max; i++ {
		assert.NotPanics(t, func() {
			b.Set(i, i)
		}, fmt.Sprintf("index:%d", i))
	}

	/*
		b.Range(func(k, v int) bool {
			fmt.Println(k, v)
			return true
		})
	*/
}

// 测试Find接口
func Test_Btree_Find(t *testing.T) {
	b := New[int, int](2)
	b.Set(0, 0)
	b.Set(1, 1)
	b.Set(2, 2)

	index, _ := b.find(b.root, 2)
	assert.Equal(t, index, 2)

	index, _ = b.find(b.root, 4)
	assert.Equal(t, index, 3)
}
