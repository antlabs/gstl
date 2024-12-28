package cmap

import (
	"github.com/antlabs/gstl/api"
	"golang.org/x/exp/constraints"
)

var _ api.Map[int, int] = (*stdmap[int, int])(nil)

type stdmap[K constraints.Ordered, V any] struct {
	m map[K]V
}

func newStdMap[K constraints.Ordered, V any]() *stdmap[K, V] {
	return &stdmap[K, V]{m: make(map[K]V)}
}

func (s *stdmap[K, V]) Get(key K) (elem V) {
	elem, _ = s.m[key]
	return
}

// 获取
func (s *stdmap[K, V]) TryGet(key K) (elem V, ok bool) {
	elem, ok = s.m[key]
	return
}

// 删除
func (s *stdmap[K, V]) Delete(key K) {
	delete(s.m, key)
}

// 设置
func (s *stdmap[K, V]) Set(key K, value V) {
	s.m[key] = value
}

// 设置值
func (s *stdmap[K, V]) SetWithPrev(key K, value V) (prev V, replaced bool) {
	prev, replaced = s.m[key]
	s.m[key] = value
	return
}

// int
func (s *stdmap[K, V]) Len() int {
	return len(s.m)
}

// 遍历
func (s *stdmap[K, V]) Range(callback func(k K, v V) bool) {
	for k, v := range s.m {
		if !callback(k, v) {
			return
		}
	}
}
