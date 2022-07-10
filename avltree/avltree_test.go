package avltree

import (
	"fmt"
	"testing"

	"github.com/guonaihong/gstl/cmp"
	"github.com/guonaihong/gstl/vec"
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

// 测试btree删除的情况, 少量数量
func Test_AVLTree_Delete1(t *testing.T) {
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

// 测试TopMax, 返回最大的几个数据降序返回
func Test_AvlTree_TopMax(t *testing.T) {

	need := [3][]int{}
	count10 := 10
	count100 := 100
	count1000 := 1000
	count := []int{count10, count100, count1000}

	for i := 0; i < len(count); i++ {
		for j, k := count[i]-1, count100-1; j >= 0 && k >= 0; j-- {
			need[i] = append(need[i], j)
			k--
		}
	}

	for i, b := range []*AvlTree[int, int]{
		// btree里面元素 少于 TopMax 需要返回的值
		func() *AvlTree[int, int] {
			b := New[int, int]()
			for i := 0; i < count10; i++ {
				b.Set(i, i)
			}

			b.Draw()

			assert.Equal(t, b.Len(), count10)
			return b
		}(),
		// btree里面元素 等于 TopMax 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count100; i++ {
				b.Set(int(i), i)
			}
			assert.Equal(t, b.Len(), count100)
			return b
		}(),
		// btree里面元素 大于 TopMax 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count1000; i++ {
				b.Set(int(i), i)
			}
			assert.Equal(t, b.Len(), count1000)
			return b
		}(),
	} {
		var key, val []int
		b.TopMax(count100, func(k int, v int) bool {
			key = append(key, int(k))
			val = append(val, v)
			return true
		})
		length := cmp.Min(count[i], len(need[i]))
		assert.Equal(t, key, need[i][:length])
		assert.Equal(t, val, need[i][:length])
	}

}

// 测试TopMin, 它返回最小的几个值
func Test_AvlTree_TopMin(t *testing.T) {

	need := []int{}
	count10 := 10
	count100 := 100
	count1000 := 1000

	for i := 0; i < count1000; i++ {
		need = append(need, i)
	}

	needCount := []int{count10, count100, count100}
	for i, b := range []*AvlTree[int, int]{
		// btree里面元素 少于 TopMin 需要返回的值
		func() *AvlTree[int, int] {
			b := New[int, int]()
			for i := 0; i < count10; i++ {
				b.Set(i, i)
			}

			assert.Equal(t, b.Len(), count10)
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count100; i++ {
				b.Set(i, i)
			}
			assert.Equal(t, b.Len(), count100)
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count1000; i++ {
				b.Set(i, i)
			}
			assert.Equal(t, b.Len(), count1000)
			return b
		}(),
	} {
		var key, val []int
		b.TopMin(count100, func(k, v int) bool {
			key = append(key, k)
			val = append(val, v)
			return true
		})
		assert.Equal(t, key, need[:needCount[i]])
		assert.Equal(t, val, need[:needCount[i]])
	}

}

func Test_RanePrev(t *testing.T) {
	a := New[int, int]()
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	dataRev := vec.New(data...).Clone().Rev().ToSlice()
	for i := len(data) / 2; i >= 0; i-- {
		a.Set(i, i)
	}

	for i := len(data)/2 + 1; i < len(data); i++ {
		a.Set(i, i)
	}

	//a.Draw()

	var gotKey []int
	var gotVal []int
	a.RangePrev(func(k, v int) bool {
		gotKey = append(gotKey, k)
		gotVal = append(gotVal, v)

		return true
	})

	assert.Equal(t, gotKey, dataRev)
	assert.Equal(t, gotVal, dataRev)
}
