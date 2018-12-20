package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	twos, threes := 0, 0
	for scanner.Scan() {
		id := scanner.Text()

		counts := make(map[rune]int)
		for _, r := range id {
			counts[r] += 1
		}

		twice, thrice := false, false
		for r, count := range counts {
			switch count {
			case 2:
				twice = true
				fmt.Printf("%q appears twice in %q\n", r, id)
			case 3:
				thrice = true
				fmt.Printf("%q appears thrice in %q\n", r, id)
			}
		}
		if twice {
			twos++
		}
		if thrice {
			threes++
		}
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
	fmt.Printf("checksum: %d * %d = %d", twos, threes, twos*threes)
}
