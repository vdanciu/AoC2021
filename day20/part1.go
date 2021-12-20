package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type Coord struct {
	row, col int
}

type Bitmap map[Coord]bool

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	algo := ""
	bitmap := Bitmap{}
	row := 0
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		if len(algo) > 0 {
			for col := 0; col < len(text); col++ {
				bitmap[Coord{row, col}] = text[col] == '#'
			}
			row++
		} else {
			algo = text
		}
	}
	_, bmax := bitmap.bounds()
	bitmap[Coord{bmax.row + 1, bmax.col + 1}] = false

	for step := 0; step < 50; step++ {
		bitmap = enhance(bitmap, algo)
	}
	count := 0
	for _, v := range bitmap {
		if v {
			count++
		}
	}
	fmt.Printf("%v", count)
}

func enhance(bitmap Bitmap, algo string) Bitmap {
	result := Bitmap{}
	bmin, bmax := bitmap.bounds()
	beyonder := bitmap[Coord{bmax.row, bmax.col}]

	for row := bmin.row - 9; row < bmax.row+9; row++ {
		for col := bmin.col - 9; col < bmax.col+9; col++ {
			coord := Coord{row, col}
			result[coord] = enhance_pixel(bitmap, coord, algo, beyonder)
		}
	}

	return result
}

func enhance_pixel(bitmap Bitmap, coord Coord, algo string, beyonder bool) bool {
	scan := []Coord{
		Coord{-1, -1},
		Coord{-1, 0},
		Coord{-1, 1},
		Coord{0, -1},
		Coord{0, 0},
		Coord{0, 1},
		Coord{1, -1},
		Coord{1, 0},
		Coord{1, 1},
	}
	binary := ""
	for _, offset := range scan {
		boolToStr := map[bool]string{true: "1", false: "0"}
		value := boolToStr[read_coord(bitmap, Coord{coord.row + offset.row, coord.col + offset.col}, beyonder)]
		binary += value
	}

	index, _ := strconv.ParseInt(binary, 2, 16)
	pixel := algo[int(index)]
	byteToBool := map[byte]bool{'#': true, '.': false}
	return byteToBool[pixel]
}

func read_coord(bitmap Bitmap, coord Coord, beyonder bool) bool {
	if value, in := bitmap[coord]; in {
		return value
	}
	return beyonder
}

func (b Bitmap) bounds() (Coord, Coord) {
	max := Coord{math.MinInt, math.MinInt}
	min := Coord{math.MaxInt, math.MaxInt}
	for c := range b {
		if c.row > max.row {
			max.row = c.row
		}
		if c.row < min.row {
			min.row = c.row
		}
		if c.col > max.col {
			max.col = c.col
		}
		if c.col < min.col {
			min.col = c.col
		}
	}
	return min, max
}
