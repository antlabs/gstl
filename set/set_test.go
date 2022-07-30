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
	s := From(a)
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
