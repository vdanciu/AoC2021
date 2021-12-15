package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type cell struct {
	x, y int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	var risk_map [][]int
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		row := make([]int, len(line))
		for i := range line {
			row[i], _ = strconv.Atoi(line[i])
		}
		risk_map = append(risk_map, row)
	}
	risk_map = multiply_map(risk_map)
	risk_map[0][0] = 0
	computed := make([][]int, len(risk_map))
	for i := range computed {
		computed[i] = make([]int, len(risk_map[0]))
	}
	compute_risk(risk_map, computed, 0, 0, map[cell]bool{})
	fmt.Printf("%v\n", computed[len(computed)-1][len(computed[0])-1])
}

func multiply_map(m [][]int) [][]int {
	m25 := make([][]int, len(m)*5)
	for y := range m {
		for oy := 0; oy < 5; oy++ {
			line := make([]int, len(m[y])*5)
			for x := range m[y] {
				for ox := 0; ox < 5; ox++ {
					value := m[y][x] + ox + oy
					if value > 9 {
						value -= 9
					}
					line[ox*len(m[y])+x] = value
				}
			}
			m25[oy*len(m)+y] = line
		}
	}
	return m25
}

func compute_risk(risk_map [][]int, computed [][]int, x, y int, visited map[cell]bool) {
	queue := []cell{cell{x, y}}
	for len(queue) != 0 {
		p := 0
		for i := 0; i < len(queue); i++ {
			if computed[queue[i].y][queue[i].x] < computed[queue[p].y][queue[p].x] {
				p = i
			}
		}
		c := queue[p]
		queue = append(queue[:p], queue[p+1:]...)
		if _, in := visited[c]; in {
			continue
		}
		neighbors := []cell{cell{x, y}}
		if c.x < len(risk_map[0])-1 {
			neighbors = append(neighbors, cell{c.x + 1, c.y})
		}
		if c.y < len(risk_map)-1 {
			neighbors = append(neighbors, cell{c.x, c.y + 1})
		}
		if c.x > 0 {
			neighbors = append(neighbors, cell{c.x - 1, c.y})
		}
		if c.y > 0 {
			neighbors = append(neighbors, cell{c.x, c.y - 1})
		}
		for _, n := range neighbors {
			if _, in := visited[n]; !in {
				computed[n.y][n.x] = minifnz(computed[n.y][n.x], computed[c.y][c.x]+risk_map[n.y][n.x])
				queue = append(queue, n)
			}
		}
		visited[c] = true
	}
}

func minifnz(a, b int) int {
	if a == 0 || b < a {
		return b
	}
	return a
}
