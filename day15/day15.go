package day15

import (
	"fmt"
	"io"
	"math"
	"math/rand"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/gheap"
	"github.com/MKuranowski/AdventOfCode2019/util/intmath"
	"github.com/MKuranowski/AdventOfCode2019/util/maps2"
	"golang.org/x/exp/maps"
)

type MapTile uint8

const (
	MapTileUnknown MapTile = iota
	MapTileWall
	MapTileCorridor
	MapTileOxygen
)

type Point struct {
	X, Y int
}

func (p Point) AfterDecision(d Decision) Point {
	switch d {
	case DecisionNorth:
		return Point{p.X, p.Y + 1}
	case DecisionSouth:
		return Point{p.X, p.Y - 1}
	case DecisionEast:
		return Point{p.X + 1, p.Y}
	case DecisionWest:
		return Point{p.X - 1, p.Y}
	default:
		panic("invalid decision to apply for a point")
	}
}

type ControllerState struct {
	Pos    Point
	Map    map[Point]MapTile
	Oxygen Point
}

type Decision uint8

const (
	DecisionHalt Decision = iota
	DecisionNorth
	DecisionSouth
	DecisionWest
	DecisionEast
)

type Decider interface {
	Decide(ControllerState) Decision
}

type Controller struct {
	ControllerState
	I *intcode.Interpreter
	D Decider
}

func (c *Controller) Run() {
	if c.Map == nil {
		c.Map = make(map[Point]MapTile)
	}

	// Launch the interpreter
	c.I.Input = make(chan int)
	c.I.Output = make(chan int)
	go c.I.ExecAll()

	for {
		// Get a decision from the Decider
		move := c.D.Decide(c.ControllerState)
		if move == DecisionHalt {
			break
		}

		// Send the decision
		c.I.Input <- int(move)

		// Update the state
		to := c.Pos.AfterDecision(move)
		result := <-c.I.Output
		switch result {
		case 0:
			c.Map[to] = MapTileWall
		case 1:
			c.Map[to] = MapTileCorridor
			c.Pos = to
		case 2:
			c.Map[to] = MapTileOxygen
			c.Oxygen = to
			c.Pos = to
		default:
			panic(fmt.Errorf("invalid result of move: %d", result))
		}
	}
}

type RandomWalker struct {
	Steps int
}

func (w *RandomWalker) Decide(ControllerState) Decision {
	if w.Steps == 0 {
		return DecisionHalt
	}
	w.Steps--
	return Decision(1 + rand.Intn(4))
}

func ShowMap(m map[Point]MapTile, on io.Writer) {
	// Find the bounds
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	for p := range m {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Start dumping the map
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			switch m[Point{x, y}] {
			case MapTileWall:
				fmt.Fprintf(on, "#")
			case MapTileCorridor:
				if x == 0 && y == 0 {
					fmt.Fprintf(on, "x")
				} else {
					fmt.Fprintf(on, ".")
				}
			case MapTileOxygen:
				fmt.Fprintf(on, "*")
			default:
				fmt.Fprint(on, " ")
			}
		}
		fmt.Fprintln(on)
	}
}

type aStarQElement struct {
	Pt     Point
	CostTo int
	Score  int
}

func reconstructPath(cameFrom map[Point]Point, from, to Point) (way []Point) {
	way = append(way, to)
	for from != to {
		to = cameFrom[to]
		way = append(way, to)
	}
	return
}

func ShortestPath(m map[Point]MapTile, from Point, to Point) []Point {
	q := gheap.NewGenericHeap(func(a, b aStarQElement) bool { return a.Score < b.Score })
	cheapestTo := make(map[Point]int)
	cameFrom := make(map[Point]Point)

	q.Push(aStarQElement{from, 0, intmath.DistManhattan(from.X, from.Y, to.X, to.Y)})

	for q.Len() > 0 {
		toExpand := q.Pop()
		if toExpand.Pt == to {
			return reconstructPath(cameFrom, from, to)
		}

		// Iterate over of the neighbors
		for _, d := range []Decision{DecisionNorth, DecisionSouth, DecisionWest, DecisionEast} {
			// Check if the neighbor is accessible
			toPt := toExpand.Pt.AfterDecision(d)
			toPtType := m[toPt]
			if toPtType == MapTileUnknown || toPtType == MapTileWall {
				continue
			}

			// Check if there already exists a cheaper route to `toPt`
			newCost := toExpand.CostTo + 1
			currentCost, alreadyVisited := cheapestTo[toPt]
			if alreadyVisited && currentCost < newCost {
				continue
			}

			// Push `toPt` onto the queue
			newScore := newCost + intmath.DistManhattan(toPt.X, toPt.Y, to.X, to.Y)
			q.Push(aStarQElement{toPt, newCost, newScore})
			cheapestTo[toPt] = newCost
			cameFrom[toPt] = toExpand.Pt
		}
	}

	panic("no path found")
}

func GetMap(r io.Reader) (m map[Point]MapTile, oxygen Point) {
	i := intcode.NewInterpreter(r)

	// Do a random walk of 1 million steps to map out the maze.
	// That's a pretty stupid strategy, but only takes ~2 seconds and works, lol.
	d := &RandomWalker{1_000_000}
	c := &Controller{I: i.Clone(), D: d}
	c.Run()

	return c.Map, c.Oxygen
}

func SolveA(r io.Reader) any {
	m, oxygen := GetMap(r)
	path := ShortestPath(m, Point{0, 0}, oxygen)
	return len(path) - 1
}

func SolveB(r io.Reader) any {
	m, _ := GetMap(r)
	tilesToFill := maps2.CountValues(m, MapTileCorridor)
	rounds := 0
	for tilesToFill > 0 {
		m2 := maps.Clone(m)
		for pt, tile := range m {
			if tile == MapTileOxygen {
				// Iterate over of the neighbors
				for _, d := range []Decision{DecisionNorth, DecisionSouth, DecisionWest, DecisionEast} {
					// Check if the neighbor is accessible
					toPt := pt.AfterDecision(d)
					toState := m2[toPt]
					if toState == MapTileCorridor {
						m2[toPt] = MapTileOxygen
						tilesToFill--
					}
				}
			}
		}

		m = m2
		rounds++
	}

	return rounds
}
