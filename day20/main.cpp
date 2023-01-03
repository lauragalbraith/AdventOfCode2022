// main.cpp: Laura Galbraith
// Description: Solve Day 20 of Advent of Code 2022
// Grove Positioning System
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <cmath>  // abs
#include <iostream>  // cout, endl
#include <string>  // string, stoi
#include <unordered_map>  // unordered_map
#include <vector>  // vector

using namespace std;

class CircularDoublyLinkedList {
 public:
  CircularDoublyLinkedList(const vector<int>& vals) {
    // initialize member variables
    this->head_ = nullptr;
    this->initial_order_ = vector<Element*>(vals.size(), nullptr);

    Element* zero_element = nullptr;

    // create elements from vals, maintaining order
    for (size_t i = 0; i < vals.size(); ++i) {
      // create and store element
      Element* elem = new Element(vals[i]);
      this->initial_order_[i] = elem;

      // update element behind us
      if (i > 0) {
        this->initial_order_[i-1]->next = elem;
        elem->prev = this->initial_order_[i-1];
      }

      if (vals[i] == 0) {
        zero_element = elem;
      }
    }

    // finish circle of pointers
    this->initial_order_[vals.size()-1]->next = this->initial_order_[0];
    this->initial_order_[0]->prev = this->initial_order_[vals.size()-1];

    // set "head" to be the zero element
    this->head_ = zero_element;
  }

  // copy constructor
  CircularDoublyLinkedList(const CircularDoublyLinkedList& other) {
    this->copy(other);
  }

  // assignment operator
  CircularDoublyLinkedList& operator=(const CircularDoublyLinkedList& other) {
    if (this != &other) {
      this->clear();
      this->copy(other);
    }

    return *this;
  }

  // deconstructor
  ~CircularDoublyLinkedList() {
    this->clear();
  }

  // "mix" the file by moving each number forward or backward in the file a number of positions equal to the value of the number being moved
  void Mix() {
    // For each value in the initial order...
    for (auto elem:this->initial_order_) {      
      // can ignore times looping completely around the other len-1 elements
      const int spaces_to_move = abs(elem->val) % static_cast<int>(this->initial_order_.size() - 1);
      int spaces_moved = 0;

      // effectively remove elem itself from list before trying to loop
      elem->prev->next = elem->next;
      elem->next->prev = elem->prev;

      // ... Move it within the list according to its value
      Element* destination = elem->prev;  // destination is the list element whose ->next value should point to elem
      while (spaces_moved < spaces_to_move) {
        // direction is based on sign of element's value
        if (elem->val > 0) {
          destination = destination->next;
        } else {
          destination = destination->prev;
        }

        ++spaces_moved;
      }

      // update next,prev pointers of the source, destination, and elem
      // cout << "DEBUG: element " << elem->val << " would be moved to just in front of " << destination->val << endl;

      elem->next = destination->next;
      elem->next->prev = elem;

      destination->next = elem;
      elem->prev = destination;
    }
  }

  // function to solve Part 1
  int GroveCoordinates() const {
    int sum = 0;

    Element* curr = this->head_;
    for (int num_to_hit = 0; num_to_hit < 3; ++num_to_hit) {
      // moves 1000 the first time, then 2000, then 3000
      size_t numbers_to_see = 1000 % this->initial_order_.size();

      for (size_t num_i = 0; num_i < numbers_to_see; ++num_i) {
        curr = curr->next;
      }

      // cout << "DEBUG: thousandth number is " << curr->val << endl;

      // add value to sum
      sum += curr->val;
    }

    return sum;
  }

  void PrintListFromZero() const {
    Element* curr = this->head_;
    do {
      cout << curr->val << ", ";
      curr = curr->next;
    } while (curr != this->head_);
    cout << endl;
  }

 private:
  // element in the list
  struct Element {
    // pointers ahead and behind in the list
    Element *next, *prev;
    // value stored by the element
    int val;

    // constructor
    Element(const int& value) {
      this->val = value;
      this->next = nullptr;
      this->prev = nullptr;
    }
  };

  // keep track of the "head" element
  Element* head_;

  // keep track of the initial order
  vector<Element*> initial_order_;

  void copy(const CircularDoublyLinkedList& other) {
    // track which elements may be used from initial order
    unordered_map<Element*, size_t> other_elements_to_order;

    // copy initial order
    this->initial_order_.resize(other.initial_order_.size());
    for (size_t i = 0; i < other.initial_order_.size(); ++i) {
      other_elements_to_order[other.initial_order_[i]] = i;

      this->initial_order_[i] = new Element(other.initial_order_[i]->val);
    }

    // arrange those elements in order from other's head_
    Element *other_curr = other.head_, *curr = nullptr;

    do {
      // save old current to use to point to
      Element* previous = curr;

      // Element has already been created; just pull it
      size_t idx = other_elements_to_order[other_curr];
      curr = this->initial_order_[idx];

      // update pointers
      if (previous != nullptr) {
        previous->next = curr;
      }
      curr->prev = previous;

      // save "head" of list to match other
      if (other_curr == other.head_) {
        this->head_ = curr;
      }

      // move on
      other_curr = other_curr->next;
    } while (other_curr != other.head_);

    // update list to be circular
    this->head_->prev = curr;
    curr->next = this->head_;
  }

  void clear() {
    // Treat head_ as the source of truth, not initial_order_

    // clear initial_order_
    this->initial_order_.clear();

    // clear head_, freeing memory
    Element* curr = this->head_;
    while (curr != nullptr) {
      Element* next = curr->next != curr ? curr->next : nullptr;

      // reset previous element's pointer
      if (curr->prev != curr) {
        curr->prev->next = curr->next;
      }

      // reset next element's pointer
      if (curr->next != curr) {
        curr->next->prev = curr->prev;
      }

      // update head_ (not strictly necessary - if multithreaded, will still need to introduce synchronization mechanisms)
      this->head_ = next;

      // delete the element
      delete curr;

      // move on
      curr = next;
    }

    // reset head_
    this->head_ = nullptr;
  }
};

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("input.txt");

  vector<int> values(input.size(), 0);
  for (size_t i = 0; i < input.size(); ++i) {
    values[i] = stoi(input[i]);
  }

  CircularDoublyLinkedList list = CircularDoublyLinkedList(values);

  // Part 1
  // Mix the file exactly once
  list.Mix();
  // list.PrintListFromZero();

  // What is the sum of the three numbers that form the grove coordinates?
  cout << endl << "Part 1 answer: " << list.GroveCoordinates() << endl;

  /*
  IDEAS

  - (circular?) doubly-linked list
  - could use vector, but that is less time-efficient since it has to copy following elements to squeeze stuff in; has the benefit of being easier to code
  - have some way to keep track of the initial order, since that is maintained for each mix operation (pointers to elements in the linked list?)
  - note: elements getting moved to before-the-first-element end up after-the-last-element (at least in the negative direction) (this shouldn't effect the final result, since the list "starts" at 0)
  - note: must be able to handle 0, which does not move
  - could consider 0 as the "head" of the list, since that's where the answer comes from
  - when finding the 1000th number, mod it by the length of the list, then move forward that many (not as relevant for the full input, which is 5000 elements)
  - since elements can be larger than the list length, mod them by the list length before moving them
  - note: there are duplicate elements, unlike the example input
  */

  // Part 2
  // TODO
  cout << endl << "Part 2 answer: " << endl;

  return 0;
}
