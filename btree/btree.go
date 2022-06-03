package btree

import (
	"errors"
	"github.com/guonaihong/gstl/must"
	"github.com/guonaihong/gstl/vec"
	"golang.org/x/exp/constraints"
)

var ErrNotFound = errors.New("btree not found")

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

// 找到err为nil
// 找不到err为ErrNotFound
func (b *Btree[K, V]) Get(k K) (v V, err error) {
	if b.root == nil {
		err = ErrNotFound
		return
	}

	n := b.root
	for {
		i, found := b.find(n, k)
		if found {
			return n.items.Get(i).val, nil
		}

		if n.leaf() {
			err = ErrNotFound
			return
		}

		n = (*n.children)[i]
	}
}

func (b *Btree[K, V]) Delete(k K) *Btree[K, V] {
	if b.root == nil {
		return b
	}

	return b
}

func (b *Btree[K, V]) delete(n *node[K, V]) {

	var i int
	if n.children.Get(i).items.Len() < b.minItems {
		b.rebalance(n, i)
	}
}

func (b *Btree[K, V]) rebalance(n *node[K, V], i int) {
	if i == n.items.Len() {
		panic("看下什么情况会触发")
	}

	left, right := n.children.Get(i), n.children.Get(i+1)

	if left.items.Len()+right.items.Len() < b.maxItems {
		left.items.Push(n.items.Get(i))
		left.items.Append(right.items)

		if !left.leaf() {
			left.children.Append(right.children)
		}
		n.items.Remove(i)
		n.children.Remove(i + 1)
	} else if left.items.Len() > right.items.Len() {

		right.items.Insert(0, n.items.Get(i))
		last, err := left.items.Pop()
		if err != nil {
			panic(err.Error())
		}

		n.items.Set(i, last)

		if !left.leaf() {

			last, err := left.children.Pop()
			if err != nil {
				panic(err.Error())
			}
			right.children.Insert(0, last)
		}
	} else {

		left.items.Push(n.items.Get(i))
		first, err := right.items.PopFront()
		if err != nil {
			panic(err.Error())
		}

		n.items.Set(i, first)

		if !left.leaf() {
			first, err := right.children.PopFront()
			if err != nil {
				panic(err.Error())
			}

			left.children.Push(first)
		}
	}
}

// 遍历b tree
func (b *Btree[K, V]) Range(callback func(k K, v V) bool) *Btree[K, V] {
	// 遍历
	if b.root == nil {
		return b
	}

	b.root.rangeInner(callback)
	return b
}

// 遍历b tree
func (n *node[K, V]) rangeInner(callback func(k K, v V) bool) bool {

	if n.leaf() {

		cycle := false
		n.items.Range(func(_ int, p pair[K, V]) (rv bool) {
			defer func() {
				cycle = rv
			}()

			return callback(p.key, p.val)
		})

		return cycle
	}

	cycle := false
	n.items.Range(func(index int, p pair[K, V]) (rv bool) {
		defer func() {
			cycle = rv
		}()
		// 这里不停递归向叶子节点的方向走去
		return n.children.Get(index).rangeInner(func(k K, v V) bool {
			return callback(p.key, p.val)
		})

	})

	if !cycle {
		return false
	}

	return must.TakeOne(n.children.Last()).rangeInner(callback)
}
