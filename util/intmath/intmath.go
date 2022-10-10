package intmath

import (
	"math"
)

func Gcd(a, b int) int {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

func GcdExtended(a, b int) (gcd, x, y int) {
	if a == 0 {
		return b, 0, 1
	}

	var x1, y1 int
	gcd, x1, y1 = GcdExtended(b%a, a)
	return gcd, y1 - (b / a * x1), x1
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

// Returns the Euclidian modulus of x.
//
//	Mod(13, 5) == 3 (== 13 % 5)
//	Mod(-2, 5) == 3 (but -2 % 5 == -2)
func Mod(x, mod int) int {
	return ((x % mod) + mod) % mod
}

// Returns `(a * b) % modulus`
func ModMul(a, b, modulus int) int {
	a = Mod(a, modulus)
	b = Abs(b)

	result := 0

	for b > 0 {
		if b&1 == 1 {
			result = (result + a) % modulus
		}
		a, b = (a<<1)%modulus, b>>1
	}

	return result % modulus
}

// Returns `(base ** exponent) % modulus`
func ModExp(base, exponent, modulus int) int {
	if exponent < 0 {
		panic("ModExp: exponent can't be negative")
	}

	base = Mod(base, modulus)
	r := 1

	for exponent > 0 {
		if exponent&1 == 1 {
			r = (r * base) % modulus
		}
		base, exponent = (base*base)%modulus, exponent>>1
	}

	return r
}

// Returns a number `result` such that `(result * base) % modulus == 1`.
// Panics if such a number doesn't exist
func ModMulInv(base, modulus int) int {
	gcd, x, _ := GcdExtended(base, modulus)
	if gcd != 1 {
		panic("no multiplicative inverse")
	}
	return Mod(x, modulus)
}
