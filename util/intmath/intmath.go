package intmath

import "math"

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
