package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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

	junctions := make(map[string]*junction, 0)
	circuits := make([]map[string]*junction, 0)

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

		j := &junction{
			x:               x,
			y:               y,
			z:               z,
			key:             buildKey(x, y, z),
			closestDistance: math.MaxFloat64,
			closestNeighbor: nil,
		}
		if _, ok := junctions[j.key]; ok {
			continue
		}

		junctions[j.key] = j
	}

	// Find closest neighbors
	for _, current := range junctions {
		for _, candidate := range junctions {
			if current.key == candidate.key {
				continue
			}
			dist := current.distanceTo(candidate)
			if dist < current.closestDistance {
				current.closestDistance = dist
				current.closestNeighbor = candidate
			}
		}
	}

	visited := make(map[string]struct{})
	for key, j := range junctions {
		if _, ok := visited[key]; ok {
			continue
		}

		circuit := make(map[string]*junction)
		current := j
		for {
			if _, ok := circuit[current.key]; ok {
				break
			}
			circuit[current.key] = current
			visited[current.key] = struct{}{}
			current = current.closestNeighbor
		}
		circuits = append(circuits, circuit)
	}

	topThreeSizes := []map[string]*junction{{}, {}, {}}
	for _, circuit := range circuits {
		fmt.Println("Circuit:", circuit)
		size := len(circuit)
		if size > len(topThreeSizes[0]) {
			topThreeSizes[2] = topThreeSizes[1]
			topThreeSizes[1] = topThreeSizes[0]
			topThreeSizes[0] = circuit
		} else if size > len(topThreeSizes[1]) {
			topThreeSizes[2] = topThreeSizes[1]
			topThreeSizes[1] = circuit
		} else if size > len(topThreeSizes[2]) {
			topThreeSizes[2] = circuit
		}
	}

	fmt.Println("Top three circuit sizes:", topThreeSizes)
	result := len(topThreeSizes[0]) * len(topThreeSizes[1]) * len(topThreeSizes[2])
	fmt.Println("Result:", result)
}
