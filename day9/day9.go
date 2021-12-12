package day9

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"adventofcode/utils"
)

func Day9() {
	fmt.Println("DAY9")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

type coord struct {
	row, column int
}

func (c coord) String() string {
	return fmt.Sprintf("r%d-c%d", c.row, c.column)
}

func process() error {
	file, err := os.Open("day9/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var previous, middle, after []int8

	// Initialize the middle line
	scanner.Scan()
	middle, err = utils.StringToInt8Array(scanner.Text())
	if err != nil {
		return err
	}

	floor := make([][]int8, 100)
	floor[0] = middle

	var coords []coord
	sum, rowIndex := 0, 0
	for scanner.Scan() {
		after, err = utils.StringToInt8Array(scanner.Text())
		if err != nil {
			return err
		}

		floor[rowIndex] = middle
		c, s := processWindow(rowIndex, previous, middle, after)
		sum += s
		coords = append(coords, c...)
		rowIndex++

		previous, middle, after = middle, after, nil
	}

	// process last line
	floor[rowIndex] = middle
	c, s := processWindow(rowIndex, previous, middle, after)
	coords = append(coords, c...)
	sum += s

	fmt.Printf("First part, the sum of lower elements is %d\n", sum)

	// GET BASSIN FOR ALL LOWER POINT
	bassinLength := make([]int, len(coords))
	for i, c := range coords {
		bassinLength[i] = getLenBassin(c, floor)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(bassinLength)))

	fmt.Printf("The multiplication of the three largest bassin: %d", bassinLength[0]*bassinLength[1]*bassinLength[2])
	return nil
}

func getLenBassin(c coord, floor [][]int8) int {
	queue := []*coord{&c}
	contains := map[string]bool{}
	lenght := 0
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]

		if floor[c.row][c.column] == 9 {
			continue
		}

		if contains[c.String()] {
			continue
		}

		lenght++
		contains[c.String()] = true

		if c.column > 0 {
			queue = append(queue, &coord{c.row, c.column - 1})
		}

		if c.column < len(floor)-1 {
			queue = append(queue, &coord{c.row, c.column + 1})
		}

		if c.row > 0 {
			queue = append(queue, &coord{c.row - 1, c.column})
		}

		if c.row < len(floor)-1 {
			queue = append(queue, &coord{c.row + 1, c.column})
		}
	}
	return lenght
}

func processWindow(rowIndex int, previous, middle, after []int8) ([]coord, int) {
	sum := 0
	var indexes []coord
	for i := 0; i < len(middle); i++ {
		if i >= 1 && middle[i] >= middle[i-1] {
			continue
		}

		if i < len(middle)-1 && middle[i] >= middle[i+1] {
			continue
		}

		if len(previous) != 0 && middle[i] >= previous[i] {
			i++
			continue
		}

		if len(after) != 0 && middle[i] >= after[i] {
			i++
			continue
		}

		sum += int(1 + middle[i])
		indexes = append(indexes, coord{row: rowIndex, column: i})
		i++
	}

	return indexes, sum
}
