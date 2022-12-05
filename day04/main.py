"""
main.py: Laura Galbraith
Description: Solve Day 4 of Advent of Code 2022
In how many assignment pairs does one range fully contain the other?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

# Parse input

# elf_pairs will be array of arrays of tuples, where each inner array has two
# tuples, representing an Elf range
elf_pairs = []
with open('input.txt', 'r', encoding='ascii') as f:
  for line in f:
    line_elves = []
    elf_ranges = (line.rstrip()).split(',')
    for r in elf_ranges:
      numbers = r.split('-')
      # always two numbers in a range: start, end
      line_elves.append((int(numbers[0]), int(numbers[1])))
    elf_pairs.append(line_elves)

# Count all the totally-encompassed elves
encompassed_elves = 0
for p in elf_pairs:
  # if starts are equal, or ends are equal
  if p[0][0] == p[1][0] or p[0][1] == p[1][1]:
    encompassed_elves += 1
  # first elf covers second elf
  elif p[0][0] < p[1][0] and p[0][1] > p[1][1]:
    encompassed_elves += 1
  # second elf covers first elf
  elif p[1][0] < p[0][0] and p[1][1] > p[0][1]:
    encompassed_elves += 1

print(f'Part 1 answer: {encompassed_elves}')

# Part 2: In how many assignment pairs do the ranges overlap?

overlapping_pairs = 0
for p in elf_pairs:
  # first elf ends inside of second elf
  if p[0][1] >= p[1][0] and p[0][1] <= p[1][1]:
    overlapping_pairs += 1
  # second elf ends inside of first elf
  elif p[1][1] >= p[0][0] and p[1][1] <= p[0][1]:
    overlapping_pairs += 1

print(f'Part 2 answer: {overlapping_pairs}')
