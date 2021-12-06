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
	var fish_list []string
	for scanner.Scan() {
		line := scanner.Text()
		fish_list = strings.Split(line, ",")
		fmt.Println(line)
	}

	var fish_ages [9]int
	for i := 0; i < len(fish_list); i++ {
		age, err := strconv.Atoi(fish_list[i])
		if err != nil {
			log.Fatal(err)
		}
		fish_ages[age] += 1
	}

	fmt.Println(fish_ages)
	for day := 1; day <= 256; day++ {
		var new_ages [9]int
		for age := 0; age < 9; age++ {
			if age == 0 {
				new_ages[8] = fish_ages[0]
				new_ages[6] = fish_ages[0]
			} else {
				if fish_ages[age] > 0 {
					new_ages[age-1] += fish_ages[age]
				}
			}
		}
		fish_ages = new_ages
	}

	fmt.Printf("result is = %v\n", total(fish_ages))
	fmt.Println("THE END")
}

func total(a [9]int) int {
	sum := 0
	for i := 0; i < 9; i++ {
		sum += a[i]
	}
	return sum
}
