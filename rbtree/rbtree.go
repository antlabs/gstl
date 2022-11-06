package rbtree

// apache 2.0 antlabs
// 参考资料
// https://github.com/torvalds/linux/blob/master/lib/rbtree.c
import (
	"errors"

	"github.com/antlabs/gstl/api"
	"golang.org/x/exp/constraints"
)

// 红黑树5条重要性质
// 1. 节点为红色或者黑色
// 2. 根节点是黑色(黑根)
// 3. 所有叶子(空节点)均为黑色
// 4. 每个红色节点的两个子节点均为黑色(红父黑子)
// 5. 从根到叶的每个路径包含相同数量的黑色节点(黑高相同)

var _ api.SortedMap[int, int] = (*RBTree[int, int])(nil)

var ErrNotFound = errors.New("rbtree: not found value")

type color int8

const (
	RED   color = 1
	BLACK color = 2
)

// 元素
type pair[K constraints.Ordered, V any] struct {
	val V
	key K
}

type parentColor[K constraints.Ordered, V any] struct {
	parent *node[K, V]
	color  color
}

type node[K constraints.Ordered, V any] struct {
	left  *node[K, V]
	right *node[K, V]
	pair[K, V]
	parentColor[K, V]
}

func (n *node[K, V]) setParent(parent *node[K, V]) {
	n.parent = parent
}

func (n *node[K, V]) link(parent *node[K, V], link **node[K, V]) {
	n.parent = parent
	n.color = RED
	*link = n
}

type root[K constraints.Ordered, V any] struct {
	node *node[K, V]
}

func (r *root[K, V]) rotateLeft(n *node[K, V]) {

	right := n.right

	n.right = right.left
	if right.left != nil {
		right.left.parent = n
	}
	right.left = n

	right.parent = n.parent

	if n.parent != nil {

		if n == n.parent.left {
			n.parent.left = right
		} else {
			n.parent.right = right
		}
	} else {
		r.node = right
	}
	n.parent = right
}

func (r *root[K, V]) rotateRight(n *node[K, V]) {
	left := n.left
	n.left = left.right
	if left.right != nil {
		left.right.parent = n
	}
	left.right = n

	left.parent = n.parent
	if n.parent != nil {
		if n == n.parent.right {
			n.parent.right = left
		} else {
			n.parent.left = left
		}
	} else {
		r.node = left
	}
	n.parent = left
}

func (r *root[K, V]) changeChild(old, new, parent *node[K, V]) {
	if parent != nil {
		if parent.left == old {
			parent.left = new
		} else {
			parent.right = new
		}
	} else {
		r.node = new
	}

}

func (r *root[K, V]) insert(n *node[K, V]) {

	var parent, gparent *node[K, V]

	for parent = n.parent; parent != nil && parent.color == RED; parent = n.parent {

		gparent = parent.parent
		if parent == gparent.left {

			uncle := gparent.right
			if uncle != nil && uncle.color == RED {
				uncle.color = BLACK
				parent.color = BLACK
				gparent.color = RED
				n = gparent
				continue
			}

			if parent.right == n {
				r.rotateLeft(parent)
				parent, n = n, parent
			}

			parent.color = BLACK
			gparent.color = RED
			r.rotateRight(gparent)
		} else {
			uncle := gparent.left
			if uncle != nil && uncle.color == RED {
				uncle.color = BLACK
				parent.color = BLACK
				gparent.color = RED
				n = gparent
				continue
			}

			if parent.left == n {
				r.rotateRight(parent)
				parent, n = n, parent
			}
			parent.color = BLACK
			gparent.color = RED
			r.rotateLeft(gparent)
		}
	}
	r.node.color = BLACK //黑根
}

// 红黑树
type RBTree[K constraints.Ordered, V any] struct {
	length int
	root   root[K, V]
}

// 初始化函数
func New[K constraints.Ordered, V any]() *RBTree[K, V] {
	return &RBTree[K, V]{}
}

