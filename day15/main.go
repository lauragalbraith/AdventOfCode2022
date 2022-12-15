/*
Package main solves Day 15 of Advent of Code 2022
main.go: Laura Galbraith
Beacon Exclusion Zone
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

var (
	input_re = regexp.MustCompile(`^Sensor at x=(\-{0,1}\d+), y=(\-{0,1}\d+): closest beacon is at x=(\-{0,1}\d+), y=(\-{0,1}\d+)$`)
)

// NOTE: as shown in example input, beacons and sensors can have negative coordinates
type Sensor struct {
	x, y                        int64
	beacon_x, beacon_y          int64 // data for closest beacon
	manhattan_to_closest_beacon int64
}

func ParseInputToSensor(input string) (*Sensor, error) {
	s := new(Sensor)

	// fmt.Printf("DEBUG: line: *%s*\n", input)
	matches := input_re.FindAllStringSubmatch(input, -1)
	if len(matches) < 1 || len(matches[0]) < 5 {
		return nil, fmt.Errorf("unexpected input format: '%s'", input)
	}

	// extract x,y coordinates of the sensor and its closest beacon
	var err error
	s.x, err = strconv.ParseInt(matches[0][1], 10, 64)
	if err != nil {
		return nil, err
	}

	s.y, err = strconv.ParseInt(matches[0][2], 10, 64)
	if err != nil {
		return nil, err
	}

	s.beacon_x, err = strconv.ParseInt(matches[0][3], 10, 64)
	if err != nil {
		return nil, err
	}

	s.beacon_y, err = strconv.ParseInt(matches[0][4], 10, 64)
	if err != nil {
		return nil, err
	}

	// calculate manhattan distance
	manhattan_x := s.x - s.beacon_x
	if manhattan_x < 0 {
		manhattan_x *= -1
	}

	manhattan_y := s.y - s.beacon_y
	if manhattan_y < 0 {
		manhattan_y *= -1
	}

	s.manhattan_to_closest_beacon = manhattan_x + manhattan_y

	return s, nil
}

func main() {
	input_file_name := "example_input.txt"
	input_lines, err := fileutil.GetLinesFromFile(input_file_name)
	if err != nil {
		panic(err)
	}

	// Parse input as sensor data
	sensors := make([]*Sensor, len(input_lines))
	for i, line := range input_lines {
		sensors[i], err = ParseInputToSensor(line)
		if err != nil {
			panic(err)
		}
	}

	var DESIRED_Y int64
	DESIRED_Y = 10 // 10 for example input, 2000000 for official input
	if strings.Compare(input_file_name, "input.txt") == 0 {
		DESIRED_Y = 2000000
	}

	// Part 1: counting the positions where a beacon cannot possibly be along just a single row

	// store impossible-beacon spots in a map of y coordinate to maps of x coordinates (to placeholder bools), then count up sizes of maps in O(1); easily expanded to full grid, and can handle negative coordinates
	no_beacon_spots_y_x := make(map[int64]map[int64]bool)

	// IDEA: if while processing the initial input, we determine if any y=desired point falls in the sensor's impossible area, then we know there's an odd number of applicable cells on y=desired point for this sensor; start at y=... and x=s[x] and then move right to see where the impossible area cuts off (or there's probably a mathematical way to calculate it) - O(s*M)

	// TODO make sure to subtract any cells from calculated areas that actually contain an input beacon

	// IDEA: brute-force: could store the entire map and mark off each impossible area - takes O(s*M^2) time where s is the number of sensors and M is the manhattan distance of each paired sensor-beacon
	// IDEA: if we just take the desired y value, and try each coordinate against each of the paired sensor-beacons, that's O(COLS*s) where COLS is the max x coordinate from the input

	fmt.Printf("\nPart 1 answer: %v\n", len(no_beacon_spots_y_x[DESIRED_Y]))

	// NOTE: I feel it's very likely that Part 2 is going to be "how many positions in the entire grid cannot contain a beacon"
}
