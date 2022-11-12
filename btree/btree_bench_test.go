package btree

// apache 2.0 antlabs
import (
	"fmt"
	"testing"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/antlabs/gstl/btree
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkGet-8   	1000000000	         0.5326 ns/op
// PASS
// ok  	github.com/antlabs/gstl/btree	25.315s
// 五百万数据的Get操作时间

// goos: darwin
// goarch: arm64
// pkg: github.com/antlabs/gstl/btree
// BenchmarkGetAsc-8    	17242494	        79.54 ns/op
// BenchmarkGetDesc-8   	17556082	        78.17 ns/op
// BenchmarkGetStd-8    	29304117	        50.49 ns/op
// PASS
// ok  	github.com/antlabs/gstl/btree	10.503s
func BenchmarkGetAsc(b *testing.B) {
	//max := 1000000.0 * 5
	set := New[float64, float64](0)
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
	set := New[float64, float64](0)
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
