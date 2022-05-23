package btree

import (
	"github.com/guonaihong/gstl/vec"
	"golang.org/x/exp/constraints"
)

//btree头结点
type Btree[K constraints.Ordered, V any] struct {
	count int         //当前元素个数
	root  *node[K, V] // root结点指针
}

// 元素
type pair[K constraints.Ordered, V any] struct {
	val V
	key K
}

// btree树的结点的组成
type node[K constraints.Ordered, V any] struct {
	items    *vec.Vec[pair[K, V]]  //存放元素的节点
	children *vec.Vec[*node[K, V]] //孩子节点
}

// 设置接口, 如果有这个值, 有值就替换, 没有就新加
func (b *Btree[K, V]) Set(k K, v V) *Btree[K, V] {

	_, _ = b.SetWithOld(k, v)
	return b
}

// 新建一个节点
func (b *Btree[K, V]) newNode(leaf bool) (n *node[K, V]) {
	n = &node[K, V]{}
	if !leaf {
		n.children = vec.New[*node[K, V]]()
	}
	return
}

// 新建叶子节点
func (b *Btree[K, V]) newLeaf() *node[K, V] {
	return b.newNode(true)
}

func (b *Btree[K, V]) find(n *node[K, V], key K) (index int, found bool) {

	index = n.items.SearchFunc(func(elem pair[K, V]) bool { return key <= elem.key })
	if index > 0 && n.items.Get(index-1).key > key {
		return index - 1, true
	}

	return index, false
}

//
func (b *Btree[K, V]) nodeSet() (old V, needSplit bool) {
	return
}

// 设置接口, 如果有值, 把old值带返回, 并且被替换, 没有就新加
func (b *Btree[K, V]) SetWithOld(k K, v V) (old V, replaced bool) {

	// 如果是每一个节点, 直接加入到root节点
	if b.root == nil {
		b.root = b.newLeaf()
		b.root.items.Push(pair[K, V]{key: k, val: v})
		b.count = 1
		return
	}

	return
}

func (b *Btree[K, V]) Get() {

}
