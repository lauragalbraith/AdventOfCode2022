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

// Make arrays to enable moving a single cell away from a cube, in all directions
constexpr int kDirections = 6;
constexpr int x_diffs[kDirections] = {1, -1, 0,  0, 0,  0};
constexpr int y_diffs[kDirections] = {0,  0, 1, -1, 0,  0};
constexpr int z_diffs[kDirections] = {0,  0, 0,  0, 1, -1};

class LavaDroplet {
public:
  LavaDroplet() {
    this->x_min = 20, this->x_max = 0;
    this->y_min = 20, this->y_max = 0;
    this->z_min = 20, this->z_max = 0;
  }

  void AddCube(string input_line) {
    // extract x
    size_t num_start = 0, past_num_end = 1;
    while (past_num_end < input_line.length() && input_line[past_num_end] != ',') { ++past_num_end; }
    int x = stoi(input_line.substr(num_start, past_num_end-num_start));

    this->x_min = min(this->x_min, x);
    this->x_max = max(this->x_max, x);

    // extract y
    num_start = past_num_end + 1;
    past_num_end = num_start + 1;
    while (past_num_end < input_line.length() && input_line[past_num_end] != ',') { ++past_num_end; }
    int y = stoi(input_line.substr(num_start, past_num_end-num_start));

    this->y_min = min(this->y_min, y);
    this->y_max = max(this->y_max, y);

    // extract z
    int z = stoi(input_line.substr(past_num_end+1));

    this->z_min = min(this->z_min, z);
    this->z_max = max(this->z_max, z);

    int exposed_sides = 6;

    // For two cubes to share a side, they must have an equal number in dimensions A and B, adjacent numbers in dimension C (abs diff of 1)
    for (int diff_i = 0; diff_i < kDirections; ++diff_i) {
      int no_adjacent = -1;
      int& adjacent_exposed_sides = this->AdjacentCubeInDirection(diff_i, x, y, z, no_adjacent);
      if (adjacent_exposed_sides != no_adjacent) {
        // a stored cube and this cube share a side
        adjacent_exposed_sides -= 1;
        --exposed_sides;
      }
    }

    // save this cube
    this->exposed_per_cube[x][y][z] = exposed_sides;
  }

  int TotalExposedSides() {
    int surface_area = 0;
    for (auto x_k_v:exposed_per_cube) {
      for (auto y_k_v:x_k_v.second) {
        for (auto z_k_v:y_k_v.second) {
          surface_area += z_k_v.second;
        }
      }
    }

    return surface_area;
  }

  bool IsCubeInDroplet(int x, int y, int z) {
    return exposed_per_cube.count(x) && exposed_per_cube[x].count(y) && exposed_per_cube[x][y].count(z);
  }

  bool IsSurrounded(int center_x, int center_y, int center_z) {
    // pick a direction to check in
    for (int diff_i = 0; diff_i < kDirections; ++diff_i) {
      bool cube_this_direction = false;

      // continue in that direction while there still might be cubes in that direction
      for (int x = center_x + x_diffs[diff_i], y = center_y + y_diffs[diff_i], z = center_z + z_diffs[diff_i];
           x >= x_min && x <= x_max && y >= y_min && y <= y_max && z >= z_min && z <= z_max;
           x += x_diffs[diff_i], y += y_diffs[diff_i], z += z_diffs[diff_i]) {
        // check for a cube
        if (IsCubeInDroplet(x, y, z)) {
          cube_this_direction = true;
          break;
        }
      }

      if (!cube_this_direction) {
        return false;
      }
    }

    return true;
  }

  int ExternalSurfaceArea() {
    int null_val = -1;

    int area = this->TotalExposedSides();

    // Check all spaces not occupied by a cube inside the min/max boundaries
    for (int x = this->x_min; x <= this->x_max; ++x) {
      for (int y = this->y_min; y <= this->y_max; ++y) {
        for (int z = this->z_min; z <= this->z_max; ++z) {
          // Skip if this cube is part of the droplet
          if (this->IsCubeInDroplet(x, y, z)) {
            continue;
          }

          // Check for exposure
          if (this->IsSurrounded(x, y, z)) {
            // Check all directions for directly-adjacent cubes
            for (int dir = 0; dir < kDirections; ++dir) {
              int& adjacent_exposed_sides = this->AdjacentCubeInDirection(dir, x, y, z, null_val);
              if (adjacent_exposed_sides != null_val) {
                // subtract this shared face from the external surface area
                area--;
              }
            }
          }
        }
      }
    }

    return area;
  }

private:
  // Store all cubes in a large map mapped to their number of exposed sides
  // accessed like map[x][y][z] = sides
  unordered_map<int,unordered_map<int,unordered_map<int,int>>> exposed_per_cube;

  // Track the smallest and largest values in all dimensions
  int x_min, x_max;
  int y_min, y_max;
  int z_min, z_max;

  int& AdjacentCubeInDirection(int direction, int x, int y, int z, int& null_val) {
    // For two cubes to share a side, they must have an equal number in dimensions A and B, adjacent numbers in dimension C (abs diff of 1)
    if (this->IsCubeInDroplet(x + x_diffs[direction], y + y_diffs[direction], z + z_diffs[direction])) {
      // a stored cube and this cube share a side
      return this->exposed_per_cube[x+x_diffs[direction]][y+y_diffs[direction]][z+z_diffs[direction]];
    }

    return null_val;
  }
};

int main() {
  // Parse input
  vector<string> coordinates = ReadLinesFromFile("input.txt");
  // NOTE: no negative values, and no values larger than 19

  LavaDroplet lava = LavaDroplet();

  // Input cubes into map
  for (auto coordinate:coordinates) {
    lava.AddCube(coordinate);
  }
  
  // Part 1
  // What is the surface area of your scanned lava droplet?
  cout << endl << "Part 1 answer: " << lava.TotalExposedSides() << endl;

  // Part 2
  // What is the exterior surface area of your scanned lava droplet?
  cout << endl << "Part 2 answer: " << lava.ExternalSurfaceArea() << endl;

  return 0;
}
