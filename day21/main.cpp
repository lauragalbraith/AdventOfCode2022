// main.cpp: Laura Galbraith
// Description: Solve Day 22 of Advent of Code 2022
// Monkey Math
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <iostream>  // cout, endl, ostream
#include <regex>  // regex, regex_match, smatch
#include <stack>  // stack
#include <string>  // string, stoi
#include <unordered_map>  // unordered_map
#include <vector>  // vector

using namespace std;

typedef long long int (*monkey_operation)(long long int a, long long int b);

long long int Add(long long int a, long long int b) { return a + b; }
long long int Multiply(long long int a, long long int b) { return a * b; }
long long int Subtract(long long int a, long long int b) { return a - b; }
long long int Divide(long long int a, long long int b) { return a / b; }

struct Monkey {
  string name;
  bool val_finalized;
  long long int val;

  monkey_operation operation;
  string a_name, b_name;

  // constructor to parse a monkey
  Monkey(const string& input) {
    smatch sm;
    if (regex_match(input, sm, number_monkey_rgx)) {
      // number monkey
      this->val_finalized = true;
      this->val = static_cast<long long int>(stoi(sm[2]));

      this->operation = nullptr;
      this->a_name = "";
      this->b_name = "";
    } else {
      // operation monkey
      regex_match(input, sm, operation_monkey_rgx);

      this->val_finalized = false;
      this->a_name = sm[2];
      this->b_name = sm[4];

      // determine operation
      if (sm[3] == "+") {
        this->operation = &Add;
      } else if (sm[3] == "*") {
        this->operation = &Multiply;
      } else if (sm[3] == "-") {
        this->operation = &Subtract;
      } else {
        this->operation = &Divide;
      }
    }

    this->name = sm[1];
  }

  // default constructor
  // note: does NOT create a valid monkey
  Monkey() {
    this->val_finalized = false;
    this->operation = nullptr;
  }

  // copy constructor
  Monkey(const Monkey& other) {
    this->copy(other);
  }

  // copy assignment operator
  Monkey& operator=(const Monkey& other) {
    if (this != &other) {
      this->clear();
      this->copy(other);
    }
    
    return *this;
  }

  // destructor
  ~Monkey() {
    this->clear();
  }

  static regex number_monkey_rgx, operation_monkey_rgx;

 private:
  void copy(const Monkey& other) {
    this->name = other.name;
    this->val_finalized = other.val_finalized;
    this->val = other.val;
    this->operation = other.operation;
    this->a_name = other.a_name;
    this->b_name = other.b_name;
  }

  void clear() {
    // no-op
  }
};

ostream& operator<<(ostream& os, const Monkey& m) {
  os << "Monkey " << m.name << " ";
  if (m.val_finalized) {
    os << "has value " << m.val;
  } else {
    os << "is waiting on " << m.a_name << " and " << m.b_name;
  }

  return os;
}

// note: no negative numbers in input
regex Monkey::number_monkey_rgx = regex("^([^: ]+): (\\d+)$");
regex Monkey::operation_monkey_rgx = regex("^([^: ]+): ([^: ]+) ([+\\-*\\/]) ([^: ]+)$");

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("input.txt");

  unordered_map<string,Monkey> monkeys;
  for (auto line:input) {
    monkeys[line.substr(0,4)] = Monkey(line);
  }

  /*for (auto k_v:monkeys) {
    cout << "DEBUG: " << k_v.second << " whose name better match " << k_v.first << endl;
  }*/

  /*
  IDEAS

  - binary tree, where each monkey is a node, who may have no children (be a number) or a math operation monkey (can be *,+,-,/)
  - could supplement with a map from monkey name to its node pointer if it hasn't been located in the tree yet
  - OR - better yet not have an actual pointer-based tree implementation but instead have each node store the monkey name, which can then be found in the map
  */

  // Part 1
  // Traverse the tree to finalize value at root
  stack<string> to_process;
  to_process.push("root");
  while (!to_process.empty()) {
    Monkey& next_monkey = monkeys[to_process.top()];

    // check if this monkey needs any more processing
    if (next_monkey.val_finalized) {
      to_process.pop();

    // check if its dependencies need any more processing
    } else if (monkeys[next_monkey.a_name].val_finalized && monkeys[next_monkey.b_name].val_finalized) {
      next_monkey.val = next_monkey.operation(monkeys[next_monkey.a_name].val, monkeys[next_monkey.b_name].val);
      next_monkey.val_finalized = true;
      to_process.pop();

    //   cout << "DEBUG: monkey " << next_monkey.name << " can now be calculated; its value is now " << next_monkey.val << endl;

    } else {
      // process children before us
      to_process.push(next_monkey.b_name);
      to_process.push(next_monkey.a_name);

    //   cout << "DEBUG: monkeys " << next_monkey.a_name << " and " << next_monkey.b_name << " must be solved before " << next_monkey.name << endl;
    }
  }

  // What number will the monkey named root yell?
  cout << endl << "Part 1 answer: " << monkeys["root"].val << endl;

  // Part 2
  // TODO
  cout << endl << "Part 2 answer: " << endl;

  return 0;
}
