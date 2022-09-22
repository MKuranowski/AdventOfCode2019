package day11

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/input"
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
	Program     *intcode.Interpreter
	CameraInput io.Writer
	ProgOutput  io.Reader

	Colors   map[Point]Color
	Position Point
	Heading  Direction
}

func NewPainter(intcodeProgram io.Reader) (p *Painter) {
	p = &Painter{}
	p.Program = intcode.NewInterpreter(intcodeProgram)
	p.Program.Input, p.CameraInput = io.Pipe()
	p.ProgOutput, p.Program.Output = io.Pipe()

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
	for {
		// Send the color of the current tile
		err := input.SendInteger(int(p.Colors[p.Position]), p.CameraInput)
		if errors.Is(err, io.ErrClosedPipe) {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to send data to the intcode Interpreter: %w", err))
		}

		// Read whatever the robot has sent
		color, err := input.ReceiveInteger(p.ProgOutput)
		if err != nil {
			panic(fmt.Errorf("failed to read data from the intcode Interpreter: %w", err))
		}
		rotation, err := input.ReceiveInteger(p.ProgOutput)
		if err != nil {
			panic(fmt.Errorf("failed to read data from the intcode Interpreter: %w", err))
		}

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
