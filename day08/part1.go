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
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		output := strings.Fields(strings.Split(line, "|")[1])
		count += count_easy(output)
	}
	fmt.Printf("result is %v\n", count)
}

func count_easy(a []string) int {
	count := 0
	for i := range a {
		switch len(a[i]) {
		case 2:
			count++
		case 3:
			count++
		case 4:
			count++
		case 7:
			count++
		}
	}
	return count
}
