package set

import (
	"github.com/guonaihong/gstl/api"
	"github.com/guonaihong/gstl/rbtree"
	"golang.org/x/exp/constraints"
)

type Set[K constraints.Ordered] struct {
	api.SortedMap[K, struct{}]
}

// 创建一个空的slice
func New[K constraints.Ordered]() *Set[K] {
	// 随手使用rbtree，后面压测再决定使用
	return &Set[K]{SortedMap: rbtree.New[K, struct{}]()}
}

// 从slice创建set
func From[K constraints.Ordered](s []K) *Set[K] {
	var b rbtree.RBTree[K, struct{}]
	for _, v := range s {
		b.Set(v, struct{}{})
	}

	return &Set[K]{SortedMap: &b}
}

// 给集合添加元素
func (s *Set[K]) Set(k K) {
	s.SortedMap.Set(k, struct{}{})
}

// 返回集合中元素的个数
func (s *Set[K]) Len() int {
	return s.SortedMap.Len()
}

// 深度复制一个集合
func (s *Set[K]) Clone() api.Set[K] {
	//	api.SortedMap
	return s
}

// 测试k是否在集合中
func (s *Set[K]) IsMember(k K) (b bool) {
	return
}

func (s *Set[K]) Diff(s1 *Set[K]) (new *Set[K]) {

	return
}

func (s *Set[K]) Union(s1 *Set[K]) (new *Set[K]) {

	return
}

func (s *Set[K]) Intersection(s1 *Set[K]) (new *Set[K]) {
	return
}

// 测试集合每个是否在s1里面
func (s *Set[K]) IsSubset(s1 Set[K]) (b bool) {
	return
}
