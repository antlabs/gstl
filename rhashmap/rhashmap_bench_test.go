package rhashmap

// apache 2.0 antlabs
import (
	"fmt"
	"testing"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/antlabs/gstl/rhashmap
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkGet-8       1000000000          0.4066 ns/op
// BenchmarkGetStd-8    1000000000          0.8333 ns/op
// PASS
// ok   github.com/antlabs/gstl/rhashmap 130.007s.
// 比标准库快一倍.

// goos: darwin
// goarch: amd64
// pkg: github.com/antlabs/gstl/rhashmap
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkSet-8      	1000000000	         0.1690 ns/op
// BenchmarkSetStd-8   	1000000000	         0.1470 ns/op
// PASS
// ok  	github.com/antlabs/gstl/rhashmap	3.970s
// 五百万数据的Get操作时间

// TODO 再优化下性能
// go1.19.1
// 3kw
// goos: darwin
// goarch: arm64
// pkg: github.com/antlabs/gstl/rhashmap
// BenchmarkGet-8      	34664005	        62.20 ns/op
// BenchmarkGetStd-8   	30007470	        49.40 ns/op
// BenchmarkSet-8      	14623854	       178.9 ns/op
// BenchmarkSetStd-8   	22709601	        74.71 ns/op
// PASS
// ok  	github.com/antlabs/gstl/rhashmap	16.521s
func BenchmarkGet(b *testing.B) {
	//max := 1000000.0 * 5
	max := float64(b.N)
	set := NewWithOpt[float64, float64](WithCap(int(max)))
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

func BenchmarkGetStd(b *testing.B) {

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

// gstl set
func BenchmarkSet(b *testing.B) {
	max := float64(b.N)
	set := NewWithOpt[float64, float64](WithCap(int(max)))
	for i := 0.0; i < max; i++ {
		set.Set(i, i)
	}

}

// 标准库set
func BenchmarkSetStd(b *testing.B) {

	max := float64(b.N)
	set := make(map[float64]float64, int(max))
	for i := 0.0; i < max; i++ {
		set[i] = i
	}
}
