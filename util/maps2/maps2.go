package maps2

func Pop[M ~map[K]V, K comparable, V any](m M) (K, V) {
	for k, v := range m {
		delete(m, k)
		return k, v
	}
	panic("pop from empty map")
}

func CountValuesFunc[M ~map[K]V, K comparable, V any](m M, pred func(V) bool) (count int) {
	for _, v := range m {
		if pred(v) {
			count++
		}
	}
	return
}

func CountValues[M ~map[K]V, K comparable, V comparable](m M, needle V) (count int) {
	return CountValuesFunc(m, func(v V) bool { return v == needle })
}
