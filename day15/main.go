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

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

var (
	// Problem parameters
	params = map[string][]int64{
		"input.txt":         {2000000, 4000000},
		"example_input.txt": {10, 20},
	}

	DESIRED_Y int64
	X_Y_LIMIT int64

	// Regex for parsing the problem input
	input_re = regexp.MustCompile(`^Sensor at x=(\-{0,1}\d+), y=(\-{0,1}\d+): closest beacon is at x=(\-{0,1}\d+), y=(\-{0,1}\d+)$`)
)

const (
	TUNING_FREQUENCY_X_MULTIPLIER = 4000000
)

// NOTE: as shown in example input, beacons and sensors can have negative coordinates
type Sensor struct {
	x, y                        int64
	beacon_x, beacon_y          int64 // data for closest beacon
	manhattan_to_closest_beacon int64
}

func (s *Sensor) String() string {
	return fmt.Sprintf("sensor:(%d,%d) with closest beacon:(%d,%d) at distance %d", s.x, s.y, s.beacon_x, s.beacon_y, s.manhattan_to_closest_beacon)
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
	// Choose parameters for run
	input_file_name := "input.txt"
	DESIRED_Y = params[input_file_name][0]
	X_Y_LIMIT = params[input_file_name][1]

	// Parse input as sensor data
	input_lines, err := fileutil.GetLinesFromFile(input_file_name)
	if err != nil {
		panic(err)
	}

	sensors := make([]*Sensor, len(input_lines))
	for i, line := range input_lines {
		sensors[i], err = ParseInputToSensor(line)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("DEBUG: created sensor %+v\n", sensors[i])
	}

	// Part 1: counting the positions where a beacon cannot possibly be along just a single row

	// Store possible beacon coordinates in the bounded range as a possibilities list
	bounded_possible_beacons_y_x := make(map[int64]map[int64]struct{})
	var y int64
	for y = 0; y <= X_Y_LIMIT; y++ {
		bounded_possible_beacons_y_x[y] = make(map[int64]struct{})

		var x int64
		for x = 0; x <= X_Y_LIMIT; x++ {
			bounded_possible_beacons_y_x[y][x] = struct{}{}
		}
	}

	// account for values outside of range as well, for Part 1 line
	impossible_beacon_outside_limit_x_counter := make(map[int64]bool)

	// TODO NEXT move to member function of sensor

	// Calculate the range of the sensor's impossible-area on the desired Y line
	for i, s := range sensors {
		fmt.Printf("DEBUG: calculating area with sensor %d...\n", i)
		// Area covered: 1 or 2 points satisfying equation:
		// |s.x - x| + |s.y - DESIRED_Y| <= s.md
		// |s.x - x| +  manhattan_y      <= s.md
		manhattan_y := s.y - DESIRED_Y
		if manhattan_y < 0 {
			manhattan_y *= -1
		}

		// continuing to evaluate equation:
		// (x - s.x) + manhattan_y <= s.md OR (s.x - x) + manhattan_y <= s.md
		// x <= s.md - manhattan_y + s.x OR x >= s.x + manhattan_y - s.md
		greater_x := s.manhattan_to_closest_beacon - manhattan_y + s.x
		lesser_x := s.x + manhattan_y - s.manhattan_to_closest_beacon

		// fmt.Printf("DEBUG: sensor at %d,%d covers the desired line from x=%d to x=%d\n", s.x, s.y, lesser_x, greater_x)

		if _, y_in_map := bounded_possible_beacons_y_x[DESIRED_Y]; !y_in_map {
			bounded_possible_beacons_y_x[DESIRED_Y] = make(map[int64]struct{})
		}

		// if lesser_x > greater_x, the equation is not satisfied and this block does nothing
		for x := lesser_x; x <= greater_x; x++ {
			// if the range extends outside of our limit, track it in a part-1-only counter
			if x < 0 || x > X_Y_LIMIT {
				impossible_beacon_outside_limit_x_counter[x] = true
			} else {
				delete(bounded_possible_beacons_y_x[DESIRED_Y], x)
			}
		}
	}

	// Remove the actual beacons from the result
	// (do not remove actual sensors, because it is true that a cell containing a sensor cannot contain a beacon)
	for _, s := range sensors {
		_, answer_contains_y := bounded_possible_beacons_y_x[s.beacon_y]
		if answer_contains_y {
			// fmt.Printf("DEBUG: size of bounded_possible_beacons_y_x[%d] BEFORE deletion is %v\n", s.beacon_y, len(bounded_possible_beacons_y_x[s.beacon_y]))

			delete(bounded_possible_beacons_y_x[s.beacon_y], s.beacon_x)
			delete(impossible_beacon_outside_limit_x_counter, s.beacon_x)

			// fmt.Printf("DEBUG: size of bounded_possible_beacons_y_x[%d] AFTER deletion is %v\n", s.beacon_y, len(bounded_possible_beacons_y_x[s.beacon_y]))
		}
	}

	// IDEA: brute-force: could store the entire map and mark off each impossible area - takes O(s*M^2) time where s is the number of sensors and M is the manhattan distance of each paired sensor-beacon
	// IDEA: if we just take the desired y value, and try each coordinate against each of the paired sensor-beacons, that's O(COLS*s) where COLS is the max x coordinate from the input

	fmt.Printf("\nPart 1 answer: %v\n", X_Y_LIMIT+1-int64(len(bounded_possible_beacons_y_x[DESIRED_Y]))+int64(len(impossible_beacon_outside_limit_x_counter)))

	// NOTE: I feel it's very likely that Part 2 is going to be "how many positions in the entire grid cannot contain a beacon"

	// TODO consider having the maps contain the full grid initially, then get erased over time
	// TODO consider using big.Int for the tuning frequency, if my initial number goes negative, but since int64 can store almost all 19 digits, 4*10^6 squared should fit within that
}
