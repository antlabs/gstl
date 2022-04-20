package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	n := New[int]()
	assert.NotNil(t, n)
}
