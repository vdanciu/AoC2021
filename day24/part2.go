package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type registers_ast [4]ast

type operand struct {
	register bool // if is register
	value    int64
}

type ast struct {
	code        string
	value       int64
	min         int64
	max         int64
	op1         *ast
	op2         *ast
	fingerprint string
}

type instruction struct {
	code string
	op1  operand
	op2  operand
}

type node struct {
	input [14]int
	value int64
	cost  float64
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	program := []instruction{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		program = append(program, make_instruction(tokens))
	}
	registers := simplify(program)

	input := [14]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	shortest_path(input, 0, &registers[2])

}

func run_ast(r *ast, input [14]int, cache *map[string]int64) int64 {
	if r.code == "OP" {
		return int64(r.value)
	}
	if r.code == "INP" {
		return int64(input[r.value-1])
	}
	if cache == nil {
		c := make(map[string]int64)
		cache = &c
	}
	if value, in := (*cache)[r.fingerprint]; in {
		return value
	}
	value := calculate(r.code, run_ast(r.op1, input, cache), run_ast(r.op2, input, cache))
	(*cache)[r.fingerprint] = value
	return value
}

func make_ast0() ast {
	a := ast{"OP", 0, 0, 0, nil, nil, "0"}
	update_fingerprint(&a)
	return a
}

func simplify(program []instruction) *registers_ast {
	ast0 := make_ast0()
	registers := registers_ast{ast0, ast0, ast0, ast0}
	input_index := int64(1)
	for _, ins := range program {
		v1 := ins.op1.value
		switch ins.code {
		case "inp":
			registers[v1] = ast{"INP", input_index, 1, 9, nil, nil, fmt.Sprintf("I%v", input_index)}
			input_index++
		default:
			registers[v1] = simplify_ast(registers[v1], &ins, &registers)
		}
	}
	return &registers
}

func simplify_ast(reg ast, ins *instruction, registers *registers_ast) ast {
	a := ast{ins.code, 0, 0, 0, &reg, ins.op2.eval_ast(registers), ""}
	update_fingerprint(&a)
	update_range(&a)
	if a.op1.code == "OP" && a.op2.code == "OP" {
		value := calculate(ins.code, a.op1.value, a.op2.value)
		a = ast{"OP", value, value, value, nil, nil, ""}
		update_fingerprint(&a)
		return a
	}
	if a.op1.code == "OP" {
		switch ins.code {
		case "div", "mod":
			if a.op1.value == 0 {
				a = *a.op1
			}
		case "mul":
			if a.op1.value == 0 {
				a = *a.op1
			} else {
				if a.op1.value == 1 {
					a = *a.op2
				}
			}
		case "add":
			if a.op1.value == 0 {
				return *a.op2
			}
		case "eql":
			if a.op1.value < a.op2.min || a.op1.value > a.op2.max {
				return make_ast0()
			}
		}
		return a
	}
	if a.op2.code == "OP" {
		switch ins.code {
		case "mod":
			if a.op1.max < a.op2.value {
				return *a.op1
			}
			if a.op1.code == "add" {
				if r, m := a.op1.op_max_lt(a.op2.value); r != nil {
					if d := m.multiple_of(a.op2.value, false); d != nil {
						return *r
					}
				}
			}
		case "div":
			if a.op2.value == 1 {
				return *a.op1
			}
			if a.op2.value > a.op1.max {
				return make_ast0()
			}
			if a.op1.code == "add" {
				if r, m := a.op1.op_max_lt(a.op2.value); r != nil {
					if d := m.multiple_of(a.op2.value, true); d != nil {
						return *d
					}
				}
			}
		case "eql":
			if a.op2.value < a.op1.min || a.op2.value > a.op1.max {
				return make_ast0()
			}
		case "mul":
			if a.op2.value == 0 {
				return *a.op2
			}
			if a.op2.value == 1 {
				return *a.op1
			}
		case "add":
			if a.op2.value == 0 {
				return *a.op1
			}
		}
		return a
	}
	switch ins.code {
	case "eql":
		if a.op1.min > a.op2.max || a.op2.min > a.op1.max {
			return make_ast0()
		}
	}
	return a
}

func calculate(code string, a, b int64) int64 {
	switch code {
	case "mul":
		return a * b
	case "add":
		return a + b
	case "div":
		return a / b
	case "mod":
		return a % b
	case "eql":
		if a == b {
			return 1
		} else {
			return 0
		}
	default:
		panic("unknown operation")
	}
}

