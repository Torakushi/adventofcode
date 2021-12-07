package day7

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Day7() {
	fmt.Println("DAY7")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func process() error {
	arrMM, err := readData()
	if err != nil {
		return err
	}

	// the result for the second part is always mean (rounded sup OR inf), we have to check
	fuelCounterConstant, fuelCounterRaisingDown, fuelCounterRaisingUp := 0, 0, 0
	for i := 0; i < len(arrMM.arr); i++ {
		fuelCounterConstant += getFuelWastedConstant(arrMM.arr[i], arrMM.median)
		fuelCounterRaisingDown += getFuelWastedRaising(arrMM.arr[i], arrMM.meanDown)
		fuelCounterRaisingUp += getFuelWastedRaising(arrMM.arr[i], arrMM.meanDown+1)
	}

	fuelCounterRaising, secondPos := 0, 0
	if fuelCounterRaisingDown < fuelCounterRaisingUp {
		fuelCounterRaising = fuelCounterRaisingDown
		secondPos = arrMM.meanDown
	} else {
		fuelCounterRaising = fuelCounterRaisingUp
		secondPos = arrMM.meanDown + 1
	}

	fmt.Printf("First part with constant fuel waste: all crabs meet at %d, fuel waste: %d\n", arrMM.median, fuelCounterConstant)
	fmt.Printf("Second part with raising fuel waste: all crabs meet at %d, fuel waste: %d\n", secondPos, fuelCounterRaising)

	return nil
}

// Return data and median
func readData() (arrayWithMeanMedian, error) {
	b, err := os.ReadFile("day7/data.txt")
	if err != nil {
		return arrayWithMeanMedian{}, err
	}

	args := strings.Split(string(b), ",")

	sum := 0
	crabs := make([]int, len(args))
	for i := 0; i < len(args); i++ {
		n, err := strconv.Atoi(args[i])
		if err != nil {
			return arrayWithMeanMedian{}, err
		}
		crabs[i] = n
		sum += n
	}

	arrMM := arrayWithMeanMedian{
		meanDown: int(sum / len(crabs)),
		arr:      crabs,
	}

	// Calculate median
	sort.Ints(crabs)
	med := len(crabs) / 2

	if len(crabs)%2 != 0 {
		arrMM.median = crabs[med]
		return arrMM, nil
	}

	arrMM.median = (crabs[med-1] + crabs[med]) / 2
	return arrMM, nil
}

func getFuelWastedConstant(from, to int) int {
	return int(math.Abs(float64(from - to)))
}

func getFuelWastedRaising(from, to int) int {
	n := math.Abs(float64(from - to))
	return int(n * (n + 1) / 2)
}

type arrayWithMeanMedian struct {
	arr      []int
	meanDown int
	median   int
}
