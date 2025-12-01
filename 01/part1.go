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

		fmt.Println(string(direction), steps)

		if direction == 'L' {
			steps *= -1
		}

		position += steps

		fmt.Println("New position:", position)

		//if position < 0 {
		//	position = 100 + (position % 100)
		//} else if position > 99 {
		//	position = position % 100
		//}

		position %= 100
		if position < 0 {
			position += 100
		}

		fmt.Println("Wrapped position:", position)

		if position == 0 {
			zeroCount++
		}
	}

	fmt.Println(zeroCount)
}
