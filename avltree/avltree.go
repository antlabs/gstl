package avltree

import "golang.org/x/exp/constraints"

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
