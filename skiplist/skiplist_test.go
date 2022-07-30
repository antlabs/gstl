package skiplist

// apache 2.0 guonaihong
import (
	"fmt"
	"testing"

	"github.com/guonaihong/gstl/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	n := New[int, int]()
	assert.NotNil(t, n)
}

func Test_SetGet(t *testing.T) {
	zset := New[float64, string]()
	max := 100.0
	for i := 0.0; i < max; i++ {
		zset.Set(i, fmt.Sprintf("%d", int(i)))
	}

	for i := 0.0; i < max; i++ {
		v := zset.Get(i)
		assert.Equal(t, v, fmt.Sprintf("%d", int(i)))
	}
}

// 测试插入重复
func Test_InsertRepeatingElement(t *testing.T) {
	sl := New[float64, string]()
	max := 100
	for i := 0; i < max; i++ {
		sl.Set(float64(i), fmt.Sprint(i))
	}

	for i := 0; i < max; i++ {
		sl.Set(float64(i), fmt.Sprint(i+1))
	}

	for i := 0; i < max; i++ {
		assert.Equal(t, sl.Get(float64(i)), fmt.Sprint(i+1))
	}
}

func Test_SetGetRemove(t *testing.T) {
	zset := New[float64, float64]()

	max := 100.0
	for i := 0.0; i < max; i++ {
		zset.Set(i, i)
	}

	for i := 0.0; i < max; i++ {
		zset.Remove(i)
		assert.Equal(t, float64(zset.Len()), max-1)
		for j := 0.0; j < max; j++ {
			if j == i {
				continue
			}
			v, err := zset.GetWithErr(j)
			assert.NoError(t, err, fmt.Sprintf("score:%f, i:%f, j:%f", j, i, j))
			if err != nil {
				return
			}
			assert.Equal(t, v, j)
		}
		zset.Set(i, i)
	}
}

// 测试TopMin, 它返回最小的几个值
func Test_Skiplist_TopMin(t *testing.T) {

	need := []int{}
	count10 := 10
	count100 := 100
	count1000 := 1000

	for i := 0; i < count1000; i++ {
		need = append(need, i)
	}

	needCount := []int{count10, count100, count100}
	for i, b := range []*SkipList[float64, int]{
		// btree里面元素 少于 TopMin 需要返回的值
		func() *SkipList[float64, int] {
			b := New[float64, int]()
			for i := 0; i < count10; i++ {
				b.Set(float64(i), i)
			}

			assert.Equal(t, b.Len(), count10)
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count100; i++ {
				b.Set(float64(i), i)
			}
			assert.Equal(t, b.Len(), count100)
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count1000; i++ {
				b.Set(float64(i), i)
			}
			assert.Equal(t, b.Len(), count1000)
			return b
		}(),
	} {
		var key, val []int
		b.TopMin(count100, func(k float64, v int) bool {
			key = append(key, int(k))
			val = append(val, v)
			return true
		})
		assert.Equal(t, key, need[:needCount[i]])
		assert.Equal(t, val, need[:needCount[i]])
	}
}

// 测试下负数
func Test_Skiplist_TopMin2(t *testing.T) {
	start := -10
	max := 100
	limit := 10
	sl := New[float64, int]()

	need := make([]int, 0, limit)
	for i, l := start, limit; i < max && l > 0; i++ {
		sl.Set(float64(i), i)
		need = append(need, i)
		l--
	}

	got := make([]int, 0, limit)
	sl.TopMin(10, func(k float64, v int) bool {
		got = append(got, int(k))
		return true
	})

	assert.Equal(t, need, got)
}

// debug, 指定层
func Test_SkipList_SetAndGet_Level(t *testing.T) {

	sl := New[float64, int]()

	keys := []int{5, 8, 10}
	level := []int{2, 3, 5}
	for i, key := range keys {
		sl.InsertInner(float64(key), key, level[i])
	}

	sl.Draw()
	for _, i := range keys {
		v, count, _ := sl.GetWithMeta(float64(i))
		fmt.Printf("get %v count = %v, nodes:%v, level:%v maxlevel:%v\n",
			float64(i),
			count.Total,
			count.Keys,
			count.Level,
			count.MaxLevel)
		assert.Equal(t, v, i)
	}
}

// debug, 用的入口函数
func Test_SkipList_SetAndGet2(t *testing.T) {

	sl := New[float64, int]()

	max := 1000
	start := -1
	for i := max; i >= start; i-- {
		sl.Set(float64(i), i)
	}

	sl.Draw()
	for i := start; i < max; i++ {
		v, count, _ := sl.GetWithMeta(float64(i))
		fmt.Printf("get %v count = %v, nodes:%v, level:%v maxlevel:%v\n",
			float64(i),
			count.Total,
			count.Keys,
			count.Level,
			count.MaxLevel)
		assert.Equal(t, v, i)
	}
}

// 测试TopMax, 返回最大的几个数据降序返回
func Test_Skiplist_TopMax(t *testing.T) {

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

	for i, b := range []*SkipList[float64, int]{
		// btree里面元素 少于 TopMax 需要返回的值
		func() *SkipList[float64, int] {
			b := New[float64, int]()
			for i := 0; i < count10; i++ {
				b.Set(float64(i), i)
			}

			assert.Equal(t, b.Len(), count10)
			return b
		}(),
		// btree里面元素 等于 TopMax 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count100; i++ {
				b.Set(float64(i), i)
			}
			assert.Equal(t, b.Len(), count100)
			return b
		}(),
		// btree里面元素 大于 TopMax 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count1000; i++ {
				b.Set(float64(i), i)
			}
			assert.Equal(t, b.Len(), count1000)
			return b
		}(),
	} {
		var key, val []int
		b.TopMax(count100, func(k float64, v int) bool {
			key = append(key, int(k))
			val = append(val, v)
			return true
		})
		length := cmp.Min(count[i], len(need[i]))
		assert.Equal(t, key, need[i][:length])
		assert.Equal(t, val, need[i][:length])
	}

}
