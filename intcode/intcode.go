package intcode

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Interpreter struct {
	Memory []int
	IP     int
	Halted bool
}

func (i *Interpreter) ExecOne() (more bool) {
	if i.Halted {
		return false
	}

	op := i.Memory[i.IP]
	opSize := 1

	switch op {
	case 1:
		a := i.Memory[i.IP+1]
		b := i.Memory[i.IP+2]
		result := i.Memory[i.IP+3]

		i.Memory[result] = i.Memory[a] + i.Memory[b]
		opSize = 4

	case 2:
		a := i.Memory[i.IP+1]
		b := i.Memory[i.IP+2]
		result := i.Memory[i.IP+3]

		i.Memory[result] = i.Memory[a] * i.Memory[b]
		opSize = 4

	case 99:
		i.Halted = true

	default:
		panic("unknown opcode in IntCodeInterpreter")
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

func NewInterpreter(r io.Reader) *Interpreter {
	i := &Interpreter{}
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
