package day02

import (
	"io"
	"runtime"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
)

func SolveA(r io.Reader) any {
	i := intcode.NewInterpreter(r)
	i.Memory[1] = 12
	i.Memory[2] = 2
	i.ExecAll()
	return i.Memory[0]
}

type InputB struct {
	Noun, Verb int
}

func SolveB(r io.Reader) any {
	baseInterpreter := intcode.NewInterpreter(r)
	ins := make(chan InputB)
	results := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for input := range ins {
				i := baseInterpreter.Clone()
				i.Memory[1] = input.Noun
				i.Memory[2] = input.Verb
				i.ExecAll()
				if i.Memory[0] == 19690720 {
					results <- 100*input.Noun + input.Verb
					break
				}
			}
		}()
	}

	for i := 0; i < 100*100; i++ {
		select {
		case ins <- InputB{i / 100, i % 100}:

		case result := <-results:
			close(ins)
			return result
		}
	}

	panic("no solution")
}
