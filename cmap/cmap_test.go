package cmap

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Store And Load
func Test_StoreAndLoad(t *testing.T) {
	m := New[string, string]()
	m.Store("hello", "1")
	m.Store("world", "2")
	v1, ok1 := m.Load("hello")
	assert.Equal(t, v1, "1")
	assert.True(t, ok1)

	v1, ok1 = m.Load("world")
	assert.Equal(t, v1, "2")
	assert.True(t, ok1)
}

// Store And Load
func Test_StoreDeleteLoad(t *testing.T) {
	m := New[string, string]()
	m.Store("hello", "1")
	m.Store("world", "2")

	m.Delete("hello")
	m.Delete("world")

	v1, ok1 := m.Load("hello")
	assert.Equal(t, v1, "")
	assert.False(t, ok1)

	v1, ok1 = m.Load("world")
	assert.Equal(t, v1, "")
	assert.False(t, ok1)
}

func Test_LoadAndDelete(t *testing.T) {

	m := New[string, string]()
	v1, ok1 := m.LoadAndDelete("hello")

	assert.Equal(t, v1, "")
	assert.False(t, ok1)

	m.Store("hello", "world")
	v1, ok1 = m.Load("hello")

	assert.Equal(t, v1, "world")

	v1, ok1 = m.LoadAndDelete("hello")
	assert.Equal(t, v1, "world")
	assert.True(t, ok1)
}

func Test_loadOrStore(t *testing.T) {
	m := New[string, string]()
	var m2 sync.Map
	v1, ok1 := m.LoadOrStore("hello", "world")
	v2, ok2 := m2.LoadOrStore("hello", "world")

	assert.Equal(t, ok1, ok2)
	assert.Equal(t, v1, v2.(string))
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

	assert.Equal(t, count, 1)
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

	assert.Equal(t, keyAll, gotKey)
	assert.Equal(t, valAll, gotVal)
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

	assert.Equal(t, keyAll, gotKey)
	assert.Equal(t, valAll, gotVal)
}

func Test_Len(t *testing.T) {
	m := New[string, string]()
	m.Store("1", "1")
	m.Store("2", "2")
	m.Store("3", "3")
	assert.Equal(t, m.Len(), 3)
}

func Test_New(t *testing.T) {
	m := New[string, string]()
	m.Store("1", "1")
	m.Store("2", "2")
	m.Store("3", "3")
	assert.Equal(t, m.Len(), 3)
}

func Test_Keys(t *testing.T) {
	m := New[string, string]()
	m.Store("a", "1")
	m.Store("b", "2")
	m.Store("c", "3")
	get := m.Keys()
	sort.Strings(get)
	assert.Equal(t, get, []string{"a", "b", "c"})

	m2 := New[string, string]()
	assert.Equal(t, len(m2.Values()), 0)
}

func Test_Values(t *testing.T) {
	m := New[string, string]()
	m.Store("a", "1")
	m.Store("b", "2")
	m.Store("c", "3")
	get := m.Values()
	sort.Strings(get)
	assert.Equal(t, get, []string{"1", "2", "3"})

	m2 := New[string, string]()
	assert.Equal(t, len(m2.Keys()), 0)
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
