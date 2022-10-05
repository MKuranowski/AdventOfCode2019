package matrix

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

type Real interface {
	constraints.Integer | constraints.Float
}

type Matrix[T Real] struct {
	height, width int
	data          []T
}

func New[T Real](height, width int) Matrix[T] {
	return Matrix[T]{
		height: height,
		width:  width,
		data:   make([]T, height*width),
	}
}

func (m Matrix[T]) Len() int {
	return m.height * m.width
}

func (m Matrix[T]) Height() int {
	return m.height
}

func (m Matrix[T]) Width() int {
	return m.width
}

func (m Matrix[T]) Data() []T {
	return m.data
}

func (m Matrix[T]) Get(row, col int) T {
	if row >= m.height || col >= m.width {
		panic(fmt.Errorf("matrix access out of bounds: [%d %d] (dimensions: %d %d)", row, col, m.height, m.width))
	}

	return m.data[row*m.width+col]
}

func (m Matrix[T]) Set(row, col int, x T) {
	if row >= m.height || col >= m.width {
		panic(fmt.Errorf("matrix access out of bounds: [%d %d] (dimensions: %d %d)", row, col, m.height, m.width))
	}

	m.data[row*m.width+col] = x
}

func (m Matrix[T]) String() string {
	b := strings.Builder{}
	b.WriteByte('[')
	for row := 0; row < m.height; row++ {
		if row != 0 {
			b.WriteByte('\n')
			b.WriteByte(' ')
		}
		b.WriteByte('[')

		for col := 0; col < m.width; col++ {
			b.WriteByte('\t')
			b.WriteString(fmt.Sprint(m.Get(row, col)))
		}
		b.WriteByte(']')
	}
	b.WriteByte(']')
	b.WriteByte('\n')
	return b.String()
}

func (m Matrix[T]) Copy() Matrix[T] {
	n := New[T](m.height, m.width)
	copy(n.data, m.data)
	return n
}

func (m Matrix[T]) FillScalar(x T) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] = x
	}
}

func (m Matrix[T]) AddScalar(x T) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] += x
	}
}

func (m Matrix[T]) SubScalar(x T) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] -= x
	}
}

func (m Matrix[T]) MulScalar(x T) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] *= x
	}
}

func (m Matrix[T]) Apply(f func(T) T) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] = f(m.data[i])
	}
}

func (m1 Matrix[T]) Add(m2 Matrix[T]) {
	if m1.Len() != m2.Len() {
		panic("elementwise addition on matrices with different sizes")
	}

	for i := 0; i < len(m1.data); i++ {
		m1.data[i] += m2.data[i]
	}
}

func (m1 Matrix[T]) Sub(m2 Matrix[T]) {
	if m1.Len() != m2.Len() {
		panic("elementwise subtraction on matrices with different sizes")
	}

	for i := 0; i < len(m1.data); i++ {
		m1.data[i] -= m2.data[i]
	}
}

func (m1 Matrix[T]) MulElementwise(m2 Matrix[T]) {
	if m1.Len() != m2.Len() {
		panic("elementwise multiplication on matrices with different sizes")
	}

	for i := 0; i < len(m1.data); i++ {
		m1.data[i] *= m2.data[i]
	}
}

func (m1 Matrix[T]) MatMul(m2 Matrix[T]) Matrix[T] {
	dest := New[T](m1.height, m2.width)
	m1.MatMulInto(m2, dest)
	return dest
}

func (m1 Matrix[T]) MatMulInto(m2, dest Matrix[T]) {
	if m1.width != m2.height || dest.height != m1.height || dest.width != m2.width {
		panic(fmt.Errorf("invalid sizes for matrix multiplication: [%d %d] @ [%d %d] = [%d %d]",
			m1.height, m1.width, m2.height, m2.width, dest.height, dest.width))
	}

	commonLen := m1.width

	for row := 0; row < dest.height; row++ {
		for col := 0; col < dest.width; col++ {
			field := row*dest.width + col
			dest.data[field] = 0

			for axis := 0; axis < commonLen; axis++ {
				dest.data[field] += m1.Get(row, axis) * m2.Get(axis, col)
			}
		}
	}
}
