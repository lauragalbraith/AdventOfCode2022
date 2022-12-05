"""
main.py: Laura Galbraith
Description: Solve Day 5 of Advent of Code 2022
After the rearrangement procedure completes, what crate ends up on top of
each stack?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from collections import deque
from copy import deepcopy
import re
from util.pyfileutil.fileutil import read_ascii_file_lines

# cheat a little bit, and see from input file we'll need 9 stacks
NUM_STACKS = 9

# Parse input
stacks = []  # array of deques of chars
for i in range(0, NUM_STACKS):
  stacks.append(deque())
stacks_filled = False

# since we're using a deque to emulate a stack data structure, we can fill them
# in by calling appendleft as we read the file so that pop would result in
# what's closest to the top of the file

movements = []  # array of tuples: (how_many, from, to)
movement_re = r'^move (\d+) from (\d+) to (\d+)'

for line in read_ascii_file_lines('input.txt'):
  # skip blank lines
  if len(line) < 10:
    continue
  # mark when we reach end of stacks list
  if line[0] == ' ' and line[1] == '1':
    stacks_filled = True
    continue
  # parse initial stack configuration
  if not stacks_filled:
    char_i = 0
    for stack_i in range(0,NUM_STACKS):
      if line[char_i] == '[':
        stacks[stack_i].appendleft(line[char_i+1])
      char_i += 4  # move past crate or spaces
  # parse movement instructions
  else:
    match = re.search(movement_re, line)
    movements.append(
      (int(match.group(1)),
      # change 1-indexed stack numbers to 0-indexed
      int(match.group(2))-1,
      int(match.group(3))-1))

# copy our initial stack configuration to be identical for Part 2
initial_stacks = deepcopy(stacks)

# simulate crate movements
for m in movements:
  # crates move one at a time, in LIFO order
  for i in range(0,m[0]):
    stacks[m[2]].append(stacks[m[1]].pop())

# Part 1: After the rearrangement procedure completes, what crate ends up on
# top of each stack?
top_stack_elements = []
for stack in stacks:
  top_stack_elements.append(stack[-1])

answer = ''.join(top_stack_elements)
print(f'Part 1 answer: {answer}')

# Part 2: After the rearrangement procedure completes, what crate ends up on
# top of each stack?
for m in movements:
  # crates move how_many at a time, retaining order
  removed_crates = deque()
  for i in range(0,m[0]):
    removed_crates.appendleft(initial_stacks[m[1]].pop())
  for crate in removed_crates:
    initial_stacks[m[2]].append(crate)

top_stack_elements = []
for stack in initial_stacks:
  top_stack_elements.append(stack[-1])

answer = ''.join(top_stack_elements)
print(f'Part 2 answer: {answer}')
