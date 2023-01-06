// main.cpp: Laura Galbraith
// Description: Solve Day 22 of Advent of Code 2022
// Monkey Map
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <algorithm>  // max, min
#include <iostream>  // cout, endl
#include <regex>  // regex_search, regex
#include <string>  // string, stoi
#include <tuple>  // pair, tuple, get
#include <unordered_map>  // unordered_map
#include <vector>  // vector

using namespace std;

struct CubeFace {
  CubeFace(
    const int& facing,
    const pair<int,int>& top_left)
  : up_facing(facing),
    top_left_coordinate(top_left)
  {}

  CubeFace() {}  // does not create a valid CubeFace

  int up_facing;  // indexes into diffs arrays
  pair<int,int> top_left_coordinate;  // looking at it in the original grid
};

enum class Input { kExample, kFinal };

struct InputParameters {
  // constructor
  InputParameters(
    const string file_name,
    const vector<CubeFace> face_details)
  : input_file_name(file_name),
    cube_faces(face_details.cbegin(), face_details.cend())
  {}

  InputParameters() {}

  const string input_file_name;
  // top-left row/col of the cube face, oriented around face 1 like the example
  const vector<CubeFace> cube_faces;
};

// manually calculated by physically constructing each cube
unordered_map<Input,InputParameters> input_parameters = {
  {Input::kExample, InputParameters(
    "example_input.txt",
    {CubeFace(3, {0,8}),
     CubeFace(3, {4,0}),
     CubeFace(3, {4,4}),
     CubeFace(3, {4,8}),
     CubeFace(3, {8,8}),
     CubeFace(3, {8,12})})},

  {Input::kFinal, InputParameters(
    "input.txt",
    {CubeFace(3, {0, 50}),
     CubeFace(2, {150,0}),
     CubeFace(2, {100,0}),
     CubeFace(3, {50,50}),
     CubeFace(3, {100,50}),
     CubeFace(1, {0,100})})},
};

// map 0-indexed source face idx, then facing, onto the wraparound destination face 1-indexed number and the new facing
// same for all inputs
vector<vector<pair<int,int>>> wraparound_face_facings = {
  // source face: 1
  {{6, 2},  // moving right (0): lands on face 6, facing left (2)
   {4, 1},
   {3, 1},
   {2, 1}},
  // source face: 2
  {{3, 0},
   {5, 3},
   {6, 3},
   {1, 1}},
  // source face: 3
  {{4, 0},
   {5, 0},
   {2, 2},
   {1, 0}},
  // source face: 4
  {{6, 1},
   {5, 1},
   {3, 2},
   {1, 3}},
  // source face: 5
  {{6, 0},
   {2, 3},
   {3, 3},
   {4, 3}},
  // source face: 6
  {{1, 2},
   {2, 0},
   {5, 2},
   {4, 2}},
};

// ordered based on the final password format (i.e. diffs[0] should indicate facing right)
// right, down, left, up
vector<int> row_diffs = {0, 1, 0, -1};
vector<int> col_diffs = {1, 0, -1, 0};

class Board {
 public:
  // parsing constructor
  Board(const vector<string>& input, const vector<CubeFace>& cube_faces)
    // Store cube faces to enable function to wraparound from one face to another
    : cube_faces_(cube_faces),
    // note the size of a face: face 1 is always on top of face 4
    kFaceSize(cube_faces[3].top_left_coordinate.first - cube_faces[0].top_left_coordinate.first)
  {
    // Determine the max row and column to create a consistently-sized board (to ensure correct row/col values)
    size_t max_num_rows = input.size(), max_num_cols = 0;
    for (size_t row = 0; row < max_num_rows; ++row) {
      max_num_cols = max(max_num_cols, input[row].size());
    }

    // Allocate space for board and wraparound memoizations
    this->board_ = vector<vector<Tile>>(max_num_rows, vector<Tile>(max_num_cols, Tile::kEmpty));
    this->top_bottom_rows_ = vector<pair<int,int>>(max_num_cols, {static_cast<int>(max_num_rows), 0});
    this->left_right_cols_ = vector<pair<int,int>>(max_num_rows, {static_cast<int>(max_num_cols), 0});

    // Fill in open and wall spaces from input
    for (size_t row = 0; row < max_num_rows; ++row) {
      for (size_t col = 0; col < input[row].size(); ++col) {
        // save tile type
        this->board_[row][col] = CharToTile(input[row][col]);

        // memoize existent board spaces
        if (this->board_[row][col] != Tile::kEmpty) {
          this->top_bottom_rows_[col].first = min(this->top_bottom_rows_[col].first, static_cast<int>(row));
          this->top_bottom_rows_[col].second = max(this->top_bottom_rows_[col].second, static_cast<int>(row));

          this->left_right_cols_[row].first = min(this->left_right_cols_[row].first, static_cast<int>(col));
          this->left_right_cols_[row].second = max(this->left_right_cols_[row].second, static_cast<int>(col));
        }
      }
    }

    this->ReturnToStartPosition();

    // start problem just with 2D grid as board
    this->is_3D_ = false;
    this->curr_face_ = this->cube_faces_.size();  // set out of range
  }

