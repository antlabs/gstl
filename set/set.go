package set

// apache 2.0 antlabs
import (
	"github.com/antlabs/gstl/api"
	"github.com/antlabs/gstl/rbtree"
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
func From[K constraints.Ordered](s ...K) *Set[K] {
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

func (s *Set[K]) ToSlice() (new []K) {
	new = make([]K, 0, s.Len())
	s.Range(func(k K) bool {
		new = append(new, k)
		return true
	})
	return
}

// 深度复制一个集合
func (s *Set[K]) Clone() (new *Set[K]) {
	new = New[K]()
	s.Range(func(k K) bool {
		new.Set(k)
		return true
	})
	return
}

// 测试k是否在集合中
func (s *Set[K]) IsMember(k K) (b bool) {
	_, b = s.TryGet(k)
	return
}

// 返回的是s1没有的元素, s - s1
func (s *Set[K]) Diff(s1 *Set[K]) (new *Set[K]) {

	new = New[K]()
	s.Range(func(k K) bool {
		if !s1.IsMember(k) {
			new.Set(k)
		}
		return true
	})
	return
}

// 返回两个集合的所有元素
func (s *Set[K]) Union(sets ...*Set[K]) (new *Set[K]) {

	new = New[K]()
	s.Range(func(k K) bool {
		new.Set(k)
		return true
	})

	for _, s1 := range sets {
		s1.Range(func(k K) bool {
			new.Set(k)
			return true
		})
	}

	return
}

// 返回两个集合的公共集合
func (s *Set[K]) Intersection(s1 *Set[K]) (new *Set[K]) {
	if s.Len() >= s1.Len() {
		s, s1 = s1, s
	}

	new = New[K]()
	s.Range(func(k K) bool {
		if s1.IsMember(k) {
			new.Set(k)
		}
		return true
	})
	return
}

// 测试集合s每个元素是否在s1里面, s <= s1
func (s *Set[K]) IsSubset(s1 *Set[K]) (b bool) {
	if s.Len() > s1.Len() {
		return false
	}

	b = true
	s.Range(func(k K) bool {
		if !s1.IsMember(k) {
			b = false
			return false
		}
		return true
	})
	return
}

// 测试集合s1每个元素是否在s里面 s1 <= s
func (s *Set[K]) IsSuperset(s1 *Set[K]) (b bool) {
	return s1.IsSubset(s)
}

// 遍历
func (s *Set[K]) Range(cb func(k K) bool) {
	s.SortedMap.Range(func(k K, _ struct{}) bool {
		return cb(k)
	})
}

// 两个集合是否相等
func (s *Set[K]) Equal(s1 *Set[K]) (b bool) {
	if s.Len() != s1.Len() {
		return false
	}

	b = true
	s.Range(func(k K) bool {
		_, b = s1.TryGet(k)
		return b
	})

	return
}
