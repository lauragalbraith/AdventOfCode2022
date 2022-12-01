// main.cpp: Laura Galbraith
// Description: Solve Day 1 of Advent of Code 2022
// How many total Calories is the Elf with the most Calories carrying?
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include <algorithm>  // max
#include <fstream>  // ifstream, getline
#include <iostream>  // cout, endl
#include <queue>  // priority_queue
#include <string>  // string, stoi

using namespace std;

// Time: O(L + n * log n) where L is the number of lines in the file and n is the number of Elves
int main() {
  ifstream f("input.txt");

  unsigned long int single_elf_calories = 0;
  priority_queue<unsigned long int> all_elves;

  if (f.is_open()) {
    string line;
    while (getline(f, line)) {
      if (line.length() == 0) {
        all_elves.push(single_elf_calories);
        single_elf_calories = 0;
      } else {
        single_elf_calories += stoi(line);
      }
    }
  } else {
    cout << "Unable to open file" << endl;
    return -1;
  }

  // Part 1 answer: just the top-Calorie Elf
  cout << "Part 1 answer: " << all_elves.top() << endl;

  // Part 2 answer: sum of top-3 Calorie Elves
  // Time to break out the PQ
  unsigned long int top_3_sum = 0;
  for (int i = 0; i < 3; ++i) {
    top_3_sum += all_elves.top();
    all_elves.pop();
  }

  cout << "Part 2 answer: " << top_3_sum << endl;

  return 0;
}
