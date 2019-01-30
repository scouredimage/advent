package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	GOBLIN byte = 'G'
	ELF    byte = 'E'
	OPEN   byte = '.'
	WALL   byte = '#'
)

type loc struct {
	x int
	y int
}

func (l *loc) neighbors() []loc {
	// Neighbors in reading order
	return []loc{
		loc{l.x, l.y - 1}, // up
		loc{l.x - 1, l.y}, // left
		loc{l.x + 1, l.y}, // right
		loc{l.x, l.y + 1}, // down
	}
}

type unit struct {
	kind  byte
	l     loc
	hp    int
	power int
	dead  bool
}

func (u *unit) enemy() byte {
	switch u.kind {
	case GOBLIN:
		return ELF
	case ELF:
		return GOBLIN
	}
	panic("unknown unit kind")
}

type combat struct {
	raw   [][]byte
	units []*unit
}

func (c *combat) read() {
	scanner := bufio.NewScanner(os.Stdin)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, _ := range line {
			if line[x] == GOBLIN || line[x] == ELF {
				c.units = append(c.units, &unit{line[x], loc{x, y}, 200, 3, false})
			}
		}
		c.raw = append(c.raw, []byte(line))
		y++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (c *combat) occupied(l loc) bool {
	return c.raw[l.y][l.x] != OPEN
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (c *combat) distance(a loc) [][]int {
	d := make([][]int, len(c.raw))
	q := make([]*loc, 0)
	for y, _ := range c.raw {
		d[y] = make([]int, len(c.raw[y]))
		for x, _ := range c.raw[y] {
			l := loc{x, y}
			if l != a {
				if c.raw[y][x] != OPEN {
					d[y][x] = -1
				} else {
					q = append(q, &loc{x, y})
				}
			}
		}
	}
	for len(q) > 0 {
		next := make([]*loc, 0)
		for _, l := range q {
			var min *int
			for _, n := range l.neighbors() {
				if (n == a || d[n.y][n.x] > 0) && (min == nil || *min > d[n.y][n.x]) {
					min = &d[n.y][n.x]
				}
			}
			if min == nil {
				next = append(next, l)
			} else {
				d[l.y][l.x] = *min + 1
			}
		}
		// Are the only squares left unreachable?
		// e.g. boxed in by walls/other units
		if len(q) != len(next) {
			q = next
			continue
		}
		for i := 0; i < len(q); i++ {
			if q[i] != next[i] {
				q = next
				continue
			}
		}
		// Yes. Mark them as such and bail out!
		for _, l := range q {
			d[l.y][l.x] = -1
		}
		break
	}
	return d
}

func (c *combat) move(u *unit) bool {
	targets := make([]*unit, 0)
	for j, _ := range c.units {
		if c.units[j].l != u.l && c.units[j].kind == u.enemy() && !c.units[j].dead {
			targets = append(targets, c.units[j])
		}
	}
	if len(targets) == 0 {
		return false
	}

	inRange := make([]loc, 0)
	for _, t := range targets {
		squares := t.l.neighbors()
		for j, l := range squares {
			// Already in range?
			if l == u.l {
				c.attack(u)
				return true
			}
			if j == len(squares)-1 {
				for _, s := range squares {
					if !c.occupied(s) {
						inRange = append(inRange, s)
					}
				}
			}
		}
	}

	reachable := make([]loc, 0)
	d := c.distance(u.l)
	for _, r := range inRange {
		if d[r.y][r.x] > -1 {
			reachable = append(reachable, r)
		}
	}

	sort.Slice(reachable, func(a, b int) bool {
		if d[reachable[a].y][reachable[a].x] == d[reachable[b].y][reachable[b].x] {
			if reachable[a].y == reachable[b].y {
				return reachable[a].x < reachable[b].x
			}
			return reachable[a].y < reachable[b].y
		}
		return d[reachable[a].y][reachable[a].x] < d[reachable[b].y][reachable[b].x]
	})
	if len(reachable) == 0 {
		return true
	}

	chosen := &reachable[0]
	attack := false
	if d[chosen.y][chosen.x] > 1 {
		d = c.distance(*chosen)
		chosen = nil
		neighbors := u.l.neighbors()
		for i, n := range neighbors {
			if d[n.y][n.x] > -1 && (chosen == nil || d[chosen.y][chosen.x] > d[n.y][n.x]) {
				chosen = &neighbors[i]
			}
		}
	} else {
		attack = true
	}

	c.raw[u.l.y][u.l.x] = OPEN
	u.l.x, u.l.y = chosen.x, chosen.y
	c.raw[u.l.y][u.l.x] = u.kind

	if attack {
		c.attack(u)
	}

	return true
}

func (c *combat) attack(u *unit) {
	targets := make([]*unit, 0)

	neighbors := u.l.neighbors()
	for _, n := range neighbors {
		for i, e := range c.units {
			if e.l == n && e.kind == u.enemy() && !e.dead {
				targets = append(targets, c.units[i])
			}
		}
	}
	if len(targets) == 0 {
		panic(fmt.Sprintf("%v: no targets to attack!", *u))
	}

	sort.Slice(targets, func(a, b int) bool {
		if targets[a].hp == targets[b].hp {
			if targets[a].l.y == targets[b].l.y {
				return targets[a].l.x < targets[b].l.x
			}
			return targets[a].l.y < targets[b].l.y
		}
		return targets[a].hp < targets[b].hp
	})
	enemy := targets[0]

	enemy.hp = enemy.hp - u.power
	if enemy.hp <= 0 {
		enemy.dead = true
		c.raw[enemy.l.y][enemy.l.x] = OPEN
	}
}

func (c *combat) moveAll() bool {
	for _, u := range c.units {
		if u.dead {
			continue
		}
		if !c.move(u) {
			return false
		}
	}
	sort.Slice(c.units, func(a, b int) bool {
		if c.units[a].l.y == c.units[b].l.y {
			return c.units[a].l.x < c.units[b].l.x
		}
		return c.units[a].l.y < c.units[b].l.y
	})
	return true
}

func (c *combat) print() {
	for i, _ := range c.raw {
		for j, _ := range c.raw[i] {
			if j > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%2s", string(c.raw[i][j]))
		}
		for _, u := range c.units {
			if u.l.y == i && !u.dead {
				fmt.Printf(" %2s(%d)", string(u.kind), u.hp)
			}
		}
		fmt.Println()
	}
}

func (c *combat) solve1() {
	var i int
	for i = 0; c.moveAll(); i++ {
	}
	fmt.Println("Battle ended after", i, "rounds")
	c.print()
	sum := 0
	for _, u := range c.units {
		if !u.dead {
			sum += u.hp
		}
	}
	fmt.Println("Outcome:", i, "*", sum, "=", i*sum)
}

func copyRaw(src [][]byte) [][]byte {
	cp := make([][]byte, len(src))
	for i, _ := range src {
		cp[i] = make([]byte, len(src[i]))
		copy(cp[i], src[i])
	}
	return cp
}

func copyUnits(src []*unit, elfPower int) []*unit {
	cp := make([]*unit, len(src))
	for i, u := range src {
		var power int
		switch u.kind {
		case ELF:
			power = elfPower
		case GOBLIN:
			power = u.power
		default:
			panic("unknown unit kind!")
		}
		cp[i] = &unit{u.kind, loc{u.l.x, u.l.y}, u.hp, power, u.dead}
	}
	return cp
}

func (c *combat) solve2() {
	backupRaw := copyRaw(c.raw)
	backupUnits := copyUnits(c.units, c.units[0].power)

Outer:
	for i := 4; ; i++ {
		fmt.Println("Elf attack power at", i)
		c.raw = copyRaw(backupRaw)
		c.units = copyUnits(backupUnits, i)

		var j int
		for j = 0; c.moveAll(); j++ {
			for _, u := range c.units {
				if u.dead && u.kind == ELF {
					continue Outer
				}
			}
		}
		for k, u := range c.units {
			if u.dead && u.kind == ELF {
				break
			}
			if k == len(c.units)-1 {
				fmt.Println("Battle ended after", j, "rounds")
				c.print()
				sum := 0
				for _, u := range c.units {
					if !u.dead {
						sum += u.hp
					}
				}
				fmt.Println("Outcome:", j, "*", sum, "=", j*sum)
				break Outer
			}
		}
	}
}

func main() {
	c := combat{}

	c.read()
	c.print()

	//c.solve1()
	c.solve2()
}
