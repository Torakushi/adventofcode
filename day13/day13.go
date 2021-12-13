package day13

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

type instruction struct {
	d string
	p int
}

var mapPerX = make(map[int]map[int]struct{})
var mapPerY = make(map[int]map[int]struct{})
var instructions = []*instruction{}

func Day13() {
	fmt.Println("DAY13")

	if err := process(); err != nil {
		log.Fatal(err)
	}

}

func process() error {
	readDatas()

	for i, inst := range instructions {
		processInstruction(inst)
		fmt.Printf("After the %d instructions (fold %s:%d), there are %d dots\n", i+1, inst.d, inst.p, getNumberOfPoints())
	}

	printFold()
	return nil
}

func processInstruction(inst *instruction) {
	if inst.d == "y" {
		foldY(inst.p)
		return
	}
	foldX(inst.p)
}

func foldX(p int) {
	newMapY := make(map[int]map[int]struct{})
	newMapX := make(map[int]map[int]struct{})
	for x, m := range mapPerX {
		if x == p {
			continue
		}

		if x > p {
			for y := range m {
				if newMapX[p-(x-p)] == nil {
					newMapX[p-(x-p)] = map[int]struct{}{}
				}
				if newMapY[y] == nil {
					newMapY[y] = map[int]struct{}{}
				}
				newMapX[p-(x-p)][y] = struct{}{}
				newMapY[y][p-(x-p)] = struct{}{}
			}
			continue
		}

		for y := range m {
			if newMapY[y] == nil {
				newMapY[y] = map[int]struct{}{}
			}
			if newMapX[x] == nil {
				newMapX[x] = map[int]struct{}{}
			}
			newMapX[x][y] = struct{}{}
			newMapY[y][x] = struct{}{}
		}
	}
	mapPerY = newMapY
	mapPerX = newMapX
}

func foldY(p int) {
	newMapY := make(map[int]map[int]struct{})
	newMapX := make(map[int]map[int]struct{})
	for y, m := range mapPerY {
		if y == p {
			continue
		}

		if y > p {
			for x := range m {
				if newMapY[p-(y-p)] == nil {
					newMapY[p-(y-p)] = map[int]struct{}{}
				}
				if newMapX[x] == nil {
					newMapX[x] = map[int]struct{}{}
				}
				newMapY[p-(y-p)][x] = struct{}{}
				newMapX[x][p-(y-p)] = struct{}{}
			}
			continue
		}

		for x := range m {
			if newMapY[y] == nil {
				newMapY[y] = map[int]struct{}{}
			}
			if newMapX[x] == nil {
				newMapX[x] = map[int]struct{}{}
			}
			newMapY[y][x] = struct{}{}
			newMapX[x][y] = struct{}{}
		}
	}
	mapPerY = newMapY
	mapPerX = newMapX
}

func printFold() {
	maxY := 0
	for y := range mapPerY {
		if maxY < y {
			maxY = y
		}
	}
	maxX := 0
	for x := range mapPerX {
		if maxX < x {
			maxX = x
		}
	}

	fmt.Println("CODE SECRET: ")
	for i := 0; i <= maxY; i++ {
		xs := make([]string, maxX+1)
		for i := range xs {
			xs[i] = "."
		}
		for x := range mapPerY[i] {
			xs[x] = "#"
		}
		fmt.Println(xs)
	}
}
func readDatas() error {
	file, err := os.Open("day13/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var isFold bool
	for scanner.Scan() {
		if scanner.Text() == "" {
			isFold = true
			continue
		}
		if !isFold {
			arr := strings.Split(scanner.Text(), ",")
			x, _ := strconv.Atoi(arr[0])
			y, _ := strconv.Atoi(arr[1])
			if mapPerX[x] == nil {
				mapPerX[x] = map[int]struct{}{}
			}
			if mapPerY[y] == nil {
				mapPerY[y] = map[int]struct{}{}
			}
			mapPerX[x][y] = struct{}{}
			mapPerY[y][x] = struct{}{}
			continue
		}

		re := regexp.MustCompile(`fold along ([x|y])=(\d+)`)
		ms := re.FindAllStringSubmatch(scanner.Text(), -1)
		d := ms[0][1]
		p, _ := strconv.Atoi(ms[0][2])
		instructions = append(instructions, &instruction{d: d, p: p})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}
	return nil
}

func getNumberOfPoints() int {
	sum := 0
	for _, v := range mapPerX {
		sum += len(v)
	}
	return sum
}
