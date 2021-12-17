package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Target struct {
	xmin, xmax, ymin, ymax int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	var target Target
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`-?\d+`)
		limits := re.FindAllString(line, -1)
		target.xmin = atoi(limits[0])
		target.xmax = atoi(limits[1])
		target.ymin = atoi(limits[2])
		target.ymax = atoi(limits[3])
	}
	highest := find_shot(target)
	fmt.Printf("result is: %d\n", highest)
}

func find_shot(target Target) int {

	n := (target.xmax) / 2
	for vx := ((target.xmin+target.xmax)/2 + n*(n-1)/2) / n; n >= vx; n-- {
		vx = ((target.xmin+target.xmax)/2 + n*(n-1)/2) / n
	}

	vx := ((target.xmin+target.xmax)/2 + n*(n-1)/2) / n
	vy := max(abs(target.ymin), abs(target.ymax)) - 1

	maxh, err := shoot(vx, vy, target)
	if err != nil {
		log.Fatal(err)
	}

	return maxh
}

func shoot(vx, vy int, target Target) (int, error) {
	x := 0
	y := 0
	var err error

	ymax := math.MinInt

	for {
		x += vx
		y += vy
		if vx > 0 {
			vx--
		}
		if vx < 0 {
			vx++
		}
		vy--

		if vy == 0 {
			ymax = y
		}

		if x >= target.xmin && x <= target.xmax &&
			y >= target.ymin && y <= target.ymax {
			break
		}

		if vy < 0 && y < target.ymin {
			err = errors.New("Overshoot")
			break
		}
	}
	return ymax, err
}

func atoi(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
