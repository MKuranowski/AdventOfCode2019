package day12

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/MKuranowski/AdventOfCode2019/util/input"
	"github.com/MKuranowski/AdventOfCode2019/util/intmath"
	"github.com/MKuranowski/AdventOfCode2019/util/set"
	"github.com/MKuranowski/AdventOfCode2019/util/vec3"
)

const (
	Steps = 1000
)

var (
	planetRegex = regexp.MustCompile(`<x=([0-9-]+), y=([0-9-]+), z=([0-9-]+)>`)
)

type Planet struct {
	Pos vec3.Vec3[int]
	Vel vec3.Vec3[int]
}

func (p Planet) String() string {
	return fmt.Sprintf(
		"pos=<x=%3d, y=%3d, z=%3d>, vel=<x=%3d, y=%3d, z=%3d>",
		p.Pos.X,
		p.Pos.Y,
		p.Pos.Z,
		p.Vel.X,
		p.Vel.Y,
		p.Vel.Z,
	)
}

func (p Planet) KineticEnergy() int {
	return p.Vel.Abs().Sum()
}

func (p Planet) PotentialEnergy() int {
	return p.Pos.Abs().Sum()
}

func (p Planet) Energy() int {
	return p.KineticEnergy() * p.PotentialEnergy()
}

type State []Planet

func (sp *State) Step() {
	s := *sp
	ns := s.Clone()

	// Apply gravity
	for i := range ns {
		for j := range s {
			g := s[j].Pos.Sub(s[i].Pos).Sign()
			ns[i].Vel.IAdd(g)
		}
	}

	// Apply velocity
	for i := range ns {
		ns[i].Pos.IAdd(ns[i].Vel)
	}

	*sp = ns
}

func (s State) Energy() (total int) {
	for _, p := range s {
		total += p.Energy()
	}
	return
}

func (s State) Dump(w io.Writer) {
	for _, p := range s {
		fmt.Fprintln(w, p)
	}
}

func (s State) Clone() (n State) {
	n = make(State, len(s))
	copy(n, s)
	return
}

func LoadPlanets(r io.Reader) (planets State) {
	l := input.NewLineIterator(r)
	for l.Next() {
		p := Planet{}
		m := planetRegex.FindStringSubmatch(l.Get())
		p.Pos.X, _ = strconv.Atoi(m[1])
		p.Pos.Y, _ = strconv.Atoi(m[2])
		p.Pos.Z, _ = strconv.Atoi(m[3])
		planets = append(planets, p)
	}
	return
}

func SolveA(r io.Reader) any {
	state := LoadPlanets(r)
	for i := 0; i < Steps; i++ {
		state.Step()
	}

	fmt.Println("\nAfter ", Steps, "steps")
	state.Dump(os.Stdout)

	return state.Energy()
}

func HashState(state State, dim int) [8]int {
	return [8]int{
		state[0].Pos.Dim(dim),
		state[0].Vel.Dim(dim),
		state[1].Pos.Dim(dim),
		state[1].Vel.Dim(dim),
		state[2].Pos.Dim(dim),
		state[2].Vel.Dim(dim),
		state[3].Pos.Dim(dim),
		state[3].Vel.Dim(dim),
	}
}

func FindIterationsUntilDimensionLoops(state State, dim int, result *int, wg *sync.WaitGroup) {
	defer wg.Done()
	states := set.Set[[8]int]{}
	stateHash := HashState(state, dim)
	i := 0

	for !states.Has(stateHash) {
		states.Add(stateHash)
		state.Step()
		stateHash = HashState(state, dim)
		i++
	}

	*result = i
}

func SolveB(r io.Reader) any {
	initialState := LoadPlanets(r)
	dimIterations := [3]int{}

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go FindIterationsUntilDimensionLoops(initialState.Clone(), 0, &dimIterations[0], wg)
	go FindIterationsUntilDimensionLoops(initialState.Clone(), 1, &dimIterations[1], wg)
	go FindIterationsUntilDimensionLoops(initialState.Clone(), 2, &dimIterations[2], wg)
	wg.Wait()

	return intmath.Lcm(dimIterations[0], intmath.Lcm(dimIterations[1], dimIterations[2]))
}
