package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type cave_map map[string][]string
type black_list struct {
	list      []string
	revisited bool
}

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
	count := count_routes(cave_map, "start", "end", black_list{list: []string{}, revisited: false})
	fmt.Printf("result is %v\n", count)
}

func add_route(cave_map cave_map, from string, to string) {
	if _, is := cave_map[from]; !is {
		cave_map[from] = []string{}
	}
	cave_map[from] = append(cave_map[from], to)
}

func count_routes(cave_map cave_map, from string, to string, bl black_list) int {
	count := 0
	for _, via := range cave_map[from] {
		new_bl := black_list{list: append(bl.list, from), revisited: bl.revisited}
		if via == to {
			count++
		} else {
			if allow_visit(via, &new_bl) {
				count += count_routes(cave_map, via, to, new_bl)
			}
		}
	}
	return count
}

func allow_visit(cave string, bl *black_list) bool {
	if cave == "start" {
		return false
	}
	if small(cave) && cave_in(cave, *bl) {
		revisited := bl.revisited
		bl.revisited = true
		return false || !revisited
	}
	// safety check
	if len(bl.list) > 10000 {
		return false
	}
	return true
}

func cave_in(cave string, bl black_list) bool {
	in := false
	for _, listed := range bl.list {
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
