// main.cpp: Laura Galbraith
// Description: Solve Day 23 of Advent of Code 2022
// Unstable Diffusion
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <iostream>  // cout, endl
#include <string>  // string
#include <vector>  // vector

using namespace std;

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("example_input.txt");

  // Part 1
  // Simulate the Elves' process and find the smallest rectangle that contains the Elves after 10 rounds
  // TODO
  
  // How many empty ground tiles does that rectangle contain?
  cout << endl << "Part 1 answer: " << input.size() << endl;

  // Part 2
  // TODO
  cout << endl << "Part 2 answer: " << endl;

  return 0;
}
