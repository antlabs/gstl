package vec

import (
	"testing"
)

// 每次Push 1个, pop 1个, 测试int类型
func Test_New_Push_Pop_Int(t *testing.T) {
	v := New(1, 2, 3, 4, 5, 6)
	v.Push(7)
	v.Push(8)
	n, _ := v.Pop()
	if n != 8 {
		t.Errorf("Expected 8, got %v", n)
	}
}

// 每次push 1个, pop 1个, 测试string类型
func Test_New_Push_Pop_String(t *testing.T) {
	v := New("1", "2", "3", "4", "5", "6")
	v.Push("7")
	v.Push("8")
	n, _ := v.Pop()
	if n != "8" {
		t.Errorf("Expected '8', got %v", n)
	}
}

// push一个slice, pop 1个, 测试string类型
func Test_New_Push_Slice_Pop_String(t *testing.T) {
	v := New("1", "2", "3")
	v.Push("4", "5", "6")
	if !slicesEqual(v.ToSlice(), []string{"1", "2", "3", "4", "5", "6"}) {
		t.Errorf("Expected %v, got %v", []string{"1", "2", "3", "4", "5", "6"}, v.ToSlice())
	}
}

// push一个slice, pop 1个, 测试string类型
func Test_New_Push_Slice_Pop_Int(t *testing.T) {
	v := New(1, 2, 3)
	v.Push(4, 5, 6)
	if !slicesEqual(v.ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected %v, got %v", []int{1, 2, 3, 4, 5, 6}, v.ToSlice())
	}
}

// 测试左移
func Test_RotateLeft(t *testing.T) {
	v := New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(2)
	if !slicesEqual(v.ToSlice(), []byte{'c', 'd', 'e', 'f', 'a', 'b'}) {
		t.Errorf("Expected %v, got %v", []byte{'c', 'd', 'e', 'f', 'a', 'b'}, v.ToSlice())
	}

	v = New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(8)
	if !slicesEqual(v.ToSlice(), []byte{'c', 'd', 'e', 'f', 'a', 'b'}) {
		t.Errorf("Expected %v, got %v", []byte{'c', 'd', 'e', 'f', 'a', 'b'}, v.ToSlice())
	}

	v = New[uint8]('a', 'b', 'c', 'd', 'e', 'f')
	v.RotateLeft(0)
	if !slicesEqual(v.ToSlice(), []byte{'a', 'b', 'c', 'd', 'e', 'f'}) {
		t.Errorf("Expected %v, got %v", []byte{'a', 'b', 'c', 'd', 'e', 'f'}, v.ToSlice())
	}
}

