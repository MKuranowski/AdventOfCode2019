package day24

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/bits"

	"github.com/MKuranowski/AdventOfCode2019/util/set"
)

const SideLength = 5

func CoordsToOffsetChecked(x, y int) int {
	if x < 0 || y < 0 || x >= SideLength || y >= SideLength {
		panic(fmt.Errorf("coords out-of-bounds: %d %d", x, y))
	}
	return SideLength*y + x
}

func CoordsToOffset(x, y int) int {
	return SideLength*y + x
}

type Eris uint32

func (m Eris) CountBugs() int {
	return bits.OnesCount32(uint32(m))
}

func (m *Eris) Set(x, y int) {
	*m |= 1 << CoordsToOffset(x, y)
}

func (m Eris) Has(x, y int) bool {
	return m&(1<<CoordsToOffset(x, y)) != 0
}

func (m Eris) CountNeighbors(x, y int) (neighbors int) {
	if x > 0 && m.Has(x-1, y) {
		neighbors++
	}
	if x < SideLength-1 && m.Has(x+1, y) {
		neighbors++
	}
	if y > 0 && m.Has(x, y-1) {
		neighbors++
	}
	if y < SideLength-1 && m.Has(x, y+1) {
		neighbors++
	}
	return
}

func (m Eris) Evolve() (n Eris) {
	for x := 0; x < SideLength; x++ {
		for y := 0; y < SideLength; y++ {
			hasBug := m.Has(x, y)
			neighbors := m.CountNeighbors(x, y)

			if neighbors == 1 || (neighbors == 2 && !hasBug) {
				n.Set(x, y)
			}
		}
	}
	return
}

func (m Eris) BioDiversity() int {
	result := 0
	power := 1

	for m > 0 {
		if m&1 == 1 {
			result += power
		}

		power *= 2
		m >>= 1
	}

	return result
}

func ReadEris(r io.Reader) (m Eris) {
	br := bufio.NewReader(r)
	x, y := 0, 0

	for {
		c, err := br.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to read map: %w", err))
		}

		switch c {
		case '#':
			m.Set(x, y)
			x++
		case '\n':
			y++
			x = 0
		default:
			x++
		}
	}

	return
}

func SolveA(r io.Reader) any {
	m := ReadEris(r)
	seen := make(set.Set[Eris])

	for !seen.Has(m) {
		seen.Add(m)
		m = m.Evolve()
	}

	return m.BioDiversity()
}
