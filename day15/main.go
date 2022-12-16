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
	"sort"
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

type Range struct {
	l, r int64
}

func (r Range) String() string {
	return fmt.Sprintf("%d->%d", r.l, r.r)
}

// NOTE: as shown in example input, beacons and sensors can have negative coordinates
type Sensor struct {
	x, y                        int64
	beacon_x, beacon_y          int64 // data for closest beacon
	manhattan_to_closest_beacon int64
}

// returns the beginning and end (inclusive) of the values that this sensor covers on the given y value
func (s *Sensor) CoveredRange(y int64) (int64, int64) {
	// Calculate the range of the sensor's impossible-area
	// Area covered: 1 or 2 points satisfying equation:
	// |s.x - x| + |s.y - DESIRED_Y| <= s.md
	// |s.x - x| +  manhattan_y      <= s.md
	manhattan_y := s.y - y
	if manhattan_y < 0 {
		manhattan_y *= -1
	}

	// continuing to evaluate equation:
	// (x - s.x) + manhattan_y <= s.md OR (s.x - x) + manhattan_y <= s.md
	// x <= s.md - manhattan_y + s.x OR x >= s.x + manhattan_y - s.md
	greater_x := s.manhattan_to_closest_beacon - manhattan_y + s.x
	lesser_x := s.x + manhattan_y - s.manhattan_to_closest_beacon

	// if lesser_x > greater_x, the equation is not satisfied and this block does nothing
	return lesser_x, greater_x
}

func (s *Sensor) String() string {
	return fmt.Sprintf("sensor:(%d,%d) with closest beacon:(%d,%d) at distance %d", s.x, s.y, s.beacon_x, s.beacon_y, s.manhattan_to_closest_beacon)
}

func ParseInputToSensor(input string) (*Sensor, error) {
	s := new(Sensor)

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
	}

	// Part 1: counting the positions where a beacon cannot possibly be along just a single row
	beacon_impossible_on_desired_y := make(map[int64]bool)

	// Calculate the range of the sensor's impossible-area
	for _, s := range sensors {
		lesser_x, greater_x := s.CoveredRange(DESIRED_Y)

		// if lesser_x > greater_x, the equation is not satisfied and this block does nothing
		for x := lesser_x; x <= greater_x; x++ {
			beacon_impossible_on_desired_y[x] = true
		}
	}

	// Remove the actual beacons from the result
	// (do not remove actual sensors, because it is true that a cell containing a sensor cannot contain a beacon)
	for _, s := range sensors {
		if s.beacon_y != DESIRED_Y {
			continue
		}

		delete(beacon_impossible_on_desired_y, s.beacon_x)
	}

	fmt.Printf("\nPart 1 answer: %v\n", len(beacon_impossible_on_desired_y))

	// Part 2: Find the only possible position for the distress beacon
	// For each y value, go over all ranges covered by the sensors on that row to find any gaps of 1

	var y, x int64
Y_Loop:
	for y = 0; y <= X_Y_LIMIT; y++ {
		// collect all covered ranges of this y value
		var covered []Range
		for _, s := range sensors {
			lesser_x, greater_x := s.CoveredRange(y)
			if lesser_x <= greater_x {
				covered = append(covered, Range{l: lesser_x, r: greater_x})
			}
		}

		// sort said ranges so we can look from left to right
		sort.Slice(covered, func(i, j int) bool {
			if covered[i].l == covered[j].l {
				return covered[i].r < covered[j].r
			}
			return covered[i].l < covered[j].l
		})

		// determine if there's any gaps at the left edge of the grid
		if covered[0].l > 0 {
			x = 0
			break Y_Loop
		}

		// try to find the gap in the middle
		max_covered_x := covered[0].r
		for _, r := range covered[1:] {
			x = max_covered_x + 1
			if r.l > x {
				break Y_Loop
			}

			if r.r > max_covered_x {
				max_covered_x = r.r
			}
		}

		// try to find a gap at the right edge
		if max_covered_x < X_Y_LIMIT {
			x = X_Y_LIMIT
			break Y_Loop
		}
	}

	fmt.Printf("\nPart 2 answer: %v\n", x*TUNING_FREQUENCY_X_MULTIPLIER+y)
}
