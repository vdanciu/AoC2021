package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type window struct {
	row_min int
	row_max int
	col_min int
	col_max int
}

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
	var top3_lens [3]int
	window := window{row_min: 0, col_min: 0, row_max: len(height_map) - 1, col_max: len(height_map[0]) - 1}
	for row := range height_map {
		for col := range height_map[row] {
			if height_map[row][col] == 9 {
				continue
			}
			adjacent := get_adjacent(height_map, row, col)
			is_low := true
			for i := range adjacent {
				if adjacent[i] <= height_map[row][col] {
					is_low = false
					break
				}
			}
			if is_low {
				basin_size := find_basin(height_map, row, col, window)
				for i := 0; i < 3; i++ {
					if basin_size > top3_lens[i] {
						top3_lens[0] = basin_size
						break
					}
				}
				sort.Ints(top3_lens[:])
			}
		}
	}
	fmt.Printf("size = %v from %v\n", top3_lens[0]*top3_lens[1]*top3_lens[2], top3_lens)
}

func get_adjacent(a [][]int, row int, col int) []int {
	var result []int
	height := len(a)
	width := len(a[0])
	if row > 0 {
		result = append(result, a[row-1][col])
	}
	if row < (height - 1) {
		result = append(result, a[row+1][col])
	}
	if col > 0 {
		result = append(result, a[row][col-1])
	}
	if col < (width - 1) {
		result = append(result, a[row][col+1])
	}

	return result
}

// the whole 'window' business is pointless, remanent from a previous failed
// attempt to incorrectly partition the search space
func find_basin(a [][]int, row int, col int, window window) int {
	val := a[row][col]
	size := 1
	// if already visited
	if val > 9 {
		return 0
	} else {
		// add 10 for visited
		a[row][col] += 10
	}

	w := window
	if row > window.row_min && is_going_up(val, a[row-1][col]) {
		size += find_basin(a, row-1, col, w)
	}
	if col > window.col_min && is_going_up(val, a[row][col-1]) {
		size += find_basin(a, row, col-1, w)
	}
	if row < window.row_max && is_going_up(val, a[row+1][col]) {
		size += find_basin(a, row+1, col, w)
	}
	if col < window.col_max && is_going_up(val, a[row][col+1]) {
		size += find_basin(a, row, col+1, w)
	}

	return size
}

func is_going_up(from int, to int) bool {
	if to > 9 {
		to -= 10
	}
	return (to != 9) && to >= from
}
