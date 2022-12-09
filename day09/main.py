"""
main.py: Laura Galbraith
Description: Solve Day 9 of Advent of Code 2022
How many positions does the tail of the rope visit at least once?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines

DIRECTIONS = ['R','U','L','D']  # used with row/col_dirs like list.index('L')
row_dirs = [0, -1, 0, 1]
col_dirs = [1, 0, -1, 0]

class Position:
  """Represents a position of a head or tail in the grid"""

  def __init__(self, r, c):
    self.row = r
    self.col = c

  def __hash__(self):
    # return self.row << 1 | self.col  # not deterministic
    cantor = (self.row + self.col) * (self.row + self.col + 1) / 2 + self.row
    return int(cantor)

  def __eq__(self, other):
    """This and __ne__ were added in an attempt to support hash-chaining"""
    if isinstance(other, self.__class__):
      return self.row == other.row and self.col == other.col
    return False

  def __ne__(self, other):
    return not self.__eq__(other)

  def __str__(self):
    return f'({self.row},{self.col})'

  def move_in_direction(self, direction):
    d = DIRECTIONS.index(direction)
    self.row += row_dirs[d]
    self.col += col_dirs[d]

  def move_to_touch(self, other):
    if not self.touching(other):
      # diagonal move combines:
      # vertical move
      if other.row != self.row:
        self.row += (1 if other.row > self.row else -1)
      # horizontal move
      if other.col != self.col:
        self.col += (1 if other.col > self.col else -1)

  def touching(self, other):
    return (-1 <= self.row - other.row <= 1) and \
           (-1 <= self.col - other.col <= 1)

def count_tail_positions_over_motions(motions, knots):
  """Return the number of positions the tail has been at"""
  tail_positions = set()
  tail_positions.add(knots[-1])

  # Simulate motions to count tail positions
  for motion in motions:
    for _ in range(int(motion[2:])):
      # move head
      knots[0].move_in_direction(motion[0])

      # move each following knot
      for k in range(1,len(knots)):
        knots[k].move_to_touch(knots[k-1])

      # count tail position
      tail_positions.add(knots[-1])

  return len(tail_positions)

GRID_SIDE_MAX = 2000 * 19 + 2  # number of input lines and maximum step size
GRID_MID_POINT = int(GRID_SIDE_MAX/2)

# Parse input
motions_input = read_ascii_file_lines('input.txt')

# Part 1
answer = count_tail_positions_over_motions(motions_input,
  [Position(GRID_MID_POINT, GRID_MID_POINT) for _ in range(2)])
print(f'Part 1 answer: {answer}')

# Part 2
# Simulate your complete series of motions on a larger rope with ten knots.
# How many positions does the tail of the rope visit at least once?
answer = count_tail_positions_over_motions(motions_input,
  [Position(GRID_MID_POINT, GRID_MID_POINT) for _ in range(10)])
print(f'Part 2 answer: {answer}')
