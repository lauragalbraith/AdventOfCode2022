"""
main.py: Laura Galbraith
Description: Solve Day 8 of Advent of Code 2022
How many trees are visible from outside the grid?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines

# Parse input
tree_height_strings = read_ascii_file_lines('input.txt')
ROWS = len(tree_height_strings)
COLS = len(tree_height_strings[0])  # guaranteed to be stable

# array of arrays of single ints, accessed with [row][col]
tree_heights = [ [0] * COLS for r in range(ROWS) ]
for r in range(ROWS):
  for c in range(COLS):
    tree_heights[r][c] = int(tree_height_strings[r][c])

# Part 1
# Consider your map; how many trees are visible from outside the grid?
# when looking directly along a row or column

# east, north, west, south
row_dirs = [0, -1, 0, 1]
col_dirs = [1, 0, -1, 0]

# Create a memoization of the tallest trees to the east, north, west, south
class TallestToThe:
  """Saves how tall the tallest trees are in each direction"""
  def __init__(self):
    self.e, self.n, self.w, self.s = -1,-1,-1,-1
  def __str__(self):
    return (f'Tallest to the East: {self.e}, North: {self.n}, '
            f'West: {self.w}, South: {self.s}')
  def is_taller_in_any_direction(self, h):
    return h > self.e or h > self.n or h > self.w or h > self.s
  def from_row_col_dirs(self, d):
    return [self.e, self.n, self.w, self.s][d]

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

# Compare each tree (including borders) against tallest in each direction
visible_trees = 0

for r in range(ROWS):
  for c in range(COLS):
    if tallest_to_the[r][c].is_taller_in_any_direction(tree_heights[r][c]):
      visible_trees += 1

# Time: O(trees), which is 99^2; Space: O(trees)
print(f'Part 1 answer: {visible_trees}')

# Part 2
# What is the highest scenic score possible for any tree?
max_scenic_score = 0

# do not consider edges since their scenic score is reduced to zero
for r in range(1, ROWS-1):
  for c in range(1, COLS-1):
    height = tree_heights[r][c]
    scenic_score = 1

    for direction in range(len(row_dirs)):
      visible_dir_trees = 1  # the tree at the end of our path

      # do we already know our sightline ends at the edge of the grid?
      if tallest_to_the[r][c].from_row_col_dirs(direction) < height:
        if direction % 2:
          # north or south: diff by rows
          visible_dir_trees = r if direction == 1 else ROWS - r - 1
        else:
          # east or west: diff by columns
          visible_dir_trees = c if direction == 2 else COLS - c - 1

      # otherwise, walk to the edge of our sightline
      else:
        neighbor_r = r + row_dirs[direction]
        neighbor_c = c + col_dirs[direction]
        while 0 <= neighbor_r < ROWS and 0 <= neighbor_c < COLS and \
            height > tree_heights[neighbor_r][neighbor_c]:
          visible_dir_trees += 1
          neighbor_r += row_dirs[direction]
          neighbor_c += col_dirs[direction]

      scenic_score *= visible_dir_trees

    max_scenic_score = max(max_scenic_score, scenic_score)

# Time: O(n^3) where n is 99; Space: O(n^2)
print(f'Part 2 answer: {max_scenic_score}')
