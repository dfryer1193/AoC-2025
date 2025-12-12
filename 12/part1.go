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

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
	}

	fmt.Println("Herp derp no impl yet lol")
}
