// main.cpp: Laura Galbraith
// Description: Solve Day 25 of Advent of Code 2022
// Full of Hot Air
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <cmath>  // llabs
#include <iostream>  // cout, endl
#include <string>  // string
#include <unordered_map>  // unordered_map
#include <vector>  // vector

using namespace std;

unordered_map<char,long long int> snafu_digit_value = {
  {'=', -2},
  {'-', -1},
  {'0', 0},
  {'1', 1},
  {'2', 2},
};

long long unsigned int FromSNAFU(const string& snafu) {
  long long int val = 0;
  long long int power_of_five = 1;

  for (size_t i = snafu.size()-1; i < snafu.size(); --i, power_of_five *= 5) {
    long long int digit_val = power_of_five * snafu_digit_value[snafu[i]];
    val += digit_val;
  }

  return static_cast<long long unsigned int>(val);
}

string ToSNAFU(const long long unsigned int& x) {
  // find largest power of 5 to represent
  size_t largest_digit_idx = 0;
  long long int power_of_five = 1;

  while (x > static_cast<long long unsigned int>(2*power_of_five)) {
    ++largest_digit_idx;
    power_of_five *= 5;
  }

  long long int remaining = static_cast<long long int>(x);
  string ans(largest_digit_idx+1, '0');

  // calculate each digit
  // unintuitively, the largest digit will sit at ans[0]
  for (size_t i = 0; i <= largest_digit_idx; ++i, power_of_five /= 5) {
    long long int digit_val = 0;

    if (remaining > 0) {
      // subtract until we've on the closer side of the number-line-teeter-totter
      // i.e. we are getting as close to 0 as possible by subtracting the power of five
      while (digit_val < 2 && !((remaining-power_of_five) < 0 && llabs(remaining-power_of_five) > remaining)) {
        remaining -= power_of_five;
        ++digit_val;
      }
    } else if (remaining < 0) {
      // add until we're on the closer side of 0
      while (digit_val > -2 && !((remaining+power_of_five) > 0 && (remaining+power_of_five) > llabs(remaining))) {
        remaining += power_of_five;
        --digit_val;
      }
    }

    // find digit
    for (auto d_v:snafu_digit_value) {
      if (d_v.second == digit_val) {
        ans[i] = d_v.first;
        break;
      }
    }
  }

  return ans;
}

int main() {
  // Parse input: fuel requirements
  vector<string> input = ReadLinesFromFile("input.txt");

  // Part 1
  long long unsigned int fuel_req_sum = 0;
  for (auto req:input) {
    long long unsigned int val = FromSNAFU(req);
    fuel_req_sum += val;
  }

  // What SNAFU number do you supply to Bob's console?
  string fuel_req_sum_snafu = ToSNAFU(fuel_req_sum);
  cout << endl << "Part 1 answer: " << fuel_req_sum_snafu << endl;

  return 0;
}
