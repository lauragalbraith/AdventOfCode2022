"""
main.py: Laura Galbraith
Description: Solve Day 6 of Advent of Code 2022
How many characters need to be processed before the first start-of-packet
marker is detected?
Compile and run: pylint --rcfile ~/pylintrc main.py && python3 main.py
pylintrc sourced from https://google.github.io/styleguide/pyguide.html
"""

from util.pyfileutil.fileutil import read_ascii_file_lines

# returns the number of characters that ends with x unique characters
def find_x_unique_chars(stream, x):
  char_count = {}
  for i in range(0,len(stream)):
    # remove (x+1)th oldest character from count
    if i >= x:
      dropped_char = stream[i-x]
      char_count[dropped_char] -= 1
      if char_count[dropped_char] == 0:
        char_count.pop(dropped_char)
    # add in newest character
    # (dictionary value default is 0)
    if stream[i] not in char_count:
      char_count[stream[i]] = 0
    char_count[stream[i]] += 1
    # check how many unique characters have been seen
    if len(char_count) == x:
      return i+1

  return -1

# Parse input
input_stream = read_ascii_file_lines('input.txt')[0]

# Part 1: How many characters need to be processed before the first
# start-of-packet marker is detected?
print(f'Part 1 answer: {find_x_unique_chars(input_stream, 4)}')

# Part 2: How many characters need to be processed before the first
# start-of-message marker is detected?
print(f'Part 2 answer: {find_x_unique_chars(input_stream, 14)}')
