package linkedlist

// apache 2.0 antlabs
import (
	"testing"

	"github.com/antlabs/gstl/must"
	"github.com/stretchr/testify/assert"
)

func Test_Push(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").ToSlice(), []string{"1", "2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).ToSlice(), []int{1, 2, 3, 4})
}

func Test_RPush(t *testing.T) {
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4").ToSlice(), []string{"1", "2", "3", "4"})
	assert.Equal(t, New[int]().RPush(1, 2, 3, 4).ToSlice(), []int{1, 2, 3, 4})
}

func Test_Len(t *testing.T) {
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").Len(), 4)
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).Len(), 4)
}

func Test_RPop(t *testing.T) {
	// 不正常的索引
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(-3), []string(nil))

	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(3), []string{"2", "3", "4"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).RPop(3), []int{2, 3, 4})
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").RPop(10), []string{"1", "2", "3", "4"})
}

func Test_LPop(t *testing.T) {
	// 不正常的索引
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(-3), []string(nil))

	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(3), []string{"1", "2", "3"})
	assert.Equal(t, New[int]().PushBack(1, 2, 3, 4).LPop(3), []int{1, 2, 3})
	assert.Equal(t, New[string]().PushBack("1", "2", "3", "4").LPop(10), []string{"1", "2", "3", "4"})
}

func Test_RangeSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangeSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	assert.Equal(t, all, []string{"1", "2", "3", "4"})
}

func Test_RangePrevSafe(t *testing.T) {
	var all []string
	New[string]().PushBack("1", "2", "3", "4").RangePrevSafe(func(n *Node[string]) bool {
		all = append(all, n.Element)
		return false
	})

	assert.Equal(t, all, []string{"4", "3", "2", "1"})
}

func Test_First(t *testing.T) {
	// 没有值
	assert.Error(t, must.TakeOneErr(New[string]().First()))
	// 有值
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").First()), "1")
	assert.Equal(t, must.TakeOne(New[int]().RPush(1, 2).First()), 1)
}

func Test_Last(t *testing.T) {
	// 没有值
	assert.Error(t, must.TakeOneErr(New[string]().Last()))
	//有值
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").Last()), "4")
	assert.Equal(t, must.TakeOne(New[int]().RPush(1, 2).Last()), 2)
}

func Test_Get(t *testing.T) {
	// 正索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(0)), "1")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(1)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(2)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(3)), "4")

	// 负索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(-1)), "4")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(-2)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(-3)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4").GetWithErr(-4)), "1")
}

func Test_Set(t *testing.T) {
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(11, "1991").ToSlice(), []string{"1", "2", "3", "4", "5", "6"})

	// 正索引
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(0, "1991").ToSlice(), []string{"1991", "2", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(1, "1991").ToSlice(), []string{"1", "1991", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(2, "1991").ToSlice(), []string{"1", "2", "1991", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(3, "1991").ToSlice(), []string{"1", "2", "3", "1991", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(4, "1991").ToSlice(), []string{"1", "2", "3", "4", "1991", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Set(5, "1991").ToSlice(), []string{"1", "2", "3", "4", "5", "1991"})

	// 负索引
}

func Test_Index(t *testing.T) {
	assert.Error(t, must.TakeOneErr(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(11)))

	// 正索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(0)), "1")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(1)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(2)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(3)), "4")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(4)), "5")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(5)), "6")

	// 负索引
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-1)), "6")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-2)), "5")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-3)), "4")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-4)), "3")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-5)), "2")
	assert.Equal(t, must.TakeOne(New[string]().RPush("1", "2", "3", "4", "5", "6").Index(-6)), "1")
}

