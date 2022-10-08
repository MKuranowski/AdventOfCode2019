package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
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

func ReadLines(r io.Reader) (lines []string) {
	l := NewLineIterator(r)
	for l.Next() {
		lines = append(lines, l.Get())
	}
	return
}

func StaticSender(ch chan<- int, wg *sync.WaitGroup, numbers ...int) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	defer close(ch)

	for _, num := range numbers {
		ch <- num
	}
}

func StdoutReceiver(ch <-chan int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	for num := range ch {
		fmt.Println(num)
	}
}

func AsciiStdoutReceiver(ch <-chan int, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	for c := range ch {
		fmt.Printf("%c", c)
	}
}

func AsciiStaticSender(ch chan<- int, wg *sync.WaitGroup, txt string) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	defer close(ch)

	for _, c := range txt {
		ch <- int(c)
	}
}
