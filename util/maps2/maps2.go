package maps2

func Pop[M ~map[K]V, K comparable, V any](m M) (K, V) {
	for k, v := range m {
		delete(m, k)
		return k, v
	}
	panic("pop from empty map")
}
