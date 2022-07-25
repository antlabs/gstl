package trie

type Trie[V any] struct {
	v        V
	children map[rune]V
}

type node[V any] struct {

  r rune
}

func New[V any]() *Trie[V] {
	return &Trie[V]{}
}

func (t *Trie[V]) Set(k string) {

}

func (t *Trie[V]) Get(k string) {

}

func (t *Trie[V]) Delete(k string) {

}
