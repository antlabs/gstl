package set

import (
	"github.com/guonaihong/gstl/api"
	"github.com/guonaihong/gstl/rbtree"
	"golang.org/x/exp/constraints"
)

type Set[K constraints.Ordered] struct {
	api.SortedMap[K, struct{}]
}

func New[K constraints.Ordered]() *Set[K] {
	// 随手使用rbtree，后面压测再决定使用
	return &Set[K]{SortedMap: rbtree.New[K, struct{}]()}
}

func (s *Set[K]) Set(k K) {
	s.SortedMap.Set(k, struct{}{})
}

func (s *Set[K]) Len() int {
	return s.SortedMap.Len()
}

func (s *Set[K]) Clone() Set[K] {

	//	api.SortedMap
	return
}
