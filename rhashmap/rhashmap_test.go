package rhashmap

// apache 2.0 antlabs
import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1. set get测试
// key string value bool
func Test_SetGet_StringBool(t *testing.T) {
	hm := New[string, bool]()
	hm.Set("hello", true)
	hm.Set("world", true)
	hm.Set("ni", true)
	hm.Set("hao", true)

	assert.True(t, hm.Get("hello"))
	assert.True(t, hm.Get("world"))
	assert.True(t, hm.Get("ni"))
	assert.True(t, hm.Get("hao"))
}

// 1. set get测试
// key string, value string
func Test_SetGet_StringString(t *testing.T) {
	hm := New[string, string]()
	hm.Set("hello", "hello")
	hm.Set("world", "world")
	hm.Set("ni", "ni")
	hm.Set("hao", "hao")

	assert.Equal(t, hm.Get("hello"), "hello")
	assert.Equal(t, hm.Get("world"), "world")
	assert.Equal(t, hm.Get("ni"), "ni")
	assert.Equal(t, hm.Get("hao"), "hao")
}

// 1. set get测试
// key string value string
func Test_SetGet_IntString(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(2, "world")
	hm.Set(3, "ni")
	hm.Set(4, "hao")

	assert.Equal(t, hm.Get(1), "hello")
	assert.Equal(t, hm.Get(2), "world")
	assert.Equal(t, hm.Get(3), "ni")
	assert.Equal(t, hm.Get(4), "hao")
}

// 1. set get测试
func Test_SetGet_IntString_Lazyinit(t *testing.T) {
	var hm HashMap[int, string]
	hm.Set(1, "hello")
	hm.Set(2, "world")
	hm.Set(3, "ni")
	hm.Set(4, "hao")

	assert.Equal(t, hm.Get(1), "hello")
	assert.Equal(t, hm.Get(2), "world")
	assert.Equal(t, hm.Get(3), "ni")
	assert.Equal(t, hm.Get(4), "hao")
}

// 1. set get测试
// 设计重复key
func Test_SetGet_Replace_IntString(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(1, "world")

	assert.Equal(t, hm.Get(1), "world")
}

// 1. set get测试
// 获取空值数据
func Test_SetGet_Zero(t *testing.T) {
	hm := New[int, int]()
	for i := 0; i < 10; i++ {
		assert.Equal(t, hm.Get(i), 0)
	}

	for i := 0; i < 10; i++ {
		v, err := hm.GetWithErr(i)
		assert.Error(t, err)
		assert.Equal(t, v, 0)
	}
}

// 1. set get测试
// 测试重复key
func Test_SetGet_NotFound(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(1, "world")

	_, err := hm.GetWithErr(3)

	assert.Error(t, err)
	assert.Equal(t, hm.Get(1), "world")
}

// 1. set get测试
// 测试
func Test_SetGet_Rehashing(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(2, "world")
	hm.Set(3, "hello")
	hm.Set(4, "world")
	hm.Set(5, "world")

	_, err := hm.GetWithErr(7)

	assert.Error(t, err)
	assert.Equal(t, hm.Get(1), "hello")

}

// 测试Len接口
func Test_Len(t *testing.T) {
	hm := New[int, int]()
	max := 3333
	for i := 0; i < max; i++ {
		hm.Set(i, i)
	}
	assert.Equal(t, hm.Len(), max)
}

// 2.测试删除功能
func Test_Delete(t *testing.T) {
	hm := New[int, int]()

	max := 3333
	for i := 0; i < max; i++ {
		hm.Set(i, i)
	}

	for i := 0; i < max; i++ {
		hm.Delete(i)
	}
	assert.Equal(t, hm.Len(), 0)
}

// 2. 测试删除功能
func Test_Delete_NotFound(t *testing.T) {
	hm := New[int, int]()

	max := 4 //不要修改4
	for i := 0; i < max; i++ {
		hm.Set(i, i)
	}

	hm.Delete(max + 1)
	for i := 0; i < max; i++ {
		hm.Delete(i)
	}

	assert.Equal(t, hm.Len(), 0)
}

// 2. 测试删除功能
func Test_Delete_Empty(t *testing.T) {
	hm := New[int, int]()

	err := hm.Remove(0)
	assert.Error(t, err)
	assert.Equal(t, hm.Len(), 0)
}

// 3. 测试Range
func Test_Range(t *testing.T) {
	max := 100
	hm := NewWithOpt[int, int](WithCap(max))
	need := []int{}
	for i := 0; i < max; i++ {
		need = append(need, i, i)
		hm.Set(i, i)
	}

	assert.Equal(t, hm.Len(), max)
	got := make([]int, 0, max)
	hm.Range(func(key int, val int) bool {
		got = append(got, key, val)
		return true
	})

	sort.Ints(got)
	assert.Equal(t, need, got)
}

// 3. 测试Range
func Test_Range_Zero(t *testing.T) {
	max := 0
	hm := New[int, int]()
	need := []int{}

	assert.Equal(t, hm.Len(), max)
	got := make([]int, 0, max)
	hm.Range(func(key int, val int) bool {
		got = append(got, key, val)
		return true
	})

	sort.Ints(got)
	assert.Equal(t, need, got)
}
func Test_Range_Rehasing(t *testing.T) {
	max := 5
	hm := New[int, int]()
	need := []int{}
	for i := 0; i < max; i++ {
		hm.Set(i, i)
		need = append(need, i, i)
	}

	assert.Equal(t, hm.Len(), max)
	got := make([]int, 0, max)
	hm.Range(func(key int, val int) bool {
		got = append(got, key, val)
		return true
	})

	sort.Ints(got)
	assert.Equal(t, need, got)
}

// 测试shrink
func Test_Range_ShrinkToFit(t *testing.T) {
	hm := New[int, int]()

	max := 3333
	for i := 0; i < max; i++ {
		hm.Set(i, i)
	}

	for i := 0; i < max; i++ {
		hm.Delete(i)
	}

	err := hm.ShrinkToFit()
	assert.NoError(t, err)
	assert.Equal(t, hm.Len(), 0)
}
