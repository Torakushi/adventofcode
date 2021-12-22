package day22

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	cubes orderedCube
	re    = regexp.MustCompile(`(\w+) x=(-?\d+)..(-?\d+),y=(-?\d+)..(-?\d+),z=(-?\d+)..(-?\d+)`)
)

type cube struct {
	xmin, xmax int
	ymin, ymax int
	zmin, zmax int
	onCounter  int
	index      int
}

// Important here !! We use a heap (~ priority queue) in this case.
// All our cubes are sorted by xmin, so that, when checking if some cubes collapses, we can stop quickly
// until x coordonates are not in the same scale.
// Len(), Less(), Swap(), Push(), and Pop() are here to implement the 'heap' interface.
type orderedCube []*cube

func (oc orderedCube) Len() int { return len(oc) }

func (oc orderedCube) Less(i, j int) bool {
	if oc[i].xmin < oc[j].xmin {
		return true
	}
	if oc[i].xmin == oc[j].xmin {
		return oc[i].xmax > oc[i].xmin
	}
	return false
}

func (oc orderedCube) Swap(i, j int) {
	oc[i], oc[j] = oc[j], oc[i]
	oc[i].index = i
	oc[j].index = j
}

func (oc *orderedCube) Push(x interface{}) {
	n := len(*oc)
	item := x.(*cube)
	item.index = n
	*oc = append(*oc, item)
}

func (oc *orderedCube) Pop() interface{} {
	old := *oc
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*oc = old[0 : n-1]
	return item
}

func Day22() {
	fmt.Println("DAY22")
	t := time.Now()
	process(false)
	d := time.Since(t)
	fmt.Printf("First Part: Initialization took %d and %d cubes are on !! \n\n", d, getOnCubes())

	t = time.Now()
	process(false)
	d = time.Since(t)
	fmt.Printf("Second Part: Reboot took %d and %d cubes are on !! \n\n", d, getOnCubes())

}

func NewCube(xmin, xmax, ymin, ymax, zmin, zmax int) *cube {
	return &cube{
		xmin:      xmin,
		xmax:      xmax,
		ymin:      ymin,
		ymax:      ymax,
		zmin:      zmin,
		zmax:      zmax,
		onCounter: int((math.Abs(float64(zmax-zmin)) + 1) * (math.Abs(float64(xmax-xmin)) + 1) * (math.Abs(float64(ymin-ymax)) + 1)),
	}
}

func process(partOne bool) error {
	file, err := os.Open("day22/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	cubes = orderedCube{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		applyInstruction(scanner.Text(), partOne)
	}
	return nil
}

func applyInstruction(s string, isPartOne bool) {
	m := re.FindAllStringSubmatch(s, -1)[0]
	var value int8 = 0
	if m[1] == "on" {
		value = 1
	}
	xmin, _ := strconv.Atoi(m[2])
	xmax, _ := strconv.Atoi(m[3])
	ymin, _ := strconv.Atoi(m[4])
	ymax, _ := strconv.Atoi(m[5])
	zmin, _ := strconv.Atoi(m[6])
	zmax, _ := strconv.Atoi(m[7])

	// Take only -50..50
	if isPartOne && !isForPartOne(xmin, xmax, ymin, ymax, zmin, zmax) {
		return
	}

	// Init heap if it is not done
	if len(cubes) == 0 && value == 1 {
		heap.Push(&cubes, NewCube(xmin, xmax, ymin, ymax, zmin, zmax))
		return
	}

	// Check for all existing cube, if it has superposition.
	// As we use a ordered list (sort with xmin), we check until c.xmin is greater than the given xmax
	var cs []*cube
	for len(cubes) > 0 {
		c := heap.Pop(&cubes).(*cube)
		if c.xmin > xmax {
			cs = append(cs, c)
			break
		}

		// Get overtaking part
		if c.colapseWith(xmin, xmax, ymin, ymax, zmin, zmax) {
			cs = append(cs, getCollapsedCubes(c, xmin, xmax, ymin, ymax, zmin, zmax, value)...)
		} else {
			cs = append(cs, c)
		}
	}

	// As we get all overtaking parts in "getCollapsedCubes", add the new cube if the value is 1 (on)
	if value == 1 {
		cs = append(cs, NewCube(xmin, xmax, ymin, ymax, zmin, zmax))
	}

	for _, c := range cs {
		heap.Push(&cubes, c)
	}
}

