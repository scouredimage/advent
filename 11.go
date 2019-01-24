package main

import (
	"fmt"
	"strconv"
)

const serial = 1955

type cell struct {
	x int
	y int
	p *int
}

func (c *cell) power() int {
	if c.p == nil {
		rack := c.x + 10
		level := rack * c.y
		s := strconv.Itoa((level + serial) * rack)
		result, err := strconv.Atoi(string(s[len(s)-3]))
		if err != nil {
			panic(err)
		}
		result -= 5
		c.p = &result
	}
	return *c.p
}

func total(grid [][]*cell, x, y, size int) int {
	t := 0
	for i := y; i < y+size; i++ {
		for j := x; j < x+size; j++ {
			t += grid[i][j].power()
		}
	}
	return t
}

func main() {
	grid := make([][]*cell, 300)
	for y := 0; y < 300; y++ {
		grid[y] = make([]*cell, 300)
		for x := 0; x < 300; x++ {
			grid[y][x] = &cell{x: y + 1, y: x + 1}
		}
	}
	max := struct {
		val *int
		x   int
		y   int
		s   int
	}{
		nil,
		-1,
		-1,
		-1,
	}
	for s := 1; s <= 300; s++ {
		for y := 0; y <= 300-s; y++ {
			for x := 0; x <= 300-s; x++ {
				t := total(grid, x, y, s)
				if max.val == nil || *max.val < t {
					max.val, max.x, max.y, max.s = &t, grid[y][x].x, grid[y][x].y, s
				}
			}
		}
	}
	fmt.Printf("Total power=%d @ %d,%d,%d\n", *max.val, max.x, max.y, max.s)
}
