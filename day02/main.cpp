// main.cpp: Laura Galbraith
// Description: Solve Day 2 of Advent of Code 2022
// What would your total score be if everything goes exactly according to your strategy guide?
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

/*
A,X: rock
B,Y: paper
C,Z: scissors

Score = your play (1 rock or 2 paper or 3 scissors) + outcome (0-3-6 loss-draw-win)

There is a B Y line, so the instructions must not represent a guaranteed win

Part 1: interpret results of following strat guide exactly
*/

#include <iostream>  // cout, endl
#include <fstream>  // ifstream, getline
#include <stdexcept>  // invalid_argument
#include <string>  // string

using namespace std;

enum class Play {kRock, kPaper, kScissors, kInvalidPlay};

Play CharToPlay(const char c) {
  switch (c) {
    case 'A':
    case 'X':
      return Play::kRock;
    case 'B':
    case 'Y':
      return Play::kPaper;
    case 'C':
    case 'Z':
      return Play::kScissors;
    default:
      throw invalid_argument("invalid input character");
  }

  return Play::kInvalidPlay;
}

unsigned long int ScoreOfPlay(const Play p) {
  return static_cast<unsigned long int>(p) + 1;
}

int main() {
  unsigned long int score_sum = 0;
  
  // parse input
  ifstream f("input.txt");
  if (f.is_open()) {
    string line;
    while (getline(f, line)) {
      Play opponent_play = CharToPlay(line[0]);
      Play my_play = CharToPlay(line[2]);

      // calculate score from outcome of round
      // wrap difference between plays around
      int outcome = (static_cast<int>(my_play) - static_cast<int>(opponent_play) + static_cast<int>(Play::kInvalidPlay)) % static_cast<int>(Play::kInvalidPlay);
      if (outcome == 0) {
        // draw
        score_sum += 3;
      } else if (outcome == 1) {
        // I win
        score_sum += 6;
      } // else outcome is 2: they win, score increases by 0

      // add value of the play I made
      score_sum += ScoreOfPlay(my_play);
    }
  } else {
    cout << "Unable to open file" << endl;
    return -1;
  }

  cout << "Part 1 answer: " << score_sum << endl;

  // Part 2: interpret second character as end result
  score_sum = 0;

  ifstream f2("input.txt");
  string line;
  while (getline(f2, line)) {
    Play opponent_play = CharToPlay(line[0]);

    // determine my play
    int desired_outcome = static_cast<int>(CharToPlay(line[2]));  // X: lose,0;  Y: draw,1;  Z: win,2

    // same logic as last time, just shifted by one
    int desired_play_value = (desired_outcome - 1 + static_cast<int>(opponent_play) + static_cast<int>(Play::kInvalidPlay)) % static_cast<int>(Play::kInvalidPlay);
    Play desired_play = static_cast<Play>(desired_play_value);

    // calculate score from outcome of round
    score_sum += desired_outcome * 3;

    // add value of the play I made
    score_sum += ScoreOfPlay(desired_play);
  }

  cout << "Part 2 answer: " << score_sum << endl;

  return 0;
}
