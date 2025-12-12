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
	}

	nodes["out"] = &node{
		name:     "out",
		children: []string{},
	}

	n := getPathsOut(nodes)
	fmt.Println("Total paths from 'svr' to 'out' including 'fft' and 'dac':", n)
}

func getNodeName(s string) string {
	return s[:len(s)-1] // Remove the colon at the end
}

func getPathsOut(nodes map[string]*node) int {
	svr := nodes["svr"]
	if svr == nil {
		return 0
	}

	totalPaths := 0
	svrToFft := getTo(svr, "fft", nodes, make(map[string]struct{}))
	fftToDac := getTo(nodes["fft"], "dac", nodes, make(map[string]struct{}))
	dacToOut := getTo(nodes["dac"], "out", nodes, make(map[string]struct{}))
	//
	//	svrToDac := getTo(svr, "dac", nodes, make(map[string]struct{}))
	//	fmt.Println("Paths from svr to dac:", svrToDac)
	//	dacToFft := getTo(nodes["dac"], "fft", nodes, make(map[string]struct{}))
	//	fmt.Println("Paths from dac to fft:", dacToFft)
	//	fftToOut := getTo(nodes["fft"], "out", nodes, make(map[string]struct{}))
	//	fmt.Println("Paths from fft to out:", fftToOut)

	totalPaths += svrToFft * fftToDac * dacToOut
	//	totalPaths += svrToDac * dacToFft * fftToOut

	return totalPaths
}

func getTo(n *node, target string, nodes map[string]*node, noPathToTarget map[string]struct{}) int {
	if n.name == target {
		return 1
	}

	if len(n.children) == 0 {
		noPathToTarget[n.name] = struct{}{}
		return 0
	}

	childrenLeadingOut := 0
	for _, childName := range n.children {
		if _, ok := noPathToTarget[childName]; ok {
			continue
		}

		childNode, ok := nodes[childName]
		if !ok {
			continue
		}

		childrenLeadingOut += getTo(childNode, target, nodes, noPathToTarget)
	}

	if childrenLeadingOut == 0 {
		noPathToTarget[n.name] = struct{}{}
	}

	return childrenLeadingOut
}
