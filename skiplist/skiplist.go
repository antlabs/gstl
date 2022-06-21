package skiplist

// apache 2.0 guonaihong

// 参考文档如下
// https://github.com/redis/redis/blob/unstable/src/t_zset.c
// https://redis.io/commands/zcount/
// ZADD
// ZCARD
// ZCOUNT
// ZDIFFSTORE
// ZINCRBY
// ZINTER
// ZINTERCARD
// ZINTERSTORE
// ZLEXCOUNT
// ZMPOP
// ZMSCORE
// ZPOPMAX
// ZPOPMIN
// ZRANDMEMBER
// ZRANGE
// ZRANGEBYLEX
// ZRANGEBYSCORE
// ZRANGESTORE
// ZRANK
// ZREMRANGEBYLEX
// ZREMRANGEBYRANK
// ZREMRANGEBYSCORE
// ZREVRANGE
// ZREVRANGEBYLEX
// ZREVRANGEBYSCORE
// ZREVRANK
// ZSCAN
// ZUNION
// ZUNIONSTORE
import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	SKIPLIST_MAXLEVEL = 32
)

var (
	ErrNotFound = errors.New("not found element")
)

type Node[T any] struct {
	score float64
	elem  T
	// 后退指针
	backward  *Node[T]
	NodeLevel []struct {
		// 指向前进节点, 是指向tail的方向
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

	//compare func(T, T) int
}

// 初始化skiplist
//func New[T any](compare func(T, T) int) *SkipList[T] {
func New[T any]() *SkipList[T] {
	s := &SkipList[T]{
		level: 1,
	}

	//s.compare = compare
	s.resetRand()
	s.head = newNode(SKIPLIST_MAXLEVEL, 0, *new(T))
	return s
}

func (s *SkipList[T]) resetRand() {

	s.r = rand.New(rand.NewSource(time.Now().UnixNano()))
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

// 设置值, 和Insert是同义词
func (s *SkipList[T]) Set(score float64, elem T) *SkipList[T] {
	return s.InsertInner(score, elem, s.rand())
}

// 设置值
func (s *SkipList[T]) Insert(score float64, elem T) *SkipList[T] {
	return s.InsertInner(score, elem, s.rand())
}

// 方便给作者调试用的函数
func (s *SkipList[T]) InsertInner(score float64, elem T, level int) *SkipList[T] {
	var (
		update [SKIPLIST_MAXLEVEL]*Node[T]
		rank   [SKIPLIST_MAXLEVEL]int
	)

	x := s.head
	var x2 *Node[T]
	for i := s.level - 1; i >= 0; i-- {
		if i == s.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.NodeLevel[i].forward != nil &&
			(x.NodeLevel[i].forward.score < score) {
			// 暂时不支持重复的key, 后面再说 TODO
			//|| x.NodeLevel[i].forward.score == score &&
			//s.compare(elem, x.NodeLevel[i].forward.elem) < 0) {

			//TODO span的含义是?
			rank[i] += x.NodeLevel[i].span
			x = x.NodeLevel[i].forward
		}

		// 保存插入位置的前一个节点, 保存的数量就是level的值
		update[i] = x
	}

	// 这个score已经存在直接返回
	x2 = x.NodeLevel[0].forward
	if x2 != nil && score == x2.score {
		x2.elem = elem
		return s
	}

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
	x = newNode(level, score, elem)
	for i := 0; i < level; i++ {
		// x.NodeLevel[i]的节点假设等于a, 需要插入的节点x在a之后,
		// a, x, a.forward三者的关系就是[a, x, a.forward]
		// 那就变成了x.forward = a.forward, 修改x.forward的指向
		// a.forward = x
		// 看如下两行代码
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
	return s
}

// 获取
func (s *SkipList[T]) GetWithErr(score float64) (elem T, err error) {

	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.NodeLevel[i].forward != nil && (x.NodeLevel[i].forward.score < score) {
			x = x.NodeLevel[i].forward
		}

		/*
			// 效果不大
						if x != s.head && x.score == score {
							elem = x.elem
							return
						}
		*/

		/*
			// 效果不大
			if x.NodeLevel[i].forward != nil && x.NodeLevel[i].forward.score == score {
				elem = x.NodeLevel[i].forward.elem
				return
							}
		*/
	}

	x = x.NodeLevel[0].forward
	if x != nil && score == x.score {
		return x.elem, nil
	}

	err = ErrNotFound
	return
}

// debug 使用
type Number struct {
	Total    int
	Keys     []float64
	Level    []int
	MaxLevel []int
}

// debug使用, 返回查找某个key 比较的次数+经过的节点数
func (s *SkipList[T]) GetWithMeta(score float64) (elem T, number Number, err error) {

	x := s.head
	fmt.Println()
	for i := s.level - 1; i >= 0; i-- {
		if x.NodeLevel[i].forward != nil {
			fmt.Printf("x.NodeLevel[%d].score:%f, score:%f\n", i, x.NodeLevel[i].forward.score, score)
		}
		for x.NodeLevel[i].forward != nil && (x.NodeLevel[i].forward.score < score) {
			number.Total++
			number.Keys = append(number.Keys, x.score)
			number.Level = append(number.Level, i)
			number.MaxLevel = append(number.MaxLevel, len(x.NodeLevel))
			x = x.NodeLevel[i].forward
		}

		if x != nil && x.score == score {
			elem = x.elem
			return
		}

		/*
			if x.NodeLevel[i].forward != nil && x.NodeLevel[i].forward.score == score {
				elem = x.NodeLevel[i].forward.elem
				return
			}
		*/
	}

	x = x.NodeLevel[0].forward
	if x != nil && score == x.score {
		return x.elem, number, nil
	}

	err = ErrNotFound
	return
}

// 根据score获取value值
func (s *SkipList[T]) Get(score float64) (elem T) {
	elem, _ = s.GetWithErr(score)
	return elem
}

func (s *SkipList[T]) removeNode(x *Node[T], update []*Node[T]) {
	for i := 0; i < s.level; i++ {
		if update[i].NodeLevel[i].forward == x {
			update[i].NodeLevel[i].span += x.NodeLevel[i].span - 1
			update[i].NodeLevel[i].forward = x.NodeLevel[i].forward
		} else {
			update[i].NodeLevel[i].span -= 1
		}
	}

	if x.NodeLevel[0].forward != nil {
		// 原来右边节点backward指针直接指各左边节点, 现在指向左左节点
		x.NodeLevel[0].forward.backward = x.backward
	} else {
		// 最后一个元素, 修改下尾指针的位置
		s.tail = x.backward
	}

	for s.level > 1 && s.head.NodeLevel[s.level-1].forward == nil {
		s.level--
	}
	s.length--
}

// 根据score删除元素
func (s *SkipList[T]) Remove(score float64) *SkipList[T] {

	var update [SKIPLIST_MAXLEVEL]*Node[T]
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.NodeLevel[i].forward != nil && (x.NodeLevel[i].forward.score < score) {
			x = x.NodeLevel[i].forward
		}
		update[i] = x
	}

	x = x.NodeLevel[0].forward
	if x != nil && score == x.score {
		s.removeNode(x, update[:])
		return s
	}

	return s
}

func (s *SkipList[T]) Draw() *SkipList[T] {
	if s.head == nil {
		return s
	}

	fmt.Printf("maxlevel:%d, head level:%d \n", s.level, len(s.head.NodeLevel))
	i := 1
	for h := s.head.NodeLevel[0].forward; h != nil; h = h.NodeLevel[0].forward {
		fmt.Printf("score:%f, level:%d -> ", h.score, len(h.NodeLevel))
		if i%6 == 0 {
			fmt.Printf("\n")
		}
		i++
	}

	fmt.Printf("\n")
	return s
}

// 遍历
func (s *SkipList[T]) Range(callback func(score float64, v T) bool) *SkipList[T] {
	if s.head == nil {
		return s
	}

	for h := s.head.NodeLevel[0].forward; h != nil; h = h.NodeLevel[0].forward {
		if !callback(h.score, h.elem) {
			return s
		}
	}
	return s
}

// 返回最小的n个值, 升序返回, 比如0,1,2,3
func (s *SkipList[T]) TopMin(limit int, callback func(score float64, v T) bool) *SkipList[T] {
	s.Range(func(score float64, v T) bool {
		if limit <= 0 {
			return false
		}
		callback(score, v)
		limit--
		return true
	})
	return s
}

// 返回长度
func (s *SkipList[T]) Len() int {
	return s.length
}

// 从后向前倒序遍历b tree
func (s *SkipList[T]) RangePrev(callback func(k float64, v T) bool) *SkipList[T] {
	// 遍历
	if s.tail == nil {
		return s
	}

	for t := s.tail; t != nil; t = t.backward {
		if !callback(t.score, t.elem) {
			return s
		}
	}

	return s
}

// 返回最大的n个值, 降序返回, 10, 9, 8, 7
func (s *SkipList[T]) TopMax(limit int, callback func(k float64, v T) bool) *SkipList[T] {
	s.RangePrev(func(k float64, v T) bool {
		if limit <= 0 {
			return false
		}
		callback(k, v)
		limit--
		return true
	})
	return s
}
