package vecdeque

import (
	"errors"
	"fmt"
	"github.com/guonaihong/gstl/cmp"
	"math"
)

// 参考文档如下
// https://doc.rust-lang.org/std/collections/struct.VecDeque.html
// https://doc.rust-lang.org/src/alloc/collections/vec_deque/mod.rs.html
// 翻译好的中文文档
// https://rustwiki.org/zh-CN/src/alloc/collections/vec_deque/mod.rs.html

const (
	INITIAL_CAPACITY uint = 7 // 2^3 - 1
	MINIMUM_CAPACITY uint = 1 // 2 - 1
)

var (
	ErrNoData = errors.New("no data")
)

type VecDeque[T any] struct {
	// tail 总是指向可以读取的第一个元素
	// head 只是指向应该写入数据的位置
	// 如果tail == head, 则缓存区为空. 环形缓冲区的长度定义为两者之间的距离
	tail uint
	head uint
	buf  []T
}

// 初始化
func New[T any]() *VecDeque[T] {
	return WithCapacity[T](int(INITIAL_CAPACITY))
}

// 初始VecDeque, 并设置实际需要的容量
func WithCapacity[T any](capacity int) *VecDeque[T] {
	cap := nextPowOfTwo(cmp.Max(uint(capacity)+1, MINIMUM_CAPACITY+1))
	return &VecDeque[T]{buf: make([]T, cap, cap)}
}

// 如果缓冲区满了. 就返回true
func (v *VecDeque[T]) IsFull() bool {
	return v.Cap()-v.Len() == 1
}

// TODO
func (v *VecDeque[T]) Len() int {
	return int(count(v.tail, v.head, uint(v.cap())))
}

func count(tail, head, size uint) uint {
	// 结果和 math.Abs(head - tail) & (size -1) 一样
	return (head - tail) & (size - 1)
}

// 扩容
func (v *VecDeque[T]) grow() *VecDeque[T] {
	if v.IsFull() {
		oldCap := v.cap()
		newBuf := make([]T, oldCap*2)
		copy(newBuf, v.buf)
		v.buf = newBuf
		v.handleCapIncrease(uint(oldCap))
	}
	return v
}

// 扩容
func (v *VecDeque[T]) handleCapIncrease(oldCap uint) {
	// Move the shortest contiguous section of the ring buffer
	//    T             H
	//   [o o o o o o o . ]
	//    T             H
	// A [o o o o o o o . . . . . . . . . ]
	//        H T
	//   [o o . o o o o o ]
	//          T             H
	// B [. . . o o o o o o o . . . . . . ]
	//              H T
	//   [o o o o o . o o ]
	//              H                 T
	// C [o o o o o . . . . . . . . . o o ]
	if v.tail <= v.head {
		// 不需要做啥
		return
	}

	// 把前面的数据移到后面, 合并起来, 中间没有空隙
	if v.head < oldCap-v.tail {
		copy(v.buf[oldCap:], v.buf[:v.head])
		v.head += oldCap
		return
	}

	// 把老的cap右边的数据放到新的cap的最右端
	newTail := oldCap + v.tail
	copy(v.buf[newTail:], v.buf[v.tail:oldCap])
	v.tail = newTail
}

// 判断VecDeque
func (v *VecDeque[T]) IsEmpty() bool {
	return v.tail == v.head
}

// 删除最后一个元素, 并且返回它. 如果为空, 返回ErrNoData
func (v *VecDeque[T]) PopBack() (value T, err error) {
	if v.IsEmpty() {
		err = ErrNoData
		return
	}

	v.head = v.wrapSub(v.head, 1)
	value = v.buf[v.head]
	return
}

// 删除第一个元素, 并且返回它, 如果为空, 返回ErrNoData
func (v *VecDeque[T]) PopFront() (value T, err error) {
	if v.IsEmpty() {
		err = ErrNoData
		return
	}

	value = v.buf[v.tail]
	v.tail = v.wrapAdd(v.tail, 1)
	return
}

// 将一个元素添加到VecDeque 后面
func (v *VecDeque[T]) PushBack(value T) {

	// 先检查是否满了
	if v.IsFull() {
		// 满了就扩容
		v.grow()
	}

	head := v.head

	// 修改head的值
	v.head = v.wrapAdd(v.head, uint(1))

	v.buf[head] = value
	// 修改head值
}

// 将一个元素添加到VecDeque的前面
func (v *VecDeque[T]) PushFront(value T) {
	if v.IsFull() {
		v.grow()
	}

	v.tail = v.wrapSub(v.tail, 1)
	v.buf[v.tail] = value
}

// 根据索引获取指定的值
func (v *VecDeque[T]) Get(i uint) T {
	idx := v.wrapAdd(v.tail, uint(i))
	return v.buf[idx]
}

func (v *VecDeque[T]) cap() int {
	return len(v.buf)
}

func (v *VecDeque[T]) Cap() int {
	return v.cap() - 1
}

func (v *VecDeque[T]) wrapSub(index uint, subtrahend uint) uint {
	return v.wrapIndex(index - subtrahend)
}

func (v *VecDeque[T]) wrapAdd(index uint, addend uint) uint {
	return v.wrapIndex(index + addend)
}

func (v *VecDeque[T]) wrapIndex(index uint) uint {
	return wrapIndex(index, uint(v.cap()))
}

func wrapIndex(index uint, size uint) uint {
	// 判断size是否是2的n次方
	if n := (size & (size - 1)); n != 0 {
		panic(fmt.Sprintf("size is always a power of 2, the current size is %d", size))
	}

	return index & (size - 1)
}

// TODO 优化下
// 使用更好的算法计算
func nextPowOfTwo(n uint) uint {

	for i := 1; i < 32; i++ {

		if nextPowOfTwoNum := math.Pow(2, float64(i)); nextPowOfTwoNum > float64(n) {
			return uint(nextPowOfTwoNum)
		}
	}

	return 0
}
