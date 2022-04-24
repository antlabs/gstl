package skiplist

import (
	"fmt"
	"strings"
	"testing"

	"github.com/guonaihong/gstl/cmp"
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

func Test_SetGetRemove(t *testing.T) {
	zset := New[float64](cmp.Compare[float64])

	for i := 0.0; i < 100.0; i++ {
		zset.Set(i, i)
	}

	for i := 0.0; i < 100.0; i++ {
		zset.Remove(i)
		for j := 0.0; i < 100.0; j++ {
			if j == i {
				continue
			}
			v, err := zset.Get(j)
			assert.NoError(t, err)
			assert.Equal(t, v, j)
		}
		zset.Set(i, i)
	}
}
