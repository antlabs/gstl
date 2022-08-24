package radix

import (
	"github.com/guonaihong/gstl/api"
	"github.com/guonaihong/gstl/vec"
)

var _ api.Trie[int] = (*Radix[int])(nil)

// 健值对
type pair[V any] struct {
	val V
	key string
}

// 边
type edge[V any] struct {
	label rune
	node  *node[V]
}

// 节点
type node[V any] struct {
	pair[V]
	prefix string
	edges  vec.Vec[edge[V]]
}

// 头节点
type Radix[V any] struct {
	root   *node[V]
	length int
}

// 获取
func (r *Radix[V]) Get(k string) (v V) {

	return
}

// 设置
func (r *Radix[V]) SetWithPrev(k string, v V) (prev V, replaced bool) {
	/*
		n := r.root
		var parent *node

		for _, rune := range k {
			n.edges.SearchFunc()
		}
	*/

	return
}

// 是否有这个前缀串
func (r *Radix[V]) HasPrefix(k string) (ok bool) {
	return
}

// 获取返回bool
func (r *Radix[V]) GetWithBool(k string) (v V, found bool) {
	return
}

// 删除
func (r *Radix[V]) Delete(k string) {

}

// 返回长度
func (r *Radix[V]) Len() int {

	return r.length
}
