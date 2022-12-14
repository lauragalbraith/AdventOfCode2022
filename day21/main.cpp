// main.cpp: Laura Galbraith
// Description: Solve Day 21 of Advent of Code 2022
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

  bool ancestor_of_human;

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
    this->ancestor_of_human = false;
  }

  // default constructor
  // note: does NOT create a valid monkey
  Monkey() {
    this->val_finalized = false;
    this->operation = nullptr;
    this->ancestor_of_human = false;
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

// constants in the problem description
string kRootMonkey("root");
string kHuman("humn");

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("input.txt");

  unordered_map<string,Monkey> monkeys;
  for (auto line:input) {
    monkeys[line.substr(0,4)] = Monkey(line);
  }

  // Part 1
  // Traverse the tree to finalize value at root
  stack<string> to_process;
  to_process.push(kRootMonkey);
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

    } else {
      // process children before us
      to_process.push(next_monkey.b_name);
      to_process.push(next_monkey.a_name);
    }
  }

  // What number will the monkey named root yell?
  cout << endl << "Part 1 answer: " << monkeys[kRootMonkey].val << endl;

  // Part 2
  // determine all ancestors of the human
  string monkey_name = kHuman;
  while (monkey_name != kRootMonkey) {
    // find the current monkey's parent
    for (auto m_k_v:monkeys) {
      if (m_k_v.second.a_name == monkey_name || m_k_v.second.b_name == monkey_name) {
        monkey_name = m_k_v.first;
        break;
      }
    }

    monkeys[monkey_name].ancestor_of_human = true;
  }

  // track the result we want to get
  Monkey& curr_monkey = monkeys[monkeys[kRootMonkey].a_name];
  long long int desired_val = monkeys[monkeys[kRootMonkey].b_name].val;

  if (!curr_monkey.ancestor_of_human) {
    desired_val = curr_monkey.val;
    curr_monkey = monkeys[monkeys[kRootMonkey].b_name];
  }

  // walk down from root setting up the answers we want
  while (curr_monkey.name != kHuman) {
    if (monkeys[curr_monkey.a_name].ancestor_of_human || curr_monkey.a_name == kHuman) {
      monkey_operation reverse_operation;
      if (curr_monkey.operation == &Add) {
        reverse_operation = &Subtract;
      } else if (curr_monkey.operation == &Multiply) {
        reverse_operation = &Divide;
      } else if (curr_monkey.operation == &Subtract) {
        reverse_operation = &Add;
      } else {
        reverse_operation = &Multiply;
      }

      desired_val = reverse_operation(
        desired_val,
        monkeys[curr_monkey.b_name].val);

      curr_monkey = monkeys[curr_monkey.a_name];
    } else {
      if (curr_monkey.operation == &Add) {
        desired_val = Subtract(desired_val, monkeys[curr_monkey.a_name].val);
      } else if (curr_monkey.operation == &Multiply) {
        desired_val = Divide(desired_val, monkeys[curr_monkey.a_name].val);
      } else if (curr_monkey.operation == &Subtract) {
        desired_val = Subtract(monkeys[curr_monkey.a_name].val, desired_val);
      } else {
        desired_val = Divide(monkeys[curr_monkey.a_name].val, desired_val);
      }

      curr_monkey = monkeys[curr_monkey.b_name];
    }
  }

  // What number do you yell to pass root's equality test?
  cout << endl << "Part 2 answer: " << desired_val << endl;

  return 0;
}
