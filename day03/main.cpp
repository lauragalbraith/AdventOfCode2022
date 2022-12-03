// main.cpp: Laura Galbraith
// Description: Solve Day 3 of Advent of Code 2022
// What is the sum of the priorities of those item types?
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include <iostream>  // cout, endl
#include <fstream>  // ifstream, getline
#include <string>  // string
#include <vector>  // vector

using namespace std;

constexpr unsigned long kLettersInAlphabet = 26;
constexpr unsigned long kNumberItemTypes = kLettersInAlphabet * 2;

int main() {
  // Part 1
  // NOTE: input space is limited to lowercase and uppercase letters, so if we need a hashmap we can just use vector[26*2]
  // first compartment is line[:line.length()/2 - 1], second compartment is line[line.length()/2:]
  // we care about item types shared across both compartments of the same rucksack

  // parse input
  vector<string> input_lines;
  ifstream f("input.txt");
  if (f.is_open()) {
    string line;
    while (getline(f, line)) {
      input_lines.push_back(line);
    }
  } else {
    cout << "Unable to open input file" << endl;
    return -1;
  }

  // collect cumulative priority sum of duplicate item types
  unsigned long priority_sum = 0;

  // for each line, find the first duplicate item type
  for (auto rucksack:input_lines) {
    vector<bool> seen(kNumberItemTypes, false);

    unsigned long dividing_line = rucksack.length() / 2;
    for (unsigned long i = 0; i < dividing_line; ++i) {
      // unintuitively, lowercase letters have greater ASCII values than uppercase
      unsigned long pos = rucksack[i] >= 'a' ? rucksack[i] - 'a' : rucksack[i] - 'A' + kLettersInAlphabet;
      seen[pos] = true;
    }

    for (unsigned long i = dividing_line; i < rucksack.length(); ++i) {
      unsigned long pos = rucksack[i] >= 'a' ? rucksack[i] - 'a' : rucksack[i] - 'A' + kLettersInAlphabet;
      if (seen[pos]) {
        priority_sum += pos + 1;  // 1-indexed
        break;
      }
    }
  }

  cout << "Part 1 answer: " << priority_sum << endl;

  return 0;
}
