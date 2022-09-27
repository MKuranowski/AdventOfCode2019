package day09

import (
	"io"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

func SolveA(r io.Reader) any {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	i := intcode.NewInterpreterNewIO(r)
	go input.StaticSender(i.Input, wg, 1)
	go input.StdoutReceiver(i.Output, wg)

	i.ExecAll()
	wg.Wait()

	return nil
}

func SolveB(r io.Reader) any {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	i := intcode.NewInterpreterNewIO(r)
	go input.StaticSender(i.Input, wg, 2)
	go input.StdoutReceiver(i.Output, wg)

	i.ExecAll()
	wg.Wait()

	return nil
}
