package api

import "golang.org/x/exp/constraints"

type Set[K constraints.Ordered, V any] interface {
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
}

type SortedSet[K constraints.Ordered, V any] interface {
	Set[K, V]
	TopMin(limit int, callback func(k K, v V) bool)
	TopMax(limit int, callback func(k K, v V) bool)
}
