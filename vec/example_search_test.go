package vec

// apache 2.0 antlabs
import "fmt"

func Example_search() {
	vec := New(1, 2, 3, 4, 5, 6, 7)
	index := vec.SearchFunc(func(e int) bool {
		return 7 <= e
	})

	fmt.Println(index)
}
