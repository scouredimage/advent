package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type dim struct {
	l int // left-offset
	t int // top-offset
	w int // width
	h int // height
}

type claim struct {
	id int
	dim
}

type overlap struct {
	rows   int
	cols   int
	fabric [][]int
}

func (o *overlap) area() int {
	start, total := -1, 0
	for i := 0; i < o.rows; i++ {
		in := false
		for j := 0; j < o.cols; j++ {
			switch o.fabric[i][j] {
			case -1:
				if !in {
					start = j
				}
				in = true
			default:
				if in {
					total += j - start
				}
				in = false
				start = -1
			}
		}
	}
	return total
}

func (o *overlap) plus(c *claim) {
	// init
	if o.fabric == nil {
		o.rows = 10
		o.cols = 10
		o.fabric = make([][]int, o.rows)
		for i := range o.fabric {
			o.fabric[i] = make([]int, o.cols)
		}
	}

	// resize?
	if o.cols < c.l+c.w {
		for i := range o.fabric {
			for j := o.cols; j < c.l+c.w; j++ {
				o.fabric[i] = append(o.fabric[i], 0)
			}
		}
		o.cols = c.l + c.w
	}
	if o.rows < c.t+c.h {
		for i := o.rows; i < c.t+c.h; i++ {
			o.fabric = append(o.fabric, make([]int, o.cols))
		}
		o.rows = c.t + c.h
	}

	for i := c.t; i < c.t+c.h; i++ {
		for j := c.l; j < c.l+c.w; j++ {
			switch o.fabric[i][j] {
			case 0:
				o.fabric[i][j] = c.id
			default:
				o.fabric[i][j] = -1
			}
		}
	}
}

func (o *overlap) print() {
	for i := range o.fabric {
		for j := range o.fabric[i] {
			val := strconv.Itoa(o.fabric[i][j])
			switch val {
			case `0`:
				val = `.`
			case `-1`:
				val = `X`
			}
			fmt.Printf("%s", val)
		}
		fmt.Println("")
	}
}

func main() {
	claims := make([]*claim, 0)
	o := overlap{}

	// #<id> @ <left-offset>,<top-offset>: <width>x<height>
	re := regexp.MustCompile(`^#(?P<id>\d+) @ (?P<l>\d+),(?P<t>\d+): (?P<w>\d+)x(?P<h>\d+)$`)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		groups := make(map[string]int)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				val, err := strconv.Atoi(match[i])
				if err != nil {
					panic(fmt.Sprintf("Not a number? %v", err))
				}
				groups[name] = val
			}
		}
		if groups["id"] <= 0 {
			panic("claim id must be > 0")
		}
		c := claim{groups["id"], dim{groups["l"], groups["t"], groups["w"], groups["h"]}}
		claims = append(claims, &c)
		o.plus(&c)
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	fmt.Println("overlap:", o.area())

	for _, claim := range claims {
		pristine := true
	Loop:
		for i := claim.t; i < claim.t+claim.h; i++ {
			for j := claim.l; j < claim.l+claim.w; j++ {
				switch o.fabric[i][j] {
				case -1:
					pristine = false
					break Loop
				}
			}
		}
		if pristine {
			fmt.Println("claim w/o overlap found! id:", claim.id)
			break
		}
	}
}
