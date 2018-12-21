package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type point struct {
	x  int
	y  int
	id int
}

func (p *point) dist(q point) int {
	return abs(p.x-q.x) + abs(p.y-q.y)
}

type grid struct {
	points []*point
	matrix [][]*point
	rows   int
	cols   int
}

func (g *grid) append(p point) {
	// init
	if g.matrix == nil {
		g.rows = 0
		g.cols = 0
		g.matrix = make([][]*point, g.rows)
		for i := range g.matrix {
			g.matrix[i] = make([]*point, g.cols)
		}
		g.points = make([]*point, 0)
	}

	// resize?
	if g.cols < p.x+1 {
		for i := range g.matrix {
			for j := g.cols; j < p.x+1; j++ {
				g.matrix[i] = append(g.matrix[i], nil)
			}
		}
		g.cols = p.x + 1
	}
	if g.rows < p.y+1 {
		for i := g.rows; i < p.y+1; i++ {
			g.matrix = append(g.matrix, make([]*point, g.cols))
		}
		g.rows = p.y + 1
	}

	g.matrix[p.y][p.x] = &p
	g.points = append(g.points, &p)
}

func (g *grid) closest(p *point) *point {
	min := struct {
		val        int
		coordinate *point
	}{math.MaxInt16, nil}
	for _, q := range g.points {
		d := p.dist(*q)
		if min.val > d {
			min.val = d
			min.coordinate = q
		} else if min.val == d {
			min.coordinate = nil
		}
	}
	return min.coordinate
}

func (g *grid) assign() {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			switch g.matrix[i][j] {
			case nil:
				g.matrix[i][j] = g.closest(&point{j, i, -1})
			}
		}
	}
}

func (g *grid) isInfArea(p *point) bool {
	if p == g.matrix[p.y][0] || p == g.matrix[p.y][g.cols-1] || p == g.matrix[g.rows-1][p.x] || p == g.matrix[0][p.x] {
		return true
	}
	return false
}

func (g *grid) maxArea() int {
	g.assign()

	max := struct {
		area       int
		coordinate *point
	}{-1, nil}

	byCoordinate := make(map[*point]int)
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			area := -1
			ok := false
			p := g.matrix[i][j]
			switch p {
			case nil:
			default:
				if area, ok = byCoordinate[p]; ok {
					area++
				} else {
					area = 1
				}
				byCoordinate[p] = area
				if max.area < area && !g.isInfArea(p) {
					max.area = area
					max.coordinate = p
				}
			}
		}
	}
	return max.area
}

func (g *grid) markLocationsInRange(cutoff int) {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			total := 0
			for _, p := range g.points {
				total += p.dist(point{j, i, -1})
			}
			if total < cutoff {
				g.matrix[i][j] = &point{0, 0, 0}
			} else {
				g.matrix[i][j] = nil
			}
		}
	}
}

func (g *grid) areaWithinRange() int {
	g.markLocationsInRange(10000)
	area := 0
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			p := g.matrix[i][j]
			if p != nil {
				area++
			}
		}
	}
	return area
}

func main() {
	grid := grid{}
	re := regexp.MustCompile(`^(?P<x>\d+),\s+(?P<y>\d+)$`)

	scanner := bufio.NewScanner(os.Stdin)
	id := 1
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		groups := make(map[string]int)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				val, err := strconv.Atoi(match[i])
				if err != nil {
					panic(fmt.Sprintf("not a number? %v", err))
				}
				groups[name] = val
			}
		}
		grid.append(point{groups["x"], groups["y"], id})
		id++
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	// What is the size of the largest area that isn't infinite?
	fmt.Println(grid.maxArea())

	// What is the size of the region containing all locations which have
	// a total distance to all given coordinates of less than 10000?
	fmt.Println(grid.areaWithinRange())
}
