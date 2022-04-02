package vec

// 参考文件如下
// https://doc.rust-lang.org/src/alloc/vec/mod.rs.html
import (
	"errors"
	"fmt"
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
func (v *Vec[T]) Dedup() {

}

// 从尾巴插入
func (v *Vec[T]) Push(e T) {
	v.slice = append(v.slice, e)
}

// 设置新长度
func (v *Vec[T]) SetLen(newLen int) {
	if newLen > v.Cap() {
		panic(ErrLenGreaterCap)
	}

	v.slice = v.slice[:newLen]
}

// 添加other类型的vec到v里面, 并且让other里面的数据为0
func (v *Vec[T]) Append(other Vec[T]) {
	v.slice = append(v.slice, other.slice...)
	other.SetLen(0)
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

//
func (v *Vec[T]) Insert() {

}

// 获取指定索引的值
func (v *Vec[T]) Get(index int) (e T) {
	return v.slice[index]
}

// 删除指定索引的元素, 空缺的位置, 使用最后一个元素替换上去
func (v *Vec[T]) SwapRemove(index int) {

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

//
func (v *Vec[T]) Reserve() {
}

func (v *Vec[T]) Truncate(newLen int) {
	v.slice = v.slice[:newLen]
}

//
func (v *Vec[T]) ExtendWith(newLen int, value T) {

}

// 缩容
func (v *Vec[T]) ShrinkToFit() {

}

//
//
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

// 原地操作, 回调函数会返回的元素值
func (v *Vec[T]) Map(m func(e T) T) {

	l := v.Len()

	for i := 0; i < l; i++ {
		newValue := m(v.slice[i])
		v.slice[i] = newValue
	}
}

// 原地操作, 回调函数返回true的元素保留
func (v *Vec[T]) Filter(filter func(e T) bool) {

	l := v.Len()
	left := 0

	for i := 0; i < l; i++ {
		if filter(v.slice[i]) {
			if left != i {
				v.slice[left] = v.slice[i]
				left++
			}
		}
	}
	v.SetLen(left)
}
