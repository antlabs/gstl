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

type parentColor[K constraints.Ordered, V any] struct {
	parent *node[K, V]
	// red/black 可以用一个变量实现, 为了代码简单清晰, 这里不考虑一个字节的收益, 放弃使用位操作
	red   bool
	black bool
}

type node[K constraints.Ordered, V any] struct {
	left  *node[K, V]
	right *node[K, V]
	pair[K, V]
	parentColor[K, V]
}

func (n *node[K, V]) setBlack() {
	n.black = true
	n.red = false
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

func (n *node[K, V]) setParent(parent *node[K, V]) {
	n.parent = parent
}

func (n *node[K, V]) setParentWithColor(parent *node[K, V], red bool, black bool) {
	n.parent = parent
	n.red = red
	n.black = black
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

func (r *root[K, V]) eraseAugmented(n *node[K, V]) *node[K, V] {
	child := n.right
	tmp := n.left

	var parent, rebalance *node[K, V]

	var pc parentColor[K, V]

	if tmp == nil {
		/*
		 * Case 1: node to erase has no more than 1 child (easy!)
		 *
		 * Note that if there is one child it must be red due to 5)
		 * and node must be black due to 4). We adjust colors locally
		 * so as to bypass __rb_erase_color() later on.
		 */

		pc = n.parentColor
		parent = n.parent

		r.changeChild(n, child, parent)
		if child != nil {
			child.parentColor = pc
		} else if pc.black {
			rebalance = parent
		}

		tmp = parent
	} else if child == nil {
		/* Still case 1, but this time the child is node->rb_left */
		tmp.parentColor = n.parentColor
		parent = n.parent
		r.changeChild(n, tmp, parent)
		tmp = parent
	} else {
		successor := child
		var child2 *node[K, V]

		tmp = child.left
		if tmp == nil {
			/*
			 * Case 2: node's successor is its right child
			 *
			 *    (n)          (s)
			 *    / \          / \
			 *  (x) (s)  ->  (x) (c)
			 *        \
			 *        (c)
			 */

			parent = successor
			child2 = successor.right
		} else {
			/*
			 * Case 3: node's successor is leftmost under
			 * node's right child subtree
			 *
			 *    (n)          (s)
			 *    / \          / \
			 *  (x) (y)  ->  (x) (y)
			 *      /            /
			 *    (p)          (p)
			 *    /            /
			 *  (s)          (c)
			 *    \
			 *    (c)
			 */
			for {
				parent = successor
				successor = tmp
				tmp = tmp.right
				if tmp == nil {
					break
				}
			}

			child2 = successor.right
			parent.left = child2
			successor.right = child
			child.setParent(successor)
		}

		tmp = n.left
		successor.left = tmp
		tmp.setParent(successor)

		pc = n.parentColor
		tmp = n.parent
		r.changeChild(n, successor, tmp)

		if child2 != nil {
			successor.parentColor = pc
			child2.setParentBlack(parent)
			rebalance = nil
		} else {
			pc2 := successor.parentColor
			successor.parentColor = pc
			if pc2.black {
				rebalance = parent
			}
		}
		tmp = successor
	}
	return rebalance
}

func (r *root[K, V]) eraseColor(parent *node[K, V]) {
	var n, sibling, tmp1, tmp2 *node[K, V]

	for {
		/*
		 * Loop invariants:
		 * - node is black (or NULL on first iteration)
		 * - node is not the root (parent is not NULL)
		 * - All leaf paths going through parent and node have a
		 *   black node count that is 1 lower than other leaf paths.
		 */
		sibling = parent.right
		if n != sibling { /* node == parent->rb_left */
			if sibling.red {
				/*
				 * Case 1 - left rotate at parent
				 *
				 *     P               S
				 *    / \             / \
				 *   N   s    -->    p   Sr
				 *      / \         / \
				 *     Sl  Sr      N   Sl
				 */
				tmp1 = sibling.left
				parent.right = tmp1
				sibling.left = parent
				tmp1.setParentBlack(parent)
				r.rotateSetParents(parent, sibling, RED)
				sibling = tmp1
			}
			tmp1 = sibling.right

			if tmp1 == nil || tmp1.black {
				tmp2 = sibling.left
				if tmp2 == nil || tmp2.black {
					/*
					 * Case 2 - sibling color flip
					 * (p could be either color here)
					 *
					 *    (p)           (p)
					 *    / \           / \
					 *   N   S    -->  N   s
					 *      / \           / \
					 *     Sl  Sr        Sl  Sr
					 *
					 * This leaves us violating 5) which
					 * can be fixed by flipping p to black
					 * if it was red, or by recursing at p.
					 * p is red when coming from Case 1.
					 */
					sibling.setParentRed(parent)
					if parent.red {
						parent.setBlack()
					} else {
						n = parent
						parent = n.parent
						if parent != nil {
							continue
						}
					}
					break
				}
				/*
				 * Case 3 - right rotate at sibling
				 * (p could be either color here)
				 *
				 *   (p)           (p)
				 *   / \           / \
				 *  N   S    -->  N   sl
				 *     / \             \
				 *    sl  Sr            S
				 *                       \
				 *                        Sr
				 *
				 * Note: p might be red, and then both
				 * p and sl are red after rotation(which
				 * breaks property 4). This is fixed in
				 * Case 4 (in __rb_rotate_set_parents()
				 *         which set sl the color of p
				 *         and set p RB_BLACK)
				 *
				 *   (p)            (sl)
				 *   / \            /  \
				 *  N   sl   -->   P    S
				 *       \        /      \
				 *        S      N        Sr
				 *         \
				 *          Sr
				 */
				tmp1 = tmp2.right
				sibling.left = tmp1
				tmp2.right = sibling
				parent.right = tmp2
				if tmp1 != nil {
					tmp1.setParentBlack(sibling)
				}
				tmp1 = sibling
				sibling = tmp2
			}
			/*
			 * Case 4 - left rotate at parent + color flips
			 * (p and sl could be either color here.
			 *  After rotation, p becomes black, s acquires
			 *  p's color, and sl keeps its color)
			 *
			 *      (p)             (s)
			 *      / \             / \
			 *     N   S     -->   P   Sr
			 *        / \         / \
			 *      (sl) sr      N  (sl)
			 */
			tmp2 = sibling.left
			parent.right = tmp2
			sibling.left = parent
			tmp1.setParentRed(sibling)
			if tmp2 != nil {
				tmp2.setParent(parent)
			}
			r.rotateSetParents(parent, sibling, RED)
			break
		} else {
			sibling = parent.left
			if sibling.red {
				tmp1 = sibling.right
				parent.left = tmp1
				sibling.right = parent
				tmp1.setParentBlack(parent)
				r.rotateSetParents(parent, sibling, RED)
				sibling = tmp1
			}

			tmp1 = sibling.left
			if tmp1 == nil || tmp1.black {
				tmp2 = sibling.right
				if tmp2 == nil || tmp2.black {
					sibling.setParentRed(parent)
					if parent.red {
						parent.setBlack()
					} else {
						n = parent
						parent = n.parent
						if parent != nil {
							continue
						}
					}
					break
				}
				/* Case 3 - left rotate at sibling */
				tmp1 = tmp2.left
				sibling.right = tmp1
				tmp2.left = sibling
				parent.left = tmp2

				if tmp1 != nil {
					tmp1.setParentBlack(sibling)
				}

				tmp1 = sibling
				sibling = tmp2
			}
			/* Case 4 - right rotate at parent + color flips */
			tmp2 = sibling.right
			parent.left = tmp2
			sibling.right = parent
			tmp1.setParentBlack(sibling)
			if tmp2 != nil {
				tmp2.setParent(parent)
			}
			r.rotateSetParents(parent, sibling, BLACK)
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
	rebalance := r.root.eraseAugmented(n)
	if rebalance != nil {
		r.root.eraseColor(n)
	}
	return
}
