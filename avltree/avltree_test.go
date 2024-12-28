package avltree

import (
	"testing"

	"github.com/antlabs/gstl/cmp"
	"github.com/antlabs/gstl/vec"
)

// 从小到大, 插入
func Test_SetAndGet(t *testing.T) {
	b := New[int, int]()
	max := 1000
	for i := 0; i < max; i++ {
		b.Swap(i, i)
	}

	for i := 0; i < max; i++ {
		v, ok := b.TryGet(i)
		if !ok {
			t.Errorf("expected true, got false for index %d", i)
		}
		if v != i {
			t.Errorf("expected %d, got %d for index %d", i, v, i)
		}
	}
}

// 从大到小, 插入
func Test_SetAndGet2(t *testing.T) {
	b := New[int, int]()
	max := 1000
	for i := max; i >= 0; i-- {
		b.Swap(i, i)
	}

	for i := max; i >= 0; i-- {
		v, ok := b.TryGet(i)
		if !ok {
			t.Errorf("expected true, got false for index %d", i)
		}
		if v != i {
			t.Errorf("expected %d, got %d for index %d", i, v, i)
		}
	}
}

// 测试avltree删除的情况, 少量数量
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
			v, ok := b.TryGet(i)
			if !ok {
				t.Errorf("expected true, got false for index %d", i)
			}
			if v != i {
				t.Errorf("expected %d, got %d for index %d", i, v, i)
			}
		}

		// 0-max/2应该找不到
		for i := 0; i < max/2; i++ {
			v, ok := b.TryGet(i)
			if ok {
				t.Errorf("expected false, got true for index %d", i)
			}
			if v != 0 {
				t.Errorf("expected 0, got %d for index %d", v, i)
			}
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

			if b.Len() != count10 {
				t.Errorf("expected length %d, got %d", count10, b.Len())
			}
			return b
		}(),
		// btree里面元素 等于 TopMax 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count100; i++ {
				b.Set(int(i), i)
			}
			if b.Len() != count100 {
				t.Errorf("expected length %d, got %d", count100, b.Len())
			}
			return b
		}(),
		// btree里面元素 大于 TopMax 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count1000; i++ {
				b.Set(int(i), i)
			}
			if b.Len() != count1000 {
				t.Errorf("expected length %d, got %d", count1000, b.Len())
			}
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
		if !equalSlices(key, need[i][:length]) {
			t.Errorf("expected keys %v, got %v", need[i][:length], key)
		}
		if !equalSlices(val, need[i][:length]) {
			t.Errorf("expected values %v, got %v", need[i][:length], val)
		}
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

			if b.Len() != count10 {
				t.Errorf("expected length %d, got %d", count10, b.Len())
			}
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count100; i++ {
				b.Set(i, i)
			}
			if b.Len() != count100 {
				t.Errorf("expected length %d, got %d", count100, b.Len())
			}
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *AvlTree[int, int] {

			b := New[int, int]()
			for i := 0; i < count1000; i++ {
				b.Set(i, i)
			}
			if b.Len() != count1000 {
				t.Errorf("expected length %d, got %d", count1000, b.Len())
			}
			return b
		}(),
	} {
		var key, val []int
		b.TopMin(count100, func(k, v int) bool {
			key = append(key, k)
			val = append(val, v)
			return true
		})
		if !equalSlices(key, need[:needCount[i]]) {
			t.Errorf("expected keys %v, got %v", need[:needCount[i]], key)
		}
		if !equalSlices(val, need[:needCount[i]]) {
			t.Errorf("expected values %v, got %v", need[:needCount[i]], val)
		}
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

	if !equalSlices(gotKey, dataRev) {
		t.Errorf("expected keys %v, got %v", dataRev, gotKey)
	}
	if !equalSlices(gotVal, dataRev) {
		t.Errorf("expected values %v, got %v", dataRev, gotVal)
	}
}

func Test_AvlTree_InsertOrUpdate(t *testing.T) {
	b := New[int, int]()
	max := 100

	// Insert elements
	for i := 0; i < max; i++ {
		b.InsertOrUpdate(i, i, func(prev, new int) int {
			return prev + new
		})
	}

	// Update elements
	for i := 0; i < max; i++ {
		b.InsertOrUpdate(i, i, func(prev, new int) int {
			return prev + new
		})
	}

	// Verify elements
	for i := 0; i < max; i++ {
		v, ok := b.TryGet(i)
		if !ok || v != i*2 {
			t.Errorf("expected %d, got %v", i*2, v)
		}
	}
}

func Test_AvlTree_InsertOrUpdate2(t *testing.T) {
	b := New[int, int]()
	max := 100

	// Insert elements
	for i := 0; i < max; i++ {
		b.InsertOrUpdate(i, i, func(prev, new int) int {
			return prev + new
		})
	}

	// Update elements
	for i := 0; i < max; i++ {
		b.InsertOrUpdate(i, i*2, func(prev, new int) int {
			return prev + new
		})
	}

	// Verify elements
	for i := 0; i < max; i++ {
		v, ok := b.TryGet(i)
		if !ok || v != i*3 {
			t.Errorf("expected %d, got %v", i*3, v)
		}
	}
}

// 辅助函数，用于比较两个切片是否相等
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
