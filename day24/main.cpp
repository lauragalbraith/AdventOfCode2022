// main.cpp: Laura Galbraith
// Description: Solve Day 24 of Advent of Code 2022
// Blizzard Basin
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <algorithm>  // min
#include <cmath>  // abs
#include <iostream>  // cout, endl, ostream
#include <string>  // string
#include <unordered_map>  // unordered_map
#include <vector>  // vector

using namespace std;

// down, right, up, left
vector<int> row_diffs = {1, 0, -1, 0};
vector<int> col_diffs = {0, 1, 0, -1};
vector<char> direction_characters = {'v', '>', '^', '<'};

struct CellState {
  CellState() {
    state = State::kEmpty;
  }

  CellState(const char& input) {
    if (input == '#')
      this->state = State::kWall;
    else if (input == '.') {
      this->state = State::kEmpty;
    } else {
      this->state = State::kBlizzard;

      for (size_t dir = 0; dir < direction_characters.size(); ++dir) {
        if (direction_characters[dir] == input) {
          this->blizzard_dirs = vector<size_t>(1, dir);
          break;
        }
      }
    }
  }

  // copy constructor
  CellState(const CellState& other) {
    this->copy(other);
  }

  // copy assignment operator
  CellState& operator=(const CellState& other) {
    if (this != &other) {
      this->copy(other);
    }

    return *this;
  }

  // destructor
  ~CellState() {}

  enum class State { kEmpty, kWall, kBlizzard };
  State state;

  // if this cell contains one or more blizzards, this stores their directions
  vector<size_t> blizzard_dirs;

 private:
  void copy(const CellState& other) {
    this->state = other.state;
    this->blizzard_dirs = vector<size_t>(other.blizzard_dirs.cbegin(), other.blizzard_dirs.cend());
  }
};

ostream& operator<<(ostream& os, const CellState& c) {
  switch (c.state) {
    case CellState::State::kWall: {
      os << '#';
      break;
    }
    case CellState::State::kEmpty: {
      os << '.';
      break;
    }
    default: {
      if (c.blizzard_dirs.size() > 1) {
        os << c.blizzard_dirs.size();
      } else {
        os << direction_characters[c.blizzard_dirs[0]];
      }
    }
  }

  return os;
}

int ManhattanDistanceToDestination(
  const int& r,
  const int& c,
  const int& destination_row,
  const int& destination_col) {
  return abs(destination_row - r) + abs(destination_col - c);
}

class ValleyState {
 public:
  ValleyState(const vector<string>& input)
    : kRows(static_cast<int>(input.size())), kColumns(static_cast<int>(input[0].size()))
  {
    // initialize state
    this->state_ = vector<vector<CellState>>(this->kRows, vector<CellState>(this->kColumns, CellState()));

    // parse state
    for (int r = 0; r < this->kRows; ++r) {
      for (int c = 0; c < this->kColumns; ++c) {
        this->state_[r][c] = CellState(input[r][c]);
      }
    }
  }

  // copy constructor
  ValleyState(const ValleyState& other) {
    this->copy(other);
  }

  // copy assignment operator
  ValleyState& operator=(const ValleyState& other) {
    if (this != &other) {
      this->clear();
      this->copy(other);
    }

    return *this;
  }

  // destructor
  ~ValleyState() {
    this->clear();
  }

