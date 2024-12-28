package radix

// apache 2.0 antlabs
import (
	"strings"
	"unicode/utf8"

	"github.com/antlabs/gstl/api"
	"github.com/antlabs/gstl/cmp"
	"github.com/antlabs/gstl/vec"
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
	v, _ = r.TryGet(k)
	return
}

// 返回共同的前缀
func commonPrefix(k1, k2 string) (i int) {
	min := cmp.Min(len(k1), len(k2))

	for i = 0; i < min; i++ {
		if k1[i] != k2[i] {
			return i
		}
	}
	return
}

func (r *Radix[V]) newEdge(label rune, p pair[V], prefix string) edge[V] {
	return edge[V]{
		label: label,
		node: &node[V]{
			pair:   p,
			prefix: prefix,
		},
	}
}

// 设置
func (r *Radix[V]) Swap(k string, v V) (prev V, replaced bool) {

	var parent *node[V]
	var found bool
	n := r.root
	remaining := k

	for {

		if len(k) == 0 {
			if n.isSet {
				prev = n.val
				n.val = v
				return prev, true
			}

			n.key, n.val = k, v
			n.isSet = true
			r.length++
			replaced = true
			return
		}

		rune, _ := utf8.DecodeLastRuneInString(remaining)
		parent = n
		n, found = n.children(rune)
		if !found {
			parent.edges.Push(r.newEdge(rune, pair[V]{
				key:   k,
				val:   v,
				isSet: true,
			}, remaining))
			r.length++
			return
		}

		// 待插入节点 貌似和当前节点有共同的路径，先continue看看情况，等会插入
		commonPrefixLen := commonPrefix(remaining, n.prefix)
		if commonPrefixLen == len(n.prefix) {
			remaining = remaining[commonPrefixLen:]
			continue
		}

		subRune, _ := utf8.DecodeLastRuneInString(remaining[commonPrefixLen:])
		// 这里遇到分叉
		// 这里的节点只加上，比如原来节点是/helloaxx, 现在要插入/hellobxx
		// /hello 会变成child
		r.length++
		child := &node[V]{
			// 共同路径成为两个分裂节点的父节点
			prefix: remaining[:commonPrefixLen],
		}

		// axx 变成child的儿子1
		child.edges.Push(edge[V]{
			label: subRune,
			node:  n,
		})
		// 把/helloaxx 变成axx
		n.prefix = n.prefix[commonPrefixLen:]

		// 如果以前 parent指向/helloaxx，那么现在parent就指向/hello(即child节点)
		// axx 和bxx都将成为/hello的两个子节点
		parent.setChildren(rune, child)

		remaining = remaining[commonPrefixLen:]
		pairKV := pair[V]{
			key:   k,
			val:   v,
			isSet: true,
		}

		// 如果新插入路径只是原路径的子集，就走起
		// 比如原路径是/helloaxx, 本次插入/hello
		if len(remaining) == 0 {
			child.pair = pairKV
			return
		}

		subRune, _ = utf8.DecodeLastRuneInString(remaining)

		// 把bxx的路径是在newEdge函数里面设置的
		child.insertChildren(subRune, r.newEdge(subRune, pairKV, remaining))
	}

}

// 是否有这个前缀串
func (r *Radix[V]) HasPrefix(k string) (ok bool) {
	return
}

func (n *node[V]) insertChildren(r rune, new edge[V]) {

	index, found := n.find(r)
	if found {
		panic("这个节点已经存在过???")
	}

	n.edges.Insert(index, new)
}

func (n *node[V]) setChildren(r rune, new *node[V]) {
	index, found := n.find(r)
	if !found {
		panic("没找到这个节点")
	}

	n.edges.GetPtr(index).node = new
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
func (r *Radix[V]) TryGet(k string) (v V, found bool) {
	n := r.root

	for {

		// k 消费完，找到，或者找不到
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
