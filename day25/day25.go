package day25

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func Day25() {
	fmt.Println("DAY25")
	readDatas()
	step := 0
	run := true
	for run {
		run = processEastBunchInParallel()
		run = processSouthBunchInParallel() || run
		step++
	}
	fmt.Printf("FINAL DAY! cucumbers stop moving after step %d\n", step)
}

// Two differents map to avoid race condition
var (
	cucumberEast   map[int]map[int]bool // For parallel
	cucumbersouth  map[int]map[int]bool // For parallel
	maxCol, maxRow int
)

func readDatas() error {
	file, err := os.Open("day25/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cucumbersouth = map[int]map[int]bool{}
	cucumberEast = map[int]map[int]bool{}
	row := 0
	for scanner.Scan() {
		txt := scanner.Text()
		arr := strings.Split(txt, "")
		maxCol = len(arr)
		for col, v := range arr {
			if v == "." {
				continue
			}

			if v == ">" {
				if cucumberEast[row] == nil {
					cucumberEast[row] = map[int]bool{}
				}
				cucumberEast[row][col] = true
				continue
			}

			if cucumbersouth[col] == nil {
				cucumbersouth[col] = map[int]bool{}
			}
			cucumbersouth[col][row] = true
		}
		row++
	}
	maxRow = row
	return err
}

// No concurency probleme as only one map of the map of map is manipulated
// by goroutines.
// Same go for EastBunch
func processSouthBunchInParallel() bool {
	var hasMove bool
	var wg sync.WaitGroup
	for i := 0; i < maxCol; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			doCycle := !cucumbersouth[i][0] && !cucumberEast[0][i] && cucumbersouth[i][maxRow-1]
			for j := 1; j < maxRow; j++ {
				if !cucumbersouth[i][j] && !cucumberEast[j][i] && cucumbersouth[i][j-1] {
					hasMove = true
					cucumbersouth[i][j] = true
					cucumbersouth[i][j-1] = false
					j++

				}
			}
			if doCycle {
				hasMove = true
				cucumbersouth[i][0] = true
				cucumbersouth[i][maxRow-1] = false
			}
		}(i)
	}
	wg.Wait()
	return hasMove
}

func processEastBunchInParallel() bool {
	var hasMove bool
	var wg sync.WaitGroup
	for i := 0; i < maxRow; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			doCycle := !cucumberEast[i][0] && !cucumbersouth[0][i] && cucumberEast[i][maxCol-1]
			for j := 1; j < maxCol; j++ {
				if !cucumberEast[i][j] && !cucumbersouth[j][i] && cucumberEast[i][j-1] {
					hasMove = true
					cucumberEast[i][j] = true
					cucumberEast[i][j-1] = false
					j++
				}
			}
			if doCycle {
				hasMove = true
				cucumberEast[i][0] = true
				cucumberEast[i][maxCol-1] = false
			}
		}(i)
	}
	wg.Wait()
	return hasMove
}
