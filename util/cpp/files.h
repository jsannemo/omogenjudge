#pragma once

#include <string>
#include <vector>

namespace omogen {
namespace util {

// Returns whether a directory exists or not.
bool DirectoryExists(const std::string& path);

// Creates a single directory with mode 755 specified by the path.
// Returns true if the directory already existed.
bool MakeDir(const std::string& path);

// Creates a directory together with all its parents with mode 755.
// Assumes path components are separated by /.
void MakeDirParents(const std::string& path);

// Create a new directory in a location meant for temporary files and returns
// its path.
std::string MakeTempDir();

// Remove a single directory, assuming it is non-empty. This function is
// idempotent.
void RemoveDir(const std::string& path);

// Remove a single directory, assuming it is non-empty. This function is
// idempotent. Ignores errors.
void TryRemoveDir(const std::string& path);

// Destroy an entire directory tree, including any files in it.
void RemoveTree(const std::string& path);

// Overwrite the file given by the path with the given contents.
void WriteToFile(const std::string& path, const std::string& contents);

// Overwrite the file given by the path with the given contents, ignoring
// failures.
void TryWriteToFile(const std::string& path, const std::string& contents);

// Write the given contents to the file descriptor.
void WriteToFd(int fd, const std::string& contents);

// Split the contents of a given file into space-separated tokens.
std::vector<std::string> TokenizeFile(const std::string& path);

// Close all file descriptors except those in a given list.
void CloseFdsExcept(std::vector<int> fdsToKeep);

// Check if a certain file exists and is executable.
bool FileIsExecutable(const std::string& path);

// Read a string containing a given number of bytes from a file descriptor.
// Note that fewer bytes may be returned in case the file descriptor closes.
std::string ReadFromFd(int bytes, int fd);

// Write bytes to a file descriptor.
void WriteToFd(int bytes, unsigned char* ptr, int fd);

// Write an integer in network byte order.
void WriteIntToFd(int value, int fd);

// Read an integer in network byte order. If enough bytes could
// not read for an integer, false is returned. Otherwise, true is returned.
bool ReadIntFromFd(int* val, int fd);

}  // namespace util
}  // namespace omogen
