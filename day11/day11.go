package day11

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Day11() {
	fmt.Println("DAY11")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

var datas [][]*octopus

type octopus struct {
	row, col int
	level    int
}

func (c *octopus) String() string {
	return fmt.Sprintf("%d-%d", c.row, c.col)
}

func process() error {
	file, err := os.Open("day11/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	row := 0
	for scanner.Scan() {
		r := make([]*octopus, 10)
		arrStr := strings.Split(scanner.Text(), "")
		for i, v := range arrStr {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			r[i] = &octopus{row: row, col: i, level: n}
		}
		datas = append(datas, r)
		row++
	}

	flashCounter, flashTotal, step := 0, 0, 0
	for step < 100 || flashCounter != 100 {
		step++
		flashCounter = processStep()
		if step <= 100 {
			flashTotal += flashCounter
		}
	}

	fmt.Printf("Number of octopuses that flashed: %d\n", flashTotal)
	fmt.Printf("The first step that octopuses flash together: %d\n", step)

	return nil
}

func processStep() int {
	flashed := make(map[string]bool)
	flashNumber := 0
	for row := 0; row < len(datas); row++ {
		for col := 0; col < len(datas[row]); col++ {
			flashNumber += levelUp(datas[row][col], flashed)
		}
	}
	return flashNumber
}

func levelUp(octo *octopus, flashed map[string]bool) int {
	queue := []*octopus{octo}
	count := 0
	for len(queue) > 0 {
		o := queue[0]
		queue = queue[1:]
		if flashed[o.String()] {
			continue
		}
		o.level++

		if o.level > 9 {
			flashed[o.String()] = true
			o.level = 0
			count++
			if o.row > 0 {
				queue = append(queue, datas[o.row-1][o.col])
				if o.col > 0 {
					queue = append(queue, datas[o.row-1][o.col-1])
				}
				if o.col < len(datas)-1 {
					queue = append(queue, datas[o.row-1][o.col+1])
				}
			}
			if o.col > 0 {
				queue = append(queue, datas[o.row][o.col-1])
				if o.row < len(datas)-1 {
					queue = append(queue, datas[o.row+1][o.col-1])
				}
			}

			if o.col < len(datas)-1 {
				queue = append(queue, datas[o.row][o.col+1])
				if o.row < len(datas)-1 {
					queue = append(queue, datas[o.row+1][o.col+1])
				}
			}
			if o.row < len(datas)-1 {
				queue = append(queue, datas[o.row+1][o.col])
			}
		}
	}
	return count
}

func displayMap() {
	for _, v := range datas {
		lvlArr := []int{}
		for _, r := range v {
			lvlArr = append(lvlArr, r.level)
		}
		fmt.Println(lvlArr)
	}
}