// 测试右移
func Test_RotateRight(t *testing.T) {
	v := New[uint8](1, 2, 3, 4, 5, 6)
	v.RotateRight(2)
	if !slicesEqual(v.ToSlice(), []byte{5, 6, 1, 2, 3, 4}) {
		t.Errorf("Expected %v, got %v", []byte{5, 6, 1, 2, 3, 4}, v.ToSlice())
	}

	v = New[uint8](1, 2, 3, 4, 5, 6)
	v.RotateRight(8)
	if !slicesEqual(v.ToSlice(), []byte{5, 6, 1, 2, 3, 4}) {
		t.Errorf("Expected %v, got %v", []byte{5, 6, 1, 2, 3, 4}, v.ToSlice())
	}

	v = New[uint8](1, 2, 3, 4, 5, 6)
	v.RotateRight(0)
	if !slicesEqual(v.ToSlice(), []byte{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected %v, got %v", []byte{1, 2, 3, 4, 5, 6}, v.ToSlice())
	}

	if !slicesEqual(New[byte](1, 2, 3, 4, 5, 6).RotateRight(0).ToSlice(), []byte{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected %v, got %v", []byte{1, 2, 3, 4, 5, 6}, New[byte](1, 2, 3, 4, 5, 6).RotateRight(0).ToSlice())
	}

	if !slicesEqual(New[byte](1, 2, 3, 4, 5, 6).RotateRight(6).ToSlice(), []byte{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected %v, got %v", []byte{1, 2, 3, 4, 5, 6}, New[byte](1, 2, 3, 4, 5, 6).RotateRight(6).ToSlice())
	}
}

// 测试填充
func Test_Repeat(t *testing.T) {
	if !slicesEqual(New(1, 2).Repeat(3).ToSlice(), []int{1, 2, 1, 2, 1, 2}) {
		t.Errorf("Expected %v, got %v", []int{1, 2, 1, 2, 1, 2}, New(1, 2).Repeat(3).ToSlice())
	}
	if !slicesEqual(New("hello").Repeat(2).ToSlice(), []string{"hello", "hello"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "hello"}, New("hello").Repeat(2).ToSlice())
	}
}

// 测试删除
func Test_Delete(t *testing.T) {
	if !slicesEqual(New[int](1, 2, 3, 4, 5).Delete(1, 2).ToSlice(), []int{1, 3, 4, 5}) {
		t.Errorf("Expected %v, got %v", []int{1, 3, 4, 5}, New[int](1, 2, 3, 4, 5).Delete(1, 2).ToSlice())
	}
	if !slicesEqual(New("hello", "world", "12345").Delete(1, 2).ToSlice(), []string{"hello", "12345"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "12345"}, New("hello", "world", "12345").Delete(1, 2).ToSlice())
	}
}

// 向指定位置插件数据
func Test_Insert(t *testing.T) {
	if !slicesEqual(New[int](1, 7).Insert(1, 2, 3, 4, 5, 6).ToSlice(), []int{1, 2, 3, 4, 5, 6, 7}) {
		t.Errorf("Expected %v, got %v", []int{1, 2, 3, 4, 5, 6, 7}, New[int](1, 7).Insert(1, 2, 3, 4, 5, 6).ToSlice())
	}
	if !slicesEqual(New("world", "12345").Insert(0, "hello").ToSlice(), []string{"hello", "world", "12345"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "world", "12345"}, New("world", "12345").Insert(0, "hello").ToSlice())
	}
	if !slicesEqual(New("hello", "12345").Insert(1, "world").ToSlice(), []string{"hello", "world", "12345"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "world", "12345"}, New("hello", "12345").Insert(1, "world").ToSlice())
	}
	if !slicesEqual(New("hello", "12345").Insert(2, "world").ToSlice(), []string{"hello", "12345", "world"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "12345", "world"}, New("hello", "12345").Insert(2, "world").ToSlice())
	}
}

// map函数, 修改函数里面的值, 不修改长度
func Test_Map(t *testing.T) {
	if !slicesEqual(New[int](1, 2, 3, 4, 5).Map(func(e int) int { return e * 2 }).ToSlice(), []int{2, 4, 6, 8, 10}) {
		t.Errorf("Expected %v, got %v", []int{2, 4, 6, 8, 10}, New[int](1, 2, 3, 4, 5).Map(func(e int) int { return e * 2 }).ToSlice())
	}
	if !slicesEqual(New[string]("world", "12345").Map(func(e string) string { return "#" + e }).ToSlice(), []string{"#world", "#12345"}) {
		t.Errorf("Expected %v, got %v", []string{"#world", "#12345"}, New[string]("world", "12345").Map(func(e string) string { return "#" + e }).ToSlice())
	}
}

// filter函数, 不修改函数里面的值, 只留满足条件的
func Test_Filter(t *testing.T) {
	if !slicesEqual(New(1, 2, 3, 4, 5).Filter(func(e int) bool { return e%2 == 0 }).ToSlice(), []int{2, 4}) {
		t.Errorf("Expected %v, got %v", []int{2, 4}, New(1, 2, 3, 4, 5).Filter(func(e int) bool { return e%2 == 0 }).ToSlice())
	}
	if !slicesEqual(New[int](1, 2, 3, 4, 5).Filter(func(e int) bool { return e%2 != 0 }).ToSlice(), []int{1, 3, 5}) {
		t.Errorf("Expected %v, got %v", []int{1, 3, 5}, New[int](1, 2, 3, 4, 5).Filter(func(e int) bool { return e%2 != 0 }).ToSlice())
	}
}

