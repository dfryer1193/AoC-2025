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

	hasMore := true
	scanner := bufio.NewScanner(f)
	for i := 0; hasMore; {
		hasMore = scanner.Scan()
		if hasMore {
			line := scanner.Text()
			warehouse = append(warehouse, line)
		}
		if len(warehouse) < 2 { // Need at least two rows to start checking
			continue
		}

		fmt.Println("Checking row:", i, warehouse[i])
		for j, c := range warehouse[i] {
			if c != '@' {
				continue
			}

			neighbors := countNeighbors(warehouse, j, i)
			if neighbors < 4 {
				fmt.Println(i, j, "is accessible")
				accessibleBoxes++
			}
		}

		i++
	}

	fmt.Println("Total accessible boxes:", accessibleBoxes)
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
