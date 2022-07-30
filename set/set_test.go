package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, got, a)
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

	assert.Equal(t, got, a)
}

func Test_Len(t *testing.T) {
	s := New[int]()
	max := 1000
	for i := 0; i < max; i++ {
		s.Set(i)
	}
	assert.Equal(t, s.Len(), max)
}

func Test_Equal(t *testing.T) {
	s := New[int]()
	max := 1000
	for i := 0; i < max; i++ {
		s.Set(i)
	}
	assert.Equal(t, s.Len(), max)

	s2 := s.Clone()

	assert.Equal(t, s.Len(), s2.Len())
	assert.True(t, s.Equal(s2))
}

func Test_Not_Equal(t *testing.T) {
	s := New[int]()
	s2 := New[int]()
	max := 1000

	for i := 0; i < max; i++ {
		s.Set(i)
	}

	assert.Equal(t, s.Len(), max)

	for i := 0; i < max; i++ {
		s2.Set(i - 1)
	}

	assert.Equal(t, s.Len(), s2.Len())
	assert.False(t, s.Equal(s2))
}

func Test_Not_Equal2(t *testing.T) {
	s := New[int]()
	s2 := New[int]()
	max := 1000

	for i := 0; i < max; i++ {
		s.Set(i)
	}

	assert.Equal(t, s.Len(), max)

	for i := 0; i < max/2; i++ {
		s2.Set(i - 1)
	}

	assert.False(t, s.Equal(s2))
}

func Test_IsMember(t *testing.T) {
	s := New[int]()
	for i := 0; i < 10; i++ {
		s.Set(i)
	}
	assert.True(t, s.IsMember(1))
}

func Test_Union(t *testing.T) {
	s := From("1111")
	s1 := From("2222")
	s2 := From("3333")

	newSet := s.Union(s1, s2)
	assert.Equal(t, newSet.ToSlice(), []string{"1111", "2222", "3333"})
}
