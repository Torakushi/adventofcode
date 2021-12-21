package day20

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Day20() {
	fmt.Println("DAY20")
	if err := process(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func process() error {
	readDatas()

	o := img.getOutputImage().getOutputImage()
	fmt.Printf("After enhance twice we have %d pixels lit\n", o.getLightPixels())

	s := img
	for i := 1; i <= 50; i++ {
		s = s.getOutputImage()
	}

	fmt.Printf("After enhance 50 times we have %d pixels lit\n", s.getLightPixels())
	return nil
}

var (
	decoder []rune
	img     *image
)

type image struct {
	minRow, maxRow  int
	minCol, maxCol  int
	mapping         map[int]map[int]struct{}
	infiniteIsLight bool
}

func (i *image) getOutputImage() *image {
	newRowMin, newRowMax, newColMin, newColMax := 0, 0, 0, 0
	newMapping := map[int]map[int]struct{}{}
	for row := i.minRow - 3; row <= i.maxRow+3; row++ {
		for col := i.minCol - 3; col <= i.maxCol+3; col++ {
			r := i.getOutput(row, col)
			if r == '#' {
				if newMapping[row] == nil {
					newMapping[row] = map[int]struct{}{}
				}
				newMapping[row][col] = struct{}{}
				if row < newRowMin {
					newRowMin = row
				}
				if row > newRowMax {
					newRowMax = row
				}
				if col < newColMin {
					newColMin = col
				}
				if col > newColMax {
					newColMax = col
				}
			}
		}
	}

	return &image{
		minRow:          newRowMin,
		minCol:          newColMin,
		maxRow:          newRowMax,
		maxCol:          newColMax,
		mapping:         newMapping,
		infiniteIsLight: (!i.infiniteIsLight && decoder[0] == '#'),
	}
}

func (i *image) getOutput(row, col int) rune {
	var sb strings.Builder
	for j := -1; j <= 1; j++ {
		for k := -1; k <= 1; k++ {
			if i.infiniteIsLight && (row+j < i.minRow || row+j > i.maxRow || col+k < i.minCol || col+k > i.maxCol) {
				sb.WriteString("1")
				continue
			}
			if _, ok := i.mapping[row+j][col+k]; ok {
				sb.WriteString("1")
				continue
			}
			sb.WriteString("0")
		}
	}
	v, err := strconv.ParseInt(sb.String(), 2, 64)
	if err != nil {
		panic(err)
	}
	return decoder[v]
}

func (i *image) getLightPixels() int {
	sum := 0
	for _, v := range i.mapping {
		sum += len(v)
	}
	return sum
}

func readDatas() error {
	file, err := os.Open("day20/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// get decoder, put in rune for memory
	scanner.Scan()
	arr := strings.Split(scanner.Text(), "")
	decoder = make([]rune, 512)
	for i, s := range arr {
		decoder[i] = []rune(s)[0]
	}
	scanner.Scan()

	m := map[int]map[int]struct{}{}
	maxRow, maxCol := 0, 0
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), "")
		for col, s := range arr {
			if s == "." {
				continue
			}
			if m[maxRow] == nil {
				m[maxRow] = map[int]struct{}{}
			}
			m[maxRow][col] = struct{}{}
		}
		maxRow++
		maxCol = len(arr) - 1
	}

	img = &image{
		minRow:  -3,
		minCol:  -3,
		maxCol:  maxCol + 3,
		maxRow:  maxRow + 2,
		mapping: m,
	}
	return nil
}
