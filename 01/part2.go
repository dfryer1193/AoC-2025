package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	zeroCount := 0
	position := 50

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		direction := line[0] // L or R - L subtracts, R adds
		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			log.Fatal(err)
		}

		var newPos int
		var count int
		if direction == 'R' {
			newPos, count = rotateRight(position, steps)
		} else if direction == 'L' {
			newPos, count = rotateLeft(position, steps)
		}
		zeroCount += count
		position = newPos
	}

	fmt.Println(zeroCount)
}

func rotateRight(position, steps int) (int, int) {
	newPos := position + steps
	fullRotations := steps / 100
	remainder := steps % 100

	if position+remainder > 99 {
		fullRotations++
	}

	return newPos % 100, fullRotations
}

func rotateLeft(position, steps int) (int, int) {
	newPos := position - steps
	fullRotations := steps / 100
	remainder := steps % 100

	if remainder >= position && position != 0 {
		fullRotations++
	}

	return ((newPos % 100) + 100) % 100, fullRotations
}
