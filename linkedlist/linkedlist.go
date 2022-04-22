package linkedlist

// 参考文档如下
// https://cs.opensource.google/go/go/+/go1.18.1:src/container/list/list.go
// https://github.com/torvalds/linux/blob/master/tools/include/linux/list.h
// https://redis.io/commands/?group=list
// LTRIM -->done
// LINDEX -->done
// LINSERT -->done
// LLEN -->done
// LPOP -->done
// LPOS #暂不实现
// LPUSH LPUSHX -->done
// LRANGE -->done
// LREM -->done
// LSET -->done
// RPOP -->done
// RPUSH RPUSHX -->done

import (
	"errors"
	"github.com/guonaihong/gstl/cmp"
)

var ErrListElemEmpty = errors.New("list is empty")
var ErrNotFound = errors.New("element not found")

type LinkedList[T any] struct {
	root   Node[T]
	length int
}

// 每个Node节点, 包含前向和后向两个指针和数据域
type Node[T any] struct {
	next    *Node[T]
	prev    *Node[T]
	Element T
}

// 返回一个双向循环链表
func New[T any]() *LinkedList[T] {
	return new(LinkedList[T]).Init()
}

// 指向自己, 组成一个环
func (l *LinkedList[T]) Init() *LinkedList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.length = 0
	return l
}

// 返回长度
func (l *LinkedList[T]) Len() int {
	return l.length
}

