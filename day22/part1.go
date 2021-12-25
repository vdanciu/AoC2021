package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type CuboidCoord struct {
	xmin, xmax,
	ymin, ymax,
	zmin, zmax int
}

type CuboidInfo struct {
	state           bool
	intersection_of []*CuboidInfo
	intersects_with []*CuboidInfo
}

type Cuboid struct {
	CuboidCoord
	CuboidInfo
}

type Coord struct{ x, y, z int }

type CuboidsMap map[CuboidCoord]*CuboidInfo
type Cuboids []*Cuboid

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	cuboids := Cuboids{}
	for scanner.Scan() {
		text := scanner.Text()
		re := regexp.MustCompile(`(on|off)|(\-?\d+)`)
		c := re.FindAllString(text, -1)
		bool_map := map[string]bool{"on": true, "off": false}
		cuboid := Cuboid{
			CuboidCoord{
				xmin: atoi(c[1]), xmax: atoi(c[2]),
				ymin: atoi(c[3]), ymax: atoi(c[4]),
				zmin: atoi(c[5]), zmax: atoi(c[6]),
			},
			CuboidInfo{state: bool_map[c[0]]},
		}
		cuboids = append(cuboids, &cuboid)
		fmt.Printf("cuboid=%v:%v\n", cuboid, cuboid)
	}
	all := map[Coord]bool{}
	for _, c := range cuboids {
		for x := max(-50, c.xmin); x <= min(50, c.xmax); x++ {
			for y := max(-50, c.ymin); y <= min(50, c.ymax); y++ {
				for z := max(-50, c.zmin); z <= min(50, c.zmax); z++ {
					all[Coord{x, y, z}] = c.state
					if !c.state {
						delete(all, Coord{x, y, z})
					}
				}
			}
		}
	}
	fmt.Printf("result is %v", len(all))
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
