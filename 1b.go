package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	changes := make([]int, 0)
	for scanner.Scan() {
		change, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("Not a number? %v", err))
		}
		changes = append(changes, change)
	}
	if scanner.Err() != nil {
		panic(fmt.Sprintf("scan error! %v", scanner.Err()))
	}

	freq := 0
	seen := make(map[int]bool)
	for i := 0; ; i++ {
		seen[freq] = true
		if i >= len(changes) {
			i = 0
		}
		change := changes[i]
		freq += change
		fmt.Printf("change = %d, freq = %d\n", change, freq)
		if seen[freq] {
			fmt.Println("calibrated at freq", freq)
			break
		}
	}
}
