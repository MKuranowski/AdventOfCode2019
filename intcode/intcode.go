package intcode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/util/deque"
)

var (
	ErrInputOverClosed = errors.New("input over closed channel")
	ErrInputOverNil    = errors.New("input over nil channel")
	ErrOutputOverNil   = errors.New("output over nil channel")
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

func NewMemoryReference(i *Interpreter, idx int) MemoryReference {
	if idx >= len(i.Memory) {
		newMemory := make([]int, idx+1)
		copy(newMemory, i.Memory)
		i.Memory = newMemory
	}

	return MemoryReference{&i.Memory[idx]}
}

func NewMemoryReferenceSync(i *SyncInterpreter, idx int) MemoryReference {
	if idx >= len(i.Memory) {
		newMemory := make([]int, idx+1)
		copy(newMemory, i.Memory)
		i.Memory = newMemory
	}

	return MemoryReference{&i.Memory[idx]}
}

func (r MemoryReference) Get() int  { return *r.int }
func (r MemoryReference) Set(x int) { *r.int = x }

type Immediate int

func (r Immediate) Get() int { return int(r) }
func (r Immediate) Set(int)  { panic("intcode.Interpreter - can't set an immediate value") }

type Interpreter struct {
	Memory       []int
	IP           int
	Halted       chan struct{}
	RelativeBase int

	Input  chan int
	Output chan int
}

// getArgument figures out the correct parameter mode for a specific argument.
// The argument index is one-based.
func (i *Interpreter) getArgument(modes int, argIdx int) OPArgument {
	// Do some math to extract the mode
	mode := (modes / powerOfTen(argIdx-1)) % 10

	switch mode {
	case 0:
		return NewMemoryReference(i, i.Memory[i.IP+argIdx])
	case 1:
		return Immediate(i.Memory[i.IP+argIdx])
	case 2:
		return NewMemoryReference(i, i.RelativeBase+i.Memory[i.IP+argIdx])
	default:
		panic("intcode.Interpreter - unsupported parameter mode")
	}
}

func (i *Interpreter) performIn() int {
	if i.Input == nil {
		panic(ErrInputOverNil)
	}
	x, ok := <-i.Input
	if !ok {
		panic(ErrInputOverClosed)
	}
	return x
}

func (i *Interpreter) performOut(x int) {
	if i.Output == nil {
		panic(ErrOutputOverNil)
	}
	i.Output <- x
}

func (i *Interpreter) ExecOne() (more bool) {
	if i.IsHalted() {
		return false
	}

	modes, op := i.Memory[i.IP]/100, i.Memory[i.IP]%100
	opSize := 1

	// fmt.Fprintf(os.Stderr, "Executing OP %d from %d (modes %d)\n", op, i.IP, modes)

	switch op {
	case 1:
		// ADD
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		dest.Set(src1.Get() + src2.Get())

	case 2:
		// MUL
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		dest.Set(src1.Get() * src2.Get())

	case 3:
		// INPUT
		opSize = 2
		dest := i.getArgument(modes, 1)

		dest.Set(i.performIn())

	case 4:
		// OUTPUT
		opSize = 2
		src := i.getArgument(modes, 1)

		i.performOut(src.Get())

	case 5:
		// JUMP-IF-TRUE
		opSize = 3
		src := i.getArgument(modes, 1)
		dest := i.getArgument(modes, 2)

		if src.Get() != 0 {
			i.IP = dest.Get()
			opSize = 0
		}

	case 6:
		// JUMP-IF-FALSE
		opSize = 3
		src := i.getArgument(modes, 1)
		dest := i.getArgument(modes, 2)

		if src.Get() == 0 {
			i.IP = dest.Get()
			opSize = 0
		}

	case 7:
		// LESS-THAN
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		if src1.Get() < src2.Get() {
			dest.Set(1)
		} else {
			dest.Set(0)
		}

	case 8:
		// EQ
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		if src1.Get() == src2.Get() {
			dest.Set(1)
		} else {
			dest.Set(0)
		}

	case 9:
		// ADJUST RELATIVE BASE
		opSize = 2
		src := i.getArgument(modes, 1)
		i.RelativeBase += src.Get()

	case 99:
		// HALT
		close(i.Halted)

	default:
		panic(fmt.Errorf("unknown opcode in IntCodeInterpreter: %d", op))
	}

	i.IP += opSize
	return !i.IsHalted()
}

func (i *Interpreter) ExecAll() {
	for i.ExecOne() {
	}

	if i.Output != nil {
		close(i.Output)
	}
}

func (i *Interpreter) Clone() (new *Interpreter) {
	new = &Interpreter{
		Memory: append([]int(nil), i.Memory...),
		IP:     i.IP,
		Halted: make(chan struct{}),
	}
	if i.IsHalted() {
		close(new.Halted)
	}
	return
}

func (i *Interpreter) IsHalted() bool {
	select {
	case <-i.Halted:
		return true
	default:
		return false
	}
}

func NewInterpreter(program io.Reader) *Interpreter {
	return NewInterpreterWithIO(program, nil, nil)
}

func NewInterpreterNewIO(program io.Reader) *Interpreter {
	return NewInterpreterWithIO(program, make(chan int), make(chan int))
}

func NewInterpreterWithIO(program io.Reader, input, output chan int) *Interpreter {
	i := &Interpreter{Input: input, Output: output, Halted: make(chan struct{})}
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

// Required for day 23, thought this will be unnecessary :<

type SyncExecutionState uint8

const (
	SyncExecutionStateReady = SyncExecutionState(iota)
	SyncExecutionStateHalted
	SyncExecutionStateBlockedOnInput
)

type SyncInterpreter struct {
	Memory       []int
	IP           int
	RelativeBase int

	Halted        bool
	Input, Output deque.Deque[int]
}

// getArgument figures out the correct parameter mode for a specific argument.
// The argument index is one-based.
func (i *SyncInterpreter) getArgument(modes int, argIdx int) OPArgument {
	// Do some math to extract the mode
	mode := (modes / powerOfTen(argIdx-1)) % 10

	switch mode {
	case 0:
		return NewMemoryReferenceSync(i, i.Memory[i.IP+argIdx])
	case 1:
		return Immediate(i.Memory[i.IP+argIdx])
	case 2:
		return NewMemoryReferenceSync(i, i.RelativeBase+i.Memory[i.IP+argIdx])
	default:
		panic("intcode.SyncInterpreter - unsupported parameter mode")
	}
}

func (i *SyncInterpreter) performIn() (int, SyncExecutionState) {
	if i.Input.Len() == 0 {
		return 0, SyncExecutionStateBlockedOnInput
	}
	return i.Input.PopFront(), SyncExecutionStateReady
}

func (i *SyncInterpreter) performOut(x int) {
	i.Output.PushBack(x)
}

func (i *SyncInterpreter) ExecOne() SyncExecutionState {
	if i.Halted {
		return SyncExecutionStateHalted
	}

	modes, op := i.Memory[i.IP]/100, i.Memory[i.IP]%100
	opSize := 1

	// fmt.Fprintf(os.Stderr, "Executing OP %d from %d (modes %d)\n", op, i.IP, modes)

	switch op {
	case 1:
		// ADD
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		dest.Set(src1.Get() + src2.Get())

	case 2:
		// MUL
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		dest.Set(src1.Get() * src2.Get())

	case 3:
		// INPUT
		opSize = 2
		dest := i.getArgument(modes, 1)
		value, state := i.performIn()
		if state == SyncExecutionStateBlockedOnInput {
			return SyncExecutionStateBlockedOnInput
		}
		dest.Set(value)

	case 4:
		// OUTPUT
		opSize = 2
		src := i.getArgument(modes, 1)
		i.performOut(src.Get())

	case 5:
		// JUMP-IF-TRUE
		opSize = 3
		src := i.getArgument(modes, 1)
		dest := i.getArgument(modes, 2)

		if src.Get() != 0 {
			i.IP = dest.Get()
			opSize = 0
		}

	case 6:
		// JUMP-IF-FALSE
		opSize = 3
		src := i.getArgument(modes, 1)
		dest := i.getArgument(modes, 2)

		if src.Get() == 0 {
			i.IP = dest.Get()
			opSize = 0
		}

	case 7:
		// LESS-THAN
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		if src1.Get() < src2.Get() {
			dest.Set(1)
		} else {
			dest.Set(0)
		}

	case 8:
		// EQ
		opSize = 4
		src1 := i.getArgument(modes, 1)
		src2 := i.getArgument(modes, 2)
		dest := i.getArgument(modes, 3)

		if src1.Get() == src2.Get() {
			dest.Set(1)
		} else {
			dest.Set(0)
		}

	case 9:
		// ADJUST RELATIVE BASE
		opSize = 2
		src := i.getArgument(modes, 1)
		i.RelativeBase += src.Get()

	case 99:
		// HALT
		opSize = 0
		i.Halted = true

	default:
		panic(fmt.Errorf("unknown opcode in IntCodeInterpreter: %d", op))
	}

	i.IP += opSize
	if i.Halted {
		return SyncExecutionStateHalted
	}
	return SyncExecutionStateReady
}

func (i *SyncInterpreter) ExecAll() SyncExecutionState {
	state := SyncExecutionStateReady
	for state == SyncExecutionStateReady {
		state = i.ExecOne()
	}
	return state
}

func (i *SyncInterpreter) Clone() (new *SyncInterpreter) {
	new = &SyncInterpreter{
		Memory:       append([]int(nil), i.Memory...),
		IP:           i.IP,
		RelativeBase: i.RelativeBase,
		Halted:       i.Halted,
		Input:        deque.NewDeque[int](),
		Output:       deque.NewDeque[int](),
	}
	return
}

func NewSyncInterpreter(program io.Reader) *SyncInterpreter {
	i := NewInterpreter(program)
	s := &SyncInterpreter{
		Memory:       i.Memory,
		IP:           0,
		RelativeBase: 0,
		Halted:       false,
		Input:        deque.NewDeque[int](),
		Output:       deque.NewDeque[int](),
	}
	return s
}
