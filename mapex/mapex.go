package mapex

type Map[K comparable, V any] map[K]V

func Keys[K comparable, V any](m map[K]V) (keys []K) {
	return Map[K, V](m).Keys()
}

func Values[K comparable, V any](m map[K]V) (values []V) {
	return Map[K, V](m).Values()
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
