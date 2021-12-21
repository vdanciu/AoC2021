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
	count := 0

	configurations := map[[5]int]int{{positions[0], positions[1], 0, 0, 0}: 1}
	//dice_roll := [2]int{1, 2}
	dice_roll := [3]int{1, 2, 3}
	wins := [2]int{0, 0}
	for {
		//fmt.Printf("\n[%v,%v]configuration=%v\n", playing, count, configurations)
		count++
		for _, r1 := range dice_roll {
			for _, r2 := range dice_roll {
				for _, r3 := range dice_roll {
					for config := range configurations {
						if config[4] == playing {
							old_config := config
							config[0+playing] = (r1+r2+r3+config[0+playing]-1)%10 + 1
							config[2+playing] += config[0+playing]
							config[4] = (playing + 1) % 2
							configurations[config] += configurations[old_config]
						}
					}
				}
			}
		}
		for config := range configurations {
			if config[2+playing] >= 21 {
				wins[playing] += configurations[config]
				delete(configurations, config)
			}
			if config[4] == playing {
				delete(configurations, config)
			}
		}

		if len(configurations) == 0 {
			break
		}
		playing = (playing + 1) % 2
	}
	fmt.Printf("max wins = %v\n", max(wins[:]))
}

func max(a []int) int {
	min := math.MinInt
	for _, n := range a {
		if n > min {
			min = n
		}
	}
	return min
}
