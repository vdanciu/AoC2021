package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	line_no := 0
	template := ""
	pairs := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if line_no > 0 {
			if len(line) == 0 {
				continue
			}
			components := strings.Split(line, " -> ")
			pairs[components[0]] = components[1]
		} else {
			template = line
		}
		line_no++
	}

	new_template := template
	for step := 0; step < 10; step++ {
		new_template = replace(new_template, pairs)
	}

	histogram := make(map[byte]int)
	for i := 0; i < len(new_template); i++ {
		if _, is := histogram[new_template[i]]; !is {
			histogram[new_template[i]] = 0
		}
		histogram[new_template[i]]++
	}
	min := math.MaxInt
	max := 0
	for i := 0; i < len(new_template); i++ {
		if min > histogram[new_template[i]] {
			min = histogram[new_template[i]]
		}
		if max < histogram[new_template[i]] {
			max = histogram[new_template[i]]
		}
	}

	fmt.Printf("result is %v", max-min)
}

func replace(template string, pairs map[string]string) string {
	new_template := string(template[0])
	for i := 0; i < len(template)-1; i++ {
		new_template = fmt.Sprintf("%v%v%v",
			new_template,
			pairs[template[i:i+2]],
			string(template[i+1]))
	}
	return new_template
}
