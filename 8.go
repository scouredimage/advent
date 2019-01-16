package main

import (
	"fmt"
)

type node struct {
	children []*node
	metadata []int
}

func read() *node {
	n := node{make([]*node, 0), make([]int, 0)}
	var children, metadata int
	fmt.Scan(&children)
	fmt.Scan(&metadata)
	for i := 0; i < children; i++ {
		n.children = append(n.children, read())
	}
	for i := 0; i < metadata; i++ {
		var m int
		fmt.Scan(&m)
		n.metadata = append(n.metadata, m)
	}
	return &n
}

func (n *node) sum() int {
	val := 0
	for _, m := range n.metadata {
		val += m
	}
	for _, c := range n.children {
		val += c.sum()
	}
	return val
}

func (n *node) value() int {
	val := 0
	if len(n.children) == 0 {
		for _, m := range n.metadata {
			val += m
		}
	}
	else {
		for _, m := range n.metadata {
			if m > 0 && m <= len(n.children) {
				val += n.children[m-1].value()
			}
		}
	}
	return val
}

func main() {
	root := read()

	// What is the sum of all metadata entries?
	fmt.Println(root.sum())

	// What is the value of the root node?
	fmt.Println(root.value())
}
