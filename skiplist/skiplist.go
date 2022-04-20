package skiplist

import (
	"math/rand"
	"time"
)

const (
	SKIPLIST_MAXLEVEL = 32
)

type Node[T any] struct {
	score     float64
	ele       T
	backward  *Node[T]
	NodeLevel []struct {
		forward *Node[T]
		span    int
	}
}

type SkipList[T any] struct {
	head *Node[T]
	tail *Node[T]

	r      *rand.Rand
	length int
	level  int
}

func New[T any]() *SkipList[T] {
	sl := &SkipList[T]{
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		level: 1,
	}
	sl.head = &Node[T]{
		//NodeLevel: make(nodeLevel, SKIPLIST_MAXLEVEL),
		NodeLevel: make([]struct {
			forward *Node[T]
			span    int
		}, SKIPLIST_MAXLEVEL),
	}
	return sl
}

func (s *SkipList[T]) rand() int {
	level := 1
	for {

		if s.r.Int()%2 == 0 {
			break
		}
		level++
	}

	if level < SKIPLIST_MAXLEVEL {
		return level
	}

	return SKIPLIST_MAXLEVEL
}
