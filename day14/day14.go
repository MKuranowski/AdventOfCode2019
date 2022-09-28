package day14

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/MKuranowski/AdventOfCode2019/util/input"
	"github.com/MKuranowski/AdventOfCode2019/util/intmath"
	"github.com/MKuranowski/AdventOfCode2019/util/maps2"
)

type RecipePart struct {
	Units      int
	Ingredient string
}

type Recipe struct {
	Input  []RecipePart
	Output RecipePart
}

type RecipeDatabase map[string]Recipe

var (
	RecipePartRegex = regexp.MustCompile(`(\d+) (\w+)`)
)

func ParseRecipeParts(x string) (parts []RecipePart) {
	for _, m := range RecipePartRegex.FindAllStringSubmatch(x, -1) {
		units, err := strconv.Atoi(m[1])
		if err != nil {
			panic(fmt.Errorf("failed to parse units (%q) from part %v: %w", m[1], m, err))
		}
		parts = append(parts, RecipePart{units, m[2]})
	}
	return
}

func ParseRecipe(line string) Recipe {
	s := strings.Split(line, " => ")
	if len(s) != 2 {
		panic(fmt.Errorf("can't parse recipe string: %q", line))
	}

	return Recipe{
		Input:  ParseRecipeParts(s[0]),
		Output: ParseRecipeParts(s[1])[0],
	}
}

func ParseRecipes(r io.Reader) (d RecipeDatabase) {
	d = make(RecipeDatabase)
	lines := input.NewLineIterator(r)
	for lines.Next() {
		l := lines.Get()
		if l == "" {
			continue
		}
		r := ParseRecipe(l)
		d[r.Output.Ingredient] = r
	}
	return
}

func OrePerFuel(db RecipeDatabase, fuel int) (ore int) {
	toGo := map[string]int{"FUEL": fuel}
	leftovers := make(map[string]int)

	for len(toGo) > 0 {
		ingredient, units := maps2.Pop(toGo)

		if ingredient == "ORE" {
			ore += units
		} else {
			// Use up "ingredient" from leftovers
			useFromLeftovers := intmath.Min(units, leftovers[ingredient])
			leftovers[ingredient] -= useFromLeftovers
			units -= useFromLeftovers

			// Apply the recipe for "ingredient"
			if units > 0 {
				r, ok := db[ingredient]
				if !ok {
					panic(fmt.Errorf("no recipe for %s", ingredient))
				}

				applyCount := intmath.CeilDiv(units, r.Output.Units)

				// Save leftovers
				leftovers[ingredient] = r.Output.Units*applyCount - units

				// Enqueue recipe inputs to be produced
				for _, i := range r.Input {
					toGo[i.Ingredient] += i.Units * applyCount
				}
			}
		}
	}

	return
}

func SolveA(r io.Reader) any {
	db := ParseRecipes(r)
	return OrePerFuel(db, 1)
}

func SolveB(r io.Reader) any {
	db := ParseRecipes(r)

	ore := 1000000000000
	fuel := 1

	for {
		// Try to find a better guess
		newFuel := int(float64(ore) / float64(OrePerFuel(db, fuel)) * float64(fuel))
		if newFuel == fuel {
			return fuel
		}
		fuel = newFuel
	}
}
