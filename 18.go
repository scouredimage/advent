package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	OPEN       byte = '.'
	TREES      byte = '|'
	LUMBERYARD byte = '#'
)

func read() [][]byte {
	scan := make([][]byte, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		scan = append(scan, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return scan
}

func print(scan [][]byte) {
	for i, _ := range scan {
		fmt.Println(string(scan[i]))
	}
}

func match(scan []byte, kind byte) int {
	n := 0
	for _, b := range scan {
		if b == kind {
			n++
		}
	}
	return n
}

func neighbors(x, y, rows, cols int, f func(i, j int)) {
	for i := y - 1; i <= y+1; i++ {
		if i >= 0 && i < rows {
			for j := x - 1; j <= x+1; j++ {
				if i == y && j == x {
					continue
				}
				if j >= 0 && j < cols {
					f(i, j)
				}
			}
		}
	}
}

func val(scan [][]byte, until int) {
	type pair struct {
		y int
		x int
		i int
		j int
	}

	chans := make(map[pair]chan byte)
	var done sync.WaitGroup
	for y, _ := range scan {
		for x, _ := range scan[y] {
			neighbors(x, y, len(scan), len(scan[y]), func(i, j int) {
				chans[pair{y, x, i, j}] = make(chan byte, 1)
			})
			done.Add(1)
		}
	}

	for y, _ := range scan {
		for x, _ := range scan[y] {
			go func(x, y int) {
				b := scan[y][x]
				neighborhood := make([]byte, 9)
				for k := 1; k <= until; k++ {
					idx := 0
					neighbors(x, y, len(scan), len(scan[y]), func(i, j int) {
						chans[pair{y, x, i, j}] <- b
						neighborhood[idx] = <-chans[pair{i, j, y, x}]
						idx++
					})
					if b == OPEN && match(neighborhood, TREES) >= 3 {
						b = TREES
					} else if b == TREES && match(neighborhood, LUMBERYARD) >= 3 {
						b = LUMBERYARD
					} else if b == LUMBERYARD {
						if match(neighborhood, LUMBERYARD) >= 1 && match(neighborhood, TREES) >= 1 {
							b = LUMBERYARD
						} else {
							b = OPEN
						}
					}
					if k%1000000 == 0 {
						fmt.Println(k, y, x)
					}
				}
				scan[y][x] = b
				done.Done()
			}(x, y)
		}
	}
	done.Wait()
}

func equal(a, b [][]byte) bool {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	scan := read()
	print(scan)

	//What will the total resource value of the lumber collection area be after 10 minutes?
	//val(scan, 10)

	// After 1,000,000,000 minutes?
	// Evaluating all 1 billion minutes is computationally impossible!
	// Instead, try to see if the landscape reaches a steady state within a large enough (1000 mins) timeframe.
	/*
		after := make([][][]byte, 1000)
		after[0] = scan
		for i := 1; i < 1000; i++ {
			tmp := make([][]byte, len(after[i-1]))
			for j := 0; j < len(after[i-1]); j++ {
				tmp[j] = make([]byte, len(after[i-1][j]))
				copy(tmp[j], after[i-1][j])
			}
			val(tmp, 1)
			after[i] = tmp
			fmt.Println("After", i, "minutes:")
			print(after[i])
		}

		for i := 0; i < len(after); i++ {
			fmt.Printf("%3d: ", i)
			for j := 0; j < len(after); j++ {
				if i == j {
					continue
				}
				if equal(after[i], after[j]) {
					fmt.Printf("%3d ", j)
				}
			}
			fmt.Println()
		}
	*/
	//Relevant output:
	// ...
	// 568:
	// 569: 597 625 653 681 709 737 765 793 821 849 877 905 933 961 989
	// ...
	// 579: 607 635 663 691 719 747 775 803 831 859 887 915 943 971 999
	// ...
	// i.e., starting at minute 569, the landscape pattern repeats itself at 28 minute periods.
	// 579 + 28 * 35714265 = 999,999,999
	// Essentially, the state of the landcape after minute 580 will be the same as after 1 billion.
	val(scan, 580)
	print(scan)

	wooded := 0
	for _, l := range scan {
		wooded += match(l, TREES)
	}
	lumberyards := 0
	for _, l := range scan {
		lumberyards += match(l, LUMBERYARD)
	}
	fmt.Println(wooded, "*", lumberyards, "=", wooded*lumberyards)
}
