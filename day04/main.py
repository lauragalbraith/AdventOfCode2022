# main.py: Laura Galbraith
# Description: Solve Day 4 of Advent of Code 2022
# In how many assignment pairs does one range fully contain the other?
# Compile and run: python3 main.py

# Parse input
elf_pairs = []  # will be array of arrays of tuples, where each inner array has two tuples, representing an Elf range
with open("input.txt", "r") as f:
  for line in f:
    line_elves = []
    ranges = (line.rstrip()).split(',')
    for range in ranges:
      numbers = range.split('-')
      # always two numbers in a range: start, end
      line_elves.append((int(numbers[0]), int(numbers[1])))
    elf_pairs.append(line_elves)

# Count all the totally-encompassed elves
encompassed_elves = 0
for p in elf_pairs:
  if p[0][0] == p[1][0] or p[0][1] == p[1][1]:  # if starts are equal, or ends are equal
    encompassed_elves += 1
  elif p[0][0] < p[1][0] and p[0][1] > p[1][1]:  # first elf covers second elf
    encompassed_elves += 1
  elif p[1][0] < p[0][0] and p[1][1] > p[0][1]:  # second elf covers first elf
    encompassed_elves += 1

print("Part 1 answer: {}".format(encompassed_elves))

# Part 2: In how many assignment pairs do the ranges overlap?

overlapping_pairs = 0
for p in elf_pairs:
  if p[0][1] >= p[1][0] and p[0][1] <= p[1][1]:  # first elf ends inside of second elf
    overlapping_pairs += 1
  elif p[1][1] >= p[0][0] and p[1][1] <= p[0][1]:  # second elf ends inside of first elf
    overlapping_pairs += 1

print("Part 2 answer: {}".format(overlapping_pairs))
