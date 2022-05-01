package rhash

import (
	"github.com/cespare/xxhash/v2"
	"reflect"
	"unsafe"
)

// 元素
type entry[K comparable, V any] struct {
	key  K
	val  V
	next *entry[K, V]
}

// hash 表头
type Hash[K comparable, V any] struct {
	table     [2][]entry[K, V]
	htUsed    [2]int // 记录每个table里面存在的元素个数
	rehashidx int
	keySize   int //key的长度
	hashFunc  func(str string) uint64
	isKeyStr  bool //是string类型的key, 或者不是
}

// 初始化一个hashtable
func New[K comparable, V any]() *Hash[K, V] {
	return &Hash[K, V]{
		rehashidx: -1,
		hashFunc:  xxhash.Sum64String,
	}
}

// 初始化一个hashtable并且可以设置值
func NewWithHashFunc[K comparable, V any](hashFunc func(str string) uint64) *Hash[K, V] {
	h := New[K, V]()
	h.hashFunc = hashFunc
	return h
}

// 保存key的类型和key的长度
func (h *Hash[K, V]) keyTypeAndKeySize() {
	var k K
	switch (interface{})(k).(type) {
	case string:
	default:
		h.keySize = int(unsafe.Sizeof(k))
	}
}

// 计算hash值
func (h *Hash[K, V]) calHash(k K) uint64 {
	var key string

	if h.isKeyStr {
		// 直接赋值会报错, 使用unsafe绕过编译器检查
		key = *(*string)(unsafe.Pointer(&k))
	} else {
		key = *(*string)(unsafe.Pointer(&reflect.StringHeader{
			Data: uintptr(unsafe.Pointer(&k)),
			Len:  h.keySize,
		}))
	}

	return xxhash.Sum64String(key)
}

func (h *Hash[K, V]) isRehashing() bool {
	return h.rehashidx != -1
}

func (h *Hash[K, V]) Set(k K, v V) {
}

func (h *Hash[K, V]) Delete(key K) {

}

func (h *Hash[K, V]) findIndexAndEntry() (i int, e *entry[K, V]) {

	return
}

func (h *Hash[K, V]) Get(key K) (v V, err error) {

	return
}

func (h *Hash[K, V]) GetOrZero(key K) (v V) {
	return
}

func (h *Hash[K, V]) Range() {

}
