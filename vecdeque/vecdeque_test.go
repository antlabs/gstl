package vecdeque

// apache 2.0 antlabs
import "testing"

func Test_PushBack(t *testing.T) {
	v := New[int]()

	max := 100
	need := make([]int, 0, max)
	got := make([]int, 0, max)

	for i := 0; i < max; i++ {
		need = append(need, i)

		v.PushBack(i)
		v2, err := v.PopFront()
		if err != nil {
			break
		}
		got = append(got, v2)
	}
}
