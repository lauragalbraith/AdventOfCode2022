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

CRT_WIDTH = 40

def add_signal_strength(
  desired_cycles,
  current_cycle,
  current_x,
  signal_strength_sum):

  if current_cycle in desired_cycles:
    signal_strength_sum += current_cycle * current_x

  return signal_strength_sum

def draw_crt(crt_image, current_cycle, current_x, _):
  # NOTE: cycle is 1-indexed and crt_image is 0-indexed
  # pixel 0 is drawn during cycle 1, then pixel 1 during cycle 2, etc.
  pixel_row = (current_cycle-1) // CRT_WIDTH
  pixel_col = (current_cycle-1) % CRT_WIDTH

  # check if sprite is close to the pixel being drawn
  # X register value is the location of the middle sprite pixel
  if current_x >= 0 and -1 <= (pixel_col - current_x) <= 1:
    crt_image[pixel_row * CRT_WIDTH + pixel_col] = '#'

def emulate_cpu(instructions, desired_output, during_cycle_func):
  # single register X, starts as 1
  x = 1
  cycle = 1
  signal_sum = 0

  for instr in instructions:
    final_x = x

    # addx
    if instr != 'noop':
      # save result of instruction for later
      add_x_match = re.match(ADD_X_INSTR_RE, instr)
      final_x = x + int(add_x_match.group(1))

      # run func "during" cycle
      signal_sum = during_cycle_func(desired_output, cycle, x, signal_sum)

      # addx's second cycle
      cycle += 1

    # run func "during" cycle
    signal_sum = during_cycle_func(desired_output, cycle, x, signal_sum)

    x = final_x
    cycle += 1  # happens for both noop and addx

  return signal_sum

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

part_1_answer = emulate_cpu(
  instructions_input,
  part_1_cycles,
  add_signal_strength)
print(f'Part 1 answer: {part_1_answer}')

# Part 2
crt_array = ['.' for _ in range(240)]
emulate_cpu(instructions_input, crt_array, draw_crt)
print('Part 2 answer: decipher capital letters below:')
row_start = 0
for _ in range(6):
  crt_slice = crt_array[row_start:row_start+CRT_WIDTH]
  crt_str = ''.join(crt_slice)
  print(crt_str)
  row_start += CRT_WIDTH
