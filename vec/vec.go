package vec

// apache 2.0 guonaihong
// 参考文档如下
// https://doc.rust-lang.org/src/alloc/vec/mod.rs.html
// https://doc.rust-lang.org/std/vec/struct.Vec.html

import (
	"errors"
	"fmt"

	"github.com/guonaihong/gstl/cmp"
)

var (
	ErrVecElemEmpty  = errors.New("vec is empty")
	ErrLenGreaterCap = errors.New("len is too long > length of cap")
	ErrIndex         = errors.New("Illegal value of index")
)

const coefficient = 1.5

// vec类型
type Vec[T any] []T

// 初始化一个vec
func New[T any](a ...T) *Vec[T] {
	return (*Vec[T])(&a)
}

// 初始化函数, 可以把slice指针转成Vec类型
func FromSlicePtr[T any](ptr *[]T) *Vec[T] {
	return (*Vec[T])(ptr)
}

// 初始化一个vec, 并指定底层的slice 容量
func WithCapacity[T any](capacity int) *Vec[T] {
	p := make([]T, 0, capacity)
	return (*Vec[T])(&p)
}

// 清空vec里面的所有值
// TODO 需要看下效率. 如果效率不行,使用reflect.SliceHeader, 强转, 然后挨个置空
func (v *Vec[T]) Clear() {
	*v = []T{}
}

// 删除连续重复值
// TODO 优化. 寻找更优做法
func (v *Vec[T]) DedupFunc(cmp func(a, b T) bool) *Vec[T] {
	if v.Len() <= 1 {
		return v
	}

	slice := v.ToSlice()
	i := 0
	for i < len(slice) {
		j := i + 1
		for j < len(slice) && cmp(slice[i], slice[j]) {
			j++
		}

		if j != i+1 {
			copy(slice[i+1:], slice[j:])
			slice = slice[:len(slice)-(j-i-1)]

			//fmt.Printf("i = %d:%v\n", i, slice)
		}
		i++
	}

	*v = *New(slice...)
	return v
}

// 从尾巴插入
// 支持插入一个值或者多个值
func (v *Vec[T]) Push(e ...T) *Vec[T] {
	*v = append(*v, e...)
	return v
}

// 设置新长度
func (v *Vec[T]) SetLen(newLen int) {
	slice := []T(*v)
	if newLen > v.Cap() {
		panic(ErrLenGreaterCap)
	}

	slice = slice[:newLen]
	*v = Vec[T](slice)
}

// 添加other类型的vec到v里面
func (v *Vec[T]) Append(other *Vec[T]) *Vec[T] {
	*v = append(*v, other.ToSlice()...)
	return v
}

// 删除vec第一个元素, 并返回它, 和TakeFirst是同义词的关系
func (v *Vec[T]) PopFront() (e T, err error) {
	return v.TakeFirst()
}

// 从尾巴弹出
func (v *Vec[T]) Pop() (e T, err error) {
	l := v.Len()
	if l == 0 {
		return e, ErrVecElemEmpty
	}

	slice := v.ToSlice()
	e = slice[l-1]
	v = New(slice[:l-1]...)

	// 缩容
	if v.Len()*2 < v.Cap() {
		newSlice := make([]T, v.Len())
		copy(newSlice, slice)
		v = New(newSlice...)
	}

	return e, nil
}

// 返回slice底层的slice
func (v *Vec[T]) ToSlice() []T {
	return []T(*v)
}

