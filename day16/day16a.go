package day16

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/MKuranowski/AdventOfCode2019/util/intmath"
	"github.com/MKuranowski/AdventOfCode2019/util/matrix"
)

func RotateLeft(x []int) {
	t := x[0]
	copy(x, x[1:])
	x[len(x)-1] = t
}

func GenerateTransformTable(digits int) (t matrix.Matrix[int]) {
	t = matrix.New[int](digits, digits)
	streakDigits := [...]int{0, 1, 0, -1}

	for destIdx := 0; destIdx < digits; destIdx++ {
		streakLen := destIdx + 1

		for srcIdx := 0; srcIdx < digits; srcIdx++ {
			streakIdx := (srcIdx + 1) / streakLen
			t.Set(srcIdx, destIdx, streakDigits[streakIdx%4])
		}
	}
	return
}

func RunPhase(msg, transformTable matrix.Matrix[int]) matrix.Matrix[int] {
	r := msg.MatMul(transformTable)
	r.Apply(func(i int) int { return intmath.Abs(i % 10) })
	return r
}

func ReadMessage(r io.Reader) matrix.Matrix[int] {
	br := bufio.NewReader(r)
	msg := []uint8(nil)
	for {
		c, err := br.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to read message: %w", err))
		} else if c == '\n' {
			break
		} else if c < '0' || c > '9' {
			panic("non-digit byte in message")
		} else {
			msg = append(msg, c-'0')
		}
	}

	m := matrix.New[int](1, len(msg))
	for i, c := range msg {
		m.Set(0, i, int(c))
	}
	return m
}

func SolveA(r io.Reader) any {
	msg := ReadMessage(r)
	transformTable := GenerateTransformTable(msg.Width())

	for i := 0; i < 100; i++ {
		msg = RunPhase(msg, transformTable)
	}

	// Get the first 8 digits of the message
	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[i] = byte(msg.Get(0, i)) + '0'
	}
	return string(result)
}
