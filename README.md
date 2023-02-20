# gstl
支持泛型的数据结构库   
[![Go](https://github.com/antlabs/gstl/workflows/Go/badge.svg)](https://github.com/antlabs/gstl/actions)
[![codecov](https://codecov.io/gh/antlabs/gstl/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/gstl)

## 一、`vec`
```go

```
## 二、`Listked`
```go

```

## 三、`rhashmap`
和标准库不同的地方是有序hash
```go
```

## 四、`btree`
```go
```
## 五、`rbtree`
```go
```

## 六、`avltree`
```go
```

## 七、`trie`
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

## 八、`set`
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

## 九、`ifop`
ifop是弥补下golang没有三目运算符，使用函数模拟

```go
// 如果该值不为0, 返回原来的值，否则默认值
val = IfElse(len(val) != 0, val, "default")
```
## 十、`mapex`
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
## 十一、`rwmap`
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
## 十二、`cmap`
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
```