package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func addr(a, b, c int, register []int) {
	register[c] = register[a] + register[b]
}

func addi(a, b, c int, register []int) {
	register[c] = register[a] + b
}

func mulr(a, b, c int, register []int) {
	register[c] = register[a] * register[b]
}

func muli(a, b, c int, register []int) {
	register[c] = register[a] * b
}

func banr(a, b, c int, register []int) {
	register[c] = register[a] & register[b]
}

func bani(a, b, c int, register []int) {
	register[c] = register[a] & b
}

func borr(a, b, c int, register []int) {
	register[c] = register[a] | register[b]
}

func bori(a, b, c int, register []int) {
	register[c] = register[a] | b
}

func setr(a, _, c int, register []int) {
	register[c] = register[a]
}

func seti(a, _, c int, register []int) {
	register[c] = a
}

func gtir(a, b, c int, register []int) {
	if a > register[b] {
		register[c] = 1
	} else {
		register[c] = 0
	}
}

func gtri(a, b, c int, register []int) {
	if register[a] > b {
		register[c] = 1
	} else {
		register[c] = 0
	}
}

func gtrr(a, b, c int, register []int) {
	if register[a] > register[b] {
		register[c] = 1
	} else {
		register[c] = 0
	}
}

func eqir(a, b, c int, register []int) {
	if a == register[b] {
		register[c] = 1
	} else {
		register[c] = 0
	}
}

func eqri(a, b, c int, register []int) {
	if register[a] == b {
		register[c] = 1
	} else {
		register[c] = 0
	}
}

func eqrr(a, b, c int, register []int) {
	if register[a] == register[b] {
		register[c] = 1
	} else {
		register[c] = 0
	}
}

type opcode func(a, b, c int, register []int)

var opcodes map[string]opcode = map[string]opcode{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}

func eq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func match(before, instruction, after []int) []string {
	test := make([]int, len(before))
	matches := make([]string, 0)
	for n, f := range opcodes {
		copy(test, before)
		f(instruction[1], instruction[2], instruction[3], test)
		if eq(test, after) {
			matches = append(matches, n)
		}
	}
	return matches
}

func parse(s, sep string) []int {
	parsed := make([]int, 0)
	for _, a := range strings.Split(s, sep) {
		i, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		parsed = append(parsed, i)
	}
	return parsed
}

func read1(scanner *bufio.Scanner) ([][]int, [][]int, [][]int) {
	before := make([][]int, 0)
	instruction := make([][]int, 0)
	after := make([][]int, 0)

	empties := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			empties++
			if empties > 1 {
				break
			}
			continue
		} else if strings.HasPrefix(line, "Before:") {
			before = append(before, parse(string(line[9:len(line)-1]), ", "))
		} else if strings.HasPrefix(line, "After:") {
			after = append(after, parse(string(line[9:len(line)-1]), ", "))
		} else {
			instruction = append(instruction, parse(line, " "))
		}
		empties = 0
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return before, instruction, after
}

func read2(scanner *bufio.Scanner) [][]int {
	instruction := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		instruction = append(instruction, parse(line, " "))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return instruction
}

func read() ([][]int, [][]int, [][]int, [][]int) {
	scanner := bufio.NewScanner(os.Stdin)
	before, instruction, after := read1(scanner)
	program := read2(scanner)
	return before, instruction, after, program
}

func main() {
	before, instruction, after, program := read()
	resolved, unresolved := make(map[int]string), make(map[int][]string)

	// How many samples in your puzzle input behave like three or more opcodes?
	samples := 0
	for i := 0; i < len(before); i++ {
		matches := match(before[i], instruction[i], after[i])
		if len(matches) == 1 {
			resolved[instruction[i][0]] = matches[0]
		} else {
			if len(matches) >= 3 {
				samples++
			}
			unresolved[instruction[i][0]] = matches
		}
	}
	fmt.Println(samples)

	// What value is contained in register 0 after executing the test program?
	for len(unresolved) > 0 {
		for _, r := range resolved {
			for op, codes := range unresolved {
				removed := make([]string, 0)
				for _, c := range codes {
					if c != r {
						removed = append(removed, c)
					}
				}
				if len(removed) == 1 {
					resolved[op] = removed[0]
					delete(unresolved, op)
				} else {
					unresolved[op] = removed
				}
			}
		}
	}
	registers := []int{0, 0, 0, 0}
	for _, i := range program {
		f := opcodes[resolved[i[0]]]
		f(i[1], i[2], i[3], registers)
	}
	fmt.Println(registers)
}
