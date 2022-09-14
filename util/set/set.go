package set

type Set[T comparable] map[T]struct{}

// Has checks whether element `x` is present in the set
func (s Set[T]) Has(x T) bool {
	_, has := s[x]
	return has
}

// Add ensures that element `x` is present in the set
func (s Set[T]) Add(x T) { s[x] = struct{}{} }

// Remove ensures that element `x` is not present in the set
func (s Set[T]) Remove(x T) { delete(s, x) }

// Len returns the amount of elements in the set
func (s Set[T]) Len() int { return len(s) }

// Clear ensures that no elements are present in the set
func (s Set[T]) Clear() {
	for x := range s {
		delete(s, x)
	}
}

// Clone returns a shallow copy of the Set
func (s Set[T]) Clone() Set[T] {
	copy := make(Set[T], len(s))
	for x := range s {
		copy.Add(x)
	}
	return copy
}

// Union ensures `s1` contains all elements from `s2`
func (s1 Set[T]) Union(s2 Set[T]) {
	for x := range s2 {
		s1.Add(x)
	}
}

// Intersection ensures `s1` only contains elements that are present in `s2`
func (s1 Set[T]) Intersection(s2 Set[T]) {
	for x := range s1 {
		if !s2.Has(x) {
			s1.Remove(x)
		}
	}
}

// Difference ensures `s1` does not contain any element from `s2`
func (s1 Set[T]) Difference(s2 Set[T]) {
	for x := range s2 {
		s1.Remove(x)
	}
}
