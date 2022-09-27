package day13

import (
	"fmt"
	"io"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
)

type TileType uint8

const (
	TileEmpty = TileType(iota)
	TileWall
	TileBlock
	TilePaddle
	TileBall
)

func (t TileType) AsChar() byte {
	switch t {
	case TileEmpty:
		return ' '
	case TileWall:
		return '#'
	case TileBlock:
		return 'x'
	case TilePaddle:
		return '-'
	case TileBall:
		return 'o'
	default:
		panic(fmt.Errorf("invalid TileType value: %d", t))
	}
}

func receiveScreenCommand(ch <-chan int) (col, row, data int, closed bool) {
	col = <-ch
	row = <-ch
	data, closed = <-ch
	return
}

type Screen = [][]TileType

func DummyScreen(r <-chan int) (s Screen) {
	for {
		// Receive the data
		col, row, data, ok := receiveScreenCommand(r)
		if !ok {
			break
		}

		// Expand the amount of rows if necessary
		if row >= len(s) {
			ns := make(Screen, row+1)
			copy(ns, s)
			s = ns
		}

		// Expand the amount of cols in s[row] if necessary
		if col >= len(s[row]) {
			nr := make([]TileType, col+1)
			copy(nr, s[row])
			s[row] = nr
		}

		// Set the block
		s[row][col] = TileType(data)
	}
	return
}

func SolveA(r io.Reader) any {
	// Prepare the interpreter
	screenInput := make(chan int)
	i := intcode.NewInterpreterWithIO(r, nil, screenInput)

	// Launch the interpreter and the screen
	go i.ExecAll()
	endScreen := DummyScreen(screenInput)

	// Count the amount of blocks
	blocks := 0
	for _, row := range endScreen {
		for _, tile := range row {
			if tile == TileBlock {
				blocks++
			}
		}
	}

	return blocks
}

type Arcade struct {
	S            Screen
	BallColumn   int
	PaddleColumn int
	Score        int

	Draw     <-chan int
	Joystick chan<- int
	Halt     <-chan struct{}
}

func (a *Arcade) Dump(w io.Writer) {
	fmt.Fprintln(w, "=======")
	fmt.Fprintf(w, "Score: %d\n", a.Score)
	for _, row := range a.S {
		for _, tile := range row {
			fmt.Fprintf(w, "%c", tile.AsChar())
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "=======")
}

func (a *Arcade) continueDraw(col int) {
	// Receive the rest of the draw command
	row := <-a.Draw
	data, ok := <-a.Draw
	if !ok {
		return
	}

	// Special case for the score
	if col == -1 {
		a.Score = data
		return
	}

	// Expand the amount of rows if necessary
	if row >= len(a.S) {
		ns := make(Screen, row+1)
		copy(ns, a.S)
		a.S = ns
	}

	// Expand the amount of cols in s[row] if necessary
	if col >= len(a.S[row]) {
		nr := make([]TileType, col+1)
		copy(nr, a.S[row])
		a.S[row] = nr
	}

	// Set the block
	tile := TileType(data)
	a.S[row][col] = tile

	if tile == TileBall {
		a.BallColumn = col
	} else if tile == TilePaddle {
		a.PaddleColumn = col
	}
}

func (a *Arcade) getJoystickInput() int {
	if a.PaddleColumn > a.BallColumn {
		return -1
	} else if a.PaddleColumn < a.BallColumn {
		return 1
	}
	return 0
}

func (a *Arcade) Run() {
	for {
		select {
		case <-a.Halt:
			return
		case col, ok := <-a.Draw:
			if !ok {
				return
			}
			a.continueDraw(col)
		case a.Joystick <- a.getJoystickInput():
		}
	}
}

func SolveB(r io.Reader) any {
	// Prepare the interpreter
	i := intcode.NewInterpreterNewIO(r)
	a := Arcade{
		Draw:     i.Output,
		Joystick: i.Input,
		Halt:     i.Halted,
	}

	// "Hack the amount of coins"
	i.Memory[0] = 2

	// Launch the interpreter and the screen
	go i.ExecAll()
	a.Run()
	close(i.Input)

	return a.Score
}
