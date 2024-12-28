package btree

import (
	"testing"

	"github.com/antlabs/gstl/cmp"
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
		v, ok := b.TryGet(i)
		if !ok {
			t.Errorf("Expected true, got false for key %d", i)
		}
		if v != i {
			t.Errorf("Expected %d, got %d for key %d", i, v, i)
		}
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
		v, ok := b.TryGet(i)
		if !ok {
			t.Errorf("Expected true, got false for key %d", i)
		}
		if v != i {
			t.Errorf("Expected %d, got %d for key %d", i, v, i)
		}
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
		v, ok := b.TryGet(i)
		if !ok {
			t.Errorf("Expected true, got false for key %d", i)
		}
		if v != i {
			t.Errorf("Expected %d, got %d for key %d", i, v, i)
		}
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
		if !replace {
			t.Errorf("Expected true, got false for key %d", i)
		}
		if prev != i {
			t.Errorf("Expected %d, got %d for key %d", i, prev, i)
		}
	}

	for i := 0; i < max; i++ {
		v, ok := b.TryGet(i)
		if !ok {
			t.Errorf("Expected true, got false for key %d", i)
		}
		if v != i+1 {
			t.Errorf("Expected %d, got %d for key %d", i+1, v, i)
		}
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
		b.Set(i, i)
	}

	for i := 0; i < max; i++ {
		need = append(need, i)
	}

	b.Range(func(k, v int) bool {
		key = append(key, k)
		val = append(val, k)
		return true
	})

	if !slicesEqual(key, need) {
		t.Errorf("Expected %v, got %v", need, key)
	}
	if !slicesEqual(val, need) {
		t.Errorf("Expected %v, got %v", need, val)
	}
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
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *Btree[int, int] {
			b := New[int, int](2)
			for i := 0; i < count100; i++ {
				b.Set(i, i)
			}
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *Btree[int, int] {
			b := New[int, int](2)
			for i := 0; i < count1000; i++ {
				b.Set(i, i)
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
		if !slicesEqual(key, need[:needCount[i]]) {
			t.Errorf("Expected %v, got %v", need[:needCount[i]], key)
		}
		if !slicesEqual(val, need[:needCount[i]]) {
			t.Errorf("Expected %v, got %v", need[:needCount[i]], val)
		}
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
		b.Set(i, i)
	}

	for i := max - 1; i >= 0; i-- {
		need = append(need, i)
	}

	b.RangePrev(func(k, v int) bool {
		key = append(key, k)
		val = append(val, k)
		return true
	})

	if !slicesEqual(key, need) {
		t.Errorf("Expected %v, got %v", need, key)
	}
}

func Test_Btree_RangePrev2(t *testing.T) {
	b := New[int, int](2)
	max := 1000
	key := make([]int, 0, max)
	val := make([]int, 0, max)
	need := make([]int, 0, max)
	for i := max - 1; i >= 0; i-- {
		b.Set(i, i)
	}

	for i := max - 1; i >= 0; i-- {
		need = append(need, i)
	}

	b.RangePrev(func(k, v int) bool {
		key = append(key, k)
		val = append(val, k)
		return true
	})

	if !slicesEqual(key, need) {
		t.Errorf("Expected %v, got %v", need, key)
	}
}

// 测试Find接口
func Test_Btree_Find(t *testing.T) {
	b := New[int, int](2)
	b.Set(0, 0)
	b.Set(1, 1)
	b.Set(2, 2)

	index, _ := b.find(b.root, 2)
	if index != 2 {
		t.Errorf("Expected 2, got %d", index)
	}

	index, _ = b.find(b.root, 4)
	if index != 3 {
		t.Errorf("Expected 3, got %d", index)
	}
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
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *Btree[int, int] {
			b := New[int, int](2)
			for i := 0; i < count100; i++ {
				b.Set(i, i)
			}
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *Btree[int, int] {
			b := New[int, int](2)
			for i := 0; i < count1000; i++ {
				b.Set(i, i)
			}
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
		if !slicesEqual(key, need[i][:length]) {
			t.Errorf("Expected %v, got %v", need[i][:length], key)
		}
		if !slicesEqual(val, need[i][:length]) {
			t.Errorf("Expected %v, got %v", need[i][:length], val)
		}
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
			v, ok := b.TryGet(i)
			if !ok {
				t.Errorf("Expected true, got false for key %d", i)
			}
			if v != i {
				t.Errorf("Expected %d, got %d for key %d", i, v, i)
			}
		}

		// 0-max/2应该找不到
		for i := 0; i < max/2; i++ {
			v, ok := b.TryGet(i)
			if ok {
				t.Errorf("Expected false, got true for key %d", i)
			}
			if v != 0 {
				t.Errorf("Expected 0, got %d for key %d", v, i)
			}
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

	for max := 0; max <= 500; max++ {
		for i := 0; i < max; i++ {
			b.Set(i, i)
		}

		start := max / 2
		// 删除后半段
		for i := start; i < max; i++ {
			prev, ok := b.DeleteWithPrev(i)
			if !ok {
				t.Errorf("Expected true, got false for key %d", i)
			}
			if prev != i {
				t.Errorf("Expected %d, got %d for key %d", i, prev, i)
			}
		}

		// 查找后半段, 应该找不到
		for i := start; i < max; i++ {
			v, ok := b.TryGet(i)
			if ok {
				t.Errorf("Expected false, got true for key %d", i)
			}
			if v != 0 {
				t.Errorf("Expected 0, got %d for key %d", v, i)
			}
		}

		// 查找前半段
		for i := 0; i < start; i++ {
			v, ok := b.TryGet(i)
			if !ok {
				t.Errorf("Expected true, got false for key %d", i)
			}
			if v != i {
				t.Errorf("Expected %d, got %d for key %d", i, v, i)
			}
		}
	}
}

// Helper function to compare slices
func slicesEqual[T comparable](a, b []T) bool {
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
