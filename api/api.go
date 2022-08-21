package api

import "golang.org/x/exp/constraints"

type Map[K constraints.Ordered, V any] interface {
	// 获取
	Get(k K) (elem V)
	// 获取
	GetWithErr(k K) (elem V, err error)
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
	GetWithBool(k string) (v V, found bool)
	Delete(k string)
	Len() int
}
