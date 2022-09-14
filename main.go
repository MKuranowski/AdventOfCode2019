package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/day01"
	"github.com/MKuranowski/AdventOfCode2019/day02"
	"github.com/MKuranowski/AdventOfCode2019/day03"
)

var solutions = map[string]func(io.Reader) any{
	"01a": day01.SolveA,
	"01b": day01.SolveB,
	"02a": day02.SolveA,
	"02b": day02.SolveB,
	"03a": day03.SolveA,
	"03b": day03.SolveB,
}

func loadInput(day string, test bool) io.ReadCloser {
	// Try to read a file with "a" or "b" suffix
	var fileName string
	if test {
		fileName = fmt.Sprintf("input/%s-test", day)
	} else {
		fileName = fmt.Sprintf("input/%s", day)
	}
	f, err := os.Open(fileName)
	if err == nil {
		return f
	} else if !errors.Is(err, fs.ErrNotExist) {
		panic(fmt.Errorf("failed to read input from file %s: %w", fileName, err))
	}

	// Second: try to read a file without the "a" or "b" suffix
	day = strings.TrimRight(day, "ab")
	if test {
		fileName = fmt.Sprintf("input/%s-test", day)
	} else {
		fileName = fmt.Sprintf("input/%s", day)
	}

	f, err = os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to read input from file %s: %w", fileName, err))
	}
	return f
}

func main() {
	// Parse arguments
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s DAY-NUMBER [test]\n", os.Args[0])
		os.Exit(1)
	}

	day := os.Args[1]
	// Enable test data?
	test := len(os.Args) == 3 && os.Args[2] == "test"

	// Open the input file
	f := loadInput(day, test)
	defer f.Close()

	// Get the solver function
	solver, ok := solutions[day]
	if !ok {
		panic(fmt.Errorf("no solution for day %q", day))
	}

	// Perform the solution and print the result
	result := solver(f)
	fmt.Println(result)
}
