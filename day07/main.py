"""
main.py: Laura Galbraith
Description: Solve Day 7 of Advent of Code 2022
What is the sum of the total sizes of those directories?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines
import re

class Directory:
  """Directory represents a dir that contains files and other dirs"""
  SMALL_DIR_LIMIT = 100000

  def __init__(self, name, parent):
    self.name = name  # name of the directory
    self.total_size = -1  # mark as unknown
    self.parent = parent  # keep parent pointer
    self.child_dirs = {}  # map of directory names to Directory objects
    self.files = {}  # map of file names to their size

  def add_file(self, file_name, file_size):
    self.files[file_name] = file_size
    # NOTE: we could update self.total_size, but due to parsing it's unknown

  def add_empty_child_directory(self, dir_name):
    self.child_dirs[dir_name] = Directory(dir_name, self)

  def get_parent(self):
    return self.parent

  def access_child_directory(self, dir_name):
    return self.child_dirs[dir_name]

  def calculate_total_size(self):
    self.total_size = 0

    # calculate size of all child directories
    for child_dir in self.child_dirs.items():
      child_dir[1].calculate_total_size()
      self.total_size += child_dir[1].total_size

    # add size of all files in this directory
    for file_size in self.files.items():
      self.total_size += file_size[1]

  def sum_small_dirs(self):
    small_sum = 0

    # sum size of all small child directories
    for child_dir in self.child_dirs.items():
      small_sum += child_dir[1].sum_small_dirs()

    # check if this directory qualifies as small
    if self.total_size <= Directory.SMALL_DIR_LIMIT:
      small_sum += self.total_size

    return small_sum


# Keep dirs as a tree structure, anchored at '/'
root = Directory('/', None)
current_dir = root

# Parse input
cd_re = r'^\$ cd (.+)$'
file_re = r'^(\d+) (.+)$'
dir_re = r'^dir (.+)$'
terminal_output = read_ascii_file_lines('input.txt')

# Identify all directories & files by reading through terminal output
# use a while loop instead of for to manipulate inside block
i = -1
while i < (len(terminal_output) - 1):
  i += 1
  cmd = terminal_output[i]
  # print(f'TESTING: i is {i}, cmd is *{cmd}*')

  # change directory command
  cd_match = re.match(cd_re, cmd)
  if cd_match is not None:
    next_dir = cd_match.group(1)
    # print(f'TESTING: line {i} *{cmd}* is a cd cmd to {next_dir}')
    if next_dir == '/':
      current_dir = root
    elif next_dir == '..':
      current_dir = current_dir.get_parent()
    else:
      current_dir = current_dir.access_child_directory(next_dir)

    continue

  # list command
  if cmd != '$ ls':
    print(f'Terminal output line "{cmd}" is not an expected command')
    raise SyntaxError('Unexpected input', cmd)

  # move through ls output, accounting for empty directories
  while (i+1 < len(terminal_output)) and (terminal_output[i+1][0] != '$'):
    i += 1
    ls_output_line = terminal_output[i]
    # print(f'TESTING: i is now {i} and line from ls is *{ls_output_line}*')

    # parse directory
    dir_match = re.match(dir_re, ls_output_line)
    if dir_match is not None:
      current_dir.add_empty_child_directory(dir_match.group(1))
      continue

    # parse file
    file_match = re.match(file_re, ls_output_line)
    if file_match is None:
      print(f'Terminal output line "{ls_output_line}" is not expected from ls')
      raise SyntaxError('Unexpected input', ls_output_line)

    current_dir.add_file(file_match.group(2), int(file_match.group(1)))

  # print(f'TESTING: Done with ls; i is now {i}')

# For each directory, calculate its size
root.calculate_total_size()

# Find all of the directories with a total size of at most 100000.
# What is the sum of the total sizes of those directories?
small_dir_sum = root.sum_small_dirs()
print(f'Part 1 answer: {small_dir_sum}')

# TODO after completing puzzle, install apt upgrades: apt list --upgradable
