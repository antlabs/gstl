package btree

import (
	"fmt"
	"testing"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/guonaihong/gstl/btree
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkGet-8   	1000000000	         0.5326 ns/op
// PASS
// ok  	github.com/guonaihong/gstl/btree	25.315s
// 五百万数据的Get操作时间
func BenchmarkGet(b *testing.B) {
	max := 1000000 * 5.0
	bt := New[float64, float64](0)
	for i := 0.0; i < max; i++ {
		bt.Set(i, i)
	}

	b.ResetTimer()

	for i := 0.0; i < max; i++ {
		v := bt.Get(i)
		if v != i {
			panic(fmt.Sprintf("need:%f, got:%f", i, v))
		}
	}
}
