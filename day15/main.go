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
	r_top_x, r_top_y            int64
	r_bottom_x, r_bottom_y      int64
	r_left_x, r_left_y          int64
	r_right_x, r_right_y        int64
}

// Calculate the range of the sensor's impossible-area
func (s *Sensor) CalculateRhombusArea() { // TODO FINALLY remove this function and variables if I don't end up using them
	s.r_top_x, s.r_bottom_x = s.x, s.x
	s.r_left_y, s.r_right_y = s.y, s.y

	s.r_top_y = s.y - s.manhattan_to_closest_beacon
	s.r_bottom_y = s.y + s.manhattan_to_closest_beacon

	s.r_left_x = s.x - s.manhattan_to_closest_beacon
	s.r_right_x = s.x + s.manhattan_to_closest_beacon
}

func (s *Sensor) IsPointInRhombus(x, y int64) bool {
	// Area covered: |s.x - x| + |s.y - y| <= s.md
	manhattan_y := s.y - y
	if manhattan_y < 0 {
		manhattan_y *= -1
	}

	manhattan_x := s.x - x
	if manhattan_x < 0 {
		manhattan_x *= -1
	}

	return (manhattan_x + manhattan_y) <= s.manhattan_to_closest_beacon
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

	// calculate rhombus
	s.CalculateRhombusArea()

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
	// TODO remove this implementation that took too much memory
	/*bounded_possible_beacons_y_x := make(map[int64]map[int64]struct{})
	var y int64
	for y = 0; y <= X_Y_LIMIT; y++ {
		bounded_possible_beacons_y_x[y] = make(map[int64]struct{})

		var x int64
		for x = 0; x <= X_Y_LIMIT; x++ {
			bounded_possible_beacons_y_x[y][x] = struct{}{}
		}
	}*/

	// Store possible beacon coordinates in a packed bit list
	// TODO HERE this is still enough memory that my process is just getting killed - could I instead have the cost be in time? iterate over all options? iterate over one y at a time and if its number of ones ever reaches 64, continue on
	// save some time by having all sensors store a 4-point rhombus of their manhattan distance areas; then have an O(1) function check if a given point is in its rhombus
	/*known_y_x_coordinates := make([][]uint64, X_Y_LIMIT+1)
	for row, _ := range known_y_x_coordinates {
		// will initialize all uint64s to 0000...
		known_y_x_coordinates[row] = make([]uint64, X_Y_LIMIT/64+1)
	}*/

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

	// IDEA: brute-force: could store the entire map and mark off each impossible area - takes O(s*M^2) time where s is the number of sensors and M is the manhattan distance of each paired sensor-beacon
	// IDEA: if we just take the desired y value, and try each coordinate against each of the paired sensor-beacons, that's O(COLS*s) where COLS is the max x coordinate from the input

	fmt.Printf("\nPart 1 answer: %v\n", len(beacon_impossible_on_desired_y))

	// Part 2: Find the only possible position for the distress beacon
	// TODO consider having the maps contain the full grid initially, then get erased over time
	// TODO consider using big.Int for the tuning frequency, if my initial number goes negative, but since int64 can store almost all 19 digits, 4*10^6 squared should fit within that

	var x, y int64
Y_Loop:
	for y = 0; y <= X_Y_LIMIT; y++ {
		if y%1000 == 0 {
			fmt.Printf("DEBUG: processing y=%d\n", y)
		}
	X_Loop:
		for x = 0; x <= X_Y_LIMIT; x++ {
		Sensor_Loop:
			for _, s := range sensors {
				// skip sensor if it definitely doesn't intersect this point
				if y < s.r_top_y || y > s.r_bottom_y || x < s.r_left_x || x > s.r_right_x {
					continue Sensor_Loop
				}

				// skip coordinate if it's covered by this sensor
				if s.IsPointInRhombus(x, y) {
					continue X_Loop
				}
			}

			// break out of everything if we've found the answer
			break Y_Loop
		}
	}

	// TODO NEXT if this is too slow, try calling CoveredRange and filling in a per-row map of size 400000 to check if it ends up being size 400000 or if not what the missing number is
	fmt.Printf("\nPart 2 answer: %v\n", x*TUNING_FREQUENCY_X_MULTIPLIER+y)
}
