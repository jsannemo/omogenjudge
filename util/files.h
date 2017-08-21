#pragma once

#include <string>
#include <vector>

namespace omogenexec {

// Returns whether a directory exists or not.
bool DirectoryExists(const std::string& path);

// Creates a single directory with mode 755 specified by the path.
// Returns true if the directory already existed.
bool MakeDir(const std::string& path);

// Creates a directory together with all its parents with mode 755.
void MakeDirParents(const std::string& path);

// Create a new directory in a location meant for temporary files and returns its path.
std::string MakeTempDir();

// Remove a single directory, assuming it is non-empty.
void RemoveDir(const std::string& path);

// Destroy an entire directory tree, including any files in it.
void RemoveTree(const std::string& path);

// Overwrite the file given by the path with the given contents.
void WriteToFile(const std::string& path, const std::string& contents);

// Split the contents of a given file into space-separated tokens.
std::vector<std::string> TokenizeFile(const std::string& path);

// Close all file descriptors except those in a given list.
void CloseFdsExcept(std::vector<int> fdsToKeep);

}
