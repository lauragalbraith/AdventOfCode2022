#include <fstream>  // ifstream
#include <stdexcept>
#include <string>  // string
#include <vector>  // vector

using namespace std;

vector<string> ReadLinesFromFile(const string file_name) {
  // Read in file
  ifstream f(file_name);
  vector<string> contents;
  if (f.is_open()) {
    // Collect lines from file
    string line;
    while (getline(f, line)) {
      contents.push_back(line);
    }
  }
  else {
    throw invalid_argument("unable to open file");
  }

  return contents;
}