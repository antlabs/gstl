package vec

// 参考文件如下
// https://doc.rust-lang.org/src/alloc/vec/mod.rs.html
// https://doc.rust-lang.org/std/vec/struct.Vec.html
import (
	"errors"
	"fmt"
	"github.com/guonaihong/gstl/cmp"
)

var ErrVecElemEmpty = errors.New("vec is empty")
var ErrLenGreaterCap = errors.New("len is too long > length of cap")

// vec结构体
type Vec[T any] struct {
	slice []T
}

// 初始化一个vec
func New[T any](a ...T) *Vec[T] {
	return &Vec[T]{slice: a}
}

// 初始化一个vec, 并指定底层的slice 容量
func WithCapacity[T any](capacity int) *Vec[T] {
	return &Vec[T]{slice: make([]T, 0, capacity)}
}

// 清空vec里面的所有值
func (v *Vec[T]) Clear() {
	v.slice = nil
}

// 删除连续重复值
func (v *Vec[T]) DedupFunc(cmp func(a, b T) bool) {
	l := v.Len()
	if l <= 1 {
		return
	}

	for i := 0; i < l; {
		right := -1
		for j := i + 1; j < l; {
			if !cmp(v.slice[i], v.slice[j]) {
				right = j
				break
			}
		}

		copy(v.slice[i:], v.slice[right:])
	}
}

// 从尾巴插入
// 支持插入一个值或者多个值
func (v *Vec[T]) Push(e ...T) {
	v.slice = append(v.slice, e...)
}

// 设置新长度
func (v *Vec[T]) SetLen(newLen int) {
	if newLen > v.Cap() {
		panic(ErrLenGreaterCap)
	}

	v.slice = v.slice[:newLen]
}

// 添加other类型的vec到v里面
func (v *Vec[T]) Append(other Vec[T]) {
	v.slice = append(v.slice, other.slice...)
}

// 添加other类型的slice到v里面
func (v *Vec[T]) Extend(other []T) {
	v.slice = append(v.slice, other...)
}

// 从尾巴弹出
func (v *Vec[T]) Pop() (e T, err error) {
	if v.Len() == 0 {
		return e, ErrVecElemEmpty
	}

	l := len(v.slice)
	e = v.slice[l-1]
	v.slice = v.slice[:l-1]

	// 缩容
	if v.Len()*2 < v.Cap() {
		slice := make([]T, v.Len())
		copy(slice, v.slice)
		v.slice = slice
	}

	return e, nil
}

// 返回slice底层的slice
func (v *Vec[T]) ToSlice() []T {
	return v.slice
}

// 往指定位置插入元素, 后面的元素往右移动
func (v *Vec[T]) Insert(i int, es ...T) *Vec[T] {
	l := v.Len()
	if i == l {
		v.Push(es...)
		return v
	}

	if i > l {
		panic(fmt.Sprintf("insertion index (is %d) should be <= len (is %d)", i, l))
	}

	need := l + len(es)
	if need > v.Cap() {
		v.Reserve(len(es))
	}

	s := v.slice
	//重置下s2 len
	newSlice := v.slice[:l+len(es)]
	copy(newSlice[i+len(es):], s[i:])
	copy(newSlice[i:], es)

	v.slice = newSlice
	return v
}

// 删除指定范围内的元素
func (v *Vec[T]) Delete(i, j int) *Vec[T] {
	copy(v.slice[i:], v.slice[j:])
	v.slice = v.slice[:v.Len()-(j-i)]
	return v
}

// 获取指定索引的值
func (v *Vec[T]) Get(index int) (e T) {
	return v.slice[index]
}

// 删除指定索引的元素, 空缺的位置, 使用最后一个元素替换上去
func (v *Vec[T]) SwapRemove(index int) (rv T) {
	l := v.Len()
	if index >= l {
		panic(fmt.Sprintf("swap_remove index (is %d) should be < len (is %d)", index, l))
	}

	rv = v.slice[index]
	v.slice[index] = v.slice[l-1]
	v.SetLen(l - 1)
	return
}

// 在给定索引处将vec拆分为两个
// 返回一个新的vec, 范围是[at, len)
// 原始的vec的范围是[0, at], 不改变原先的容量
func (v *Vec[T]) SplitOff(at int) (new *Vec[T]) {
	l := v.Len()

	if at > l {
		panic(fmt.Sprintf("`at` split index (is %d) should be <= len (is %d)", at, l))
	}

	if at == 0 {
		new = v.Clone()
		v.SetLen(0)
		return
	}

	newSlice := make([]T, l-at)
	copy(newSlice, v.slice[at:])
	return New(newSlice...)
}

// 删除指定索引的元素
func (v *Vec[T]) Remove(index int) int {
	l := v.Len()
	if index >= l {
		panic(fmt.Sprintf("removal index (is %d) should be < len (is %d)", index, l))
	}

	copy(v.slice[index:], v.slice[index+1:])
	v.slice = v.slice[:l-1]

	return v.Len()
}

