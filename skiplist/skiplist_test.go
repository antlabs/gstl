package skiplist

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	n := New[string](strings.Compare)
	assert.NotNil(t, n)
}
