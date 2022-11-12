package avltree

import (
	"fmt"
	"testing"
)

// b.N = 3kw
// pkg: github.com/antlabs/gstl/avltree
// BenchmarkGetAsc-8    	33178270	        41.07 ns/op
// BenchmarkGetDesc-8   	33488839	        39.91 ns/op
// BenchmarkGetStd-8    	29553132	        49.34 ns/op

func BenchmarkGetAsc(b *testing.B) {
	//max := 1000000.0 * 5
	max := float64(b.N)
	set := New[float64, float64]()
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
	//max := 1000000.0 * 5
	max := float64(b.N)
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

	//max := 1000000.0 * 5
	max := float64(b.N)
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
