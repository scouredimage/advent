package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type point struct {
	x  int
	y  int
	dx int
	dy int
}

func (p *point) move(t int) *point {
	return &point{p.x + t*p.dx, p.y + t*p.dy, p.dx, p.dy}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type bounds struct {
	maxX int
	maxY int
	minX int
	minY int
}

func (b *bounds) area() int {
	return b.maxX - b.minX + b.maxY - b.minY
}

func findMinBounds(points []*point, limit int) int {
	var minB *bounds
	var minT int
	for t := 0; t < limit; t++ {
		b := bounds{}
		for i, p := range points {
			q := p.move(t)
			if i == 0 || b.maxX < q.x {
				b.maxX = q.x
			}
			if i == 0 || b.maxY < q.y {
				b.maxY = q.y
			}
			if i == 0 || b.minX > q.x {
				b.minX = q.x
			}
			if i == 0 || b.minY > q.y {
				b.minY = q.y
			}
		}
		if minB == nil || minB.area() > b.area() {
			minB, minT = &b, t
		}
	}
	return minT
}

func main() {
	points := make([]*point, 0)
	min, max := math.MaxInt16, 0

	re := regexp.MustCompile(`^position=<\s*(?P<x>-?\d+),\s+(?P<y>-?\d+)>\s+velocity=<\s*(?P<dx>-?\d+),\s+(?P<dy>-?\d+)>$`)

	scanner := bufio.NewScanner(os.Stdin)
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
		p := point{groups["x"], groups["y"], groups["dx"], groups["dy"]}
		points = append(points, &p)
		if min > p.x {
			min = p.x
		}
		if max < p.x {
			max = p.x
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Find time instant, t, at which points are closest to each other
	// i.e., the bounding box containing the points has smallest area
	t := findMinBounds(points, 6*60*60) // over 6 hours
	fmt.Println(t)

	// Pretty print
	minX := math.MaxInt16
	for i := 0; i < len(points); i++ {
		points[i] = points[i].move(t)
		if minX > points[i].x {
			minX = points[i].x
		}
	}
	sort.Slice(points, func(i, j int) bool {
		if points[i].y == points[j].y {
			return points[i].x < points[j].x
		}
		return points[i].y < points[j].y
	})
	for i, p := range points {
		if i > 0 && points[i-1].y < p.y {
			fmt.Println()
		}
		j := minX
		if i > 0 && points[i-1].y == p.y {
			j = points[i-1].x + 1
		}
		for ; j < p.x; j++ {
			fmt.Print(" ")
		}
		if i == 0 || points[i-1].y != p.y || points[i-1].x != p.x {
			fmt.Print("#")
		}
	}
	fmt.Println()
}
