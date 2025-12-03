package main

import (
	"bufio"
	"fmt"
	"log"
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

	sum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		ranges := strings.Split(line, ",")
		for _, r := range ranges {
			sum += processRange(r)
		}
	}

	fmt.Println("Total sum:", sum)
}

func processRange(r string) int {
	parts := strings.Split(r, "-")
	if len(parts) != 2 {
		log.Fatalf("Invalid range format: %s", r)
		return 0
	}

	start, _ := strconv.Atoi(parts[0])
	end, _ := strconv.Atoi(parts[1])

	sum := 0

	for i := start; i <= end; i++ {
		numStr := strconv.Itoa(i)

		possibleSizes := possibleGramSizes(numStr)
		for j := len(possibleSizes) - 1; j >= 0; j-- { // Check larger ngram sizes first for fewer comparisons
			size := possibleSizes[j]
			if !isValidID(numStr, size) {
				sum += i
				break // As soon as we find an invalid number, we can stop checking further gram sizes
			}
		}
	}

	return sum
}

func possibleGramSizes(numStr string) []int {
	length := len(numStr)
	sizes := make([]int, 0)
	if length > 1 {
		sizes = append(sizes, 1)
	}

	// We can only have ngrams of sizes that divide the length evenly
	for size := 2; size <= length/2; size++ {
		if length%size == 0 {
			sizes = append(sizes, size)
		}
	}

	return sizes
}

func isValidID(numStr string, gramSize int) bool {
	// If the ngram repeats for the	 entire length, it's invalid
	pattern := numStr[:gramSize]
	for i := gramSize; i < len(numStr); i += gramSize {
		if numStr[i:i+gramSize] != pattern {
			return true
		}
	}
	return false
}
