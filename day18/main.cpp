// main.cpp: Laura Galbraith
// Description: Solve Day 18 of Advent of Code 2022
// Boiling Boulders
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile
#include <iostream>  // cout, endl
#include <string>  // string, stoi
#include <unordered_map>  // unordered_map
#include <unordered_set>  // unordered_set
#include <vector>  // vector

using namespace std;

int main() {
  // Parse input
  vector<string> coordinates = ReadLinesFromFile("input.txt");

  // NOTE: no negative values, and no values larger than 19

  // Part 1
  // Store all cubes in a large map mapped to their number of exposed sides
  // accessed like map[x][y][z] = sides
  unordered_map<int,unordered_map<int,unordered_map<int,int>>> exposed_per_cube;

  // Input cubes into map
  for (auto coordinate:coordinates) {
    // extract x
    size_t num_start = 0, past_num_end = 1;
    while (past_num_end < coordinate.length() && coordinate[past_num_end] != ',') { ++past_num_end; }
    int x = stoi(coordinate.substr(num_start, past_num_end-num_start));

    // extract y
    num_start = past_num_end + 1;
    past_num_end = num_start + 1;
    while (past_num_end < coordinate.length() && coordinate[past_num_end] != ',') { ++past_num_end; }
    int y = stoi(coordinate.substr(num_start, past_num_end-num_start));

    // extract z
    int z = stoi(coordinate.substr(past_num_end+1));

    int exposed_sides = 6;

    // For two cubes to share a side, they must have an equal number in dimensions A and B, adjacent numbers in dimension C (abs diff of 1)
    vector<int> diffs = {1, -1};

    // find adjacent cube with matching x,y values
    for (auto diff:diffs) {
      if (exposed_per_cube.count(x) && exposed_per_cube[x].count(y) && exposed_per_cube[x][y].count(z+diff)) {
        exposed_per_cube[x][y][z+diff] -= 1;
        --exposed_sides;
      }
    }

    // find adjacent cube with matching x,z values
    for (auto diff:diffs) {
      if (exposed_per_cube.count(x) && exposed_per_cube[x].count(y+diff) && exposed_per_cube[x][y+diff].count(z)) {
        exposed_per_cube[x][y+diff][z] -= 1;
        --exposed_sides;
      }
    }

    // find adjacent cube with matching y,z values
    for (auto diff:diffs) {
      if (exposed_per_cube.count(x+diff) && exposed_per_cube[x+diff].count(y) && exposed_per_cube[x+diff][y].count(z)) {
        exposed_per_cube[x+diff][y][z] -= 1;
        --exposed_sides;
      }
    }

    // save this cube
    exposed_per_cube[x][y][z] = exposed_sides;
  }

  // Total the number of exposed sides across all cubes
  int surface_area = 0;
  for (auto x_k_v:exposed_per_cube) {
    for (auto y_k_v:x_k_v.second) {
      for (auto z_k_v:y_k_v.second) {
        surface_area += z_k_v.second;
      }
    }
  }

  /*
  Any two cubes can share 0 or 1 sides, and no three cubes share a single side between all 3
  For two cubes to share a side, they must have an equal number in dimensions A and B, adjacent numbers in dimension C (abs diff of 1)

  IDEAS
  - store each cube as an object, which keeps a counter of its exposed sides (starts at 6 and gets decremented when a cube is added that has a shared side)
  - find the max/min values in all 3 dimensions
  - count all exposed faces for +/- x/y/z (i.e. 6 perspectives to check from)
    - allocate 3D grid of memory and fill in the cubes, then when you pass from a 0->1 or vice versa that is an exposed side
  - hash coordinate into a single number
  - KD tree?
  - for each new cube that gets read in, rather than iterating over all previous cubes to account for its sides, try +/-1 on all 3 dimensions to search for adjacent cubes (map[x]map[y]set[z] ?)
  */

  // What is the surface area of your scanned lava droplet?
  cout << endl << "Part 1 answer: " << surface_area << endl;
  return 0;
}
