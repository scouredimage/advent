package main

import (
	"fmt"
)

func print(recipes, elves []int) {
	for i, r := range recipes {
		if i == elves[0] {
			fmt.Printf("(%d)", r)
		} else if i == elves[1] {
			fmt.Printf("[%d]", r)
		} else {
			fmt.Printf(" %d ", r)
		}
	}
	fmt.Println()
}

func next(recipes, elves []int) []int {
	n := recipes[elves[0]] + recipes[elves[1]]
	if n/10 > 0 {
		recipes = append(recipes, n/10)
		n %= 10
	}
	recipes = append(recipes, n)

	for i := 0; i < len(elves); i++ {
		for j, k := recipes[elves[i]]+1, elves[i]; j > 0; j-- {
			k = (k + 1) % len(recipes)
			elves[i] = k
		}
	}

	return recipes
}

func part1(recipes, elves []int, target int) {
	for len(recipes) < 10+target {
		//print(recipes, elves)
		recipes = next(recipes, elves)
	}
	//print(recipes, elves)
	fmt.Println(recipes[target : target+10])
}

func part2(recipes, elves, target []int) {
	for i := 1; ; i++ {
		recipes = next(recipes, elves)
		if len(recipes) >= len(target) {
			// Only need to check the (up to) 2 last windows of length=target
			// that (might) have been added in the last step.
			for j := 0; j < 2; j++ {
				for k, l := len(target)-1, 0; k >= 0; k, l = k-1, l+1 {
					if recipes[len(recipes)-j-l-1] == target[k] {
						if k == 0 {
							fmt.Println(len(recipes) - len(target) - j)
							return
						}
					} else {
						break
					}
				}
			}
		}
	}
}

func main() {
	recipes := []int{3, 7}
	elves := []int{0, 1}
	//part1(recipes, elves, 330121)
	part2(recipes, elves, []int{3, 3, 0, 1, 2, 1})
}
