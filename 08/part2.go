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
	if len(args) < 1 {
		fmt.Println("Usage: go run part2.go <filename>")
		return
	}

	filename := args[0]
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

	// Keep track of the last pair that successfully merges two circuits
	var lastConnectedPair pair

	// Process all pairs, the last one to cause a merge is our answer
	for _, p := range pairs {
		if find(parents, p.a.key) != find(parents, p.b.key) {
			union(parents, p.a.key, p.b.key)
			lastConnectedPair = p
		}
	}

	// The lastConnectedPair holds the two junctions that made the final connection
	lastJunctionA := lastConnectedPair.a
	lastJunctionB := lastConnectedPair.b
	result := lastJunctionA.x * lastJunctionB.x

	fmt.Printf("Last connection made between junction %s and %s\n", lastJunctionA.key, lastJunctionB.key)
	fmt.Printf("Multiplying their X coordinates (%d * %d)\n", lastJunctionA.x, lastJunctionB.x)
	fmt.Println("Final Result:", result)
}
