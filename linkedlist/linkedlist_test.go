package linkedlist

import (
	"testing"

	"github.com/antlabs/gstl/must"
)

func Test_Push(t *testing.T) {
	if got := New[string]().PushBack("1", "2", "3", "4").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("PushBack() = %v, want %v", got, []string{"1", "2", "3", "4"})
	}
	if got := New[int]().PushBack(1, 2, 3, 4).ToSlice(); !sliceEqual(got, []int{1, 2, 3, 4}) {
		t.Errorf("PushBack() = %v, want %v", got, []int{1, 2, 3, 4})
	}
}

func Test_RPush(t *testing.T) {
	if got := New[string]().RPush("1", "2", "3", "4").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("RPush() = %v, want %v", got, []string{"1", "2", "3", "4"})
	}
	if got := New[int]().RPush(1, 2, 3, 4).ToSlice(); !sliceEqual(got, []int{1, 2, 3, 4}) {
		t.Errorf("RPush() = %v, want %v", got, []int{1, 2, 3, 4})
	}
}

func Test_Len(t *testing.T) {
	if got := New[string]().PushBack("1", "2", "3", "4").Len(); got != 4 {
		t.Errorf("Len() = %v, want %v", got, 4)
	}
	if got := New[int]().PushBack(1, 2, 3, 4).Len(); got != 4 {
		t.Errorf("Len() = %v, want %v", got, 4)
	}
}

func Test_RPop(t *testing.T) {
	// 不正常的索引
	if got := New[string]().PushBack("1", "2", "3", "4").RPop(-3); !sliceEqual(got, []string(nil)) {
		t.Errorf("RPop() = %v, want %v", got, []string(nil))
	}

	if got := New[string]().PushBack("1", "2", "3", "4").RPop(3); !sliceEqual(got, []string{"2", "3", "4"}) {
		t.Errorf("RPop() = %v, want %v", got, []string{"2", "3", "4"})
	}
	if got := New[int]().PushBack(1, 2, 3, 4).RPop(3); !sliceEqual(got, []int{2, 3, 4}) {
		t.Errorf("RPop() = %v, want %v", got, []int{2, 3, 4})
	}
	if got := New[string]().PushBack("1", "2", "3", "4").RPop(10); !sliceEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("RPop() = %v, want %v", got, []string{"1", "2", "3", "4"})
	}
}

func Test_LPop(t *testing.T) {
	// 不正常的索引
	if got := New[string]().PushBack("1", "2", "3", "4").LPop(-3); !sliceEqual(got, []string(nil)) {
		t.Errorf("LPop() = %v, want %v", got, []string(nil))
	}

	if got := New[string]().PushBack("1", "2", "3", "4").LPop(3); !sliceEqual(got, []string{"1", "2", "3"}) {
		t.Errorf("LPop() = %v, want %v", got, []string{"1", "2", "3"})
	}
	if got := New[int]().PushBack(1, 2, 3, 4).LPop(3); !sliceEqual(got, []int{1, 2, 3}) {
		t.Errorf("LPop() = %v, want %v", got, []int{1, 2, 3})
	}
	if got := New[string]().PushBack("1", "2", "3", "4").LPop(10); !sliceEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("LPop() = %v, want %v", got, []string{"1", "2", "3", "4"})
	}
}

func Test_RangeSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangeSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	if !sliceEqual(all, []string{"1", "2", "3", "4"}) {
		t.Errorf("RangeSafe() = %v, want %v", all, []string{"1", "2", "3", "4"})
	}
}

func Test_RangePrevSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangePrevSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	if !sliceEqual(all, []string{"4", "3", "2", "1"}) {
		t.Errorf("RangePrevSafe() = %v, want %v", all, []string{"4", "3", "2", "1"})
	}
}

func Test_First(t *testing.T) {
	// 没有值
	if got := must.TakeOneBool(New[string]().First()); got != false {
		t.Errorf("First() = %v, want %v", got, false)
	}
	// 有值
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").First()); got != "1" {
		t.Errorf("First() = %v, want %v", got, "1")
	}
	if got := must.TakeOneDiscardBool(New[int]().RPush(1, 2).First()); got != 1 {
		t.Errorf("First() = %v, want %v", got, 1)
	}
}

