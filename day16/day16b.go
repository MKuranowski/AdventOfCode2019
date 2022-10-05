package day16

import (
	"io"
	"math/big"
)

func getOffset(msg []int) (offset int) {
	offset += msg[6]
	offset += msg[5] * 10
	offset += msg[4] * 100
	offset += msg[3] * 1000
	offset += msg[2] * 10000
	offset += msg[1] * 100000
	offset += msg[0] * 1000000
	return
}

func SolveB(r io.Reader) any {
	// I have no fucking clue how are you supposed to come up with the solution
	// Idea from https://github.com/encse/adventofcode/blob/master/2019/Day16/Solution.cs
	msg := ReadMessage(r).Data()
	offset := getOffset(msg)
	columnsToEnd := len(msg)*10_000 - offset

	coefficientsMod10 := make([]int64, columnsToEnd+1)

	coefficient := big.NewInt(1)
	bigTen := big.NewInt(10)

	for i := 1; i <= columnsToEnd; i++ {
		coefficientsMod10[i] = big.NewInt(0).Mod(coefficient, bigTen).Int64()
		coefficient = coefficient.Mul(coefficient, big.NewInt(int64(i)+99)).Div(coefficient, big.NewInt(int64(i)))
	}

	result := int64(0)
	for i := 1; i <= 8; i++ {
		sum := int64(0)
		for j := i; j <= columnsToEnd; j++ {
			input := int64(msg[(offset+j-1)%len(msg)])
			sum += input * coefficientsMod10[j-i+1]
		}
		result *= 10
		result += sum % 10
	}

	return result
}
