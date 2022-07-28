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

		fmt.Println(key[:i])
		assert.True(t, tm.HasPrefix(key[:i]))
	}
	_, ok := tm.GetWithBool("/ha")
	assert.False(t, ok)
}

func Test_TrieMap_Delete(t *testing.T) {

}
