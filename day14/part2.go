package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	line_no := 0
	template := ""
	pairs := make(map[string]string)
	formula := make(map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		if line_no > 0 {
			if len(line) == 0 {
				continue
			}
			components := strings.Split(line, " -> ")
			pairs[components[0]] = components[1]
		} else {
			for i := 0; i < len(line)-1; i++ {
				add_value(formula, line[i:i+2], 1)
			}
			template = line
		}
		line_no++
	}

	new_formula := formula
	for step := 0; step < 40; step++ {
		new_formula = replace(new_formula, pairs)
	}

	histogram := make(map[byte]int)
	histogram[template[len(template)-1]] = 1
	for pair, value := range new_formula {
		if _, is := histogram[pair[0]]; !is {
			histogram[pair[0]] = 0
		}
		histogram[pair[0]] += value
	}
	min := math.MaxInt
	max := 0
	for _, v := range histogram {
		if min > v {
			min = v
		}
		if max < v {
			max = v
		}
	}

	fmt.Printf("result is %v", max-min)
}

func replace(formula map[string]int, pairs map[string]string) map[string]int {
	new_formula := make(map[string]int)
	for pair, value := range formula {
		if value > 0 {
			insert := pairs[pair]
			add_value(new_formula, string(pair[0])+insert, value)
			add_value(new_formula, insert+string(pair[1]), value)
		}
	}
	return new_formula
}

func add_value(formula map[string]int, key string, value int) {
	if _, in := formula[key]; !in {
		formula[key] = 0
	}
	formula[key] += value
}

func sum(formula map[string]int) int {
	sum := 1 //last letter is not accounted for in the formula
	for _, v := range formula {
		sum += v
	}
	return sum
}
