package main

import (
	"bufio"
	"fmt"
	"log"
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
	var height_map [][]int
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		height_line := make([]int, len(line))
		for i := range height_line {
			height_line[i], _ = strconv.Atoi(line[i])
		}
		height_map = append(height_map, height_line)
	}
	risk := 0
	for row := range height_map {
		for col := range height_map[row] {
			adjacent := get_adjacent(height_map, row, col)
			is_low := true
			for i := range adjacent {
				if adjacent[i] <= height_map[row][col] {
					is_low = false
					break
				}
			}
			if is_low {
				risk += height_map[row][col] + 1
			}
		}
	}
	fmt.Printf("risk = %v\n", risk)
}

func get_adjacent(a [][]int, row int, col int) []int {
	var result []int
	height := len(a)
	width := len(a[0])
	if row > 0 {
		result = append(result, a[row-1][col])
	}
	if col > 0 {
		result = append(result, a[row][col-1])
	}
	if row < (height - 1) {
		result = append(result, a[row+1][col])
	}
	if col < (width - 1) {
		result = append(result, a[row][col+1])
	}

	return result
}
