package day09

import (
	"io"
	"os"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
)

func SolveA(r io.Reader) any {
	i := intcode.NewInterpreterWithIO(r, strings.NewReader("1\n"), os.Stdout)
	i.ExecAll()
	return nil
}

func SolveB(r io.Reader) any {
	i := intcode.NewInterpreterWithIO(r, strings.NewReader("2\n"), os.Stdout)
	i.ExecAll()
	return nil
}
