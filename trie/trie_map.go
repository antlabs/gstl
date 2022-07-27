package trie

import "unicode/utf8"

// apache 2.0 guonaihong

type Trie[V any] struct {
	v V
	// 这里也可以换成多少数据结构, 压测下性能
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
			n.children[r] = c
		}

		n = c
	}

	prev = n.v
	n.v = v

	replaced = n.set
	n.set = true
	return
}

func (t *Trie[V]) HasPrefix(k string) bool {

	n := t
	for _, r := range k {
		n = n.children[r]
		if n == nil {
			return false
		}
	}

	return true
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

func (t *Trie[V]) isLeaf() bool {
	return len(t.children) == 0
}

// 记录删除的过程
type recogNode[V any] struct {
	r rune
	n *Trie[V]
}

// 删除有两种方法, 这里先选择第1种，后面有时间再压测下第二种效率如何
// 1.记录rune和节点，删除这个节点。如果是子节点，再回溯删除
// 2.声明一个parent指针，不记录过程节点，直接p = n.parent; p != nil; p=p.parent 回溯删除
func (t *Trie[V]) Delete(k string) {
	recog := make([]recogNode[V], 0, utf8.RuneCountInString(k))

	var v V
	n := t

	for _, r := range k {
		recog = append(recog, recogNode[V]{r, n})
		n = n.children[r]
		if n == nil {
			return
		}
	}

	n.v = v
	n.set = false

	if !n.isLeaf() {
		return
	}

	for last := len(recog) - 1; last >= 0; last-- {
		p := recog[last].n
		delete(p.children, recog[last].r)

		if !p.isLeaf() {
			return
		}
	}
}
