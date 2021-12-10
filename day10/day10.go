package day10

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

func Day10() {
	fmt.Println("DAY10")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

var openingMapWithClosure = map[rune]rune{
	'(': ')',
	'[': ']',
	'<': '>',
	'{': '}',
}

var illegalScoreMap = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var autoScoreMap = map[rune]int{
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,
}

func process() error {
	file, err := os.Open("day10/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	illegaScore := 0
	var autoScore []int
	for scanner.Scan() {
		line := scanner.Text()

		illegalChar, isIllegal, auto := processLine(line)
		if isIllegal {
			illegaScore += illegalScoreMap[illegalChar]
			continue
		}
		autoScore = append(autoScore, auto)
	}

	fmt.Printf("Score for illegal lines: %d\n", illegaScore)

	sort.Ints(autoScore)
	fmt.Printf("Score for auto completion lines: %d\n", autoScore[int(math.Round(float64(len(autoScore)/2)))])
	return nil
}

func processLine(s string) (rune, bool, int) {
	closeCount, openCount := 0, 0
	var open []rune
	for _, r := range s {
		if isClosure(r) {
			// Closing non-existing
			if closeCount == openCount {
				return r, true, 0
			}

			lastOpen := open[len(open)-1]
			// Closing with wrong rune
			if r != openingMapWithClosure[lastOpen] {
				return r, true, 0
			}

			open = open[:len(open)-1]
			closeCount++
			continue
		}
		openCount++
		open = append(open, r)
	}

	autoScore := 0
	for i := len(open) - 1; i >= 0; i-- {
		autoScore = autoScore*5 + autoScoreMap[open[i]]
	}

	return 0, false, autoScore
}

func isClosure(r rune) bool {
	return r == ')' || r == ']' || r == '>' || r == '}'
}