func Test_Last(t *testing.T) {
	// 没有值
	if got := must.TakeOneBool(New[string]().Last()); got != false {
		t.Errorf("Last() = %v, want %v", got, false)
	}
	//有值
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").Last()); got != "4" {
		t.Errorf("Last() = %v, want %v", got, "4")
	}
	if got := must.TakeOneDiscardBool(New[int]().RPush(1, 2).Last()); got != 2 {
		t.Errorf("Last() = %v, want %v", got, 2)
	}
}

func Test_Get(t *testing.T) {
	// 正索引
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(0)); got != "1" {
		t.Errorf("Get() = %v, want %v", got, "1")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(1)); got != "2" {
		t.Errorf("Get() = %v, want %v", got, "2")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(2)); got != "3" {
		t.Errorf("Get() = %v, want %v", got, "3")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(3)); got != "4" {
		t.Errorf("Get() = %v, want %v", got, "4")
	}

	// 负索引
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(-1)); got != "4" {
		t.Errorf("Get() = %v, want %v", got, "4")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(-2)); got != "3" {
		t.Errorf("Get() = %v, want %v", got, "3")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(-3)); got != "2" {
		t.Errorf("Get() = %v, want %v", got, "2")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4").GetWithBool(-4)); got != "1" {
		t.Errorf("Get() = %v, want %v", got, "1")
	}
}

func Test_Set(t *testing.T) {
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(11, "1991").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "5", "6"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1", "2", "3", "4", "5", "6"})
	}

	// 正索引
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(0, "1991").ToSlice(); !sliceEqual(got, []string{"1991", "2", "3", "4", "5", "6"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1991", "2", "3", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(1, "1991").ToSlice(); !sliceEqual(got, []string{"1", "1991", "3", "4", "5", "6"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1", "1991", "3", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(2, "1991").ToSlice(); !sliceEqual(got, []string{"1", "2", "1991", "4", "5", "6"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1", "2", "1991", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(3, "1991").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "1991", "5", "6"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1", "2", "3", "1991", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(4, "1991").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "1991", "6"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1", "2", "3", "4", "1991", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Set(5, "1991").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "5", "1991"}) {
		t.Errorf("Set() = %v, want %v", got, []string{"1", "2", "3", "4", "5", "1991"})
	}

	// 负索引
}

func Test_Index(t *testing.T) {
	if got := must.TakeOneBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(11)); got != false {
		t.Errorf("Index() = %v, want %v", got, false)
	}

	// 正索引
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(0)); got != "1" {
		t.Errorf("Index() = %v, want %v", got, "1")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(1)); got != "2" {
		t.Errorf("Index() = %v, want %v", got, "2")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(2)); got != "3" {
		t.Errorf("Index() = %v, want %v", got, "3")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(3)); got != "4" {
		t.Errorf("Index() = %v, want %v", got, "4")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(4)); got != "5" {
		t.Errorf("Index() = %v, want %v", got, "5")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(5)); got != "6" {
		t.Errorf("Index() = %v, want %v", got, "6")
	}

	// 负索引
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-1)); got != "6" {
		t.Errorf("Index() = %v, want %v", got, "6")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-2)); got != "5" {
		t.Errorf("Index() = %v, want %v", got, "5")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-3)); got != "4" {
		t.Errorf("Index() = %v, want %v", got, "4")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-4)); got != "3" {
		t.Errorf("Index() = %v, want %v", got, "3")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-5)); got != "2" {
		t.Errorf("Index() = %v, want %v", got, "2")
	}
	if got := must.TakeOneDiscardBool(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-6)); got != "1" {
		t.Errorf("Index() = %v, want %v", got, "1")
	}
}

func Test_Remove(t *testing.T) {
	// 正索引
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(0).ToSlice(); !sliceEqual(got, []string{"2", "3", "4", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"2", "3", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(1).ToSlice(); !sliceEqual(got, []string{"1", "3", "4", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "3", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(2).ToSlice(); !sliceEqual(got, []string{"1", "2", "4", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(3).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "3", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(4).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "3", "4", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(5).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "5"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "3", "4", "5"})
	}

	// 负索引
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-1).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "5"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "3", "4", "5"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-2).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "3", "4", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-3).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "3", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-4).ToSlice(); !sliceEqual(got, []string{"1", "2", "4", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "2", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-5).ToSlice(); !sliceEqual(got, []string{"1", "3", "4", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"1", "3", "4", "5", "6"})
	}
	if got := New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-6).ToSlice(); !sliceEqual(got, []string{"2", "3", "4", "5", "6"}) {
		t.Errorf("Remove() = %v, want %v", got, []string{"2", "3", "4", "5", "6"})
	}
}

