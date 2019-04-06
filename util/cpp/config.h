#pragma once

#include <string>

namespace omogen {
namespace util {

// Returns the integer config value with the given key.
// If there is no such file or key, or the value is not an int,
// the program will crash.
int ReadConfigInt(const std::string& filePath, const std::string& key);

}  // namespace util
}  // namespace omogen