func run_program(program []instruction) [4]int64 {
	registers := [4]int64{}
	//input := [14]int64{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9}
	input := [14]int64{9, 9, 9, 9, 9, 9, 9, 4, 8, 9, 4, 9, 9, 3}

	//input := [14]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	input_index := 0
	for _, ins := range program {
		v1 := ins.op1.value
		switch ins.code {
		case "inp":
			registers[v1] = input[input_index]
			input_index++
		case "add":
			registers[v1] = registers[v1] + ins.op2.eval(registers)
		case "mul":
			registers[v1] = registers[v1] * ins.op2.eval(registers)
		case "div":
			registers[v1] = registers[v1] / ins.op2.eval(registers)
		case "mod":
			registers[v1] = registers[v1] % ins.op2.eval(registers)
		case "eql":
			if registers[v1] == ins.op2.eval(registers) {
				registers[v1] = 1
			} else {
				registers[v1] = 0
			}
		}
		//fmt.Printf("%v\n", ins)
	}
	return registers
}

func make_instruction(tokens []string) instruction {
	ins := instruction{}
	ins.code = tokens[0]
	ins.op1 = make_operand(tokens[1])
	if len(tokens) == 3 {
		ins.op2 = make_operand(tokens[2])
	}
	return ins
}

func make_operand(token string) operand {
	op := operand{}
	if reg, is := is_register(token); is {
		op = operand{true, int64(reg)}
	} else {
		op = operand{false, int64(atoi(token))}
	}
	return op
}

