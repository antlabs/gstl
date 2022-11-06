package vecdeque

// apache 2.0 antlabs
import (
	"errors"
	"fmt"
	"math"

	"github.com/antlabs/gstl/cmp"
)

// 参考文档如下
// https://doc.rust-lang.org/std/collections/struct.VecDeque.html
// https://doc.rust-lang.org/src/alloc/collections/vec_deque/mod.rs.html
// https://doc.rust-lang.org/beta/src/alloc/collections/vec_deque/ring_slices.rs.html
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

// 返回当前使用的容量
func (v *VecDeque[T]) Len() int {
	return int(count(v.tail, v.head, uint(v.cap())))
}

// 统计数据
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

// 内存里面的物理容量
func (v *VecDeque[T]) cap() int {
	return len(v.buf)
}

// 业务意义上的容量, 有一个格式是空的
func (v *VecDeque[T]) Cap() int {
	return v.cap() - 1
}

// 对index 减去一些值
func (v *VecDeque[T]) wrapSub(index uint, subtrahend uint) uint {
	return v.wrapIndex(index - subtrahend)
}

// 对index 增加一些值
func (v *VecDeque[T]) wrapAdd(index uint, addend uint) uint {
	return v.wrapIndex(index + addend)
}

// 操作index的包装函数
func (v *VecDeque[T]) wrapIndex(index uint) uint {
	return wrapIndex(index, uint(v.cap()))
}

// 操作index的核心函数
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

// 交换索引为i和j的元素
func (v *VecDeque[T]) Swap(i, j uint) {
	ri := v.wrapAdd(v.tail, i)
	rj := v.wrapAdd(v.tail, j)
	v.buf[ri], v.buf[rj] = v.buf[rj], v.buf[ri]
}

// 向左旋转
func (v *VecDeque[T]) RotateLeftInner(k uint) {
	v.head = v.wrapAdd(v.head, k)
	v.tail = v.wrapAdd(v.tail, k)
}

// 向左旋转
func (v *VecDeque[T]) RotateLeft(k uint) {
	other := uint(v.Len()) - k

	if k <= other {
		v.RotateLeftInner(k)
		return
	}

	v.RotateRightInner(other)
}

// 向右旋转
func (v *VecDeque[T]) RotateRightInner(k uint) {
	//v.wrapCopy()
	v.head = v.wrapSub(v.head, k)
	v.tail = v.wrapSub(v.tail, k)
}

// 向右旋转
func (v *VecDeque[T]) RotateRight(k uint) {
	other := uint(v.Len()) - k
	if k <= other {
		// k = k
		// other = o
		// kkkkkkkkkkooo
		v.RotateRightInner(k)
		return
	}

	v.RotateLeftInner(other)

}

// 尽可能缩小VecDeque的容量
// 它将尽可能接近Len的位置
func (v *VecDeque[T]) ShrinkToFit() {
	v.ShrinkTo(0)
}

// 缩容
func (v *VecDeque[T]) ShrinkTo(minCapacity uint) {
	minCapacity = cmp.Min(minCapacity, uint(v.Cap()))
	minCapacity = cmp.Max(minCapacity, uint(v.Len()))
	targetCap := nextPowOfTwo(cmp.Max(minCapacity+1, MINIMUM_CAPACITY+1))

	if targetCap < uint(v.cap()) {

		//有三种情况值得关注：

		//所有元素都超出了预期范围

		//元素是连续的，head超出了所需的边界

		//元素是不连续的，尾部超出了期望的界限

		//

		//在所有其他时间，元素位置不受影响。

		//

		//指示应移动头部的元素。
		headOutside := v.head == 0 || v.head >= targetCap
		if v.tail >= targetCap && headOutside {
			//                    T             H
			//   [. . . . . . . . o o o o o o o . ]
			//    T             H
			//   [o o o o o o o . ]
			copy(v.buf, v.buf[v.tail:v.head])
			v.tail = 0
			v.head = uint(v.Len())
		} else if v.tail != 0 && v.tail < targetCap && headOutside {

			//          T             H
			//   [. . . o o o o o o o . . . . . . ]
			//        H T
			//   [o o . o o o o o ]
			length := v.wrapSub(v.head, targetCap)
			copy(v.buf, v.buf[targetCap:v.head])
			v.head = length
		} else if v.tail >= targetCap {

			//              H                 T
			//   [o o o o o . . . . . . . . . o o ]
			//              H T
			//   [o o o o o . o o ]
			length := uint(len(v.buf)) - v.tail
			newTail := targetCap - length
			copy(v.buf[newTail:], v.buf[v.tail:])
			v.tail = newTail
		}

		newBuf := make([]T, targetCap)
		copy(newBuf, v.buf)
		v.buf = newBuf
	}
}

func (v *VecDeque[T]) Truncate() {

}

func (v *VecDeque[T]) ToSlices() (first []T, second []T) {
	return
}

func (v *VecDeque[T]) wrapCopy(dst, src, length uint) {
	if src == dst || length == 0 {
		return
	}

}

func (v *VecDeque[T]) ReserveExact() {

}

func (v *VecDeque[T]) Reserve() {

}

func (v *VecDeque[T]) Contains(x T) bool {
	return false
}

// 获取第1个元素, 第二个参数返回错误
func (v *VecDeque[T]) Front() (e T, err error) {
	if v.Len() == 0 {
		err = ErrNoData
		return
	}

	return v.Get(0), nil
}

