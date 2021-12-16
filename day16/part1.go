package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	//line = "A0016C880162017C3686B18A3D4780"
	bit_stream := ""
	for i := range line {
		bit_stream += decode_hex(line[i])
	}
	acc := 0
	process_packet(bit_stream, &acc)
	fmt.Printf("version_sum=%v\n", acc)

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

func process_packet(bs string, acc *int) int {
	if len(bs) < 7 {
		return len(bs)
	}
	p := 0

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

	for state != END {
		if p == len(bs) {
			break
		}
		switch state {
		case VERSION:
			if len(bs)-p < 7 {
				state = END
				break
			}
			version, _ := strconv.ParseInt(bs[p:p+3], 2, 8)
			*acc += int(version)
			p += 3
			state = ID
		case ID:
			id, _ := strconv.ParseInt(bs[p:p+3], 2, 8)
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
			p += process_packet(bs[p:p+int(size)], acc)
			state = VERSION
		case COUNT:
			count, _ := strconv.ParseInt(bs[p:p+11], 2, 16)
			p += 11
			for i := int64(0); i < count; i++ {
				p += process_packet(bs[p:], acc)
			}
			state = END
		case PACKETS:
			state = END
		}
	}

	return p
}
