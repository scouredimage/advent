package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	ids := make([]string, 0)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

Loop:
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if len(ids[i]) != len(ids[j]) {
				panic(fmt.Sprintf("can't compare varying lengths: %q and %q", ids[i], ids[j]))
			}

			diffs, diff := 0, -1
			for k := 0; k < len(ids[i]); k++ {
				if ids[i][k] != ids[j][k] {
					diffs++
					diff = k
				}
			}
			switch diffs {
			case 1:
				fmt.Printf("Found match: %q and %q\n", ids[i], ids[j])
				common := make([]rune, len(ids[i])-1)
				p := 0
				for q, r := range ids[i] {
					if q == diff {
						continue
					}
					common[p] = r
					p++
				}
				fmt.Println(string(common))
				break Loop
			}
		}
	}
}
