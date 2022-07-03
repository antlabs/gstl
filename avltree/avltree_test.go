package avltree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetAndGet(t *testing.T) {
	b := New[int, int]()
	max := 2
	for i := 0; i < max; i++ {
		b.SetWithPrev(i, i)
	}

	for i := 0; i < max; i++ {
		v, err := b.Get(i)
		assert.NoError(t, err)
		assert.Equal(t, v, i)
	}
}
