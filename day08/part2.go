package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		input := strings.Fields(strings.Split(line, "|")[0])
		output := strings.Fields(strings.Split(line, "|")[1])

		key := decypher(input)
		number := read(output, key)
		total += number
	}
	fmt.Printf("result is %v\n", total)
}

func decypher(input []string) map[string]int {
	var segments [7]string
	len_map := make(map[int]string)
	for i := range input {
		l := len(input[i])
		if _, in := len_map[l]; in {
			len_map[l] = intersection(len_map[l], input[i])
		} else {
			len_map[l] = strings.TrimSpace(input[i])
		}
	}

	//  0000
	// 6    1
	// 6    1
	//  2222
	// 5    3
	// 5    3
	//  4444

	digit_map := make(map[int]string)
	digit_map[1] = len_map[2]
	digit_map[4] = len_map[4]
	digit_map[5] = len_map[6]
	digit_map[7] = len_map[3]
	digit_map[8] = len_map[7]

	segments[0] = difference(digit_map[7], digit_map[1])
	segments[3] = intersection(digit_map[5], digit_map[1])
	segments[1] = difference(digit_map[1], segments[3])
	segments[4] = difference(difference(len_map[5], digit_map[4]), digit_map[7])
	segments[2] = difference(difference(len_map[5], segments[4]), segments[0])

	digit_map[3] = uniform(segments[0] + segments[1] + segments[2] + segments[3] + segments[4])

	segments[6] = difference(digit_map[4], digit_map[3])
	segments[5] = difference(difference(digit_map[8], digit_map[3]), segments[6])

	key := make(map[string]int)
	key[make_digit(segments, 0, 1, 3, 4, 5, 6)] = 0
	key[make_digit(segments, 1, 3)] = 1
	key[make_digit(segments, 0, 1, 2, 5, 4)] = 2
	key[make_digit(segments, 0, 1, 2, 3, 4)] = 3
	key[make_digit(segments, 6, 2, 1, 3)] = 4
	key[make_digit(segments, 0, 6, 2, 3, 4)] = 5
	key[make_digit(segments, 0, 6, 5, 4, 3, 2)] = 6
	key[make_digit(segments, 0, 1, 3)] = 7
	key[make_digit(segments, 0, 1, 2, 3, 4, 5, 6)] = 8
	key[make_digit(segments, 2, 6, 0, 1, 3, 4)] = 9

	return key
}

func read(output []string, key map[string]int) int {
	result := 0
	for _, digit := range output {
		result *= 10
		result += key[uniform(digit)]
	}
	return result
}

func make_digit(segments [7]string, indices ...int) string {
	var digit string
	for _, index := range indices {
		digit += segments[index]
	}
	return uniform(digit)
}

func difference(s1 string, s2 string) string {
	var result string
	for _, r := range s1 {
		if !strings.Contains(s2, string(r)) {
			result += string(r)
		}
	}
	result = uniform(result)
	return result
}

func intersection(s1 string, s2 string) string {
	var result string
	for _, r := range s1 {
		if strings.Contains(s2, string(r)) {
			result += string(r)
		}
	}
	result = uniform(result)
	return result
}

func uniform(s string) string {
	sorted := strings.Split(s, "")
	sort.Strings(sorted)
	return strings.Join(sorted, "")
}
