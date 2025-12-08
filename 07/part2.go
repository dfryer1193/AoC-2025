package main

import (
	"bufio"
	"fmt"
	"os"
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

	// beamLocs keeps track of possible beam locations at each row and how many timelines led to that location
	beamLocs := make([]map[int]int, 0)

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if len(beamLocs) == 0 {
			beamLocs = append(beamLocs, make(map[int]int))
			startLoc := strings.Index(line, "S")
			beamLocs[0][startLoc] = 1
			continue
		}

		beamLocs = append(beamLocs, make(map[int]int))
		for idx := range beamLocs[i-1] {
			if line[idx] == '.' {
				beamLocs[i][idx] += beamLocs[i-1][idx]
			}

			if line[idx] == '^' {
				if idx > 0 && line[idx-1] == '.' {
					beamLocs[i][idx-1] += beamLocs[i-1][idx]
				}

				if idx < len(line)-1 && line[idx+1] == '.' {
					beamLocs[i][idx+1] += beamLocs[i-1][idx]
				}
			}
		}
	}

	timelines := 0
	for _, count := range beamLocs[len(beamLocs)-1] {
		timelines += count
	}

	fmt.Println("Number of timelines:", timelines)
}