func Test_LPush(t *testing.T) {
	if got := New[string]().ToSlice(); !sliceEqual(got, []string(nil)) {
		t.Errorf("LPush() = %v, want %v", got, []string(nil))
	}
	if got := New[string]().LPush("1").ToSlice(); !sliceEqual(got, []string{"1"}) {
		t.Errorf("LPush() = %v, want %v", got, []string{"1"})
	}
	if got := New[string]().LPush("2", "1").ToSlice(); !sliceEqual(got, []string{"1", "2"}) {
		t.Errorf("LPush() = %v, want %v", got, []string{"1", "2"})
	}
	if got := New[string]().LPush("3", "2", "1").ToSlice(); !sliceEqual(got, []string{"1", "2", "3"}) {
		t.Errorf("LPush() = %v, want %v", got, []string{"1", "2", "3"})
	}
	if got := New[string]().LPush("4", "3", "2", "1").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("LPush() = %v, want %v", got, []string{"1", "2", "3", "4"})
	}
}

func Test_PushFront(t *testing.T) {
	if got := New[string]().ToSlice(); !sliceEqual(got, []string(nil)) {
		t.Errorf("PushFront() = %v, want %v", got, []string(nil))
	}
	if got := New[string]().PushFront("1").ToSlice(); !sliceEqual(got, []string{"1"}) {
		t.Errorf("PushFront() = %v, want %v", got, []string{"1"})
	}
	if got := New[string]().PushFront("2", "1").ToSlice(); !sliceEqual(got, []string{"1", "2"}) {
		t.Errorf("PushFront() = %v, want %v", got, []string{"1", "2"})
	}
	if got := New[string]().PushFront("3", "2", "1").ToSlice(); !sliceEqual(got, []string{"1", "2", "3"}) {
		t.Errorf("PushFront() = %v, want %v", got, []string{"1", "2", "3"})
	}
	if got := New[string]().PushFront("4", "3", "2", "1").ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("PushFront() = %v, want %v", got, []string{"1", "2", "3", "4"})
	}
}

func Test_Clear(t *testing.T) {
	if got := New[string]().PushFront("1").Clear().ToSlice(); !sliceEqual(got, []string(nil)) {
		t.Errorf("Clear() = %v, want %v", got, []string(nil))
	}
	if got := New[string]().PushFront("2", "1").Clear().ToSlice(); !sliceEqual(got, []string(nil)) {
		t.Errorf("Clear() = %v, want %v", got, []string(nil))
	}
	if got := New[string]().PushFront("3", "2", "1").Clear().ToSlice(); !sliceEqual(got, []string(nil)) {
		t.Errorf("Clear() = %v, want %v", got, []string(nil))
	}
	if got := New[string]().PushFront("4", "3", "2", "1").Clear().ToSlice(); !sliceEqual(got, []string(nil)) {
		t.Errorf("Clear() = %v, want %v", got, []string(nil))
	}
}

func Test_IsEmpty(t *testing.T) {
	if got := New[string]().IsEmpty(); got != true {
		t.Errorf("IsEmpty() = %v, want %v", got, true)
	}
	if got := New[string]().PushFront("1").Clear().IsEmpty(); got != true {
		t.Errorf("IsEmpty() = %v, want %v", got, true)
	}
	if got := New[string]().PushFront("2", "1").Clear().IsEmpty(); got != true {
		t.Errorf("IsEmpty() = %v, want %v", got, true)
	}
	if got := New[string]().PushFront("3", "2", "1").Clear().IsEmpty(); got != true {
		t.Errorf("IsEmpty() = %v, want %v", got, true)
	}
	if got := New[string]().PushFront("4", "3", "2", "1").Clear().IsEmpty(); got != true {
		t.Errorf("IsEmpty() = %v, want %v", got, true)
	}
}

