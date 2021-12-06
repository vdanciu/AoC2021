package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	num_codes := 0
	counters := make([]int, 0, 100)
	for scanner.Scan() {
		num_codes++
		line := scanner.Text()
		for i := 0; i < len(line); i++ {
			if line[i] == '1' {
				if len(counters)-1 < i {
					counters = counters[:(i + 1)]
					counters[i] = 0
				}
				counters[i]++
			}
		}
	}
	gamma := 0
	epsilon := 0

	for i := 0; i < len(counters); i++ {
		gamma <<= 1
		epsilon <<= 1
		gamma_bit := 0
		if counters[i] > num_codes/2 {
			gamma_bit = 1
		}
		gamma += gamma_bit
		epsilon += 1 & ^gamma_bit
	}
	fmt.Printf("gamma = %v\n", gamma)
	fmt.Printf("epsilon = %v\n", epsilon)
	fmt.Printf("result = %v\n", epsilon*gamma)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
