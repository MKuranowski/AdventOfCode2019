package day24

import (
	"io"

	"github.com/MKuranowski/AdventOfCode2019/util/set"
)

type Coordinates struct{ X, Y, Depth int }

func RecursiveNeighbors(c Coordinates) (n []Coordinates) {
	// Add vertical neighbors
	if c.Y == 0 {
		n = append(
			n,
			Coordinates{2, 1, c.Depth + 1},     // Top (outer level)
			Coordinates{c.X, c.Y + 1, c.Depth}, // Bottom
		)
	} else if c.Y == 1 && c.X == 2 {
		n = append(
			n,
			Coordinates{c.X, c.Y - 1, c.Depth}, // Top
			Coordinates{0, 0, c.Depth - 1},     // Bottom (inner level)
			Coordinates{1, 0, c.Depth - 1},
			Coordinates{2, 0, c.Depth - 1},
			Coordinates{3, 0, c.Depth - 1},
			Coordinates{4, 0, c.Depth - 1},
		)
	} else if c.Y == 3 && c.X == 2 {
		n = append(
			n,
			Coordinates{0, 4, c.Depth - 1}, // Top (inner level)
			Coordinates{1, 4, c.Depth - 1},
			Coordinates{2, 4, c.Depth - 1},
			Coordinates{3, 4, c.Depth - 1},
			Coordinates{4, 4, c.Depth - 1},
			Coordinates{c.X, c.Y + 1, c.Depth}, // Bottom
		)
	} else if c.Y == 4 {
		n = append(
			n,
			Coordinates{c.X, c.Y - 1, c.Depth}, // Top
			Coordinates{2, 3, c.Depth + 1},     // Bottom (outer level)
		)
	} else {
		n = append(
			n,
			Coordinates{c.X, c.Y - 1, c.Depth}, // Top
			Coordinates{c.X, c.Y + 1, c.Depth}, // Bottom
		)
	}

	// Add horizontal neighbors
	if c.X == 0 {
		n = append(
			n,
			Coordinates{1, 2, c.Depth + 1},     // Left (outer level)
			Coordinates{c.X + 1, c.Y, c.Depth}, // Right
		)
	} else if c.X == 1 && c.Y == 2 {
		n = append(
			n,
			Coordinates{c.X - 1, c.Y, c.Depth}, // Left
			Coordinates{0, 0, c.Depth - 1},     // Right (inner level)
			Coordinates{0, 1, c.Depth - 1},
			Coordinates{0, 2, c.Depth - 1},
			Coordinates{0, 3, c.Depth - 1},
			Coordinates{0, 4, c.Depth - 1},
		)
	} else if c.X == 3 && c.Y == 2 {
		n = append(
			n,
			Coordinates{4, 0, c.Depth - 1}, // Left (inner level)
			Coordinates{4, 1, c.Depth - 1},
			Coordinates{4, 2, c.Depth - 1},
			Coordinates{4, 3, c.Depth - 1},
			Coordinates{4, 4, c.Depth - 1},
			Coordinates{c.X + 1, c.Y, c.Depth}, // Right
		)
	} else if c.X == 4 {
		n = append(
			n,
			Coordinates{c.X - 1, c.Y, c.Depth}, // Left
			Coordinates{3, 2, c.Depth + 1},     // Right (outer level)
		)
	} else {
		n = append(
			n,
			Coordinates{c.X - 1, c.Y, c.Depth}, // Left
			Coordinates{c.X + 1, c.Y, c.Depth}, // Right
		)
	}

	return
}

type State = set.Set[Coordinates]

func CountNeighbors(s State, c Coordinates) (count int) {
	for _, n := range RecursiveNeighbors(c) {
		if s.Has(n) {
			count++
		}
	}
	return
}

func Evolve(s State) (n State) {
	n = make(State)
	emptyNeighbors := make(set.Set[Coordinates])

	// Check if bugs die
	for bug := range s {
		bugNeighbors := 0

		for _, neighborCoord := range RecursiveNeighbors(bug) {
			if s.Has(neighborCoord) {
				bugNeighbors++
			} else {
				emptyNeighbors.Add(neighborCoord)
			}
		}

		if bugNeighbors == 1 {
			n.Add(bug)
		}
	}

	// Check if bug spawn
	for emptySpot := range emptyNeighbors {
		count := CountNeighbors(s, emptySpot)
		if count == 1 || count == 2 {
			n.Add(emptySpot)
		}
	}

	return
}

func ReadState(r io.Reader) (s State) {
	e := ReadEris(r)
	s = make(State)

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if e.Has(x, y) {
				s.Add(Coordinates{x, y, 0})
			}
		}
	}

	return
}

func SolveB(r io.Reader) any {
	s := ReadState(r)

	for i := 0; i < 200; i++ {
		s = Evolve(s)
	}

	return s.Len()
}
