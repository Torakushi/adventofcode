package day19

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gonum.org/v1/gonum/mat"
)

var (
	scanners     []*scanner
	transfMatrix []*mat.Dense
	totalBeacons = map[string]struct{}{}
)

func Day19() {
	fmt.Println("DAY19")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func process() error {
	readDatas()
	buildRotationsMatrix()

	var lock sync.Mutex
	done := map[int]bool{0: true}
	todo := []int{0}
	t := time.Now()
	for len(todo) > 0 {
		s := scanners[todo[0]]
		todo = todo[1:]
		var wg sync.WaitGroup
		for i := 0; i < len(scanners); i++ {
			if done[i] {
				continue
			}
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				pos := compare(s, scanners[i])
				if pos != nil {
					lock.Lock()
					todo = append(todo, i)
					done[i] = true
					lock.Unlock()
				}
			}(i)
		}
		wg.Wait()
	}

	fmt.Printf("Finish processing in %s\n", time.Since(t))
	fmt.Printf("The total number of beacons are %d\n", len(totalBeacons))

	scI, scJ, maxD := getLargestManhattan()
	fmt.Printf("The largest manhattan disance is between Scanner%d and Scanner%d and is: %d\n", scI, scJ, maxD)
	return nil
}

type scanner struct {
	researchMap map[string]bool
	beacons     *mat.Dense
	pos         []float64
}

func readDatas() error {
	file, err := os.Open("day19/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	arr := []float64{}
	rMap := map[string]bool{}
	for sc.Scan() {
		txt := sc.Text()
		if txt == "" {
			continue
		}
		if strings.HasPrefix(txt, "---") {
			if len(arr) != 0 {
				scanners = append(scanners,
					&scanner{
						researchMap: rMap,
						beacons:     mat.NewDense(len(arr)/3, 3, arr),
						pos:         []float64{0, 0, 0},
					},
				)
			}
			arr = []float64{}
			rMap = map[string]bool{}
			continue
		}

		coordAsStr := strings.Split(txt, ",")
		rMap[txt] = true
		// Initialize 'totalBeacons' with scanner 0 beacons
		if len(scanners) == 0 {
			totalBeacons[txt] = struct{}{}
		}
		for _, s := range coordAsStr {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			arr = append(arr, i)
		}
	}

	scanners = append(scanners,
		&scanner{
			researchMap: rMap,
			beacons:     mat.NewDense(len(arr)/3, 3, arr),
		},
	)

	return nil
}

func buildRotationsMatrix() {
	permutes := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	p := []*mat.Dense{}
	for _, perm := range permutes {
		arr := make([]float64, 9)
		for i := 0; i < 3; i++ {
			arr[perm[i]+(3*i)] = 1
		}
		p = append(p, mat.NewDense(3, 3, arr))
	}

	basicRotation := []*mat.Dense{}
	for _, i := range [...]float64{1, -1} {
		for _, j := range [...]float64{1, -1} {
			for _, k := range [...]float64{1, -1} {
				basicRotation = append(basicRotation,
					mat.NewDense(3, 3, []float64{i, 0, 0, 0, j, 0, 0, 0, k}))
			}
		}
	}

	transfMatrix = []*mat.Dense{}
	for _, permMat := range p {
		for _, rot := range basicRotation {
			var c mat.Dense
			c.Mul(permMat, rot)
			if mat.Det(&c) == 1 {
				transfMatrix = append(transfMatrix, &c)
			}
		}
	}
}

func compare(scBase, scOther *scanner) []float64 {
	rNumBase, _ := scBase.beacons.Dims()
	rNumOther, _ := scOther.beacons.Dims()
	var pos []float64
	var wg0 sync.WaitGroup
	for _, tm := range transfMatrix {
		wg0.Add(1)
		go func(tm *mat.Dense) {
			defer wg0.Done()
			if pos != nil {
				return
			}
			var tOther mat.Dense
			tOther.Mul(scOther.beacons, mat.Transpose{Matrix: tm})
			var l sync.Mutex
			var wg1 sync.WaitGroup
			for i := 0; i < rNumBase; i++ {
				wg1.Add(1)
				go func(i int) {
					defer wg1.Done()
					if pos != nil {
						return
					}
					bForTranslation := createMatrixWithSameRaw(rNumOther, scBase.beacons.RawRowView(i))
					var wg2 sync.WaitGroup
					for j := 0; j < rNumOther; j++ {
						wg2.Add(1)
						go func(j int) {
							defer wg2.Done()
							if pos != nil {
								return
							}
							oForTranslation := createMatrixWithSameRaw(rNumOther, tOther.RawRowView(j))
							var tMatrix, translatedMatrix mat.Dense
							tMatrix.Sub(bForTranslation, oForTranslation)
							translatedMatrix.Add(&tOther, &tMatrix)

							count := 0
							for k := 0; k < rNumOther; k++ {
								r := translatedMatrix.RawRowView(k)
								if scBase.researchMap[arrAsString(r)] {
									count++
								}
								if count >= 12 {
									l.Lock()
									defer l.Unlock()
									if pos != nil {
										return
									}
									scOther.changeBasis(&translatedMatrix)
									scOther.pos = tMatrix.RawRowView(0)
									pos = scOther.pos
									return
								}
							}
						}(j)
					}
					wg2.Wait()
				}(i)
				wg1.Wait()
			}
		}(tm)
		wg0.Wait()
	}
	return pos
}

func (s *scanner) changeBasis(m *mat.Dense) {
	s.beacons = m
	s.researchMap = map[string]bool{}
	r, _ := m.Dims()
	for i := 0; i < r; i++ {
		arrStr := arrAsString(m.RawRowView(i))
		totalBeacons[arrStr] = struct{}{}
		s.researchMap[arrStr] = true
	}
}

func createMatrixWithSameRaw(numRow int, r []float64) *mat.Dense {
	res := []float64{}
	for i := 0; i < numRow; i++ {
		res = append(res, r...)
	}
	return mat.NewDense(numRow, 3, res)
}

func getLargestManhattan() (int, int, int) {
	manhattan := func(a, b []float64) int {
		sum := 0
		for i := 0; i < len(a); i++ {
			sum += int(math.Abs(b[i] - a[i]))
		}
		return sum
	}

	scI, scJ, max := 0, 0, 0
	for i, sc := range scanners {
		for j, oSc := range scanners {
			if i == j {
				continue
			}
			d := manhattan(sc.pos, oSc.pos)
			if d > max {
				scI, scJ, max = i, j, d
			}
		}
	}
	return scI, scJ, max
}

func arrAsString(arr []float64) string {
	return fmt.Sprintf("%d,%d,%d", int(arr[0]), int(arr[1]), int(arr[2]))
}
