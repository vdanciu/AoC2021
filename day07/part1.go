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

	var positions []int
	min_position := math.MaxInt64
	max_position := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		components := strings.Split(scanner.Text(), ",")
		for i := 0; i < len(components); i++ {
			num, _ := strconv.Atoi(components[i])
			positions = append(positions, num)
			if num < min_position {
				min_position = num
			}
			if num > max_position {
				max_position = num
			}
		}
	}

	min_move := math.MaxInt
	target := 0
	for t := min_position; t <= max_position; t++ {
		move := compute_move(positions, t)

		if move < min_move {
			min_move = move
			target = t
		}
	}
	fmt.Printf("t = %v, with move = %v\n", target, min_move)
}

func compute_move(a []int, t int) int {
	move := 0
	for i := 0; i < len(a); i++ {
		if a[i] > t {
			move += a[i] - t
		} else {
			move += t - a[i]
		}
	}
	return move
}
