package day25

import (
	"io"

	"github.com/MKuranowski/AdventOfCode2019/day17"
	"github.com/MKuranowski/AdventOfCode2019/intcode"
	"github.com/MKuranowski/AdventOfCode2019/util/input"
)

// Mine map:
//
//  SCIENCE LAB--------CORRIDOR-------STABLES---------HOT CHOC.-------CREW QUART
//                      (mutex)     (astrolabe)     (deh. water)       (wreath)
//                         |                             |
//                         |                             |
//                         |                             |
//                         |          SECURITY CH.       |
//                         |               |             |
//                         |               |             |
//                         |               |             |
//  OBSERVATORY-------GIFT WRAPPING     PASSAGES-----HOLO DECK
//                         |                       (escape pod?)
//                         |
//                         |
//                         |         SICK BAY----NAVIGATION
//                         |         (ice cream)   |
//                         |                       |
// STORAGE             HULL BREACH---------------ARCADE
//    |                    |                     (coin)
//    |                    |                       |
//    |                    |                       |
//  KITCHEN------------ENGINEERING             WARP DRIVE M.
// (asterisk)          (monolith)
//
// SECURITY CHECKPOINT:
// - wreath
// - asterisk
// - monolith
// - astrolabe

func SolveA(r io.Reader) any {
	i := intcode.NewInterpreterNewIO(r)
	s := day17.Screen{}

	go input.AsciiStdinSender(i.Input, nil)
	go s.Run(i.Output, nil)
	i.ExecAll()
	return s.LastNonASCII
}
