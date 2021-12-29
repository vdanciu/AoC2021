package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type position struct {
	row, col int8
}

type state struct {
	positions [8]position
	status    [8]int8
}

const (
	STATUS_INIT int8 = iota
	STATUS_STARTED
	STATUS_PAUSED
	STATUS_RESUMED
	STATUS_HOME
)

type move struct {
	state
	cost int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	state := state{}
	row := int8(1)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		if text[3] != '.' && text[3] != '#' {
			state.add_initial_position(text[3], position{row, 2})
			state.add_initial_position(text[5], position{row, 4})
			state.add_initial_position(text[7], position{row, 6})
			state.add_initial_position(text[9], position{row, 8})
			row++
		}
	}
	for i := range state.status {
		state.status[i] = state.check_home(i)
	}
	fmt.Printf("%v\n", state)
	v := cheapest_shuffle(state)
	fmt.Printf("%v", v)
}

func (state state) String() string {
	s := "\n###########"
	for row := int8(0); row < 3; row++ {
		s += "\n"
		for col := int8(0); col < 11; col++ {
			if a, in := state.at_position(row, col); in {
				s += string(a)
			} else {
				if row > 0 && (col != 2 && col != 4 && col != 6 && col != 8) {
					s += "#"
				} else {
					s += "."
				}
			}
		}
	}
	return s + "\n" + fmt.Sprintf("%v\n", state.status)
}

func (state *state) add_initial_position(a byte, position position) {
	position_map := map[byte][2]int{
		'A': {0, 1},
		'B': {2, 3},
		'C': {4, 5},
		'D': {6, 7}}
	i := position_map[a][1]
	if state.positions[position_map[a][0]].row == 0 {
		i = position_map[a][0]
	}
	state.positions[i] = position
}

func (state *state) at_position(row, col int8) (byte, bool) {
	index_to_char := [8]byte{'A', 'A', 'B', 'B', 'C', 'C', 'D', 'D'}
	what := position{row, col}
	for index, position := range state.positions {
		if position == what {
			return index_to_char[index], true
		}
	}
	return '-', false
}

func cheapest_shuffle(s state) int {
	next_moves := []move{{s, 0}}
	visited := map[state]bool{}
	cost := 0
	max_moves := 0
	max_homes := 0
	for {
		if len(next_moves) == 0 {
			break
		}
		if max_moves < len(next_moves) {
			max_moves = len(next_moves)
		}

		move := pop_cheapest_move(&next_moves)
		if _, in := visited[move.state]; in {
			continue
		}
		visited[move.state] = true
		if move.count_home() == 8 {
			cost = move.cost
			break
		}
		if move.homes() > max_homes {
			max_homes = move.homes()
			fmt.Printf("visited=%v\n", len(visited))
			fmt.Printf("%v", move.state)
		}

		for i := range s.positions {
			add_move(&next_moves, move, i, -1, 0)
			add_move(&next_moves, move, i, 1, 0)
			add_move(&next_moves, move, i, 0, -1)
			add_move(&next_moves, move, i, 0, 1)
		}
	}
	fmt.Printf("visited=%v\n", len(visited))
	fmt.Printf("next_moves=%v\n", max_moves)

	return cost
}

func pop_cheapest_move(moves *[]move) move {
	result := (*moves)[0]
	(*moves) = (*moves)[1:]
	return result
}

func add_move(moves *[]move, move move, i int, move_row, move_col int8) {
	if move.status[i] == STATUS_HOME {
		return
	}
	new_position := position{move.positions[i].row + move_row, move.positions[i].col + move_col}
	if new_position.valid() &&
		new_position.free(move.state) &&
		move.can_move(i, move_row, move_col) {
		new_move := move
		new_move.positions[i] = new_position
		costs := map[int]int{0: 1, 1: 1, 2: 10, 3: 10, 4: 100, 5: 100, 6: 1000, 7: 1000}
		new_move.cost += costs[i]
		new_move.status[i] = status_move(new_move.status[i])
		new_move.status[i] = new_move.check_home(i)
		for j := range new_move.status {
			if j != i {
				new_move.status[j] = status_stop(new_move.status[j])
			}
		}
		insert_before := 0
		for ; insert_before < len(*moves) && (*moves)[insert_before].sort_cost() < new_move.sort_cost(); insert_before++ {
		}

		*moves = append(*moves, new_move)
		copy((*moves)[insert_before+1:], (*moves)[insert_before:])
		(*moves)[insert_before] = new_move
	}
}

func (p position) valid() bool {
	return (p.row == 0 && p.col >= 0 && p.col <= 10) ||
		((p.row == 1 || p.row == 2) && (p.col == 2 || p.col == 4 || p.col == 6 || p.col == 8))
}

