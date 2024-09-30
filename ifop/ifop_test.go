package ifop

// apache 2.0 antlabs
import (
	"reflect"
	"testing"
)

func TestIf(t *testing.T) {
	a := ""
	if result := If(len(a) == 0, "default"); result != "default" {
		t.Errorf("expected 'default', got '%s'", result)
	}
}

func TestIfElse(t *testing.T) {
	a := ""
	if result := IfElse(len(a) != 0, a, "default"); result != "default" {
		t.Errorf("expected 'default', got '%s'", result)
	}
	a = "hello"
	if result := IfElse(len(a) != 0, a, "default"); result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestIfElse2(t *testing.T) {
	o := map[string]any{"hello": "hello"}
	a := []any{"hello", "world"}

	if result := IfElseAny(o != nil, o, a); !reflect.DeepEqual(result, o) {
		t.Errorf("expected %v, got %v", o, result)
	}
	o = nil
	if result := IfElseAny(o != nil, o, a); !reflect.DeepEqual(result, a) {
		t.Errorf("expected %v, got %v", a, result)
	}
}
