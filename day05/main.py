"""
main.py: Laura Galbraith
Description: Solve Day 5 of Advent of Code 2022
After the rearrangement procedure completes, what crate ends up on top of
each stack?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from collections import deque
import re

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

with open('input.txt', 'r', encoding='ascii') as f:
  # parse initial stack configuration
  for line in f:
    # line = line.rstrip()  # not removing newlines b/c spaces are demarking
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
        # if len(line) <= char_i:  # spaces continue until end of line...
          # break
        # print(f'on stack_i {stack_i} and char_i is {char_i}')
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

# simulate crate movements
for m in movements:
  # crates move one at a time, in LIFO order
  for i in range(0,m[0]):
    stacks[m[2]].append(stacks[m[1]].pop())

# After the rearrangement procedure completes, what crate ends up on top of
# each stack?
top_stack_elements = []
for stack in stacks:
  top_stack_elements.append(stack[-1])

answer = ''.join(top_stack_elements)
print(f'Part 1 answer: {answer}')

# TODO after solving Day 5, write a python module to load puzzle input
