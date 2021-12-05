package day3

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func Day3() {
	fmt.Println("DAY3:")

	if err := firstPart(); err != nil {
		log.Fatal(err)
	}

	if err := secondPart(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

// 3923414
func firstPart() error {
	file, err := os.Open("day3/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records [][]string
	for scanner.Scan() {
		records = append(records, strings.Split(scanner.Text(), ""))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	df := dataframe.LoadRecords(records,
		dataframe.DetectTypes(false),
		dataframe.DefaultType(series.Int),
	)

	dfCounts := df.Capply(maxValue)

	gam, eps := 0, 0
	for i := 0; i < dfCounts.Ncol(); i++ {
		elem, _ := dfCounts.Elem(0, i).Int()
		gam += elem * int(math.Pow(2, float64(11-i)))
		eps += int(^uint8(elem)&1) * int(math.Pow(2, float64(11-i)))
	}

	fmt.Printf("Gamma: %d, Epsilon: %d, Multiplication: %d\n", gam, eps, gam*eps)
	return nil
}

//5852595
func secondPart() error {
	file, err := os.Open("day3/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records [][]string
	for scanner.Scan() {
		records = append(records, strings.Split(scanner.Text(), ""))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	df := dataframe.LoadRecords(records,
		dataframe.DetectTypes(false),
		dataframe.DefaultType(series.Int),
	)

	dfOx, dfCo2 := df, df
	for i := 0; i < dfOx.Ncol(); i++ {
		if dfOx.Nrow() > 1 {
			dfCounts := dfOx.Capply(maxValue)
			max, _ := dfCounts.Elem(0, i).Int()
			dfOx = dfOx.Filter(dataframe.F{
				Colidx:     i,
				Comparator: series.Eq,
				Comparando: max,
			})
		}

		if dfCo2.Nrow() > 1 {
			dfCounts := dfCo2.Capply(minValue)
			min, _ := dfCounts.Elem(0, i).Int()
			dfCo2 = dfCo2.Filter(dataframe.F{
				Colidx:     i,
				Comparator: series.Eq,
				Comparando: min,
			})
		}

		if dfOx.Nrow() == 1 && dfCo2.Nrow() == 1 {
			break
		}
	}

	ox, co2 := 0, 0
	for i := 0; i < dfOx.Ncol(); i++ {
		elem, _ := dfOx.Elem(0, i).Int()
		ox += elem * int(math.Pow(2, float64(11-i)))

		elem, _ = dfCo2.Elem(0, i).Int()
		co2 += elem * int(math.Pow(2, float64(11-i)))
	}

	fmt.Printf("Oxygene: %d, CO2: %d, Multiplication: %d\n", ox, co2, ox*co2)
	return nil
}

func maxValue(s series.Series) series.Series {
	zeros, ones := oneZeroCounter(s)

	if ones >= zeros {
		return series.Ints(1)
	}
	return series.Ints(0)
}

func minValue(s series.Series) series.Series {
	zeros, ones := oneZeroCounter(s)

	if ones >= zeros {
		return series.Ints(0)
	}
	return series.Ints(1)
}

func oneZeroCounter(s series.Series) (int, int) {
	ints, err := s.Int()
	if err != nil {
		panic("can't transform into int !")
	}

	zeros, ones := 0, 0
	for _, i := range ints {
		if i == 0 {
			zeros++
		} else {
			ones++
		}
	}

	return zeros, ones
}