  void FoldIntoCube() {
    this->is_3D_ = true;

    this->ReturnToStartPosition();

    // Convert the current position to be relative to its face
    this->curr_face_ = 0;  // cube is drawn so starting face is always 1
    this->curr_row_ -= this->cube_faces_[curr_face_].top_left_coordinate.first;
    this->curr_col_ -= this->cube_faces_[curr_face_].top_left_coordinate.second;
  }

  void WalkSteps(const int num_steps) {
    for (int steps_taken = 0; steps_taken < num_steps; ++steps_taken) {
      tuple<int,int,int,size_t> next_tile = this->NextTile();

      // If you run into a wall, you stop moving forward and continue with the next instruction
      if (this->BoardValue(get<0>(next_tile), get<1>(next_tile), get<3>(next_tile)) == Tile::kWall) {
        break;
      }

      // Take step
      this->curr_row_ = get<0>(next_tile);
      this->curr_col_ = get<1>(next_tile);
      this->curr_facing_ = get<2>(next_tile);
      this->curr_face_ = get<3>(next_tile);
    }
  }

  void Turn(const int turn_dir) {
    this->curr_facing_ = (static_cast<int>(row_diffs.size()) + this->curr_facing_ + turn_dir) % static_cast<int>(row_diffs.size());
  }

  // return the final password after the path is walked
  int Password() const {
    // password is the sum of 1000 times the row, 4 times the column, and the facing
    int row = this->curr_row_;
    int col = this->curr_col_;
    int facing = this->curr_facing_;

    // translate from 3D
    if (this->is_3D_) {
      int num_dirs = static_cast<int>(row_diffs.size());
      // effectively turn counterclockwise the number of times we'd have to turn clockwise to match
      // TODO put the logic to rotate clockwise/counterclockwise in a function (named From3D vs To3D?)
      int turns = (num_dirs - (3 - this->cube_faces_[this->curr_face_].up_facing)) + num_dirs % num_dirs;

      for (int t = 0; t < turns; ++t) {
        RotateClockwiseOnFace(row, col);
        facing = (facing+1) % num_dirs;
      }

      // translate from face
      row += this->cube_faces_[this->curr_face_].top_left_coordinate.first;
      col += this->cube_faces_[this->curr_face_].top_left_coordinate.second;
    }

    // row and column being 1-indexed
    return 1000*(row+1) + 4*(col+1) + facing;
  }

 private:
  enum class Tile { kOpen, kWall, kEmpty };  // kEmpty represents a space off the board

  vector<vector<Tile>> board_;
  // memoize the rows/cols of the board to enable O(1) wraparound
  // NOTE: this assumes that the input is structured in a way so there's not multiple juts for any row/column
  vector<pair<int,int>> top_bottom_rows_;
  vector<pair<int,int>> left_right_cols_;

  int curr_row_, curr_col_, curr_facing_;

  bool is_3D_;
  const vector<CubeFace>& cube_faces_;
  const int kFaceSize;
  size_t curr_face_;

  Tile CharToTile(const char& c) {
    switch(c) {
      case '.': return Tile::kOpen;
      case '#': return Tile::kWall;
      default:  return Tile::kEmpty;
    }
  }

  void ReturnToStartPosition() {
    // Set starting position
    // leftmost . in top row of tiles, facing right (positive-col dir)
    this->curr_row_ = 0;
    this->curr_facing_ = 0;

    for (size_t col = 0; col < this->board_[this->curr_row_].size(); ++col) {
      // determine if this should be the start space
      if (this->board_[this->curr_row_][col] == Tile::kOpen) {
        this->curr_col_ = static_cast<int>(col);
        break;
      }
    }
  }

  void RotateClockwiseOnFace(int& row, int& col) const {
    // col becomes row
    int temp_row = col;
    // col is calculated by calculating distance to the axis and then shifting in front of next axis
    col = this->kFaceSize - row - 1 + this->kFaceSize % this->kFaceSize;

    row = temp_row;
  }

