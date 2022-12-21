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
	for row := len(tiles) - 1; row >= 0; row-- {
		for col := 0; col < COLS; col++ {
			if tiles[row][col] == ROCK {
				return row
			}
		}
	}

	return -1
}

func main() {
	// Parse input
	input_lines, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	jet_pattern := input_lines[0]

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

	// How many units tall will the tower of rocks be after 2022 rocks have stopped falling?
	for i := 1; i <= 2022; i, rock_type_i = i+1, (rock_type_i+1)%len(ROCKS) {
		// Determine starting position of lower-left corner of rock
		rock_row := get_tallest_rock_row(tiles) + ROCK_START_ROW_BUFFER + 1
		rock_col := ROCK_START_COL_BUFFER

		// fmt.Printf("DEBUG: rock %d starts falling at row:%d,col:%d\n", i, rock_row, rock_col)

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

			// if fall fails, store in tiles and move onto next rock
			if fall_impeded {
				// fmt.Printf("DEBUG: rock %d comes to a rest at row:%d,col:%d\n", i, rock_row, rock_col)

				// allocate new rows in tiles if needed
				for rock_row+len(ROCKS[rock_type_i])-1 >= len(tiles) {
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

				// handle new rock
				break
			} else {
				rock_row -= 1
				// fmt.Printf("DEBUG: rock %d is now at row:%d,col:%d\n", i, rock_row, rock_col)
			}
		}
	}

	/*
		IDEAS:
		- "fall" upward (to smaller row values)
		- answer will be the value I need to save for the rock falls anyway: the highest rock in the room

		Algorithms and time/space used:
		- store all rows and columns, adding new rows as needed
			- O(7*2022) space
		- for each column value, store the highest row that is occupied by a rock/floor
			- O(7) space
	*/

	// Determine final tallest rock
	fmt.Printf("\nPart 1 answer: %v\n", get_tallest_rock_row(tiles)+1)
}
