package intmath

import "math"

type Point struct {
	X, Y int
}

func Gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return Gcd(b, a%b)
}

func Lcm(a, b int) int {
	return a * b / Gcd(a, b)
}

func CeilDiv(x, div int) int {
	d, m := x/div, x%div
	if m != 0 {
		d++
	}
	return d
}

func Min(xs ...int) (min int) {
	min = math.MaxInt
	for _, x := range xs {
		if x < min {
			min = x
		}
	}
	return
}

func Max(xs ...int) (max int) {
	max = math.MinInt
	for _, x := range xs {
		if x > max {
			max = x
		}
	}
	return
}

func Abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func DistSquared(x1, y1, x2, y2 int) int {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func DistManhattan(x1, y1, x2, y2 int) int {
	return Abs(x2-x1) + Abs(y2-y1)
}
