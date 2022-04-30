package rhash

import (
	//"github.com/cespare/xxhash/v2"
	"unsafe"
)

type entry[K comparable, V any] struct {
	key  K
	val  V
	next *entry[K, V]
}

type Hash[K comparable, V any] struct {
	table     [2][]entry[K, V]
	htUsed    [2]int
	rehashidx int

	isKeyStr bool
	keySize  int
}

func New[K comparable, V any]() *Hash[K, V] {
	return &Hash[K, V]{rehashidx: -1}
}

func (h *Hash[K, V]) keyTypeAndKeySize() {
	var k K
	switch (interface{})(k).(type) {
	case string:
	default:
		h.keySize = int(unsafe.Sizeof(k))
	}
}

func (h *Hash[K, V]) isRehashing() bool {
	return h.rehashidx != -1
}

func (h *Hash[K, V]) Set(k K, v V) {

}

func (h *Hash[K, V]) Delete(key K) {

}

func (h *Hash[K, V]) Get(key K) (v V, err error) {

	return
}

func (h *Hash[K, V]) GetOrZero(key K) (v V) {
	return
}

func (h *Hash[K, V]) Range() {

}
