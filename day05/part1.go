package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type segment struct {
	p1 point
	p2 point
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	segments := []segment{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		components := strings.Split(line, " -> ")
		segments = append(
			segments,
			segment{
				newPoint(components[0]),
				newPoint(components[1])})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	traces := make(map[point]int)
	for i := 0; i < len(segments); i++ {
		segments[i].trace(traces)
	}

	count := 0
	for point := range traces {
		if traces[point] > 1 {
			count++
		}
	}
	fmt.Printf("result is %v", count)
}

func newPoint(p string) point {
	a := strings.Split(p, ",")
	x, _ := strconv.Atoi(a[0])
	y, _ := strconv.Atoi(a[1])
	return point{x, y}
}

func (s segment) trace(traces map[point]int) {
	low, high := s.bb()
	if low.x == high.x {
		for y := low.y; y <= high.y; y++ {
			point{low.x, y}.trace(traces)
		}
	}
	if low.y == high.y {
		for x := low.x; x <= high.x; x++ {
			point{x, low.y}.trace(traces)
		}
	}
	if high.x-low.x == high.y-low.y {
		for o := 0; o <= (high.x - low.x); o++ {
			x := s.p1.x + o*sign(s.p1.x, s.p2.x)
			y := s.p1.y + o*sign(s.p1.y, s.p2.y)
			point{x, y}.trace(traces)
		}
	}
}

func (s segment) bb() (low point, high point) {
	low = point{min(s.p1.x, s.p2.x), min(s.p1.y, s.p2.y)}
	high = point{max(s.p1.x, s.p2.x), max(s.p1.y, s.p2.y)}
	return
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func sign(a int, b int) int {
	if a < b {
		return 1
	}
	return -1
}

func (p point) trace(traces map[point]int) {
	_, is := traces[p]
	if !is {
		traces[p] = 0
	}
	traces[p] += 1
}
