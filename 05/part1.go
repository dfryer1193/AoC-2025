package main

import (
	"bufio"
	"fmt"
	"os"
)

type parseState int

const (
	parsingFreshnessRanges parseState = iota
	parsingItemIDs
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

	parserState := parsingFreshnessRanges
	ranges := make([][2]int, 0)
	freshItemCount := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parserState = parsingItemIDs
			continue
		}

		switch parserState {
		case parsingFreshnessRanges:
			min, max, err := parseFreshnessRange(line)
			if err != nil {
				fmt.Println("Error parsing freshness range:", err)
				return
			}
			ranges = append(ranges, [2]int{min, max})
		case parsingItemIDs:
			id, err := parseItemID(line)
			if err != nil {
				fmt.Println("Error parsing item ID:", err)
				return
			}
			if isFresh(id, ranges) {
				freshItemCount++
			}
		}
	}

	fmt.Printf("Total fresh items: %d\n", freshItemCount)
}

func parseFreshnessRange(line string) (int, int, error) {
	var min, max int
	_, err := fmt.Sscanf(line, "%d-%d", &min, &max)
	if err != nil {
		return 0, 0, err
	}

	return min, max, nil
}

func parseItemID(line string) (int, error) {
	var id int
	_, err := fmt.Sscanf(line, "%d", &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func isFresh(id int, ranges [][2]int) bool {
	for _, r := range ranges {
		if !(id < r[0] || id > r[1]) {
			return true
		}
	}

	return false
}
