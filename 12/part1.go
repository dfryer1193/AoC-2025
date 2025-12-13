package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type shape struct {
	area int
}

type grid struct {
	width  int
	height int
	area   int

	shapeCount []int
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

	shapes := make([]*shape, 0)
	grids := make([]*grid, 0)
	idx := -1

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		if cidx := strings.Index(line, ":"); cidx != -1 {
			if !strings.Contains(line, "x") {
				idx, _ = strconv.Atoi(strings.TrimSpace(line[:cidx]))
				shapes = append(shapes, &shape{})
				continue
			}

			sizes := strings.Split(line[:cidx], "x")

			w, _ := strconv.Atoi(sizes[0])
			h, _ := strconv.Atoi(sizes[1])

			rawShapeCounts := strings.Fields(strings.TrimSpace(line[cidx+1:]))
			shapeCounts := make([]int, len(rawShapeCounts))
			for j, c := range rawShapeCounts {
				shapeCounts[j], _ = strconv.Atoi(string(c))
			}

			grids = append(grids, &grid{
				width:      w,
				height:     h,
				area:       w * h,
				shapeCount: shapeCounts,
			})
		}

		for _, c := range line {
			if c == '#' {
				shapes[idx].area++
			}
		}
	}

	canFitAllCount := 0

	for _, g := range grids {
		totalShapesArea := 0
		for i, count := range g.shapeCount {
			totalShapesArea += count * shapes[i].area
		}

		if totalShapesArea <= g.area {
			canFitAllCount++
		}
	}

	fmt.Println("Total number of grids that can fit all shapes:", canFitAllCount)
}
