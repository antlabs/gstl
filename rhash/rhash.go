package rhash

import (
	"github.com/cespare/xxhash/v2"
	"math"
	"reflect"
	"unsafe"
)

const (
	HT_INITIAL_EXP = 2
)

// 元素
type entry[K comparable, V any] struct {
	key  K
	val  V
	next *entry[K, V]
}

// hash 表头
type Hash[K comparable, V any] struct {
	table   [2][]entry[K, V] //hash table
	used    [2]uint64        // 记录每个table里面存在的元素个数
	sizeExp [2]uint64        //记录exp

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
		// 因为xxhash.Sum64String 接收string, 所以要把非string类型变量当成string类型来处理
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

// TODO 这个函数可以优化下
func nextExp(size uint64) uint64 {
	if size >= math.MaxUint64 {
		return 63
	}
	e := uint64(HT_INITIAL_EXP)
	for {
		if 1<<e >= size {
			return e
		}
		e++
	}
	return e
}

// 手动扩容
func (h *Hash[K, V]) Expand(size int) {
	// 如果正在扩容中, 或者需要扩容的数据小于已存在的元素, 直接返回
	if h.isRehashing() || h.used[0] > uint64(size) {
		return
	}

	newSizeExp := nextExp(uint64(size))
	// 新大小比需要的大小还小
	newSize := 1 << newSizeExp
	if newSize < size {
		return
	}

	// 新扩容大小和以前的一样
	if uint64(newSizeExp) == h.sizeExp[0] {
		return
	}

	newTable := make([]entry[K, V], newSize)

	// 第一次初始化
	if h.table[0] == nil {
		h.sizeExp[0] = newSizeExp
		h.table[0] = newTable
		return
	}

	// 把新hash表放到table[1]里面
	h.sizeExp[1] = newSizeExp
	h.used[1] = 0
	h.table[1] = newTable
	h.rehashidx = 0
	return
}

// 返回索引值和entry
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
