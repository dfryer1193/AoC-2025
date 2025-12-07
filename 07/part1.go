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

	beamLocs := make([]map[int]struct{}, 0)
	splits := 0

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if len(beamLocs) == 0 {
			beamLocs = append(beamLocs, make(map[int]struct{}))
			startLoc := strings.Index(line, "S")
			beamLocs[0][startLoc] = struct{}{}
			continue
		}

		beamLocs = append(beamLocs, make(map[int]struct{}))
		for idx := range beamLocs[i-1] {
			if line[idx] == '.' {
				beamLocs[i][idx] = struct{}{}
			}

			if line[idx] == '^' {
				splits++
				if idx > 0 && line[idx-1] == '.' {
					beamLocs[i][idx-1] = struct{}{}
				}

				if idx < len(line)-1 && line[idx+1] == '.' {
					beamLocs[i][idx+1] = struct{}{}
				}
			}
		}
	}

	fmt.Println("Number of splits:", splits)
}
