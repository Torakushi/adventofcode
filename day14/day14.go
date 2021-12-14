package day14

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var polymere string
var pairs = map[string]string{}

func Day14() {
	fmt.Println("DAY14")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func process() error {
	if err := readDatas(); err != nil {
		return nil
	}

	var counter10 = map[string]int{}
	var counter40 = map[string]int{}
	for _, v := range polymere {
		counter10[string(v)]++
		counter40[string(v)]++
	}

	// Get pairs for initial string
	pairs := []string{}
	for i := 0; i < len(polymere)-1; i++ {
		var p string
		if i+2 > len(polymere)-1 {
			p = polymere[i:]
		} else {
			p = polymere[i : i+2]
		}
		pairs = append(pairs, p)
	}

	for _, p := range pairs {
		counter10 = merge(counter10, processStep(p, 10))
	}
	fmt.Printf("First part: Max occurences - Min occurences after 10 steps: %d\n", getScore(counter10))

	for _, p := range pairs {
		counter40 = merge(counter40, processStep(p, 40))
	}
	fmt.Printf("Second part: Max occurences - Min occurences after 40 steps: %d\n", getScore(counter40))
	return nil
}

func readDatas() error {
	file, err := os.Open("day14/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var initPolymere bool
	for scanner.Scan() {
		if !initPolymere {
			polymere = scanner.Text()
			initPolymere = true
			scanner.Scan()
			continue
		}

		arr := strings.Split(scanner.Text(), " -> ")
		pairs[arr[0]] = arr[1]
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}
	return nil
}

// Cache that stores the number of letters for a given pair at a given step
var cache = map[string]map[int]map[string]int{}

func processStep(s string, step int) map[string]int {
	if step == 0 {
		return nil
	}
	if cache[s] != nil {
		if _, ok := cache[s][step]; ok {
			return cache[s][step]

		}
	}

	n := map[string]int{}
	n[pairs[s]]++
	n = merge(n, processStep(string(s[0])+pairs[s], step-1))
	n = merge(n, processStep(pairs[s]+string(s[1]), step-1))

	if cache[s] == nil {
		cache[s] = map[int]map[string]int{}
	}
	if cache[s][step] == nil {
		cache[s][step] = map[string]int{}
	}
	cache[s][step] = n
	return n
}

func getScore(counter map[string]int) int {
	min, max := 0, 0
	for _, v := range counter {
		if min == 0 {
			min = v
		}

		if v > max {
			max = v
		}

		if v < min {
			min = v
		}
	}
	return max - min
}

func merge(m1, m2 map[string]int) map[string]int {
	for k, v := range m2 {
		m1[k] += v
	}
	return m1
}
