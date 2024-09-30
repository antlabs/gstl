package rhashmap

import (
	"sort"
	"testing"
)

// 1. set get测试
// key string value bool
func Test_SetGet_StringBool(t *testing.T) {
	hm := New[string, bool]()
	hm.Set("hello", true)
	hm.Set("world", true)
	hm.Set("ni", true)
	hm.Set("hao", true)

	if !hm.Get("hello") {
		t.Errorf("Expected true, got false for key 'hello'")
	}
	if !hm.Get("world") {
		t.Errorf("Expected true, got false for key 'world'")
	}
	if !hm.Get("ni") {
		t.Errorf("Expected true, got false for key 'ni'")
	}
	if !hm.Get("hao") {
		t.Errorf("Expected true, got false for key 'hao'")
	}
}

// 1. set get测试
// key string, value string
func Test_SetGet_StringString(t *testing.T) {
	hm := New[string, string]()
	hm.Set("hello", "hello")
	hm.Set("world", "world")
	hm.Set("ni", "ni")
	hm.Set("hao", "hao")

	if hm.Get("hello") != "hello" {
		t.Errorf("Expected 'hello', got %v for key 'hello'", hm.Get("hello"))
	}
	if hm.Get("world") != "world" {
		t.Errorf("Expected 'world', got %v for key 'world'", hm.Get("world"))
	}
	if hm.Get("ni") != "ni" {
		t.Errorf("Expected 'ni', got %v for key 'ni'", hm.Get("ni"))
	}
	if hm.Get("hao") != "hao" {
		t.Errorf("Expected 'hao', got %v for key 'hao'", hm.Get("hao"))
	}
}

// 1. set get测试
// key string value string
func Test_SetGet_IntString(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(2, "world")
	hm.Set(3, "ni")
	hm.Set(4, "hao")

	if hm.Get(1) != "hello" {
		t.Errorf("Expected 'hello', got %v for key 1", hm.Get(1))
	}
	if hm.Get(2) != "world" {
		t.Errorf("Expected 'world', got %v for key 2", hm.Get(2))
	}
	if hm.Get(3) != "ni" {
		t.Errorf("Expected 'ni', got %v for key 3", hm.Get(3))
	}
	if hm.Get(4) != "hao" {
		t.Errorf("Expected 'hao', got %v for key 4", hm.Get(4))
	}
}

// 1. set get测试
func Test_SetGet_IntString_Lazyinit(t *testing.T) {
	var hm HashMap[int, string]
	hm.Set(1, "hello")
	hm.Set(2, "world")
	hm.Set(3, "ni")
	hm.Set(4, "hao")

	if hm.Get(1) != "hello" {
		t.Errorf("Expected 'hello', got %v for key 1", hm.Get(1))
	}
	if hm.Get(2) != "world" {
		t.Errorf("Expected 'world', got %v for key 2", hm.Get(2))
	}
	if hm.Get(3) != "ni" {
		t.Errorf("Expected 'ni', got %v for key 3", hm.Get(3))
	}
	if hm.Get(4) != "hao" {
		t.Errorf("Expected 'hao', got %v for key 4", hm.Get(4))
	}
}

// 1. set get测试
// 设计重复key
func Test_SetGet_Replace_IntString(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(1, "world")

	if hm.Get(1) != "world" {
		t.Errorf("Expected 'world', got %v for key 1", hm.Get(1))
	}
}

// 1. set get测试
// 获取空值数据
func Test_SetGet_Zero(t *testing.T) {
	hm := New[int, int]()
	for i := 0; i < 10; i++ {
		if hm.Get(i) != 0 {
			t.Errorf("Expected 0, got %v for key %d", hm.Get(i), i)
		}
	}

	for i := 0; i < 10; i++ {
		v, err := hm.GetWithBool(i)
		if err {
			t.Errorf("Expected false, got true for key %d", i)
		}
		if v != 0 {
			t.Errorf("Expected 0, got %v for key %d", v, i)
		}
	}
}

// 1. set get测试
// 测试重复key
func Test_SetGet_NotFound(t *testing.T) {
	hm := New[int, string]()
	hm.Set(1, "hello")
	hm.Set(1, "world")

	_, err := hm.GetWithBool(3)

	if err {
		t.Errorf("Expected false, got true for key 3")
	}
	if hm.Get(1) != "world" {
		t.Errorf("Expected 'world', got %v for key 1", hm.Get(1))
	}
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

	_, err := hm.GetWithBool(7)

	if err {
		t.Errorf("Expected false, got true for key 7")
	}
	if hm.Get(1) != "hello" {
		t.Errorf("Expected 'hello', got %v for key 1", hm.Get(1))
	}
}

// 测试Len接口
func Test_Len(t *testing.T) {
	hm := New[int, int]()
	max := 3333
	for i := 0; i < max; i++ {
		hm.Set(i, i)
	}
	if hm.Len() != max {
		t.Errorf("Expected %d, got %v", max, hm.Len())
	}
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
	if hm.Len() != 0 {
		t.Errorf("Expected 0, got %v", hm.Len())
	}
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

	if hm.Len() != 0 {
		t.Errorf("Expected 0, got %v", hm.Len())
	}
}

// 2. 测试删除功能
func Test_Delete_Empty(t *testing.T) {
	hm := New[int, int]()

	err := hm.Remove(0)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if hm.Len() != 0 {
		t.Errorf("Expected 0, got %v", hm.Len())
	}
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

	if hm.Len() != max {
		t.Errorf("Expected %d, got %v", max, hm.Len())
	}
	got := make([]int, 0, max)
	hm.Range(func(key int, val int) bool {
		got = append(got, key, val)
		return true
	})

	sort.Ints(got)
	if !slicesEqual(need, got) {
		t.Errorf("Expected %v, got %v", need, got)
	}
}

// 3. 测试Range
func Test_Range_Zero(t *testing.T) {
	max := 0
	hm := New[int, int]()
	need := []int{}

	if hm.Len() != max {
		t.Errorf("Expected %d, got %v", max, hm.Len())
	}
	got := make([]int, 0, max)
	hm.Range(func(key int, val int) bool {
		got = append(got, key, val)
		return true
	})

	sort.Ints(got)
	if !slicesEqual(need, got) {
		t.Errorf("Expected %v, got %v", need, got)
	}
}

func Test_Range_Rehasing(t *testing.T) {
	max := 5
	hm := New[int, int]()
	need := []int{}
	for i := 0; i < max; i++ {
		hm.Set(i, i)
		need = append(need, i, i)
	}

	if hm.Len() != max {
		t.Errorf("Expected %d, got %v", max, hm.Len())
	}
	got := make([]int, 0, max)
	hm.Range(func(key int, val int) bool {
		got = append(got, key, val)
		return true
	})

	sort.Ints(got)
	if !slicesEqual(need, got) {
		t.Errorf("Expected %v, got %v", need, got)
	}
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
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hm.Len() != 0 {
		t.Errorf("Expected 0, got %v", hm.Len())
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
