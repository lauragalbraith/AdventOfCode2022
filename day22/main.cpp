// main.cpp: Laura Galbraith
// Description: Solve Day 22 of Advent of Code 2022
// Monkey Map
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <algorithm>  // max, min
#include <iostream>  // cout, endl
#include <regex>  // regex_search, regex
#include <string>  // string, stoi
#include <tuple>  // pair
#include <vector>  // vector

using namespace std;

// ordered based on the final password format (i.e. diffs[0] should indicate facing right)
// right, down, left, up
vector<int> row_diffs = {0, 1, 0, -1};
vector<int> col_diffs = {1, 0, -1, 0};

class Board {
 public:
  // parsing constructor
  Board(const vector<string>& input) {
    // Determine the max row and column to create a consistently-sized board (to ensure correct row/col values)
    size_t max_num_rows = input.size(), max_num_cols = 0;
    for (size_t row = 0; row < max_num_rows; ++row) {
      max_num_cols = max(max_num_cols, input[row].size());
    }

    // Allocate space for board and wraparound memoizations
    this->board = vector<vector<Tile>>(max_num_rows, vector<Tile>(max_num_cols, Tile::kEmpty));
    this->top_bottom_rows = vector<pair<int,int>>(max_num_cols, {static_cast<int>(max_num_rows), 0});
    this->left_right_cols = vector<pair<int,int>>(max_num_rows, {static_cast<int>(max_num_cols), 0});

    // Set starting position
    // leftmost . in top row of tiles, facing right (positive-col dir)
    this->curr_row = 0;
    this->curr_col = -1;
    this->curr_facing = 0;

    // Fill in open and wall spaces from input
    for (size_t row = 0; row < max_num_rows; ++row) {
      for (size_t col = 0; col < input[row].size(); ++col) {
        // save tile type
        this->board[row][col] = CharToTile(input[row][col]);

        // memoize existent board spaces
        if (this->board[row][col] != Tile::kEmpty) {
          this->top_bottom_rows[col].first = min(this->top_bottom_rows[col].first, static_cast<int>(row));
          this->top_bottom_rows[col].second = max(this->top_bottom_rows[col].second, static_cast<int>(row));

          this->left_right_cols[row].first = min(this->left_right_cols[row].first, static_cast<int>(col));
          this->left_right_cols[row].second = max(this->left_right_cols[row].second, static_cast<int>(col));
        }

        // determine if this should be the start space
        if (static_cast<int>(row) == curr_row && curr_col < 0 && this->board[row][col] == Tile::kOpen) {
          this->curr_col = static_cast<int>(col);
        }
      }
    }

    // cout << "DEBUG: starting space in the board is row:" << this->curr_row << ", col:" << this->curr_col << ", facing " << this->curr_facing << endl;
  }

  void WalkSteps(const int num_steps) {
    // cout << "DEBUG: starting at " << this->curr_row << "," << this->curr_col << endl;

    for (int steps_taken = 0; steps_taken < num_steps; ++steps_taken) {
      // Look at step forward
      int next_row = this->curr_row + row_diffs[this->curr_facing];
      int next_col = this->curr_col + col_diffs[this->curr_facing];
      // cout << "DEBUG: before checking any bounds, next_row is " << next_row << " and next_col is " << next_col << endl;

      // Wraparound if applicable
      // NOTE: we must check wraparound before checking wall since wall at wraparound edge means movement stops at the initial position without taking any wraparound effect
      bool off_right = false, off_left = false, off_top = false, off_bottom = false;
      if (next_row < 0) {
        off_top = true;
      } else if (next_row >= static_cast<int>(this->board.size())) {
        off_bottom = true;
      } else if (next_col < 0) {
        off_left = true;
      } else if (next_col >= static_cast<int>(this->board[next_row].size())) {
        off_right = true;
      } else if (this->board[next_row][next_col] == Tile::kEmpty) {
        switch(this->curr_facing) {
          case 0: {
            off_right = true;
            break;
          }
          case 1: {
            off_bottom = true;
            break;
          }
          case 2: {
            off_left = true;
            break;
          }
          default: {
            off_top = true;
            break;
          }
        }
      }

      if (off_right) {
        next_col = this->left_right_cols[next_row].first;
      } else if (off_left) {
        next_col = this->left_right_cols[next_row].second;
        // cout << "DEBUG: stepped off the left edge of the board; instead coming around to column " << next_col << endl;
      } else if (off_top) {
        next_row = this->top_bottom_rows[next_col].second;
      } else if (off_bottom) {
        next_row = this->top_bottom_rows[next_col].first;
      }

      // If you run into a wall, you stop moving forward and continue with the next instruction
      if (this->board[next_row][next_col] == Tile::kWall) {
        // cout << "DEBUG: blocked by a wall at " << next_row << "," << next_col << endl;
        break;
      }

      // Take step
      // cout << "DEBUG: stepping from " << this->curr_row << "," << this->curr_col << " to " << next_row << "," << next_col << endl;
      this->curr_row = next_row;
      this->curr_col = next_col;
    }

    if (this->board[this->curr_row][this->curr_col] != Tile::kOpen) {
      cout << "DEBUG : THIS SHOULD NEVER HAPPEN - we are not in an open space after movement" << endl;
    }

    // cout << "DEBUG: after attempting " << num_steps << " steps, we are now at row:" << this->curr_row << ", col:" << this->curr_col << endl;
  }

