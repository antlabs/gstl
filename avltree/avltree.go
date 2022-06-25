package avltree

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var ErrNotFound = errors.New("avltree: not found value")

type Node[K constraints.Ordered, V any] struct {
	left   *Node[K, V]
	right  *Node[K, V]
	parent *Node[K, V]
	pair[K, V]
	height int
}

// 元素
type pair[K constraints.Ordered, V any] struct {
	val V
	key K
}

type AvlTree[K constraints.Ordered, V any] struct {
	length int
	root   *Node[K, V]
}

func (a *AvlTree[K, V]) New() {

}

func (a *AvlTree[K, V]) First() (v V) {
	return
}

func (a *AvlTree[K, V]) Last() (v V) {
	return
}

// 从avl tree找到需要的值
func (a *AvlTree[K, V]) Get(k K) (v V, err error) {
	n := a.root
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
