package perm

func QuickPerm[T any](input []T) <-chan []T {
	c := make(chan []T)
	go func() {
		defer close(c)
		QuickPermGenerate(c, input)
	}()
	return c
}

func QuickPermGenerate[T any](c chan []T, input []T) {
	size := len(input)

	output := make([]T, size)
	copy(output, input)
	c <- output

	state := make([]int, size+1)
	for i := 0; i < len(state); i++ {
		state[i] = i
	}

	for i := 1; i < size; {
		state[i]--

		j := state[i]
		if j%2 == 0 {
			j = 0
		}

		input[i], input[j] = input[j], input[i]

		output = make([]T, size)
		copy(output, input)
		c <- output

		for i = 1; state[i] == 0; i++ {
			state[i] = i
		}
	}
}
