package day17

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
)

const input = "target area: x=185..221, y=-122..-74"

func Day17() {
	fmt.Println("DAY17")
	if err := process(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

var xMin, xMax, yMin, yMax int

func readInput() {
	re := regexp.MustCompile(`target area: x=(-?\d+)..(-?\d+), y=(-?\d+)..(-?\d+)`)
	s := re.FindAllStringSubmatch(input, -1)[0]
	xMin, _ = strconv.Atoi(s[1])
	xMax, _ = strconv.Atoi(s[2])
	yMin, _ = strconv.Atoi(s[3])
	yMax, _ = strconv.Atoi(s[4])
}

func getConsecutiveSum(n int) int {
	return int(n * (n + 1) / 2)
}

func possibleX(x int) bool {
	if getConsecutiveSum(x) < xMin {
		return false
	}

	xPos := 0
	for xPos < xMax {
		xPos += x
		if xPos >= xMin && xPos <= xMax {
			return true
		}
		x = int(math.Max(float64(x)-1, 0))
	}
	return false
}

func getXInterval() []int {
	xInf := int((-1+math.Sqrt(float64(1+8*xMin)))/2) - 1
	arr := []int{}
	for i := xInf; i <= xMax; i++ {
		if possibleX(i) {
			arr = append(arr, i)
		}
	}
	return arr
}

var print bool

func lauchProbe(xVel, yVel int) (bool, int) {
	xPos, yPos := 0, 0
	yHigh := 0
	for xPos <= xMax && yPos >= yMin {
		xPos += xVel
		yPos += yVel
		if yPos >= yHigh {
			yHigh = yPos
		}

		if xPos >= xMin && xPos <= xMax && yPos >= yMin && yPos <= yMax {
			return true, yHigh
		}
		if xPos > xMax || yPos < yMin {
			return false, 0
		}
		xVel = int(math.Max(float64(xVel)-1, 0))
		yVel--
	}
	return false, 0
}

func process() error {
	count := 0
	high := 0
	readInput()
	xs := getXInterval()
	fmt.Println(xs)
	for _, xVel := range xs {
		for yVel := -1 * getConsecutiveSum(yMax); yVel < -1*yMin; yVel++ {

			if ok, h := lauchProbe(xVel, yVel); ok {
				count++
				if high < h {
					high = h
				}
			}
		}
	}

	fmt.Println(high)
	fmt.Println(count)
	return nil
}
