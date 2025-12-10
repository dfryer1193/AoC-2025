package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

	// 1. Read points (red tiles)
	redTiles := make([][2]int, 0)
	xCoordsSet := make(map[int]bool)
	yCoordsSet := make(map[int]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
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
		redTiles = append(redTiles, coords)
		xCoordsSet[coords[0]] = true
		yCoordsSet[coords[1]] = true
	}

	// 2. Coordinate Compression
	xCoords := make([]int, 0, len(xCoordsSet))
	for x := range xCoordsSet {
		xCoords = append(xCoords, x)
	}
	sort.Ints(xCoords)
	xMap := make(map[int]int, len(xCoords))
	for i, x := range xCoords {
		xMap[x] = i
	}

	yCoords := make([]int, 0, len(yCoordsSet))
	for y := range yCoordsSet {
		yCoords = append(yCoords, y)
	}
	sort.Ints(yCoords)
	yMap := make(map[int]int, len(yCoords))
	for i, y := range yCoords {
		yMap[y] = i
	}

	// 3. Build pathTiles set and shapes map for all path points
	pathTiles := make(map[[2]int]bool)
	shapes := make(map[[2]int]rune) // Now this will contain shapes for all path points

	for i := 0; i < len(redTiles); i++ {
		p1 := redTiles[i]
		p2 := redTiles[(i+1)%len(redTiles)]

		// Fill pathTiles and initial shapes for segments
		if p1[0] == p2[0] { // Vertical segment
			for y := min(p1[1], p2[1]); y <= max(p1[1], p2[1]); y++ {
				tile := [2]int{p1[0], y}
				pathTiles[tile] = true
				if _, exists := shapes[tile]; !exists { // Only set if not already determined as a corner
					shapes[tile] = '|'
				}
			}
		} else { // Horizontal segment
			for x := min(p1[0], p2[0]); x <= max(p1[0], p2[0]); x++ {
				tile := [2]int{x, p1[1]}
				pathTiles[tile] = true
				if _, exists := shapes[tile]; !exists { // Only set if not already determined as a corner
					shapes[tile] = '-'
				}
			}
		}
	}

	// Now, determine and overwrite shapes for actual corners (red tiles)
	for i := 0; i < len(redTiles); i++ {
		pCurr := redTiles[i]
		pPrev := redTiles[(i+len(redTiles)-1)%len(redTiles)] // Wraps around
		pNext := redTiles[(i+1)%len(redTiles)]               // Wraps around

		dx1 := pCurr[0] - pPrev[0]
		dy1 := pCurr[1] - pPrev[1]
		dx2 := pNext[0] - pCurr[0]
		dy2 := pNext[1] - pCurr[1]

		var shape rune
		if (dy1 > 0 && dx2 > 0) || (dx1 < 0 && dy2 < 0) { // Segment goes S then E, or W then N  (S-E or W-N)
			shape = 'F'
		} else if (dy1 > 0 && dx2 < 0) || (dx1 > 0 && dy2 < 0) { // Segment goes S then W, or E then N (S-W or E-N)
			shape = '7'
		} else if (dy1 < 0 && dx2 < 0) || (dx1 > 0 && dy2 > 0) { // Segment goes N then W, or E then S (N-W or E-S)
			shape = 'J'
		} else if (dy1 < 0 && dx2 > 0) || (dx1 < 0 && dy2 > 0) { // Segment goes N then E, or W then S (N-E or W-S)
			shape = 'L'
		} else if dy1 != 0 { // This is a straight vertical segment at a red tile (shouldn't happen with distinct corners)
			shape = '|'
		} else if dx1 != 0 { // This is a straight horizontal segment at a red tile (shouldn't happen with distinct corners)
			shape = '-'
		}
		shapes[pCurr] = shape
	}

	// 4. Scanline on compressed grid
	// isInsideCellGrid[iy_cell][ix_cell] will be true if the cell [xCoords[ix_cell], xCoords[ix_cell+1]) x [yCoords[iy_cell], yCoords[iy_cell+1]) is inside
	isInsideCellGrid := make([][]bool, len(yCoords)-1) // iy_cell: 0 to len(yCoords)-2
	for i := range isInsideCellGrid {
		isInsideCellGrid[i] = make([]bool, len(xCoords)-1) // ix_cell: 0 to len(xCoords)-2
	}

	for iy_scanline := 0; iy_scanline < len(yCoords)-1; iy_scanline++ { // Loop for each horizontal scanline that defines cells
		y := yCoords[iy_scanline]
		isInside := false // Reset for each scanline

		// Iterate through all potential vertical grid lines (x-coordinates)
		for ix_boundary := 0; ix_boundary < len(xCoords); ix_boundary++ {
			x := xCoords[ix_boundary]

			// Check if we cross a path segment that flips `isInside`
			if pathTiles[[2]int{x, y}] {
				shape := shapes[[2]int{x, y}]
				// Only these shapes have an "upward" component for the scanline logic
				if shape == '|' || shape == 'L' || shape == 'J' {
					isInside = !isInside
				}
			}

			// The 'isInside' status *after* crossing xCoords[ix_boundary] applies to the cell to its right.
			// This cell has compressed x-index `ix_boundary`.
			// So, if `ix_boundary` is less than `len(xCoords)-1`, there is a cell to its right.
			if ix_boundary < len(xCoords)-1 {
				isInsideCellGrid[iy_scanline][ix_boundary] = isInside
			}
		}
	}

	// 5. Build Summed-Area Table for forbidden cells
	numForbiddenCells := make([][]int, len(yCoords)-1)
	for i := range numForbiddenCells {
		numForbiddenCells[i] = make([]int, len(xCoords)-1)
	}

	for iy_cell := 0; iy_cell < len(yCoords)-1; iy_cell++ {
		for ix_cell := 0; ix_cell < len(xCoords)-1; ix_cell++ {
			if !isInsideCellGrid[iy_cell][ix_cell] {
				numForbiddenCells[iy_cell][ix_cell] = 1
			}
		}
	}

	sat := make([][]int, len(yCoords)) // Size is (len(yCoords)-1)+1 x (len(xCoords)-1)+1
	for i := range sat {
		sat[i] = make([]int, len(xCoords))
	}

	for iy_cell := 0; iy_cell < len(yCoords)-1; iy_cell++ { // Loop up to len-2
		for ix_cell := 0; ix_cell < len(xCoords)-1; ix_cell++ { // Loop up to len-2
			val := numForbiddenCells[iy_cell][ix_cell]
			sat[iy_cell+1][ix_cell+1] = val + sat[iy_cell][ix_cell+1] + sat[iy_cell+1][ix_cell] - sat[iy_cell][ix_cell]
		}
	}

	queryForbidden := func(ix1, iy1, ix2, iy2 int) int {
		if ix1 > ix2 || iy1 > iy2 {
			return 0
		}
		return sat[iy2+1][ix2+1] - sat[iy1][ix2+1] - sat[iy2+1][ix1] + sat[iy1][ix1]
	}

	// 6. Find max rectangle using SAT
	maxArea := 0
	for _, p1 := range redTiles {
		for _, p2 := range redTiles {
			if p1[0] == p2[0] || p1[1] == p2[1] {
				continue
			}

			rectMinX, rectMaxX := min(p1[0], p2[0]), max(p1[0], p2[0])
			rectMinY, rectMaxY := min(p1[1], p2[1]), max(p1[1], p2[1])

			ix1 := xMap[rectMinX]
			ix2 := xMap[rectMaxX]
			iy1 := yMap[rectMinY]
			iy2 := yMap[rectMaxY]

			// Check the cells strictly inside the rectangle boundaries
			if queryForbidden(ix1, iy1, ix2-1, iy2-1) == 0 {
				area := (rectMaxX - rectMinX) * (rectMaxY - rectMinY)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	// The problem asks for area inclusive of boundaries.
	// The logic for area calculation was also simplified. Let's fix it.
	// The check is correct. If the cells between the boundary points are all allowed, the rect is valid.
	// The path tiles are always allowed.
	for _, p1 := range redTiles {
		for _, p2 := range redTiles {
			if p1[0] == p2[0] || p1[1] == p2[1] {
				continue
			}

			rectMinX, rectMaxX := min(p1[0], p2[0]), max(p1[0], p2[0])
			rectMinY, rectMaxY := min(p1[1], p2[1]), max(p1[1], p2[1])

			ix1 := xMap[rectMinX]
			ix2 := xMap[rectMaxX]
			iy1 := yMap[rectMinY]
			iy2 := yMap[rectMaxY]

			if queryForbidden(ix1, iy1, ix2-1, iy2-1) == 0 {
				area := (rectMaxX - rectMinX + 1) * (rectMaxY - rectMinY + 1)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	fmt.Println("Maximum rectangle area:", maxArea)
}
