package set

// apache 2.0 antlabs
import (
	"testing"
)

func Test_Range_New(t *testing.T) {
	s := New[string]()
	a := []string{"1111", "2222", "3333"}
	for _, v := range a {
		s.Set(v)
	}
	var got []string
	s.Range(func(k string) bool {
		got = append(got, k)
		return true
	})

	if !equalSlices(got, a) {
		t.Errorf("expected %v, got %v", a, got)
	}
}

func Test_Range_From(t *testing.T) {
	a := []string{"1111", "2222", "3333"}
	s := From(a...)
	for _, v := range a {
		s.Set(v)
	}
	var got []string
	s.Range(func(k string) bool {
		got = append(got, k)
		return true
	})

	if !equalSlices(got, a) {
		t.Errorf("expected %v, got %v", a, got)
	}
}

func Test_Len(t *testing.T) {
	s := New[int]()
	max := 1000
	for i := 0; i < max; i++ {
		s.Set(i)
	}
	if s.Len() != max {
		t.Errorf("expected %d, got %d", max, s.Len())
	}
}

func Test_Equal(t *testing.T) {
	s := New[int]()
	max := 1000
	for i := 0; i < max; i++ {
		s.Set(i)
	}
	if s.Len() != max {
		t.Errorf("expected %d, got %d", max, s.Len())
	}

	s2 := s.Clone()

	if s.Len() != s2.Len() {
		t.Errorf("expected %d, got %d", s.Len(), s2.Len())
	}
	if !s.Equal(s2) {
		t.Errorf("expected sets to be equal")
	}
}

func Test_Not_Equal(t *testing.T) {
	s := New[int]()
	s2 := New[int]()
	max := 1000

	for i := 0; i < max; i++ {
		s.Set(i)
	}

	if s.Len() != max {
		t.Errorf("expected %d, got %d", max, s.Len())
	}

	for i := 0; i < max; i++ {
		s2.Set(i - 1)
	}

	if s.Len() != s2.Len() {
		t.Errorf("expected %d, got %d", s.Len(), s2.Len())
	}
	if s.Equal(s2) {
		t.Errorf("expected sets to be not equal")
	}
}

func Test_Not_Equal2(t *testing.T) {
	s := New[int]()
	s2 := New[int]()
	max := 1000

	for i := 0; i < max; i++ {
		s.Set(i)
	}

	if s.Len() != max {
		t.Errorf("expected %d, got %d", max, s.Len())
	}

	for i := 0; i < max/2; i++ {
		s2.Set(i - 1)
	}

	if s.Equal(s2) {
		t.Errorf("expected sets to be not equal")
	}
}

func Test_IsMember(t *testing.T) {
	s := New[int]()
	for i := 0; i < 10; i++ {
		s.Set(i)
	}
	if !s.IsMember(1) {
		t.Errorf("expected true, got false")
	}
}

func Test_Union(t *testing.T) {
	s := From("1111")
	s1 := From("2222")
	s2 := From("3333")

	newSet := s.Union(s1, s2)
	expected := []string{"1111", "2222", "3333"}
	if !equalSlices(newSet.ToSlice(), expected) {
		t.Errorf("expected %v, got %v", expected, newSet.ToSlice())
	}
}

func Test_Diff(t *testing.T) {
	s := From("hello", "world", "1234", "4567")
	s2 := From("1234", "4567")

	newSet := s.Diff(s2)
	expected := []string{"hello", "world"}
	if !equalSlices(newSet.ToSlice(), expected) {
		t.Errorf("expected %v, got %v", expected, newSet.ToSlice())
	}
}

func Test_Intersection(t *testing.T) {
	s := From("1234", "5678", "9abc")
	s2 := From("abcde", "5678", "9abc")

	v := s.Intersection(s2).ToSlice()
	expected := []string{"5678", "9abc"}
	if !equalSlices(v, expected) {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func Test_IsSubset(t *testing.T) {
	s := From("5678", "9abc")
	s2 := From("abcde", "5678", "9abc")

	if !s.IsSubset(s2) {
		t.Errorf("expected true, got false")
	}
}

func Test_IsSubset_Not(t *testing.T) {
	s := From("aa", "5678", "9abc")
	s2 := From("abcde", "5678", "9abc")

	if s.IsSubset(s2) {
		t.Errorf("expected false, got true")
	}
}

func Test_IsSubset_Not2(t *testing.T) {
	s := From("aa", "5678", "9abc", "33333")
	s2 := From("abcde", "5678", "9abc")

	if s.IsSubset(s2) {
		t.Errorf("expected false, got true")
	}
}

func Test_IsSuperset(t *testing.T) {
	s2 := From("5678", "9abc")
	s := From("abcde", "5678", "9abc")

	if !s.IsSuperset(s2) {
		t.Errorf("expected true, got false")
	}
}

func Test_IsSuperset_Not(t *testing.T) {
	s2 := From("aa", "5678", "9abc")
	s := From("abcde", "5678", "9abc")

	if s.IsSuperset(s2) {
		t.Errorf("expected false, got true")
	}
}

func Test_IsSuperset_Not2(t *testing.T) {
	s2 := From("aa", "5678", "9abc", "33333")
	s := From("abcde", "5678", "9abc")

	if s.IsSuperset(s2) {
		t.Errorf("expected false, got true")
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
