package day6

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Day6() {
	fmt.Println("DAY6")

	if err := process(80); err != nil {
		log.Fatal(err)
	}

	if err := process(256); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func process(nbDays int) error {
	fishCounter, err := parseData()
	if err != nil {
		return err
	}

	for i := 1; i <= nbDays; i++ {
		newFishes := fishCounter[0]
		fishCounter[0] = 0
		for j := 1; j < len(fishCounter); j++ {
			fishCounter[j-1] = fishCounter[j]
		}
		// Old fishes + New fishes
		fishCounter[6], fishCounter[8] = fishCounter[6]+newFishes, newFishes
	}

	sum := 0
	for _, c := range fishCounter {
		sum += c
	}

	fmt.Printf("%d fishes after %d days !\n", sum, nbDays)
	return nil
}

func parseData() ([]int, error) {
	b, err := os.ReadFile("day6/data.txt")
	if err != nil {
		return nil, err
	}

	args := strings.Split(string(b), ",")
	pop := make([]int, 9)
	for i := 0; i < len(args); i++ {
		n, err := strconv.Atoi(args[i])
		if err != nil {
			return nil, err
		}
		pop[n]++
	}
	return pop, nil
}
