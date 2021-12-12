package day12

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Day12() {
	fmt.Println("DAY12")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

var mapPath = make(map[string][]string)

func process() error {
	file, err := os.Open("day12/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		points := strings.Split(scanner.Text(), "-")
		if points[0] == "end" {
			mapPath[points[1]] = append(mapPath[points[1]], points[0])
			continue
		}
		if points[1] == "end" {
			mapPath[points[0]] = append(mapPath[points[0]], points[1])
			continue
		}
		if points[0] == "start" {
			mapPath[points[0]] = append(mapPath[points[0]], points[1])
			continue
		}
		if points[1] == "start" {
			mapPath[points[1]] = append(mapPath[points[1]], points[0])
			continue
		}
		mapPath[points[1]] = append(mapPath[points[1]], points[0])
		mapPath[points[0]] = append(mapPath[points[0]], points[1])
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	fmt.Printf("First part: there are %d paths\n", len(findAllPathes("start", map[string]int{}, true)))
	fmt.Printf("Second part: there are %d paths\n", len(findAllPathes("start", map[string]int{}, false)))

	return nil
}

func findAllPathes(start string, alreadyDone map[string]int, hasVisitSmallCaveTwice bool) [][]string {
	possiblePaths := mapPath[start]
	var result [][]string
	for _, p := range possiblePaths {
		if p == "end" {
			result = append(result, []string{"end"})
			continue
		}

		hasVisited := hasVisitSmallCaveTwice
		if strings.ToUpper(p) != p && alreadyDone[p] >= 1 {
			if !hasVisited {
				hasVisited = true
			} else {
				continue
			}
		}

		newMap := make(map[string]int)
		for k, v := range alreadyDone {
			newMap[k] = v
		}
		newMap[p]++

		r := findAllPathes(p, newMap, hasVisited)
		for _, v := range r {
			n := []string{p}
			result = append(result, append(n, v...))
		}
	}
	return result
}
