package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type Sfnumber struct {
	number              int
	left, right, parent *Sfnumber
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	list := []*Sfnumber{}
	for scanner.Scan() {
		text := scanner.Text()
		sfnumber := parseSfnumber(text)
		list = append(list, sfnumber)
	}
	max := math.MinInt
	for _, n1 := range list {
		for _, n2 := range list {
			if n1 == n2 {
				continue
			}
			mag := magnitude(add(n1, n2))
			if max < mag {
				max = mag
			}
		}
	}
	fmt.Printf("magnitude: %v", max)
}

func parseSfnumber(text string) *Sfnumber {
	var sfnumber Sfnumber
	if text[0] != '[' {
		number, _ := strconv.Atoi(text)
		sfnumber.number = number
	} else {
		count := 0
		s := 1
		for i, r := range text[1 : len(text)-1] {
			if r == '[' {
				count++
			}
			if r == ']' {
				count--
			}
			if r == ',' && count == 0 {
				sfnumber.left = parseSfnumber(text[s : i+1])
				s = i + 2
			}
		}
		sfnumber.right = parseSfnumber(text[s : len(text)-1])
		sfnumber.left.parent = &sfnumber
		sfnumber.right.parent = &sfnumber
	}

	return &sfnumber
}

func inspect(sfnumber Sfnumber) string {
	if sfnumber.left != nil {
		return fmt.Sprintf("[%v,%v]", inspect(*sfnumber.left), inspect(*sfnumber.right))
	} else {
		return fmt.Sprintf("%v", sfnumber.number)
	}
}

func explode(sfnumber *Sfnumber) bool {
	explodable, found := drill_leftmost(sfnumber, 5)
	if found {
		left, found_left := find_left_neighbor(explodable)
		if found_left {
			left.number += explodable.left.number
		}

		right, found_right := find_right_neighbor(explodable)
		if found_right {
			right.number += explodable.right.number
		}

		explodable.left = nil
		explodable.right = nil
		explodable.number = 0
	}
	return found
}

func split(sfnumber *Sfnumber) bool {
	if sfnumber.left == nil {
		if sfnumber.number >= 10 {
			sfnumber.left = &Sfnumber{number: int(math.Floor(float64(sfnumber.number) / 2))}
			sfnumber.right = &Sfnumber{number: int(math.Ceil(float64(sfnumber.number) / 2))}
			sfnumber.left.parent = sfnumber
			sfnumber.right.parent = sfnumber
			return true
		}
		return false
	}
	if split(sfnumber.left) {
		return true
	} else {
		return split(sfnumber.right)
	}
}

func drill_leftmost(sfn *Sfnumber, depth int) (*Sfnumber, bool) {
	if sfn.left != nil {
		if depth <= 1 {
			if sfn.left.left == nil && sfn.right.left == nil {
				return sfn, true
			}
		}
		left, found := drill_leftmost(sfn.left, depth-1)
		if found {
			return left, found
		} else {
			return drill_leftmost(sfn.right, depth-1)
		}
	} else {
		return sfn, false
	}
}

func find_left_neighbor(sfn *Sfnumber) (*Sfnumber, bool) {
	if sfn.parent == nil {
		return nil, false
	}
	if sfn.parent.left == sfn {
		return find_left_neighbor(sfn.parent)
	}
	return find_rightmost_value(sfn.parent.left)
}

func find_right_neighbor(sfn *Sfnumber) (*Sfnumber, bool) {
	if sfn.parent == nil {
		return nil, false
	}
	if sfn.parent.right == sfn {
		return find_right_neighbor(sfn.parent)
	}
	return find_leftmost_value(sfn.parent.right)
}

func find_rightmost_value(sfn *Sfnumber) (*Sfnumber, bool) {
	if sfn.left == nil {
		return sfn, true
	}
	return find_rightmost_value(sfn.right)
}

func find_leftmost_value(sfn *Sfnumber) (*Sfnumber, bool) {
	if sfn.left == nil {
		return sfn, true
	}
	return find_leftmost_value(sfn.left)
}

func add(a, b *Sfnumber) *Sfnumber {
	if a == nil {
		return reduce(b)
	}
	var r Sfnumber
	r.left = a
	a.parent = &r
	r.right = b
	b.parent = &r
	return reduce(&r)
}

func reduce(sfnumber *Sfnumber) *Sfnumber {
	clone := dup(sfnumber)
	for {
		if explode(clone) {
			continue
		}
		if !split(clone) {
			break
		}
	}
	return clone
}

func magnitude(sfnumber *Sfnumber) int {
	if sfnumber.left == nil {
		return sfnumber.number
	}
	return 3*magnitude(sfnumber.left) + 2*magnitude(sfnumber.right)
}

func dup(src *Sfnumber) *Sfnumber {
	if src == nil {
		return nil
	}
	var sfnumber Sfnumber
	sfnumber = *src
	sfnumber.left = dup(src.left)
	sfnumber.right = dup(src.right)
	if sfnumber.left != nil {
		sfnumber.left.parent = &sfnumber
		sfnumber.right.parent = &sfnumber
	}

	return &sfnumber
}