func Test_Remove(t *testing.T) {
	// 正索引
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(0).ToSlice(), []string{"2", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(1).ToSlice(), []string{"1", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(2).ToSlice(), []string{"1", "2", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(3).ToSlice(), []string{"1", "2", "3", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(4).ToSlice(), []string{"1", "2", "3", "4", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(5).ToSlice(), []string{"1", "2", "3", "4", "5"})

	// 负索引
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-1).ToSlice(), []string{"1", "2", "3", "4", "5"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-2).ToSlice(), []string{"1", "2", "3", "4", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-3).ToSlice(), []string{"1", "2", "3", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-4).ToSlice(), []string{"1", "2", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-5).ToSlice(), []string{"1", "3", "4", "5", "6"})
	assert.Equal(t, New[string]().RPush("1", "2", "3", "4", "5", "6").Remove(-6).ToSlice(), []string{"2", "3", "4", "5", "6"})
}

func Test_LPush(t *testing.T) {
	assert.Equal(t, New[string]().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().LPush("1").ToSlice(), []string{"1"})
	assert.Equal(t, New[string]().LPush("2", "1").ToSlice(), []string{"1", "2"})
	assert.Equal(t, New[string]().LPush("3", "2", "1").ToSlice(), []string{"1", "2", "3"})
	assert.Equal(t, New[string]().LPush("4", "3", "2", "1").ToSlice(), []string{"1", "2", "3", "4"})
}

func Test_PushFront(t *testing.T) {
	assert.Equal(t, New[string]().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("1").ToSlice(), []string{"1"})
	assert.Equal(t, New[string]().PushFront("2", "1").ToSlice(), []string{"1", "2"})
	assert.Equal(t, New[string]().PushFront("3", "2", "1").ToSlice(), []string{"1", "2", "3"})
	assert.Equal(t, New[string]().PushFront("4", "3", "2", "1").ToSlice(), []string{"1", "2", "3", "4"})
}

func Test_Clear(t *testing.T) {
	assert.Equal(t, New[string]().PushFront("1").Clear().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("2", "1").Clear().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("3", "2", "1").Clear().ToSlice(), []string(nil))
	assert.Equal(t, New[string]().PushFront("4", "3", "2", "1").Clear().ToSlice(), []string(nil))
}

func Test_IsEmpty(t *testing.T) {
	assert.True(t, New[string]().IsEmpty())
	assert.True(t, New[string]().PushFront("1").Clear().IsEmpty())
	assert.True(t, New[string]().PushFront("2", "1").Clear().IsEmpty())
	assert.True(t, New[string]().PushFront("3", "2", "1").Clear().IsEmpty())
	assert.True(t, New[string]().PushFront("4", "3", "2", "1").Clear().IsEmpty())
}

func Test_Trim(t *testing.T) {
	// 返回空值的情况
	assert.Equal(t, New[string]().RPush("one", "two", "three").Trim(100, -2).ToSlice(), []string{"one", "two", "three"})

	assert.Equal(t, New[string]().RPush("one", "two", "three").Trim(1, -1).ToSlice(), []string{"two", "three"})
	assert.Equal(t, New[string]().RPush("one", "two", "three").Trim(0, 1).ToSlice(), []string{"one", "two"})
	assert.Equal(t, New[string]().RPush("one", "two", "three").Trim(-2, -1).ToSlice(), []string{"two", "three"})
	assert.Equal(t, New[string]().RPush("one", "two", "three").Trim(-100, -2).ToSlice(), []string{"one", "two"})
}

func Test_Range(t *testing.T) {
	var all []string

	// 1.1遍历全部
	all = make([]string, 0, 3)
	New[string]().RPush("one", "two", "three").Range(func(s string) {
		all = append(all, s)
	}, 0, -1)

	assert.Equal(t, all, []string{"one", "two", "three"})

	// 1.2遍历全部
	all = make([]string, 0, 3)
	New[string]().RPush("one", "two", "three").Range(func(s string) {
		all = append(all, s)
	})

	assert.Equal(t, all, []string{"one", "two", "three"})

	// 2.遍历部分
	all = make([]string, 0, 3)
	New[string]().RPush("one", "two", "three").Range(func(s string) {
		all = append(all, s)
	}, 0, -2)

	assert.Equal(t, all, []string{"one", "two"})
}

func Test_PushBackList(t *testing.T) {
	assert.Equal(t, New[string]().RPush("one", "two", "three").PushBackList(New[string]().RPush("1", "2", "3")).ToSlice(), []string{"one", "two", "three", "1", "2", "3"})
}

func Test_PushFrontList(t *testing.T) {
	assert.Equal(t, New[string]().RPush("one", "two", "three").PushFrontList(New[string]().RPush("1", "2", "3")).ToSlice(), []string{"1", "2", "3", "one", "two", "three"})
}

func Test_ContainsFunc(t *testing.T) {
	assert.True(t, New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "one" == v
	}))

	assert.True(t, New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "two" == v
	}))

	assert.True(t, New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "three" == v
	}))

	assert.False(t, New[string]().RPush("one", "two", "three").ContainsFunc(func(v string) bool {
		return "xx" == v
	}))
}

