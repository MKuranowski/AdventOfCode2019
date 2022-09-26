package vec3

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Vec3[T constraints.Float | constraints.Integer] struct {
	X, Y, Z T
}

func (v Vec3[T]) String() string {
	return fmt.Sprintf("[%v %v %v]", v.X, v.Y, v.Z)
}

func (v1 Vec3[T]) Add(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

func (v1 *Vec3[T]) IAdd(v2 Vec3[T]) {
	v1.X += v2.X
	v1.Y += v2.Y
	v1.Z += v2.Z
}

func (v1 Vec3[T]) Sub(v2 Vec3[T]) Vec3[T] {
	return Vec3[T]{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

func (v1 *Vec3[T]) ISub(v2 Vec3[T]) {
	v1.X -= v2.X
	v1.Y -= v2.Y
	v1.Z -= v2.Z
}

func (v Vec3[T]) Sign() (r Vec3[int]) {
	if v.X < 0 {
		r.X = -1
	} else if v.X > 0 {
		r.X = 1
	}

	if v.Y < 0 {
		r.Y = -1
	} else if v.Y > 0 {
		r.Y = 1
	}

	if v.Z < 0 {
		r.Z = -1
	} else if v.Z > 0 {
		r.Z = 1
	}

	return
}

func (v Vec3[T]) Abs() (r Vec3[T]) {
	if v.X < 0 {
		r.X = -v.X
	} else {
		r.X = v.X
	}

	if v.Y < 0 {
		r.Y = -v.Y
	} else {
		r.Y = v.Y
	}

	if v.Z < 0 {
		r.Z = -v.Z
	} else {
		r.Z = v.Z
	}

	return
}

func (v Vec3[T]) Sum() T {
	return v.X + v.Y + v.Z
}

func (v Vec3[T]) Dim(idx int) T {
	switch idx {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	default:
		panic(fmt.Errorf("invalid dimension for Vec3: %d", idx))
	}
}
