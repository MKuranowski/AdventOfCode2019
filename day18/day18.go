package day18

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"os"

	"github.com/MKuranowski/AdventOfCode2019/util/ascii"
	"github.com/MKuranowski/AdventOfCode2019/util/gheap"
)

type Point struct{ X, Y int }

func NeighborsOf(p Point) [4]Point {
	return [...]Point{
		{p.X - 1, p.Y - 1},
		{p.X + 1, p.Y - 1},
		{p.X - 1, p.Y + 1},
		{p.X + 1, p.Y + 1},
	}
}

type Map map[Point]byte

type KeySet uint32

func keyToIdx(k byte) int {
	if ascii.IsLower(k) {
		return int(k - 'a')
	} else if ascii.IsUpper(k) {
		return int(k - 'A')
	}
	panic("key must be a letter")
}

func (s KeySet) Has(k byte) bool {
	return s&(1<<keyToIdx(k)) != 0
}

func (s *KeySet) Add(k byte) {
	*s |= 1 << keyToIdx(k)
}

func (s KeySet) Len() int {
	return bits.OnesCount32(uint32(s))
}

func (s KeySet) IsSubset(o KeySet) bool {
	return s&o == s
}

func (s KeySet) IsSuperset(o KeySet) bool {
	return s&o == o
}

type MazeData struct {
	M     Map
	Start Point

	Keys  map[byte]Point
	Doors map[byte]Point
}

func LoadMaze(r io.Reader) (d MazeData) {
	d.M = make(Map)
	d.Keys = make(map[byte]Point)
	d.Doors = make(map[byte]Point)

	x, y := 0, 0
	br := bufio.NewReader(r)
	for {
		c, err := br.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to load maze: %w", err))
		}

		if c == '\n' {
			x = 0
			y++
		} else if c == '.' {
			d.M[Point{x, y}] = c
			x++
		} else if c == '@' {
			d.M[Point{x, y}] = '.'
			d.Start.X, d.Start.Y = x, y
			x++
		} else if ascii.IsUpper(c) {
			d.M[Point{x, y}] = c
			d.Doors[c] = Point{x, y}
			x++
		} else if ascii.IsLower(c) {
			d.M[Point{x, y}] = c
			d.Keys[c] = Point{x, y}
			x++
		} else {
			x++
		}
	}

	return
}

type MazeGraphEdge struct {
	Steps int
	Doors KeySet
}

type MazeGraph map[byte]map[byte]MazeGraphEdge

type simplifierQueueEntry struct {
	P     Point
	C     int
	Doors KeySet
}

func generateSimpleEdgesFrom(m MazeData, from Point) (t map[byte]MazeGraphEdge) {
	t = make(map[byte]MazeGraphEdge)

	costs := make(map[Point]int)
	costs[from] = 0

	q := gheap.NewGenericHeap(func(a, b simplifierQueueEntry) bool { return a.C < b.C })
	q.Push(simplifierQueueEntry{from, 0, 0})

	for q.Len() > 0 {
		elem := q.Pop()
		elemType := m.M[elem.P]

		if ascii.IsLower(elemType) {
			t[elemType] = MazeGraphEdge{elem.C, elem.Doors}
		}

		for _, neighbor := range NeighborsOf(elem.P) {
			neighborType := m.M[neighbor]
			if neighborType == 0 {
				continue
			}

			costToNeighbor := elem.C + 1
			knownCost, visited := costs[neighbor]
			if visited && knownCost < costToNeighbor {
				continue
			}

			newDoors := elem.Doors
			if ascii.IsUpper(neighborType) {
				newDoors.Add(neighborType)
			}

			costs[neighbor] = costToNeighbor
			q.Push(simplifierQueueEntry{neighbor, costToNeighbor, newDoors})
		}
	}

	return t
}

func MazeToSimpleGraph(m MazeData) (g MazeGraph) {
	g = make(MazeGraph)
	g['@'] = generateSimpleEdgesFrom(m, m.Start)
	for key, pt := range m.Keys {
		g[key] = generateSimpleEdgesFrom(m, pt)
	}
	return
}

type shortestPathQueueEntryHash struct {
	Keys KeySet
	At   byte
}

type shortestPathQueueEntry struct {
	shortestPathQueueEntryHash
	Cost int
}

func FindShortestPath(g MazeGraph, allKeysCount int) int {
	costs := make(map[shortestPathQueueEntryHash]int)
	queue := gheap.NewGenericHeap(func(a, b shortestPathQueueEntry) bool { return a.Cost < b.Cost })

	queue.Push(shortestPathQueueEntry{shortestPathQueueEntryHash{0, '@'}, 0})
	costs[shortestPathQueueEntryHash{0, '@'}] = 0

	for queue.Len() > 0 {
		e := queue.Pop()

		if e.Keys.Len() == allKeysCount {
			return e.Cost
		}

		for neighbor, edge := range g[e.At] {
			if !e.Keys.IsSuperset(edge.Doors) {
				// Don't have the keys to open doors en-route to neighbor!
				continue
			}

			n := shortestPathQueueEntry{
				shortestPathQueueEntryHash{e.Keys, neighbor},
				e.Cost + e.Cost,
			}
			n.Keys.Add(neighbor)

			// Check if a cheaper route to this state exists
			knownCost, ok := costs[n.shortestPathQueueEntryHash]
			if ok && knownCost < n.Cost {
				continue
			}

			// Push the new state onto the queue
			costs[n.shortestPathQueueEntryHash] = n.Cost
			queue.Push(n)
		}
	}

	panic("no solution")
}

func DumpGraph(g MazeGraph, w io.Writer) {
	for from, edges := range g {
		fmt.Fprintf(w, "%c\t->:\n", from)
		for to, edge := range edges {
			fmt.Fprintf(w, "\t%c\t%d\n", to, edge.Steps)
		}
	}
}

func SolveA(r io.Reader) any {
	m := LoadMaze(r)
	g := MazeToSimpleGraph(m)
	DumpGraph(g, os.Stderr)
	return FindShortestPath(g, len(m.Keys))
}
