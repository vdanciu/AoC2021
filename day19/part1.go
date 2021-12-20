package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MappingCoordsEnum int

const (
	MAP_NIL MappingCoordsEnum = iota
	MAP_XYZ
	MAP_XZY
	MAP_ZYX
	MAP_YXZ
	MAP_ZXY
	MAP_YZX
)

type TranslateFun func(Coordinate) Coordinate

type Scanner struct {
	id         int
	coords     Coordinates
	absDiffs   AbsDiffs
	transforms Transforms
	transref   []TranslateFun
}

type Transforms map[int]TranslateFun

type CoordPair struct {
	first, second *Coordinate
}

type AbsDiffs map[Coordinate]CoordPair

type Coordinates []Coordinate

type Coordinate struct {
	x, y, z int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	scanners := []Scanner{}
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		if text[0:3] == "---" {
			scanners = append(scanners, Scanner{id: len(scanners), transref: []TranslateFun{}})
			continue
		}
		components := strings.Split(text, ",")
		scanners[len(scanners)-1].coords = append(scanners[len(scanners)-1].coords,
			Coordinate{
				atoi(components[0]),
				atoi(components[1]),
				atoi(components[2])})
	}
	for i := range scanners {
		for j := i + 1; j < len(scanners); j++ {
			if compare(&scanners[i], &scanners[j]) {
				compare(&scanners[j], &scanners[i])
			}
		}

	}
	beacons := map[Coordinate]bool{}
	build_transref(&scanners, &scanners[0])
	transform_beacons(beacons, &scanners, &scanners[0])
	fmt.Printf("beacons(%v):\n", len(beacons))
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func compare(s1, s2 *Scanner) bool {
	build_abs_diff(s1)
	build_abs_diff(s2)
	for diff1, pair1 := range s1.absDiffs {
		if _, in := s2.absDiffs[diff1]; in {
			t1, t2 := map_coords(pair1, s2.absDiffs[diff1])
			count_total := 0
			for _, t := range []TranslateFun{t1, t2} {
				if t == nil {
					continue
				}
				transformed := apply(t, s2.coords)
				count := count_coords(s1.coords, transformed)
				count_total += count
				if count >= 12 {
					if s2.transforms == nil {
						s2.transforms = make(map[int]TranslateFun)
					}
					s2.transforms[s1.id] = t
					return true
				}
			}
			if count_total == 0 {
				_, _ = map_coords(pair1, s2.absDiffs[diff1])
				panic(fmt.Sprintf("it does count! %v, %v\n", s1.id, s2.id))
			}
		}
	}
	return false
}

