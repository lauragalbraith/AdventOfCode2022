"""
main.py: Laura Galbraith
Description: Solve Day 6 of Advent of Code 2022
How many characters need to be processed before the first start-of-packet
marker is detected?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines

# Parse input
stream = read_ascii_file_lines('input.txt')[0]

# Part 1: How many characters need to be processed before the first
# start-of-packet marker is detected?
START_MARKER_LEN = 4

char_count = {}
for i in range(0,len(stream)):
  # remove 5th oldest character from count
  if i >= START_MARKER_LEN:
    dropped_char = stream[i-START_MARKER_LEN]
    char_count[dropped_char] -= 1
    if char_count[dropped_char] == 0:
      char_count.pop(dropped_char)
  # add in newest character
  # (dictionary value default is 0)
  if stream[i] not in char_count:
    char_count[stream[i]] = 0
  char_count[stream[i]] += 1
  # check how many unique characters have been seen
  if len(char_count) == START_MARKER_LEN:
    print(f'Part 1 answer: {i+1}')
    break
