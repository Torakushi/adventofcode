package day2

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Day2() {
	fmt.Println("DAY2:")

	if err := firstPart(); err != nil {
		log.Fatal(err)
	}

	if err := secondPart(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func firstPart() error {
	file, err := os.Open("day2/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	position, depth := 0, 0
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " ")
		dir := arr[0]
		pow, err := strconv.Atoi(arr[1])
		if err != nil {
			return fmt.Errorf("error while parsing power %q into int: %s", arr[1], err)
		}
		switch dir {
		case "forward":
			position += pow
		case "down":
			depth += pow
		case "up":
			depth -= pow
		default:
			return fmt.Errorf("Unknown dir %q", dir)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	fmt.Printf("First part: Position %d, Depth %d, multiplication %d\n", position, depth, position*depth)
	return nil
}

func secondPart() error {
	file, err := os.Open("day2/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	position, depth, aim := 0, 0, 0
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " ")
		dir := arr[0]
		pow, err := strconv.Atoi(arr[1])
		if err != nil {
			return fmt.Errorf("error while parsing power %q into int: %s", arr[1], err)
		}
		switch dir {
		case "forward":
			position += pow
			depth += aim * pow
		case "down":
			aim += pow
		case "up":
			aim -= pow
		default:
			return fmt.Errorf("Unknown dir %q", dir)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	fmt.Printf("Second part: Position %d, Depth %d, multiplication %d\n", position, depth, position*depth)
	return nil
}
