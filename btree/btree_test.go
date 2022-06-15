package btree

import (
	"fmt"
	"testing"

	"github.com/guonaihong/gstl/cmp"
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
		v, err := b.GetWithErr(i)
		assert.NoError(t, err)
		assert.Equal(t, v, i)
	}
}

// 测试get set
// 分裂逻辑
func Test_Btree_SetAndGet_Split(t *testing.T) {
	b := New[int, int](2)

	max := 10
	for i := 0; i < max; i++ {
		b.Set(i, i)
	}

	for i := 0; i < max; i++ {
		v, err := b.GetWithErr(i)
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
		v, err := b.GetWithErr(i)
		assert.NoError(t, err, fmt.Sprintf("index:%d", i))
		assert.Equal(t, v, i, fmt.Sprintf("index:%d", i))
	}
}

// 测试get set, 小数据量下面的替换
func Test_Btree_SetAndGet_Replace(t *testing.T) {
	max := 10
	b := New[int, int](max)

	for i := 0; i < max; i++ {
		b.Set(i, i)
	}

	for i := 0; i < max; i++ {
		prev, replace := b.SetWithPrev(i, i+1)
		assert.True(t, replace)
		assert.Equal(t, prev, i)
	}

	for i := 0; i < max; i++ {
		v, err := b.GetWithErr(i)
		assert.NoError(t, err, fmt.Sprintf("index:%d", i))
		assert.Equal(t, v, i+1, fmt.Sprintf("index:%d", i))
	}
}

// 测试Range, 小数据量测试
func Test_Btree_Range(t *testing.T) {
	b := New[int, int](2)
	max := 100
	key := make([]int, 0, max)
	val := make([]int, 0, max)
	need := make([]int, 0, max)
	for i := max - 1; i >= 0; i-- {
		assert.NotPanics(t, func() {
			b.Set(i, i)
		}, fmt.Sprintf("index:%d", i))
	}

	for i := 0; i < max; i++ {
		need = append(need, i)
	}

	b.Range(func(k, v int) bool {
		key = append(key, k)
		val = append(val, k)
		return true
	})

	assert.Equal(t, key, need)
	assert.Equal(t, val, need)
}

// 测试TopMin, 它返回最小的几个值
func Test_Btree_TopMin(t *testing.T) {

	need := []int{}
	count10 := 10
	count100 := 100
	count1000 := 1000

	for i := 0; i < count1000; i++ {
		need = append(need, i)
	}

	needCount := []int{count10, count100, count100}
	for i, b := range []*Btree[int, int]{
		// btree里面元素 少于 TopMin 需要返回的值
		func() *Btree[int, int] {
			b := New[int, int](2)
			for i := 0; i < count10; i++ {
				b.Set(i, i)
			}

			assert.Equal(t, b.Len(), count10)
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *Btree[int, int] {

			b := New[int, int](2)
			for i := 0; i < count100; i++ {
				b.Set(i, i)
			}
			assert.Equal(t, b.Len(), count100)
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *Btree[int, int] {

			b := New[int, int](2)
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

// 测试倒序输出
func Test_Btree_RangePrev(t *testing.T) {

	b := New[int, int](2)
	max := 1000
	key := make([]int, 0, max)
	val := make([]int, 0, max)
	need := make([]int, 0, max)
	for i := 0; i < max; i++ {
		assert.NotPanics(t, func() {
			b.Set(i, i)
		}, fmt.Sprintf("index:%d", i))
	}

	for i := max - 1; i >= 0; i-- {
		need = append(need, i)
	}

	b.RangePrev(func(k, v int) bool {
		key = append(key, k)
		val = append(val, k)
		return true
	})

	assert.Equal(t, key, need)
}

func Test_Btree_RangePrev2(t *testing.T) {

	b := New[int, int](2)
	max := 1000
	key := make([]int, 0, max)
	val := make([]int, 0, max)
	need := make([]int, 0, max)
	for i := max - 1; i >= 0; i-- {
		assert.NotPanics(t, func() {
			b.Set(i, i)
		}, fmt.Sprintf("index:%d", i))
	}

	for i := max - 1; i >= 0; i-- {
		need = append(need, i)
	}

	b.RangePrev(func(k, v int) bool {
		key = append(key, k)
		val = append(val, k)
		return true
	})

	assert.Equal(t, key, need)
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

// 测试TopMax, 返回最大的几个数据降序返回
func Test_Btree_TopMax(t *testing.T) {

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

	for i, b := range []*Btree[int, int]{
		// btree里面元素 少于 TopMin 需要返回的值
		func() *Btree[int, int] {
			b := New[int, int](2)
			for i := 0; i < count10; i++ {
				b.Set(i, i)
			}

			assert.Equal(t, b.Len(), count10)
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *Btree[int, int] {

			b := New[int, int](2)
			for i := 0; i < count100; i++ {
				b.Set(i, i)
			}
			assert.Equal(t, b.Len(), count100)
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *Btree[int, int] {

			b := New[int, int](2)
			for i := 0; i < count1000; i++ {
				b.Set(i, i)
			}
			assert.Equal(t, b.Len(), count1000)
			return b
		}(),
	} {
		var key, val []int
		b.TopMax(count100, func(k, v int) bool {
			key = append(key, k)
			val = append(val, v)
			return true
		})
		length := cmp.Min(count[i], len(need[i]))
		assert.Equal(t, key, need[i][:length])
		assert.Equal(t, val, need[i][:length])
	}

}

// 测试btree删除的情况, 少量数量
func Test_Btree_Delete1(t *testing.T) {
	for max := 3; max < 1000; max++ {

		b := New[int, int](64)

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

// 测试draw
func Test_Btree_Draw(t *testing.T) {

	b := New[int, int](2)
	for i := 0; i < 10; i++ {
		b.Set(i, i)
	}

	b.Draw()
}

func Test_Btree_Delete2(t *testing.T) {
	b := New[int, int](2)

	for max := 0; max <= 33; max++ {
		//for max := 0; max <= 22; max++ {
		for i := 0; i < max; i++ {
			b.Set(i, i)
		}

		start := max / 2
		// 删除后半段
		for i := start; i < max; i++ {
			prev, ok := b.DeleteWithPrev(i)
			assert.True(t, ok, fmt.Sprintf("max:%d, i:%d", max, i))
			assert.Equal(t, prev, i, fmt.Sprintf("max:%d, i:%d", max, i))

			if !ok {
				return
			}
		}

		// 查找后半段, 应该找不到
		for i := start; i < max; i++ {
			v, err := b.GetWithErr(i)
			assert.Error(t, err, fmt.Sprintf("index:%d", i))
			assert.Equal(t, v, 0, fmt.Sprintf("index:%d", i))
		}

		// 查找前半段
		for i := 0; i < start; i++ {
			v, err := b.GetWithErr(i)
			assert.NoError(t, err, fmt.Sprintf("index:%d, max:%d, delete-start:%d", i, max, start))
			if err != nil {
				fmt.Println(b.GetWithErr(i))
				return
			}

			assert.Equal(t, v, i, fmt.Sprintf("index:%d", i))
		}

	}
}
