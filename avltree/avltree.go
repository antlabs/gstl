package avltree

import (
	"errors"

	"github.com/guonaihong/gstl/cmp"
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
func (n *Node[K, V]) leftHeight() int {
	if n.left != nil {
		return n.left.height
	}

	return 0
}

// 返回右子树高度
func (n *Node[K, V]) rightHeight() int {
	if n.right != nil {
		return n.right.height
	}
	return 0
}

func (n *Node[K, V]) heightUpdate() {
	lh := n.left.leftHeight()
	rh := n.right.rightHeight()
	n.height = cmp.Max(lh, rh) + 1
}

type root[K constraints.Ordered, V any] struct {
	node *Node[K, V]
}

func (r *root[K, V]) fixLeft(node *Node[K, V]) *Node[K, V] {
	right := node.right
	rlh := node.right.leftHeight()
	rrh := node.right.rightHeight()

	if rlh > rrh {
		right = r.rotateRight(right)
		right.right.heightUpdate()
		right.heightUpdate()
	}
	node = r.rotateLeft(node)
	node.left.heightUpdate()
	node.heightUpdate()

	return node
}

func (r *root[K, V]) fixRight(node *Node[K, V]) *Node[K, V] {
	return node
}

func (r *root[K, V]) postInsert(node *Node[K, V]) {
	node.height = 1

	for node = node.parent; node != nil; node = node.parent {
		lh := node.leftHeight()
		lr := node.rightHeight()
		height := cmp.Max(lh, lr) + 1

		diff := lh - lr
		if node.height == height {
			break
		}
		node.height = height

		if diff <= -2 {
			node = r.fixLeft(node)
		} else {
			node = r.fixRight(node)
		}
	}
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

func (r *root[K, V]) rotateRight(node *Node[K, V]) *Node[K, V] {
	left := node.left
	parent := node.parent
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}

	left.right = node
	left.parent = parent
	r.childReplace(node, left, parent)
	node.parent = left

	return node
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
