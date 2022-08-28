package radix

import (
	"strings"
	"unicode/utf8"

	"github.com/guonaihong/gstl/api"
	"github.com/guonaihong/gstl/vec"
)

var _ api.Trie[int] = (*Radix[int])(nil)

// 健值对
type pair[V any] struct {
	val   V
	key   string
	isSet bool
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
	v, _ = r.GetWithBool(k)
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

func (n *node[V]) children(r rune) (v *node[V], found bool) {
	index, found := n.find(r)
	if !found {
		return
	}
	return n.edges.Get(index).node, true
}

func (n *node[V]) find(r rune) (index int, found bool) {

	index = n.edges.SearchFunc(func(elem edge[V]) bool { return r < elem.label })
	if index > 0 && n.edges.Get(index-1).label >= r {
		return index - 1, true
	}

	return index, false
}

// 获取返回bool
func (r *Radix[V]) GetWithBool(k string) (v V, found bool) {
	n := r.root

	for {

		// k 消费完，要不找到，要不找不到
		if len(k) == 0 {
			if n.isSet {
				return n.val, true
			}
			return
		}

		rune, _ := utf8.DecodeLastRuneInString(k)
		n, found = n.children(rune)
		if !found {
			return
		}

		if strings.HasPrefix(k, n.prefix) {
			k = k[len(n.prefix):]
			continue
		}
		return
	}

}

// 删除
func (r *Radix[V]) Delete(k string) {

}

// 返回长度
func (r *Radix[V]) Len() int {

	return r.length
}
