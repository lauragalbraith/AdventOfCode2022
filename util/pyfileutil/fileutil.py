"""
fileutil.py: Laura Galbraith
Description: Implement methods for interacting with files
"""

def read_ascii_file_lines(filename):
  lines = []
  try:
    with open(filename, 'r', encoding='ascii') as f:
      for line in f:
        # remove ending newline
        lines.append(line[0:-1])
  except OSError as e:
    print(f'read_ascii_file_lines: Unable to open file: {e}')
    raise e

  return lines
