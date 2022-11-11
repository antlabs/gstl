// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// guonaihong: 如何如下
// 1. interface的换成泛型语法

package cmap

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/antlabs/gstl/api"
	xxhash "github.com/cespare/xxhash/v2"
)

type syncmap[K comparable, V any] struct {
	m sync.Map
}

func (c *syncmap[K, V]) Delete(key K) {
	c.m.Delete(key)
}

func (c *syncmap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return
	}

	return v.(V), ok
}

func (c *syncmap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := c.m.LoadAndDelete(key)
	if !ok {
		return
	}

	return v.(V), ok
}

func (c *syncmap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, ok := c.m.LoadOrStore(key, value)
	if !ok {
		return
	}
	return v.(V), ok
}

func (c *syncmap[K, V]) Range(f func(key K, value V) bool) {
	c.m.Range(func(key any, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (c *syncmap[K, V]) Store(key K, value V) {
	c.m.Store(key, value)
}

type bench[K comparable, V any] struct {
	setup func(*testing.B, api.CMaper[K, V])
	perG  func(b *testing.B, pb *testing.PB, i int, m api.CMaper[K, V])
}

func benchMap(b *testing.B, bench bench[int, int]) {
	for _, m := range [...]api.CMaper[int, int]{New[int, int](), &syncmap[int, int]{}} {
		b.Run(fmt.Sprintf("%T", m), func(b *testing.B) {
			m = reflect.New(reflect.TypeOf(m).Elem()).Interface().(api.CMaper[int, int])
			if m2, ok := m.(*CMap[int, int]); ok {
				m2.init(1)
			}

			if bench.setup != nil {
				bench.setup(b, m)
			}

			b.ResetTimer()

			var i int64
			b.RunParallel(func(pb *testing.PB) {
				id := int(atomic.AddInt64(&i, 1) - 1)
				bench.perG(b, pb, id*b.N, m)
			})
		})
	}
}

func BenchmarkLoadMostlyHits(b *testing.B) {
	const hits, misses = 1023, 1

	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		},
	})
}

func BenchmarkLoadMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		},
	})
}

func BenchmarkLoadOrStoreBalanced(b *testing.B) {
	const hits, misses = 128, 128

	benchMap(b, bench[int, int]{
		setup: func(b *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				j := i % (hits + misses)
				if j < hits {
					if _, ok := m.LoadOrStore(j, i); !ok {
						b.Fatalf("unexpected miss for %v", j)
					}
				} else {
					if v, loaded := m.LoadOrStore(i, i); loaded {
						b.Fatalf("failed to store %v: existing value %v", i, v)
					}
				}
			}
		},
	})
}

func BenchmarkLoadOrStoreUnique(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(b *testing.B, m api.CMaper[int, int]) {
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(i, i)
			}
		},
	})
}

func BenchmarkDelete(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(b *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < 1000000; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Delete(i)
			}
		},
	})
}

func BenchmarkStore(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			//m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Store(i, i)
			}
		},
	})
}

func BenchmarkLoadOrStoreCollision(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(0, 0)
			}
		},
	})
}

func BenchmarkLoadAndDeleteBalanced(b *testing.B) {
	const hits, misses = 128, 128

	benchMap(b, bench[int, int]{
		setup: func(b *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				j := i % (hits + misses)
				if j < hits {
					m.LoadAndDelete(j)
				} else {
					m.LoadAndDelete(i)
				}
			}
		},
	})
}

func BenchmarkLoadAndDeleteUnique(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(b *testing.B, m api.CMaper[int, int]) {
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.LoadAndDelete(i)
			}
		},
	})
}

func BenchmarkLoadAndDeleteCollision(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.LoadAndDelete(0)
			}
		},
	})
}

func BenchmarkRange(b *testing.B) {
	const mapSize = 1 << 10

	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < mapSize; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Range(func(_, _ int) bool { return true })
			}
		},
	})
}

// BenchmarkAdversarialAlloc tests performance when we store a new value
// immediately whenever the map is promoted to clean and otherwise load a
// unique, missing key.
//
// This forces the Load calls to always acquire the map's mutex.
func BenchmarkAdversarialAlloc(b *testing.B) {
	benchMap(b, bench[int, int]{
		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			var stores, loadsSinceStore int
			for ; pb.Next(); i++ {
				m.Load(i)
				if loadsSinceStore++; loadsSinceStore > stores {
					m.LoadOrStore(i, stores)
					loadsSinceStore = 0
					stores++
				}
			}
		},
	})
}

// BenchmarkAdversarialDelete tests performance when we periodically delete
// one key and add a different one in a large map.
//
// This forces the Load calls to always acquire the map's mutex and periodically
// makes a full copy of the map despite changing only one entry.

//  这个case不测试, 锁分区的方式这么使用会死锁
/*
func BenchmarkAdversarialDelete(b *testing.B) {
	const mapSize = 1 << 10

	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			for i := 0; i < mapSize; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Load(i)

				if i%mapSize == 0 {
					m.Range(func(k, _ int) bool {
						m.Delete(k)
						return false
					})
					m.Store(i, i)
				}
			}
		},
	})
}
*/

func BenchmarkDeleteCollision(b *testing.B) {
	benchMap(b, bench[int, int]{
		setup: func(_ *testing.B, m api.CMaper[int, int]) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m api.CMaper[int, int]) {
			for ; pb.Next(); i++ {
				m.Delete(0)
			}
		},
	})
}

func BenchmarkXXHash(b *testing.B) {

	for i := 0; i < b.N; i++ {
		key := *(*string)(unsafe.Pointer(&reflect.StringHeader{
			Data: uintptr(unsafe.Pointer(&i)),
			Len:  8}))

		xxhash.Sum64String(key)
	}

}
