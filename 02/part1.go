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

		// If we don't have an even number of digits, the number can't be made of two runs
		if len(numStr)%2 != 0 {
			continue
		}

		mid := len(numStr) / 2
		firstHalf := numStr[:mid]
		secondHalf := numStr[mid:]

		if firstHalf == secondHalf {
			fmt.Printf("Found matching number: %d\n", i)
			sum += i
		}
	}

	return sum
}
