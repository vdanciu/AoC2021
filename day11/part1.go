package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type matrix [10][10]int

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	var matrix matrix
	scanner := bufio.NewScanner(f)
	row := 0
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		for col, num := range line {
			val, _ := strconv.Atoi(num)
			matrix[row][col] = val
		}
		row++
	}
	fmt.Printf("%v", matrix)
	flashes := 0
	for step := 0; step < 100; step++ {
		for row := 0; row < 10; row++ {
			for col := 0; col < 10; col++ {
				matrix[row][col] = (matrix[row][col] + 1) % 10
			}
		}

		frozen := matrix
		for row := 0; row < 10; row++ {
			for col := 0; col < 10; col++ {
				if frozen[row][col] == 0 {
					flash(&matrix, row, col, &flashes)
				}
			}
		}
	}
	fmt.Printf("result is %v", flashes)
}

func flash(matrix *matrix, row, col int, flashes *int) {
	(*flashes)++
	r := []int{-1, 0, 1}
	for _, i := range r {
		for _, j := range r {
			if !(i == 0 && j == 0) && row+i >= 0 && col+j >= 0 && row+i < 10 && col+j < 10 {
				if (*matrix)[row+i][col+j] != 0 {
					(*matrix)[row+i][col+j] = ((*matrix)[row+i][col+j] + 1) % 10
					if (*matrix)[row+i][col+j] == 0 {
						flash(matrix, row+i, col+j, flashes)
					}
				}
			}
		}
	}
}
