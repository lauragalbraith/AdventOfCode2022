// main.cpp: Laura Galbraith
// Description: Solve Day 23 of Advent of Code 2022
// Unstable Diffusion
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <algorithm>  // min, max
#include <iostream>  // cout, endl, ostream
#include <string>  // string
#include <tuple>  // pair
#include <unordered_map>  // unordered_map
#include <utility>  // hash
#include <vector>  // vector

using namespace std;

struct Coordinate {
  size_t row, col;

  // default constructor; does not necessarily represent a valid coordinate
  Coordinate() {
    this->row = 0;
    this->row = 0;
  }

  // expected constructor
  Coordinate(const size_t& r, const size_t& c): row(r), col(c) {}

  // copy assignment operator
  Coordinate& operator=(const Coordinate& other) {
    if (this != &other) {
      this->copy(other);
    }
    return *this;
  }

  // copy constructor
  Coordinate(const Coordinate& other) {
    this->copy(other);
  }

  // destructor
  ~Coordinate() {}

  // necessary to use in unordered_map as key: equality methodology
  bool operator==(const Coordinate& other) const {
    return this->row == other.row && this->col == other.col;
  }

 private:
  void copy(const Coordinate& other) {
    this->row = other.row;
    this->col = other.col;
  }
};

// necessary to use in unordered_map as key: hashing methodology
namespace std {
  template<> struct hash<Coordinate> {
    size_t operator()(const Coordinate& c) const {
      return c.row << 5 || c.col;
    }
  };
}

ostream& operator<<(ostream& os, const Coordinate& c) {
  os << "(" << c.row << "," << c.col << ")";
  return os;
}

// directions: north, south, west, east
vector<int> row_diffs = {-1, 1, 0, 0};
vector<int> col_diffs = {0, 0, -1, 1};

class Grid {
 public:
  Grid(const int& num_rows, const int& num_cols)
    : kRows(static_cast<size_t>(num_rows)), kColumns(static_cast<size_t>(num_cols))
  {
    this->occupied_ = vector<vector<bool>>(kRows, vector<bool>(kColumns, false));
    
    this->occupied_row_min_ = kRows;
    this->occupied_row_max_ = 0;
    this->occupied_col_min_ = kColumns;
    this->occupied_col_max_ = 0;

    this->curr_dir_ = 0;
  }

  Coordinate Center() const {
    return Coordinate(kRows/2, kColumns/2);
  }

  void MarkOccupied(const Coordinate& c) {
    this->occupied_[c.row][c.col] = true;

    this->occupied_row_min_ = min(this->occupied_row_min_, c.row);
    this->occupied_row_max_ = max(this->occupied_row_max_, c.row);
    this->occupied_col_min_ = min(this->occupied_col_min_, c.col);
    this->occupied_col_max_ = max(this->occupied_col_max_, c.col);
  }

  void MarkUnoccupied(const Coordinate& c) {
    this->occupied_[c.row][c.col] = false;
    // NOTE: while technically this may change what the min-max of the occupied spaces are, we expect to grow to that size again
  }

  size_t GetStartingDirectionIndex() const {
    return this->curr_dir_;
  }

  bool AnyNeighboringCellsInDirectionOccupied(const Coordinate& c, size_t dir) const {
    vector<int> row_variations, col_variations;
    if (row_diffs[dir] == 0) {
      // direction is east or west
      row_variations = row_diffs;
      col_variations = vector<int>(col_diffs.size(), col_diffs[dir]);
    } else {
      // direction is north or south
      row_variations = vector<int>(row_diffs.size(), row_diffs[dir]);
      col_variations = col_diffs;
    }

    // note: orthogonal coordinate is checked twice
    for (size_t var = 0; var < row_variations.size(); ++var) {
      Coordinate neighbor(c.row + row_variations[var], c.col + col_variations[var]);

      if (this->occupied_[neighbor.row][neighbor.col]) {
        return true;
      }
    }

    return false;
  }

  // returns true if any of the 8 neighboring cells to (r,c) are occupied
  bool AnyNeighboringCellsOccupied(const Coordinate& c) const {
    // note: diagonal coordinates are checked twice
    for (size_t dir = 0; dir < row_diffs.size(); ++dir) {
      if (this->AnyNeighboringCellsInDirectionOccupied(c, dir)) {
        return true;
      }
    }

    return false;
  }

  size_t SmallestRectangleAroundOccupations() const {
    return (this->occupied_row_max_ - this->occupied_row_min_ + 1) *
      (this->occupied_col_max_ - this->occupied_col_min_ + 1);
  }

  void NextRound() {
    this->curr_dir_ = (this->curr_dir_ + 1) % row_diffs.size();
  }

 private:
  size_t kRows, kColumns;
  vector<vector<bool>> occupied_;

