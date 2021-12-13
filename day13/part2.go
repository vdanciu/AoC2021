package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}
type fold struct {
	axis  string
	value int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	const (
		READ_DOTS  = iota
		READ_FOLDS = iota
	)
	mode := READ_DOTS
	matrix := make(map[point]bool)
	var folds []fold
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		switch mode {
		case READ_DOTS:
			if len(line) > 0 {
				matrix[make_point(line)] = true
			} else {
				mode = READ_FOLDS
			}
		case READ_FOLDS:
			folds = append(folds, make_fold(line))
		}
	}

	for _, fold := range folds {
		folded := make(map[point]bool)
		for p := range matrix {
			if fold.axis == "x" {
				if p.x > fold.value {
					folded[point{x: p.x + (fold.value-p.x)*2, y: p.y}] = true
					continue
				}
			} else {
				if p.y > fold.value {
					folded[point{x: p.x, y: p.y + (fold.value-p.y)*2}] = true
					continue
				}
			}
			folded[point{x: p.x, y: p.y}] = true
		}
		// print_matrix(folded)
		// fmt.Printf("folded %v\n", folded)
		matrix = folded
	}

	// fmt.Printf("matrix = %v\n", matrix)
	// fmt.Printf("folds = %v\n", folds)
	print_matrix(matrix)
}

func make_point(line string) point {
	components := strings.Split(line, ",")
	return point{x: atoi(components[0]), y: atoi(components[1])}
}

func make_fold(line string) fold {
	components := strings.Split(line, "=")
	return fold{axis: string(components[0][len(components[0])-1]), value: atoi(components[1])}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func print_matrix(matrix map[point]bool) {
	width := -1
	height := -1
	for p := range matrix {
		if p.x > width {
			width = p.x
		}
		if p.y > height {
			height = p.y
		}
	}
	width++
	height++

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			_, present := matrix[point{x: col, y: row}]
			if present {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}