// 提前在现有基础上再额外申请 additional 长度空间
// 可以避免频繁的重新分配
// 如果容量已经满足, 则什么事也不做
func (v *Vec[T]) Reserve(additional int) {
	v.reserve(additional, 1.2)
}

// 如果容量已经满足, 则什么事也不做
// 保留最小容量, 提前在现有基础上再额外申请 additional 长度空间
func (v *Vec[T]) ReserveExact(additional int) {
	v.reserve(additional, 1)
}

func (v *Vec[T]) reserve(additional int, factor float64) {
	l := v.Len()
	if l+additional <= v.Cap() {
		return
	}

	newSlice := make([]T, l, int(float64(l+additional)*factor))
	copy(newSlice, v.slice)
	v.slice = newSlice

}

// 向下收缩vec的容器
func (v *Vec[T]) ShrinkTo() {
	l := v.Len()
	if v.Cap() > l {
		v.ShrinkToFit(l)
	}
}

// 向下收缩vec的容器, 会重新分配底层的slice
func (v *Vec[T]) ShrinkToFit(minCapacity int) {
	cap := v.Cap()
	if cap > minCapacity {
		max := cmp.Max(cap, minCapacity)
		newSlice := append([]T{}, v.slice[:max]...)
		v.slice = newSlice
	}
}

// 修改vec可访问的容量, 但是不会修改底层的slice, 只是修改slice的len
func (v *Vec[T]) Truncate(newLen int) {
	v.slice = v.slice[:newLen]
}

// 在vec后面追加newLen 长度的value
func (v *Vec[T]) ExtendWith(newLen int, value T) {

	oldLen := v.Len()
	v.Reserve(newLen)
	v.slice = v.slice[:oldLen+newLen]
	for i, l := newLen, oldLen; i < l; i++ {
		v.slice[i] = value
	}
}

// 调整vec的大小, 使用len等于newLen
// 如果newLen > len, 差值部分会填充value
// 如果newLen < len, 多余的部分会被截断
func (v *Vec[T]) Resize(newLen int, value T) {
	l := v.Len()
	if newLen > l {
		v.ExtendWith(newLen-l, value)
		return
	}

	v.Truncate(newLen)
}

// 深度拷贝一份
func (v *Vec[T]) Clone() *Vec[T] {
	newSlice := make([]T, v.Len())
	copy(newSlice, v.slice)
	return New(newSlice...)
}

// 如果为空
func (v *Vec[T]) IsEmpty() bool {
	return len(v.slice) == 0
}

// len
func (v *Vec[T]) Len() int {
	return len(v.slice)
}

// cap
func (v *Vec[T]) Cap() int {
	return cap(v.slice)
}

// 返回第1个元素
func (v *Vec[T]) First() (n T, err error) {
	if v.Len() == 0 {
		return n, ErrVecElemEmpty
	}

	return v.slice[0], nil
}

// 删除vec第一个元素, 并返回它
func (v *Vec[T]) TakeFirst() (n T, err error) {
	if v.Len() == 0 {
		return n, ErrVecElemEmpty
	}

	n = v.slice[0]
	v.Remove(0)
	return n, nil
}

// 返回最后一个元素
func (v *Vec[T]) Last() (n T, err error) {
	if v.Len() == 0 {
		return n, ErrVecElemEmpty
	}

	return v.slice[v.Len()-1], nil
}

// 原地操作, 回调函数会返回的元素值
func (v *Vec[T]) Map(m func(e T) T) *Vec[T] {

	l := v.Len()

	for i := 0; i < l; i++ {
		v.slice[i] = m(v.slice[i])
	}

	return v
}

// 原地操作, 回调函数返回true的元素保留
func (v *Vec[T]) Filter(filter func(e T) bool) *Vec[T] {

	l := v.Len()
	left := 0

	for i := 0; i < l; i++ {
		if filter(v.slice[i]) {
			if left != i {
				v.slice[left] = v.slice[i]
			}
			left++
		}
	}
	v.SetLen(left)
	return v
}

// 原地旋转vec, 向左边旋转
func (v *Vec[T]) RotateLeft(n int) {
	l := v.Len()
	n %= l

	if n == 0 {
		return
	}

	left := make([]T, n)
	// 先备份左边
	copy(left, v.slice[:n])
	// 备下的往左拷贝
	copy(v.slice, v.slice[n:])
	// 右边需要被替换的空间
	copy(v.slice[l-n:], left)
}

//原地旋转vec, 向右边旋转
func (v *Vec[T]) RotateRight(n int) {
	l := v.Len()
	n %= l
	if n == 0 {
		return
	}

	rightVec := v.SplitOff(n)
	for right, left := l, n; right > 0 && left > 0; {
		v.slice[right] = v.slice[left]
		right--
		left--
	}

	copy(v.slice, rightVec.slice)
}

// 用于写入重复的值, 返回新的内存块, 来创建新的vec
func (v *Vec[T]) Repeat(count int) *Vec[T] {
	need := v.Len() * count
	rv := WithCapacity[T](need)

	for i := 0; i < count; i++ {
		rv.Append(*v)
	}

	return rv
}
