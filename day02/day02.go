package day02

import (
	"bufio"
	"io"
	"runtime"
	"strconv"
	"strings"
)

type IntCodeInterpreter struct {
	Memory []int
	IP     int
	Halted bool
}

func (i *IntCodeInterpreter) ExecOne() (more bool) {
	if i.Halted {
		return false
	}

	op := i.Memory[i.IP]

	// Special case for the halt
	if op == 99 {
		i.Halted = true
		return false
	}

	switch op {
	case 1:
		a := i.Memory[i.IP+1]
		b := i.Memory[i.IP+2]
		result := i.Memory[i.IP+3]
		i.Memory[result] = i.Memory[a] + i.Memory[b]

	case 2:
		a := i.Memory[i.IP+1]
		b := i.Memory[i.IP+2]
		result := i.Memory[i.IP+3]
		i.Memory[result] = i.Memory[a] * i.Memory[b]

	default:
		panic("unknown opcode in IntCodeInterpreter")
	}

	i.IP += 4
	return true
}

func (i *IntCodeInterpreter) ExecAll() {
	for i.ExecOne() {
	}
}

func (i *IntCodeInterpreter) Clone() *IntCodeInterpreter {
	return &IntCodeInterpreter{
		Memory: append([]int(nil), i.Memory...),
		IP:     i.IP,
		Halted: i.Halted,
	}
}

func NewIntCodeInterpreter(r io.Reader) *IntCodeInterpreter {
	i := &IntCodeInterpreter{}
	br := bufio.NewReader(r)

	eof := false
	for !eof {
		str, err := br.ReadString(',')
		eof = err != nil

		op, _ := strconv.Atoi(strings.TrimRight(str, "\n,"))
		i.Memory = append(i.Memory, op)
	}

	return i
}

func SolveA(r io.Reader) any {
	i := NewIntCodeInterpreter(r)
	i.Memory[1] = 12
	i.Memory[2] = 2
	i.ExecAll()
	return i.Memory[0]
}

type InputB struct {
	Noun, Verb int
}

func SolveB(r io.Reader) any {
	baseInterpreter := NewIntCodeInterpreter(r)
	ins := make(chan InputB)
	results := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for input := range ins {
				i := baseInterpreter.Clone()
				i.Memory[1] = input.Noun
				i.Memory[2] = input.Verb
				i.ExecAll()
				if i.Memory[0] == 19690720 {
					results <- 100*input.Noun + input.Verb
					break
				}
			}
		}()
	}

	for i := 0; i < 100*100; i++ {
		select {
		case ins <- InputB{i / 100, i % 100}:

		case result := <-results:
			close(ins)
			return result
		}
	}

	panic("no solution")
}
