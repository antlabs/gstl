# gstl
支持泛型的数据结构库   
[![Go](https://github.com/antlabs/gstl/workflows/Go/badge.svg)](https://github.com/antlabs/gstl/actions)
[![codecov](https://codecov.io/gh/antlabs/gstl/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/gstl)

## 一、`vec`
```go

```
## 二、`Listked`

`Listked` 是一个支持泛型的双向链表容器，提供了加锁和不加锁的实现。

#### 不加锁的使用方式

```go
package main

import (
    "fmt"
    "github.com/antlabs/gstl/linkedlist"
)

func main() {
    // 创建一个不加锁的链表
    list := linkedlist.New[int]()

    // 插入元素
    list.PushBack(1)
    list.PushFront(0)

    // 遍历链表
    list.Range(func(value int) {
        fmt.Println(value)
    })

    // 删除元素
    list.Remove(0)
}
```

#### 加锁的使用方式

```go
package main

import (
    "fmt"
    "sync"
    "github.com/antlabs/gstl/linkedlist"
)

func main() {
    // 创建一个加锁的链表
    list := linkedlist.NewConcurrent[int]()

    var wg sync.WaitGroup
    wg.Add(2)

    // 并发插入元素
    go func() {
        defer wg.Done()
        list.PushBack(1)
        list.PushFront(0)
    }()

    // 并发遍历链表
    go func() {
        defer wg.Done()
        list.Range(func(value int) {
            fmt.Println(value)
        })
    }()

    wg.Wait()

    // 删除元素
    list.Remove(0)
}
```

### 区别

- **不加锁的链表**：适用于单线程环境，性能更高。
- **加锁的链表**：适用于多线程环境，保证线程安全。

## 三、`rhashmap`
和标准库不同的地方是有序hash
```go
```

## 四、`btree`
```go
```
## 五、`SkipList`

`SkipList` 是一种高效的有序数据结构，支持快速的插入、删除和查找操作。

#### 基本使用

```go
package main

import (
    "fmt"
    "github.com/antlabs/gstl/skiplist"
)

func main() {
    // 创建一个新的 SkipList
    sl := skiplist.New[int, string]()

    // 插入元素
    sl.Insert(1, "one")
    sl.Insert(2, "two")

    // 获取元素
    if value, ok := sl.Get(1); ok {
        fmt.Println("Key 1:", value)
    }

    // 删除元素
    sl.Delete(1)
}
```

#### 并发安全的使用

```go
package main

import (
    "fmt"
    "sync"
    "github.com/antlabs/gstl/skiplist"
)

func main() {
    // 创建一个并发安全的 SkipList
    csl := skiplist.NewConcurrent[int, string]()
    var wg sync.WaitGroup

    // 并发插入元素
    wg.Add(2)
    go func() {
        defer wg.Done()
        csl.Insert(1, "one")
    }()
    go func() {
        defer wg.Done()
        csl.Insert(2, "two")
    }()
    wg.Wait()

    // 并发获取元素
    wg.Add(2)
    go func() {
        defer wg.Done()
        if value, ok := csl.Get(1); ok {
            fmt.Println("Key 1:", value)
        }
    }()
    go func() {
        defer wg.Done()
        if value, ok := csl.Get(2); ok {
            fmt.Println("Key 2:", value)
        }
    }()
    wg.Wait()
}
```

## 六、`rbtree`
```go
```

## 七、`avltree`
```go
```

## 八、`trie`
```go
// 声明一个bool类型的trie tree
t := trie.New[bool]()

// 新增一个key
t.Set("hello", true)

// 获取值
v := t.Get("hello")

// 检查trie中是有hello前缀的数据
ok := t.HasPrefix("hello") 

// 删除键
t.Delete(k string)

// 返回trie中保存的元素个数
t.Len()
```

## 九、`set`
```go
// 声明一个string类型的set
s := set.New[string]()

// 新加成员
s.Set("1")
s.Set("2")
s.Set("3")

// 查看某个变量是否存在set中
s.IsMember(1)

// 长度
s.Len()

// set转slice
s.ToSlice()

// 深度复制一份
newSet := s.Close()

// 集合取差集 s - s2
s := From("hello", "world", "1234", "4567")
s2 := From("1234", "4567")

newSet := s.Diff(s2)
assert.Equal(t, newSet.ToSlice(), []string{"hello", "world"})

// 集合取交集
s := From("1234", "5678", "9abc")
s2 := From("abcde", "5678", "9abc")

v := s.Intersection(s2).ToSlice()
assert.Equal(t, v, []string{"5678", "9abc"})

// 集合取并集
s := From("1111")
s1 := From("2222")
s2 := From("3333")

newSet := s.Union(s1, s2)
assert.Equal(t, newSet.ToSlice(), []string{"1111", "2222", "3333"})

// 测试集合s每个元素是否在s1里面, s <= s1
s := From("5678", "9abc")
s2 := From("abcde", "5678", "9abc")

assert.True(t, s.IsSubset(s2))

// 测试集合s1每个元素是否在s里面 s1 <= s
s2 := From("5678", "9abc")
s := From("abcde", "5678", "9abc")

assert.True(t, s.IsSuperset(s2))

// 遍历某个集合
a := []string{"1111", "2222", "3333"}
s := From(a...)
for _, v := range a {
  s.Set(v)
}

s.Range(func(k string) bool {
  fmt.Println(k)
  return true
})

// 测试两个集合是否相等 Equal
s := New[int]()
max := 1000
for i := 0; i < max; i++ {
  s.Set(i)
}

s2 := s.Clone()

assert.True(t, s.Equal(s2))
```

## 十、`ifop`
ifop是弥补下golang没有三目运算符，使用函数模拟
### 10.1 if else部分类型相同
```go
// 如果该值不为0, 返回原来的值，否则默认值
val = IfElse(len(val) != 0, val, "default")
```
### 10.2 if else部分类型不同
```go
o := map[string]any{"hello": "hello"}
a := []any{"hello", "world"}
fmt.Printf("%#v", IfElseAny(o != nil, o, a))
```
## 十一、`mapex`
薄薄一层包装，增加标准库map的接口
* mapex.Keys()
```go
m := make(map[string]string)
m["a"] = "1"
m["b"] = "2"
m["c"] = "3"
get := mapex.Keys(m)// 返回map的所有key

```
* mapex.Values()
```go
m := make(map[string]string)
m["a"] = "1"
m["b"] = "2"
m["c"] = "3"
get := mapex.Values(m)
```
## 十二、`rwmap`
rwmap与sync.Map类似支持并发访问，只解决sync.Map 2个问题.  
1. 没有Len成员函数  
2. 以及没有使用泛型语法，有运行才发现类型使用错误的烦恼
```go
var m rwmap.RWMap[string, string] // 声明一个string, string的map
m.Store("hello", "1") // 保存
v1, ok1 := m.Load("hello") // 获取值
v1, ok1 = m.LoadAndDelete("hello") //返回hello对应值，然后删除hello
Delete("hello") // 删除
v1, ok1 = m.LoadOrStore("hello", "world")

// 遍历，使用回调函数
m.Range(func(key, val string) bool {
	fmt.Printf("k:%s, val:%s\n"i, key, val)
	return true
})

// 遍历，迭代器
for pair := range m.Iter() {
  fmt.Printf("k:%s, val:%s\n", pair.Key, pair.Val)
}

m.Len()// 获取长度
allKeys := m.Keys() //返回所有的key
allValues := m.Values()// 返回所有的value
```
## 十三、`cmap`
cmap是用锁分区的方式实现的，(TODO优化，目前只有几个指标比sync.Map快)
```go
var m cmap.CMap[string, string] // 声明一个string, string的map
m.Store("hello", "1") // 保存
v1, ok1 := m.Load("hello") // 获取值
v1, ok1 = m.LoadAndDelete("hello") //返回hello对应值，然后删除hello
Delete("hello") // 删除
v1, ok1 = m.LoadOrStore("hello", "world")

// 遍历，使用回调函数
m.Range(func(key, val string) bool {
	fmt.Printf("k:%s, val:%s\n"i, key, val)
	return true
})

// 遍历，迭代器
for pair := range m.Iter() {
  fmt.Printf("k:%s, val:%s\n", pair.Key, pair.Val)
}

m.Len()// 获取长度
allKeys := m.Keys() //返回所有的key
allValues := m.Values()// 返回所有的value
