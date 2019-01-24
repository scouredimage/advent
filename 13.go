package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	direction_UP       string = `^`
	direction_DOWN     string = `v`
	direction_LEFT     string = `<`
	direction_RIGHT    string = `>`
	direction_STRAIGHT string = `%`
)

type loc struct {
	x int
	y int
}

type cart struct {
	l        loc
	d        string
	lastTurn string
	crashed  bool
}

func (c *cart) turn() string {
	switch c.lastTurn {
	case direction_LEFT:
		c.lastTurn = direction_STRAIGHT
	case direction_STRAIGHT:
		c.lastTurn = direction_RIGHT
	case direction_RIGHT:
		fallthrough
	default:
		c.lastTurn = direction_LEFT
	}
	return c.lastTurn
}

func (c *cart) move(tracks []string) {
	switch c.d {
	case direction_UP:
		c.l.y -= 1
	case direction_DOWN:
		c.l.y += 1
	case direction_LEFT:
		c.l.x -= 1
	case direction_RIGHT:
		c.l.x += 1
	}
	switch string(tracks[c.l.y][c.l.x]) {
	case `+`:
		switch c.turn() {
		case direction_LEFT:
			switch c.d {
			case direction_UP:
				c.d = direction_LEFT
			case direction_DOWN:
				c.d = direction_RIGHT
			case direction_LEFT:
				c.d = direction_DOWN
			case direction_RIGHT:
				c.d = direction_UP
			}
		case direction_STRAIGHT:
			// no op
		case direction_RIGHT:
			switch c.d {
			case direction_UP:
				c.d = direction_RIGHT
			case direction_DOWN:
				c.d = direction_LEFT
			case direction_LEFT:
				c.d = direction_UP
			case direction_RIGHT:
				c.d = direction_DOWN
			}
		}
	case `/`:
		switch c.d {
		case direction_UP:
			c.d = direction_RIGHT
		case direction_DOWN:
			c.d = direction_LEFT
		case direction_LEFT:
			c.d = direction_DOWN
		case direction_RIGHT:
			c.d = direction_UP
		}
	case `\`:
		switch c.d {
		case direction_UP:
			c.d = direction_LEFT
		case direction_DOWN:
			c.d = direction_RIGHT
		case direction_LEFT:
			c.d = direction_UP
		case direction_RIGHT:
			c.d = direction_DOWN
		}
	}
}

type cartMap struct {
	tracks []string
	carts  []*cart
	crash  *loc
}

func (cm *cartMap) read() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cm.tracks = append(cm.tracks, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (cm *cartMap) active() int {
	n := 0
	for _, c := range cm.carts {
		if !c.crashed {
			n++
		}
	}
	return n
}

func (cm *cartMap) parse() {
	for y := 0; y < len(cm.tracks); y++ {
		for x := 0; x < len(cm.tracks[y]); x++ {
			switch string(cm.tracks[y][x]) {
			case direction_UP:
				cm.carts = append(cm.carts, &cart{l: loc{x: x, y: y}, d: direction_UP})
				cm.tracks[y] = cm.tracks[y][:x] + `|` + cm.tracks[y][x+1:]
			case direction_DOWN:
				cm.carts = append(cm.carts, &cart{l: loc{x: x, y: y}, d: direction_DOWN})
				cm.tracks[y] = cm.tracks[y][:x] + `|` + cm.tracks[y][x+1:]
			case direction_LEFT:
				cm.carts = append(cm.carts, &cart{l: loc{x: x, y: y}, d: direction_LEFT})
				cm.tracks[y] = cm.tracks[y][:x] + `-` + cm.tracks[y][x+1:]
			case direction_RIGHT:
				cm.carts = append(cm.carts, &cart{l: loc{x: x, y: y}, d: direction_RIGHT})
				cm.tracks[y] = cm.tracks[y][:x] + `-` + cm.tracks[y][x+1:]
			default:
			}
		}
	}
}

func (cm *cartMap) part1() {
	for i := 1; ; i++ {
		if cm.crash != nil {
			fmt.Printf("#%d: %d,%d", i, cm.crash.x, cm.crash.y)
			break
		}
		cm.sort()
		for i, p := range cm.carts {
			p.move(cm.tracks)
			for j, q := range cm.carts {
				if i != j && p.l == q.l {
					cm.crash = &loc{x: p.l.x, y: p.l.y}
					continue
				}
			}
		}
	}
}

func (cm *cartMap) part2() {
	for i := 1; ; i++ {
		cm.sort()
		for i, p := range cm.carts {
			if !p.crashed {
				p.move(cm.tracks)
				for j, q := range cm.carts {
					if i != j && !q.crashed && p.l == q.l {
						p.crashed, q.crashed = true, true
					}
				}
			}
		}
		if cm.active() <= 1 {
			break
		}
	}
	for _, c := range cm.carts {
		if !c.crashed {
			fmt.Println(*c)
			break
		}
	}
}

func (cm *cartMap) sort() {
	sort.Slice(cm.carts, func(i, j int) bool {
		if cm.carts[i].l.y == cm.carts[j].l.y {
			return cm.carts[i].l.x < cm.carts[j].l.x
		}
		return cm.carts[i].l.y < cm.carts[j].l.y
	})
}

func (cm *cartMap) print() {
	locations := make(map[loc]*cart)
	for _, c := range cm.carts {
		locations[loc{x: c.l.x, y: c.l.y}] = c
	}
	for y, _ := range cm.tracks {
		for x, _ := range cm.tracks[y] {
			if cm.crash != nil && cm.crash.x == x && cm.crash.y == y {
				fmt.Printf("X")
			} else if c, p := locations[loc{x: x, y: y}]; p {
				fmt.Printf("%s", c.d)
			} else {
				fmt.Printf("%s", string(cm.tracks[y][x]))
			}
		}
		fmt.Println()
	}
}

func main() {
	cm := cartMap{}
	cm.read()
	cm.parse()
	//cm.part1()
	cm.part2()
}
