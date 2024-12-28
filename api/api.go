package api

import "golang.org/x/exp/constraints"

type Map[K constraints.Ordered, V any] interface {
	// 获取
	Get(k K) (elem V)
	// 获取
	TryGet(k K) (elem V, ok bool)
	// 删除
	Delete(k K)
	// 设置
	Set(k K, v V)
	// 设置值
	SetWithPrev(k K, v V) (prev V, replaced bool)
	// int
	Len() int
	// 遍历
	Range(callback func(k K, v V) bool)
}

type SortedMap[K constraints.Ordered, V any] interface {
	Map[K, V]
	TopMin(limit int, callback func(k K, v V) bool)
	TopMax(limit int, callback func(k K, v V) bool)
}

// TODO
type Set[K constraints.Ordered] interface {
	Set(k K)
}

type Trie[V any] interface {
	Get(k string) (v V)
	SetWithPrev(k string, v V) (prev V, replaced bool)
	HasPrefix(k string) bool
	TryGet(k string) (v V, found bool)
	Delete(k string)
	Len() int
}

type CMaper[K comparable, V any] interface {
	Delete(key K)
	Load(key K) (value V, ok bool)
	LoadAndDelete(key K) (value V, loaded bool)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Range(f func(key K, value V) bool)
	Store(key K, value V)
}
