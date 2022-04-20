package skiplist

import (
	"math/rand"
	"time"
)

const (
	SKIPLIST_MAXLEVEL = 32
)

type Node[T any] struct {
}

type SkipList[T any] struct {
	r *rand.Rand
}

func New[T any]() *SkipList[T] {
	return &SkipList[T]{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
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
