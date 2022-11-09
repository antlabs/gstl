// apache 2.0 antlabs

package rwmap

import "sync"

type RWMap[K comparable, V any] struct {
	rw sync.RWMutex
	m  map[K]V
}

type Pair[K comparable, V any] struct {
	Key K
	Val V
}

// 通过new函数分配可以指定map的长度
func New[K comparable, V any](l int) *RWMap[K, V] {
	return &RWMap[K, V]{
		m: make(map[K]V, l),
	}
}

// 删除
func (r *RWMap[K, V]) Delete(key K) {
	r.rw.Lock()
	delete(r.m, key)
	r.rw.Unlock()
}

// 加载
func (r *RWMap[K, V]) Load(key K) (value V, ok bool) {
	r.rw.RLock()
	value, ok = r.m[key]
	r.rw.RUnlock()
	return

}

// 获取值，然后并删除
func (r *RWMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	r.rw.Lock()
	if r.m == nil {
		r.rw.Unlock()
		return
	}
	value, loaded = r.m[key]
	delete(r.m, key)
	r.rw.Unlock()
	return
}

// 存在返回现有的值，loaded 为true
// 不存在就保存现在的值，loaded为false
func (r *RWMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	r.rw.Lock()
	if r.m == nil {
		r.m = make(map[K]V)
	}
	actual, loaded = r.m[key]
	if !loaded {
		actual = value
		r.m[key] = actual
	}
	r.rw.Unlock()
	return
}

func (r *RWMap[K, V]) Range(f func(key K, value V) bool) {
	r.rw.RLock()
	for k, v := range r.m {
		if !f(k, v) {
			break
		}
	}
	r.rw.RUnlock()
}

func (r *RWMap[K, V]) Iter() <-chan Pair[K, V] {
	p := make(chan Pair[K, V])
	go func() {
		r.rw.RLock()
		for k, v := range r.m {
			p <- Pair[K, V]{Key: k, Val: v}
		}
		close(p)
		r.rw.RUnlock()

	}()
	return p
}

// 保存值
func (r *RWMap[K, V]) Store(key K, value V) {
	r.rw.Lock()
	if r.m == nil {
		r.m = make(map[K]V)
	}
	r.m[key] = value
	r.rw.Unlock()
}

// 返回长度
func (r *RWMap[K, V]) Len() (l int) {
	r.rw.RLock()
	l = len(r.m)
	r.rw.RUnlock()
	return
}
