package mapex

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type Map[K comparable, V any] map[K]V

func Keys[K comparable, V any](m map[K]V) (keys []K) {
	return Map[K, V](m).Keys()
}

func SortKeys[K constraints.Ordered, V any](m map[K]V) (keys []K) {
	keys = Keys(m)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func Values[K comparable, V any](m map[K]V) (values []V) {
	return Map[K, V](m).Values()
}

func SortValues[K comparable, V constraints.Ordered](m map[K]V) (values []V) {
	values = Values(m)
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}

func (m Map[K, V]) Keys() (keys []K) {
	keys = make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func (m Map[K, V]) Values() (values []V) {
	values = make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return
}