  // move blizzards over a minute
  // takes a reference to a valley rather than creating it and then having to copy it out b/c valleys take up a large amount of space
  void FillInAdvancedState(ValleyState& next_valley) const {
    // empty out blizzard lists in next state
    for (int r = 0; r < next_valley.kRows; ++r) {
      for (int c = 0; c < next_valley.kColumns; ++c) {
        if (next_valley.state_[r][c].state != CellState::State::kWall) {
          next_valley.state_[r][c].state = CellState::State::kEmpty;
        }

        next_valley.state_[r][c].blizzard_dirs.resize(0);
      }
    }

    // move all blizzards represented in this valley to the next one
    for (int r = 0; r < this->kRows; ++r) {
      for (int c = 0; c < this->kColumns; ++c) {
        if (this->state_[r][c].state == CellState::State::kBlizzard) {
          for (auto dir:this->state_[r][c].blizzard_dirs) {
            // calculate where this blizzard will end up
            int next_blizzard_row = r + row_diffs[dir];
            int next_blizzard_col = c + col_diffs[dir];

            // check if the blizzard wraps around the valley
            if (next_valley.state_[next_blizzard_row][next_blizzard_col].state == CellState::State::kWall) {
              if (next_blizzard_row == 0) {
                next_blizzard_row = this->kRows - 2;
              } else if (next_blizzard_row == this->kRows-1) {
                next_blizzard_row = 1;
              }

              if (next_blizzard_col == 0) {
                next_blizzard_col = this->kColumns - 2;
              } else if (next_blizzard_col == this->kColumns-1) {
                next_blizzard_col = 1;
              }
            }

            // save the blizzard in the next valley
            next_valley.state_[next_blizzard_row][next_blizzard_col].state = CellState::State::kBlizzard;
            next_valley.state_[next_blizzard_row][next_blizzard_col].blizzard_dirs.push_back(dir);
          }
        }
      }
    }
  }

  // also returns false if cell is outside the bounds of the valley
  bool IsCellEmpty(const int& row, const int& col) const {
    if (row < 0 || row >= this->kRows || col < 0 || col >= this->kColumns) {
      return false;
    }

    return this->state_[row][col].state == CellState::State::kEmpty;
  }

  friend ostream& operator<<(ostream& os, const ValleyState& v);  // declare as friend so the function, defined elsewhere, can access private members of this class

 private:
  int kRows, kColumns;
  vector<vector<CellState>> state_;

  void copy(const ValleyState& other) {
    this->kRows = other.kRows;
    this->kColumns = other.kColumns;

    this->state_ = vector<vector<CellState>>(other.state_.size(), vector<CellState>(other.state_[0].size(), CellState()));
    for (int r = 0; r < this->kRows; ++r) {
      for (int c = 0; c < this->kColumns; ++c) {
        this->state_[r][c] = other.state_[r][c];
      }
    }
  }

  void clear() {
    this->kRows = -1;
    this->kColumns = -1;
    this->state_.resize(0);
  }
};

ostream& operator<<(ostream& os, const ValleyState& v) {
  os << "valley with " << v.kRows << " rows and " << v.kColumns << " columns" << endl;
  for (int r = 0; r < v.kRows; ++r) {
    for (int c = 0; c < v.kColumns; ++c) {
      os << v.state_[r][c];
    }
    os << endl;
  }

  return os;
}

void DFS(
  int curr_row,
  int curr_col,
  int curr_minute,
  const int& destination_row,
  const int& destination_col,
  int& min_minutes,
  vector<ValleyState>& valley_over_time,
  vector<vector<vector<bool>>>& solved)
{
  // Check for answer
  if (curr_row == destination_row && curr_col == destination_col) {
    min_minutes = min(min_minutes, curr_minute);
    return;
  }

  // Check if this path can feasibly beat the current best answer
  if (ManhattanDistanceToDestination(curr_row, curr_col, destination_row, destination_col) + curr_minute >= min_minutes) {
    return;
  }

  // Mark that we are starting to solve this path
  solved[curr_minute][curr_row][curr_col] = true;

  // Determine what the valley will look like in the next minute
  int next_minute = curr_minute+1;
  if (valley_over_time.size() <= static_cast<size_t>(next_minute)) {
    ValleyState next_valley(valley_over_time[curr_minute]);
    valley_over_time[curr_minute].FillInAdvancedState(next_valley);

    valley_over_time.push_back(next_valley);
  }

  // Try moving
  // if our destination is above our current position, try positions in reverse
  for (size_t move = destination_row > curr_row ? 0 : row_diffs.size()-1;
    move < row_diffs.size();
    move = destination_row > curr_row ? move+1 : move-1) {
    int next_row = curr_row + row_diffs[move];
    int next_col = curr_col + col_diffs[move];

    // check move is in bounds and empty (i.e. not a wall and not a blizzard)
    if (!valley_over_time[next_minute].IsCellEmpty(next_row, next_col)) {
      continue;
    }

    // check if the next minute/position has already been solved
    if (solved[next_minute][next_row][next_col]) {
      continue;
    }

    // take move
    DFS(
      next_row,
      next_col,
      next_minute,
      destination_row,
      destination_col,
      min_minutes,
      valley_over_time,
      solved);
  }

  // Try waiting in place
  // check this space will be open next minute
  if (!valley_over_time[next_minute].IsCellEmpty(curr_row, curr_col)) {
    return;
  }

  // check if the next minute/position has already been solved
  if (solved[next_minute][curr_row][curr_col]) {
    return;
  }

  DFS(
    curr_row,
    curr_col,
    next_minute,
    destination_row,
    destination_col,
    min_minutes,
    valley_over_time,
    solved);
}

