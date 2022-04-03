package vec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_Push_Pop_Int(t *testing.T) {
	v := New(1, 2, 3, 4, 5, 6)
	v.Push(7)
	v.Push(8)
	n, _ := v.Pop()
	assert.Equal(t, n, 8)
}

func Test_New_Push_Pop_String(t *testing.T) {
	v := New("1", "2", "3", "4", "5", "6")
	v.Push("7")
	v.Push("8")
	n, _ := v.Pop()
	assert.Equal(t, n, "8")
}

func Test_RotateLeft(t *testing.T) {
	v := New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(2)
	assert.Equal(t, v.ToSlice(), []byte{'c', 'd', 'e', 'f', 'a', 'b'})

	v = New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(8)
	assert.Equal(t, v.ToSlice(), []byte{'c', 'd', 'e', 'f', 'a', 'b'})

	v = New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(0)
	assert.Equal(t, v.ToSlice(), []byte{'a', 'b', 'c', 'd', 'e', 'f'})
}

// 测试填充
func Test_Repeat(t *testing.T) {
	assert.Equal(t, New(1, 2).Repeat(3).ToSlice(), []int{1, 2, 1, 2, 1, 2})
	assert.Equal(t, New("hello").Repeat(2).ToSlice(), []string{"hello", "hello"})
}

// 测试删除
func Test_Delete(t *testing.T) {
	assert.Equal(t, New[int](1, 2, 3, 4, 5).Delete(1, 2).ToSlice(), []int{1, 3, 4, 5})
	assert.Equal(t, New("hello", "world", "12345").Delete(1, 2).ToSlice(), []string{"hello", "12345"})
}

// 向指定位置插件数据
func Test_Insert(t *testing.T) {
	assert.Equal(t, New[int](1, 7).Insert(1, 2, 3, 4, 5, 6).ToSlice(), []int{1, 2, 3, 4, 5, 6, 7})
	assert.Equal(t, New("world", "12345").Insert(0, "hello").ToSlice(), []string{"hello", "world", "12345"})
}

// map函数, 修改函数里面的值, 不修改长度
func Test_Map(t *testing.T) {
	assert.Equal(t, New[int](1, 2, 3, 4, 5).Map(func(e int) int { return e * 2 }).ToSlice(), []int{2, 4, 6, 8, 10})
	assert.Equal(t, New[string]("world", "12345").Map(func(e string) string { return "#" + e }).ToSlice(), []string{"#world", "#12345"})
}

// filter函数, 不修改函数里面的值, 只留满足条件的
func Test_Filter(t *testing.T) {
	assert.Equal(t, New(1, 2, 3, 4, 5).Filter(func(e int) bool { return e%2 == 0 }).ToSlice(), []int{2, 4})
	assert.Equal(t, New[int](1, 2, 3, 4, 5).Filter(func(e int) bool { return e%2 != 0 }).ToSlice(), []int{1, 3, 5})
}
