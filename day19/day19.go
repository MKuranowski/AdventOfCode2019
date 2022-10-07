package day19

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
)

func Hits(prog *intcode.Interpreter, x, y int) bool {
	wg := &sync.WaitGroup{}
	progClone := prog.Clone()
	progClone.Input = make(chan int, 2)
	progClone.Output = make(chan int)

	wg.Add(1)
	go func(i *intcode.Interpreter) {
		defer wg.Done()
		i.ExecAll()
	}(progClone)

	progClone.Input <- x
	progClone.Input <- y
	hit := <-progClone.Output == 1
	wg.Wait()

	return hit
}

func SolveA(r io.Reader) any {
	prog := intcode.NewInterpreterNewIO(r)
	wg := &sync.WaitGroup{}
	result := 0

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if Hits(prog, x, y) {
				result++
				fmt.Fprint(os.Stderr, "#")
			} else {
				fmt.Fprint(os.Stderr, ".")
			}
		}
		fmt.Fprintln(os.Stderr)
	}

	wg.Wait()
	return result
}

func SolveB(r io.Reader) any {
	p := intcode.NewInterpreter(r)

	// NOTE: The beam is so narrow that it doesn't hit the first few rows -
	// that's why search is started at 10.
	beamStartX := 0
	for y := 10; y < 10_000; y++ {
		// Move to the column where beam starts
		for !Hits(p, beamStartX, y) {
			beamStartX++
		}

		for x := beamStartX; Hits(p, x+99, y); x++ {
			if Hits(p, x, y+99) && Hits(p, x+99, y+99) {
				return 10_000*x + y
			}
		}
	}

	panic("no solution")
}