  void Turn(const int turn_dir) {
    this->curr_facing = (static_cast<int>(row_diffs.size()) + this->curr_facing + turn_dir) % static_cast<int>(row_diffs.size());
    // cout << "DEBUG: now facing " << this->curr_facing << " after turning " << turn_dir << endl;
  }

  // return the final password after the path is walked
  int Password() const {
    // password is the sum of 1000 times the row, 4 times the column, and the facing
    // row and column being 1-indexed
    return 1000*(this->curr_row+1) + 4*(this->curr_col+1) + this->curr_facing;
  }

 private:
  enum class Tile { kOpen, kWall, kEmpty };  // kEmpty represents a space off the board

  vector<vector<Tile>> board;
  // memoize the rows/cols of the board to enable O(1) wraparound
  // NOTE: this assumes that the input is structured in a way so there's not multiple juts for any row/column
  vector<pair<int,int>> top_bottom_rows;
  vector<pair<int,int>> left_right_cols;

  int curr_row, curr_col, curr_facing;

  Tile CharToTile(const char& c) {
    switch(c) {
      case '.': return Tile::kOpen;
      case '#': return Tile::kWall;
      default:  return Tile::kEmpty;
    }
  }
};

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("input.txt");

  // create board from first n-2 lines of input
  vector<string> board_input(input.cbegin(), input.cend() - 2);
  Board board(board_input);

  // form the path steps from the last line of input
  // pair.first: number of steps to take
  // pair.second: turn-diff-change (1 or -1 to change diff_i, or 0 for last step)
  vector<pair<int,int>> path_steps;
  regex path_rgx("(\\d+)([LR]*)");

  string remaining_input = input[input.size()-1];
  smatch sm;
  while (regex_search(remaining_input, sm, path_rgx)) {
    int num_steps = stoi(sm[1]);

    int turn_dir = 0;
    if (sm.size() > 2 && sm[2] != "") {
      turn_dir = sm[2] == "L" ? -1 : 1;
    }

    path_steps.push_back({num_steps, turn_dir});
    remaining_input = sm.suffix().str();
  }

  /*cout << "DEBUG: instructions:" << endl;
  for (auto instr:path_steps) {
    cout << "move " << instr.first << " steps then turn " << instr.second << endl;
  }*/

  // Part 1
  // Follow the path given in the monkeys' notes
  for (auto step:path_steps) {
    board.WalkSteps(step.first);
    board.Turn(step.second);
  }

  /*
  IDEAS
  - movement wraps around the board, so empty tiles are never occupied, but traversed in the wraparound movement
    - idea: could implement as brute-force, walking backward until we reach the correct edge
  */

  // What is the final password?
  cout << endl << "Part 1 answer: " << board.Password() << endl;

  // Part 2
  // TODO
  cout << endl << "Part 2 answer: " << endl;

  return 0;
}