// 测试clear函数
func Test_Clear(t *testing.T) {
	v := WithCapacity[bool](10)
	v.Clear()
	if v.Cap() != 0 {
		t.Errorf("Expected 0, got %v", v.Cap())
	}
}

// 测试对vec去重
func Test_DedupFunc(t *testing.T) {
	v := New("1", "2", "2", "3")
	v2 := v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1", "2", "3"}) {
		t.Errorf("Expected %v, got %v", []string{"1", "2", "3"}, v2)
	}

	v = New("1", "2", "2")
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1", "2"}) {
		t.Errorf("Expected %v, got %v", []string{"1", "2"}, v2)
	}

	v = New("1", "2", "2", "3", "4", "5", "5", "6")
	bk := v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1", "2", "3", "4", "5", "6"}) {
		t.Errorf("Expected %v, got %v", []string{"1", "2", "3", "4", "5", "6"}, v2)
	}

	v = New("1", "2", "2", "3", "4", "5", "5", "6", "7")
	bk = v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1", "2", "3", "4", "5", "6", "7"}) {
		t.Errorf("Expected %v, got %v", []string{"1", "2", "3", "4", "5", "6", "7"}, v2)
	}

	v = New("1")
	bk = v.Clone()
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1"}) {
		t.Errorf("Expected %v, got %v", []string{"1"}, v2)
	}
	if !slicesEqual(bk.ToSlice(), v.ToSlice()) {
		t.Errorf("Expected %v, got %v", bk.ToSlice(), v.ToSlice())
	}

	v = New("1", "1")
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1"}) {
		t.Errorf("Expected %v, got %v", []string{"1"}, v2)
	}

	v = New("1", "1", "1", "1")
	v2 = v.DedupFunc(func(a, b string) bool { return a == b }).ToSlice()
	if !slicesEqual(v2, []string{"1"}) {
		t.Errorf("Expected %v, got %v", []string{"1"}, v2)
	}
}

// 测试SetLen测试
func Test_SetLen(t *testing.T) {
	v := WithCapacity[int](10)
	v.Push(1, 2, 3, 4, 5)
	if v.Len() != 5 {
		t.Errorf("Expected 5, got %v", v.Len())
	}
	if v.Cap() != 10 {
		t.Errorf("Expected 10, got %v", v.Cap())
	}
	v.SetLen(3)
	if v.Len() != 3 {
		t.Errorf("Expected 3, got %v", v.Len())
	}
}

// 测试append接口
func Test_Append(t *testing.T) {
	if !slicesEqual(New(1, 2, 3).Append(New(4, 5, 6)).ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected %v, got %v", []int{1, 2, 3, 4, 5, 6}, New(1, 2, 3).Append(New(4, 5, 6)).ToSlice())
	}
	if !slicesEqual(New("hello").Append(New("world")).ToSlice(), []string{"hello", "world"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "world"}, New("hello").Append(New("world")).ToSlice())
	}
}

// 测试Set接口
func Test_Set(t *testing.T) {
	if !slicesEqual(New(1, 2, 3).Set(0, 2).ToSlice(), []int{2, 2, 3}) {
		t.Errorf("Expected %v, got %v", []int{2, 2, 3}, New(1, 2, 3).Set(0, 2).ToSlice())
	}
}

// 测试Set接口
func Test_Get(t *testing.T) {
	if New(1, 2, 3).Get(0) != 1 {
		t.Errorf("Expected 1, got %v", New(1, 2, 3).Get(0))
	}
	if New(1, 2, 3).Get(1) != 2 {
		t.Errorf("Expected 2, got %v", New(1, 2, 3).Get(1))
	}
	if New(1, 2, 3).Get(2) != 3 {
		t.Errorf("Expected 3, got %v", New(1, 2, 3).Get(2))
	}
}

