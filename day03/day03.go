package day03

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

type SegmentDirection uint8

const (
	SegmentDirectionHorizontal = SegmentDirection(iota)
	SegmentDirectionVertical
	SegmentDirectionDiagonal
)

func det(a, b Point) int {
	return a.X*b.Y - a.Y*b.X
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sortedPair(a, b int) (int, int) {
	if b > a {
		return a, b
	}
	return b, a
}

// Point represents 2 numbers od a 2D plane.
//
// X represents the horizontal axis,
// Y represents the vertical axis.
type Point struct {
	X, Y int
}

func (p Point) String() string { return fmt.Sprintf("(%d, %d)", p.X, p.Y) }

func (p1 Point) ManhattanDistance(p2 Point) int {
	return abs(p2.X-p1.X) + abs(p2.Y-p1.Y)
}

type Segment struct {
	From, To Point
}

func (s Segment) String() string { return fmt.Sprintf("%v -> %v", s.From, s.To) }

func (s Segment) Direction() SegmentDirection {
	if s.From.X == s.To.X {
		return SegmentDirectionVertical
	} else if s.From.Y == s.To.Y {
		return SegmentDirectionHorizontal
	}
	return SegmentDirectionDiagonal
}

func (s Segment) ManhattanDistance() int {
	return s.From.ManhattanDistance(s.To)
}

func (s Segment) Contains(p Point) bool {
	xMin, xMax := sortedPair(s.From.X, s.To.X)
	yMin, yMax := sortedPair(s.From.Y, s.To.Y)

	return p.X >= xMin && p.X <= xMax && p.Y >= yMin && p.Y <= yMax
}

// Checks if 2 segments intersect
//
// Source: https://stackoverflow.com/a/20677983
// Under CC-BY-SA 4.0, by Paul Draper
// Modified from Python to Go and to check whether _segments_ intersect
func (s1 Segment) Intersects(s2 Segment) (Point, bool) {
	deltaX := Point{s1.From.X - s1.To.X, s2.From.X - s2.To.X}
	deltaY := Point{s1.From.Y - s1.To.Y, s2.From.Y - s2.To.Y}

	div := det(deltaX, deltaY)
	if div == 0 {
		// No intersection
		return Point{0, 0}, false
	}

	d := Point{det(s1.From, s1.To), det(s2.From, s2.To)}
	intersection := Point{det(d, deltaX) / div, det(d, deltaY) / div}

	// Check whether both lines contain the intersection
	if s1.Contains(intersection) && s2.Contains(intersection) {
		return intersection, true
	}

	return Point{0, 0}, false
}

type Wire []Segment

func ParseWire(line string) Wire {
	w := Wire{}
	start := Point{0, 0}
	end := Point{0, 0}

	for _, segmentString := range strings.Split(line, ",") {
		dir := segmentString[0]
		steps, _ := strconv.Atoi(segmentString[1:])

		switch dir {
		case 'R':
			end = Point{start.X + steps, start.Y}

		case 'L':
			end = Point{start.X - steps, start.Y}

		case 'U':
			end = Point{start.X, start.Y + steps}

		case 'D':
			end = Point{start.X, start.Y - steps}

		default:
			panic(fmt.Errorf("invalid direction"))
		}

		w = append(w, Segment{start, end})
		start = end
	}

	return w
}

func (w Wire) StepsTo(p Point) (steps int) {
	for _, s := range w {
		if s.Contains(p) {
			steps += s.From.ManhattanDistance(p)
			break
		} else {
			steps += s.ManhattanDistance()
		}
	}
	return
}

func SolveA(r io.Reader) any {
	// Load wires
	wires := []Wire{}
	inLines := input.NewLineIterator(r)
	for inLines.Next() {
		wires = append(wires, ParseWire(inLines.Get()))
	}

	// Iterate over every pair of segments of every pair of wires
	closest := math.MaxInt
	for i, w1 := range wires {
		for j, w2 := range wires {
			// Don't intersect a wire with itself
			if i == j {
				continue
			}

			for _, s1 := range w1 {
				for _, s2 := range w2 {
					intersection, ok := s1.Intersects(s2)
					intersectionDistance := intersection.ManhattanDistance(Point{0, 0})
					if ok && intersectionDistance > 0 && intersectionDistance < closest {
						closest = intersectionDistance
					}
				}
			}
		}
	}

	return closest
}

func SolveB(r io.Reader) any {
	// Load wires
	wires := []Wire{}
	inLines := input.NewLineIterator(r)
	for inLines.Next() {
		wires = append(wires, ParseWire(inLines.Get()))
	}

	// Iterate over every pair of segments of every pair of wires
	cheapest := math.MaxInt
	for i, w1 := range wires {
		for j, w2 := range wires {
			// Don't intersect a wire with itself
			if i == j {
				continue
			}

			for _, s1 := range w1 {
				for _, s2 := range w2 {
					intersection, ok := s1.Intersects(s2)
					if ok {
						cost := w1.StepsTo(intersection) + w2.StepsTo(intersection)
						if cost > 0 && cost < cheapest {
							cheapest = cost
						}
					}
				}
			}
		}
	}

	return cheapest
}