// 第一个节点
func (r *RBTree[K, V]) First() (v V, err error) {
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
func (r *RBTree[K, V]) Last() (v V, err error) {
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

func (r *RBTree[K, V]) Set(k K, v V) {
	_, _ = r.SetWithPrev(k, v)
}

// 设置
func (r *RBTree[K, V]) SetWithPrev(k K, v V) (prev V, replaced bool) {
	link := &r.root.node
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
	r.root.insert(node)
	r.length++
	return
}

// Get
func (r *RBTree[K, V]) Get(k K) (v V) {
	v, _ = r.GetWithErr(k)
	return
}

// 从rbtree 找到需要的值
func (r *RBTree[K, V]) GetWithErr(k K) (v V, err error) {
	n := r.root.node
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

// 删除
func (r *root[K, V]) erase(n *node[K, V]) {

	var child, parent *node[K, V]
	var color color
	if n.left == nil {
		child = n.right
	} else if n.right == nil {
		child = n.left
	} else {
		old := n
		n = n.right
		left := n.left
		for ; left != nil; left = n.left {
		}
		child = n.right
		parent = n.parent
		color = n.color

		if child != nil {
			child.parent = parent
		}

		if parent != nil {
			if parent.left == n {
				parent.left = child
			} else {
				parent.right = child
			}
		} else {
			r.node = child
		}

		if n.parent == old {
			parent = n
		}

		n.parent = old.parent
		n.color = old.color
		n.right = old.right
		n.left = old.left

		if old.parent != nil {
			if old.parent.left == old {
				old.parent.left = n
			} else {
				old.parent.right = n
			}
		} else {
			r.node = n
		}
		old.left.parent = n
		if old.right != nil {
			old.right.parent = n
		}
		goto color
	}
	parent = n.parent
	color = n.color

	if child != nil {
		child.parent = parent
	}
	if parent != nil {
		if parent.left == n {
			parent.left = child
		} else {
			parent.right = child
		}
	} else {
		r.node = child
	}

color:
	if color == BLACK {
		r.eraseColor(child, parent)
	}
}

func (r *root[K, V]) eraseColor(n *node[K, V], parent *node[K, V]) {

	var other *node[K, V]
	for (n == nil || n.color == BLACK) && n != r.node {
		if parent.left == n {
			other = parent.right
			if other.color == RED {
				other.color = BLACK
				parent.color = RED
				r.rotateLeft(parent)
				other = parent.right
			}
			if (other.left == nil || other.left.color == BLACK) && (other.right == nil || other.right.color == BLACK) {

				other.color = RED
				n = parent
				parent = n.parent
			} else {

				if other.right == nil || other.right.color == BLACK {

					oleft := other.left
					if oleft != nil {
						oleft.color = BLACK
					}
					other.color = RED
					r.rotateRight(other)
					other = parent.right
				}
				other.color = parent.color
				parent.color = BLACK
				if other.right != nil {
					other.right.color = BLACK
				}
				r.rotateLeft(parent)
				n = r.node
				break
			}
		} else {

			other = parent.left
			if other.color == RED {

				other.color = BLACK
				parent.color = RED
				r.rotateRight(parent)
				other = parent.left
			}

			if (other.left == nil || other.left.color == BLACK) && (other.right == nil || other.right.color == BLACK) {
				other.color = RED
				n = parent
				parent = n.parent
			} else {
				if other.left == nil || other.left.color == BLACK {

					oright := other.right
					if oright != nil {
						oright.color = BLACK
					}
					other.color = RED
					r.rotateLeft(other)
					other = parent.left
				}
				other.color = parent.color
				parent.color = BLACK
				if other.left != nil {
					other.left.color = BLACK
				}
				r.rotateRight(parent)

				n = r.node
				break
			}

		}
	}
	if n != nil {
		n.color = BLACK
	}

}

func (r *RBTree[K, V]) Delete(k K) {
	n := r.root.node
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
	return

found:
	r.root.erase(n)
	return
}

func (r *RBTree[K, V]) Len() int {
	return r.length
}

func (r *RBTree[K, V]) TopMin(limit int, callback func(k K, v V) bool) {

	r.Range(func(k K, v V) bool {

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

// 遍历rbtree
func (r *RBTree[K, V]) RangePrev(callback func(k K, v V) bool) {
	// 遍历
	if r.root.node == nil {
		return
	}

	r.root.node.rangePrevInner(callback)
	return
}

func (r *RBTree[K, V]) TopMax(limit int, callback func(k K, v V) bool) {

	r.RangePrev(func(k K, v V) bool {

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

// 遍历rbtree
func (a *RBTree[K, V]) Range(callback func(k K, v V) bool) {
	// 遍历
	if a.root.node == nil {
		return
	}

	a.root.node.rangeInner(callback)
	return
}
