package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type LineIterator struct {
	r    *bufio.Reader
	line string
}

func NewLineIterator(r io.Reader) *LineIterator {
	return &LineIterator{
		r: bufio.NewReader(r),
	}
}

func (i *LineIterator) Get() string { return i.line }

func (i *LineIterator) Next() bool {
	l, err := i.r.ReadString('\n')

	if errors.Is(err, io.EOF) {
		return false
	} else if err != nil {
		panic(fmt.Errorf("failed to move to the next line: %w", err))
	}

	i.line = strings.TrimRight(l, "\n")
	return true
}
