package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	freq := 0
	for scanner.Scan() {
		change, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("Not a number? %v", err))
		}
		freq += change
		fmt.Printf("change = %d, freq = %d\n", change, freq)
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
}
