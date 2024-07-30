package cmap

import (
	"reflect"
	"runtime"
	"sync"
	"unsafe"

	"github.com/antlabs/gstl/api"
	xxhash "github.com/cespare/xxhash/v2"
	"golang.org/x/exp/constraints"
)

var _ api.CMaper[int, int] = (*CMap[int, int])(nil)

type Pair[K constraints.Ordered, V any] struct {
	Key K
	Val V
}

type CMap[K constraints.Ordered, V any] struct {
	bucket   []Item[K, V]
	keySize  int
	isKeyStr bool
}

type Item[K constraints.Ordered, V any] struct {
	rw sync.RWMutex
	m  api.Map[K, V]
}

func New[K constraints.Ordered, V any]() (c *CMap[K, V]) {
	c = &CMap[K, V]{}
	c.init(0)
	return c
}

func (c *CMap[K, V]) init(n int) {
	np := runtime.GOMAXPROCS(0)
	if np <= 0 {
		np = 8
	}

	if n > 0 {
		np = n
	}

	c.bucket = make([]Item[K, V], np)

	for i := range c.bucket {
		c.bucket[i].m = newStdMap[K, V]()
	}

}

// 计算hash值
func (c *CMap[K, V]) calHash(k K) uint64 {
	var key string

	if c.isKeyStr {
		// 直接赋值会报错, 使用unsafe绕过编译器检查
		key = *(*string)(unsafe.Pointer(&k))
	} else {
		// 因为xxhash.Sum64String 接收string, 所以要把非string类型变量当成string类型来处理
		key = *(*string)(unsafe.Pointer(&reflect.StringHeader{
			Data: uintptr(unsafe.Pointer(&k)),
			Len:  c.keySize,
		}))
	}

	return xxhash.Sum64String(key)
}

// 保存key的类型和key的长度
func (h *CMap[K, V]) keyTypeAndKeySize() {
	var k K
	switch (interface{})(k).(type) {
	case string:
		h.isKeyStr = true
	default:
		h.keySize = int(unsafe.Sizeof(k))
	}
}

// 找到索引
func (c *CMap[K, V]) findIndex(key K) *Item[K, V] {
	index := c.calHash(key) % uint64(len(c.bucket))
	return &c.bucket[index]
}

// 删除
func (c *CMap[K, V]) Delete(key K) {
	item := c.findIndex(key)
	item.rw.Lock()
	item.m.Delete(key)
	item.rw.Unlock()
}

type UpdataOrInsertCb[K constraints.Ordered, V any] func(exist bool, old V) (newVal V)

// 删除或者更新
func (c *CMap[K, V]) UpdateOrInsert(k K, cb UpdataOrInsertCb[K, V]) {
	item := c.findIndex(k)
	item.rw.Lock()
	old, ok := item.m.GetWithBool(k)
	newVal := cb(ok, old)
	item.m.Set(k, newVal)
	item.rw.Unlock()

}

func (c *CMap[K, V]) Load(key K) (value V, ok bool) {
	item := c.findIndex(key)
	item.rw.RLock()
	value, ok = item.m.GetWithBool(key)
	item.rw.RUnlock()
	return
}

func (c *CMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	item := c.findIndex(key)
	item.rw.Lock()
	value, loaded = item.m.GetWithBool(key)
	if !loaded {
		item.rw.Unlock()
		return
	}
	item.m.Delete(key)
	item.rw.Unlock()
	return
}

func (c *CMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	item := c.findIndex(key)
	item.rw.Lock()
	actual, loaded = item.m.GetWithBool(key)
	if !loaded {
		actual = value
		item.m.Set(key, actual)
		item.rw.Unlock()
		return
	}

	actual, loaded = item.m.GetWithBool(key)
	item.rw.Unlock()
	return
}

func (c *CMap[K, V]) Range(f func(key K, value V) bool) {
	for i := 0; i < len(c.bucket); i++ {
		item := &c.bucket[i]
		item.rw.RLock()
		item.m.Range(f)
		item.rw.RUnlock()
	}
}

func (c *CMap[K, V]) Iter() (rv chan Pair[K, V]) {

	rv = make(chan Pair[K, V])
	var wg sync.WaitGroup

	wg.Add(len(c.bucket))

	go func() {
		wg.Wait()
		close(rv)
	}()

	for i := 0; i < len(c.bucket); i++ {
		item := &c.bucket[i]
		go func(item *Item[K, V]) {

			defer wg.Done()
			item.rw.RLock()
			item.m.Range(func(key K, value V) bool {
				rv <- Pair[K, V]{Key: key, Val: value}
				return true
			})
			item.rw.RUnlock()

		}(item)
	}
	return rv

}

func (c *CMap[K, V]) Store(key K, value V) {
	item := c.findIndex(key)
	item.rw.Lock()
	item.m.Set(key, value)
	item.rw.Unlock()
	return
}

// TODO 优化
func (c *CMap[K, V]) Keys() []K {
	l := c.Len()
	all := make([]K, 0, l)
	if l == 0 {
		return nil
	}

	for i := 0; i < len(c.bucket); i++ {

		item := &c.bucket[i]
		item.rw.RLock()
		item.m.Range(func(key K, value V) bool {
			all = append(all, key)
			return true
		})
		item.rw.RUnlock()
	}
	return all
}

func (c *CMap[K, V]) Values() []V {
	l := c.Len()
	all := make([]V, 0, l)
	if l == 0 {
		return nil
	}

	for i := 0; i < len(c.bucket); i++ {

		item := &c.bucket[i]
		item.rw.RLock()
		item.m.Range(func(key K, value V) bool {
			all = append(all, value)
			return true
		})
		item.rw.RUnlock()
	}
	return all
}

func (c *CMap[K, V]) Len() int {
	l := 0
	for i := 0; i < len(c.bucket); i++ {
		item := &c.bucket[i]
		item.rw.RLock()
		l += item.m.Len()
		item.rw.RUnlock()
	}
	return l
}
