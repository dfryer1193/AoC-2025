package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type node struct {
	name     string
	children []string
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Please provide an input filename.")
		return
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	childrenToParents := make(map[string][]string)
	nodes := make(map[string]*node)
	var graph *node

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		segments := strings.Split(line, " ")
		nodeName := getNodeName(segments[0])

		var n *node
		n, ok := nodes[nodeName]
		if !ok {
			n = &node{
				name:     nodeName,
				children: segments[1:],
			}

			nodes[nodeName] = n
		}

		for _, child := range segments[1:] {
			childrenToParents[child] = append(childrenToParents[child], n.name)
		}

		if n.name == "you" {
			graph = n
		}
	}

	nodes["out"] = &node{
		name:     "out",
		children: []string{},
	}

	n := getPathsOut(graph, nodes)
	fmt.Println("Total paths from 'you' to 'out':", n)
}

func getNodeName(s string) string {
	return s[:len(s)-1] // Remove the colon at the end
}

func getPathsOut(n *node, nodes map[string]*node) int {
	return getOut(n, nodes)
}

func getOut(n *node, nodes map[string]*node) int {
	if n.name == "out" {
		return 1
	}

	if len(n.children) == 0 {
		return 0
	}

	childrenLeadingOut := 0
	for _, childName := range n.children {
		childNode, ok := nodes[childName]
		if !ok {
			continue
		}

		childrenLeadingOut += getOut(childNode, nodes)
	}

	return childrenLeadingOut
}
