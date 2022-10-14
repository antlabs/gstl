package ifop

// apache 2.0 guonaihong
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIf(t *testing.T) {
	a := ""
	assert.Equal(t, If(len(a) == 0, "default"), "default")
}

func TestIfElse(t *testing.T) {
	a := ""
	assert.Equal(t, IfElse(len(a) != 0, a, "default"), "default")
	a = "hello"
	assert.Equal(t, IfElse(len(a) != 0, a, "default"), "hello")
}
