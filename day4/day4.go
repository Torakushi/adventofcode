package day4

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/mat"
)

func Day4() {
	fmt.Println("DAY4:")

	if err := firstPart(); err != nil {
		log.Fatal(err)
	}
	if err := secondPart(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func firstPart() error {
	df, input, err := readDatas()
	if err != nil {
		return err
	}

	nums := make(map[int]bool)

	// Check if any columns/rows is BINGO !
	bingo := func(s series.Series) series.Series {
		ints, _ := s.Int()
		for _, i := range ints {
			if !nums[i] {
				return series.Bools(false)
			}
		}
		return series.Bools(true)
	}

	// Sum unmarked number of a column/row
	sum := func(s series.Series) series.Series {
		ints, _ := s.Int()
		sum := 0
		for _, i := range ints {
			if !nums[i] {
				sum += i
			}
		}
		return series.Ints(sum)
	}

	var numGrids []int
	count := 0
	for i := 0; i < df.Nrow(); i++ {
		if i != 0 && i%10 == 0 {
			count++
		}
		numGrids = append(numGrids, count)
	}

	var dfWinner dataframe.DataFrame
	var lastCall int
	for i, n := range input {
		num, _ := strconv.Atoi(n)
		nums[num] = true
		lastCall = num
		if i < 4 {
			continue
		}

		dfWinner = df.Rapply(bingo)
		dfWinner = dfWinner.Mutate(series.Ints(numGrids))

		dfWinner = dfWinner.Filter(
			dataframe.F{
				Colidx:     0,
				Comparator: series.Eq,
				Comparando: true,
			},
		)
		if dfWinner.Nrow() != 0 {
			break
		}
	}

	winnerGridNumber, _ := dfWinner.Elem(0, dfWinner.Ncol()-1).Int()
	winnerGrid := df.Subset(
		[]int{
			winnerGridNumber * 10,
			winnerGridNumber*10 + 1,
			winnerGridNumber*10 + 2,
			winnerGridNumber*10 + 3,
			winnerGridNumber*10 + 4,
		})
	unmarkedNumberSum := winnerGrid.Capply(sum)
	unmarkedNumberSum = unmarkedNumberSum.Rapply(sum)

	uSum, _ := unmarkedNumberSum.Elem(0, 0).Int()
	fmt.Printf("First Part: The sum of unmarked number with last num call is %d\n", uSum*lastCall)

	return nil
}

func secondPart() error {
	df, input, err := readDatas()
	if err != nil {
		return err
	}

	nums := make(map[int]bool)
	// Check if any columns/rows is BINGO !
	bingo := func(s series.Series) series.Series {
		ints, _ := s.Int()
		for _, i := range ints {
			if !nums[i] {
				return series.Bools(false)
			}
		}
		return series.Bools(true)
	}

	// Sum unmarked number of a column/row
	sum := func(s series.Series) series.Series {
		ints, _ := s.Int()
		sum := 0
		for _, i := range ints {
			if !nums[i] {
				sum += i
			}
		}
		return series.Ints(sum)
	}

	alreadyWin := make(map[int]bool)
	// Check all new winners
	newWinner := func(s series.Series) series.Series {
		ints, _ := s.Int()
		var arr []int
		for _, i := range ints {
			if !alreadyWin[i] {
				arr = append(arr, i)
			}
		}
		return series.Ints(arr)
	}

	var numGrids []int
	var count int
	for i := 0; i < df.Nrow(); i++ {
		if i != 0 && i%10 == 0 {
			count++
		}
		numGrids = append(numGrids, count)
	}

	var lastCall, lastWinner int
	for j, num := range input {
		number, _ := strconv.Atoi(num)
		nums[number] = true
		lastCall = number

		if j < 4 {
			continue
		}

		dfWinner := df.Rapply(bingo)
		dfWinner = dfWinner.Mutate(series.Ints(numGrids))

		dfWinner = dfWinner.Filter(
			dataframe.F{
				Colidx:     0,
				Comparator: series.Eq,
				Comparando: true,
			},
		)

		dfWinner = dfWinner.Select(dfWinner.Ncol() - 1)
		dfWinner = dfWinner.Capply(newWinner)

		for i := 0; i < dfWinner.Nrow(); i++ {
			j, _ := dfWinner.Elem(i, 0).Int()
			alreadyWin[j] = true
			lastWinner = j
		}

		if len(alreadyWin) == 100 {
			break
		}
	}

	winnerGrid := df.Subset(
		[]int{
			lastWinner * 10,
			lastWinner*10 + 1,
			lastWinner*10 + 2,
			lastWinner*10 + 3,
			lastWinner*10 + 4,
		})
	unmarkedNumberSum := winnerGrid.Capply(sum)
	unmarkedNumberSum = unmarkedNumberSum.Rapply(sum)

	uSum, _ := unmarkedNumberSum.Elem(0, 0).Int()
	fmt.Printf("LAST WINNER: The sum of unmarked number with last num call is %d\n", uSum*lastCall)

	return nil
}

func readDatas() (dataframe.DataFrame, []string, error) {
	file, err := os.Open("day4/data.txt")
	if err != nil {
		return dataframe.DataFrame{}, nil, fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	input := strings.Split(scanner.Text(), ",")
	var df dataframe.DataFrame
	for scanner.Scan() {
		var records [][]string
		for i := 0; i < 5; i++ {
			if scanner.Scan() {
				records = append(records, strings.Split(
					strings.ReplaceAll(
						strings.TrimSpace(
							scanner.Text()),
						"  ", " "),
					" "))
			}
		}
		dftemp := dataframe.LoadRecords(records,
			dataframe.DetectTypes(false),
			dataframe.DefaultType(series.Int),
			dataframe.HasHeader(false),
		)

		// Create transpose for columns
		var t [][]string
		transpose := matrix{dftemp}.T()
		for i := 0; i < dftemp.Ncol(); i++ {
			var arr []string
			for j := 0; j < dftemp.Ncol(); j++ {
				arr = append(arr, strconv.Itoa(int(transpose.At(i, j))))
			}
			t = append(t, arr)
		}
		dft := dataframe.LoadRecords(t,
			dataframe.DetectTypes(false),
			dataframe.DefaultType(series.Int),
			dataframe.HasHeader(false),
		)

		dftemp = dftemp.Concat(dft)
		df = df.Concat(dftemp)
	}
	if err := scanner.Err(); err != nil {
		return dataframe.DataFrame{}, nil, fmt.Errorf("scanner error: %s", err)
	}

	return df, input, nil
}

type matrix struct {
	dataframe.DataFrame
}

func (m matrix) At(i, j int) float64 {
	return m.Elem(i, j).Float()
}

func (m matrix) T() mat.Matrix {
	return mat.Transpose{Matrix: m}
}
