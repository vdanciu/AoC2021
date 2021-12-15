package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	var risk_map [][]int
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		row := make([]int, len(line))
		for i := range line {
			row[i], _ = strconv.Atoi(line[i])
		}
		risk_map = append(risk_map, row)
	}
	risk := search_path(risk_map, 0, 0)
	fmt.Printf("%v\n", risk)
}

func search_path(risk_map [][]int, x, y int) int {
	if risk_map[y][x] < 0 {
		//already been here and calculated the risk saved as a negative value
		return -risk_map[y][x]
	}
	risk := math.MaxInt
	if x == len(risk_map[0])-1 && y == len(risk_map)-1 {
		return risk_map[y][x]
	}
	if x < len(risk_map[0])-1 {
		risk = min(risk, search_path(risk_map, x+1, y))
	}
	if y < len(risk_map)-1 {
		risk = min(risk, search_path(risk_map, x, y+1))
	}
	total_risk := risk
	if x != 0 || y != 0 {
		total_risk += +risk_map[y][x]
	}
	risk_map[y][x] = -total_risk
	return total_risk
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dup(src [][]int) [][]int {
	dup := make([][]int, len(src))
	for i := range src {
		dup[i] = make([]int, len(src[i]))
		copy(dup[i], src[i])
	}
	return dup
}

func print(a [][]int) {
	for i := range a {
		for j := range a[i] {
			fmt.Printf("% d", a[i][j])
		}
		fmt.Printf("\n")
	}
}
