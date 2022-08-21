package radix

// 健值对
type pair[V any] struct {
	val V
	key string
}

// 边
type edge[V any] struct {
	label rune
	node  *node[V]
}

// 节点
type node[V any] struct {
	pair[V]
	prefix string
}

// 头节点
type Tree[V any] struct {
	root   *node[V]
	length int
}

func (t *Tree[V]) Get(k string) (v V) {

	return
}

func (t *Tree[V]) SetWithPrev(k string, v V) (prev V, replaced bool) {

	return
}

func (t *Tree[V]) HasPrefix(k string) (ok bool) {
	return
}

func (t *Tree[V]) GetWithBool(k string) (v V, found bool) {

	return
}
func (t *Tree[V]) Delete(k string) {

}

func (t *Tree[V]) Len() int {

	return t.length
}
