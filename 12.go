package main

import "fmt"

func match(state string, notes map[string]string, index int, prefix string) string {
	if index >= len(state)+4 {
		return prefix
	}

	sub := ""
	for i := index - 2; i <= index+2; i++ {
		if i < 0 || i >= len(state) {
			sub += "."
		} else {
			sub += string(state[i])
		}
	}

	r, f := notes[sub]
	if !f {
		r = "."
	}
	return match(state, notes, index+1, prefix+r)
}

func sum(state string, start, end int) int {
	s := 0
	for i, j := start, 0; i < end; i, j = i+1, j+1 {
		if state[j] == '#' {
			s += i
		}
	}
	return s
}

func main() {
	notes := map[string]string{
		".#.#.": `.`,
		"...#.": `#`,
		"..##.": `.`,
		"....#": `.`,
		"##.#.": `#`,
		".##.#": `#`,
		".####": `#`,
		"#.#.#": `#`,
		"#..#.": `#`,
		"##..#": `.`,
		"#####": `.`,
		"...##": `.`,
		".#...": `.`,
		"###..": `#`,
		"#..##": `.`,
		"#...#": `.`,
		".#..#": `#`,
		".#.##": `.`,
		"#.#..": `#`,
		".....": `.`,
		"####.": `.`,
		"##.##": `#`,
		"..###": `#`,
		"#....": `.`,
		"###.#": `.`,
		".##..": `#`,
		"#.###": `#`,
		"..#.#": `.`,
		".###.": `#`,
		"##...": `#`,
		"#.##.": `#`,
		"..#..": `#`,
	}
	state := "#.##.###.#.##...##..#..##....#.#.#.#.##....##..#..####..###.####.##.#..#...#..######.#.....#..##...#"

	// After 20 generations, what is the sum of the numbers of all pots which contain a plant?
	for i, length, last := 1, len(state), 0; i <= 20; i++ {
		state = match(state, notes, -3, "")
		start, end := -i*3, length+(i*4)
		s := sum(state, start, end)
		fmt.Printf("gen = %d, sum = %d, delta = %d\n", i, s, s-last)
		last = s
	}

	// After fifty billion generations, what is the sum of the numbers of all pots which contain a plant?
	// ...
	// gen = 97, sum = 6651, delta = 40
	// gen = 98, sum = 6731, delta = 80
	// gen = 99, sum = 6793, delta = 62
	// gen = 100, sum = 6855, delta = 62
	// gen = 101, sum = 6917, delta = 62
	// gen = 102, sum = 6979, delta = 62
	// gen = 103, sum = 7041, delta = 62
	// ...
	// Answer = 6855+((50000000000-100)*62)

}
