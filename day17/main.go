/*
Package main solves Day 17 of Advent of Code 2022
main.go: Laura Galbraith
TODO title
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"fmt"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

func main() {
	input_lines, err := fileutil.GetLinesFromFile("example_input.txt")
	if err != nil {
		panic(err)
	}

	// TODO solve day 17

	fmt.Printf("\nPart 1 answer: %v\n", len(input_lines))
}
