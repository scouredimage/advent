package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"unicode"
)

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func react(polymer *[]rune) *[]rune {
	for i, r := range *polymer {
		if i > 0 && abs(r-(*polymer)[i-1]) == 32 {
			res := make([]rune, len(*polymer)-2)
			copy(res, (*polymer)[:i-1])
			copy(res[i-1:], (*polymer)[i+1:])
			return react(&res)
		}
	}
	return polymer
}

func remove(polymer *[]rune, unit rune) *[]rune {
	pruned := make([]rune, 0)
	for _, r := range *polymer {
		if r != unit && abs(r-unit) != 32 {
			pruned = append(pruned, r)
		}
	}
	return &pruned
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	in := []rune(scanner.Text())

	// How many units remain after fully reacting the polymer you scanned?
	out := react(&in)
	fmt.Printf("%s, %d\n", string(*out), len(*out))

	// What is the length of the shortest polymer you can produce by removing
	// all units of exactly one type and fully reacting the result?
	seen := make(map[rune]bool)
	min := math.MaxInt16
	for _, r := range in {
		if _, ok := seen[unicode.ToLower(r)]; !ok {
			out := react(remove(&in, r))
			fmt.Printf("%s, %d\n", string(r), len(*out))
			if min > len(*out) {
				min = len(*out)
			}
			seen[unicode.ToLower(r)] = true
		}
	}
	fmt.Println(min)
}
