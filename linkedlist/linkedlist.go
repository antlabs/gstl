package linkedlist

// https://cs.opensource.google/go/go/+/go1.18.1:src/container/list/list.go
// https://github.com/torvalds/linux/blob/master/tools/include/linux/list.h
type LinkedList[T any] struct {
	root   Node[T]
	length int
}

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
func (l *LinkedList[T]) PushFront(e T) {
	l.lazyInit()
	l.insert(&l.root, &Node[T]{element: e})
}

// 往尾部的位置插入
func (l *LinkedList[T]) PushBack(e T) {
	l.lazyInit()
	l.insert(l.root.prev, &Node[T]{element: e})
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
