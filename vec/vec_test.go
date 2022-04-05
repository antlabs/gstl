package vec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 每次Push 1个, pop 1个, 测试int类型
func Test_New_Push_Pop_Int(t *testing.T) {
	v := New(1, 2, 3, 4, 5, 6)
	v.Push(7)
	v.Push(8)
	n, _ := v.Pop()
	assert.Equal(t, n, 8)
}

// 每次push 1个, pop 1个, 测试string类型
func Test_New_Push_Pop_String(t *testing.T) {
	v := New("1", "2", "3", "4", "5", "6")
	v.Push("7")
	v.Push("8")
	n, _ := v.Pop()
	assert.Equal(t, n, "8")
}

// push一个slice, pop 1个, 测试string类型
func Test_New_Push_Slice_Pop_String(t *testing.T) {
	v := New("1", "2", "3")
	v.Push("4", "5", "6")
	assert.Equal(t, v.ToSlice(), []string{"1", "2", "3", "4", "5", "6"})
}

// push一个slice, pop 1个, 测试string类型
func Test_New_Push_Slice_Pop_Int(t *testing.T) {
	v := New(1, 2, 3)
	v.Push(4, 5, 6)
	assert.Equal(t, v.ToSlice(), []int{1, 2, 3, 4, 5, 6})
}

// 测试左移
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

// 测试右移
func Test_RotateRight(t *testing.T) {
	v := New[uint8](1, 2, 3, 4, 5, 6)
	v.RotateRight(2)
	assert.Equal(t, v.ToSlice(), []byte{5, 6, 1, 2, 3, 4}, "test case:0")

	v = New[uint8](1, 2, 3, 4, 5, 6)
	v.RotateRight(8)
	assert.Equal(t, v.ToSlice(), []byte{5, 6, 1, 2, 3, 4}, "test case:1")

	v = New[uint8](1, 2, 3, 4, 5, 6)
	v.RotateRight(0)
	assert.Equal(t, v.ToSlice(), []byte{1, 2, 3, 4, 5, 6}, "test case:2")

	assert.Equal(t, New[byte](1, 2, 3, 4, 5, 6).RotateRight(0).ToSlice(), []byte{1, 2, 3, 4, 5, 6}, "test case:3")
	assert.Equal(t, New[byte](1, 2, 3, 4, 5, 6).RotateRight(6).ToSlice(), []byte{1, 2, 3, 4, 5, 6}, "test case:3")

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

// 测试clear函数
func Test_Clear(t *testing.T) {
	v := WithCapacity[bool](10)
	v.Clear()
	assert.Equal(t, v.Cap(), 0)
}

// 测试对vec去重
func Test_DedupFunc(t *testing.T) {
	v := New("1", "2", "2", "3")
	v2 := v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1", "2", "3"})

	v = New("1", "2", "2")
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1", "2"})

	v = New("1", "2", "2", "3", "4", "5", "5", "6")
	bk := v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1", "2", "3", "4", "5", "6"}, bk)

	v = New("1", "2", "2", "3", "4", "5", "5", "6", "7")
	bk = v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1", "2", "3", "4", "5", "6", "7"}, bk)

	v = New("1")
	bk = v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1"}, bk)

	v = New("1", "1")
	bk = v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1"}, bk)

	v = New("1", "1", "1", "1")
	bk = v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	assert.Equal(t, v2, []string{"1"}, bk)
}

// 测试SetLen测试
func Test_SetLen(t *testing.T) {
	v := WithCapacity[int](10)
	v.Push(1, 2, 3, 4, 5)
	assert.Equal(t, v.Len(), 5)
	assert.Equal(t, v.Cap(), 10)
	v.SetLen(3)
	assert.Equal(t, v.Len(), 3)
}

// 测试append接口
func Test_Append(t *testing.T) {
	assert.Equal(t, New(1, 2, 3).Append(New(4, 5, 6)).ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, New("hello").Append(New("world")).ToSlice(), []string{"hello", "world"})
}

// 测试Extend接口
func Test_Extend(t *testing.T) {
	assert.Equal(t, New(1, 2, 3).Extend([]int{4, 5, 6}).ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, New("hello").Extend([]string{"world"}).ToSlice(), []string{"hello", "world"})
}

// 测试Set接口
func Test_Set(t *testing.T) {
	assert.Equal(t, New(1, 2, 3).Set(0, 2).ToSlice(), []int{2, 2, 3})
}

// 测试Set接口
func Test_Get(t *testing.T) {
	assert.Equal(t, New(1, 2, 3).Get(0), 1)
	assert.Equal(t, New(1, 2, 3).Get(1), 2)
	assert.Equal(t, New(1, 2, 3).Get(2), 3)
}

// 测试SwapRemove
func Test_SwapRemove(t *testing.T) {
	v := New(1, 2, 3)

	assert.Equal(t, v.SwapRemove(0), 1)
	assert.Equal(t, v.ToSlice(), []int{3, 2})

	v = New(1, 2, 3)

	assert.Equal(t, v.SwapRemove(2), 3)
	assert.Equal(t, v.ToSlice(), []int{1, 2})
}

// 测试SplitOff接口
func Test_SplitOff(t *testing.T) {
	v := New(1, 2, 3)
	v2 := v.SplitOff(0)

	assert.Equal(t, v2.ToSlice(), []int{1, 2, 3})
	assert.Equal(t, v.ToSlice(), []int{})
	v2.Set(0, 5)
	assert.Equal(t, v.ToSlice(), []int{})

	v = New(1, 2, 3)
	v2 = v.SplitOff(1)
	assert.Equal(t, v2.ToSlice(), []int{2, 3})
	assert.Equal(t, v.ToSlice(), []int{1})
}

// 测试Remove接口
func Test_Remove(t *testing.T) {
	assert.Equal(t, New(1, 2, 3).Remove(0).ToSlice(), []int{2, 3})
	assert.Equal(t, New(1, 2, 3).Remove(1).ToSlice(), []int{1, 3})
	assert.Equal(t, New(1, 2, 3).Remove(2).ToSlice(), []int{1, 2})
}

// 测试扩容
func Test_Reserve(t *testing.T) {
	assert.Equal(t, WithCapacity[int](3).Reserve(1).Cap(), 3)
	assert.Equal(t, WithCapacity[int](3).Reserve(2).Cap(), 3)
	assert.Equal(t, WithCapacity[int](3).Reserve(3).Cap(), 3)
	assert.Greater(t, WithCapacity[int](3).Reserve(4).Cap(), 4)
	assert.NotEqual(t, WithCapacity[int](3).Reserve(4).Cap(), 4)
}

// 测试扩容
func Test_ReserveExact(t *testing.T) {
	assert.Equal(t, WithCapacity[int](3).ReserveExact(1).Cap(), 3)
	assert.Equal(t, WithCapacity[int](3).ReserveExact(2).Cap(), 3)
	assert.Equal(t, WithCapacity[int](3).ReserveExact(3).Cap(), 3)
	assert.Equal(t, WithCapacity[int](3).ReserveExact(4).Cap(), 4)
}

// 测试IsEmpty接口
func Test_IsEmpty(t *testing.T) {
	assert.True(t, WithCapacity[int](3).IsEmpty())
	assert.True(t, WithCapacity[int](4).IsEmpty())
	assert.True(t, WithCapacity[int](5).IsEmpty())
}
