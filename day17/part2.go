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

type Speed struct {
	vx, vy int
}

type SpeedSet map[Speed]bool

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
	count := count_shots(target)
	fmt.Printf("result is: %d\n", count)
}

func atoi(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}

func count_shots(target Target) int {
	speed_set := SpeedSet{}
	for n := 1; n*(n-1)/2 < abs(target.xmax); n++ {
		vx1 := (target.xmin + sign(target.xmin)*n*(n-1)/2) / n
		vx2 := (target.xmax + sign(target.xmax)*n*(n-1)/2) / n
		for vx := vx1; vx <= vx2; vx++ {
			vy1 := (target.ymin + n*(n-1)/2) / n
			vy2 := max(abs(target.ymin), abs(target.ymax)) - 1
			for vy := vy1; vy <= vy2; vy++ {
				_, err := shoot(vx, vy, target)
				if err == nil {
					speed_set[Speed{vx, vy}] = true
				}
			}
		}
	}
	return len(speed_set)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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

func sign(a int) int {
	if a < 0 {
		return -1
	}
	return 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
