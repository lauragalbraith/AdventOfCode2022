"""
main.py: Laura Galbraith
Description: Solve Day 10 of Advent of Code 2022
What is the sum of these six signal strengths?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines
import re

ADD_X_INSTR_RE = r'^addx (.+)$'

def emulate_cpu(instructions, desired_cycles):
  # single register X, starts as 1
  x = 1
  cycle = 1
  signal_strength_sum = 0

  for instr in instructions:
    final_x = x

    # addx
    if instr != 'noop':
      # save result of instruction for later
      add_x_match = re.match(ADD_X_INSTR_RE, instr)
      final_x = x + int(add_x_match.group(1))

      # check for desired cycle
      if cycle in desired_cycles:
        signal_strength_sum += cycle * x

      # addx's second cycle
      cycle += 1

    # check for desired cycle
    if cycle in desired_cycles:
      signal_strength_sum += cycle * x

    x = final_x
    cycle += 1  # happens for both noop and addx

  return signal_strength_sum

# Parse input
instructions_input = read_ascii_file_lines('input.txt')

# Part 1
# Find the signal strength during the 20th, 60th, 100th, 140th, 180th, and
# 220th cycles. What is the sum of these six signal strengths?
part_1_cycles = set()
part_1_cycles.add(20)
part_1_cycles.add(60)
part_1_cycles.add(100)
part_1_cycles.add(140)
part_1_cycles.add(180)
part_1_cycles.add(220)

signal_sum = emulate_cpu(instructions_input, part_1_cycles)
print(f'Part 1 answer: {signal_sum}')