func Test_Trim(t *testing.T) {
	// 返回空值的情况
	if got := New[string]().RPush("one", "two", "three").Trim(100, -2).ToSlice(); !sliceEqual(got, []string{"one", "two", "three"}) {
		t.Errorf("Trim() = %v, want %v", got, []string{"one", "two", "three"})
	}

	if got := New[string]().RPush("one", "two", "three").Trim(1, -1).ToSlice(); !sliceEqual(got, []string{"two", "three"}) {
		t.Errorf("Trim() = %v, want %v", got, []string{"two", "three"})
	}
	if got := New[string]().RPush("one", "two", "three").Trim(0, 1).ToSlice(); !sliceEqual(got, []string{"one", "two"}) {
		t.Errorf("Trim() = %v, want %v", got, []string{"one", "two"})
	}
	if got := New[string]().RPush("one", "two", "three").Trim(-2, -1).ToSlice(); !sliceEqual(got, []string{"two", "three"}) {
		t.Errorf("Trim() = %v, want %v", got, []string{"two", "three"})
	}
	if got := New[string]().RPush("one", "two", "three").Trim(-100, -2).ToSlice(); !sliceEqual(got, []string{"one", "two"}) {
		t.Errorf("Trim() = %v, want %v", got, []string{"one", "two"})
	}
}

func Test_Range(t *testing.T) {
	var all []string

	// 1.1遍历全部
	all = make([]string, 0, 3)
	New[string]().RPush("one", "two", "three").Range(func(s string) {
		all = append(all, s)
	}, 0, -1)

	if !sliceEqual(all, []string{"one", "two", "three"}) {
		t.Errorf("Range() = %v, want %v", all, []string{"one", "two", "three"})
	}

	// 1.2遍历全部
	all = make([]string, 0, 3)
	New[string]().RPush("one", "two", "three").Range(func(s string) {
		all = append(all, s)
	})

	if !sliceEqual(all, []string{"one", "two", "three"}) {
		t.Errorf("Range() = %v, want %v", all, []string{"one", "two", "three"})
	}

	// 2.遍历部分
	all = make([]string, 0, 3)
	New[string]().RPush("one", "two", "three").Range(func(s string) {
		all = append(all, s)
	}, 0, -2)

	if !sliceEqual(all, []string{"one", "two"}) {
		t.Errorf("Range() = %v, want %v", all, []string{"one", "two"})
	}
}

func Test_PushBackList(t *testing.T) {
	if got := New[string]().RPush("one", "two", "three").PushBackList(New[string]().RPush("1", "2", "3")).ToSlice(); !sliceEqual(got, []string{"one", "two", "three", "1", "2", "3"}) {
		t.Errorf("PushBackList() = %v, want %v", got, []string{"one", "two", "three", "1", "2", "3"})
	}
}

func Test_PushFrontList(t *testing.T) {
	if got := New[string]().RPush("one", "two", "three").PushFrontList(New[string]().RPush("1", "2", "3")).ToSlice(); !sliceEqual(got, []string{"1", "2", "3", "one", "two", "three"}) {
		t.Errorf("PushFrontList() = %v, want %v", got, []string{"1", "2", "3", "one", "two", "three"})
	}
}

func Test_ContainsFunc(t *testing.T) {
	if got := New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "one" == v
	}); got != true {
		t.Errorf("ContainsFunc() = %v, want %v", got, true)
	}

	if got := New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "two" == v
	}); got != true {
		t.Errorf("ContainsFunc() = %v, want %v", got, true)
	}

	if got := New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "three" == v
	}); got != true {
		t.Errorf("ContainsFunc() = %v, want %v", got, true)
	}

	if got := New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "xx" == v
	}); got != false {
		t.Errorf("ContainsFunc() = %v, want %v", got, false)
	}
}

func Test_InsertAfter(t *testing.T) {
	if got := New[string]().RPush("one", "two").InsertAfter("three", func(s string) bool { return s == "two" }).ToSlice(); !sliceEqual(got, []string{"one", "two", "three"}) {
		t.Errorf("InsertAfter() = %v, want %v", got, []string{"one", "two", "three"})
	}
	if got := New[string]().RPush("one", "three").InsertAfter("two", func(s string) bool { return s == "one" }).ToSlice(); !sliceEqual(got, []string{"one", "two", "three"}) {
		t.Errorf("InsertAfter() = %v, want %v", got, []string{"one", "two", "three"})
	}
}

