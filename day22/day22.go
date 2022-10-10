package day22

import (
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/util/input"
	"github.com/MKuranowski/AdventOfCode2019/util/intmath"
)

var lastPartRegex = regexp.MustCompile(`\S+$`)

func lastPart(s string) string { return lastPartRegex.FindString(s) }

// Unfortunately, this task is a math task.
// Instead of applying every move independently, we can think about the shuffle
// as a single linear function in the deckSize-integer ring.
// The function (`f(n) = an + b`) returns transforms an index after performing a single shuffle.

// ReadMoves returns the coefficients a and b of the shuffle function
func ReadMoves(r io.Reader, deckSize int) (a, b int) {
	// Start with the identity function (f(n) = n)
	a, b = 1, 0

	it := input.NewLineIterator(r)

	for it.Next() {
		line := it.Get()

		if strings.HasPrefix(line, "cut ") {
			cutOffset, err := strconv.Atoi(lastPart(line))
			if err != nil {
				panic(fmt.Errorf("can't extract offset from %q: %w", line, err))
			}

			// Deal with increment has the following mapping: f(n) = n - o
			// Applying to the existing function: f'(n) = f(n) - o = an + b - o
			// Thus: a' = a (unchanged); b' = b - o

			b = intmath.Mod(b-cutOffset, deckSize)

		} else if strings.HasPrefix(line, "deal with increment ") {
			increment, err := strconv.Atoi(lastPart(line))
			if err != nil {
				panic(fmt.Errorf("can't extract increment from %q: %w", line, err))
			}

			// Deal with increment has the following mapping: f(n) = i * n
			// Applying to the existing function: f'(n) = i * f(n) = i * (an + b) = ain + bi
			// Thus: a' = ai; b' = bi

			a = intmath.ModMul(a, increment, deckSize)
			b = intmath.ModMul(b, increment, deckSize)

		} else if strings.HasPrefix(line, "deal into new stack") {
			// Deal into new stack (reverse) has the following mapping: f(n) = -n-1
			// Applying to existing function: f'(n) = -f(n) -1 = -(an + b) - 1 = -an - b - 1
			// Thus: a' = -a; b' = -b-1

			a = intmath.Mod(-a, deckSize)
			b = intmath.Mod(-b-1, deckSize)

		} else if line != "" {
			panic(fmt.Errorf("invalid line: %q", line))
		}
	}

	return
}

func SolveA(r io.Reader) any {
	deckSize := 10_007
	a, b := ReadMoves(r, deckSize)
	return intmath.Mod(a*2019+b, deckSize)
}

// ReadMovesBig is the same as ReadMoves, expect that it returns *big.Int
func ReadMovesBig(r io.Reader, deckSize *big.Int) (a, b *big.Int) {
	ai, bi := ReadMoves(r, int(deckSize.Int64()))
	return big.NewInt(int64(ai)), big.NewInt(int64(bi))
}

var bigOne = big.NewInt(1)

// ApplyIterations returns a new linear function
// which is the equivalent of performing the shuffle `iterations` times.
//
// This is equivalent of nesting the function (f(f(f(n)))) `iterations` times.
func ApplyIterations(a, b, iterations, deckSize *big.Int) (newA, newB *big.Int) {
	newA, newB = &big.Int{}, &big.Int{}

	// f(f(n)) = f(an + b) = a(an + b) + b = a²n + ba + b = a²n + b(a + 1)
	// f(f(f(n))) = a(a²n + b(a + 1)) + b
	//            = a³n + b(a² + a) + b = a³n + b(a² + a + 1)
	// fⁱ(n) = aⁱn + b(aⁱ⁻¹ + aⁱ⁻² + ... + a² + a¹ + 1)

	newA.Exp(a, iterations, deckSize)

	// From the geometric series formula we can rewrite (aⁱ⁻¹ + aⁱ⁻² + ... + a² + a¹ + 1)
	// as the following: (1 - aⁱ) / (1 - a)
	p := big.NewInt(0).Sub(bigOne, newA)
	p.Mod(p, deckSize)

	q := big.NewInt(0).Sub(bigOne, a)
	q.ModInverse(q, deckSize)

	newB.Set(b).Mul(newB, p).Mod(newB, deckSize).Mul(newB, q).Mod(newB, deckSize)

	return
}

func SolveB(r io.Reader) any {
	position := big.NewInt(2020)
	deckSize := big.NewInt(119315717514047)
	iterations := big.NewInt(101741582076661)

	// Read the single-shuffle function
	singleA, singleB := ReadMovesBig(r, deckSize)

	// Apply all the iterations
	a, b := ApplyIterations(singleA, singleB, iterations, deckSize)

	// Solve for the following equation:
	// 2020 = a*n + b
	// n = (2020-b) / a
	n := &big.Int{}
	n.Sub(position, b)
	n.Mod(n, deckSize)
	n.Mul(n, big.NewInt(0).ModInverse(a, deckSize))
	n.Mod(n, deckSize)

	return n
}
