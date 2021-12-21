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
	positions := []int{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		components := strings.Split(text, ": ")
		components[1] = strings.TrimSpace(components[1])
		position, err := strconv.Atoi(components[1])
		if err != nil {
			log.Fatal(err)
		}
		positions = append(positions, position)
	}

	playing := 0
	dice := 1
	count := 0
	score := []int{0, 0}
	for {
		count++
		roll := dice + dice + 1 + dice + 2
		dice += 3

		position := (roll+positions[playing]-1)%10 + 1

		positions[playing] = position
		score[playing] += position
		if score[playing] >= 1000 {
			break
		}
		playing = (playing + 1) % 2
	}
	fmt.Printf("result=%v\n", min(score[:])*count*3)

}

func min(a []int) int {
	min := math.MaxInt
	for _, n := range a {
		if n < min {
			min = n
		}
	}
	return min
}
