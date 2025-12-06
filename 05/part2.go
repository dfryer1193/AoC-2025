package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	ranges := make([][2]int, 0)

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

		ranges = append(ranges, [2]int{min, max})
	}

	for {
		merged, mergecount := mergeRanges(ranges)
		if mergecount == 0 {
			break
		}

		ranges = merged
	}

	for _, r := range ranges {
		freshItemCount += r[1] - r[0] + 1
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

func mergeRanges(ranges [][2]int) ([][2]int, int) {
	if len(ranges) == 0 {
		return ranges, 0
	}
	// Sort ranges by start so a single linear pass can merge all overlaps
	sort.Slice(ranges, func(i, j int) bool { return ranges[i][0] < ranges[j][0] })

	merged := make([][2]int, 0, len(ranges))
	mergecount := 0

	cur := ranges[0]
	for i := 1; i < len(ranges); i++ {
		r := ranges[i]
		if r[0] > cur[1]+1 {
			// disjoint
			merged = append(merged, cur)
			cur = r
		} else {
			// overlapping or adjacent: merge
			if r[1] > cur[1] {
				cur[1] = r[1]
			}
			mergecount++
		}
	}
	merged = append(merged, cur)
	return merged, mergecount
}
