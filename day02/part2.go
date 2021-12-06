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
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	horizontal := 0
	depth := 0
	aim := 0
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		value, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		switch line[0] {
		case "forward":
			horizontal += value
			depth += aim * value
		case "up":
			aim -= value
		case "down":
			aim += value
		}
	}

	position := horizontal * depth

	fmt.Printf("result = %v", position)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