  Tile BoardValue(int row, int col, size_t zero_indexed_face) const {
    if (this->is_3D_) {
      // rotate coordinate so perspective is up
      int facing = this->cube_faces_[zero_indexed_face].up_facing;
      while (facing != 3) {  // up
        RotateClockwiseOnFace(row, col);
        facing = (facing-1 + static_cast<int>(row_diffs.size())) % static_cast<int>(row_diffs.size());
      }

      // translate face back into board
      row += this->cube_faces_[zero_indexed_face].top_left_coordinate.first;
      col += this->cube_faces_[zero_indexed_face].top_left_coordinate.second;
    }

    return this->board_[row][col];
  }

  // determine what the row,col, facing,face of the next tile to occupy would be (can be a wall)
  tuple<int,int,int,size_t> NextTile() const {
    int next_row = this->curr_row_ + row_diffs[this->curr_facing_];
    int next_col = this->curr_col_ + col_diffs[this->curr_facing_];
    int next_facing = this->curr_facing_;
    size_t next_face = this->curr_face_;

    // Check for wraparound conditions
    // NOTE: we must check wraparound before checking wall since wall at wraparound edge means movement stops at the initial position without taking any wraparound effect
    bool off_right = false, off_left = false, off_top = false, off_bottom = false;

    // check if we have fallen off the relevant grid itself
    int row_upper_limit = this->is_3D_ ? this->kFaceSize : static_cast<int>(this->board_.size());
    if (next_row < 0) {
      off_top = true;
    } else if (next_row >= row_upper_limit) {
      off_bottom = true;
    }

    if (!off_top && !off_bottom) {
      int col_upper_limit = this->is_3D_ ? this->kFaceSize : static_cast<int>(this->board_[next_row].size());

      if (next_col < 0) {
        off_left = true;
      } else if (next_col >= col_upper_limit) {
        off_right = true;
      }
    }

    // check if the tile is within the 2D bounds but is not on the cube
    if (!this->is_3D_ && !off_top && !off_bottom && !off_left && !off_right && this->BoardValue(next_row, next_col, this->curr_face_) == Tile::kEmpty) {
      switch(this->curr_facing_) {
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

    // Wraparound
    if (!this->is_3D_) {
      if (off_right) {
        next_col = this->left_right_cols_[next_row].first;
      } else if (off_left) {
        next_col = this->left_right_cols_[next_row].second;
      } else if (off_top) {
        next_row = this->top_bottom_rows_[next_col].second;
      } else if (off_bottom) {
        next_row = this->top_bottom_rows_[next_col].first;
      }
    } else if (off_right || off_left || off_bottom || off_top) {
      // calculate next face
      pair<int,int> next_face_facing_1 = wraparound_face_facings[this->curr_face_][this->curr_facing_];

      next_face = static_cast<size_t>(next_face_facing_1.first - 1);
      next_facing = next_face_facing_1.second;

      // calculate position on next face
      next_row = this->curr_row_;  // reset so they're valid coordinates
      next_col = this->curr_col_;

      int temp_facing = this->curr_facing_;
      while (temp_facing != next_facing) {
        RotateClockwiseOnFace(next_row, next_col);
        temp_facing = (temp_facing+1) % static_cast<int>(row_diffs.size());
      }

      // set known position relative to edge
      switch (next_facing) {
        case 0: {
          // facing right, so we must be on the left edge
          next_col = 0;
          break;
        }
        case 1: {
          // facing down, so we must be on the top edge
          next_row = 0;
          break;
        }
        case 2: {
          // facing left, so we must be on right edge
          next_col = this->kFaceSize - 1;
          break;
        }
        default: {
          // facing up, so we must be on bottom edge
          next_row = this->kFaceSize - 1;
        }
      }
    }

    return {next_row, next_col, next_facing, next_face};
  }
};

int main() {
  Input input_choice = Input::kFinal;

  // Parse input
  vector<string> input = ReadLinesFromFile(input_parameters[input_choice].input_file_name);

  // create board from first n-2 lines of input
  vector<string> board_input(input.cbegin(), input.cend() - 2);
  Board board(board_input, input_parameters[input_choice].cube_faces);

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

  // Part 1
  // Follow the path given in the monkeys' notes
  for (auto step:path_steps) {
    board.WalkSteps(step.first);
    board.Turn(step.second);
  }

  // What is the final password?
  cout << endl << "Part 1 answer: " << board.Password() << endl;

  // Part 2
  // Fold the map into a cube
  board.FoldIntoCube();

  // Follow the path given in the monkeys' notes
  for (auto step:path_steps) {
    board.WalkSteps(step.first);
    board.Turn(step.second);
  }

  // What is the final password?
  cout << endl << "Part 2 answer: " << board.Password() << endl;

  return 0;
}
