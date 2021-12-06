package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	prev_depth := -1
	window_buff := []int{0, 0, 0}
	count := 0
	for i := 0; scanner.Scan(); i++ {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		window_buff[0] += depth
		window_buff[1] += depth
		window_buff[2] += depth

		buf_idx := (i + 1) % 3

		if i >= 2 {
			window_depth := window_buff[buf_idx]
			if prev_depth != -1 && prev_depth < window_depth {
				count++
			}
			prev_depth = window_depth
		}
		window_buff[buf_idx] = 0 // reset the window
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", count)
}
