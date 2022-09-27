package day07

import (
	"io"
	"math"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/perm"
)

func SolveA(r io.Reader) any {
	prog := intcode.NewInterpreter(r)
	amps := 5
	maxPower := math.MinInt

	// Prepare the input for generating permutations
	permInput := make([]int, amps)
	for i := 0; i < len(permInput); i++ {
		permInput[i] = i
	}

	// Iterate over every possible permutation of phase settings
	for phases := range perm.QuickPerm(permInput) {
		wg := &sync.WaitGroup{}

		// The input of the first amp is 0
		// NOTE: Channels need to be buffered to send the phase setting before starting the amplifier
		firstInput := make(chan int, 1)
		lastOutput := firstInput

		// Launch all amps chaining their inputs and outputs
		for i := 0; i < amps; i++ {
			amp := prog.Clone()

			// Set the amp's input - previous amplifier's output
			lastOutput <- phases[i]
			amp.Input = lastOutput

			// Connect output
			lastOutput = make(chan int, 1)
			amp.Output = lastOutput

			// Run the amplifier
			wg.Add(1)
			go func(amp *intcode.Interpreter) {
				defer wg.Done()
				amp.ExecAll()
			}(amp)
		}

		firstInput <- 0
		wg.Wait()
		power := <-lastOutput

		if power > maxPower {
			maxPower = power
		}
	}

	return maxPower
}

func SolveB(r io.Reader) any {
	prog := intcode.NewInterpreter(r)
	amps := 5
	maxPower := math.MinInt

	// Prepare the input for generating permutations
	permInput := make([]int, amps)
	for i := 0; i < len(permInput); i++ {
		permInput[i] = 5 + i
	}

	// Iterate over every possible permutation of phase settings
	for phases := range perm.QuickPerm(permInput) {
		wg := &sync.WaitGroup{}

		// Create the very first channel which loops around
		// NOTE: Channels need to be buffered to send the phase setting before starting the amplifier
		loop := make(chan int, 1)
		lastOutput := loop

		// Launch all amps chaining their inputs and outputs
		for i := 0; i < amps; i++ {
			amp := prog.Clone()

			// Set the amp's input - previous amplifier's output
			amp.Input = lastOutput
			lastOutput <- phases[i]

			// Connect output
			if i == amps-1 {
				lastOutput = loop
				amp.Output = loop
			} else {
				lastOutput = make(chan int, 1)
				amp.Output = lastOutput
			}

			// Run the amplifier
			wg.Add(1)
			go func(amp *intcode.Interpreter, i int) {
				defer wg.Done()
				amp.ExecAll()
			}(amp, i)
		}

		loop <- 0
		wg.Wait()
		power := <-loop

		if power > maxPower {
			maxPower = power
		}
	}

	return maxPower
}
