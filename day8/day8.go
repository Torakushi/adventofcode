package day8

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

func Day8() {
	fmt.Println("DAY8")

	if err := firstPart(); err != nil {
		log.Fatal(err)
	}

	if err := secondPart(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func readDatas() ([][]string, [][]string, error) {
	file, err := os.Open("day8/data.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	reInput := regexp.MustCompile(`^([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+)`)
	reOutput := regexp.MustCompile(`([\w+]+) ([\w+]+) ([\w+]+) ([\w+]+)$`)

	var inputs, outputs [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		inputs = append(inputs, reInput.FindStringSubmatch(text)[1:])
		outputs = append(outputs, reOutput.FindStringSubmatch(text)[1:])
	}
	return inputs, outputs, nil
}

func firstPart() error {
	_, outputs, err := readDatas()
	if err != nil {
		return err
	}

	counter := 0
	for _, output := range outputs {
		for i := 0; i < len(output); i++ {
			if (len(output[i]) >= 2 && len(output[i]) <= 4) || len(output[i]) == 7 {
				counter++
			}
		}
	}
	fmt.Printf("First part (1,7,4,8 count): %d\n", counter)
	return nil
}

func secondPart() error {
	inputs, outputs, err := readDatas()
	if err != nil {
		return err
	}

	count := 0
	for i := 0; i < len(inputs); i++ {
		count += decode(inputs[i], outputs[i])
	}

	fmt.Printf("Second part (sum of outputs): %d\n", count)
	return nil
}

func decode(input, output []string) int {
	mapByLen := make(map[int][]string)
	for _, in := range input {
		mapByLen[len(in)] = append(mapByLen[len(in)], in)
	}

	parsing := map[string]int{
		mapByLen[2][0]: 1,
		mapByLen[3][0]: 7,
		mapByLen[4][0]: 4,
		mapByLen[7][0]: 8,
	}

	one := mapByLen[2][0]
	four := mapByLen[4][0]

	var nine string
	for _, sixLetter := range mapByLen[6] {
		if len(substractString(sixLetter, one)) == 5 {
			parsing[sixLetter] = 6
		} else if len(substractString(sixLetter, four)) == 3 {
			parsing[sixLetter] = 0
		} else {
			nine = sixLetter
			parsing[sixLetter] = 9
		}
	}

	for _, fiveLetter := range mapByLen[5] {
		if len(substractString(fiveLetter, one)) == 3 {
			parsing[fiveLetter] = 3
		} else if len(substractString(fiveLetter, nine)) == 0 {
			parsing[fiveLetter] = 5
		} else {
			parsing[fiveLetter] = 2
		}
	}

	result := 0
	for i := 0; i < len(output); i++ {
		if _, ok := parsing[output[i]]; ok {
			result += parsing[output[i]] * int(math.Pow(10, float64(3-i)))
			continue
		}

	L:
		for _, v := range mapByLen[len(output[i])] {
			for _, c := range v {
				if !strings.ContainsRune(output[i], c) {
					continue L
				}
			}
			result += parsing[v] * int(math.Pow(10, float64(3-i)))
			break
		}
	}

	return result
}

func substractString(s, m string) string {
	for _, c := range m {
		s = strings.Replace(s, string(c), "", 1)
	}
	return s
}
