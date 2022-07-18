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

	for i := 0; i < max; i++ {
		v, err := b.GetWithErr(i)
		assert.NoError(t, err)
		assert.Equal(t, v, i)
	}
}

// 从大到小, 插入
func Test_SetAndGet2(t *testing.T) {
	b := New[int, int]()
	max := 1000
	for i := max; i >= 0; i-- {
		b.SetWithPrev(i, i)
	}

	for i := max; i >= 0; i-- {
		v, err := b.GetWithErr(i)
		assert.NoError(t, err)
		assert.Equal(t, v, i)
	}
}

// 测试avltree删除的情况, 少量数量
func Test_RBTree_Delete1(t *testing.T) {
	for max := 3; max < 1000; max++ {

		b := New[int, int]()

		// 设置0-max
		for i := 0; i < max; i++ {
			b.Set(i, i)
		}

		// 删除0-max/2
		for i := 0; i < max/2; i++ {
			b.Delete(i)
		}

		// max/2-max应该能找到
		for i := max / 2; i < max; i++ {
			v, err := b.GetWithErr(i)
			assert.NoError(t, err, fmt.Sprintf("index:%d", i))
			assert.Equal(t, v, i, fmt.Sprintf("index:%d", i))
		}

		// 0-max/2应该找不到
		for i := 0; i < max/2; i++ {
			v, err := b.GetWithErr(i)
			assert.Error(t, err, fmt.Sprintf("index:%d", i))
			assert.Equal(t, v, 0, fmt.Sprintf("index:%d", i))
		}
	}
}
