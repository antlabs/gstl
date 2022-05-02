package rhashmap

import (
	"fmt"
	"testing"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/guonaihong/gstl/rhashmap
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkGet-8       1000000000          0.4066 ns/op
// BenchmarkGetStd-8    1000000000          0.8333 ns/op
// PASS
// ok   github.com/guonaihong/gstl/rhashmap 130.007s.
// 比标准库快一倍.


// 五百万数据的Get操作时间
func BenchmarkGet(b *testing.B) {
	max := 1000000.0 * 5
	set := NewWithOpt[float64, float64](WithCap(int(max)))
	for i := 0.0; i < max; i++ {
		set.Set(i, i)
	}

	b.ResetTimer()

	for i := 0.0; i < max; i++ {
		v := set.GetOrZero(i)
		if v != i {
			panic(fmt.Sprintf("need:%f, got:%f", i, v))
		}
	}
}

func BenchmarkGetStd(b *testing.B) {

	max := 1000000.0 * 5
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