func (p position) free(s state) bool {
	_, b := s.at_position(p.row, p.col)
	return !b
}

func status_move(old_status int8) int8 {
	switch old_status {
	case STATUS_INIT:
		return STATUS_STARTED
	case STATUS_STARTED:
		return STATUS_STARTED
	case STATUS_PAUSED:
		return STATUS_RESUMED
	case STATUS_RESUMED:
		return STATUS_RESUMED
	case STATUS_HOME:
		panic("where do you think you're going")
	}
	panic("this really should be all")
}

func status_stop(old_status int8) int8 {
	switch old_status {
	case STATUS_INIT:
		return STATUS_INIT
	case STATUS_STARTED:
		return STATUS_PAUSED
	case STATUS_PAUSED:
		return STATUS_PAUSED
	case STATUS_RESUMED:
		panic("you should be home now")
	case STATUS_HOME:
		return STATUS_HOME
	}
	panic("this really should be all")
}

func (s *state) check_home(i int) int8 {
	status := s.status[i]
	other := s.positions[other(i)]
	if s.positions[i].row == 1 &&
		home_col(i) == s.positions[i].col &&
		other.col == s.positions[i].col && other.row == 2 {
		return STATUS_HOME
	}
	if s.positions[i].row == 2 && home_col(i) == s.positions[i].col {
		return STATUS_HOME
	}
	return status
}

func (move *move) can_move(i int, move_row, move_col int8) bool {
	// only move down home
	if move_row > 0 && home_col(i) != move.positions[i].col {
		return false
	}
	// don't move if another is forced to
	for j := range move.status {
		if j == i {
			continue
		}
		if move.status[j] == STATUS_RESUMED {
			return false
		}
		if move.positions[j].row == 0 && is_sideroom(move.positions[j].col) {
			return false
		}
	}
	// if paused don't start moving unless there is a free spot home
	if move.status[i] == STATUS_PAUSED {
		if !move.home_is_free(i) {
			return false
		}
	}
	// don't move if another is on the verge to go home
	for j := range move.status {
		if j != i {
			if move.can_go_home(j) && !move.can_go_home(i) {
				return false
			}
		}
	}
	// don't move sideways when above home
	if move_row == 0 && move.positions[i].row == 0 && home_col(i) == move.positions[i].col && move.home_is_free(i) {
		return false
	}
	return true
}

func is_sideroom(col int8) bool {
	return col == 2 || col == 4 || col == 6 || col == 8
}

func other(i int) int {
	if i%2 == 0 {
		return i + 1
	} else {
		return i - 1
	}
}

func home_col(i int) int8 {
	return map[int]int8{0: 2, 1: 2, 2: 4, 3: 4, 4: 6, 5: 6, 6: 8, 7: 8}[i]
}

func (s *state) count_home() int {
	count := 0
	for _, status := range s.status {
		if status == STATUS_HOME {
			count++
		}
	}
	return count
}

func (state *state) home_is_free(i int) bool {
	pos_up := position{row: int8(1), col: home_col(i)}
	pos_down := position{int8(2), home_col(i)}
	if !pos_up.free(*state) && !(pos_up == state.positions[i] && pos_down.free(*state)) {
		return false
	}
	if !pos_down.free(*state) && state.positions[other(i)] != pos_down {
		return false
	}
	return true
}

// func (move *move) String() string {
// 	return fmt.Sprintf("%v:%v", move.state.status, move.cost)
// }

func (m *move) sort_cost() int {
	return m.cost
}

func (m *move) homes() int {
	count_homes := 0
	for _, s := range m.status {
		count_homes += int(s / 4)
	}
	return count_homes
}

func (s *state) can_go_home(i int) bool {
	if !s.home_is_free(i) {
		return false
	}
	src := s.positions[i]
	trg := position{2, home_col(i)}

	if trg.free(*s) {
		return way_clear(s, src, trg)
	} else {
		trg = position{1, home_col(i)}
		if trg.free(*s) {
			return way_clear(s, src, trg)
		}
	}
	return false
}

func way_clear(s *state, src, trg position) bool {
	if src == trg {
		return true
	}
	var new_src position
	if src.col != trg.col {
		if src.row == 0 {
			new_src = position{0, src.col + sign(trg.col-src.col)}
		} else {
			new_src = position{src.row - 1, src.col}
		}
	} else {
		new_src = position{src.row + 1, src.col}
	}
	if new_src.free(*s) {
		return way_clear(s, new_src, trg)
	}
	return false
}

func sign(a int8) int8 {
	if a < 0 {
		return -1
	}
	return 1
}
