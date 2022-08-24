package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// set get 预期是设置进去， 也能读出来
func Test_TrimeMap_SetGet(t *testing.T) {

	tm := New[string]()
	max := 1000

	for i := 1; i < max; i++ {

		key := fmt.Sprint(i)

		tm.Set(key, key)
		val := tm.Get(key)
		assert.Equal(t, key, val)
	}
}

// HasPrefix 找到
func Test_TrieMap_HasPrefix(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		assert.True(t, tm.HasPrefix(key[:i]))
	}

}

// HasPrefix 找不到
func Test_TrieMap_HasPrefix_notFound(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		assert.True(t, tm.HasPrefix(key[:i]))
	}

	assert.False(t, tm.HasPrefix("/ha"))
}

func Test_TrieMap_GetWithBool_notFound(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		assert.True(t, tm.HasPrefix(key[:i]))
	}
	_, ok := tm.GetWithBool("/ha")
	assert.False(t, ok)
	_, ok = tm.GetWithBool("/he")
	assert.False(t, ok)
}

func Test_TrieMap_Delete(t *testing.T) {

	tm := New[string]()
	max := 1000

	for i := 1; i < max; i++ {

		key := fmt.Sprint(i)

		tm.Set(key, key)
		val := tm.Get(key)
		assert.Equal(t, key, val)
		tm.Delete(key)
		val, ok := tm.GetWithBool(key)
		assert.False(t, ok)
		assert.Equal(t, "", val)
	}

	key := fmt.Sprint(max + 1)
	tm.Delete(key)
	val, ok := tm.GetWithBool(key)
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

// 删除长的
func Test_TrieMap_Delete2(t *testing.T) {

	tm := New[string]()

	tm.Set("/1", "/1")
	tm.Set("/12", "/12")
	tm.Delete("/12")
	assert.Equal(t, tm.Get("/12"), "")
	assert.Equal(t, tm.Get("/1"), "/1")
}

// 删除短的
func Test_TrieMap_Delete3(t *testing.T) {

	tm := New[string]()

	tm.Set("/1", "/1")
	tm.Set("/12", "/12")
	tm.Delete("/1")
	assert.Equal(t, tm.Get("/12"), "/12")
	assert.Equal(t, tm.Get("/1"), "")
}

// 删除带中文
func Test_TrieMap_Delete4(t *testing.T) {

	tm := New[string]()

	tm.Set("中", "中")
	tm.Set("中国", "中国")
	tm.Delete("中")
	assert.Equal(t, tm.Get("中国"), "中国")
	assert.Equal(t, tm.Get("中"), "")
}

func Test_TrieMap_Delete5(t *testing.T) {

	tm := New[string]()

	tm.Set("中", "中")
	tm.Set("中国", "中国")
	tm.Delete("中")
	tm.Delete("中国")
	assert.Equal(t, tm.Get("中国"), "")
	assert.Equal(t, tm.Get("中"), "")
}

func Test_TrieMap_Delete6(t *testing.T) {

	tm := New[string]()

	tm.Set("/1", "/1")
	tm.Set("/12", "/12")
	tm.Set("/13", "/13")
	tm.Delete("/12")
	assert.Equal(t, tm.Get("/12"), "")
	assert.Equal(t, tm.Get("/1"), "/1")
	assert.Equal(t, tm.Get("/13"), "/13")
}
