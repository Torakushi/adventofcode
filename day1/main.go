package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// SOLUTION: 1665
func main() {
	if err := firstPart(); err != nil {
		log.Fatal(err)
	}

	if err := secondPart(); err != nil {
		log.Fatal(err)
	}
}

// the number of times a measurement increases
func firstPart() error {
	file, err := os.Open("data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	increaseCounter := 0
	actualDepth := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("error while parsing to int: %s", err)
		}

		// Not initialized
		if actualDepth < 0 {
			actualDepth = depth
			continue
		}

		if depth >= actualDepth {
			increaseCounter++
		}

		actualDepth = depth
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	fmt.Printf("%d measurements are larger than previous one for the part part\n", increaseCounter)
	return nil
}

type window struct {
	sum int
	len int
}

// the number of times the sum of measurements in a 3 - sliding window increases
func secondPart() error {
	file, err := os.Open("data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	increaseCounter := 0
	var wA, wB, wC window
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("error while parsing to int: %s", err)
		}

		if wA.len < 3 {
			wA.sum += depth
			wA.len++
		}

		if wB.len < 3 && wA.len > 1 {
			wB.sum += depth
			wB.len++
		}

		if wC.len < 3 && wB.len > 1 {
			wC.sum += depth
			wC.len++
		}

		if wB.len == 3 {
			if wB.sum > wA.sum {
				increaseCounter++
			}
			wA, wB, wC = wB, wC, window{sum: depth, len: 1}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	fmt.Printf("%d measurements are larger than previous one for the second part\n", increaseCounter)
	return nil
}
