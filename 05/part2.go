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

	freshItemCount := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		min, max, err := parseFreshnessRange(line)
		if err != nil {
			fmt.Println("Error parsing freshness range:", err)
			return
		}
		// TODO: Track ranges and avoid double counting
		freshItemsInRange := max - min + 1
		fmt.Printf("Freshness range %d-%d has %d possible fresh items.\n", min, max, freshItemsInRange)
	}

	fmt.Printf("Total possible fresh items: %d\n", freshItemCount)
}

func parseFreshnessRange(line string) (int, int, error) {
	var min, max int
	_, err := fmt.Sscanf(line, "%d-%d", &min, &max)
	if err != nil {
		return 0, 0, err
	}

	return min, max, nil
}
