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
	close_symbol := map[string]string{
		")": "(",
		"]": "[",
		"}": "{",
		">": "<"}
	symbol_value := map[string]int{
		"(": 1,
		"[": 2,
		"{": 3,
		"<": 4}
	autocomplete_scores := []int{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		symbols := strings.Split(line, "")
		open_symbols := ""
		error_present := false
		for _, symbol := range symbols {
			if _, is := close_symbol[symbol]; is {
				if len(open_symbols) > 0 {
					if close_symbol[symbol] != string(open_symbols[len(open_symbols)-1]) {
						error_present = true
						break
					} else {
						open_symbols = open_symbols[:len(open_symbols)-1]
					}
				}
			} else {
				open_symbols += symbol
			}
		}
		if !error_present {
			score := 0
			for i := len(open_symbols) - 1; i >= 0; i-- {
				score *= 5
				score += symbol_value[string(open_symbols[i])]
			}
			autocomplete_scores = append(autocomplete_scores, score)
		}
	}
	sort.Ints(autocomplete_scores)
	fmt.Printf("result is %v\n", autocomplete_scores[len(autocomplete_scores)/2])
}
