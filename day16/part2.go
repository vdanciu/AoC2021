package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Operation int

const (
	OpSum     Operation = 0
	OpProduct           = 1
	OpMin               = 2
	OpMax               = 3
	OpVal               = 4
	OpGT                = 5
	OpLT                = 6
	OpEq                = 7
	OpNoop              = 99
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	line := ""
	for scanner.Scan() {
		line = scanner.Text()
	}
	bit_stream := ""
	for i := range line {
		bit_stream += decode_hex(line[i])
	}
	_, v, _ := process_packet(bit_stream, OpVal, math.MaxInt)

	fmt.Printf("v=%v\n", v)
}

func decode_hex(r byte) string {
	n, err := strconv.ParseInt(string(r), 16, 8)
	if err != nil {
		log.Fatal(err)
	}
	s := strconv.FormatInt(n, 2)
	if len(s) < 4 {
		s = strings.Repeat("0", 4-len(s)) + s
	}
	return s[len(s)-4:]
}

func process_packet(bs string, top_op Operation, max_reads int) (pointer, value int, err error) {
	if len(bs) < 7 {
		return len(bs), 0, errors.New("Invalid packet")
	}
	p := 0
	reads := 0

	const (
		VERSION = iota
		ID
		I
		COUNT
		SIZE
		LITERAL
		PACKETS
		END
	)
	state := VERSION

	operation := Operation(OpNoop)
	operands := []int{}

	for state != END {
		if p == len(bs) {
			break
		}
		switch state {
		case VERSION:
			if len(bs)-p < 7 || reads == max_reads {
				state = END
				break
			}
			p += 3
			reads++
			state = ID
		case ID:
			id, _ := strconv.ParseInt(bs[p:p+3], 2, 8)
			operation = Operation(id)
			p += 3
			if id == 4 {
				state = LITERAL
			} else {
				state = I
			}
		case LITERAL:
			literal := ""
			for bs[p] == '1' {
				p += 1
				literal += bs[p : p+4]
				p += 4
			}
			p += 1
			literal += bs[p : p+4]
			p += 4
			number, _ := strconv.ParseInt(literal, 2, 64)
			operands = append(operands, int(number))
			state = VERSION
		case I:
			i, _ := strconv.ParseInt(bs[p:p+1], 2, 8)
			p += 1
			if i == 0 {
				state = SIZE
			} else {
				state = COUNT
			}
		case SIZE:
			size, _ := strconv.ParseInt(bs[p:p+15], 2, 16)
			p += 15
			offset, value, err := process_packet(bs[p:p+int(size)], operation, math.MaxInt)
			p += offset
			if err == nil {
				operands = append(operands, value)
			}
			state = VERSION
		case COUNT:
			count, _ := strconv.ParseInt(bs[p:p+11], 2, 16)
			p += 11
			offset, value, err := process_packet(bs[p:], operation, int(count))
			p += offset
			if err == nil {
				operands = append(operands, value)
			}
			state = VERSION
		case PACKETS:
			state = END
		}
	}

	switch top_op {
	case OpVal:
		value = operands[0]
	case OpSum:
		for _, o := range operands {
			value += o
		}
	case OpProduct:
		value = 1
		for _, o := range operands {
			value *= o
		}
	case OpMin:
		value = math.MaxInt
		for _, o := range operands {
			if value > o {
				value = o
			}
		}
	case OpMax:
		value = math.MinInt
		for _, o := range operands {
			if value < o {
				value = o
			}
		}
	case OpLT:
		if operands[0] < operands[1] {
			value = 1
		} else {
			value = 0
		}
	case OpGT:
		if operands[0] > operands[1] {
			value = 1
		} else {
			value = 0
		}
	case OpEq:
		if operands[0] == operands[1] {
			value = 1
		} else {
			value = 0
		}
	}

	return p, value, nil
}
