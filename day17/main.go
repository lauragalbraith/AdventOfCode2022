/*
Package main solves Day 17 of Advent of Code 2022
main.go: Laura Galbraith
Pyroclastic Flow
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"fmt"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

// Problem constants
const (
	// The tall, vertical chamber is exactly seven units wide
	COLS = 7

	// Each rock appears so that its left edge is two units away from the left wall...
	ROCK_START_COL_BUFFER = 2
	// ...and its bottom edge is three units above the highest rock in the room (or the floor, if there isn't one)
	ROCK_START_ROW_BUFFER = 3

	ROCK = '#'
	AIRE = '.'
)

// From problem description
var (
	ROCKS = [][][]byte{
		{
			{ROCK, ROCK, ROCK, ROCK},
		},
		{
			{AIRE, ROCK, AIRE},
			{ROCK, ROCK, ROCK},
			{AIRE, ROCK, AIRE},
		},
		{
			{AIRE, AIRE, ROCK},
			{AIRE, AIRE, ROCK},
			{ROCK, ROCK, ROCK},
		},
		{
			{ROCK},
			{ROCK},
			{ROCK},
			{ROCK},
		},
		{
			{ROCK, ROCK},
			{ROCK, ROCK},
		},
	}

	// specifies the columnar difference when a rock is pushed by a jet
	JET_DIRECTIONS = map[byte]int{'>': 1, '<': -1}
)

// Constants for hashing chamber state
const (
	ROWS_TO_HASH = 5 // inspired by solution from Reddit: https://github.com/vss2sn/advent_of_code/blob/master/2022/cpp/day_17b.cpp
)

// returns the row, column of where a single tile of the rock is
func get_rock_tile_indeces(
	sw_corner_row, sw_corner_col int,
	rock_type int,
	row_in_rock, col_in_rock int) (int, int) {

	tile_row := sw_corner_row + len(ROCKS[rock_type]) - 1 - row_in_rock
	tile_col := sw_corner_col + col_in_rock

	return tile_row, tile_col
}

// returns -1 if there are no rocks
func get_tallest_rock_row(tiles [][]byte) int {
	// look at the topmost row, continuing down until we see a rock
	for tiles_row := len(tiles) - 1; tiles_row >= 0; tiles_row-- {
		for col := 0; col < COLS; col++ {
			if tiles[tiles_row][col] == ROCK {
				return tiles_row
			}
		}
	}

	return -1
}

// in Go, I'm not sure there's a way to do native hashing for a type (like size_t operator() in C++); could just create an INT out of hashing together the occupied tiles in the top R rows
func hash_chamber_top(tiles [][]byte) uint64 {
	var hash uint64 // Go initializes to 0

	for i := 0; i < ROWS_TO_HASH; i++ {
		row := get_tallest_rock_row(tiles) - i
		if get_tallest_rock_row(tiles) < ROWS_TO_HASH-1 {
			row = ROWS_TO_HASH - 1 - i
		}
		for col := 0; col < COLS; col++ {
			var val uint64 // Go initializes to 0
			if row < len(tiles) && tiles[row][col] == ROCK {
				val = 1
			}

			// it's a row of bytes uniquely representing the state; NOTE: only works if ROWS_TO_HASH*COLS <= 64
			hash = (hash << 1) + val
		}
	}

	return hash
}

// returns the height of the tallest rock tower after x many rocks have fallen
func TallestTowerHeightAfterXFalls(
	x int64,
	jet_pattern string) int64 {

	// Initialize falling state
	jet_i := 0
	rock_type_i := 0

	// Store all relevant tiles: a tower could have gauges in its side that perfectly fit a rock being pushed into it, but the tallest rock at all columns wouldn't be able to tell you that
	tiles := make([][]byte, ROCK_START_ROW_BUFFER)
	for row := range tiles {
		tiles[row] = make([]byte, COLS)
		for col := range tiles[row] {
			tiles[row][col] = AIRE
		}
	}

	// Store hashed state of the upper (most relevant) chamber to detect repeated states (stable cycle)
	seen_jet_type_chamber_heights := make(map[int]map[int]map[uint64][]int64)
	var check_cycle_after_rocks int64
	check_cycle_after_rocks = -1

	// Simulate X rock falls
	var i int64
	for i = 1; i <= x; i, rock_type_i = i+1, (rock_type_i+1)%len(ROCKS) {
		// Determine starting position of lower-left corner of rock
		rock_row := get_tallest_rock_row(tiles) + ROCK_START_ROW_BUFFER + 1
		rock_col := ROCK_START_COL_BUFFER

		// Get pushed by jet, then fall
		// Operation: pushed by jet (which may fail), then falling one unit down (if this fails, go to next rock)
		for true {
			// determine if jet can push the rock over by one space
			push_dir := JET_DIRECTIONS[jet_pattern[jet_i]]
			jet_i = (jet_i + 1) % len(jet_pattern)

			push_impeded := false
			for i := range ROCKS[rock_type_i] {
				for j, tile := range ROCKS[rock_type_i][i] {
					if tile == AIRE {
						continue
					}

					tile_row, tile_col := get_rock_tile_indeces(rock_row, rock_col, rock_type_i, i, j)

					// check if it hits a wall
					new_tile_col := tile_col + push_dir
					if new_tile_col < 0 || new_tile_col >= COLS {
						push_impeded = true
						break
					}

					// check if it is blocked by another rock
					// if the row isn't even accounted for in the tiles, it cannot be a rock
					if tile_row < len(tiles) && tiles[tile_row][new_tile_col] == ROCK {
						push_impeded = true
						break
					}
				}
			}

			// take effect of jet
			if !push_impeded {
				rock_col += push_dir
			}

			// determine if fall can occur
			fall_impeded := false
			for i := range ROCKS[rock_type_i] {
				for j, tile := range ROCKS[rock_type_i][i] {
					if tile == AIRE {
						continue
					}

					tile_row, tile_col := get_rock_tile_indeces(rock_row, rock_col, rock_type_i, i, j)

					new_tile_row := tile_row - 1

					// check if it hits the floor
					if new_tile_row < 0 {
						fall_impeded = true
						break
					}

					// check if it is blocked by another rock
					// if the row isn't even accounted for in the tiles, it cannot be a rock
					if new_tile_row < len(tiles) && tiles[new_tile_row][tile_col] == ROCK {
						fall_impeded = true
						break
					}
				}
			}

			if !fall_impeded {
				rock_row -= 1
				continue
			}

			// if fall fails, note state and move onto next rock

			// allocate new rows in tiles if needed
			for int(rock_row)+len(ROCKS[rock_type_i])-1 >= len(tiles) {
				tiles = append(tiles, make([]byte, COLS))
				for col := 0; col < COLS; col++ {
					tiles[len(tiles)-1][col] = AIRE
				}
			}

			// store all of this rock's tiles in tiles
			for i := range ROCKS[rock_type_i] {
				for j, tile := range ROCKS[rock_type_i][i] {
					if tile == AIRE {
						continue
					}

					tile_row, tile_col := get_rock_tile_indeces(rock_row, rock_col, rock_type_i, i, j)
					tiles[tile_row][tile_col] = ROCK
				}
			}

			// Store hashed state of the upper (most relevant) chamber

			// technically this is the jet to affect the next state, but it should be same difference
			if _, jet_seen := seen_jet_type_chamber_heights[jet_i]; !jet_seen {
				seen_jet_type_chamber_heights[jet_i] = make(map[int]map[uint64][]int64)
			}

			if _, rock_type_seen := seen_jet_type_chamber_heights[jet_i][rock_type_i]; !rock_type_seen {
				seen_jet_type_chamber_heights[jet_i][rock_type_i] = make(map[uint64][]int64)
			}

			// Check if hashed state has been seen before, where we can immediately calculate and return result
			chamber_hash := hash_chamber_top(tiles)
			current_tallest_rock_row := int64(get_tallest_rock_row(tiles))
			previous_state, state_seen := seen_jet_type_chamber_heights[jet_i][rock_type_i][chamber_hash]

			// for some reason, full input incorrectly detects a cycle early, so delay it to double-check the cycle theory
			if state_seen && check_cycle_after_rocks < 0 {
				cycle_length := i - previous_state[1]
				check_cycle_after_rocks = i + cycle_length
			} else if state_seen && i > check_cycle_after_rocks {
				// Extrapolate pattern at x rocks by this result
				cycle_length := i - previous_state[1]
				complete_cycles_to_go := (x - i) / cycle_length
				growth_during_cycle := current_tallest_rock_row - previous_state[0]

				partial_cycle_size_to_go := (x - i) % cycle_length

				// find how much growth occurs during this partial cycle
				for _, jet_info := range seen_jet_type_chamber_heights {
					for _, type_info := range jet_info {
						for _, state_info := range type_info {
							if state_info[1] == previous_state[1]+partial_cycle_size_to_go {
								growth_during_partial_cycle := state_info[0] - previous_state[0]

								// The height of the tower will be 1 greater than the row of the extrapolated tallest rock
								return current_tallest_rock_row + (complete_cycles_to_go * growth_during_cycle) + growth_during_partial_cycle + 1
							}
						}
					}
				}
			}

			// save this state
			seen_jet_type_chamber_heights[jet_i][rock_type_i][chamber_hash] = []int64{current_tallest_rock_row, i}

			// Handle new rock
			break
		}
	}

	// Determine final tallest rock if we never encountered a cycle
	return int64(get_tallest_rock_row(tiles) + 1)
}

func main() {
	// Parse input
	input_lines, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	// How many units tall will the tower of rocks be after 2022 rocks have stopped falling?
	answer := TallestTowerHeightAfterXFalls(2022, input_lines[0])
	fmt.Printf("\nPart 1 answer: %v\n", answer)

	// How tall will the tower be after 1000000000000 rocks have stopped?
	answer = TallestTowerHeightAfterXFalls(1000000000000, input_lines[0])
	fmt.Printf("\nPart 2 answer: %v\n", answer)
}
