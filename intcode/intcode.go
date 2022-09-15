package intcode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing/iotest"
)

func powerOfTen(x int) (result int) {
	result = 1
	for i := 0; i < x; i++ {
		result *= 10
	}
	return
}

type OPArgument interface {
	Get() int
	Set(int)
}

type MemoryReference struct {
	*int
}

func (r MemoryReference) Get() int  { return *r.int }
func (r MemoryReference) Set(x int) { *r.int = x }

type Immediate int

func (r Immediate) Get() int { return int(r) }
func (r Immediate) Set(int)  { panic("intcode.Interpreter - can't set an immediate value") }

type Interpreter struct {
	Memory []int
	IP     int
	Halted bool

	Input  io.Reader
	Output io.Writer
}

// getArgument figures out the correct parameter mode for a specific argument.
// The argument index is one-based.
func (i *Interpreter) getArgument(modes int, argIdx int, opSize int) OPArgument {
	// Do some math to extract the mode
	mode := (modes / powerOfTen(argIdx-1)) % 10

	switch mode {
	case 0:
		return MemoryReference{&i.Memory[i.Memory[i.IP+argIdx]]}
	case 1:
		return Immediate(i.Memory[i.IP+argIdx])
	default:
		panic("intcode.Interpreter - unsupported parameter mode")
	}
}

func (i *Interpreter) performIn() (x int) {
	_, err := fmt.Fscan(i.Input, &x)
	if err != nil {
		panic(fmt.Errorf("intcode.Interpreter - INPUT op failed: %w", err))
	}
	return x
}

func (i *Interpreter) performOut(x int) {
	_, err := fmt.Fprintln(i.Output, x)
	if err != nil {
		panic(fmt.Errorf("intcode.Interpreter - OUTPUT op failed: %w", err))
	}
}

func (i *Interpreter) ExecOne() (more bool) {
	if i.Halted {
		return false
	}

	modes, op := i.Memory[i.IP]/100, i.Memory[i.IP]%100
	opSize := 1

	// fmt.Fprintf(os.Stderr, "Executing OP %d from %d (modes %d)\n", op, i.IP, modes)

	switch op {
	case 1:
		// ADD
		opSize = 4
		src1 := i.getArgument(modes, 1, opSize)
		src2 := i.getArgument(modes, 2, opSize)
		dest := i.getArgument(modes, 3, opSize)

		dest.Set(src1.Get() + src2.Get())

	case 2:
		// MUL
		opSize = 4
		src1 := i.getArgument(modes, 1, opSize)
		src2 := i.getArgument(modes, 2, opSize)
		dest := i.getArgument(modes, 3, opSize)

		dest.Set(src1.Get() * src2.Get())

	case 3:
		// INPUT
		opSize = 2
		dest := i.getArgument(modes, 1, opSize)

		dest.Set(i.performIn())

	case 4:
		// OUTPUT
		opSize = 2
		src := i.getArgument(modes, 1, opSize)

		i.performOut(src.Get())

	case 99:
		// HALT
		i.Halted = true

	default:
		panic(fmt.Errorf("unknown opcode in IntCodeInterpreter: %d", op))
	}

	i.IP += opSize
	return !i.Halted
}

func (i *Interpreter) ExecAll() {
	for i.ExecOne() {
	}
}

func (i *Interpreter) Clone() *Interpreter {
	return &Interpreter{
		Memory: append([]int(nil), i.Memory...),
		IP:     i.IP,
		Halted: i.Halted,
	}
}

func NewInterpreter(program io.Reader) *Interpreter {
	return NewInterpreterWithIO(program, iotest.ErrReader(io.EOF), io.Discard)
}

func NewInterpreterWithIO(program io.Reader, input io.Reader, output io.Writer) *Interpreter {
	i := &Interpreter{Input: input, Output: output}
	br := bufio.NewReader(program)

	eof := false
	for !eof {
		str, err := br.ReadString(',')
		eof = err != nil

		op, _ := strconv.Atoi(strings.TrimRight(str, "\n,"))
		i.Memory = append(i.Memory, op)
	}

	return i
}
