package skiplist

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	n := New[string](strings.Compare)
	assert.NotNil(t, n)
}

func Test_SetGet(t *testing.T) {
	zset := New[string](strings.Compare)
	max := 100.0
	for i := 0.0; i < max; i++ {
		zset.Set(i, fmt.Sprintf("%d", int(i)))
	}

	for i := 0.0; i < max; i++ {
		v := zset.GetOrZero(i)
		assert.Equal(t, v, fmt.Sprintf("%d", int(i)))
	}
}
