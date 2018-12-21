package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type step struct {
	id           rune
	dependencies []*step
	dependents   []*step
}

func (s *step) print() {
	fmt.Printf("%c [", s.id)
	for i, d := range s.dependencies {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%c", d.id)
	}
	fmt.Printf("] [")
	for i, d := range s.dependents {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%c", d.id)
	}
	fmt.Println("]")
}

func mkStep(id rune) *step {
	return &step{id, make([]*step, 0), make([]*step, 0)}
}

func printDAG(steps *map[rune]*step) {
	for _, s := range *steps {
		s.print()
	}
}

func remove(steps *map[rune]*step, s *step, pending []*step) []*step {
	for _, d := range s.dependents {
		for i, x := range d.dependencies {
			if x.id == s.id {
				d.dependencies[i] = d.dependencies[len(d.dependencies)-1]
				d.dependencies = d.dependencies[:len(d.dependencies)-1]
				if len(d.dependencies) == 0 {
					pending = append(pending, d)
				}
				break
			}
		}
	}
	delete(*steps, s.id)
	return pending
}

func work(steps *map[rune]*step, pending []*step) []rune {
	if len(pending) == 0 {
		found := false
		for _, s := range *steps {
			if len(s.dependencies) == 0 {
				pending = append(pending, s)
				found = true
			}
		}
		if !found {
			return make([]rune, 0)
		}
	}

	fmt.Printf("pending: [")
	for i, s := range pending {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%c", s.id)
	}
	fmt.Println("]")

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].id < pending[j].id
	})
	chosen := pending[0]
	pending = pending[1:]

	fmt.Printf("processing: %q\n", chosen.id)
	pending = remove(steps, chosen, pending)
	return append([]rune{chosen.id}, work(steps, pending)...)
}

type worker struct {
	s     *step
	start int
}

func tick(steps *map[rune]*step, pending []*step, workers []*worker, t int) int {
	fmt.Println("clock tick: ", t)

	if len(pending) == 0 {
		found := false
	Loop:
		for _, s := range *steps {
			if len(s.dependencies) == 0 {
				// already being processed?
				for _, w := range workers {
					if w.s != nil && w.s.id == s.id {
						continue Loop
					}
				}
				pending = append(pending, s)
				found = true
			}
		}
		if !found {
			processing := false
			for _, w := range workers {
				if w.s != nil {
					processing = true
					break
				}
			}
			if !processing {
				return t - 1
			}
		}
	}

	fmt.Printf("pending: [")
	for i, s := range pending {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%c", s.id)
	}
	fmt.Println("]")

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].id < pending[j].id
	})

	// process
	for _, w := range workers {
		if w.s != nil {
			if t >= w.start+int(w.s.id-'A')+61 {
				pending = remove(steps, w.s, pending)
				w.s = nil
				w.start = -1
			}
		}
	}

	// assign
	for _, w := range workers {
		if w.s == nil && len(pending) > 0 {
			w.s = pending[0]
			w.start = t
			pending = pending[1:]
		}
	}

	fmt.Printf("worker states: [")
	for i, w := range workers {
		if i > 0 {
			fmt.Printf(", ")
		}
		if w.s != nil {
			fmt.Printf("%c", w.s.id)
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Println("]")

	return tick(steps, pending, workers, t+1)
}

func main() {
	steps := make(map[rune]*step)
	re := regexp.MustCompile(`^Step (?P<id>[A-Z]) must be finished before step (?P<dep>[A-Z]) can begin.$`)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		groups := make(map[string]rune)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				groups[name] = []rune(match[i])[0]
			}
		}
		var s, d *step = nil, nil
		ok := false
		if s, ok = steps[groups["id"]]; !ok {
			s = mkStep(groups["id"])
			steps[groups["id"]] = s
		}
		if d, ok = steps[groups["dep"]]; !ok {
			d = mkStep(groups["dep"])
			steps[groups["dep"]] = d
		}
		s.dependents = append(s.dependents, d)
		d.dependencies = append(d.dependencies, s)
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	printDAG(&steps)

	// In what order should the steps in your instructions be completed?
	// fmt.Println(string(work(&steps, make([]*step, 0))))

	// With 5 workers and the 60+ second step durations, how long will
	// it take to complete all of the steps
	workers := []*worker{
		&worker{nil, -1},
		&worker{nil, -1},
		&worker{nil, -1},
		&worker{nil, -1},
		&worker{nil, -1},
	}
	fmt.Println(tick(&steps, make([]*step, 0), workers, 0))
}
