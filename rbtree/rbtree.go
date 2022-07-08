package rbtree

// apache 2.0 guonaihong
// 参考资料
// https://github.com/torvalds/linux/blob/master/lib/rbtree.c
import (
	"errors"

	"golang.org/x/exp/constraints"
)

// 红黑树5条重要性质
// 1. 节点为红色或者黑色
// 2. 根节点是黑色(黑根)
// 3. 所有叶子(空节点)均为黑色
// 4. 每个红色节点的两个子节点均为黑色(红父黑子)
// 5. 从根到叶的每个路径包含相同数量的黑色节点(黑高相同)

var ErrNotFound = errors.New("rbtree: not found value")

// 元素
type pair[K constraints.Ordered, V any] struct {
	val V
	key K
}

type node[K constraints.Ordered, V any] struct {
	left   *node[K, V]
	right  *node[K, V]
	parent *node[K, V]
	pair[K, V]
	// red/black 可以用一个变量实现, 为了代码简单清晰, 这里不考虑一个字节的收益, 放弃使用位操作
	red   bool
	black bool
}

func (n *node[K, V]) setParentRed(parent *node[K, V]) {
	parent.red = true
	n.parent = parent
}

func (n *node[K, V]) setParentBlack(parent *node[K, V]) {
	parent.black = true
	n.parent = parent
}

func (n *node[K, V]) link(parent *node[K, V], link **node[K, V]) {
	n.parent = parent
	n.red = true
	n.black = false
	*link = n
}

type root[K constraints.Ordered, V any] struct {
	node *node[K, V]
}

func (r *root[K, V]) insert(n *node[K, V]) {
}

// 红黑树
type RBtree[K constraints.Ordered, V any] struct {
	length int
	root   root[K, V]
}

// 第一个节点
func (r *RBtree[K, V]) First() (v V, err error) {
	n := r.root.node
	if n == nil {
		err = ErrNotFound
		return
	}

	for n.left != nil {
		n = n.left
	}

	return n.val, nil
}

// 最后一个节点
func (r *RBtree[K, V]) Last() (v V, err error) {
	n := r.root.node
	if n == nil {
		err = ErrNotFound
		return
	}

	for n.right != nil {
		n = n.right
	}

	return n.val, nil
}

// 设置
func (r *RBtree[K, V]) SetWithPrev(k K, v V) (prev V, replaced bool) {
	link := &r.root.node
	var parent *node[K, V]

	node := &node[K, V]{pair: pair[K, V]{key: k, val: v}}

	for *link != nil {
		parent = *link
		if parent.key == k {
			prev = parent.val
			parent.val = v
			return prev, true
		}

		if parent.key < k {
			*link = parent.right
		} else {
			*link = parent.left
		}
	}

	node.link(parent, link)
	r.root.insert(node)
	r.length++
	return
}
