# gstl
支持泛型的数据结构库
## 一、`vec`
```go

```
## 二、`Listked`
```go

```

## 三、`rhashmap`
和标准库不同的地址是有序hash
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
