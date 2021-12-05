package day5

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/gonum/mat"
)

func Day5() {
	fmt.Println("DAY5")

	if err := withDataframe(false); err != nil {
		log.Fatal(err)
	}

	if err := withDataframe(true); err != nil {
		log.Fatal(err)
	}

	if err := withMat(false); err != nil {
		log.Fatal(err)
	}
	if err := withMat(true); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func withDataframe(withDiagonal bool) error {
	file, err := os.Open("day5/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	twosCounter := 0
	t := time.Now()
	addOne := func(x, y int) func(s series.Series) series.Series {
		return func(s series.Series) series.Series {
			from, to := 0, 0
			if x < y {
				from = x
				to = y
			} else {
				from = y
				to = x
			}
			ints, _ := s.Int()

			for i := 0; i < len(ints); i++ {
				if i >= from && i <= to {
					ints[i]++
					if ints[i] == 2 {
						twosCounter++
					}
				}
			}
			return series.Ints(ints)
		}
	}

	df := initDataFrame()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " -> ")
		coordStr := strings.Split(args[0], ",")
		coordStr = append(coordStr, strings.Split(args[1], ",")...)

		var x0, y0, x1, y1 int
		for i := 0; i < 4; i++ {
			n, err := strconv.Atoi(coordStr[i])
			if err != nil {
				return err
			}

			switch i {
			case 0:
				x0 = n
			case 1:
				y0 = n
			case 2:
				x1 = n
			case 3:
				y1 = n
			}
		}

		// row
		if y0 == y1 {
			row := df.Subset(y0)
			row = row.Rapply(addOne(x0, x1))
			df = df.Set(y0, row)
		}

		if x0 == x1 {
			//Column
			col := df.Select(x0)
			col = col.Capply(addOne(y0, y1))
			df = df.Mutate(col.Col(col.Names()[0]))
		}

		if withDiagonal && isDiagonal(x0, y0, x1, y1) {
			var xFrom, xTo, yFrom, yTo int
			if x0 < x1 {
				xFrom, xTo, yFrom, yTo = x0, x1, y0, y1
			} else {
				xFrom, xTo, yFrom, yTo = x1, x0, y1, y0
			}

			var yMin, yMax int
			if y0 < y1 {
				yMin, yMax = y0, y1
			} else {
				yMin, yMax = y1, y0
			}

			colnames := make(map[string]bool)
			for i := xFrom; i <= xTo; i++ {
				colnames[df.Names()[i]] = true
			}

			var rowsIndexes []int
			for i := yMin; i <= yMax; i++ {
				rowsIndexes = append(rowsIndexes, i)
			}

			incrementer := float64(yTo-yFrom) / math.Abs(float64(yTo-yFrom))
			colCounter := yFrom
			df = df.Capply(
				func(s series.Series) series.Series {
					if !colnames[s.Name] {
						return s
					}

					ints, _ := s.Int()
					ints[colCounter] += 1
					if ints[colCounter] == 2 {
						twosCounter++
					}
					colCounter += int(incrementer)

					return series.Ints(ints)
				})
		}
	}
	fmt.Printf("Using dataframe, with diagonal : %v, number of points (>=2): %d, took %s\n\n",
		withDiagonal, twosCounter, time.Since(t))

	return nil
}

func initDataFrame() dataframe.DataFrame {
	records := make([][]string, 1000)
	for i := 0; i < 1000; i++ {
		arr := make([]string, 1000)
		for j := 0; j < 1000; j++ {
			arr[j] = "0"
		}
		records[i] = arr
	}

	return dataframe.LoadRecords(records,
		dataframe.DetectTypes(false),
		dataframe.DefaultType(series.Int),
	)
}

func isDiagonal(x0, y0, x1, y1 int) bool {
	return x0 != x1 &&
		y0 != y1 &&
		math.Abs(float64(x0-x1)) == math.Abs(float64(y0-y1))
}

func withMat(withDiagonal bool) error {
	file, err := os.Open("day5/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	m := mat.NewDense(1000, 1000, nil)
	scanner := bufio.NewScanner(file)
	twosCounter := 0
	t := time.Now()
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " -> ")
		coordStr := strings.Split(args[0], ",")
		coordStr = append(coordStr, strings.Split(args[1], ",")...)

		var x0, y0, x1, y1 int
		for i := 0; i < 4; i++ {
			n, err := strconv.Atoi(coordStr[i])
			if err != nil {
				return err
			}

			switch i {
			case 0:
				x0 = n
			case 1:
				y0 = n
			case 2:
				x1 = n
			case 3:
				y1 = n
			}
		}
		// row
		if y0 == y1 {
			var xMin, xMax int
			if x0 < x1 {
				xMin, xMax = x0, x1
			} else {
				xMin, xMax = x1, x0
			}
			m.Apply(func(i, j int, v float64) float64 {
				if j != y0 || i < xMin || i > xMax {
					return v
				}
				if v+1 == 2 {
					twosCounter++
				}
				return v + 1
			}, m)
		}

		if x0 == x1 {
			var yMin, yMax int
			if y0 < y1 {
				yMin, yMax = y0, y1
			} else {
				yMin, yMax = y1, y0
			}
			m.Apply(func(i, j int, v float64) float64 {
				if i != x0 || j < yMin || j > yMax {
					return v
				}
				if v+1 == 2 {
					twosCounter++
				}
				return v + 1
			}, m)
		}

		if withDiagonal && isDiagonal(x0, y0, x1, y1) {
			var xFrom, yFrom, xTo, yTo int
			if x0 < x1 {
				xFrom, yFrom, xTo, yTo = x0, y0, x1, y1
			} else {
				xFrom, yFrom, xTo, yTo = x1, y1, x0, y0
			}

			incrementer := int(float64(yTo-yFrom) / math.Abs(float64(yTo-yFrom)))
			colCounter := yFrom
			done := make(map[int]bool)
			m.Apply(func(i, j int, v float64) float64 {
				if j != colCounter || done[i] || i < xFrom || i > xTo {
					return v
				}
				done[i] = true
				if v+1 == 2 {
					twosCounter++
				}
				colCounter += incrementer
				return v + 1
			}, m)
		}
	}
	fmt.Printf("Using matrix, with diagonal : %v, number of points (>=2): %d took %s\n\n",
		withDiagonal, twosCounter, time.Since(t))
	return nil
}
