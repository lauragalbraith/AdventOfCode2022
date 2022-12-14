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
	MAX_ROW   = 0
	MAX_COL   = 0
	FLOOR_ROW int
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

// GetRockLineCoordinatesFromInput returns a list of rock lines, which are a list of turning points, which are a pair of row,col values
func GetRockLineCoordinatesFromInput(file_name string) ([][][]int, error) {
	cave_input, err := fileutil.GetLinesFromFile(file_name)
	if err != nil {
		return nil, err
	}

	coordinate_re := regexp.MustCompile(`(\d+),(\d+)`)

	ret := make([][][]int, len(cave_input))
	for line_i, line := range cave_input {
		coordinates_found := coordinate_re.FindAllStringSubmatch(line, -1)

		ret[line_i] = make([][]int, len(coordinates_found))
		for coord_i, coord := range coordinates_found {
			if len(coord) < 3 {
				return nil, fmt.Errorf("Unexpected input line format: '%s'", line)
			}

			// NOTE: x = col (distance right), y = row (distance down)
			x, err := strconv.Atoi(coord[1])
			if err != nil {
				return nil, err
			}

			y, err := strconv.Atoi(coord[2])
			if err != nil {
				return nil, err
			}

			ret[line_i][coord_i] = make([]int, 2)
			ret[line_i][coord_i][0] = y
			ret[line_i][coord_i][1] = x
		}
	}

	return ret, nil
}

// FillInRocks populates spaces that have rock lines, as specified by cave input
func FillInRocks(grid [][]int, rock_lines [][][]int) {
	for _, rock_line := range rock_lines {
		var rock_line_start_row, rock_line_start_col int

		for rock_point_i, rock_point := range rock_line {
			// draw rock line if its start point is known
			if rock_point_i > 0 {
				// lines can point in any direction
				row_start := rock_line_start_row
				row_end := rock_point[0]
				if rock_line_start_row > rock_point[0] {
					row_start = rock_point[0]
					row_end = rock_line_start_row
				}

				col_start := rock_line_start_col
				col_end := rock_point[1]
				if rock_line_start_col > rock_point[1] {
					col_start = rock_point[1]
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
			rock_line_start_row = rock_point[0]
			rock_line_start_col = rock_point[1]
		}
	}
}

// Emulate sand: flows one unit (cell/tile) at a time, comes to rest, and then the next sand is produced
func EmulateSand(grid [][]int) int {
	still_sand := 0

	// Let sand fall until the sand source is blocked
	for ; grid[SAND_SOURCE_ROW][SAND_SOURCE_COL] == Air; still_sand++ {
		sand_row := SAND_SOURCE_ROW
		sand_col := SAND_SOURCE_COL
		for sand_row < FLOOR_ROW {
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

		// stop processing sand when a sand coordinate reaches the lower grid edge, where we know we have no rock lines to stop it
		if sand_row >= FLOOR_ROW {
			break
		}

		// Note where sand landed still
		grid[sand_row][sand_col] = Sand
	}

	return still_sand
}

func main() {
	rock_line_coordinates, err := GetRockLineCoordinatesFromInput("input.txt")
	if err != nil {
		panic(err)
	}

	// Determine size of grid
	for _, rock_line := range rock_line_coordinates {
		for _, rock_point := range rock_line {
			if rock_point[1] > MAX_COL {
				MAX_COL = rock_point[1]
			}

			if rock_point[0] > MAX_ROW {
				MAX_ROW = rock_point[0]
			}
		}
	}

	// Initialize grid

	// (Part 2 lets us know there's an infinite floor at 2+MAX_ROW)
	FLOOR_ROW = 2 + MAX_ROW

	// NOTE: both inputs do not have any lines that jut up near the left or top of the cave (0), so we shouldn't have to worry about emulating negative indeces
	grid := make([][]int, FLOOR_ROW+1) // floor is edge
	for row := range grid {
		// if we can pile up on the floor, then make sure that COLS also stretches out to MAX_COL+MAX_ROW (give a couple more) so that sand can pile up in a diagonal line to the sand source
		grid[row] = make([]int, MAX_COL+4+FLOOR_ROW)
		// each tile is set to int's default value: 0, which is Air
	}

	// since grid is a slice (created by "make"), it is effectively passed by reference
	FillInRocks(grid, rock_line_coordinates)

	// since grid is a slice (created by "make"), it is effectively passed by reference
	still_sand := EmulateSand(grid)

	// How many units of sand come to rest before sand starts flowing into the abyss below?
	fmt.Printf("Part 1 answer: %v\n", still_sand)

	// Part 2
	// draw infinite floor
	FillInRocks(grid, [][][]int{{{FLOOR_ROW, 0}, {FLOOR_ROW, len(grid[0]) - 1}}})

	// clear grid of sand
	for r, row := range grid {
		for c, tile := range row {
			if tile == Sand {
				grid[r][c] = Air
			}
		}
	}

	// re-emulate to get correct number of sand units
	still_sand = EmulateSand(grid)
	fmt.Printf("Part 2 answer: %v\n", still_sand)
}
