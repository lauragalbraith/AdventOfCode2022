/*
Package main solves Day 12 of Advent of Code 2022
main.go: Laura Galbraith
What is the fewest steps required to move from your current position to the location that should get the best signal?
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"container/heap"
	"fmt"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

type CellToVisit struct {
	row              int
	col              int
	dist_from_source uint
}

func (c *CellToVisit) GetNeighbors() []CellToVisit {
	// directions: up,down,left,right
	row_diffs := []int{0, 1, 0, -1}
	col_diffs := []int{1, 0, -1, 0}
	neighbors := make([]CellToVisit, len(row_diffs))

	for i, _ := range row_diffs {
		neighbors[i] = CellToVisit{
			row: c.row + row_diffs[i],
			col: c.col + col_diffs[i],
		}
	}

	return neighbors
}

type CellPriorityQueue struct {
	cells []CellToVisit
}

func (pq *CellPriorityQueue) Len() int {
	return len(pq.cells)
}

// needed for sort.Interface
func (pq *CellPriorityQueue) Less(i, j int) bool {
	// return true if i has higher priority than j
	return pq.cells[i].dist_from_source < pq.cells[j].dist_from_source
}

// needed for sort.Interface
func (pq *CellPriorityQueue) Swap(i, j int) {
	pq.cells[i], pq.cells[j] = pq.cells[j], pq.cells[i]
}

func (pq *CellPriorityQueue) Push(x any) {
	new_cell := x.(CellToVisit)
	pq.cells = append(pq.cells, new_cell)
}

func (pq *CellPriorityQueue) Pop() any {
	cell_to_return := pq.cells[len(pq.cells)-1]
	pq.cells = pq.cells[0 : len(pq.cells)-1]
	return cell_to_return
}

func FewestStepsFromSource(heightmap_lines []string, source_row, source_col int, ascending bool) [][]uint {
	ROWS := len(heightmap_lines)
	COLS := len(heightmap_lines[0])

	// Track the minimum distance to source found
	path_len := make([][]uint, ROWS)
	// Keep track of which cell is each cell's predecessor
	predecessors := make([][]CellToVisit, ROWS)
	for row, _ := range path_len {
		path_len[row] = make([]uint, COLS)
		predecessors[row] = make([]CellToVisit, COLS)

		for col, _ := range path_len[row] {
			path_len[row][col] = uint(ROWS*COLS + 1)
			predecessors[row][col] = CellToVisit{row: -1, col: -1}
		}
	}
	path_len[source_row][source_col] = 0

	// Form a priority queue of cells to try next
	var pq CellPriorityQueue
	heap.Init(&pq)
	heap.Push(&pq, CellToVisit{row: source_row, col: source_col, dist_from_source: 0})

	// Until we can no longer improve, keep trying the best cell in the PQ
	for len(pq.cells) > 0 {
		curr_cell := heap.Pop(&pq).(CellToVisit)

		// make sure current cell has most up-to-date best distance to source
		curr_cell.dist_from_source = path_len[curr_cell.row][curr_cell.col]

		// fmt.Printf("DEBUG: current cell is [%d,%d], %v from source\n", curr_cell.row, curr_cell.col, curr_cell.dist_from_source)

		// add closer neighbors to the list to be considered
		curr_height := int(heightmap_lines[curr_cell.row][curr_cell.col])

		for _, n := range curr_cell.GetNeighbors() {
			// check valid are bounds
			if n.row < 0 || n.row >= ROWS || n.col < 0 || n.col >= COLS {
				continue
			}

			// check visiting this neighbor is possible
			n_height := int(heightmap_lines[n.row][n.col])
			ascent := n_height - curr_height
			if !ascending {
				ascent = curr_height - n_height
			}

			// valid edges: destination cell can be at most one higher than elevation of current cell, and can be as low as you want
			// fmt.Printf("DEBUG: neighbor %c - current %c = %d; will it continue? %v\n", heightmap_lines[n.row][n.col], heightmap_lines[curr_cell.row][curr_cell.col], int(heightmap_lines[n.row][n.col]-heightmap_lines[curr_cell.row][curr_cell.col]), int(heightmap_lines[n.row][n.col])-int(heightmap_lines[curr_cell.row][curr_cell.col]) > 1)
			if ascent > 1 {
				continue
			}

			// check if visiting from the current cell is an improvement
			n.dist_from_source = curr_cell.dist_from_source + 1
			if path_len[n.row][n.col] > n.dist_from_source {
				// update paths
				path_len[n.row][n.col] = n.dist_from_source
				predecessors[n.row][n.col] = curr_cell // TODO FINALLY remove if we never need to reconstruct the path

				// add neighbor to PQ
				// do not bother removing old value in PQ; it should not amount to anything
				heap.Push(&pq, n)
			}
		}
	}

	return path_len
}

func main() {
	// Get input
	heightmap_lines, err := fileutil.GetLinesFromFile("example_input.txt")
	if err != nil {
		panic(err)
	}

	// Part 1: What is the fewest steps required to move from your current position to the location that should get the best signal?

	// Save, then overwrite source and end cells
	var s_row, s_col, e_row, e_col int
	for r, row := range heightmap_lines {
		for c, height := range row {
			// S (at height a) is current position, E (at height z) is best signal location
			if height == 'S' {
				s_row, s_col = r, c

				byte_arr := []byte(heightmap_lines[s_row])
				byte_arr[s_col] = byte('a')
				heightmap_lines[s_row] = string(byte_arr)
			}
			if height == 'E' {
				e_row, e_col = r, c

				byte_arr := []byte(heightmap_lines[e_row])
				byte_arr[e_col] = byte('z')
				heightmap_lines[e_row] = string(byte_arr)
			}
		}
	}

	// Compute all path lengths, treating 'S' as source
	path_lengths := FewestStepsFromSource(heightmap_lines, s_row, s_col, true)

	// Answer is path length to 'E' cell
	fmt.Printf("\nPart 1 answer: %d\n", path_lengths[e_row][e_col])

	// Part 2: What is the fewest steps required to move starting from any square with elevation a to the location that should get the best signal?
	path_lengths = FewestStepsFromSource(heightmap_lines, e_row, e_col, false)

	// find the minimum among all 'a' cells
	min_path_len := path_lengths[s_row][s_col]
	for r, row := range heightmap_lines {
		for c, height := range row {
			if height == 'a' && path_lengths[r][c] < min_path_len {
				min_path_len = path_lengths[r][c]
			}
		}
	}

	fmt.Printf("\nPart 2 answer: %d\n", min_path_len)
}
