package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		")": 3,
		"]": 57,
		"}": 1197,
		">": 25137}
	error_score := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		symbols := strings.Split(line, "")
		open_symbols := ""
		for _, symbol := range symbols {
			if _, is := close_symbol[symbol]; is {
				if len(open_symbols) > 0 {
					if close_symbol[symbol] != string(open_symbols[len(open_symbols)-1]) {
						error_score += symbol_value[symbol]
						break
					} else {
						open_symbols = open_symbols[:len(open_symbols)-1]
					}
				}
			} else {
				open_symbols += symbol
			}
		}
	}
	fmt.Printf("result is %v\n", error_score)
}
