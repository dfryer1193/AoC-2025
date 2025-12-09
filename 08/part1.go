package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type junction struct {
	x int
	y int
	z int

	key string

	closestNeighbor *junction
	closestDistance float64
}

func buildKey(x, y, z int) string {
	return fmt.Sprintf("%d,%d,%d", x, y, z)
}

func (j *junction) distanceTo(other *junction) float64 {
	dx := math.Pow(float64(j.x-other.x), 2)
	dy := math.Pow(float64(j.y-other.y), 2)
	dz := math.Pow(float64(j.z-other.z), 2)
	return math.Sqrt(dx + dy + dz)
}

type pair struct {
	a, b     *junction
	distance float64
}

// DSU (Disjoint Set Union) functions
func find(parents map[string]string, key string) string {
	if parents[key] == key {
		return key
	}
	parents[key] = find(parents, parents[key]) // Path compression
	return parents[key]
}

func union(parents map[string]string, a, b string) {
	rootA := find(parents, a)
	rootB := find(parents, b)
	if rootA != rootB {
		parents[rootB] = rootA
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage: go run part1.go <filename> <connections_to_make>")
		return
	}

	filename := args[0]
	limitStr := args[1]
	mergeLimit, err := strconv.Atoi(limitStr)
	if err != nil {
		fmt.Println("Error: Invalid number for connections to make:", err)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	junctions := make([]*junction, 0)
	junctionMap := make(map[string]*junction)

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		coords := strings.Split(line, ",")
		if len(coords) != 3 {
			fmt.Println("Invalid line:", line)
			return
		}

		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])

		key := buildKey(x, y, z)
		if _, ok := junctionMap[key]; ok {
			continue
		}

		j := &junction{
			x:   x,
			y:   y,
			z:   z,
			key: key,
		}
		junctions = append(junctions, j)
		junctionMap[key] = j
	}

	// Generate all unique pairs
	pairs := make([]pair, 0)
	for i := 0; i < len(junctions); i++ {
		for k := i + 1; k < len(junctions); k++ {
			p := pair{
				a:        junctions[i],
				b:        junctions[k],
				distance: junctions[i].distanceTo(junctions[k]),
			}
			pairs = append(pairs, p)
		}
	}

	// Sort pairs by distance
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].distance < pairs[j].distance
	})

	// Initialize DSU
	parents := make(map[string]string)
	for _, j := range junctions {
		parents[j.key] = j.key
	}

	// Process the N shortest connections, where N is the mergeLimit
	for i, p := range pairs {
		if i >= mergeLimit {
			break
		}
		if find(parents, p.a.key) != find(parents, p.b.key) {
			union(parents, p.a.key, p.b.key)
		}
	}

	// Group junctions by their circuit root
	circuits := make(map[string][]string)
	for _, j := range junctions {
		root := find(parents, j.key)
		circuits[root] = append(circuits[root], j.key)
	}

	// Find the sizes of all circuits
	allSizes := make([]int, 0, len(circuits))
	for _, circuit := range circuits {
		allSizes = append(allSizes, len(circuit))
	}

	// Sort sizes in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(allSizes)))

	if len(allSizes) < 3 {
		fmt.Println("Error: Less than three circuits found.")
		return
	}

	// Multiply the sizes of the three largest circuits
	topThree := allSizes[:3]
	result := topThree[0] * topThree[1] * topThree[2]

	fmt.Println("Sizes of the three largest circuits:", topThree)
	fmt.Println("Final Result:", result)
}
