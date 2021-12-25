package main

import (
	"bufio"
	"errors"
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
	state   bool
	deleted bool
	exclude Cuboids
}

type Cuboid struct {
	CuboidCoord
	CuboidInfo
}

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
		if len(text) == 0 {
			continue
		}
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
	}
	count := cuboids.volume()
	fmt.Printf("count=%v\n", count)
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("fail to convert")
	}
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

func (c *CuboidCoord) volume() int64 {
	return abs(1+c.xmax-c.xmin) * abs(1+c.ymax-c.ymin) * abs(1+c.zmax-c.zmin)
}

func abs(a int) int64 {
	if a < 0 {
		return int64(-a)
	}
	return int64(a)
}

func (c1 *Cuboid) intersection(c2 *Cuboid) (Cuboid, error) {
	intersection := Cuboid{
		CuboidCoord{
			xmin: max(c1.xmin, c2.xmin),
			xmax: min(c1.xmax, c2.xmax),
			ymin: max(c1.ymin, c2.ymin),
			ymax: min(c1.ymax, c2.ymax),
			zmin: max(c1.zmin, c2.zmin),
			zmax: min(c1.zmax, c2.zmax),
		},
		CuboidInfo{state: true},
	}
	if intersection.valid() {
		return intersection, nil
	}
	return Cuboid{}, errors.New("No intersect")
}

func (c Cuboid) valid() bool {
	return c.xmin <= c.xmax &&
		c.ymin <= c.ymax &&
		c.zmin <= c.zmax
}

func (list Cuboids) String() (s string) {
	for _, c := range list {
		s += fmt.Sprintf("%v\n", *c)
	}
	return
}

func (cuboids Cuboids) volume() int64 {
	result := int64(0)
	for c1 := range cuboids {
		for c2 := c1 - 1; c2 >= 0; c2-- {
			if cuboids[c2].deleted {
				continue
			}
			intersection, err := cuboids[c1].intersection(cuboids[c2])
			if err != nil {
				continue //nothing to do if they don't intersect
			}
			if cuboids[c2].CuboidCoord == intersection.CuboidCoord {
				cuboids[c2].deleted = true
				cuboids[c2].exclude = Cuboids{}
			} else {
				cuboids[c2].exclude = append(cuboids[c2].exclude, &intersection)
			}
		}
		cuboids[c1].deleted = !cuboids[c1].state
	}
	for _, cuboid := range cuboids {
		volume := int64(0)
		if !cuboid.deleted {
			volume = cuboid.volume()
		}
		exclude := cuboid.exclude.volume()
		if !cuboid.deleted && volume < exclude {
			panic("this is not possible")
		}
		result += (volume - exclude)
	}
	return result
}
