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

  def move(self, direction):
    d = DIRECTIONS.index(direction)
    self.row += row_dirs[d]
    self.col += col_dirs[d]

  def touching(self, other_position):
    return (-1 <= self.row - other_position.row <= 1) and \
           (-1 <= self.col - other_position.col <= 1)

GRID_SIDE_MAX = 2000 * 19 + 2  # number of input lines and maximum step size
GRID_MID_POINT = int(GRID_SIDE_MAX/2)
# using list rather than tuple so it's mutable
head_pos = Position(GRID_MID_POINT, GRID_MID_POINT)
tail_pos = Position(GRID_MID_POINT, GRID_MID_POINT)
tail_positions = set()
tail_positions.add(tail_pos)

# Parse input
motions = read_ascii_file_lines('input.txt')

# Part 1
# Simulate motions to count tail positions
for motion in motions:
  # print(f'Moving {motion[0]} by {int(motion[2:])}')
  for _ in range(int(motion[2:])):
    # move head
    head_pos.move(motion[0])

    # check if tail is no longer "touching" head, and move it if so
    if not tail_pos.touching(head_pos):
      # print(f'non-touching Head: {head_pos}, Tail: {tail_pos}')

      # diagonal case: move one step diagonally
      if head_pos.row != tail_pos.row and head_pos.col != tail_pos.col:
        tail_pos.row += (1 if head_pos.row > tail_pos.row else -1)
        tail_pos.col += (1 if head_pos.col > tail_pos.col else -1)

      # lateral case: move in same direction as head
      else:
        tail_pos.move(motion[0])

    # count tail position
    # print(f'Final Head: {head_pos}, Tail: {tail_pos}')
    tail_positions.add(tail_pos)
    # print(f'Added {tail_pos} to set and size is now {len(tail_positions)}')

print(f'Part 1 answer: {len(tail_positions)}')
