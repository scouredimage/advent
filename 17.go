package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type loc struct {
	x int
	y int
}

type scan struct {
	x      []*int // [xmin, xmax, stream1x, stream2x,...]
	y      []*int // [ymin, ymax, stream1y, stream2y,...]
	clay   map[loc]bool
	stream map[loc]bool
	water  map[loc]bool
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

		if s.x[0] == nil || *s.x[0] > x1 {
			s.x[0] = &x1
		}
		if s.x[1] == nil || *s.x[1] < x2 {
			s.x[1] = &x2
		}
		if s.y[0] == nil || *s.y[0] > y1 {
			s.y[0] = &y1
		}
		if s.y[1] == nil || *s.y[1] < y2 {
			s.y[1] = &y2
		}

		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				s.clay[loc{x, y}] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (s *scan) print() {
	for y := *s.y[0]; y <= *s.y[1]; y++ {
		for x := *s.x[0]; x <= *s.x[1]; x++ {
			if y == *s.y[0] && x == *s.x[0] {
				fmt.Println(x)
			}
			if _, ok := s.clay[loc{x, y}]; ok {
				fmt.Printf("#")
			} else {
				if _, ok := s.water[loc{x, y}]; ok {
					fmt.Printf("~")
				} else {
					if _, ok := s.stream[loc{x, y}]; ok {
						fmt.Printf("|")
					} else {
						fmt.Printf(".")
					}
				}
			}
			if x == *s.x[1] && (y == *s.y[0] || y == *s.y[1]) {
				fmt.Printf("  %d", y)
			}
		}
		fmt.Println()
		if y == *s.y[1] {
			fmt.Println(*s.x[1])
		}
	}
}

func (s *scan) count() (int, int) {
	total, atRest := 0, 0
	for y := *s.y[0]; y <= *s.y[1]; y++ {
		for x := *s.x[0]; x <= *s.x[1]; x++ {
			if _, ok := s.clay[loc{x, y}]; !ok {
				if _, ok := s.water[loc{x, y}]; ok {
					total++
					atRest++
				} else {
					if _, ok := s.stream[loc{x, y}]; ok {
						total++
					}
				}
			}
		}
	}
	return total, atRest
}

func (s *scan) sand(x, y int) bool {
	l := loc{x, y}
	if _, ok := s.clay[l]; !ok {
		if _, ok := s.water[l]; !ok {
			return true
		}
	}
	return false
}

func (s *scan) tick() bool {
	if len(s.y) <= 2 {
		return false
	}
	for i := 2; i < len(s.y); {
		y, x := *s.y[i], *s.x[i]
		if !s.sand(x, y) || y > *s.y[1] { // streams merge || stream overflows beyond bottom of scan
			s.x[i], s.y[i] = s.x[len(s.x)-1], s.y[len(s.y)-1]
			s.x, s.y = s.x[:len(s.x)-1], s.y[:len(s.y)-1]
			continue
		}
		if !s.sand(x, y+1) { // clay or water on bottom
			var x1, x2 *int
			for j := x - 1; j >= *s.x[0]; j-- {
				if s.sand(j, y+1) {
					break
				} else if !s.sand(j, y) {
					x1 = &j
					break
				}
			}
			for j := x + 1; j <= *s.x[1]; j++ {
				if s.sand(j, y+1) {
					break
				} else if !s.sand(j, y) {
					x2 = &j
					break
				}
			}
			if x1 != nil && x2 != nil {
				for j := *x1; j <= *x2; j++ {
					s.water[loc{j, y}] = true
				}
				y -= 1
			} else if x1 == nil && x2 != nil { // clay wall on right
				for j := *x2 - 1; ; j-- {
					if s.sand(j, y+1) {
						x = j
						break
					}
					s.stream[loc{j, y}] = true
				}
			} else if x1 != nil && x2 == nil { // clay wall on left
				for j := *x1 + 1; ; j++ {
					if s.sand(j, y+1) {
						x = j
						break
					}
					s.stream[loc{j, y}] = true
				}
			} else { // no walls on either side
				for j := x; ; j-- {
					if s.sand(j, y+1) {
						x = j
						break
					}
					s.stream[loc{j, y}] = true
				}
				for j := x + 1; ; j++ {
					if s.sand(j, y+1) {
						s.x = append(s.x, &j)
						s.y = append(s.y, &y)
						break
					}
					s.stream[loc{j, y}] = true
				}
			}
		} else {
			s.stream[loc{x, y}] = true
			y += 1
		}
		s.x[i], s.y[i] = &x, &y
		if *s.x[0] > x {
			s.x[0] = &x
		}
		if *s.x[1] < x {
			s.x[1] = &x
		}
		i++
	}
	return true
}

func main() {
	s := scan{
		[]*int{nil, nil},
		[]*int{nil, nil},
		make(map[loc]bool),
		make(map[loc]bool),
		make(map[loc]bool),
	}
	s.read()

	spring := loc{500, 0}
	s.x = []*int{s.x[0], s.x[1], &spring.x}
	s.y = []*int{s.y[0], s.y[1], &spring.y}

	for s.tick() {
	}
	s.print()
	fmt.Println(s.count())
}
