package rbtree

// apache 2.0 guonaihong
// 参考资料
// https://github.com/torvalds/linux/blob/master/lib/rbtree.c
import (
	"errors"

	"golang.org/x/exp/constraints"
)

// 红黑树5条重要性质
// 1. 节点为红色或者黑色
// 2. 根节点是黑色(黑根)
// 3. 所有叶子(空节点)均为黑色
// 4. 每个红色节点的两个子节点均为黑色(红父黑子)
// 5. 从根到叶的每个路径包含相同数量的黑色节点(黑高相同)

var ErrNotFound = errors.New("rbtree: not found value")

const (
	RED   = 1
	BLACK = 2
)

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
	// red/black 可以用一个变量实现, 为了代码简单清晰, 这里不考虑一个字节的收益, 放弃使用位操作
	red   bool
	black bool
}

func (n *node[K, V]) setParentRed(parent *node[K, V]) {
	parent.red = true
	parent.black = false
	n.parent = parent
}

func (n *node[K, V]) setParentBlack(parent *node[K, V]) {
	parent.red = false
	parent.black = true
	n.parent = parent
}

func (n *node[K, V]) link(parent *node[K, V], link **node[K, V]) {
	n.parent = parent
	n.red = true
	n.black = false
	*link = n
}

type root[K constraints.Ordered, V any] struct {
	node *node[K, V]
}

func (r *root[K, V]) rotateSetParents(old, new *node[K, V], color int) {
	parent := old.parent
	new.parent = old.parent
	if color == RED {
		old.setParentRed(parent)
	} else {
		old.setParentBlack(parent)
	}
	r.changeChild(old, new, parent)
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

	parent := n.parent
	var gparent, tmp *node[K, V]

	for {

		if parent == nil {
			n.setParentBlack(parent)
			break
		}

		if parent.black {
			break
		}

		gparent = parent.parent
		tmp = gparent.right
		if parent != tmp {
			if tmp != nil && tmp.red {
				/*
				 * Case 1 - node's uncle is red (color flips).
				 *
				 *       G            g
				 *      / \          / \
				 *     p   u  -->   P   U
				 *    /            /
				 *   n            n
				 *
				 * However, since g's parent might be red, and
				 * 4) does not allow this, we need to recurse
				 * at g.
				 */
				tmp.setParentBlack(gparent)
				parent.setParentRed(gparent)
				n = gparent
				parent = n.parent
				n.setParentRed(parent)
				continue
			}

			tmp = parent.right
			if n == tmp {
				/*
				 * Case 2 - node's uncle is black and node is
				 * the parent's right child (left rotate at parent).
				 *
				 *      G             G
				 *     / \           / \
				 *    p   U  -->    n   U
				 *     \           /
				 *      n         p
				 *
				 * This still leaves us in violation of 4), the
				 * continuation into Case 3 will fix that.
				 */
				tmp = n.left
				parent.right = tmp
				n.left = parent
				if tmp != nil {
					tmp.setParentBlack(parent)
				}
				parent.setParentRed(n)
				parent = n
				tmp = n.right
			}

			/*
			 * Case 3 - node's uncle is black and node is
			 * the parent's left child (right rotate at gparent).
			 *
			 *        G           P
			 *       / \         / \
			 *      p   U  -->  n   g
			 *     /                 \
			 *    n                   U
			 */
			gparent.left = tmp
			parent.right = gparent
			if tmp != nil {
				tmp.setParentBlack(gparent)
			}
			break
		} else {
			tmp = gparent.left
			if tmp != nil && tmp.red {
				tmp.setParentBlack(gparent)
				parent.setParentRed(gparent)
				n = gparent
				parent = n.parent
				n.setParentRed(parent)
				continue
			}
			tmp = parent.left
			if n == tmp {
				tmp = n.right
				parent.left = tmp
				parent.right = parent
				if tmp != nil {
					tmp.setParentBlack(parent)
				}
				parent.setParentRed(n)
				parent = n
				tmp = n.left
			}

			/* Case 3 - left rotate at gparent */
			gparent.right = tmp
			parent.left = gparent
			if tmp != nil {
				tmp.setParentBlack(gparent)
			}
			r.rotateSetParents(gparent, parent, RED)
			break
		}
	}
}

// 红黑树
type RBTree[K constraints.Ordered, V any] struct {
	length int
	root   root[K, V]
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
			*link = parent.right
		} else {
			*link = parent.left
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
