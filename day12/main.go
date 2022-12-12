/*
Package main solves Day 12 of Advent of Code 2022
main.go: Laura Galbraith
What is the fewest steps required to move from your current position to the location that should get the best signal?
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"fmt"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil" // GetLinesFromFile
)

func main() {
	// Get input
	input, err := fileutil.GetLinesFromFile("example_input.txt")
	if err != nil {
		panic(err)
	}

	for _, i := range input {
		fmt.Println(i)
	}

	fmt.Println("Part 1 answer: TODO")
	// TODO IDEA: Dijkstra's (from Reddit comments)
}
