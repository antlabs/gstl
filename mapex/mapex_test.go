package mapex

import (
	"sort"
	"testing"
)

func Test_Keys(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "1"
	m["b"] = "2"
	m["c"] = "3"
	get := Keys(m)
	sort.Strings(get)
	expected := []string{"a", "b", "c"}
	if !equalSlices(get, expected) {
		t.Errorf("expected %v, got %v", expected, get)
	}
	get = Map[string, string](m).Keys()
	sort.Strings(get)
	if !equalSlices(get, expected) {
		t.Errorf("expected %v, got %v", expected, get)
	}
}

func Test_Values(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "1"
	m["b"] = "2"
	m["c"] = "3"
	get := Values(m)
	sort.Strings(get)
	expected := []string{"1", "2", "3"}
	if !equalSlices(get, expected) {
		t.Errorf("expected %v, got %v", expected, get)
	}

	get = Map[string, string](m).Values()
	sort.Strings(get)
	if !equalSlices(get, expected) {
		t.Errorf("expected %v, got %v", expected, get)
	}
}

// 辅助函数，用于比较两个切片是否相等
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
