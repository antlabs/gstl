package cmap

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

// Store And Load
func Test_StoreAndLoad(t *testing.T) {
	m := New[string, string]()
	m.Store("hello", "1")
	m.Store("world", "2")
	v1, ok1 := m.Load("hello")
	if v1 != "1" {
		t.Errorf("expected '1', got '%s'", v1)
	}
	if !ok1 {
		t.Errorf("expected true, got false")
	}

	v1, ok1 = m.Load("world")
	if v1 != "2" {
		t.Errorf("expected '2', got '%s'", v1)
	}
	if !ok1 {
		t.Errorf("expected true, got false")
	}
}

// Store And Load
func Test_StoreDeleteLoad(t *testing.T) {
	m := New[string, string]()
	m.Store("hello", "1")
	m.Store("world", "2")

	m.Delete("hello")
	m.Delete("world")

	v1, ok1 := m.Load("hello")
	if v1 != "" {
		t.Errorf("expected '', got '%s'", v1)
	}
	if ok1 {
		t.Errorf("expected false, got true")
	}

	v1, ok1 = m.Load("world")
	if v1 != "" {
		t.Errorf("expected '', got '%s'", v1)
	}
	if ok1 {
		t.Errorf("expected false, got true")
	}
}

func Test_LoadAndDelete(t *testing.T) {
	m := New[string, string]()
	v1, ok1 := m.LoadAndDelete("hello")

	if v1 != "" {
		t.Errorf("expected '', got '%s'", v1)
	}
	if ok1 {
		t.Errorf("expected false, got true")
	}

	m.Store("hello", "world")
	v1, ok1 = m.Load("hello")

	if v1 != "world" {
		t.Errorf("expected 'world', got '%s'", v1)
	}

	v1, ok1 = m.LoadAndDelete("hello")
	if v1 != "world" {
		t.Errorf("expected 'world', got '%s'", v1)
	}
	if !ok1 {
		t.Errorf("expected true, got false")
	}
}

func Test_loadOrStore(t *testing.T) {
	m := New[string, string]()
	var m2 sync.Map
	v1, ok1 := m.LoadOrStore("hello", "world")
	v2, ok2 := m2.LoadOrStore("hello", "world")

	if ok1 != ok2 {
		t.Errorf("expected %v, got %v", ok2, ok1)
	}
	if v1 != v2.(string) {
		t.Errorf("expected '%s', got '%s'", v2.(string), v1)
	}
}

func Test_RangeBreak(t *testing.T) {
	m := New[string, string]()
	m.Store("1", "1")
	m.Store("2", "2")

	count := 0
	m.Range(func(key, val string) bool {
		count++
		return false
	})

	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func Test_Range(t *testing.T) {
	m := New[string, string]()
	max := 5
	keyAll := []string{}
	valAll := []string{}

	for i := 1; i < max; i++ {
		key := fmt.Sprintf("%dk", i)
		val := fmt.Sprintf("%dv", i)
		keyAll = append(keyAll, key)
		valAll = append(valAll, val)
		m.Store(key, val)
	}

	gotKey := []string{}
	gotVal := []string{}
	m.Range(func(key, val string) bool {
		gotKey = append(gotKey, key)
		gotVal = append(gotVal, val)
		return true
	})

	sort.Strings(gotKey)
	sort.Strings(gotVal)

	if !equalSlices(keyAll, gotKey) {
		t.Errorf("expected keys %v, got %v", keyAll, gotKey)
	}
	if !equalSlices(valAll, gotVal) {
		t.Errorf("expected values %v, got %v", valAll, gotVal)
	}
}

func Test_Iter(t *testing.T) {
	m := New[string, string]()
	max := 5
	keyAll := []string{}
	valAll := []string{}

	for i := 1; i < max; i++ {
		key := fmt.Sprintf("%dk", i)
		val := fmt.Sprintf("%dv", i)
		keyAll = append(keyAll, key)
		valAll = append(valAll, val)
		m.Store(key, val)
	}

	gotKey := []string{}
	gotVal := []string{}
	for pair := range m.Iter() {
		gotKey = append(gotKey, pair.Key)
		gotVal = append(gotVal, pair.Val)
	}

	sort.Strings(gotKey)
	sort.Strings(gotVal)

	if !equalSlices(keyAll, gotKey) {
		t.Errorf("expected keys %v, got %v", keyAll, gotKey)
	}
	if !equalSlices(valAll, gotVal) {
		t.Errorf("expected values %v, got %v", valAll, gotVal)
	}
}

func Test_Len(t *testing.T) {
	m := New[string, string]()
	m.Store("1", "1")
	m.Store("2", "2")
	m.Store("3", "3")
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
}

func Test_New(t *testing.T) {
	m := New[string, string]()
	m.Store("1", "1")
	m.Store("2", "2")
	m.Store("3", "3")
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
}

func Test_Keys(t *testing.T) {
	m := New[string, string]()
	m.Store("a", "1")
	m.Store("b", "2")
	m.Store("c", "3")
	get := m.Keys()
	sort.Strings(get)
	if !equalSlices(get, []string{"a", "b", "c"}) {
		t.Errorf("expected keys %v, got %v", []string{"a", "b", "c"}, get)
	}

	m2 := New[string, string]()
	if len(m2.Values()) != 0 {
		t.Errorf("expected 0, got %d", len(m2.Values()))
	}
}

func Test_Values(t *testing.T) {
	m := New[string, string]()
	m.Store("a", "1")
	m.Store("b", "2")
	m.Store("c", "3")
	get := m.Values()
	sort.Strings(get)
	if !equalSlices(get, []string{"1", "2", "3"}) {
		t.Errorf("expected values %v, got %v", []string{"1", "2", "3"}, get)
	}

	m2 := New[string, string]()
	if len(m2.Keys()) != 0 {
		t.Errorf("expected 0, got %d", len(m2.Keys()))
	}
}

func Test_UpdateOrInsert(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		m := New[string, string]()
		m.Store("a", "1")
		m.Store("b", "2")
		m.Store("c", "3")
		m.UpdateOrInsert("a", func(exist bool, old string) string {
			if !exist {
				t.Error("should exist")
			}
			if exist {
				return "4"
			}
			return old
		})
		get, _ := m.Load("a")
		if get != "4" {
			t.Error("should be 4")
		}
	})

	t.Run("Insert", func(t *testing.T) {
		m := New[string, string]()
		m.Store("a", "1")
		m.UpdateOrInsert("b", func(exist bool, old string) string {
			if !exist {
				return "2"
			}
			return ""
		})

		get, _ := m.Load("b")
		if get != "2" {
			t.Error("should be 2")
		}
	})
}

// 辅助函数，用于比较两个切片是否相等
func equalSlices(a, b []string) bool {
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
