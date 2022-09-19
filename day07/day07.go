package day07

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
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
		lastOutput := io.Reader(strings.NewReader("0\n"))

		// Launch all amps chaining their inputs and outputs
		for i := 0; i < amps; i++ {
			amp := prog.Clone()

			// Set the amp's input - previous amplifier's output
			amp.Input = io.MultiReader(
				strings.NewReader(fmt.Sprintln(phases[i])),
				lastOutput,
			)

			// Connect output
			lastOutput, amp.Output = io.Pipe()

			// Run the amplifier
			wg.Add(1)
			go func(amp *intcode.Interpreter) {
				defer wg.Done()
				amp.ExecAll()
			}(amp)
		}

		powerStr, err := io.ReadAll(lastOutput)
		if err != nil {
			panic(fmt.Errorf("failed to read from last amplifier: %w", err))
		}

		wg.Wait()

		power, err := strconv.Atoi(strings.TrimRight(string(powerStr), "\n"))
		if err != nil {
			panic(fmt.Errorf("last amplifier didn't send a number (%q): %w", string(powerStr), err))
		}

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
		aDone := make(chan struct{})

		// The input of the first amp is 0
		loopReader, loopWriter := io.Pipe()
		lastOutput := io.MultiReader(strings.NewReader("0\n"), loopReader)

		// Launch all amps chaining their inputs and outputs
		for i := 0; i < amps; i++ {
			amp := prog.Clone()

			// Set the amp's input - previous amplifier's output
			amp.Input = io.MultiReader(
				strings.NewReader(fmt.Sprintln(phases[i])),
				lastOutput,
			)

			// Connect output
			if i == amps-1 {
				amp.Output = loopWriter
				lastOutput = loopReader
			} else {
				lastOutput, amp.Output = io.Pipe()
			}

			// Run the amplifier
			wg.Add(1)
			if i == 0 {
				go func(amp *intcode.Interpreter) {
					defer wg.Done()
					defer close(aDone)
					amp.ExecAll()
				}(amp)
			} else {
				go func(amp *intcode.Interpreter) {
					defer wg.Done()
					amp.ExecAll()
				}(amp)
			}
		}

		<-aDone

		powerStr, err := io.ReadAll(lastOutput)
		if err != nil {
			panic(fmt.Errorf("failed to read from last amplifier: %w", err))
		}

		wg.Wait()

		power, err := strconv.Atoi(strings.TrimRight(string(powerStr), "\n"))
		if err != nil {
			panic(fmt.Errorf("last amplifier didn't send a number (%q): %w", string(powerStr), err))
		}

		if power > maxPower {
			maxPower = power
		}
	}

	return maxPower
}
