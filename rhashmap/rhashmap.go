package rhashmap

// 参考资料
// https://github.com/redis/redis/blob/unstable/src/dict.c
import (
	"errors"
	"github.com/cespare/xxhash/v2"
	"math"
	"reflect"
	"unsafe"
)

const (
	HT_INITIAL_EXP  = 2
	HT_INITIAL_SIZE = (1 << (HT_INITIAL_EXP))
)

var forceResizeRatio = 5

var (
	ErrHashing  = errors.New("rehashing...")
	ErrSize     = errors.New("wrong size")
	ErrNotFound = errors.New("not found")
)

// 元素
type entry[K comparable, V any] struct {
	key  K
	val  V
	next *entry[K, V]
}

type config struct {
	hashFunc func(str string) uint64
	cap      int
}

// hash 表头
type HashMap[K comparable, V any] struct {
	// 大多数情况, table[0]里就存在hash表元素的数据
	// 大小一尘不变hash随着数据的增强效率会降低, rhashmap的实现是超过某阈值时
	// table[1] 会先放新申请的hash表元素, 当table[0]都移动到table[1]时, table[1]赋值给table[0], 完成一次hash扩容
	// 移动的操作都分摊到Get, Set, Delete操作中, 每次移动一个槽位, 或者跳运100个空桶(TODO修改代码, 需要修改这边的注释)
	table   [2][]*entry[K, V] //hash table
	used    [2]uint64         // 记录每个table里面存在的元素个数
	sizeExp [2]int8           //记录exp

	rehashidx int // rehashid目前的槽位
	keySize   int //key的长度
	config
	isKeyStr bool //是string类型的key, 或者不是
	init     bool
}

// 初始化一个hashtable
func New[K comparable, V any]() *HashMap[K, V] {
	h := &HashMap[K, V]{}
	h.Init()
	return h
}

func (h *HashMap[K, V]) Init() {

	h.rehashidx = -1
	h.hashFunc = xxhash.Sum64String
	h.init = true

	h.reset(0)
	h.reset(1)
	h.keyTypeAndKeySize()
}

func (h *HashMap[K, V]) lazyinit() {
	if !h.init {
		h.Init()
	}
}

// 初始化一个hashtable并且可以设置值
func NewWithOpt[K comparable, V any](opts ...Option) *HashMap[K, V] {
	h := New[K, V]()
	for _, o := range opts {
		o.apply(&h.config)
	}

	if h.cap > 0 {
		h.Resize(uint64(h.cap))
	}
	return h
}

// 保存key的类型和key的长度
func (h *HashMap[K, V]) keyTypeAndKeySize() {
	var k K
	switch (interface{})(k).(type) {
	case string:
		h.isKeyStr = true
	default:
		h.keySize = int(unsafe.Sizeof(k))
	}
}

