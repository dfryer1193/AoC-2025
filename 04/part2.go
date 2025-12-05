package main

import (
	"bufio"
	"fmt"
	"os"
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

	accessibleBoxes := 0
	warehouse := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		warehouse = append(warehouse, line)
	}

	for {
		removedBoxes := 0
		warehouse, removedBoxes = removeAccessible(warehouse)
		accessibleBoxes += removedBoxes

		if removedBoxes == 0 {
			break
		}
	}

	fmt.Println("Total accessible boxes:", accessibleBoxes)
}

func removeAccessible(warehouse []string) ([]string, int) {
	accessibleBoxes := 0
	next := make([]string, len(warehouse))

	for i, row := range warehouse {
		nextRow := []rune(row)
		for j, c := range row {
			if c != '@' {
				continue
			}

			neighbors := countNeighbors(warehouse, j, i)
			if neighbors < 4 {
				nextRow[j] = '.'
				accessibleBoxes++
			}
		}
		next[i] = string(nextRow)
	}

	return next, accessibleBoxes
}

func countNeighbors(warehouse []string, x int, y int) int {
	neighbors := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			nx, ny := x+i, y+j

			if ny >= 0 && ny < len(warehouse) && nx >= 0 && nx < len(warehouse[ny]) {
				if warehouse[ny][nx] == '@' {
					neighbors++
				}
			}
		}
	}

	return neighbors
}