int main() {
  // Parse input: map of the valley and the blizzards
  vector<string> input = ReadLinesFromFile("input.txt");

  const int kSourceRow = 0, kSourceCol = 1;  // same across all inputs
  const int kDestinationRow = static_cast<int>(input.size()) - 1;
  const int kDestinationCol = static_cast<int>(input[0].size()) - 2;

  ValleyState start_state(input);

  // we should be able to beat travelling up every non-wall row and down every column twice, accounting for the source/destination positions
  const int kMaxTravelTime = 2 * ((kDestinationRow+1 - 2) * (kDestinationCol+1 - 2)) + 2;

  // Part 1
  int min_minutes_to_bottom_first_time = kMaxTravelTime;

  // credit to https://github.com/ColasNahaboo/advent-of-code-my-solutions/blob/main/go/2022/days/d24/d24.go for the inspiration to:
  // memoize the possible states of the valley, since they are all merely a function of time
  vector<ValleyState> valley_over_time(1, start_state);

  // credit to https://github.com/ColasNahaboo/advent-of-code-my-solutions/blob/main/go/2022/days/d24/d24.go for the inspiration to:
  // track which states we have already tried to find a path from: a function of current position and time
  vector<vector<vector<bool>>> solved(kMaxTravelTime+1, vector<vector<bool>>(kDestinationRow+1, vector<bool>(kDestinationCol+1, false)));  // accessed like [min][row][col]

  DFS(
    kSourceRow,
    kSourceCol,
    0,
    kDestinationRow,
    kDestinationCol,
    min_minutes_to_bottom_first_time,
    valley_over_time,
    solved);

  // What is the fewest number of minutes required to avoid the blizzards and reach the goal?
  cout << endl << "Part 1 answer: " << min_minutes_to_bottom_first_time << endl;

  // Part 2
  // Travel back to the start from the goal
  int min_minutes_to_top = kMaxTravelTime;

  valley_over_time = vector<ValleyState>(1, valley_over_time[min_minutes_to_bottom_first_time]);

  solved = vector<vector<vector<bool>>>(kMaxTravelTime+1, vector<vector<bool>>(kDestinationRow+1, vector<bool>(kDestinationCol+1, false)));

  DFS(
    kDestinationRow,
    kDestinationCol,
    0,
    kSourceRow,
    kSourceCol,
    min_minutes_to_top,
    valley_over_time,
    solved);

  // Travel to the goal again
  int min_minutes_to_bottom_second_time = kMaxTravelTime;

  valley_over_time = vector<ValleyState>(1, valley_over_time[min_minutes_to_top]);

  solved = vector<vector<vector<bool>>>(kMaxTravelTime+1, vector<vector<bool>>(kDestinationRow+1, vector<bool>(kDestinationCol+1, false)));

  DFS(
    kSourceRow,
    kSourceCol,
    0,
    kDestinationRow,
    kDestinationCol,
    min_minutes_to_bottom_second_time,
    valley_over_time,
    solved);

  // What is the fewest number of minutes required to reach the goal, go back to the start, then reach the goal again?
  cout << endl << "Part 2 answer: " << min_minutes_to_bottom_first_time + min_minutes_to_top + min_minutes_to_bottom_second_time << endl;

  return 0;
}