// 延迟初始化
func (l *LinkedList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// at e at.next
// at <- e
//       e -> at.next
// at -> e
//       e <- at.next
func (l *LinkedList[T]) insert(at, e *Node[T]) {
	e.prev = at
	e.next = at.next
	e.next.prev = e
	at.next = e
	l.length++
}

func (l *LinkedList[T]) OtherMoveToBackList(other *LinkedList[T]) *LinkedList[T] {
	l.lazyInit()
	if l == other || other.length == 0 {
		return l
	}

	l.length += other.length

	tail := l.root.prev
	otherHead := other.root.next
	otherTail := other.root.prev

	// 第一个链表的尾巴next指针接上第二个链表的头
	tail.next = otherHead
	// 第二个链表的头prev指针接上第一个链表的尾部
	otherHead.prev = tail

	otherTail.next = &l.root
	l.root.prev = otherTail

	other.Init()
	return l
}

func (l *LinkedList[T]) OtherMoveToFrontList(other *LinkedList[T]) *LinkedList[T] {
	l.lazyInit()
	if l == other || other.length == 0 {
		return l
	}

	l.length += other.length

	head := l.root.next
	otherHead := other.root.next
	otherTail := other.root.prev

	otherTail.next = head
	head.prev = otherTail

	l.root.next = otherHead
	otherHead.prev = &l.root

	other.Init()
	return l
}

// 类似redis lpop命令
// O(n)
func (l LinkedList[T]) LPop(count int) []T {
	if count <= 0 {
		return nil
	}

	count = cmp.Min(count, l.length)
	all := make([]T, count)
	i := 0

	l.RangeSafe(func(n *Node[T]) bool {
		if i == count {
			return true
		}

		all[i] = n.Element
		l.remove(n)
		i++
		return false
	})
	return all
}

// 类似redis rpop命令
// O(n)
func (l *LinkedList[T]) RPop(count int) []T {
	if count <= 0 {
		return nil
	}

	count = cmp.Min(count, l.length)
	all := make([]T, count)
	l.RangePrevSafe(func(n *Node[T]) bool {
		if count <= 0 {
			return true
		}

		all[count-1] = n.Element
		l.remove(n)
		count--
		return false
	})
	return all
}

// 从后向前遍历
func (l *LinkedList[T]) RangePrevSafe(callback func(n *Node[T]) bool) {

	pos := l.root.prev
	n := pos.prev

	for pos != &l.root {
		if callback(pos) {
			break
		}
		pos = n
		n = pos.prev
	}
}

// 从前向后遍历
// callback 返回truek就退出遍历
func (l *LinkedList[T]) RangeSafe(callback func(n *Node[T]) (exit bool)) {

	pos := l.root.next
	n := pos.next

	for pos != &l.root {
		if callback(pos) {
			break
		}

		pos = n
		n = pos.next
	}
}

// 类似redis lpush命令
// PushFront的同义词
func (l *LinkedList[T]) LPush(elems ...T) *LinkedList[T] {
	l.PushFront(elems...)
	return l
}

// PushFrontList在列表l前面插入一个新的列表other的副本
func (l *LinkedList[T]) PushFrontList(other *LinkedList[T]) *LinkedList[T] {
	other.RangePrevSafe(func(n *Node[T]) bool {
		l.insert(&l.root, &Node[T]{Element: n.Element})
		return false
	})

	return l

}

// 往头位置插入
func (l *LinkedList[T]) PushFront(elems ...T) *LinkedList[T] {
	l.lazyInit()
	for _, e := range elems {
		l.insert(&l.root, &Node[T]{Element: e})
	}
	return l
}

// RPush是PushBack的同义词, 类似redis的RPush命令
func (l *LinkedList[T]) RPush(elems ...T) *LinkedList[T] {
	l.PushBack(elems...)
	return l
}

// PushBackList往尾部的位置插入一个新的列表other的副本
func (l *LinkedList[T]) PushBackList(other *LinkedList[T]) *LinkedList[T] {
	l.lazyInit()

	other.RangeSafe(func(n *Node[T]) bool {
		l.insert(l.root.prev, &Node[T]{Element: n.Element})
		return false
	})
	return l
}

// 往尾部的位置插入
func (l *LinkedList[T]) PushBack(elems ...T) *LinkedList[T] {
	l.lazyInit()
	for _, e := range elems {
		l.insert(l.root.prev, &Node[T]{Element: e})
	}
	return l
}

// 返回第1个元素
func (l *LinkedList[T]) First() (e T, err error) {
	if l.length == 0 {
		err = ErrListElemEmpty
		return
	}
	return l.root.next.Element, nil
}

// 返回最后1个元素
func (l *LinkedList[T]) Last() (e T, err error) {
	if l.length == 0 {
		err = ErrListElemEmpty
		return
	}
	return l.root.prev.Element, nil
}

// 链表是否为空
func (l *LinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

// 清空链表 O(n)
func (l *LinkedList[T]) Clear() *LinkedList[T] {

	l.RangeSafe(func(n *Node[T]) bool {
		l.remove(n)
		return false
	})

	return l
}

// 类似于redis linsert after 命令
func (l *LinkedList[T]) InsertAfter(value T, equal func(value T) bool) *LinkedList[T] {
	l.RangeSafe(func(n *Node[T]) bool {
		if equal(n.Element) {
			l.insert(n, &Node[T]{Element: value})
			return true
		}
		return false
	})
	return l
}

//  类似于redis linsert before 命令
func (l *LinkedList[T]) InsertBefore(value T, equal func(value T) bool) *LinkedList[T] {
	l.RangeSafe(func(n *Node[T]) bool {
		if equal(n.Element) {
			l.insert(n.prev, &Node[T]{Element: value})
			return true
		}
		return false
	})

	return l
}

// 查找是否包含这个value
func (l *LinkedList[T]) ContainsFunc(cb func(value T) bool) bool {
	for pos := l.root.next; pos != &l.root; pos = pos.next {
		if cb(pos.Element) {
			return true
		}
	}
	return false
}

// 通过索引查找是否包含这个value
// Get是Index的同义词
func (l *LinkedList[T]) Get(idx int) (e T, err error) {
	return l.Index(idx)
}

// 删除这个元素
// n.prev n n.next
// n.prev --> n.next
// n.prev <-- n.next
func (l *LinkedList[T]) remove(n *Node[T]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	l.length--
}

// 类似redis lrem命令
// count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
// count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
// count = 0 : 移除表中所有与 VALUE 相等的值。

// 返回值
// 被删除元素个数
func (l *LinkedList[T]) RemFunc(count int, cb func(value T) bool) (ndel int) {
	var (
		pos *Node[T]
		n   *Node[T]
		i   int
	)

	if count >= 0 {
		pos = l.root.next
		n = pos.next
		for pos != &l.root {
			if count == 0 || i <= count {
				if cb(pos.Element) {
					l.remove(pos)
					ndel++
				}
				i++
			}
			pos = n
			n = pos.next
		}

		return
	}

	count = -count
	pos = l.root.prev
	n = pos.prev
	for pos != &l.root {
		if count == 0 || i <= count {
			if cb(pos.Element) {
				l.remove(pos)
				ndel++
			}
			i++
		}
		pos = n
		n = pos.prev
	}

	return
}

// Index 通过索引查找是否包含这个value
// 和redis lindex命令类似,
// idx >= 0 形为和Get(idx int) 一样, 获取指向索引的元素
// idx < 0 获取倒数第几个元素
// O(min(index, length - index))
func (l *LinkedList[T]) Index(idx int) (e T, err error) {
	var n *Node[T]
	n, err = l.indexInner(idx)
	if err != nil {
		return
	}
	return n.Element, nil
}

// 类型redis lset命令
// index >= 0 正着数
// index < 0 倒着数
// On(min(index, length - index))
func (l *LinkedList[T]) Set(index int, value T) *LinkedList[T] {
	n, err := l.indexInner(index)
	if err != nil {
		return l
	}
	n.Element = value
	return l
}

func (l *LinkedList[T]) indexInner(idx int) (*Node[T], error) {
	idx, front := l.index(idx)

	if front {
		for pos, i := l.root.next, 0; pos != &l.root; pos, i = pos.next, i+1 {
			if i == idx {
				return pos, nil
			}
		}
	} else {
		for pos, i := l.root.prev, idx; pos != &l.root; pos, i = pos.prev, i-1 {
			if i == 0 {
				return pos, nil
			}
		}

	}

	return nil, ErrNotFound

}

// 计算索引
func (l *LinkedList[T]) index(idx int) (newIdx int, front bool) {
	length := l.length

	// 转正索引
	if idx < 0 {
		idx = idx + length
	}

	// 如果倒序遍历元素更少
	other := length - idx - 1
	if idx > other {
		return other, false
	}

	// 如果正序遍历元素更少
	return idx, true
}

// 删除指定索引的元素, 效率 min(O(index), O(len - index))
func (l *LinkedList[T]) Remove(index int) *LinkedList[T] {
	l.removeInner(index)
	return l
}

// list 转成slice , 效率O(n)
func (l *LinkedList[T]) ToSlice() []T {
	if l.length == 0 {
		return nil
	}

	rv := make([]T, 0, l.length)
	for pos := l.root.next; pos != &l.root; pos = pos.next {
		rv = append(rv, pos.Element)
	}

	return rv
}

func (l *LinkedList[T]) removeInner(index int) {
	var (
		pos *Node[T]
		n   *Node[T]
		i   int
	)

	idx, front := l.index(index)

	if front {

		pos = l.root.next
		n = pos.next
		for pos != &l.root {
			if i == idx {
				l.remove(pos)
				return
			}

			pos = n
			n = pos.next
			i++
		}
	}

	i = idx
	pos = l.root.prev
	n = pos.prev
	for pos != &l.root {
		if i == 0 {
			l.remove(pos)
			return
		}
		pos = n
		n = pos.prev
		i--
	}

}

// range 类似redis lrange命令
func (l *LinkedList[T]) Range(callback func(value T), startAndEnd ...int) {

	start := 0
	end := 0

	if len(startAndEnd) > 0 {
		start = startAndEnd[0]
	}

	if len(startAndEnd) > 1 {
		end = startAndEnd[1]
	}

	i := 0
	l.rangeStartEndSafe(start, end, func(start, end int, n *Node[T]) bool {
		if len(startAndEnd) != 0 {
			if i >= start && i <= end {
				callback(n.Element)
			}

			if i > end {
				return true
			}
			i++
			return false
		}

		i++
		callback(n.Element)
		return false
	})
}

func (l *LinkedList[T]) rangeStartEndSafe(start, end int, callback func(start, end int, n *Node[T]) (exit bool)) {

	if start < 0 {
		start += l.length
		if start < 0 {
			start = 0
		}
	}

	if end < 0 {
		end += l.length
	}

	if start > end || start >= l.length {
		return
	}
	var pos *Node[T]
	var n *Node[T]

	pos = l.root.next
	n = pos.next

	for pos != &l.root {
		if callback(start, end, pos) {
			break
		}
		pos = n
		n = pos.next
	}
}

// 类似于redis ltrim命令, 对列表进行裁剪
func (l *LinkedList[T]) Trim(start, end int) *LinkedList[T] {

	i := 0
	l.rangeStartEndSafe(start, end, func(start, end int, n *Node[T]) bool {
		if i < start || i > end {
			l.remove(n)
		}
		i++
		return false
	})
	return l
}
