package day15

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"adventofcode/utils"
)

type node struct {
	row, col  int
	val       int8
	dist      int
	isInf     bool
	previous  *node
	heuristic int
	index     int
}

type PriorityQueue []*node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].heuristic < pq[j].heuristic
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func Day15() {
	fmt.Println("DAY15")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func process() error {
	for i, f := range []func(n *node) int{noHeuristic, distFromEnd} {
		if err := readDatas(); err != nil {
			return err
		}
		t := time.Now()
		AStar(f)
		timer := time.Since(t)
		if i == 0 {
			fmt.Printf("Part one: Without heuristic (Djikstra): %s\n", timer)
		} else {
			fmt.Printf("Part one: With heuristic (Astar): %s\n", timer)
		}
	}
	path, risk := retrievePath(last)
	fmt.Printf("First part: risk: %d, path: %v\n", risk-1, path)

	for i, f := range []func(n *node) int{noHeuristic, distFromEnd} {
		if err := generateMapSizeFive(); err != nil {
			return nil
		}
		t := time.Now()
		AStar(f)
		path, risk = retrievePath(last)
		timer := time.Since(t)
		if i == 0 {
			fmt.Printf("Part two: Without heuristic (Djikstra): %s\n", timer)
		} else {
			fmt.Printf("Part two: With heuristic (Astar): %s\n", timer)
		}
	}
	fmt.Printf("Second part: risk: %d, path: %v\n", risk-1, path)
	return nil
}

var (
	nodeMap     [][]*node
	seen        = map[int]map[int]bool{}
	start, last *node
)

func AStar(heuristic func(n *node) int) {
	pq := PriorityQueue{}
	heap.Push(&pq, start)
	heap.Init(&pq)

	for len(pq) > 0 {
		n := heap.Pop(&pq).(*node)
		if n == last {
			return
		}

		seen[n.row][n.col] = true

		for _, v := range getNeighboor(n) {
			if seen[v.row][v.col] {
				continue
			}

			alt := n.dist + int(v.val)
			if v.isInf || alt < v.dist {
				v.dist = alt
				v.previous = n
				v.isInf = false
				v.heuristic = v.dist + heuristic(v)
				heap.Push(&pq, v)
			}
		}
	}
}

func getNeighboor(n *node) []*node {
	arr := []*node{}
	if n.row > 0 {
		arr = append(arr, nodeMap[n.row-1][n.col])
	}
	if n.col > 0 {
		arr = append(arr, nodeMap[n.row][n.col-1])
	}
	if n.row < len(nodeMap)-1 {
		arr = append(arr, nodeMap[n.row+1][n.col])
	}
	if n.col < len(nodeMap[0])-1 {
		arr = append(arr, nodeMap[n.row][n.col+1])
	}
	return arr
}

func retrievePath(end *node) ([]int8, int) {
	path := []int8{}
	node := end
	risk := 0
	for node != nil {
		path = append(path, node.val)
		risk += int(node.val)
		node = node.previous
	}
	return path, risk
}

func readDatas() error {
	file, err := os.Open("day15/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	row := 0
	nodeMap = [][]*node{}
	for scanner.Scan() {
		arr, err := utils.StringToInt8Array(scanner.Text())
		if err != nil {
			return err
		}
		arrNode := make([]*node, len(arr))
		for i, v := range arr {
			n := &node{row: row, col: i, val: v, isInf: true}
			arrNode[i] = n
			if seen[n.row] == nil {
				seen[n.row] = map[int]bool{}
			}
			seen[n.row][n.col] = false
		}
		nodeMap = append(nodeMap, arrNode)
		row++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %s", err)
	}

	nodeMap[0][0].isInf = false
	start, last = nodeMap[0][0], nodeMap[len(nodeMap)-1][len(nodeMap[0])-1]
	seen[0][0] = false
	return nil
}

func generateMapSizeFive() error {
	if err := readDatas(); err != nil {
		return err
	}
	var result [][]*node
	for i := 0; i < len(nodeMap); i++ {
		result = append(result, make([]*node, 5*len(nodeMap[0])))
		for j := 0; j < len(nodeMap[0]); j++ {
			n := nodeMap[i][j]
			newVal := n.val
			for k := 0; k < 5; k++ {
				if newVal == 10 {
					newVal = 1
				}
				newNode := &node{
					row:   i,
					col:   j + k*len(nodeMap[0]),
					val:   newVal,
					isInf: true,
				}
				result[i][j+k*len(nodeMap[0])] = newNode
				newVal++
			}
		}
	}

	template := result
	fmt.Println()
	for k := 1; k < 5; k++ {
		arr := [][]*node{}
		for i := 0; i < len(template); i++ {
			a := []*node{}
			for j := 0; j < len(template[i]); j++ {
				newVal := template[i][j].val + 1
				if newVal == 10 {
					newVal = 1
				}
				a = append(a, &node{row: i + k*len(template), col: j, val: newVal, isInf: true})
			}
			arr = append(arr, a)
		}
		template = arr
		for _, v := range arr {
			result = append(result, v)
		}
	}

	nodeMap = result
	seen = map[int]map[int]bool{}
	for i := 0; i < len(nodeMap); i++ {
		for j := 0; j < len(nodeMap[0]); j++ {
			nodeMap[i][j].dist = 0
			nodeMap[i][j].isInf = true
			if seen[i] == nil {
				seen[i] = map[int]bool{}
			}
			seen[i][j] = false
		}
	}
	start, last = nodeMap[0][0], nodeMap[len(nodeMap)-1][len(nodeMap[0])-1]
	start.isInf = false
	seen[0][0] = false

	return nil
}

func distFromEnd(n *node) int {
	return int(math.Abs(float64(n.col)-float64(last.col)) + math.Abs(float64(n.row)-float64(last.row)))
}

func noHeuristic(n *node) int {
	return 0
}
