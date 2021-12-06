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

	numbers := []string{}
	numbers_width := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		number := scanner.Text()
		if len(number) > numbers_width {
			numbers_width = len(number)
		}
		numbers = append(numbers, number)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	o2rate := find_number(numbers, numbers_width, o2compare)
	co2rate := find_number(numbers, numbers_width, co2compare)

	fmt.Printf("o2rate = %v\n", o2rate)
	fmt.Printf("co2rate = %v\n", co2rate)
	fmt.Printf("result = %v\n", co2rate*o2rate)
}

func find_number(src []string, numbers_width int, fn compare) int {
	rate := make([]string, len(src))
	copy(rate, src)
	for p := 0; p < numbers_width; p++ {
		a_1 := []string{}
		a_0 := []string{}
		for i := 0; i < len(rate); i++ {
			if rate[i][p] == '1' {
				a_1 = append(a_1, rate[i])
			} else {
				a_0 = append(a_0, rate[i])
			}
		}

		rate = a_1
		if fn(len(a_0), len(a_1)) {
			rate = a_0
		}

		if len(rate) == 1 {
			return strbit_to_number(rate[0])
		}
	}
	panic("this is not suppposed to happen")
}

type compare func(int, int) bool

func o2compare(a int, b int) bool {
	return a > b
}

func co2compare(a int, b int) bool {
	return a <= b
}

func strbit_to_number(number string) int {
	result := 0
	for i := 0; i < len(number); i++ {
		bit := 0
		if number[i] == '1' {
			bit = 1
		}
		result <<= 1
		result += bit
	}
	return result
}
