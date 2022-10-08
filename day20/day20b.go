package day20

import (
	"container/heap"
	"fmt"
	"io"
	"math"

	"github.com/MKuranowski/AdventOfCode2019/util/ascii"
	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

type RecursiveMapEdge struct {
	To         Point
	LevelDelta int
}

type RecursiveMap struct {
	G map[Point][]RecursiveMapEdge

	Start, End Point
}

type pointAndSide struct {
	P       Point
	IsOuter bool
}

func bounds(l []string) (maxX, maxY int) {
	maxX = math.MinInt
	for _, r := range l {
		x := len(r) - 1
		if x > maxX {
			maxX = x
		}
	}
	maxY = len(l) - 1
	return
}

func ReadRecursiveMap(r io.Reader) (m RecursiveMap) {
	lines := input.ReadLines(r)
	labeledPoint := make(map[string][]pointAndSide)
	m.G = make(map[Point][]RecursiveMapEdge)

	maxOuterX, maxOuterY := bounds(lines)
	maxOuterY -= 3
	maxOuterX -= 3

	for y, row := range lines {
		for x, c := range row {
			if c != '.' {
				continue
			}

			pt := Point{x, y}
			label := ""

			for neighborIdx, neighbor := range GetNeighbors(pt) {
				nType := lines[neighbor.Y][neighbor.X]
				if nType == '.' {
					m.G[pt] = append(m.G[pt], RecursiveMapEdge{neighbor, 0})
				} else if ascii.IsUpper(nType) {
					label = readLabel(neighbor, neighborIdx, lines)
				}
			}

			if label != "" {
				isOuter := pt.X < 3 || pt.Y < 3 || pt.X > maxOuterX || pt.Y > maxOuterY
				labeledPoint[label] = append(labeledPoint[label], pointAndSide{pt, isOuter})
			}
		}
	}

	for label, points := range labeledPoint {
		if label == "AA" {
			if len(points) != 1 {
				panic("multiple start points")
			}
			m.Start = points[0].P
		} else if label == "ZZ" {
			if len(points) != 1 {
				panic("multiple end points")
			}
			m.End = points[0].P
		} else {
			if len(points) != 2 {
				panic("portals must connect exactly 2 points")
			} else if points[0].IsOuter == points[1].IsOuter {
				panic(fmt.Errorf("portals %s must connect points on outer and inner sides: %#v and %#v", label, points[0], points[1]))
			}

			// Ensure points[0] is the outer point
			if points[1].IsOuter {
				points[0], points[1] = points[1], points[0]
			}

			m.G[points[0].P] = append(m.G[points[0].P], RecursiveMapEdge{points[1].P, -1})
			m.G[points[1].P] = append(m.G[points[1].P], RecursiveMapEdge{points[0].P, 1})
		}
	}

	return
}

// rdqe = recursiveDijkstraQueueEntry
type rdqe struct {
	rdqeHash
	Cost, Index int
}

type rdqeHash struct {
	At    Point
	Level int
}

type rdq []*rdqe

func (q rdq) Len() int { return len(q) }

func (q rdq) Less(i, j int) bool { return q[i].Cost < q[j].Cost }

func (q rdq) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Index = i
	q[j].Index = j
}

func (q *rdq) Push(x any) {
	item := x.(*rdqe)
	item.Index = len(*q)
	*q = append(*q, item)
}

func (q *rdq) Pop() any {
	l := len(*q)
	r := (*q)[l-1]
	(*q)[l-1] = nil
	*q = (*q)[:l-1]
	r.Index = -1
	return r
}

const LevelLimit = 50

func RecursiveShortestPath(m RecursiveMap) int {
	q := &rdq{&rdqe{rdqeHash{m.Start, 0}, 0, 0}}
	entries := map[rdqeHash]*rdqe{(*q)[0].rdqeHash: (*q)[0]}

	for q.Len() > 0 {
		e := heap.Pop(q).(*rdqe)

		if e.At == m.End && e.Level == 0 {
			return e.Cost
		}

		for _, neighbor := range m.G[e.At] {
			n := &rdqe{rdqeHash{neighbor.To, e.Level + neighbor.LevelDelta}, e.Cost + 1, -1}
			if n.Level < 0 || n.Level > LevelLimit {
				continue
			}

			existing := entries[n.rdqeHash]
			if existing != nil && existing.Cost < n.Cost {
				continue
			} else if existing != nil {
				heap.Remove(q, existing.Index)
			}

			entries[n.rdqeHash] = n
			heap.Push(q, n)
		}
	}

	panic("no solution")
}

func SolveB(r io.Reader) any {
	m := ReadRecursiveMap(r)
	return RecursiveShortestPath(m)
}
