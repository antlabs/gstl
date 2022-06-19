package skiplist

// apache 2.0 guonaihong
import (
	"fmt"
	"testing"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/guonaihong/gstl/skiplist
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkGet-8
// 1000000000	         0.7430 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/guonaihong/gstl/skiplist	114.326s

func intCmp(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// 五百万数据的Get操作时间
func BenchmarkGet(b *testing.B) {
	max := 1000000.0 * 5
	set := New(intCmp)
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
