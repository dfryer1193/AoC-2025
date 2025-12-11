package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type machine struct {
	reqMask  int
	buttons  []int
	joltages []int
}

// Try to find the shortest sequence of button presses that lights up the required lights.
// The machine starts with all lights off. Each button toggles certain lights on/off.
func (m *machine) Start() int {
	if m.reqMask == 0 {
		return 0
	}

	level := []int{0}
	visited := map[int]bool{0: true}
	presses := 0

	for len(level) > 0 {
		presses++
		nextLevel := []int{}
		for _, mask := range level {
			for _, button := range m.buttons {
				nextMask := mask ^ button
				if nextMask == m.reqMask {
					return presses
				}
				if !visited[nextMask] {
					visited[nextMask] = true
					nextLevel = append(nextLevel, nextMask)
				}
			}
		}
		level = nextLevel
	}
	return -1
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: go run part1.go <filename>")
		return
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	machines := make([]*machine, 0)
	minPresses := 0

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		machine := &machine{}

		if len(parts) > 0 {
			machine.reqMask = parseLights(parts[0])
		}

		for i := 1; i < len(parts); i++ {
			part := parts[i]
			if strings.HasPrefix(part, "(") {
				machine.buttons = append(machine.buttons, parseButtons(part))
			} else if strings.HasPrefix(part, "{") {
				machine.joltages = parseJoltages(part)
			}
		}

		machines = append(machines, machine)
		minPresses += machine.Start()
	}

	fmt.Println("Total minimum button presses for all machines:", minPresses)
}

func parseLights(state string) int {
	lights := 0
	lightStr := state[1 : len(state)-1]
	for i, ch := range lightStr {
		if ch == '#' {
			lights |= (1 << uint(i))
		}
	}
	return lights
}

func parseButtons(state string) int {
	buttonMask := 0
	nums := strings.Split(state[1:len(state)-1], ",")
	for _, num := range nums {
		val, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Error parsing button value:", err)
			continue
		}
		buttonMask |= (1 << val)
	}
	return buttonMask
}

func parseJoltages(state string) []int {
	joltages := make([]int, 0)
	nums := strings.Split(state[1:len(state)-1], ",")
	for _, num := range nums {
		val, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Error parsing joltage value:", err)
			continue
		}
		joltages = append(joltages, val)
	}
	return joltages
}

