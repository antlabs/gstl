package skiplist

// 参考文档如下
// https://github.com/redis/redis/blob/unstable/src/t_zset.c
import (
	"math/rand"
	"time"
)

const (
	SKIPLIST_MAXLEVEL = 32
)

type Node[T any] struct {
	score     float64
	elem      T
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

	compare func(T, T) int
}

// 初始化skiplist
func New[T any](compare func(T, T) int) *SkipList[T] {
	s := &SkipList[T]{
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		level: 1,
	}
	s.head = newNode[T](SKIPLIST_MAXLEVEL, 0, *new(T))
	return s
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

func newNode[T any](level int, score float64, elem T) *Node[T] {
	return &Node[T]{
		score: score,
		elem:  elem,
		NodeLevel: make([]struct {
			forward *Node[T]
			span    int
		}, level),
	}
}

func (s *SkipList[T]) Insert(score float64, elem T) {
	var (
		update [SKIPLIST_MAXLEVEL]*Node[T]
		rank   [SKIPLIST_MAXLEVEL]int
	)

	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		if i == s.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.NodeLevel[i].forward != nil &&
			(x.NodeLevel[i].forward.score < score ||
				x.NodeLevel[i].forward.score == score &&
					s.compare(elem, x.NodeLevel[i].forward.elem) < 0) {
			//TODO span的含义是?
			rank[i] += x.NodeLevel[i].span
			x = x.NodeLevel[i].forward
		}

		// 保存插入位置的前一个节点, 保存的数量就是level的值
		update[i] = x
	}

	// 生成新节点的level
	level := s.rand()
	if level > s.level {
		// 这次新的level与老的level的差值, 给填充head指针
		for i := s.level; i < level; i++ {
			// TODO rank的含义
			rank[i] = 0
			update[i] = s.head
			update[i].NodeLevel[i].span = s.length
		}
		s.level = level
	}

	// 创建新节点
	x = newNode[T](level, score, elem)
	for i := 0; i < level; i++ {
		x.NodeLevel[i].forward = update[i].NodeLevel[i].forward
		update[i].NodeLevel[i].forward = x
		x.NodeLevel[i].span = update[i].NodeLevel[i].span - (rank[0] - rank[i])
		update[i].NodeLevel[i].span = rank[0] - rank[i] + 1
	}

	for i := level; i < s.level; i++ {
		update[i].NodeLevel[i].span++
	}

	if update[0] != s.head {
		x.backward = update[0]
	}

	if x.NodeLevel[0].forward != nil {
		x.NodeLevel[0].forward.backward = x
	} else {
		s.tail = x
	}

	s.length++
}