// 测试SwapRemove
func Test_SwapRemove(t *testing.T) {
	v := New(1, 2, 3)

	if v.SwapRemove(0) != 1 {
		t.Errorf("Expected 1, got %v", v.SwapRemove(0))
	}
	if !slicesEqual(v.ToSlice(), []int{3, 2}) {
		t.Errorf("Expected %v, got %v", []int{3, 2}, v.ToSlice())
	}

	v = New(1, 2, 3)

	if v.SwapRemove(2) != 3 {
		t.Errorf("Expected 3, got %v", v.SwapRemove(2))
	}
	if !slicesEqual(v.ToSlice(), []int{1, 2}) {
		t.Errorf("Expected %v, got %v", []int{1, 2}, v.ToSlice())
	}
}

// 测试SplitOff接口
func Test_SplitOff(t *testing.T) {
	v := New(1, 2, 3)
	v2 := v.SplitOff(0)

	if !slicesEqual(v2.ToSlice(), []int{1, 2, 3}) {
		t.Errorf("Expected %v, got %v", []int{1, 2, 3}, v2.ToSlice())
	}
	if !slicesEqual(v.ToSlice(), []int{}) {
		t.Errorf("Expected %v, got %v", []int{}, v.ToSlice())
	}
	v2.Set(0, 5)
	if !slicesEqual(v.ToSlice(), []int{}) {
		t.Errorf("Expected %v, got %v", []int{}, v.ToSlice())
	}

	v = New(1, 2, 3)
	v2 = v.SplitOff(1)
	if !slicesEqual(v2.ToSlice(), []int{2, 3}) {
		t.Errorf("Expected %v, got %v", []int{2, 3}, v2.ToSlice())
	}
	if !slicesEqual(v.ToSlice(), []int{1}) {
		t.Errorf("Expected %v, got %v", []int{1}, v.ToSlice())
	}
}

// 测试Remove接口
func Test_Remove(t *testing.T) {
	if !slicesEqual(New(1, 2, 3).Remove(0).ToSlice(), []int{2, 3}) {
		t.Errorf("Expected %v, got %v", []int{2, 3}, New(1, 2, 3).Remove(0).ToSlice())
	}
	if !slicesEqual(New(1, 2, 3).Remove(1).ToSlice(), []int{1, 3}) {
		t.Errorf("Expected %v, got %v", []int{1, 3}, New(1, 2, 3).Remove(1).ToSlice())
	}
	if !slicesEqual(New(1, 2, 3).Remove(2).ToSlice(), []int{1, 2}) {
		t.Errorf("Expected %v, got %v", []int{1, 2}, New(1, 2, 3).Remove(2).ToSlice())
	}
}

// 测试扩容
func Test_Reserve(t *testing.T) {
	if WithCapacity[int](3).Reserve(1).Cap() != 3 {
		t.Errorf("Expected 3, got %v", WithCapacity[int](3).Reserve(1).Cap())
	}
	if WithCapacity[int](3).Reserve(2).Cap() != 3 {
		t.Errorf("Expected 3, got %v", WithCapacity[int](3).Reserve(2).Cap())
	}
	if WithCapacity[int](3).Reserve(3).Cap() != 3 {
		t.Errorf("Expected 3, got %v", WithCapacity[int](3).Reserve(3).Cap())
	}
	if WithCapacity[int](3).Reserve(4).Cap() <= 4 {
		t.Errorf("Expected greater than 4, got %v", WithCapacity[int](3).Reserve(4).Cap())
	}
	if WithCapacity[int](3).Reserve(4).Cap() == 4 {
		t.Errorf("Expected not equal to 4, got %v", WithCapacity[int](3).Reserve(4).Cap())
	}
}

