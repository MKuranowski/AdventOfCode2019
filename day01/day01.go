package day01

import (
	"io"
	"strconv"

	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

func SolveA(r io.Reader) any {
	it := input.NewLineIterator(r)
	sum := int64(0)
	for it.Next() {
		i, _ := strconv.ParseInt(it.Get(), 10, 64)
		sum += i/3 - 2
	}
	return sum
}

func recursiveFuelForModule(m int64) int64 {
	sum := int64(0)
	addedFuel := m
	for addedFuel > 0 {
		addedFuel = addedFuel/3 - 2
		if addedFuel < 0 {
			addedFuel = 0
		}
		sum += addedFuel
	}
	return sum
}

func SolveB(r io.Reader) any {
	it := input.NewLineIterator(r)
	sum := int64(0)
	for it.Next() {
		i, _ := strconv.ParseInt(it.Get(), 10, 64)
		sum += recursiveFuelForModule(i)
	}
	return sum
}
