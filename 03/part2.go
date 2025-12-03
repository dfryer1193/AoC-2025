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

	sum := uint64(0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		peak := peakJoltage(line)
		fmt.Printf("Peak joltage for %s: %d\n", line, peak)
		sum += peak
	}

	fmt.Println("Total sum:", sum)
}

// peakJoltage finds the maximum twelve digit number that can be formed from the input string.
func peakJoltage(n string) uint64 {
	accumulator := uint64(0)
	peakIdx := 0

	for i := 12; i >= 0; i-- {
		peak := '0'
		for j, c := range n[peakIdx : len(n)-i] {
			if c > peak {
				peak = c
				j = peakIdx // This is wrong
			}
		}
	}
}