// 获取最后一个元素, 第二个参数返回错误
func (v *VecDeque[T]) Back() (e T, err error) {
	if v.Len() == 0 {
		err = ErrNoData
		return
	}

	newIndex := v.wrapSub(uint(v.Len()), uint(1))
	return v.Get(newIndex), nil
}

// 从 `VecDeque` 的任何位置删除一个元素并返回，并用第一个元素替换它。
func (v *VecDeque[T]) SwapRemoveFront(index uint) (e T, err error) {
	length := uint(v.Len())

	if index >= length {
		err = ErrNoData
		return
	}

	if length > 0 && index < length && index != 0 {
		v.Swap(index, 0)
	}

	return v.PopFront()

}

func (v *VecDeque[T]) SwapRemoveBack() {

}

// 在VecDeque内的index处插入一个元素, 所有索引大于或者等于'index'的元素向后移动
// TODO
func (v *VecDeque[T]) Insert(index uint, value T) {
	if v.IsFull() {
		v.grow()
	}

	// 移动环形缓冲区中最少的元素并插入
	// 给定对象
	//
	// 最多会移动len/2-1元素。O(min(n, n-i))
	//
	// 主要有三种情况：
	//  元素是连续的
	//      -尾部为0时的特殊情况
	//  元素不连续，插入部分位于尾部
	//  元素不连续，插入部分位于头部
	//
	// 对于每一种情况，还有两种情况：
	//  插入物更靠近尾部
	//  插入物更靠近头部
	//
	// key：H - v.head
	//      T - v.tail
	//      o - 有效元素
	//      I - 插入元素
	//      A - 应位于插入点之后的元素
	//      M - 表示元素已移动

	//idx := v.wrapAdd(v.tail, index)
	distanceToTail := index
	distanceToHead := uint(v.Len()) - index
	contiguous := v.isContiguous()

	if contiguous && distanceToTail < distanceToHead {

		if index == 0 {
			// push_front
			//
			//       T I             H
			//      [A o o o o o o . . . . . . .
			//      .
			//      .]
			//
			//                       H         T
			//      [A o o o o o o o . . . . . I]

			v.tail = v.wrapSub(v.tail, 1)
		} else {

		}
	}
}

func (v *VecDeque[T]) Remove(index int) {

}

func (v *VecDeque[T]) SplitOff() {

}

func (v *VecDeque[T]) Append(other *VecDeque[T]) {

}

func (v *VecDeque[T]) Retain() {

}

func (v *VecDeque[T]) ResizeWith() {

}

func (v *VecDeque[T]) isContiguous() bool {
	return v.tail <= v.head
}

func (v *VecDeque[T]) MakeContiguous() []T {
	if v.isContiguous() {
		return v.buf[v.tail:v.head]
	}

	cap := uint(v.cap())    //取出物理容量
	length := uint(v.Len()) // 取出已存元素个数
	free := v.tail - v.head // 空间的空间个数
	tailLen := cap - v.tail // tail到右顶边的个数

	if free >= tailLen {
		// 有足够的可用空间来一次性复制尾部，这意味着我们先将头向后移动，然后再将尾部复制到正确的位置。
		//
		//
		// 从: DEFGH....ABC 到: ABCDEFGH....

		// ...DEFGH.ABC
		copy(v.buf[tailLen:], v.buf[:v.head])
		// ABCDEFGH....
		copy(v.buf, v.buf[v.tail:])
		v.tail = 0
		v.head = length
		return v.buf[:v.head]
	}

	if free > v.head {
		// 有足够的自由空间可以一次性复制头部，这意味着我们先将尾部向前移动，然后再将头部复制到正确的位置。
		//
		//
		// 从: FGH....ABCDE 到: ...ABCDEFGH。
		//
		//

		// FGHABCDE....
		copy(v.buf[v.head:], v.buf[v.tail:])

		// ...ABCDEFGH.
		copy(v.buf[v.head+v.tail:], v.buf[:v.head])
		v.tail = v.head
		v.head = v.wrapAdd(v.tail, length)
	}

	// free 小于头和尾，这意味着我们必须缓慢地 "swap" 尾和头。
	//
	// 从: EFGHI...ABCD 或 HIJK.ABCDEFG
	// 到:   ABCDEFGHI... 或 ABCDEFGHIJK.
	leftEdge := uint(0)
	rightEdge := v.tail

	// The general problem looks like this
	// GHIJKLM...ABCDEF - before any swaps
	// ABCDEFM...GHIJKL - after 1 pass of swaps
	// ABCDEFGHIJM...KL - swap until the left edge reaches the temp store
	//                  - then restart the algorithm with a new (smaller) store
	// Sometimes the temp store is reached when the right edge is at the end
	// of the buffer - this means we've hit the right order with fewer swaps!
	// E.g
	// EF..ABCD
	// ABCDEF.. - after four only swaps we've finished

	// TODO 再仔细捋一捋逻辑
	for leftEdge < length && rightEdge != cap {
		rightOffset := uint(0)
		for i := leftEdge; i < rightEdge; i++ {
			rightOffset = (i - leftEdge) % (cap - rightEdge)
			src := (rightEdge + rightOffset)
			v.buf[i], v.buf[src] = v.buf[src], v.buf[i]
		}

		nOps := rightEdge - leftEdge
		leftEdge += nOps
		rightEdge += rightOffset + 1
	}

	v.tail = 0
	v.head = length

	return v.buf[v.tail:v.head]
}

func (v *VecDeque[T]) BinarySearch() {

}
