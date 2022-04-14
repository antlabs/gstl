package linkedlist

import (
	"errors"
)

var ErrListElemEmpty = errors.New("list is empty")

// https://cs.opensource.google/go/go/+/go1.18.1:src/container/list/list.go
// https://github.com/torvalds/linux/blob/master/tools/include/linux/list.h
type LinkedList[T any] struct {
	root   Node[T]
	length int
}

// 每个Node节点, 包含前向和后向两个指针和数据域
type Node[T any] struct {
	next    *Node[T]
	prev    *Node[T]
	element T
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
	at.next = e
	at.next.prev = e
	l.length++
}

func (l *LinkedList[T]) Append(other LinkedList[T]) {
	l.lazyInit()

}

// 往头位置插入
func (l *LinkedList[T]) PushFront(elems ...T) {
	l.lazyInit()
	for _, e := range elems {
		l.insert(&l.root, &Node[T]{element: e})
	}
}

// 往尾部的位置插入
func (l *LinkedList[T]) PushBack(elems ...T) {
	l.lazyInit()
	for _, e := range elems {
		l.insert(l.root.prev, &Node[T]{element: e})
	}
}

// 返回第1个元素
func (l *LinkedList[T]) First() (e T, err error) {
	if l.length == 0 {
		err = ErrListElemEmpty
		return
	}
	return l.root.next.element, nil
}

// 返回最后1个元素
func (l *LinkedList[T]) Last() (e T, err error) {
	if l.length == 0 {
		err = ErrListElemEmpty
		return
	}
	return l.root.prev.element, nil
}

// 链表是否为空
func (l *LinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

func (l *LinkedList[T]) Clear() {

}

func (l *LinkedList[T]) Contains() bool {
	return false
}

func (l *LinkedList[T]) Remove(index int) {

}
