package trie

// apache 2.0 guonaihong

type Trie[V any] struct {
	v        V
	children map[rune]*Trie[V]
	set      bool
}

func New[V any]() *Trie[V] {
	return &Trie[V]{}
}

func (t *Trie[V]) Set(k string, v V) {
	_, _ = t.SetWithPrev(k, v)
}

func (t *Trie[V]) SetWithPrev(k string, v V) (prev V, replaced bool) {
	n := t
	for _, r := range k {
		c := n.children[r]
		if c == nil {
			if n.children == nil {
				n.children = map[rune]*Trie[V]{}
			}
			c = &Trie[V]{}
		}
		n = c
	}

	prev = n.v
	n.v = v

	replaced = n.set
	n.set = true
	return
}

func (t *Trie[V]) Get(k string) (v V) {
	n := t
	for _, r := range k {
		n = n.children[r]
		if n == nil {
			return
		}
	}
	return n.v
}

func (t *Trie[V]) Delete(k string) {

}
