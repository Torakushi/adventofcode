package day18

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func Day18() {
	fmt.Println("DAY18")
	if err := processFirstPart(); err != nil {
		log.Fatal(err)
	}

	if err := processSecondPart(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func processFirstPart() error {
	file, err := os.Open("day18/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	var arr []interface{}
	json.Unmarshal([]byte(scanner.Text()), &arr)
	p := buildPairs(arr, 1)
	for scanner.Scan() {
		var arr []interface{}
		json.Unmarshal([]byte(scanner.Text()), &arr)
		np := buildPairs(arr, 1)
		p = addPairs(p, np)
		p.process()
	}

	fmt.Printf("The reduced fish is:%v\n magnitude of full sum is: %d\n\n", p.toArray(), p.getMagnitude())
	return nil
}

func processSecondPart() error {
	file, err := os.Open("day18/data.txt")
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	pairs := []*pairs{}
	for scanner.Scan() {
		var arr []interface{}
		json.Unmarshal([]byte(scanner.Text()), &arr)
		np := buildPairs(arr, 1)
		pairs = append(pairs, np)
	}

	max := 0
	goodI, goodJ := 0, 0
	for i := 0; i < len(pairs); i++ {
		for j := 0; j < len(pairs); j++ {
			if i == j {
				continue
			}
			pc := addPairs(pairs[i].Copy(), pairs[j].Copy())
			pc.process()
			mg := pc.getMagnitude()
			if mg > max {
				goodI = i
				goodJ = j
				max = mg
			}
		}
	}
	fmt.Printf("SECOND part, the best Magnitude of the sum of two != fishes is: %d\n", max)
	fmt.Printf("It is for %v + %v", pairs[goodI].toArray(), pairs[goodJ].toArray())
	return nil
}

type pairs struct {
	depth                int
	onlyInt              bool
	pair                 []interface{}
	propagateExplToRight int
	propagateExplToLeft  int
}

func addPairs(p1, p2 *pairs) *pairs {
	p1.addDepth()
	p2.addDepth()
	return &pairs{
		depth: 1,
		pair:  []interface{}{p1, p2},
	}
}

func (p *pairs) Copy() *pairs {
	res := &pairs{
		onlyInt:              p.onlyInt,
		propagateExplToRight: p.propagateExplToRight,
		propagateExplToLeft:  p.propagateExplToLeft,
		depth:                p.depth,
		pair:                 []interface{}{},
	}

	for _, v := range p.pair {
		if _, ok := v.(int); ok {
			res.pair = append(res.pair, v)
			continue
		}
		res.pair = append(res.pair, v.(*pairs).Copy())
	}

	return res
}

func (p *pairs) addDepth() {
	p.depth += 1
	for _, v := range p.pair {
		if np, ok := v.(*pairs); ok {
			np.addDepth()
		}
	}
}

func (p *pairs) toArray() []interface{} {
	res := []interface{}{}
	for _, i := range p.pair {
		if _, ok := i.(int); ok {
			res = append(res, i)
			continue
		}
		ll, _ := i.(*pairs)
		res = append(res, ll.toArray())
	}
	return res
}

func (p *pairs) getMagnitude() int {
	res := 0
	for index, p := range p.pair {
		if i, ok := p.(int); ok {
			res += (3 - index) * i
			continue
		}
		ll, _ := p.(*pairs)
		res += (3 - index) * ll.getMagnitude()
	}
	return res
}

func (p *pairs) process() {
	for {
		for p.processActionOnPairs(false) {
		}
		ok := p.processActionOnPairs(true)
		if !ok {
			break
		}
	}
}

func buildPairs(arr []interface{}, depth int) *pairs {
	res := &pairs{
		depth:   depth,
		onlyInt: true,
	}
	for _, v := range arr {
		if f, ok := v.(float64); ok {
			res.pair = append(res.pair, int(f))
		} else {
			pArr, _ := v.([]interface{})
			p := buildPairs(pArr, depth+1)
			res.pair = append(res.pair, p)
			res.onlyInt = false
		}
	}
	return res
}

func (p *pairs) processActionOnPairs(canSplit bool) bool {
	if p.depth >= 4 && !p.onlyInt {
		p.pair, p.propagateExplToLeft, p.propagateExplToRight = explode(p.pair)
		_, ok := p.pair[0].(int)
		_, ok2 := p.pair[1].(int)
		p.onlyInt = ok && ok2
		return true
	}

	var actionDone bool
	for index, v := range p.pair {
		if f, ok := v.(int); ok {
			if f >= 10 && canSplit {
				p.pair[index] = split(f, p.depth)
				p.onlyInt = false
				return true
			}
			continue
		}

		np, _ := v.(*pairs)
		actionDone = np.processActionOnPairs(canSplit)
		// No action done, continue
		if !actionDone {
			continue
		}

		// Action done but no explosion to propagate
		if np.propagateExplToLeft == 0 && np.propagateExplToRight == 0 {
			return true
		}

		// propagate to left
		if np.propagateExplToLeft != 0 {
			if index == 1 {
				if ll, ok := p.pair[0].(*pairs); ok {
					ll.addToFirstRightValue(np.propagateExplToLeft)
					np.propagateExplToLeft = 0
					return true
				} else if f, ok := p.pair[0].(int); ok {
					p.pair[0] = f + np.propagateExplToLeft
					np.propagateExplToLeft = 0
					return true
				}
			}
			if p.depth != 1 {
				p.propagateExplToLeft = np.propagateExplToLeft
			}
			np.propagateExplToLeft = 0
			return true
		}

		// propagate to right
		if np.propagateExplToRight != 0 {
			if index == 0 {
				if ll, ok := p.pair[1].(*pairs); ok {
					ll.addToFirstLeftValue(np.propagateExplToRight)
					np.propagateExplToRight = 0
					return true
				} else if f, ok := p.pair[1].(int); ok {
					p.pair[1] = f + np.propagateExplToRight
					np.propagateExplToRight = 0
					return true
				}
			}
			if p.depth != 1 {
				p.propagateExplToRight = np.propagateExplToRight
			}
			np.propagateExplToRight = 0
			return true
		}
	}
	return actionDone
}

func (p *pairs) addToFirstRightValue(value int) bool {
	if v, ok := p.pair[1].(int); ok {
		p.pair[1] = v + value
		return true
	}
	ll, _ := p.pair[1].(*pairs)
	return ll.addToFirstRightValue(value)
}

func (p *pairs) addToFirstLeftValue(value int) bool {
	if v, ok := p.pair[0].(int); ok {
		p.pair[0] = v + value
		return true
	}
	ll, _ := p.pair[0].(*pairs)
	return ll.addToFirstLeftValue(value)
}

func explode(arr []interface{}) ([]interface{}, int, int) {
	if f, ok := arr[0].(int); ok {
		p, _ := arr[1].(*pairs)
		left, _ := p.pair[0].(int)
		right, _ := p.pair[1].(int)
		return []interface{}{f + left, 0}, 0, right
	}

	if f, ok := arr[1].(int); ok {
		p, _ := arr[0].(*pairs)
		right, _ := p.pair[1].(int)
		left, _ := p.pair[0].(int)
		return []interface{}{0, f + right}, left, 0
	}

	p1, _ := arr[0].(*pairs)
	right, _ := p1.pair[1].(int)
	left, _ := p1.pair[0].(int)
	p2, _ := arr[1].(*pairs)
	p2.pair[0] = right + p2.pair[0].(int)
	return []interface{}{0, p2}, left, 0
}

func split(f, depth int) *pairs {
	return &pairs{
		depth:   depth + 1,
		pair:    []interface{}{int(f / 2), int(f) - int(f/2)},
		onlyInt: true,
	}
}
