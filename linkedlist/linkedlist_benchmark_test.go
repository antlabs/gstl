package linkedlist

// apache 2.0 guonaihong
import (
	"container/list"
	"testing"
	"time"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/guonaihong/gstl/linkedlist
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// Benchmark_ListAdd_Stdlib-8   	 5918479	       190.0 ns/op
// Benchmark_ListAdd_gstl-8     	15942064	        83.15 ns/op
// PASS
// ok  	github.com/guonaihong/gstl/linkedlist	3.157s
type timeNodeStdlib struct {
	expire     uint64
	userExpire time.Duration
	callback   func()
	isSchedule bool
	close      uint32
	lock       uint32
}

// 标准库
func Benchmark_ListAdd_Stdlib(b *testing.B) {
	head := list.New()
	for i := 0; i < b.N; i++ {
		node := timeNodeStdlib{}
		head.PushBack(node)
	}
}

func Benchmark_ListAdd_gstl(b *testing.B) {
	head := New[timeNodeStdlib]()
	for i := 0; i < b.N; i++ {
		node := timeNodeStdlib{}
		head.PushBack(node)
	}
}
