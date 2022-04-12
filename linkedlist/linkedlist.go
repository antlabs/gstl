package linkedlist

type LinkedList[T any] struct {
	head   *Node[T]
	tail   *Node[T]
	length uint
}

type Node[T any] struct {
	next    *Node[T]
	prev    *Node[T]
	element T
}

func New[T any]() {

}

func (l *LinkedList[T]) Append() {

}

func (l *LinkedList[T]) PushFront(e T) {

}

func (l *LinkedList[T]) PushBack(e T) {

}

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
