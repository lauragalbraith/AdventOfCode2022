/*
Package fileutil provides Go utility functions for interacting with Advent of Code files
fileutil.go: Laura Galbraith
*/
package fileutil

import (
	"bufio"
	"fmt"
	"os"
)

func GetLinesFromFile(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", name, err)
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	if err := s.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", name, err)
		return nil, err
	}

	return lines, nil
}