// 测试扩容
func Test_ReserveExact(t *testing.T) {
	if WithCapacity[int](3).ReserveExact(1).Cap() != 3 {
		t.Errorf("Expected 3, got %v", WithCapacity[int](3).ReserveExact(1).Cap())
	}
	if WithCapacity[int](3).ReserveExact(2).Cap() != 3 {
		t.Errorf("Expected 3, got %v", WithCapacity[int](3).ReserveExact(2).Cap())
	}
	if WithCapacity[int](3).ReserveExact(3).Cap() != 3 {
		t.Errorf("Expected 3, got %v", WithCapacity[int](3).ReserveExact(3).Cap())
	}
	if WithCapacity[int](3).ReserveExact(4).Cap() != 4 {
		t.Errorf("Expected 4, got %v", WithCapacity[int](3).ReserveExact(4).Cap())
	}
}

// 测试IsEmpty接口
func Test_IsEmpty(t *testing.T) {
	if !WithCapacity[int](3).IsEmpty() {
		t.Errorf("Expected true, got false")
	}
	if !WithCapacity[int](4).IsEmpty() {
		t.Errorf("Expected true, got false")
	}
	if !WithCapacity[int](5).IsEmpty() {
		t.Errorf("Expected true, got false")
	}
}

// 收缩内存
func Test_ShrinkTo(t *testing.T) {
	if WithCapacity[int](10).ShrinkTo(0).Cap() < 3 {
		t.Errorf("Expected greater than or equal to 3, got %v", WithCapacity[int](10).ShrinkTo(0).Cap())
	}
	if WithCapacity[int](10).Push(1, 2, 3).ShrinkTo(4).Cap() < 3 {
		t.Errorf("Expected greater than or equal to 3, got %v", WithCapacity[int](10).Push(1, 2, 3).ShrinkTo(4).Cap())
	}
	if WithCapacity[int](10).Push(1, 2, 3).ShrinkTo(4).Cap() < 3 {
		t.Errorf("Expected greater than or equal to 3, got %v", WithCapacity[int](10).Push(1, 2, 3).ShrinkTo(4).Cap())
	}
}

// 收缩内存
func Test_ShrinkToFit(t *testing.T) {
	if WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap() < 3 {
		t.Errorf("Expected greater than or equal to 3, got %v", WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap())
	}
	if WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap() < 3 {
		t.Errorf("Expected greater than or equal to 3, got %v", WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap())
	}

	if WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap() >= 10 {
		t.Errorf("Expected less than 10, got %v", WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap())
	}
	if WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap() >= 10 {
		t.Errorf("Expected less than 10, got %v", WithCapacity[int](10).Push(1, 2, 3).ShrinkToFit().Cap())
	}
}

// TODO
func Test_ExtendWith(t *testing.T) {
	if !slicesEqual(WithCapacity[string](10).ExtendWith(3, "hello").ToSlice(), []string{"hello", "hello", "hello"}) {
		t.Errorf("Expected %v, got %v", []string{"hello", "hello", "hello"}, WithCapacity[string](10).ExtendWith(3, "hello").ToSlice())
	}
}

func Test_Resize(t *testing.T) {
	if !slicesEqual(WithCapacity[string](10).Push("goto").Resize(3, "hello").ToSlice(), []string{"goto", "hello", "hello"}) {
		t.Errorf("Expected %v, got %v", []string{"goto", "hello", "hello"}, WithCapacity[string](10).Push("goto").Resize(3, "hello").ToSlice())
	}
}

func Test_Search(t *testing.T) {
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 7 <= e }) != 6 {
		t.Errorf("Expected 6, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 7 <= e }))
	}
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 6 <= e }) != 5 {
		t.Errorf("Expected 5, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 6 <= e }))
	}
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 5 <= e }) != 4 {
		t.Errorf("Expected 4, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 5 <= e }))
	}
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 4 <= e }) != 3 {
		t.Errorf("Expected 3, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 4 <= e }))
	}
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 3 <= e }) != 2 {
		t.Errorf("Expected 2, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 3 <= e }))
	}
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 2 <= e }) != 1 {
		t.Errorf("Expected 1, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 2 <= e }))
	}
	if New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 1 <= e }) != 0 {
		t.Errorf("Expected 0, got %v", New(1, 2, 3, 4, 5, 6, 7).SearchFunc(func(e int) bool { return 1 <= e }))
	}
}

// Helper function to compare slices
func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
