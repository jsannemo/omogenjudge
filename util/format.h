#pragma once

#include <string>
#include <vector>

namespace omogenexec {

// Parse a string to a long long. Throws if the result is out of range,
// returns 0 if the token could not be converted
long long StringToLL(const std::string& str);

// Split a string into the minimal parts delimited by the delim character.
std::vector<std::string> Split(const std::string &s, char delim);

}
