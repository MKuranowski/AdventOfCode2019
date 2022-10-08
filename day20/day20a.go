package day20

import (
	"fmt"
	"io"

	"github.com/MKuranowski/AdventOfCode2019/util/ascii"
	"github.com/MKuranowski/AdventOfCode2019/util/gheap"
	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

type Point struct{ X, Y int }

func GetNeighbors(p Point) [4]Point {
	return [4]Point{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}

func readLabel(pt Point, neighborIndex int, lines []string) string {
	label := [2]byte{}

	switch neighborIndex {
	case 0:
		label[0] = lines[pt.Y][pt.X-1]
		label[1] = lines[pt.Y][pt.X]
	case 1:
		label[0] = lines[pt.Y][pt.X]
		label[1] = lines[pt.Y][pt.X+1]
	case 2:
		label[0] = lines[pt.Y-1][pt.X]
		label[1] = lines[pt.Y][pt.X]
	case 3:
		label[0] = lines[pt.Y][pt.X]
		label[1] = lines[pt.Y+1][pt.X]
	default:
		panic("invalid neighbor index")
	}

	return string(label[:])
}

type Map struct {
	G map[Point][]Point

	Start Point
	End   Point
}

func ReadMap(r io.Reader) (m Map) {
	lines := input.ReadLines(r)
	labeledPoint := make(map[string][]Point)
	m.G = make(map[Point][]Point)

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
					m.G[pt] = append(m.G[pt], neighbor)
				} else if ascii.IsUpper(nType) {
					label = readLabel(neighbor, neighborIdx, lines)
				}
			}

			if label != "" {
				labeledPoint[label] = append(labeledPoint[label], pt)
			}
		}
	}

	for label, points := range labeledPoint {
		if label == "AA" {
			if len(points) != 1 {
				panic("multiple start points")
			}
			m.Start = points[0]
		} else if label == "ZZ" {
			if len(points) != 1 {
				panic("multiple end points")
			}
			m.End = points[0]
		} else {
			if len(points) != 2 {
				panic("portals must connect exactly 2 points")
			}
			m.G[points[0]] = append(m.G[points[0]], points[1])
			m.G[points[1]] = append(m.G[points[1]], points[0])
		}
	}

	return
}

func DumpGraph(g map[Point][]Point, w io.Writer) {
	for from, tos := range g {
		fmt.Fprintln(w, from, "->", tos)
	}
}

type dijkstraQueueEntry struct {
	At   Point
	Cost int
}

func ShortestPath(m Map) int {
	q := gheap.NewGenericHeap(func(a, b dijkstraQueueEntry) bool { return a.Cost < b.Cost })
	costs := map[Point]int{m.Start: 0}
	q.Push(dijkstraQueueEntry{m.Start, 0})

	for q.Len() > 0 {
		e := q.Pop()

		if e.At == m.End {
			return e.Cost
		}

		for _, neighbor := range m.G[e.At] {
			costTo := e.Cost + 1
			knownCost, ok := costs[neighbor]
			if ok && knownCost < costTo {
				continue
			}

			costs[neighbor] = costTo
			q.Push(dijkstraQueueEntry{neighbor, costTo})
		}
	}

	panic("no solution")
}

func SolveA(r io.Reader) any {
	m := ReadMap(r)
	// DumpGraph(m.G, os.Stderr)
	// fmt.Fprintln(os.Stderr, "Start:", m.Start)
	// fmt.Fprintln(os.Stderr, "End:", m.End)
	return ShortestPath(m)
}
