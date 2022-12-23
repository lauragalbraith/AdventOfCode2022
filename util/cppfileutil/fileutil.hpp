// fileutil.hpp: Laura Galbraith
// Description: header file for useful file-related functions for The Advent Of Code 2022

#ifndef FILE_UTIL_H
#define FILE_UTIL_H

#include <tuple>
#include <vector>
#include <string>

// ReadLinesFromFile takes in the name of a file and returns each line as a string
// Runtime complexity: linear in the size of the file
std::vector<std::string> ReadLinesFromFile(const std::string file_name);

#endif // FILE_UTIL_H