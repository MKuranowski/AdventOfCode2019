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
	"github.com/MKuranowski/AdventOfCode2019/day04"
	"github.com/MKuranowski/AdventOfCode2019/day05"
	"github.com/MKuranowski/AdventOfCode2019/day06"
	"github.com/MKuranowski/AdventOfCode2019/day07"
	"github.com/MKuranowski/AdventOfCode2019/day08"
	"github.com/MKuranowski/AdventOfCode2019/day09"
	"github.com/MKuranowski/AdventOfCode2019/day10"
	"github.com/MKuranowski/AdventOfCode2019/day11"
	"github.com/MKuranowski/AdventOfCode2019/day12"
	"github.com/MKuranowski/AdventOfCode2019/day13"
	"github.com/MKuranowski/AdventOfCode2019/day14"
)

var solutions = map[string]func(io.Reader) any{
	"01a": day01.SolveA,
	"01b": day01.SolveB,
	"02a": day02.SolveA,
	"02b": day02.SolveB,
	"03a": day03.SolveA,
	"03b": day03.SolveB,
	"04a": day04.SolveA,
	"04b": day04.SolveB,
	"05a": day05.SolveA,
	"05b": day05.SolveB,
	"06a": day06.SolveA,
	"06b": day06.SolveB,
	"07a": day07.SolveA,
	"07b": day07.SolveB,
	"08a": day08.SolveA,
	"08b": day08.SolveB,
	"09a": day09.SolveA,
	"09b": day09.SolveB,
	"10a": day10.SolveA,
	"10b": day10.SolveB,
	"11a": day11.SolveA,
	"11b": day11.SolveB,
	"12a": day12.SolveA,
	"12b": day12.SolveB,
	"13a": day13.SolveA,
	"13b": day13.SolveB,
	"14a": day14.SolveA,
	"14b": day14.SolveB,
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
		panic(fmt.Errorf("no solver for %q in main.go lookup table", day))
	}

	// Perform the solution and print the result
	result := solver(f)
	if result != nil {
		fmt.Println(result)
	}
}
