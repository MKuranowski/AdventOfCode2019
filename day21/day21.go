package day21

import (
	"io"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/day17"
	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

// Jump if there's a hole at A, B or C and not D; aka
// J = (~A or ~B or ~C) and D
// J = ~(A and B and C) and D
const SolutionA = "OR A T\nAND B T\nAND C T\nNOT T J\nAND D J\nWALK\n"

// J = (~A or ~B or ~C) and D and (E or H)
const SolutionB = "OR A T\nAND B T\nAND C T\nNOT T J\nAND D J\nOR E T\nOR H T\nAND T J\nRUN\n"

func Solve(r io.Reader, solution string) int {
	s := day17.Screen{}
	i := intcode.NewInterpreterNewIO(r)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go s.Run(i.Output, wg)
	go input.AsciiStaticSender(i.Input, wg, solution)
	i.ExecAll()

	wg.Wait()
	return s.LastNonASCII
}

func SolveA(r io.Reader) any { return Solve(r, SolutionA) }
func SolveB(r io.Reader) any { return Solve(r, SolutionB) }
