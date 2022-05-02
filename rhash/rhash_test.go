package rhash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1. set get测试
// key string value bool
func Test_SetGet_StringBool(t *testing.T) {
	ht := New[string, bool]()
	ht.Set("hello", true)
	ht.Set("world", true)
	ht.Set("ni", true)
	ht.Set("hao", true)

	assert.True(t, ht.GetOrZero("hello"))
	assert.True(t, ht.GetOrZero("world"))
	assert.True(t, ht.GetOrZero("ni"))
	assert.True(t, ht.GetOrZero("hao"))
}

// 1. set get测试
// key string, value string
func Test_SetGet_StringString(t *testing.T) {
	ht := New[string, string]()
	ht.Set("hello", "hello")
	ht.Set("world", "world")
	ht.Set("ni", "ni")
	ht.Set("hao", "hao")

	assert.Equal(t, ht.GetOrZero("hello"), "hello")
	assert.Equal(t, ht.GetOrZero("world"), "world")
	assert.Equal(t, ht.GetOrZero("ni"), "ni")
	assert.Equal(t, ht.GetOrZero("hao"), "hao")
}

// set get测试
// key string value string
func Test_SetGet_IntString(t *testing.T) {
	ht := New[int, string]()
	ht.Set(1, "hello")
	ht.Set(2, "world")
	ht.Set(3, "ni")
	ht.Set(4, "hao")

	assert.Equal(t, ht.GetOrZero(1), "hello")
	assert.Equal(t, ht.GetOrZero(2), "world")
	assert.Equal(t, ht.GetOrZero(3), "ni")
	assert.Equal(t, ht.GetOrZero(4), "hao")
}

// set get测试
// 设计重复key
func Test_SetGet_Replace_IntString(t *testing.T) {
	ht := New[int, string]()
	ht.Set(1, "hello")
	ht.Set(1, "world")

	assert.Equal(t, ht.GetOrZero(1), "world")
}
