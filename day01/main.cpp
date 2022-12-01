// main.cpp: Laura Galbraith
// Description: Solve Day 1 of Advent of Code 2022
// How many total Calories is the Elf with the most Calories carrying?
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

/*NOTES

input is number of Calories each Elf is carrrying; two blank lines separate each elf
seems to be limited to 5 digits (small enough for int)

IDEAS

max PQ - total O(n log n)
or just count up all and keep track of max (can use max PQ later if needed)
*/

#include <algorithm>  // max
#include <fstream>  // ifstream, getline
#include <iostream>  // cout, endl
#include <string>  // string, stoi

using namespace std;

int main() {
  ifstream f("input.txt");

  unsigned long int single_elf_calories = 0;
  unsigned long int max_elf_calories = 0;

  if (f.is_open()) {
    string line;
    while (getline(f, line)) {
      if (line.length() == 0) {
        max_elf_calories = max(max_elf_calories, single_elf_calories);
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
  cout << max_elf_calories << endl;
  return 0;
}
