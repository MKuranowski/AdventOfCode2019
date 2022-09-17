package day06

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/MKuranowski/AdventOfCode2019/util/input"
	"github.com/MKuranowski/AdventOfCode2019/util/set"
)

type Object struct {
	Name     string
	Parent   *Object
	Children []*Object

	Depth       uint
	Descendants set.Set[string]
}

func (o *Object) CountAllOrbits() int {
	if o.Parent == nil {
		return 0
	}
	return 1 + o.Parent.CountAllOrbits()
}

func (o *Object) HopsTo(name string) uint {
	if o.Name == name {
		return 0
	}

	for _, child := range o.Children {
		if child.Name == name {
			return 1
		} else if child.Descendants.Has(name) {
			return 1 + child.HopsTo(name)
		}
	}

	panic(fmt.Errorf("%s is not %s's descendant", name, o.Name))
}

func (o *Object) calculateDescendants() {
	if o.Descendants != nil {
		return
	}

	o.Descendants = make(set.Set[string])
	for _, child := range o.Children {
		child.calculateDescendants()
		o.Descendants.Add(child.Name)
		o.Descendants.Union(child.Descendants)
	}
}

func (o *Object) calculateDepth() {
	if o.Parent == nil || o.Depth > 0 {
		// Don't re-calculate
		return
	} else {
		o.Parent.calculateDepth()
		o.Depth = o.Parent.Depth + 1
	}
}

func ReadObjects(r io.Reader) (objects map[string]*Object) {
	objects = map[string]*Object{}
	inLines := input.NewLineIterator(r)

	for inLines.Next() {
		objectNames := strings.Split(inLines.Get(), ")")
		parentName, childName := objectNames[0], objectNames[1]

		parent := objects[parentName]
		if parent == nil {
			parent = &Object{Name: parentName}
			objects[parentName] = parent
		}

		child := objects[childName]
		if child == nil {
			child = &Object{Name: childName}
			objects[childName] = child
		}

		child.Parent = parent
		parent.Children = append(parent.Children, child)
	}

	for _, object := range objects {
		object.calculateDescendants()
		object.calculateDepth()
	}

	return
}

func SolveA(r io.Reader) any {
	objects := ReadObjects(r)
	todoQueue := make(chan *Object)
	sum := &atomic.Uint32{}
	wg := &sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for obj := range todoQueue {
				sum.Add(uint32(obj.CountAllOrbits()))
			}
			wg.Done()
		}()
	}

	for _, obj := range objects {
		todoQueue <- obj
	}

	close(todoQueue)
	wg.Wait()
	return sum.Load()
}

func SolveB(r io.Reader) any {
	// I have tried using Dijkstra's algorithm, but it ended up being too slow
	// only after that I have realized that we can interpret the data as a tree,
	// not only a generic graph.
	objects := ReadObjects(r)

	// Find the deepest node with both "SAN" and "YOU" as descendants
	deepestCommonAscendant := objects["COM"]
	for _, object := range objects {
		if object.Descendants.Has("SAN") && object.Descendants.Has("YOU") && object.Depth > deepestCommonAscendant.Depth {
			deepestCommonAscendant = object
		}
	}

	return deepestCommonAscendant.HopsTo("SAN") + deepestCommonAscendant.HopsTo("YOU") - 2
}
