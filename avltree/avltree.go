package avltree

// apache 2.0 antlabs

// 参考资料
// https://github.com/skywind3000/avlmini
import (
	"errors"
	"fmt"

	"github.com/antlabs/gstl/api"
	"github.com/antlabs/gstl/cmp"
	"github.com/antlabs/gstl/vec"
	"golang.org/x/exp/constraints"
)

var ErrNotFound = errors.New("avltree: not found value")

var _ api.SortedMap[int, int] = (*AvlTree[int, int])(nil)

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
	height int
}

// 返回左子树高度
func (n *node[K, V]) leftHeight() int {
	if n.left != nil {
		return n.left.height
	}

	return 0
}

// 返回右子树高度
func (n *node[K, V]) rightHeight() int {
	if n.right != nil {
		return n.right.height
	}
	return 0
}

func (n *node[K, V]) heightUpdate() {
	lh := n.leftHeight()
	rh := n.rightHeight()
	n.height = cmp.Max(lh, rh) + 1
}

func (n *node[K, V]) link(parent *node[K, V], link **node[K, V]) {
	n.parent = parent
	*link = n
}

type root[K constraints.Ordered, V any] struct {
	node *node[K, V]
}

func (r *root[K, V]) fixLeft(node *node[K, V]) *node[K, V] {
	right := node.right
	// 右节点, 左子树高度
	rlh := right.leftHeight()
	// 右节点, 右子树高度
	rrh := right.rightHeight()

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

func (r *root[K, V]) fixRight(node *node[K, V]) *node[K, V] {
	left := node.left
	// 右节点, 左子树高度
	llh := left.leftHeight()
	// 右节点, 右子树高度
	lrh := left.rightHeight()

	if llh < lrh {
		left = r.rotateLeft(left)
		left.left.heightUpdate()
		left.heightUpdate()
	}

	node = r.rotateRight(node)
	if node.right != nil {
		node.right.heightUpdate()
	}

	node.heightUpdate()

	return node
}

func (r *root[K, V]) postInsert(node *node[K, V]) {
	node.height = 1

	for node = node.parent; node != nil; node = node.parent {
		lh := node.leftHeight()
		rh := node.rightHeight()
		height := cmp.Max(lh, rh) + 1

		diff := lh - rh
		if node.height == height {
			break
		}
		node.height = height

		if diff <= -2 {
			node = r.fixLeft(node)
		} else if diff >= 2 {
			node = r.fixRight(node)
		}
	}
}

func (r *root[K, V]) childReplace(oldNode, newNode, parent *node[K, V]) {
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

// 左旋就是拽住node往左下拉, node.right升为父节点
func (r *root[K, V]) rotateLeft(node *node[K, V]) *node[K, V] {
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

// 右旋就是拽往node往右下拉, node.left升为父节点
func (r *root[K, V]) rotateRight(node *node[K, V]) *node[K, V] {
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

// Get
func (a *AvlTree[K, V]) Get(k K) (v V) {
	v, _ = a.GetWithErr(k)
	return
}

// 从avl tree找到需要的值
func (a *AvlTree[K, V]) GetWithErr(k K) (v V, err error) {
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

func (a *AvlTree[K, V]) Set(k K, v V) {
	_, _ = a.SetWithPrev(k, v)
}

// 设置接口, 如果有值, 把prev值带返回, 并且被替换, 没有就新加
func (a *AvlTree[K, V]) SetWithPrev(k K, v V) (prev V, replaced bool) {
	link := &a.root.node
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
			link = &parent.right
		} else {
			link = &parent.left
		}
	}

	node.link(parent, link)
	a.root.postInsert(node)
	a.length++
	return
}

func (r *root[K, V]) rebalance(node *node[K, V]) {

	for ; node != nil; node = node.parent {
		lh := node.leftHeight()
		lr := node.rightHeight()
		height := cmp.Max(lh, lr) + 1

		diff := lh - lr
		if node.height != height {
			node.height = height
		} else if diff >= -1 && diff <= 1 {
			break
		}

		if diff <= -2 {
			node = r.fixLeft(node)
		} else if diff >= 2 {
			node = r.fixRight(node)
		}
	}
}

func (a *AvlTree[K, V]) Delete(k K) {
	a.Remove(k)
}

func (a *AvlTree[K, V]) Remove(k K) *AvlTree[K, V] {
	n := a.root.node
	for n != nil {
		if n.key == k {
			goto found
		}

		if k > n.key {
			n = n.right
		} else {
			n = n.left
		}
	}

	return a

found:
	var child, parent *node[K, V]
	if n.left != nil && n.right != nil {
		old := n
		n = n.right
		for left := n; left != nil; left = left.left {
			n = left
		}
		// 待会儿old被删除时, 使用n贴到old原来的位置

		child = n.left
		parent = n.parent
		if child != nil {
			// child 这条线不再n 节点
			child.parent = parent
		}
		// TODO 写注释
		a.root.childReplace(n, child, parent)

		if n.parent == old {
			parent = n
		}

		// 把n节点贴到原来old的位置
		n.left = old.left
		n.right = old.right
		n.parent = old.parent
		n.height = old.height

		a.root.childReplace(old, n, old.parent)
		old.left.parent = n

		if old.right != nil {
			old.right.parent = n
		}
	} else {
		if n.left == nil {
			child = n.right
		} else {
			child = n.left
		}
		parent = n.parent
		a.root.childReplace(n, child, parent)
		if child != nil {
			child.parent = parent
		}
	}

	if parent != nil {
		a.root.rebalance(parent)
	}
	return a
}

func (n *node[K, V]) rangeInner(callback func(k K, v V) bool) bool {

	if n == nil {
		return true
	}

	if n.left != nil {
		if !n.left.rangeInner(callback) {
			return false
		}
	}

	if !callback(n.key, n.val) {
		return false
	}

	if n.right != nil {
		if !n.right.rangeInner(callback) {
			return false
		}
	}
	return true
}

func (n *node[K, V]) rangePrevInner(callback func(k K, v V) bool) bool {

	if n == nil {
		return true
	}

	if n.right != nil {
		if !n.right.rangePrevInner(callback) {
			return false
		}
	}

	if !callback(n.key, n.val) {
		return false
	}

	if n.left != nil {
		if !n.left.rangePrevInner(callback) {
			return false
		}
	}

	return true
}

// 遍历avl tree
func (a *AvlTree[K, V]) Range(callback func(k K, v V) bool) {
	// 遍历
	if a.root.node == nil {
		return
	}

	a.root.node.rangeInner(callback)
	return
}

// 遍历avl tree
func (a *AvlTree[K, V]) RangePrev(callback func(k K, v V) bool) {
	// 遍历
	if a.root.node == nil {
		return
	}

	a.root.node.rangePrevInner(callback)
	return
}

func (a *AvlTree[K, V]) TopMax(limit int, callback func(k K, v V) bool) {
	a.RangePrev(func(k K, v V) bool {

		if limit <= 0 {
			return false
		}

		if !callback(k, v) {
			return false
		}

		limit--
		return true
	})
}

func (a *AvlTree[K, V]) TopMin(limit int, callback func(k K, v V) bool) {

	a.Range(func(k K, v V) bool {

		if limit <= 0 {
			return false
		}

		if !callback(k, v) {
			return false
		}

		limit--
		return true
	})
}

func (a *AvlTree[K, V]) Len() int {
	return a.length
}

func (a *AvlTree[K, V]) Draw() {
	if a.root.node == nil {
		return
	}

	a.root.node.draw(a.root.node)
}

// 画出avl tree
// 使用层序遍历的姿势
func (n *node[K, V]) draw(root *node[K, V]) {
	if root == nil {
		return
	}

	q := vec.New(root)
	for height := 0; q.Len() > 0; height++ {
		tmp := q.ToSlice()
		q = vec.New[*node[K, V]]()

		fmt.Printf("height:%d ", height)
		for _, node := range tmp {
			fmt.Printf("%v ", node.pair)

			if node.left != nil {

				q.Push(node.left)
			}

			if node.right != nil {
				q.Push(node.right)
			}

		}
		fmt.Printf("\n")

	}
}