func is_register(token string) (int, bool) {
	if token >= "w" && token <= "z" {
		return map[string]int{"x": 0, "y": 1, "z": 2, "w": 3}[token], true
	}
	return 0, false
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func (op operand) eval(registers [4]int64) int64 {
	if !op.register {
		return op.value
	}
	return registers[op.value]
}

func (op operand) eval_ast(registers *registers_ast) *ast {
	var a ast
	if op.register {
		a = registers[op.value]
	} else {
		a = ast{"OP", op.value, op.value, op.value, nil, nil, ""}
		update_fingerprint(&a)
	}
	return &a
}

func (ins instruction) String() string {
	s := ""
	if ins.code != "inp" {

		s = fmt.Sprintf("%v = <%v %v %v>", ins.op1, ins.op1, sym(ins.code), ins.op2)
	} else {
		s = fmt.Sprintf("<input %v>", ins.op1)
	}
	return s
}

func (op operand) String() string {
	s := ""
	if !op.register {
		s = fmt.Sprintf("%v", op.value)
	} else {
		s = fmt.Sprintf("%v", map[int]string{0: "x", 1: "y", 2: "z", 3: "w"}[int(op.value)])
	}
	return s
}

func (ast ast) String() string {
	switch ast.code {
	case "":
		return "<null>"
	case "OP":
		return fmt.Sprintf("%v", ast.value)
	case "INP":
		return fmt.Sprintf("i%v", ast.value)
	default:
		//return fmt.Sprintf("(%v)", ast.fingerprint)
		return format_ast(&ast, 3, nil)
	}
}

func sym(code string) string {
	return map[string]string{"mul": "*", "add": "+", "eql": "==", "div": "/", "mod": "%"}[code]
}

func update_fingerprint(a *ast) {
	if a.op1 != nil && a.op2 != nil {
		a.fingerprint = sym(a.code) + a.op1.fingerprint + a.op2.fingerprint
	} else {
		a.fingerprint = fmt.Sprintf("C%v", a.value)
	}
	if len(a.fingerprint) > 64 {
		a.fingerprint = compute_hash(a.fingerprint)
	}
}

func update_range(a *ast) {
	if a.code == "OP" {
		a.min = a.value
		a.max = a.value
		return
	}
	if a.code == "INP" {
		a.min = 1
		a.max = 9
		return
	}
	switch a.code {
	case "add", "mul":
		a.min = calculate(a.code, a.op1.min, a.op2.min)
		a.max = calculate(a.code, a.op1.max, a.op2.max)
	case "div":
		a.min = a.op1.min / a.op2.max
		a.max = a.op1.max / a.op2.min
	case "mod":
		if a.op1.min < a.op2.min {
			a.min = a.op1.min
		} else {
			a.min = 0
		}
		a.max = a.op2.max - 1
	case "eql":
		a.min = 0
		a.max = 1
	}
}

func compute_hash(in string) string {
	first := sha256.New()
	first.Write([]byte(in))
	return fmt.Sprintf("%x", first.Sum(nil))
}

type symbol struct {
	sym string
	ast *ast
}

func print_ast(a *ast) {
	symbols := map[string]symbol{}
	println(format_ast(a, 50, &symbols) + "\n")
	queue := []symbol{}
	printed := map[string]bool{}
	enqueue_map(&queue, &symbols)
	for {
		if len(queue) == 0 {
			break
		}
		sym := queue[0]
		queue = queue[1:]
		if _, in := printed[sym.sym]; in {
			continue
		}
		printed[sym.sym] = true
		println(string(sym.sym) + "=" + format_ast(sym.ast, 4, &symbols))
		enqueue_map(&queue, &symbols)
	}
}

func enqueue_map(q *[]symbol, m *map[string]symbol) {
	for _, v := range *m {
		*q = append(*q, v)
	}
}

func format_ast(a *ast, depth int, symbols *map[string]symbol) string {
	if a.op1 == nil {
		return fmt.Sprintf("%v", a)
	}
	if symbols == nil {
		symbols = &map[string]symbol{}
	}
	if depth == 0 {
		if s, in := (*symbols)[a.fingerprint]; in {
			return s.sym
		} else {
			max := symbol{}
			for _, sym := range *symbols {
				if len(sym.sym) > len(max.sym) || (len(sym.sym) == len(max.sym) && sym.sym > max.sym) {
					max = sym
				}
			}
			s := "a"
			if max.sym != "" {
				if max.sym[len(max.sym)-1] == 'z' {
					s = max.sym[0:len(max.sym)-1] + "aa"
				} else {
					s = max.sym[0:len(max.sym)-1] + string(max.sym[len(max.sym)-1]+1)
				}
			}
			(*symbols)[a.fingerprint] = symbol{s, a}
			return string(s)
		}
	}
	if a.code == "eql" {
		return fmt.Sprintf(
			"eql(%v,%v)",
			format_ast(a.op1, depth-1, symbols),
			format_ast(a.op2, depth-1, symbols))
	}
	s := fmt.Sprintf(
		"(%v%v%v)",
		format_ast(a.op1, depth-1, symbols),
		sym(a.code),
		format_ast(a.op2, depth-1, symbols))
	return s
}

func (parent *ast) get_operands() (*ast, *ast) {
	if parent.op1.code == "OP" {
		return parent.op1, parent.op2
	}
	if parent.op2.code == "OP" {
		return parent.op2, parent.op1
	}
	return nil, nil
}

func (parent *ast) op_max_lt(value int64) (*ast, *ast) {
	if parent.op1.max < value {
		return parent.op1, parent.op2
	}
	if parent.op2.max < value {
		return parent.op2, parent.op1
	}
	return nil, nil
}

func (parent *ast) multiple_of(value int64, divide bool) *ast {
	if parent.code != "mul" {
		return nil
	}
	parent_copy := parent.deep_copy()
	new_node := first_multiple_of(parent_copy, parent_copy.op1, parent_copy.op2, value, divide)
	if new_node != nil {
		return new_node.deep_copy()
	}
	new_node = first_multiple_of(parent_copy, parent_copy.op2, parent_copy.op1, value, divide)
	if new_node != nil {
		return new_node.deep_copy()
	}
	new_node = parent_copy.op1.multiple_of(value, divide)
	if new_node != nil {
		parent_copy.op1 = new_node
		return parent_copy
	}
	new_node = parent_copy.op2.multiple_of(value, divide)
	if new_node != nil {
		parent_copy.op2 = new_node
		return parent_copy
	}

	return nil
}

func first_multiple_of(p, n1, n2 *ast, value int64, divide bool) *ast {
	if n1.value >= value && n1.value%value == 0 {
		if divide {
			n1.value /= value
		}
		if n1.value == 1 {
			return n2
		} else {
			return p
		}
	}
	return nil
}

func (src *ast) deep_copy() *ast {
	dst := *src
	if src.op1 != nil {
		dst.op1 = src.op1.deep_copy()
	}
	if src.op2 != nil {
		dst.op2 = src.op2.deep_copy()
	}
	return &dst
}

func shortest_path(start [14]int, goal int, register *ast) {
	queue := []node{{start, run_ast(register, start, nil), 0}}
	visited := map[[14]int]bool{}
	queue[0].update_cost()
	min_value := queue[0].cost
	for len(queue) > 0 {
		number := pop_queue(&queue)
		if _, in := visited[number.input]; in {
			continue
		}
		visited[number.input] = true
		if min_value > number.cost {
			fmt.Printf("best so far: %v\n", number)
			min_value = number.cost
		}
		for i := 13; i >= 0; i-- {
			num_copy := number
			for j := 9; j >= 1; j-- {
				num_copy.input[i] = j
				if num_copy.input[i] >= 1 && num_copy.input[i] <= 9 {
					num_copy.value = run_ast(register, num_copy.input, nil)
					num_copy.update_cost()
					enqueue(&queue, num_copy)
				}
			}
		}
	}
}

func pop_queue(q *[]node) node {
	min_cost := 0
	for i, n := range *q {
		if (*q)[min_cost].cost > n.cost {
			min_cost = i
		}
	}
	top := (*q)[min_cost]
	*q = append((*q)[0:min_cost], (*q)[min_cost+1:]...)
	return top
}

func enqueue(q *[]node, n node) {
	*q = append(*q, n)
}

func (n *node) update_cost() {
	input_value := float64(0)
	for i := 13; i >= 0; i-- {
		input_value *= 10
		input_value += float64(n.input[i])
	}
	n.cost = (float64(n.value) + 0.00001) * input_value
}
