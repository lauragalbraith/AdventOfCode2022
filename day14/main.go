/*
Package main solves Day 14 of Advent of Code 2022
main.go: Laura Galbraith
Regolith Reservoir
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

const (
	// provided in problem specification
	SAND_SOURCE_ROW = 0
	SAND_SOURCE_COL = 500
)

var (
	coordinate_re = regexp.MustCompile(`(\d+),(\d+)`)

	MAX_ROW = 0
	MAX_COL = 0
)

const (
	Air = iota
	Rock
	Sand
)

func PrintGrid(grid [][]int) {
	fmt.Println("DEBUG: this is the grid:")
	for _, row := range grid {
		for _, tile := range row {
			fmt.Printf("%c", (".#o")[tile])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	cave_input, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	// Determine size of grid
	for _, line := range cave_input {
		coordinates_found := coordinate_re.FindAllStringSubmatch(line, -1)
		// fmt.Printf("DEBUG: %d coordinates found in line %d\n", len(coordinates_found), i)

		// NOTE: x = col (distance right), y = row (distance down)
		for _, coord := range coordinates_found {
			if len(coord) < 3 {
				panic(fmt.Errorf("Unexpected input line format: '%s'", line))
			}

			x, err := strconv.Atoi(coord[1])
			if err != nil {
				panic(err)
			}

			if x > MAX_COL {
				MAX_COL = x
			}

			y, err := strconv.Atoi(coord[2])
			if err != nil {
				panic(err)
			}

			if y > MAX_ROW {
				MAX_ROW = y
			}

			// fmt.Printf("DEBUG: coordinate is x=%d, y=%d\n", x, y)
		}
	}

	// Initialize grid
	// NOTE: both inputs do not have any lines that jut up near the left or top of the cave (0), so we shouldn't have to worry about emulating negative indeces
	grid := make([][]int, MAX_ROW+2)
	for row := range grid {
		grid[row] = make([]int, MAX_COL+2)
		// each tile is set to int's default value: 0, which is Air
	}

	// Fill in rock lines in the grid
	// NOTE: x = col (distance right), y = row (distance down)
	for _, line := range cave_input {
		coordinates_found := coordinate_re.FindAllStringSubmatch(line, -1)

		var rock_line_start_row, rock_line_start_col int
		for coord_i, coord := range coordinates_found {
			if len(coord) < 3 {
				panic(fmt.Errorf("Unexpected input line format: '%s'", line))
			}

			x, err := strconv.Atoi(coord[1])
			if err != nil {
				panic(err)
			}

			y, err := strconv.Atoi(coord[2])
			if err != nil {
				panic(err)
			}

			// draw rock line if its start point is known
			if coord_i > 0 {
				// lines can point in any direction
				row_start := rock_line_start_row
				row_end := y
				if rock_line_start_row > y {
					row_start = y
					row_end = rock_line_start_row
				}

				col_start := rock_line_start_col
				col_end := x
				if rock_line_start_col > x {
					col_start = x
					col_end = rock_line_start_col
				}

				// draw the line, whether vertical or horizontal
				for row := row_start; row <= row_end; row++ {
					for col := col_start; col <= col_end; col++ {
						grid[row][col] = Rock
					}
				}
			}

			// save end point of this line as start point of next
			rock_line_start_row = y
			rock_line_start_col = x
		}
	}

	// Emulate sand: flows one unit (cell/tile) at a time, comes to rest, and then the next sand is produced
	still_sand := 0
	for ; ; still_sand++ {
		if still_sand%10 == 0 {
			fmt.Printf("DEBUG: %d pieces of still sand\n", still_sand)
			PrintGrid(grid)
		}

		sand_row := SAND_SOURCE_ROW
		sand_col := SAND_SOURCE_COL
		for sand_row <= MAX_ROW && sand_col <= MAX_COL {
			// Try to fall down, then diagonally 1 down & 1 left, then diagonally 1 down & 1 right
			if grid[sand_row+1][sand_col] == Air {
				sand_row++
			} else if grid[sand_row+1][sand_col-1] == Air {
				sand_row++
				sand_col--
			} else if grid[sand_row+1][sand_col+1] == Air {
				sand_row++
				sand_col++
			} else {
				// if there's no place to fall, come to rest
				break
			}
		}

		// stop processing sand when a sand coordinate reaches the lower or right grid edge, where we know we have no rock lines to stop it
		if sand_row > MAX_ROW || sand_col > MAX_COL {
			break
		}

		// Note where sand landed still
		grid[sand_row][sand_col] = Sand
	}

	// How many units of sand come to rest before sand starts flowing into the abyss below?
	fmt.Printf("Part 1 answer: %v\n", still_sand)
}
