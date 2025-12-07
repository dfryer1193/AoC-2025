package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type operation int

const (
	add operation = iota
	multiply
)

type equation struct {
	values []int
	op     operation
}

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

	accumulator := 0
	eqs := make([]*equation, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		for i, field := range strings.Fields(line) {
			if i > len(eqs)-1 {
				eqs = append(eqs, &equation{
					values: make([]int, 0),
					op:     add,
				})
			}

			eq := eqs[i]

			val, err := strconv.Atoi(field)
			if err != nil {
				if field == "+" {
					eq.op = add
				} else if field == "*" {
					eq.op = multiply
				} else {
					fmt.Println("Unknown operator:", field)
					return
				}
				continue
			}

			eq.values = append(eq.values, val)
		}
	}

	for _, eq := range eqs {
		result := 0
		if len(eq.values) == 0 {
			continue
		}

		for _, val := range eq.values {
			if eq.op == add {
				result += val
			} else if eq.op == multiply {
				if result == 0 {
					result = 1
				}
				result *= val
			}
		}
		accumulator += result
	}

	fmt.Println("Final accumulator value:", accumulator)
}
