package avltree

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var ErrNotFound = errors.New("avltree: not found value")

// 元素
type pair[K constraints.Ordered, V any] struct {
	val V
	key K
}

type Node[K constraints.Ordered, V any] struct {
	left   *Node[K, V]
	right  *Node[K, V]
	parent *Node[K, V]
	pair[K, V]
	height int
}

type root[K constraints.Ordered, V any] struct {
	node *Node[K, V]
}

// 返回左子树高度
func (n *Node[K, V]) LeftHeight() int {
	if n.left != nil {
		return n.left.height
	}

	return 0
}

// 返回右子树高度
func (n *Node[K, V]) RightHeight() int {
	if n.right != nil {
		return n.right.height
	}
	return 0
}

// avl tree的结构
type AvlTree[K constraints.Ordered, V any] struct {
	length int
	root   root[K, V]
}

// 构造函数
func New[K constraints.Ordered, V any]() *AvlTree[K, V] {
	return &AvlTree[K, V]{}
}

// 第一个节点
func (a *AvlTree[K, V]) First() (v V, err error) {
	n := a.root.node
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
func (a *AvlTree[K, V]) Last() (v V, err error) {
	n := a.root.node
	if n == nil {
		err = ErrNotFound
		return
	}

	for n.right != nil {
		n = n.right
	}

	return n.val, nil
}

// 从avl tree找到需要的值
func (a *AvlTree[K, V]) Get(k K) (v V, err error) {
	n := a.root.node
	for n != nil {
		if n.key == k {
			return n.val, nil
		}

		if k > n.key {
			n = n.right
		} else {
			n = n.left
		}
	}

	err = ErrNotFound
	return
}
