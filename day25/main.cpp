// main.cpp: Laura Galbraith
// Description: Solve Day 25 of Advent of Code 2022
// Full of Hot Air
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <cmath>  // llabs
#include <iostream>  // cout, endl
#include <string>  // string, stoi, to_string
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

    // cout << "DEBUG: for index " << i << ": '" << snafu[i] << "' we are adding " << digit_val << endl;

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

  // long long int prev_multiple = 5;  // TODO FINALLY remove if unused

  // calculate each digit
  // unintuitively, the largest digit will sit at ans[0]
  for (size_t i = 0; i <= largest_digit_idx; ++i, power_of_five /= 5) {
    long long int digit_val = 0;

    if (remaining > 0) {
      // subtract until we've on the closer side of the number-line-teeter-totter
      // i.e. we are getting as close to 0 as possible by subtracting the power of five
      while (digit_val < 2 && !((remaining-power_of_five) < 0 && llabs(remaining-power_of_five) > remaining)) {
        // cout << "DEBUG: remaining is +" << remaining << " and we could get closer to zero, so we're including another " << power_of_five << endl;
        remaining -= power_of_five;
        ++digit_val;
      }
    } else if (remaining < 0) {
      // add until we're on the closer side of 0
      while (digit_val > -2 && !((remaining+power_of_five) > 0 && (remaining+power_of_five) > llabs(remaining))) {
        // cout << "DEBUG: remaining is " << remaining << " and we could get closer to zero, so we're detracting another " << power_of_five << endl;
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

  /*
  IDEAS
  - base 5 instead of base 10
  - 0,1,2,3,4 -> =,-,0,1,2
  - "So, because ten (in normal numbers) is two fives and no ones, in SNAFU it is written 20. Since eight (in normal numbers) is two fives minus two ones, it is written 2=."
    - 

  - parsing from input to int should be easy (just apply the formula on successive powers of two), but maybe I could manually translate from int to our final answer

  - TODO FIRST write a SNAFU conversion function and test it on the SNAFU brochure values
  */

  /*cout << "DEBUG: testing:" << endl;
  vector<string> test_snafu = {
    "1",
    "2",
    "1=",
    "1-",
    "10",
    "11",
    "12",
    "2=",
    "2-",
    "20",
    "1=0",
    "1-0",
    "1=11-2",
    "1-0---0",
    "1121-1110-1=0"};

  for (auto snafu:test_snafu) {
    long long unsigned int val = FromSNAFU(snafu);
    // cout << "DEBUG: decimal value of snafu '" << snafu << "' is " << val << endl;

    string back_to_snafu = ToSNAFU(val);
    // cout << "DEBUG: original snafu value:" << snafu << " which is decimal " << val << " is back to snafu:" << back_to_snafu << endl;
  }*/
  // TODO FINALLY remove all comment blocks

  // Part 1
  long long unsigned int fuel_req_sum = 0;
  for (auto req:input) {
    long long unsigned int val = FromSNAFU(req);
    // cout << "DEBUG: val from input " << req << " is " << val << endl;
    fuel_req_sum += val;
  }

  // What SNAFU number do you supply to Bob's console?
  string fuel_req_sum_snafu = ToSNAFU(fuel_req_sum);
  cout << endl << "Part 1 answer: " << fuel_req_sum_snafu << endl;

  // Part 2
  // TODO
  cout << endl << "Part 2 answer: " << endl;

  return 0;
}
