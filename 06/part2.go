package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type operation int

const (
	add operation = iota
	multiply
)

type equation struct {
	values   []int
	op       operation
	rawVals  []string
	numCount int
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
	rawLines := make([]string, 0)
	eqs := make([]*equation, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		rawLines = append(rawLines, line)
		fields := strings.Fields(line)
		for i, field := range fields {
			if i > len(eqs)-1 {
				eqs = append(eqs, &equation{
					values:   make([]int, 0),
					op:       add,
					rawVals:  make([]string, 0),
					numCount: 0,
				})
			}

			eq := eqs[i]

			if field == "+" {
				eq.op = add
			} else if field == "*" {
				eq.op = multiply
			} else {
				if len(field) > eq.numCount {
					eq.numCount = len(field)
				}
				eq.rawVals = append(eq.rawVals, field)
			}
		}
	}

	for i := len(eqs) - 1; i >= 0; i-- {
		eq := eqs[i]
		lsum := 0
		vnums := make([]string, 0)
		eq.values = make([]int, eq.numCount)
		for j, l := range rawLines {
			vnums = append(vnums, l[len(l)-eq.numCount:])
			rawLines[j] = l[:len(l)-eq.numCount] // remove processed part
			if len(rawLines[j]) > 0 {
				rawLines[j] = rawLines[j][:len(rawLines[j])-1] // remove trailing space
			}
		}

		for j := eq.numCount - 1; j >= 0; j-- {
			for _, vnum := range vnums {
				if vnum[j] == ' ' {
					continue
				}

				if vnum[j] == '+' {
					eq.op = add
				} else if vnum[j] == '*' {
					eq.op = multiply
				} else {
					digit := int(vnum[j] - '0')
					eq.values[j] *= 10
					eq.values[j] += digit
				}
			}
		}

		for _, val := range eq.values {
			if eq.op == add {
				lsum += val
			} else if eq.op == multiply {
				if lsum == 0 {
					lsum = 1
				}
				lsum *= val
			}
		}
		accumulator += lsum
	}

	fmt.Println("Final accumulator value:", accumulator)
}
