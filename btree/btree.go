package btree

import (
	"github.com/guonaihong/gstl/vec"
	"golang.org/x/exp/constraints"
)

//btree头结点
type Btree[K constraints.Ordered, V any] struct {
	count    int         //当前元素个数
	root     *node[K, V] // root结点指针
	maxItems int
	minItems int
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

func (n *node[K, V]) leaf() bool {
	return n.children == nil || n.children.Len() == 0
}

func New[K constraints.Ordered, V any](degree int) *Btree[K, V] {

	if degree == 0 {
		degree = 128 //拍脑袋给的, 需要压测下
	}

	maxItems := degree*2 - 1 // max items per node. max children is +1
	return &Btree[K, V]{
		maxItems: maxItems,
		minItems: maxItems / 2,
	}
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

// 分裂结点
func (b *Btree[K, V]) nodeSplit(n *node[K, V]) (right *node[K, V], median pair[K, V]) {
	i := b.maxItems / 2
	median = n.items.Get(i)

	// 新的左孩子就是n节点
	rightItems := n.items.SplitOff(i + 1)
	n.items.SetLen(n.items.Len() - 1)
	// 当前节点还有下层节点, 也要左右分家
	right = b.newNode(n.leaf())
	right.items = rightItems
	if !n.leaf() {
		right.children = n.children.SplitOff(i + 2)
	}

	return
}

//
func (b *Btree[K, V]) nodeSet(n *node[K, V], item pair[K, V]) (old V, replaced bool, needSplit bool) {
	i, found := b.find(n, item.key)
	// 找到位置直接替换
	if found {
		oldPtr := n.items.GetPtr(i)
		old = oldPtr.val
		oldPtr.val = item.val
		return old, true, false
	}

	// 如果是叶子节点
	if n.leaf() {
		// 没有位置插入新元素, 上层节点需要分裂
		if n.items.Len() == b.maxItems {
			needSplit = true
			return
		}
		n.items.Insert(i, item)
		return
	}

	old, replaced, needSplit = b.nodeSet(n.children.Get(i), item)
	if needSplit {
		// 没有位置插入新元素, 上层节点需要分裂
		if n.items.Len() == b.maxItems {

			needSplit = true
			return
		}

		right, median := b.nodeSplit(n.children.Get(i))
		n.children.Insert(i, right) // TODO debug
		n.items.Insert(i, median)   // TODO debug
	}
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
