package day05

import (
	"io"
	"os"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
)

func SolveA(r io.Reader) any {
	i := intcode.NewInterpreterWithIO(r, strings.NewReader("1"), os.Stdout)
	i.ExecAll()
	return nil
}
