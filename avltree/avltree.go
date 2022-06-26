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

type root[K constraints.Ordered, V any] struct {
	node *Node[K, V]
}

func (r *root[K, V]) childReplace(oldNode, newNode, parent *Node[K, V]) {
	if parent != nil {
		if parent.left == oldNode {
			parent.left = newNode
		} else {
			parent.right = newNode
		}
	} else {
		r.node = newNode
	}

}

// 左旋就是拽住node往下拉, node.right升为父节点
func (r *root[K, V]) rotateLeft(node *Node[K, V]) *Node[K, V] {
	right := node.right
	parent := node.parent

	// node会滑成right的左节点
	// 这里安排下node.right的位置, 这里不再指向right, 再向right的左孩子
	// right.left 大于node, 小于right, 所以新的位置就是node.right
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}

	// 把node从父的位置降下来
	right.left = node
	right.parent = parent
	r.childReplace(node, right, parent)
	node.parent = right
	return right
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
