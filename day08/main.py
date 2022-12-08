"""
main.py: Laura Galbraith
Description: Solve Day 8 of Advent of Code 2022
How many trees are visible from outside the grid?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines
# from collections import namedtuple  # TODO clean up

# Parse input
tree_height_strings = read_ascii_file_lines('input.txt')
ROWS = len(tree_height_strings)
COLS = len(tree_height_strings[0])  # guaranteed to be stable

# array of arrays of single ints, accessed with [row][col]
tree_heights = [ [0] * COLS for r in range(ROWS) ]
for r in range(ROWS):
  for c in range(COLS):
    tree_heights[r][c] = int(tree_height_strings[r][c])

# testing - TODO remove
# for r in range(5):
  # row = tree_heights[r]
  # print(f'row: len {len(row)}; first elem: {row[0]} of type {type(row[0])}')

# Part 1
# Consider your map; how many trees are visible from outside the grid?
# when looking directly along a row or column

# Create a memoization of the tallest trees to the east, north, west, south
class TallestToThe:
  def __init__(self):
    self.e, self.n, self.w, self.s = -1,-1,-1,-1
  def __str__(self):
    return (f'Tallest to the East: {self.e}, North: {self.n}, '
            f'West: {self.w}, South: {self.s}')
  def is_taller_in_any_direction(self, height):
    return height > self.e or height > self.n or \
          height > self.w or height > self.s

# TallestToThe = namedtuple('TallestToThe', 'e n w s')  # TODO clean up
# array of arrays of TallestToThe types
tallest_to_the = [[TallestToThe() for c in range(COLS)] for r in range(ROWS)]

# calculate west and north numbers using traditional iteration
tallest_per_col = [-1] * COLS
for r in range(ROWS):
  tallest_per_row = -1
  for c in range(COLS):
    # tallest to the west
    tallest_to_the[r][c].w = tallest_per_row
    tallest_per_row = max(tallest_per_row, tree_heights[r][c])
    # tallest to the north
    tallest_to_the[r][c].n = tallest_per_col[c]
    tallest_per_col[c] = max(tallest_per_col[c], tree_heights[r][c])

# calculate east and south numbers using reverse iteration
tallest_per_col = [-1] * COLS
for r in reversed(range(ROWS)): # range(ROWS-1, -1, -1):
  tallest_per_row = -1
  for c in reversed(range(COLS)): # range(COLS-1, -1, -1):
    # tallest to the east
    tallest_to_the[r][c].e = tallest_per_row
    tallest_per_row = max(tallest_per_row, tree_heights[r][c])
    # tallest to the north
    tallest_to_the[r][c].s = tallest_per_col[c]
    tallest_per_col[c] = max(tallest_per_col[c], tree_heights[r][c])

# testing - TODO remove
# for r in range(ROWS):
  # row = tallest_to_the[r]
  # for c in range(COLS):
    # col = row[c]
    # print(f'row:{r}, col:{c}: Tallest to the {col}')

# Compare each tree (including borders) against tallest in each direction
visible_trees = 0

"""
# IDEAS: brute-force: check every row and column, front and back, and mark
a matrix for every visible tree whose value is not strictly ascending in
that direction; count number of ones in final matrix
# BFS, treating edges as only N-S-E-W movements; but some winding path that
is strictly increasing finds an inner tree, that tree could be blocked from
view on all the borders of that winding path and should not be counted as
visible
# DFS from each border tree, only treating >= values as edges - but if we
got walled in by really short trees we might never reach a super tall inner
tree which should be counted
# DFS all inner trees, only treating < trees as edges - but this would
incorrectly count an inner tree that is walled in by short trees but then
walled in again by really tall trees; each DFS path must be a straight line
to the edge and must reach the edge for the destination to be counted
# could I combine this with DP and store four numbers for each tree: the
tallest to the north, tallest to the east, etc. then the DFS wouldn't have
to continue too far when the answer is in the memo
"""
for r in range(ROWS):
  for c in range(COLS):
    if tallest_to_the[r][c].is_taller_in_any_direction(tree_heights[r][c]):
      visible_trees += 1

# Time: O(trees), which is 99^2; Space: O(trees)
print(f'Part 1 answer: {visible_trees}')

# Part 2
# TODO
