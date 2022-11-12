package rbtree

// apache 2.0 antlabs

// b.N = 500w
// goos: darwin
// goarch: amd64
// pkg: github.com/antlabs/gstl/rbtree
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkGetAsc-8    	1000000000	         0.3336 ns/op
// BenchmarkGetDesc-8   	1000000000	         0.3702 ns/op
// BenchmarkGetStd-8
// 1000000000	         0.8940 ns/op
// PASS
// ok  	github.com/antlabs/gstl/rbtree	139.415s

// b.N = 3kw
// goos: darwin
// goarch: arm64
// pkg: github.com/antlabs/gstl/rbtree
// BenchmarkGetAsc-8    	32662837	        40.22 ns/op
// BenchmarkGetDesc-8   	33250437	        39.52 ns/op
// BenchmarkGetStd-8    	29353758	        49.73 ns/op
// PASS
// ok  	github.com/antlabs/gstl/rbtree	13.030s

import (
	"fmt"
	"testing"
)

func BenchmarkSetAsc(b *testing.B) {
	//max := 1000000.0 * 5
	set := New[float64, float64]()
	max := float64(b.N)
	for i := 0.0; i < max; i++ {
		set.Set(i, i)
	}

}

func BenchmarkGetAsc(b *testing.B) {
	//max := 1000000.0 * 5
	set := New[float64, float64]()
	max := float64(b.N)
	for i := 0.0; i < max; i++ {
		set.Set(i, i)
	}

	b.ResetTimer()

	for i := 0.0; i < max; i++ {
		v := set.Get(i)
		if v != i {
			panic(fmt.Sprintf("need:%f, got:%f", i, v))
		}
	}
}

func BenchmarkGetDesc(b *testing.B) {
	max := float64(b.N)
	//max := 1000000.0 * 5
	set := New[float64, float64]()
	for i := max; i >= 0; i-- {
		set.Set(i, i)
	}

	b.ResetTimer()

	for i := 0.0; i < max; i++ {
		v := set.Get(i)
		if v != i {
			panic(fmt.Sprintf("need:%f, got:%f", i, v))
		}
	}
}

func BenchmarkGetStd(b *testing.B) {

	max := float64(b.N)
	//max := 1000000.0 * 5
	set := make(map[float64]float64, int(max))
	for i := 0.0; i < max; i++ {
		set[i] = i
	}

	b.ResetTimer()

	for i := 0.0; i < max; i++ {
		v := set[i]
		if v != i {
			panic(fmt.Sprintf("need:%f, got:%f", i, v))
		}
	}
}

func BenchmarkSet(b *testing.B) {
	max := float64(b.N)
	//max := 1000000.0 * 5
	set := New[float64, float64]()
	for i := max; i >= 0; i-- {
		set.Set(i, i)
	}

}

func BenchmarkSetStd(b *testing.B) {
	max := float64(b.N)
	//max := 1000000.0 * 5
	set := make(map[float64]float64)
	for i := max; i >= 0.0; i-- {
		set[i] = i
	}

}
