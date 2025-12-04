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
	lpeakIdx := 0

	for i := 11; i >= 0; i-- {
		lmax := byte('0')
		for j := lpeakIdx; j <= len(n)-i-1; j++ {
			if n[j] > lmax {
				lmax = n[j]
				lpeakIdx = j + 1
			}
		}
		val, _ := strconv.Atoi(string(lmax))
		accumulator = accumulator*10 + uint64(val)
	}

	return accumulator
}
