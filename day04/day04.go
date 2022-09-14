package day04

import (
	"io"
	"strconv"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/util/slices"
)

func ReadRange(r io.Reader) (int, int) {
	data, _ := io.ReadAll(r)
	parts := strings.Split(strings.TrimSpace(string(data)), "-")
	low, _ := strconv.Atoi(parts[0])
	high, _ := strconv.Atoi(parts[1])
	return low, high
}

func digitsOf(x int) (digits []uint8) {
	if x < 0 {
		panic("digitsOf only works with non-negative numbers")
	} else if x == 0 {
		return []uint8{0}
	}

	digit := 0
	for x > 0 {
		x, digit = x/10, x%10
		digits = append(digits, uint8(digit))
	}
	slices.Reverse(digits)
	return
}

func ruleTwoAdjacentSame(digits []uint8) bool {
	prev := digits[0]
	for _, digit := range digits[1:] {
		if digit == prev {
			return true
		}
		prev = digit
	}
	return false
}

func ruleNonDecreasing(digits []uint8) bool {
	prev := digits[0]
	for _, digit := range digits[1:] {
		if prev > digit {
			return false
		}
		prev = digit
	}
	return true
}

func ruleHasRunOfTwo(digits []uint8) bool {
	prev := digits[0]
	runLen := 1
	for _, digit := range digits[1:] {
		if digit == prev {
			runLen++
		} else if runLen == 2 {
			return true
		} else {
			runLen = 1
		}
		prev = digit
	}

	return runLen == 2
}

func SolveA(r io.Reader) any {
	start, end := ReadRange(r)
	count := 0
	for i := start; i <= end; i++ {
		iDigits := digitsOf(i)
		if ruleNonDecreasing(iDigits) && ruleTwoAdjacentSame(iDigits) {
			count++
		}
	}
	return count
}

func SolveB(r io.Reader) any {
	start, end := ReadRange(r)
	count := 0
	for i := start; i <= end; i++ {
		iDigits := digitsOf(i)
		if ruleNonDecreasing(iDigits) && ruleHasRunOfTwo(iDigits) {
			count++
		}
	}
	return count
}
