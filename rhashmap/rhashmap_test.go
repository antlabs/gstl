package rhashmap

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

	assert.True(t, hm.GetOrZero("hello"))
	assert.True(t, hm.GetOrZero("world"))
	assert.True(t, hm.GetOrZero("ni"))
	assert.True(t, hm.GetOrZero("hao"))
}

// 1. set get测试
// key string, value string
func Test_SetGet_StringString(t *testing.T) {
	hm := New[string, string]()
	hm.Set("hello", "hello")
	hm.Set("world", "world")
	hm.Set("ni", "ni")
	hm.Set("hao", "hao")

	assert.Equal(t, hm.GetOrZero("hello"), "hello")
	assert.Equal(t, hm.GetOrZero("world"), "world")
	assert.Equal(t, hm.GetOrZero("ni"), "ni")
	assert.Equal(t, hm.GetOrZero("hao"), "hao")
}

// 1. set get测试
// key string value string
func Test_SetGet_IntString(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(2, "world")
	hm.Set(3, "ni")
	hm.Set(4, "hao")

	assert.Equal(t, hm.GetOrZero(1), "hello")
	assert.Equal(t, hm.GetOrZero(2), "world")
	assert.Equal(t, hm.GetOrZero(3), "ni")
	assert.Equal(t, hm.GetOrZero(4), "hao")
}

// 1. set get测试
// 设计重复key
func Test_SetGet_Replace_IntString(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(1, "world")

	assert.Equal(t, hm.GetOrZero(1), "world")
}

// 1. set get测试
// 获取空值数据
func Test_SetGet_Zero(t *testing.T) {
	hm := New[int, int]()
	for i := 0; i < 10; i++ {
		assert.Equal(t, hm.GetOrZero(i), 0)
	}

	for i := 0; i < 10; i++ {
		v, err := hm.Get(i)
		assert.Error(t, err)
		assert.Equal(t, v, 0)
	}
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
func Test_Delete_Empty(t *testing.T) {
	hm := New[int, int]()

	err := hm.Remove(0)
	assert.Error(t, err)
	assert.Equal(t, hm.Len(), 0)
}

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
	hm.Range(func(key int, val int) {
		got = append(got, key, val)
	})

	sort.Ints(got)
	assert.Equal(t, need, got)
}
