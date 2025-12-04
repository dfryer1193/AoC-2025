package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
		peak := peakJoltage(line)
		fmt.Printf("Peak joltage for %s: %d\n", line, peak)
		sum += peak
	}

	fmt.Println("Total sum:", sum)
}

// peakJoltage finds the maximum two digit numbe that can be formed from the input string.
func peakJoltage(n string) int {
	peak := '0'
	peakIdx := -1

	// Scan the input, omitting the last character, since the output must be two digits.
	// The first number will be the tens place, and the second number will be the ones place, so we just want to use the peak as the tens place.
	for i, c := range n[:len(n)-1] {
		if c > peak {
			peak = c
			peakIdx = i
		}
	}

	nextPeak := '0'
	// Now scan the input again, starting from the character after the peak found above.
	for _, c := range n[peakIdx+1:] {
		if c > nextPeak {
			nextPeak = c
		}
	}

	tens, _ := strconv.Atoi(string(peak))
	ones, _ := strconv.Atoi(string(nextPeak))

	return (tens * 10) + ones
}