func Test_InsertBefore(t *testing.T) {
	if got := New[string]().RPush("one", "three").InsertBefore("two", func(s string) bool { return s == "three" }).ToSlice(); !sliceEqual(got, []string{"one", "two", "three"}) {
		t.Errorf("InsertBefore() = %v, want %v", got, []string{"one", "two", "three"})
	}
	if got := New[string]().RPush("two", "three").InsertBefore("one", func(s string) bool { return s == "two" }).ToSlice(); !sliceEqual(got, []string{"one", "two", "three"}) {
		t.Errorf("InsertBefore() = %v, want %v", got, []string{"one", "two", "three"})
	}
}

func Test_RemFunc(t *testing.T) {
	if got := New[int]().RPush(1, 1, 2, 2, 3, 3).RemFunc(2, func(v int) bool { return v == 1 }); got != 2 {
		t.Errorf("RemFunc() = %v, want %v", got, 2)
	}
	if got := New[int]().RPush(1, 1, 2, 2, 3, 3).RemFunc(-2, func(v int) bool { return v == 3 }); got != 2 {
		t.Errorf("RemFunc() = %v, want %v", got, 2)
	}
}

func Test_OtherMoveToBackList(t *testing.T) {
	// 第1个链表为空, 第2个链表有值
	l := New[int]()
	other := New[int]().PushBack(1, 2, 3, 4, 5, 6)
	l.OtherMoveToBackList(other)
	if !sliceEqual(l.ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("OtherMoveToBackList() = %v, want %v", l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	}
	if !sliceEqual(other.ToSlice(), []int(nil)) {
		t.Errorf("OtherMoveToBackList() = %v, want %v", other.ToSlice(), []int(nil))
	}

	// 第1个链表有值, 第2个链表也有值
	l = New[int]().PushBack(1, 2, 3)
	other = New[int]().PushBack(4, 5, 6)
	l.OtherMoveToBackList(other)
	if !sliceEqual(l.ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("OtherMoveToBackList() = %v, want %v", l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	}
	if !sliceEqual(other.ToSlice(), []int(nil)) {
		t.Errorf("OtherMoveToBackList() = %v, want %v", other.ToSlice(), []int(nil))
	}

	// 第1个链表有值, 第2个链表为空
	l = New[int]().PushBack(1, 2, 3, 4, 5, 6)
	other = New[int]()
	l.OtherMoveToBackList(other)
	if !sliceEqual(l.ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("OtherMoveToBackList() = %v, want %v", l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	}
	if !sliceEqual(other.ToSlice(), []int(nil)) {
		t.Errorf("OtherMoveToBackList() = %v, want %v", other.ToSlice(), []int(nil))
	}
}

func Test_OtherMoveToFrontList(t *testing.T) {
	// 第1个链表为空, 第2个链表有值
	l := New[int]()
	other := New[int]().PushBack(1, 2, 3, 4, 5, 6)
	l.OtherMoveToFrontList(other)
	if !sliceEqual(l.ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("OtherMoveToFrontList() = %v, want %v", l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	}
	if !sliceEqual(other.ToSlice(), []int(nil)) {
		t.Errorf("OtherMoveToFrontList() = %v, want %v", other.ToSlice(), []int(nil))
	}

	// 第1个链表有值, 第2个链表也有值
	l = New[int]().PushBack(1, 2, 3)
	other = New[int]().PushBack(4, 5, 6)
	l.OtherMoveToFrontList(other)
	if !sliceEqual(l.ToSlice(), []int{4, 5, 6, 1, 2, 3}) {
		t.Errorf("OtherMoveToFrontList() = %v, want %v", l.ToSlice(), []int{4, 5, 6, 1, 2, 3})
	}
	if !sliceEqual(other.ToSlice(), []int(nil)) {
		t.Errorf("OtherMoveToFrontList() = %v, want %v", other.ToSlice(), []int(nil))
	}

	// 第1个链表有值, 第2个链表为空
	l = New[int]().PushBack(1, 2, 3, 4, 5, 6)
	other = New[int]()
	l.OtherMoveToFrontList(other)
	if !sliceEqual(l.ToSlice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("OtherMoveToFrontList() = %v, want %v", l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	}
	if !sliceEqual(other.ToSlice(), []int(nil)) {
		t.Errorf("OtherMoveToFrontList() = %v, want %v", other.ToSlice(), []int(nil))
	}
}

// Helper function to compare slices
func sliceEqual[T comparable](a, b []T) bool {
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
