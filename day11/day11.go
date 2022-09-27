package day11

import (
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
)

type Point struct {
	X, Y int
}

type Direction uint8

const (
	DirectionUp    = Direction(iota) // Towards positive Y
	DirectionRight                   // Towards positive X
	DirectionDown                    // Towards negative Y
	DirectionLeft                    // Towards negative X
)

type Color uint8

const (
	ColorBlack = Color(iota)
	ColorWhite
)

type Painter struct {
	Program  *intcode.Interpreter
	Colors   map[Point]Color
	Position Point
	Heading  Direction
}

func NewPainter(intcodeProgram io.Reader) (p *Painter) {
	p = &Painter{}
	p.Program = intcode.NewInterpreterNewIO(intcodeProgram)

	p.Colors = make(map[Point]Color)
	return p
}

func (p *Painter) Move() {
	switch p.Heading {
	case DirectionUp:
		p.Position.Y++
	case DirectionRight:
		p.Position.X++
	case DirectionDown:
		p.Position.Y--
	case DirectionLeft:
		p.Position.X--
	}
}

func (p *Painter) ExecAll() {
	// Launch the program
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Program.ExecAll()
	}()

	// Operate the painter
PainterLoop:
	for {
		// Send the color of the current tile
		select {
		case <-p.Program.Halted:
			break PainterLoop
		case p.Program.Input <- int(p.Colors[p.Position]):
			// Successfully sent the current tile - continue waiting for robot input
		}

		// Read whatever the robot has sent
		color := <-p.Program.Output
		rotation := <-p.Program.Output

		// Paint the panel
		p.Colors[p.Position] = Color(color)

		// Rotate the robot
		if rotation == 0 {
			// Rotate counter-clockwise
			switch p.Heading {
			case DirectionUp:
				p.Heading = DirectionLeft
			case DirectionRight:
				p.Heading = DirectionUp
			case DirectionDown:
				p.Heading = DirectionRight
			case DirectionLeft:
				p.Heading = DirectionDown
			default:
				panic("invalid painter heading")
			}
		} else if rotation == 1 {
			// Rotate clockwise
			switch p.Heading {
			case DirectionUp:
				p.Heading = DirectionRight
			case DirectionRight:
				p.Heading = DirectionDown
			case DirectionDown:
				p.Heading = DirectionLeft
			case DirectionLeft:
				p.Heading = DirectionUp
			default:
				panic("invalid painter heading")
			}
		} else {
			panic("invalid direction received")
		}

		p.Move()
	}

	// Wait for the program to finish
	close(p.Program.Input)
	wg.Wait()
}

func SolveA(r io.Reader) any {
	painter := NewPainter(r)
	painter.ExecAll()
	return len(painter.Colors)
}

func PrintImage(m map[Point]Color) {
	// Find image bounds
	left := math.MaxInt
	right := math.MinInt
	top := math.MinInt
	bottom := math.MaxInt
	for p := range m {
		if p.X < left {
			left = p.X
		}
		if p.X > right {
			right = p.X
		}
		if p.Y < bottom {
			bottom = p.Y
		}
		if p.Y > top {
			top = p.Y
		}
	}

	for y := top; y >= bottom; y-- {
		for x := left; x <= right; x++ {
			if m[Point{x, y}] == ColorWhite {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func SolveB(r io.Reader) any {
	painter := NewPainter(r)
	painter.Colors[Point{0, 0}] = ColorWhite
	painter.ExecAll()
	PrintImage(painter.Colors)
	return nil
}
