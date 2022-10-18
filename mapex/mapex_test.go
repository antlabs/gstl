package mapex

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Keys(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "1"
	m["b"] = "2"
	m["c"] = "3"
	get := Keys(m)
	sort.Strings(get)
	assert.Equal(t, get, []string{"a", "b", "c"})
	get = Map[string, string](m).Keys()
	sort.Strings(get)
	assert.Equal(t, get, []string{"a", "b", "c"})
}

func Test_Values(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "1"
	m["b"] = "2"
	m["c"] = "3"
	get := Values(m)
	sort.Strings(get)
	assert.Equal(t, get, []string{"1", "2", "3"})

	get = Map[string, string](m).Values()
	sort.Strings(get)
	assert.Equal(t, get, []string{"1", "2", "3"})
}
