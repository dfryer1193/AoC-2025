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

		if direction == 'L' {
			steps *= -1
		}

		newPosition := position + steps

		zeroPasses := 0
		if position < 0 {
			zeroPasses = (position / 100) * -1
		} else if position > 99 {
			zeroPasses = position / 100
		}

		position %= 100
		if position < 0 {
			position += 100
		}

		zeroCount += zeroPasses
		if position == 0 {
			zeroCount++
		}
	}

	fmt.Println(zeroCount)
}