  // keep track of the limits of the occupied spaces
  size_t occupied_row_min_, occupied_row_max_, occupied_col_min_, occupied_col_max_;

  size_t curr_dir_;
};

class Elf {
 public:
  Elf(const Coordinate& c) {
    this->pos_ = c;
  }

  Coordinate CurrentPosition() const {
    return this->pos_;
  }

  // returns false if no proposal is made
  pair<Coordinate,bool> CalculateProposal(const Grid& grid) {
    // If no other Elves are in one of those eight positions, the Elf does not do anything during this round
    if (!grid.AnyNeighboringCellsOccupied(this->pos_)) {
      return {this->pos_, false};
    }

    // Otherwise, the Elf looks in each of four directions in the following order...
    size_t starting_dir = grid.GetStartingDirectionIndex();
    for (size_t dir_change = 0; dir_change < row_diffs.size(); ++dir_change) {
      size_t dir = (starting_dir + dir_change) % row_diffs.size();

      // ... and proposes moving one step in the first valid direction
      if (!grid.AnyNeighboringCellsInDirectionOccupied(this->pos_, dir)) {
        this->proposed_pos_ = Coordinate(this->pos_.row + row_diffs[dir], this->pos_.col + col_diffs[dir]);

        return {this->proposed_pos_, true};
      }
    }

    // cannot move
    return {this->pos_, false};
  }

  void MoveAsProposed() {
    this->pos_ = this->proposed_pos_;
  }

 private:
  Coordinate pos_, proposed_pos_;
};

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("input.txt");

  // count elves
  int elf_count = 0;
  for (auto line:input) {
    for (auto c:line) {
      if (c == '#') {
        ++elf_count;
      }
    }
  }

  // initialize grid large enough to fit all diffused elves (so all elves have a space between them even if lined up in a single row)
  Grid grid(elf_count*2, elf_count*2);
  Coordinate center_of_grid = grid.Center();

  // initialize elf positions within grid
  vector<Elf> elves;
  for (size_t r = 0; r < input.size(); ++r) {
    // center initial input in my grid by lining up its center with the center of my grid
    size_t adjusted_row = r + center_of_grid.row - (input.size()/2);
    for (size_t c = 0; c < input[r].size(); ++c) {
      if (input[r][c] == '#') {
        Coordinate elf_coord(adjusted_row, c + center_of_grid.col - (input[r].size()/2));

        elves.push_back(Elf(elf_coord));
        grid.MarkOccupied(elf_coord);
      }
    }
  }

  // Part 1
  // Simulate the Elves' process and find the smallest rectangle that contains the Elves after 10 rounds
  int round = 1;
  for (; round <= 10; ++round) {
    unordered_map<Coordinate, vector<size_t>> new_cells_proposed;

    // collect all proposals
    for (size_t elf_i = 0; elf_i < elves.size(); ++elf_i) {
      pair<Coordinate,bool> elf_proposed = elves[elf_i].CalculateProposal(grid);
      if (elf_proposed.second) {
        new_cells_proposed[elf_proposed.first].push_back(elf_i);
      }
    }

    // take all non-clashing proposals
    for (auto k_v:new_cells_proposed) {
      if (k_v.second.size() == 1) {
        grid.MarkUnoccupied(elves[k_v.second[0]].CurrentPosition());
        elves[k_v.second[0]].MoveAsProposed();
        grid.MarkOccupied(k_v.first);
      }
    }

    grid.NextRound();
  }
  
  // How many empty ground tiles does that rectangle contain?
  size_t smallest_rectangle_area = grid.SmallestRectangleAroundOccupations();
  cout << endl << "Part 1 answer: " << smallest_rectangle_area - elves.size() << endl;

  // Part 2
  // Figure out where the Elves need to go.
  size_t num_new_cells = 1;
  for (; ; ++round) {
    unordered_map<Coordinate, vector<size_t>> new_cells_proposed;

    // collect all proposals
    for (size_t elf_i = 0; elf_i < elves.size(); ++elf_i) {
      pair<Coordinate,bool> elf_proposed = elves[elf_i].CalculateProposal(grid);
      if (elf_proposed.second) {
        new_cells_proposed[elf_proposed.first].push_back(elf_i);
      }
    }

    // see if we've reached a stable state
    num_new_cells = new_cells_proposed.size();
    if (num_new_cells == 0) {
      break;
    }

    // take all non-clashing proposals
    for (auto k_v:new_cells_proposed) {
      if (k_v.second.size() == 1) {
        grid.MarkUnoccupied(elves[k_v.second[0]].CurrentPosition());
        elves[k_v.second[0]].MoveAsProposed();
        grid.MarkOccupied(k_v.first);
      }
    }

    grid.NextRound();
  }
  
  // What is the number of the first round where no Elf moves?
  cout << endl << "Part 2 answer: " << round << endl;

  return 0;
}
