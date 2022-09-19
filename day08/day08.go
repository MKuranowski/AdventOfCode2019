package day08

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"unicode"
)

const (
	Width  = 25
	Height = 6
)

type Layer = [Height][Width]uint8
type Image = []Layer

func LoadImage(r io.Reader) (i Image) {
	br := bufio.NewReader(r)

	l := Layer{}
	row := 0
	column := 0

	for {
		// Read the byte from input
		c, err := br.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to read image: %w", err))
		} else if unicode.IsSpace(rune(c)) || c == '\n' {
			continue
		}

		// Convert the read byte to a pixel value
		pixel := c - '0'
		if pixel > 9 {
			panic("Image contain something other than a digit")
		}

		// Save the pixel
		l[row][column] = pixel
		column++
		if column >= Width {
			row++
			column = 0
		}
		if row >= Height {
			row = 0
			i = append(i, l)
			l = Layer{}
		}
	}

	return
}

func CountOn(l Layer, what uint8) (count int) {
	for row := 0; row < Height; row++ {
		for col := 0; col < Width; col++ {
			if l[row][col] == what {
				count++
			}
		}
	}
	return
}

func SolveA(r io.Reader) any {
	img := LoadImage(r)

	minZeroCount := math.MaxInt
	result := 0

	for _, layer := range img {
		zeroCount := CountOn(layer, 0)
		if zeroCount < minZeroCount {
			minZeroCount = zeroCount
			result = CountOn(layer, 1) * CountOn(layer, 2)
		}
	}

	return result
}

func CollapseImage(img Image) (collapsed Layer) {
	// Initialize the background to all transparent
	for row := 0; row < Height; row++ {
		for col := 0; col < Width; col++ {
			collapsed[row][col] = 2
		}
	}

	// Collapse the layers, starting from the last one
	for i := len(img) - 1; i >= 0; i-- {
		for row := 0; row < Height; row++ {
			for col := 0; col < Width; col++ {
				if img[i][row][col] != 2 {
					collapsed[row][col] = img[i][row][col]
				}
			}
		}
	}

	return
}

func PrintLayer(l Layer, sink io.Writer) {
	for row := 0; row < Height; row++ {
		for col := 0; col < Width; col++ {
			switch l[row][col] {
			case 0:
				sink.Write([]byte{' '})
			case 1:
				sink.Write([]byte{'#'})
			case 2:
				sink.Write([]byte{'x'})
			default:
				panic("unknown color")
			}
		}
		sink.Write([]byte{'\n'})
	}
}

func SolveB(r io.Reader) any {
	img := LoadImage(r)
	PrintLayer(CollapseImage(img), os.Stdout)
	return nil
}