func isForPartOne(xmin, xmax, ymin, ymax, zmin, zmax int) bool {
	return xmin >= -50 && xmax <= 50 && ymin >= -50 && ymax <= 50 && zmin >= -50 && zmax <= 50
}

func (c *cube) containsCube(xmin, xmax, ymin, ymax, zmin, zmax int) bool {
	return (c.xmin <= xmin && c.xmax >= xmax) &&
		(c.ymin <= ymin && c.ymax >= ymax) &&
		(c.zmin <= zmin && c.zmax >= zmax)
}

func (c *cube) colapseWith(xmin, xmax, ymin, ymax, zmin, zmax int) bool {
	return ((xmin >= c.xmin && xmin <= c.xmax) || (xmax <= c.xmax && xmax >= c.xmin) || (xmin <= c.xmin && xmax >= c.xmax)) &&
		((ymin >= c.ymin && ymin <= c.ymax) || (ymax <= c.ymax && ymax >= c.ymin) || (ymin <= c.ymin && ymax >= c.ymax)) &&
		((zmin >= c.zmin && zmin <= c.zmax) || (zmax <= c.zmax && zmax >= c.zmin) || (zmin <= c.zmin && zmax >= c.zmax))
}

func getCollapsedCubes(c *cube, xmin, xmax, ymin, ymax, zmin, zmax int, value int8) []*cube {
	// c is contained
	if (xmin <= c.xmin && xmax >= c.xmax) && (ymin <= c.ymin && ymax >= c.ymax) && (zmin <= c.zmin && zmax >= c.zmax) {
		return nil
	}

	// Get the conjugate of the intersection of existing cube 'c' and the new cube  (xmin, xmax, ymin, ymax, zmin, zmax).
	// (It is, all overtaking part of this new cube).
	// We check on each side of c, if there is any overtaking.
	cs := []*cube{}
	if c.xmin < xmin {
		cs = append(cs, NewCube(c.xmin, xmin-1, c.ymin, c.ymax, c.zmin, c.zmax))
	}
	if c.xmax > xmax {
		cs = append(cs, NewCube(xmax+1, c.xmax, c.ymin, c.ymax, c.zmin, c.zmax))
	}
	if c.ymin < ymin {
		cs = append(cs, NewCube(
			int(math.Max(float64(c.xmin), float64(xmin))),
			int(math.Min(float64(c.xmax), float64(xmax))),
			c.ymin,
			ymin-1,
			c.zmin,
			c.zmax,
		))
	}
	if c.ymax > ymax {
		cs = append(cs, NewCube(
			int(math.Max(float64(c.xmin), float64(xmin))),
			int(math.Min(float64(c.xmax), float64(xmax))),
			ymax+1,
			c.ymax,
			c.zmin,
			c.zmax,
		))
	}
	if c.zmax > zmax {
		cs = append(cs, NewCube(
			int(math.Max(float64(c.xmin), float64(xmin))),
			int(math.Min(float64(c.xmax), float64(xmax))),
			int(math.Max(float64(c.ymin), float64(ymin))),
			int(math.Min(float64(c.ymax), float64(ymax))),
			zmax+1,
			c.zmax,
		))
	}
	if c.zmin < zmin {
		cs = append(cs, NewCube(
			int(math.Max(float64(c.xmin), float64(xmin))),
			int(math.Min(float64(c.xmax), float64(xmax))),
			int(math.Max(float64(c.ymin), float64(ymin))),
			int(math.Min(float64(c.ymax), float64(ymax))),
			c.zmin,
			zmin-1,
		))
	}
	return cs
}

func getOnCubes() int {
	c := 0
	for _, v := range cubes {
		c += v.onCounter
	}
	return c
}
