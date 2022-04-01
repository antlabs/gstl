package vec

// 函数取名参考如下文档
// https://doc.rust-lang.org/src/alloc/vec/mod.rs.html
import (
	"errors"
)

var ErrEmpty = errors.New("vec is empty")

type Vec[T any] struct {
	slice []T
}

// 初始化一个vec
func New[T any](a ...T) *Vec[T] {
	return &Vec[T]{slice: a}
}

// 从尾巴插入
func (v *Vec[T]) Push(a T) {
	v.slice = append(v.slice, a)
}

// 从尾巴弹出
func (v *Vec[T]) Pop() (e T, err error) {
	if v.Len() == 0 {
		return e, ErrEmpty
	}

	l := len(v.slice)
	e = v.slice[l-1]
	v.slice = v.slice[:l-1]

	// 缩容
	if v.Len()*2 < v.Cap() {
		slice := make(T, v.Len())
		copy(slice, v.slice)
		v.slice = slice
	}

	return e, nil
}

// 删除
func (V *Vec[T]) Remove(index int) int {
	l := v.Len()
	if index >= l {
		panic(fmt.Sprintf("removal index (is %d) should be < len (is %d)", index, l))
	}

	copy(v.slice[index:], v.slice[index+1:])
	v.slice = v.slice[:l-1]
}

// 反转
func (v *Vec[T]) Reserve() {
	for i, l := 0, len(v.slice)-1; i < l; {
		v.slice[i], v.slice[l] = v.slice[l], v.slice[i]
		i++
		l--
	}
}

// 深度拷贝一份
func (v *Vec[T]) Clone() *Vec[T] {
	return New(v.slice...)
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
