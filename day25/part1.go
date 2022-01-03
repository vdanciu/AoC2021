package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type seaMap [][]byte

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	sea_map := seaMap{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		map_line := line
		sea_map = append(sea_map, []byte(map_line))
	}
	i := 1
	for ; true; i++ {
		if !(move(sea_map)) {
			break
		}
	}
	fmt.Printf("step %v\n", i)
}

func print_map(m seaMap) {
	for _, line := range m {
		fmt.Println(string(line))
	}
	fmt.Println(string(""))
}

func move(m seaMap) bool {
	east := move_east(m)
	south := move_south(m)
	return east || south
}

func move_east(m seaMap) bool {
	moved := false
	copy := deep_copy(m)
	for i := range copy {
		for j := range copy[i] {
			if copy[i][j] != '>' {
				continue
			}
			next := 0
			if j < len(copy[i])-1 {
				next = j + 1
			}
			if copy[i][next] == '.' {
				m[i][j], m[i][next] = copy[i][next], copy[i][j]
				moved = true
			}
		}
	}
	return moved
}
func move_south(m seaMap) bool {
	moved := false
	copy := deep_copy(m)
	for i := range copy {
		next_i := 0
		if i < len(m)-1 {
			next_i = i + 1
		}
		for j := 0; j < len(copy[i]); j++ {
			if copy[i][j] != 'v' {
				continue
			}
			if copy[next_i][j] == '.' {
				m[i][j], m[next_i][j] = copy[next_i][j], copy[i][j]
				moved = true
			}
		}
	}
	return moved
}

func deep_copy(m seaMap) seaMap {
	r := make(seaMap, len(m))
	for i := range m {
		r[i] = make([]byte, len(m[i]))
		copy(r[i], m[i])
	}
	return r
}
