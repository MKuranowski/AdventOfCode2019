package day10

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/MKuranowski/AdventOfCode2019/util/set"
	"golang.org/x/exp/slices"
)

type Point struct {
	X, Y float64
}

func (from Point) DistanceSquared(to Point) float64 {
	dx, dy := to.X-from.X, to.Y-from.Y
	return dx*dx + dy*dy
}

func NormalizeAngle(ang float64) float64 {
	for ang < 0 {
		ang += 2 * math.Pi
	}
	for ang > 2*math.Pi {
		ang -= 2 * math.Pi
	}
	return ang
}

func AngleHash(from, to Point) int {
	angleDetailed := math.Atan2(to.Y-from.Y, to.X-from.X)
	angleDetailed += math.Pi / 2                  // Shift domain so that upwards (Pi/2) ends up at zero
	angleDetailed = NormalizeAngle(angleDetailed) // Ensure angle is between 0 and 2*pi
	return int(math.Round(angleDetailed * 1000.0))
}

func CountVisiblePoints(root Point, other []Point) int {
	angles := make(set.Set[int])

	for _, p := range other {
		// Don't count yourself
		if p != root {
			angles.Add(AngleHash(root, p))
		}
	}

	return angles.Len()
}

func ReadImage(r io.Reader) (points []Point) {
	br := bufio.NewReader(r)
	x := 0
	y := 0

	for {
		c, err := br.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to read image: %w", err))
		}

		switch c {
		case '\n':
			y++
			x = 0

		case '#':
			points = append(points, Point{float64(x), float64(y)})
			x++

		default:
			x++
		}
	}

	return
}

func SolveA(r io.Reader) any {
	bestRoot := Point{}
	max := math.MinInt
	asteroids := ReadImage(r)

	for _, root := range asteroids {
		visibleFromRoot := CountVisiblePoints(root, asteroids)
		if visibleFromRoot > max {
			max = visibleFromRoot
			bestRoot = root
		}
	}

	fmt.Printf("Monitoring station location: %v\n", bestRoot)
	return max
}

type Bucket struct {
	Angle     int
	Asteroids []Point
}

func BucketCompareAngle(a, b Bucket) int {
	if a.Angle < b.Angle {
		return -1
	} else if a.Angle > b.Angle {
		return 1
	}
	return 0
}

func SolveB(r io.Reader) any {
	// Organize the asteroids in buckets by their angles
	root := Point{23, 29} // NORMAL
	// root := Point{11, 13} // TEST
	buckets := []Bucket(nil)

	for _, asteroid := range ReadImage(r) {
		angle := AngleHash(root, asteroid)
		i, exists := slices.BinarySearchFunc(buckets, Bucket{Angle: angle}, BucketCompareAngle)
		if !exists {
			// New bucket - just insert it
			buckets = append(buckets, Bucket{})
			copy(buckets[i+1:], buckets[i:])
			buckets[i] = Bucket{angle, []Point{asteroid}}
		} else {
			// Existing bucket - add asteroid to the bucket, while keeping it sorted by order
			bucketArr := buckets[i].Asteroids
			j, _ := slices.BinarySearchFunc(
				bucketArr,
				asteroid,
				func(a, b Point) int {
					aDist, bDist := root.DistanceSquared(a), root.DistanceSquared(b)
					if aDist < bDist {
						return -1
					} else if aDist > bDist {
						return 1
					}
					return 0
				},
			)

			// Insert the asteroid to the bucket
			if j == 0 {
				bucketArr = append([]Point{asteroid}, bucketArr...)
			} else {
				bucketArr = append(bucketArr, Point{})
				copy(bucketArr[j+1:], bucketArr[j:])
				bucketArr[j] = asteroid
			}

			buckets[i].Asteroids = bucketArr
		}
	}

	// Rotate until we've made 200 shots
	i := 0
	shot := 1
	for len(buckets) > 0 {
		// Done
		target := buckets[i].Asteroids[0]
		if shot == 200 {
			return int(target.X)*100 + int(target.Y)
		}

		// Take a shot at the first asteroid from the i-th bucket
		if len(buckets[i].Asteroids) == 1 {
			// Last asteroid from bucket needs to be removed - drop the whole bucket
			buckets[i] = Bucket{}
			buckets = append(buckets[:i], buckets[i+1:]...)
			if i == len(buckets) {
				i = 0
			}
		} else {
			buckets[i].Asteroids = buckets[i].Asteroids[1:]
			i = (i + 1) % len(buckets)
		}

		shot++
	}

	panic("shot less than 200 times")
}
