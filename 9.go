package main

import "fmt"

type game struct {
	state   []int
	current int
	size    int
}

func (g *game) place(p *player, m int) {
	if m%23 == 0 {
		p.add(m)
		i := g.current
		for j := 0; j < 7; j++ {
			if i == 0 {
				i = g.size - 1
			} else {
				i--
			}
		}
		p.add(g.state[i])
		for j := i; j < g.size; j++ {
			g.state[j] = g.state[j+1]
		}
		g.current = i % g.size
		g.size--
	} else {
		g.current = (g.current + 2) % g.size
		if g.current == 0 {
			g.current = g.size
		}
		for i := g.size; i >= g.current; i-- {
			g.state[i+1] = g.state[i]
		}
		g.state[g.current] = m
		g.size++
	}
}

func (g *game) print(p *player) {
	fmt.Printf("[%d]", p.id)
	for i := 0; i < g.size; i++ {
		if i == g.current {
			fmt.Printf(" (")
		} else {
			fmt.Printf("  ")
		}
		fmt.Printf("%d", g.state[i])
		if i == g.current {
			fmt.Printf(")")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println("")
}

type player struct {
	id      int
	marbles []int
	score   int
}

func (p *player) add(m int) int {
	p.marbles = append(p.marbles, m)
	p.score += m
	return p.score
}

func main() {
	var players int
	fmt.Printf("# players: ")
	fmt.Scan(&players)

	var last int
	fmt.Printf("last marble: ")
	fmt.Scan(&last)

	g := game{make([]int, last), 0, 1}

	p := make([]player, 0)
	for i := 0; i < players; i++ {
		p = append(p, player{i + 1, make([]int, 0), 0})
	}

	for m := 1; m <= last; m++ {
		fmt.Println(m)
		g.place(&p[(m-1)%len(p)], m)
	}

	winner := -1
	for i := 0; i < len(p); i++ {
		if winner == -1 || p[winner].score < p[i].score {
			winner = i
		}
	}
	fmt.Println(p[winner])
}
