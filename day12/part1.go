package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type cave_map map[string][]string
type black_list []string

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	cave_map := make(cave_map)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		route := strings.Split(line, "-")
		add_route(cave_map, route[0], route[1])
		add_route(cave_map, route[1], route[0])
	}
	count := count_routes(cave_map, "start", "end", black_list{})
	fmt.Printf("%v\n", count)
}

func add_route(cave_map cave_map, from string, to string) {
	if _, is := cave_map[from]; !is {
		cave_map[from] = []string{}
	}
	cave_map[from] = append(cave_map[from], to)
}

func count_routes(cave_map cave_map, from string, to string, bl black_list) int {
	count := 0
	bl = append(bl, from)
	for _, via := range cave_map[from] {
		if via == to {
			count++
		} else {
			if allow_visit(via, bl) {
				count += count_routes(cave_map, via, to, bl)
			}
		}
	}
	return count
}

func allow_visit(cave string, bl black_list) bool {
	if small(cave) && cave_in(cave, bl) {
		return false
	}
	// safety check
	if len(bl) > 10000 {
		return false
	}
	return true
}

func cave_in(cave string, bl black_list) bool {
	in := false
	for _, listed := range bl {
		if cave == listed {
			in = true
			break
		}
	}
	return in
}

func small(cave string) bool {
	return strings.ToLower(cave) == cave
}
