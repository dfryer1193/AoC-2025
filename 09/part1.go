package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

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

	points := make([][2]int, 0)

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		strCoords := strings.Split(line, ",")
		var coords [2]int
		for j, strCoord := range strCoords {
			coord, err := strconv.Atoi(strCoord)
			if err != nil {
				fmt.Println("Error converting coordinate:", err)
				return
			}
			coords[j] = coord
		}
		points = append(points, coords)
	}

	maxRectArea := float64(0)
	for _, p1 := range points {
		for _, p2 := range points {
			// Skip if points are aligned vertically or horizontally
			if p1[0] == p2[0] || p1[1] == p2[1] {
				continue
			}

			length := math.Max(float64(p1[0]), float64(p2[0])) -
				math.Min(float64(p1[0]), float64(p2[0])) + 1
			width := math.Max(float64(p1[1]), float64(p2[1])) -
				math.Min(float64(p1[1]), float64(p2[1])) + 1

			area := length * width
			if area > maxRectArea {
				maxRectArea = area
			}
		}
	}

	fmt.Println("Maximum rectangle area:", maxRectArea)
}
