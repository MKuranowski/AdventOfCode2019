package day17

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/set"
)

type Point struct{ X, Y int }

type Scaffolding = set.Set[Point]

func ReceiveScreenOnce(s Scaffolding, ch <-chan int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	x := 0
	y := 0

	for c := range ch {
		switch c {
		case '\n':
			x, y = 0, y+1
		case '#', '<', '^', '>', 'v':
			s.Add(Point{x, y})
			x++
		default:
			x++
		}
	}
}

func SolveA(r io.Reader) any {
	// Run the program and map out the scaffolding
	scaffolding := make(Scaffolding)
	i := intcode.NewInterpreterNewIO(r)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	close(i.Input)
	go ReceiveScreenOnce(scaffolding, i.Output, wg)
	i.ExecAll()
	wg.Wait()

	// Calculate the alignment
	alignment := 0
	for pt := range scaffolding {
		// Count adjacent blocks
		adjacent := 0
		if scaffolding.Has(Point{pt.X - 1, pt.Y}) {
			adjacent++
		}
		if scaffolding.Has(Point{pt.X + 1, pt.Y}) {
			adjacent++
		}
		if scaffolding.Has(Point{pt.X, pt.Y - 1}) {
			adjacent++
		}
		if scaffolding.Has(Point{pt.X, pt.Y + 1}) {
			adjacent++
		}

		// At least 3 adjacent blocks - an intersection
		if adjacent > 2 {
			alignment += pt.X * pt.Y
		}
	}

	return alignment
}

// Figured out this by hand, lol
// A C A B A B C B B C
// A: L4 L4 L10 R4
// B: R4 L10 R10
// C: R4 L4 L4 R8 R10
const solution = "A,C,A,B,A,B,C,B,B,C\nL,4,L,4,L,10,R,4\nR,4,L,10,R,10\nR,4,L,4,L,4,R,8,R,10\nn\n"

func SolutionSender(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, b := range solution {
		ch <- int(b)
	}
	close(ch)
}

type Screen struct {
	LastNonASCII int
}

func (s *Screen) Run(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for c := range ch {
		if c >= 0x7F {
			s.LastNonASCII = c
		} else {
			fmt.Fprintf(os.Stderr, "%c", c)
		}
	}
}

func SolveB(r io.Reader) any {
	s := Screen{}
	i := intcode.NewInterpreterNewIO(r)
	i.Memory[0] = 2

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go s.Run(i.Output, wg)
	go SolutionSender(i.Input, wg)
	i.ExecAll()

	wg.Wait()
	return s.LastNonASCII
}
