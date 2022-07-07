package rbtree

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// 红黑树5条重要性质
// 1. 节点为红色或者黑色
// 2. 根节点是黑色
// 3. 所有叶子(空节点)均为黑色
// 4. 每个红色节点的两个子节点均为黑色
// 5. 从根到叶的每个简单路径包含相同数量的黑色节点

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

type root[K constraints.Ordered, V any] struct {
	node *node[K, V]
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