func Test_InsertAfter(t *testing.T) {

	assert.Equal(t, New[string]().RPush("one", "two").InsertAfter("three", func(s string) bool { return s == "two" }).ToSlice(), []string{"one", "two", "three"})
	assert.Equal(t, New[string]().RPush("one", "three").InsertAfter("two", func(s string) bool { return s == "one" }).ToSlice(), []string{"one", "two", "three"})
}

func Test_InsertBefore(t *testing.T) {

	assert.Equal(t, New[string]().RPush("one", "three").InsertBefore("two", func(s string) bool { return s == "three" }).ToSlice(), []string{"one", "two", "three"})
	assert.Equal(t, New[string]().RPush("two", "three").InsertBefore("one", func(s string) bool { return s == "two" }).ToSlice(), []string{"one", "two", "three"})
}

func Test_RemFunc(t *testing.T) {

	assert.Equal(t, New[int]().RPush(1, 1, 2, 2, 3, 3).RemFunc(2, func(v int) bool { return v == 1 }), 2)
	assert.Equal(t, New[int]().RPush(1, 1, 2, 2, 3, 3).RemFunc(-2, func(v int) bool { return v == 3 }), 2)
}

func Test_OtherMoveToBackList(t *testing.T) {
	// 第1个链表为空, 第2个链表有值
	l := New[int]()
	other := New[int]().PushBack(1, 2, 3, 4, 5, 6)
	l.OtherMoveToBackList(other)
	assert.Equal(t, l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, other.ToSlice(), []int(nil))

	// 第1个链表有值, 第2个链表也有值
	l = New[int]().PushBack(1, 2, 3)
	other = New[int]().PushBack(4, 5, 6)
	l.OtherMoveToBackList(other)
	assert.Equal(t, l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, other.ToSlice(), []int(nil))

	// 第1个链表有值, 第2个链表为空
	l = New[int]().PushBack(1, 2, 3, 4, 5, 6)
	other = New[int]()
	l.OtherMoveToBackList(other)
	assert.Equal(t, l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, other.ToSlice(), []int(nil))
}

func Test_OtherMoveToFrontList(t *testing.T) {
	// 第1个链表为空, 第2个链表有值
	l := New[int]()
	other := New[int]().PushBack(1, 2, 3, 4, 5, 6)
	l.OtherMoveToFrontList(other)
	assert.Equal(t, l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, other.ToSlice(), []int(nil))

	// 第1个链表有值, 第2个链表也有值
	l = New[int]().PushBack(1, 2, 3)
	other = New[int]().PushBack(4, 5, 6)
	l.OtherMoveToFrontList(other)
	assert.Equal(t, l.ToSlice(), []int{4, 5, 6, 1, 2, 3})
	assert.Equal(t, other.ToSlice(), []int(nil))

	// 第1个链表有值, 第2个链表为空
	l = New[int]().PushBack(1, 2, 3, 4, 5, 6)
	other = New[int]()
	l.OtherMoveToFrontList(other)
	assert.Equal(t, l.ToSlice(), []int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, other.ToSlice(), []int(nil))
}