// 计算hash值
func (h *HashMap[K, V]) calHash(k K) uint64 {
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

func (h *HashMap[K, V]) isRehashing() bool {
	return h.rehashidx != -1
}

// TODO 这个函数可以优化下
func nextExp(size uint64) int8 {
	if size >= math.MaxUint64 {
		return 63
	}

	e := int8(HT_INITIAL_EXP)
	for {
		if 1<<e >= size {
			return e
		}
		e++
	}

	return e
}

func (h *HashMap[K, V]) expand() error {
	if h.isRehashing() {
		return nil
	}

	if hashSize(h.sizeExp[0]) == 0 {
		return h.Resize(HT_INITIAL_SIZE)
	}

	if h.used[0] >= hashSize(h.sizeExp[0]) || h.used[0]/hashSize(h.sizeExp[0]) > uint64(forceResizeRatio) {
		return h.Resize(h.used[0] + 1)
	}

	return nil
}

// 手动修改hashtable的大小
func (h *HashMap[K, V]) Resize(size uint64) error {
	h.lazyinit()
	// 如果正在扩容中, 或者需要扩容的数据小于已存在的元素, 直接返回
	if h.isRehashing() || h.used[0] > uint64(size) {
		return ErrHashing
	}

	newSizeExp := nextExp(uint64(size))
	// 新大小比需要的大小还小
	newSize := uint64(1 << newSizeExp)
	if newSize < size {
		return ErrSize
	}

	// 新扩容大小和以前的一样
	if uint64(newSizeExp) == uint64(h.sizeExp[0]) {
		return nil
	}

	newTable := make([]*entry[K, V], newSize)

	// 第一次初始化
	if h.table[0] == nil {
		h.sizeExp[0] = newSizeExp
		h.table[0] = newTable
		return nil
	}

	// 把新hash表放到table[1]里面
	h.sizeExp[1] = newSizeExp
	h.used[1] = 0
	h.table[1] = newTable
	h.rehashidx = 0
	return nil
}

// 收缩hash table
func (h *HashMap[K, V]) ShrinkToFit() error {
	h.lazyinit()
	if h.isRehashing() {
		return ErrHashing
	}

	minimal := h.used[0]
	if minimal < HT_INITIAL_SIZE {
		minimal = HT_INITIAL_SIZE
	}

	return h.Resize(minimal)
}

// 返回索引值和entry
func (h *HashMap[K, V]) findIndexAndEntry(key K) (i uint64, e *entry[K, V], err error) {
	if err := h.expand(); err != nil {
		return 0, nil, err
	}

	hashCode := h.calHash(key)
	idx := uint64(0)
	for table := 0; table < 2; table++ {
		idx = hashCode & sizeMask(h.sizeExp[table])
		head := h.table[table][idx]
		for head != nil {
			if key == head.key {
				return idx, head, nil
			}

			head = head.next
		}

		if !h.isRehashing() {
			break
		}
	}

	return idx, nil, nil
}

func (h *HashMap[K, V]) rehash(n int) error {
	// 控制访问空槽位的个数
	emptyVisits := n * 10

	// 没有rehashing 就退出
	if !h.isRehashing() {
		return ErrHashing
	}

	// n是控制桶数
	for ; n > 0 && h.used[0] != 0; n-- {

		for h.table[0][h.rehashidx] == nil {
			h.rehashidx++
			emptyVisits--
			if emptyVisits == 0 {
				return nil
			}
		}

		// 取出hash槽中第一个元素
		head := h.table[0][h.rehashidx]
		for head != nil {
			next := head.next
			newIdx := h.calHash(head.key) & sizeMask(h.sizeExp[1])
			head.next = h.table[1][newIdx]
			h.table[1][newIdx] = head
			h.used[0]--
			h.used[1]++
			head = next
		}

		h.table[0][h.rehashidx] = nil
		h.rehashidx++
	}

	if h.used[0] == 0 {
		h.table[0] = h.table[1]
		h.used[0] = h.used[1]
		h.sizeExp[0] = h.sizeExp[1]
		h.reset(1)
		// 这里重装置为-1
		h.rehashidx = -1
	}
	return nil
}

func (h *HashMap[K, V]) reset(idx int) {
	h.table[idx] = nil
	h.sizeExp[idx] = -1
	h.used[idx] = 0
}

func hashSize(exp int8) uint64 {
	if exp == -1 {
		return 0
	}
	return 1 << exp
}

func sizeMask(exp int8) uint64 {
	if exp == -1 {
		return 0
	}

	return (1 << exp) - 1
}

// 获取
func (h *HashMap[K, V]) Get(key K) (v V, err error) {
	if h.Len() == 0 {
		err = ErrNotFound
		return
	}

	if h.isRehashing() {
		h.rehash(1)
	}

	hashCode := h.calHash(key)
	idx := uint64(0)
	for table := 0; table < 2; table++ {
		idx = hashCode & sizeMask(h.sizeExp[table])
		head := h.table[table][idx]
		for head != nil {
			if key == head.key {
				return head.val, nil
			}

			head = head.next
		}

		if !h.isRehashing() {
			break
		}
	}
	err = ErrNotFound
	return
}

// 获取
func (h *HashMap[K, V]) GetOrZero(key K) (v V) {
	v, _ = h.Get(key)
	return
}

// 遍历
func (h *HashMap[K, V]) Range(pr func(key K, val V)) (err error) {
	if h.Len() == 0 {
		err = ErrNotFound
		return
	}

	if h.isRehashing() {
		h.rehash(1)
	}

	length := h.Len()
	for table := 0; table < 2 && length > 0; table++ {

		for idx := 0; idx < len(h.table[table]); idx++ {
			head := h.table[table][idx]
			for head != nil {
				pr(head.key, head.val)
				head = head.next
			}

			length--
		}
		if !h.isRehashing() {
			break
		}
	}
	return nil
}

// 设置
func (h *HashMap[K, V]) Set(k K, v V) error {
	h.lazyinit()
	if h.isRehashing() {
		h.rehash(1)
	}

	index, e, err := h.findIndexAndEntry(k)
	if err != nil {
		return err
	}

	idx := 0
	if h.isRehashing() {
		//如果在rehasing过程中, 如果这个key是第一次存入到hash table, 优先写入到新hash table中
		idx = 1
	}

	// element存在, 这里是替换
	if e != nil {
		//e.key = k
		e.val = v
		return nil
	}

	e = &entry[K, V]{key: k, val: v}
	e.next = h.table[idx][index]
	h.table[idx][index] = e
	h.used[idx]++
	return nil
}

// Remove是delete别名
func (h *HashMap[K, V]) Remove(key K) (err error) {
	return h.Delete(key)
}

// 删除
func (h *HashMap[K, V]) Delete(key K) (err error) {
	if h.Len() == 0 {
		err = ErrNotFound
		return
	}

	if h.isRehashing() {
		h.rehash(1)
	}

	hashCode := h.calHash(key)
	idx := uint64(0)
	for table := 0; table < 2; table++ {
		idx = hashCode & sizeMask(h.sizeExp[table])
		var prev *entry[K, V]
		head := h.table[table][idx]
		for head != nil {
			if key == head.key {
				if prev != nil {
					// 使用双指针删除中间的元素
					prev.next = head.next
				} else {
					// 表头元素, 直接跳过就可以删除
					h.table[table][idx] = head.next
				}
				h.used[table]--
				return nil
			}

			prev = head
			head = head.next
		}

		if !h.isRehashing() {
			break
		}
	}
	return nil
}

// 测试长度
func (h *HashMap[K, V]) Len() int {
	return int(h.used[0] + h.used[1])
}