func build_abs_diff(s *Scanner) {
	if len(s.absDiffs) == 0 {
		s.absDiffs = make(AbsDiffs)
		for i := range s.coords {
			for j := i + 1; j < len(s.coords); j++ {
				diff := []int{
					abs(s.coords[i].x - s.coords[j].x),
					abs(s.coords[i].y - s.coords[j].y),
					abs(s.coords[i].z - s.coords[j].z),
				}
				sort.Ints(diff)
				s.absDiffs[Coordinate{diff[0], diff[1], diff[2]}] = CoordPair{&s.coords[i], &s.coords[j]}
			}
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sign(a int) int {
	if a < 0 {
		return -1
	}
	return 1
}

func map_coords(p1, p2 CoordPair) (TranslateFun, TranslateFun) {
	mapping := MAP_NIL
	if abs(p1.first.x-p1.second.x) == abs(p2.first.x-p2.second.x) {
		if abs(p1.first.y-p1.second.y) == abs(p2.first.y-p2.second.y) {
			mapping = MAP_XYZ
		} else {
			mapping = MAP_XZY
		}
	} else {
		if abs(p1.first.y-p1.second.y) == abs(p2.first.y-p2.second.y) {
			if abs(p1.first.x-p1.second.x) == abs(p2.first.z-p2.second.z) {
				mapping = MAP_ZYX
			} else {
				panic("oh no")
			}
		} else {
			if abs(p1.first.z-p1.second.z) == abs(p2.first.z-p2.second.z) {
				if abs(p1.first.x-p1.second.x) == abs(p2.first.y-p2.second.y) {
					mapping = MAP_YXZ
				} else {
					panic("oh no")
				}
			} else {
				if abs(p1.first.z-p1.second.z) == abs(p2.first.x-p2.second.x) &&
					abs(p1.first.x-p1.second.x) == abs(p2.first.y-p2.second.y) {
					mapping = MAP_ZXY
				} else {
					if abs(p1.first.y-p1.second.y) == abs(p2.first.x-p2.second.x) &&
						abs(p1.first.z-p1.second.z) == abs(p2.first.y-p2.second.y) {
						mapping = MAP_YZX
					} else {
						panic(fmt.Sprintf("oh no: %v <> %v", p1, p2))
					}
				}
			}
		}
	}
	if mapping == MAP_NIL {
		return nil, nil
	}
	return transform(mapping, p1, p2), transform(mapping, p1, reverse(p2))
}

func transform(mapping MappingCoordsEnum, p1, p2 CoordPair) TranslateFun {
	var signx, signy, signz, transx, transy, transz int
	var t TranslateFun

	switch mapping {
	case MAP_XYZ:
		signx = sign(p1.first.x-p1.second.x) * sign(p2.first.x-p2.second.x)
		signy = sign(p1.first.y-p1.second.y) * sign(p2.first.y-p2.second.y)
		signz = sign(p1.first.z-p1.second.z) * sign(p2.first.z-p2.second.z)
		transx = p1.first.x - signx*p2.first.x
		transy = p1.first.y - signy*p2.first.y
		transz = p1.first.z - signz*p2.first.z

		t = func(src Coordinate) Coordinate {
			return Coordinate{
				src.x*signx + transx,
				src.y*signy + transy,
				src.z*signz + transz,
			}
		}
	case MAP_XZY:
		signx = sign(p1.first.x-p1.second.x) * sign(p2.first.x-p2.second.x)
		signy = sign(p1.first.z-p1.second.z) * sign(p2.first.y-p2.second.y)
		signz = sign(p1.first.y-p1.second.y) * sign(p2.first.z-p2.second.z)
		transx = p1.first.x - signx*p2.first.x
		transy = p1.first.y - signz*p2.first.z
		transz = p1.first.z - signy*p2.first.y

		t = func(src Coordinate) Coordinate {
			return Coordinate{
				src.x*signx + transx,
				src.z*signz + transy,
				src.y*signy + transz,
			}
		}
	case MAP_ZYX:
		signx = sign(p1.first.z-p1.second.z) * sign(p2.first.x-p2.second.x)
		signy = sign(p1.first.y-p1.second.y) * sign(p2.first.y-p2.second.y)
		signz = sign(p1.first.x-p1.second.x) * sign(p2.first.z-p2.second.z)
		transx = p1.first.x - signz*p2.first.z
		transy = p1.first.y - signy*p2.first.y
		transz = p1.first.z - signx*p2.first.x

		t = func(src Coordinate) Coordinate {
			return Coordinate{
				src.z*signz + transx,
				src.y*signy + transy,
				src.x*signx + transz,
			}
		}
	case MAP_YXZ:
		signx = sign(p1.first.y-p1.second.y) * sign(p2.first.x-p2.second.x)
		signy = sign(p1.first.x-p1.second.x) * sign(p2.first.y-p2.second.y)
		signz = sign(p1.first.z-p1.second.z) * sign(p2.first.z-p2.second.z)
		transx = p1.first.x - signy*p2.first.y
		transy = p1.first.y - signx*p2.first.x
		transz = p1.first.z - signz*p2.first.z

		t = func(src Coordinate) Coordinate {
			return Coordinate{
				src.y*signy + transx,
				src.x*signx + transy,
				src.z*signz + transz,
			}
		}
	case MAP_ZXY:
		signx = sign(p1.first.z-p1.second.z) * sign(p2.first.x-p2.second.x)
		signy = sign(p1.first.x-p1.second.x) * sign(p2.first.y-p2.second.y)
		signz = sign(p1.first.y-p1.second.y) * sign(p2.first.z-p2.second.z)

		transx = p1.first.x - signy*p2.first.y
		transy = p1.first.y - signz*p2.first.z
		transz = p1.first.z - signx*p2.first.x

		t = func(src Coordinate) Coordinate {
			return Coordinate{
				src.y*signy + transx,
				src.z*signz + transy,
				src.x*signx + transz,
			}
		}
	case MAP_YZX:
		signx = sign(p1.first.y-p1.second.y) * sign(p2.first.x-p2.second.x)
		signy = sign(p1.first.z-p1.second.z) * sign(p2.first.y-p2.second.y)
		signz = sign(p1.first.x-p1.second.x) * sign(p2.first.z-p2.second.z)
		transx = p1.first.x - signz*p2.first.z
		transy = p1.first.y - signx*p2.first.x
		transz = p1.first.z - signy*p2.first.y

		t = func(src Coordinate) Coordinate {
			return Coordinate{
				src.z*signz + transx,
				src.x*signx + transy,
				src.y*signy + transz,
			}
		}
	}
	if t != nil {
		// fmt.Printf("1. trans(%v) %v to %v=%v\n", mapping, *p2.first, t(*p2.first), *p1.first)
		// fmt.Printf("2. trans(%v) %v to %v=%v\n", mapping, *p2.second, t(*p2.second), *p1.second)

		// if true {
		// 	fmt.Printf("sx=%v sy=%v sz=%v\n", signx, signy, signz)
		// 	fmt.Printf("tx=%v ty=%v tz=%v\n", transx, transy, transz)
		// 	fmt.Printf("1. trans(%v) %v to %v=%v\n", mapping, *p2.first, t(*p2.first), *p1.first)
		// 	fmt.Printf("2. trans(%v) %v to %v=%v\n", mapping, *p2.second, t(*p2.second), *p1.second)
		//panic("incorrect transform")
		//}
		if *p1.first != t(*p2.first) || *p1.second != t(*p2.second) {
			t = nil
		}
	}
	return t
}

func apply(t TranslateFun, src []Coordinate) map[Coordinate]bool {
	result := map[Coordinate]bool{}
	for _, c := range src {
		result[t(c)] = true
	}
	return result
}

func count_coords(what []Coordinate, where map[Coordinate]bool) int {
	count := 0
	for _, c := range what {
		if _, in := where[c]; in {
			count++
		}
	}
	return count
}

func reverse(p CoordPair) (c CoordPair) {
	c.first, c.second = p.second, p.first
	return
}

func (transforms Transforms) String() string {
	s := ""
	for k, t := range transforms {
		s += fmt.Sprintf("%v -> %v, ", k, t)
	}
	return s
}

func (p CoordPair) String() string {
	return fmt.Sprintf("%v - %v [%v]", *p.first, *p.second, Coordinate{p.first.x - p.second.x, p.first.y - p.second.y, p.first.z - p.second.z})
}

func transform_beacons(beacons map[Coordinate]bool, scanners *[]Scanner, ref *Scanner) {
	for _, scanner := range *scanners {
		transform_scanner(beacons, &scanner, ref)
	}
}

func transform_scanner(beacons map[Coordinate]bool, scanner *Scanner, ref *Scanner) {
	for _, c := range scanner.coords {
		for _, t := range scanner.transref {
			c = t(c)
		}
		beacons[c] = true
	}
}

func build_transref(scanners *[]Scanner, ref *Scanner) {
	for i := range *scanners {
		if (*scanners)[i].id != ref.id {
			(*scanners)[i].transref, _ = find_transref(&(*scanners)[i], ref, scanners, map[int]bool{})
		}
	}
}

func find_transref(scanner, ref *Scanner, scanners *[]Scanner, visited map[int]bool) ([]TranslateFun, bool) {
	visited[scanner.id] = true
	//fmt.Printf("processing for %v in %v [%v]\n", scanner.id, scanner.transforms, visited)
	if t, in := scanner.transforms[ref.id]; in {
		//fmt.Printf("processing found\n")
		return []TranslateFun{t}, true
	} else {
		for s, t := range scanner.transforms {
			if _, in := visited[s]; !in {
				partial, found := find_transref(&(*scanners)[s], ref, scanners, visited)
				if found {
					partial = append([]TranslateFun{t}, partial...)
					return partial, true
				}
			}
		}
	}
	return []TranslateFun{}, false
}