// 往指定位置插入元素, 后面的元素往右移动
// i是位置, es可以是单个值和多个值
func (v *Vec[T]) Insert(i int, es ...T) *Vec[T] {
	l := v.Len()
	if i == l {
		// slice=1 2 3 4 insert(4, 5), result=1 2 3 4 5
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

	slice := v.ToSlice()

	// 插入之前: hello world
	// 插入之后: hello es world
	newSlice := slice[:need]
	copy(newSlice[i+len(es):], slice[i:]) //先往后挪
	copy(newSlice[i:], es)                //拷贝到i指定的位置

	// TODO 需要压测下, 这种写法是否慢
	*v = *New(newSlice...)
	return v
}

// 删除指定范围内的元素
func (v *Vec[T]) Delete(i, j int) *Vec[T] {
	slice := v.ToSlice()
	copy(slice[i:], slice[j:])
	*v = *New(slice[:v.Len()-(j-i)]...)
	return v
}

// 获取指定索引的值
func (v *Vec[T]) Get(index int) (e T) {
	slice := v.ToSlice()
	return slice[index]
}

// 获取指定索引的值, 如果索引不合法会返回错误
func (v *Vec[T]) GetWithErr(index int) (e T, err error) {
	if index < 0 || index >= v.Len() {
		err = ErrIndex
		return
	}
	return v.Get(index), nil
}

// 获取指定索引的指针
func (v *Vec[T]) GetPtr(index int) (e *T) {
	slice := v.ToSlice()
	return &slice[index]
}

// 设置指定索引的值
func (v *Vec[T]) Set(index int, value T) *Vec[T] {
	v.ToSlice()[index] = value
	return v
}

// 删除指定索引的元素, 空缺的位置, 使用最后一个元素替换上去
func (v *Vec[T]) SwapRemove(index int) (rv T) {
	l := v.Len()
	if index >= l {
		panic(fmt.Sprintf("SwapRemove index (is %d) should be < len (is %d)", index, l))
	}

	rv = v.Get(index)
	v.Set(index, v.Get(l-1))
	v.SetLen(l - 1)
	return
}

// 在给定索引处将vec拆分为两个
// 返回一个新的vec, 范围是[at, len), 这里需要注意
// 原始的vec的范围是[0, at), 不改变原先的容量
func (v *Vec[T]) SplitOff(at int) (new *Vec[T]) {
	l := v.Len()

	if at > l {
		panic(fmt.Sprintf("`at` split index (is %d) should be <= len (is %d)", at, l))
	}

	if at == 0 {
		v2 := *v
		v.Clear()
		return &v2
	}

	newSlice := make([]T, l-at)
	copy(newSlice, v.ToSlice()[at:])

	*v = *New(v.ToSlice()[:at]...)
	return New(newSlice...)
}

// 删除指定索引的元素
func (v *Vec[T]) Remove(index int) *Vec[T] {
	l := v.Len()
	if index >= l {
		panic(fmt.Sprintf("removal index (is %d) should be < len (is %d)", index, l))
	}

	copy(v.ToSlice()[index:], v.ToSlice()[index+1:])
	v.SetLen(l - 1)

	return v
}

// 提前在现有基础上再额外申请 additional 长度空间
// 可以避免频繁的重新分配
// 如果容量已经满足, 则什么事也不做
func (v *Vec[T]) Reserve(additional int) *Vec[T] {
	return v.reserve(additional, coefficient)
}

// 如果容量已经满足, 则什么事也不做
// 保留最小容量, 提前在现有基础上再额外申请 additional 长度空间
func (v *Vec[T]) ReserveExact(additional int) *Vec[T] {
	return v.reserve(additional, 1)
}

func (v *Vec[T]) reserve(additional int, factor float64) *Vec[T] {
	l := v.Len()
	if l+additional <= v.Cap() {
		return v
	}

	newSlice := make([]T, l, int(float64(l+additional)*factor))
	copy(newSlice, v.ToSlice())
	*v = Vec[T](newSlice)
	return v
}

// 向下收缩vec的容器
func (v *Vec[T]) ShrinkToFit() *Vec[T] {
	l := v.Len()
	if v.Cap() > getCap(l) {
		v.ShrinkTo(l)
	}
	return v
}

// 向下收缩vec的容器, 会重新分配底层的slice
func (v *Vec[T]) ShrinkTo(minCapacity int) *Vec[T] {
	cap := v.Cap()
	minCapacity = getCap(minCapacity)
	if cap > minCapacity {
		min := cmp.Min(cap, minCapacity)
		if min == 0 {
			min = int(0.66 * float64(cap))
		}

		newSlice := append([]T{}, v.ToSlice()[:min]...)
		*v = Vec[T](newSlice)
	}
	return v
}

// 修改vec可访问的容量, 但是不会修改底层的slice, 只是修改slice的len
func (v *Vec[T]) Truncate(newLen int) {
	*v = Vec[T](v.ToSlice()[:newLen])
}

// 在vec后面追加newLen 长度的value
func (v *Vec[T]) ExtendWith(newLen int, value T) *Vec[T] {

	oldLen := v.Len()
	v.Reserve(newLen)
	slice := v.ToSlice()

	l := oldLen + newLen
	slice = slice[:l]

	for i := oldLen; i < l; i++ {
		slice[i] = value
	}
	*v = Vec[T](slice)

	return v
}

// 调整vec的大小, 使用len等于newLen
// 如果newLen > len, 差值部分会填充value
// 如果newLen < len, 多余的部分会被截断
func (v *Vec[T]) Resize(newLen int, value T) *Vec[T] {
	l := v.Len()
	if newLen > l {
		v.ExtendWith(newLen-l, value)
		return v
	}

	v.Truncate(newLen)
	return v
}

// 深度拷贝一份
func (v *Vec[T]) Clone() *Vec[T] {
	newSlice := make([]T, v.Len())
	copy(newSlice, v.ToSlice())
	return (*Vec[T])(&newSlice)
}

// 如果为空
func (v *Vec[T]) IsEmpty() bool {
	return len(v.ToSlice()) == 0
}

// len
func (v *Vec[T]) Len() int {
	return len(v.ToSlice())
}

// cap
func (v *Vec[T]) Cap() int {
	return cap(v.ToSlice())
}

// 返回第1个元素
func (v *Vec[T]) First() (n T, err error) {
	if v.Len() == 0 {
		return n, ErrVecElemEmpty
	}

	return v.Get(0), nil
}

// 删除vec第一个元素, 并返回它
func (v *Vec[T]) TakeFirst() (n T, err error) {
	if v.Len() == 0 {
		return n, ErrVecElemEmpty
	}

	n = v.Get(0)
	v.Remove(0)
	return n, nil
}

// 返回最后一个元素
func (v *Vec[T]) Last() (n T, err error) {
	if v.Len() == 0 {
		return n, ErrVecElemEmpty
	}

	return v.Get(v.Len() - 1), nil
}

// 原地操作, 回调函数会返回的元素值
func (v *Vec[T]) Map(m func(e T) T) *Vec[T] {

	l := v.Len()

	slice := v.ToSlice()
	for i := 0; i < l; i++ {
		slice[i] = m(slice[i])
	}

	return v
}

// Retain 是Filter函数的同义词
func (v *Vec[T]) Retain(filter func(e T) bool) *Vec[T] {
	return v.Filter(filter)
}

// 原地操作, 回调函数返回true的元素保留
func (v *Vec[T]) Filter(filter func(e T) bool) *Vec[T] {

	l := v.Len()
	left := 0

	slice := v.ToSlice()
	for i := 0; i < l; i++ {
		if filter(slice[i]) {
			if left != i {
				slice[left] = slice[i]
			}
			left++
		}
	}
	v.SetLen(left)
	return v
}

// 原地旋转vec, 向左边旋转
func (v *Vec[T]) RotateLeft(n int) *Vec[T] {
	l := v.Len()
	n %= l

	if n == 0 {
		return v
	}

	slice := v.ToSlice()
	left := make([]T, n)
	// 先备份左边
	copy(left, slice[:n])
	// 备下的往左拷贝
	copy(slice, slice[n:])
	// 右边需要被替换的空间
	copy(slice[l-n:], left)

	return v
}

//原地旋转vec, 向右边旋转
func (v *Vec[T]) RotateRight(n int) *Vec[T] {
	l := v.Len()
	n %= l
	if n == 0 {
		return v
	}

	at := l - n
	slice := v.ToSlice()
	rightVec := make([]T, n)
	copy(rightVec, slice[at:])

	for right, left := l-1, at-1; right >= 0 && left >= 0; {
		slice[right] = slice[left]
		right--
		left--
	}

	copy(slice[:n], rightVec)
	return v
}

// 用于写入重复的值, 返回新的内存块, 来创建新的vec
func (v *Vec[T]) Repeat(count int) *Vec[T] {
	need := v.Len() * count
	rv := WithCapacity[T](need)

	for i := 0; i < count; i++ {
		rv.Append(v)
	}

	return rv
}

// 二分搜索
func (v *Vec[T]) SearchFunc(f func(T) bool) int {

	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, v.Len()
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if !f(v.Get(h)) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}

	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

// 遍历, callback 返回false就停止遍历, 返回true继续遍历
func (v *Vec[T]) Range(callback func(index int, v T) bool) *Vec[T] {
	slice := v.ToSlice()
	for i, val := range slice {
		if !callback(i, val) {
			return v
		}
	}
	return v
}

func getCap(l int) int {
	return int(float64(l) * coefficient)
}
