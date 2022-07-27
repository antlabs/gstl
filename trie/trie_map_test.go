package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func Test_TrieMap_HashPrefix(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		fmt.Println(key[:i])
		assert.True(t, tm.HasPrefix(key[:i]))
	}

}
