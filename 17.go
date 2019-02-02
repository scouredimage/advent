package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type pair struct {
	x int
	y int
}

type scan struct {
	min     pair // min x, y
	max     pair // max x, y
	streams []pair
	clay    map[pair]bool
	running map[pair]bool // water flow from stream origin downward
	still   map[pair]bool // water that has collected in clay pockets and come to rest
}

func parse(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("Not a number? %v", err))
	}
	return val
}

func (s *scan) read() {
	// x=123, y=456..789
	re := regexp.MustCompile(`^(?P<dim1>x|y)=(?P<val1>\d+), (?P<dim2>x|y)=(?P<val21>\d+)\.\.(?P<val22>\d+)$`)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		match := re.FindStringSubmatch(line)
		groups := make(map[string]string)
		for i, name := range re.SubexpNames() {
			groups[name] = match[i]
		}

		var x1, x2, y1, y2 int
		if groups["dim1"] == "x" {
			x1 = parse(groups["val1"])
			x2 = x1
			y1 = parse(groups["val21"])
			y2 = parse(groups["val22"])
		} else if groups["dim1"] == "y" {
			y1 = parse(groups["val1"])
			y2 = y1
			x1 = parse(groups["val21"])
			x2 = parse(groups["val22"])
		} else {
			panic(fmt.Sprintf("malformed: %s", line))
		}

		if s.min.x > x1 {
			s.min.x = x1
		}
		if s.max.x < x2 {
			s.max.x = x2
		}
		if s.min.y > y1 {
			s.min.y = y1
		}
		if s.max.y < y2 {
			s.max.y = y2
		}

		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				s.clay[pair{x, y}] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (s *scan) print() {
	for y := s.min.y; y <= s.max.y; y++ {
		for x := s.min.x; x <= s.max.x; x++ {
			if y == s.min.y && x == s.min.x {
				fmt.Println(x)
			}
			if _, ok := s.clay[pair{x, y}]; ok {
				fmt.Printf("#")
			} else {
				if _, ok := s.still[pair{x, y}]; ok {
					fmt.Printf("~")
				} else {
					if _, ok := s.running[pair{x, y}]; ok {
						fmt.Printf("|")
					} else {
						fmt.Printf(".")
					}
				}
			}
			if x == s.max.x && (y == s.min.y || y == s.max.y) {
				fmt.Printf("  %d", y)
			}
		}
		fmt.Println()
		if y == s.max.y {
			fmt.Println(s.max.x)
		}
	}
}

func (s *scan) count() (int, int) {
	total, atRest := 0, 0
	for y := s.min.y; y <= s.max.y; y++ {
		for x := s.min.x; x <= s.max.x; x++ {
			if _, ok := s.clay[pair{x, y}]; !ok {
				if _, ok := s.still[pair{x, y}]; ok {
					total++
					atRest++
				} else {
					if _, ok := s.running[pair{x, y}]; ok {
						total++
					}
				}
			}
		}
	}
	return total, atRest
}

func (s *scan) sand(x, y int) bool {
	l := pair{x, y}
	if _, ok := s.clay[l]; !ok {
		if _, ok := s.still[l]; !ok {
			return true
		}
	}
	return false
}

func (s *scan) tick() {
	for i := 0; i < len(s.streams); {
		x, y := s.streams[i].x, s.streams[i].y
		if !s.sand(x, y) || y > s.max.y { // streams merge || stream overflows beyond bottom of scan
			s.streams[i] = s.streams[len(s.streams)-1]
			s.streams = s.streams[:len(s.streams)-1]
			continue
		}
		if !s.sand(x, y+1) { // clay or water on bottom
			var x1, x2 *int
			for j := x - 1; j >= s.min.x; j-- {
				if s.sand(j, y+1) {
					break
				} else if !s.sand(j, y) {
					x1 = &j
					break
				}
			}
			for j := x + 1; j <= s.max.x; j++ {
				if s.sand(j, y+1) {
					break
				} else if !s.sand(j, y) {
					x2 = &j
					break
				}
			}
			if x1 != nil && x2 != nil {
				for j := *x1; j <= *x2; j++ {
					s.still[pair{j, y}] = true
				}
				y -= 1
			} else if x1 == nil && x2 != nil { // clay wall on right
				for j := *x2 - 1; ; j-- {
					if s.sand(j, y+1) {
						x = j
						break
					}
					s.running[pair{j, y}] = true
				}
			} else if x1 != nil && x2 == nil { // clay wall on left
				for j := *x1 + 1; ; j++ {
					if s.sand(j, y+1) {
						x = j
						break
					}
					s.running[pair{j, y}] = true
				}
			} else { // no walls on either side
				for j := x; ; j-- {
					if s.sand(j, y+1) {
						x = j
						break
					}
					s.running[pair{j, y}] = true
				}
				for j := x + 1; ; j++ {
					if s.sand(j, y+1) {
						s.streams = append(s.streams, pair{j, y})
						break
					}
					s.running[pair{j, y}] = true
				}
			}
		} else {
			s.running[pair{x, y}] = true
			y += 1
		}
		s.streams[i].x, s.streams[i].y = x, y
		if s.min.x > x {
			s.min.x = x
		}
		if s.max.x < x {
			s.max.x = x
		}
		i++
	}
}

func main() {
	s := scan{
		pair{math.MaxInt16, math.MaxInt16},
		pair{0, 0},
		[]pair{pair{500, 0}},
		make(map[pair]bool),
		make(map[pair]bool),
		make(map[pair]bool),
	}
	s.read()
	s.print()
	for len(s.streams) > 0 {
		s.tick()
	}
	s.print()
	fmt.Println(s.count())
}
