package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	NUMBERS_DRAW = iota
	NEW_BOARD    = iota
	READ_BOARD   = iota
)

type board [5][5]string

type winboard struct {
	board  board
	number int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	state := NUMBERS_DRAW
	numbers := []string{}
	boards := []board{}
	var acc_board board
	acc_index := 0
	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case NUMBERS_DRAW:
			numbers = strings.Split(line, ",")
			state = NEW_BOARD
		case NEW_BOARD:
			acc_index = 0
			state = READ_BOARD
		case READ_BOARD:
			board_line := strings.Fields(line)
			if len(board_line) == 5 {
				copy(acc_board[acc_index][:], board_line)
			}
			acc_index++
			if acc_index == 5 {
				boards = append(boards, acc_board)
				state = NEW_BOARD
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rows := make([][5]int, len(boards))
	cols := make([][5]int, len(boards))
	winning_boards := make(map[int]winboard)
	var winning_order []int
	winning_number := -1

	for n := range numbers {
		for b := range boards {
			for row := 0; row < 5; row++ {
				for col := 0; col < 5; col++ {
					if _, is := winning_boards[b]; is {
						continue
					}
					if numbers[n] == boards[b][row][col] {
						boards[b][row][col] = "*" + boards[b][row][col]
						rows[b][row] += 1
						cols[b][col] += 1
					}
					if rows[b][row] == 5 || cols[b][col] == 5 {
						winning_number, err = strconv.Atoi(numbers[n])
						var win winboard
						win.board = boards[b]
						win.number = winning_number
						winning_boards[b] = win
						winning_order = append(winning_order, b)
					}
				}
			}
		}
	}

	winning_board := winning_boards[winning_order[len(winning_order)-1]]
	sum := 0
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			if winning_board.board[row][col][0] != '*' {
				v, err := strconv.Atoi(winning_board.board[row][col])
				if err != nil {
					log.Fatal(err)
				}
				sum += v
			}
		}
	}

	fmt.Printf("result is %v\n", sum*winning_board.number)
}
