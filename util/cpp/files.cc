#include "util/cpp/files.h"

#include <dirent.h>
#include <fcntl.h>
#include <ftw.h>
#include <sys/stat.h>
#include <unistd.h>

#include <algorithm>
#include <cassert>
#include <fstream>

#include "glog/logging.h"
#include "glog/raw_logging.h"

using std::string;

namespace omogen {
namespace util {

bool DirectoryExists(const std::string& path) {
  struct stat sb;
  if (stat(path.c_str(), &sb) == -1) {
    if (errno == ENOENT || errno == ENOTDIR) {
      return false;
    }
    RAW_LOG(FATAL, "Could not check if directory exists");
  }
  return S_ISDIR(sb.st_mode);
}

bool MakeDir(const std::string& path) {
  RAW_VLOG(2, "Making path %s", path.c_str());
  if (mkdir(path.c_str(), 0755) == -1) {
    if (errno == EEXIST) {
      return true;
    }
    RAW_LOG(FATAL, "Could not create path: %s", path.c_str());
  }
  if (!DirectoryExists(path)) {
    RAW_LOG(FATAL, "Tried to make directory %s but does not exist",
            path.c_str());
  }
  return false;
}

void MakeDirParents(const std::string& path) {
  if (path.empty()) {
    return;
  }
  RAW_VLOG(2, "Making recursive path %s", path.c_str());
  // Create all directories on a path by repeatedly finding the next directory
  // separator and creating the path up to this location.
  size_t at = 0;
  while (at < path.size()) {
    size_t next = path.find('/', at + 1);
    if (next == std::string::npos) {
      next = path.size();
    }
    MakeDir(path.substr(0, next));
    at = next;
  }
}

std::string MakeTempDir() {
  char rootPath[] = "/tmp/omogencontainXXXXXX";
  RAW_CHECK(mkdtemp(rootPath) != nullptr,
            "Could not create temporary directory");
  return std::string(rootPath);
}

void RemoveDir(const std::string& path) {
  RAW_VLOG(2, "Removing path %s", path.c_str());
  RAW_CHECK(rmdir(path.c_str()) != -1 || errno == ENOENT,
            "Could not remove directory");
}

void TryRemoveDir(const std::string& path) {
  RAW_VLOG(2, "Removing path %s", path.c_str());
  rmdir(path.c_str());
}

static int removeTree0(const char* filePath, const struct stat* statData,
                       int typeflag, struct FTW* ftwbuf) {
  RAW_VLOG(3, "Recursive remove path %s", filePath);
  if (S_ISDIR(statData->st_mode)) {
    RAW_CHECK(rmdir(filePath) != -1, "Could not remove folder");
  } else {
    RAW_CHECK(unlink(filePath) != -1, "Could not remove file");
  }
  return FTW_CONTINUE;
}

void RemoveTree(const std::string& directoryPath) {
  RAW_VLOG(2, "Removing tree %s", directoryPath.c_str());
  RAW_CHECK(nftw(directoryPath.c_str(), removeTree0, 32,
                 FTW_MOUNT | FTW_PHYS | FTW_DEPTH) != -1,
            "Could not tree walk path to remove");
}

void WriteToFile(const std::string& path, const std::string& contents) {
  RAW_VLOG(2, "Writing to %s", path.c_str());
  std::ofstream ofs(path);
  if (!(ofs << contents)) {
    RAW_LOG(FATAL, "Failed writing to %s", path.c_str());
  }
}

void WriteToFd(int fd, const std::string& contents) {
  size_t at = 0;
  while (at != contents.size()) {
    int wrote = write(fd, contents.data() + at, contents.size() - at);
    if (wrote == -1) {
      RAW_CHECK(errno == EINTR, "Write failed with something other than EINTR");
    } else {
      at += wrote;
    }
  }
}

std::vector<std::string> TokenizeFile(const std::string& path) {
  std::ifstream ifs(path);
  std::vector<std::string> tokens;
  std::string tok;
  while (ifs >> tok) {
    tokens.push_back(tok);
  }
  RAW_CHECK(ifs.eof(), "Failed tokenizing from file");
  return tokens;
}

void CloseFdsExcept(std::vector<int> fdsToKeep) {
  DIR* fdDir = opendir("/proc/self/fd");
  RAW_CHECK(fdDir != nullptr, "Could not open /proc/self/fd");
  // Do not accidentally close the fd directory file descriptor.
  fdsToKeep.push_back(dirfd(fdDir));
  while (true) {
    errno = 0;
    struct dirent* entry = readdir(fdDir);
    if (entry == nullptr) {
      RAW_CHECK(errno == 0, "Could not read next file descriptor");
      break;
    }

    errno = 0;
    int fd = strtol(entry->d_name, nullptr, 10);
    if (errno != 0) {
      RAW_LOG(ERROR, "Ignoring invalid fd entry: %s", entry->d_name);
    } else if (std::find(fdsToKeep.begin(), fdsToKeep.end(), fd) ==
               fdsToKeep.end()) {
      RAW_CHECK(close(fd) != -1, "Closing file descriptor %d failed");
    }
  }
  RAW_CHECK(closedir(fdDir) != -1, "Closing file descriptor folder failed");
}

bool FileIsExecutable(const std::string& path) {
  struct stat sb;
  if (stat(path.c_str(), &sb) == -1) {
    if (errno == EACCES || errno == ENAMETOOLONG || errno == ENOTDIR ||
        errno == ENOENT) {
      return false;
    }
    RAW_LOG(FATAL, "Stating file failed");
  }
  return S_ISREG(sb.st_mode) && (S_IXUSR & sb.st_mode);
}

void WriteIntToFd(int value, int fd) {
  unsigned char buf[4];
  for (int i = 0; i < 4; i++) {
    buf[3 - i] = value & 0xff;
    value >>= 8;
  }
  WriteToFd(4, buf, fd);
}

bool ReadIntFromFd(int* val, int fd) {
  std::string ret = ReadFromFd(4, fd);
  if (ret.size() != 4) {
    return false;
  }
  *val = 0;
  for (unsigned char p : ret) {
    *val = (*val << 8) | p;
  }
  return true;
}

void WriteToFd(int bytes, unsigned char* ptr, int fd) {
  int at = 0;
  while (at != bytes) {
    int r = write(fd, ptr + at, bytes - at);
    PCHECK(r != -1 || errno == EINTR) << "Failed reading return value";
    if (errno == EINTR) {
      continue;
    }
    at += r;
  }
}

std::string ReadFromFd(int bytes, int fd) {
  char buf[bytes];
  int at = 0;
  while (at != bytes) {
    int r = read(fd, buf + at, bytes - at);
    RAW_CHECK(r != -1 || errno == EINTR, "Failed reading return value");
    if (errno == EINTR) {
      continue;
    }
    if (r == 0) {
      break;
    }
    at += r;
  }
  return string(buf, buf + at);
}

}  // namespace util
}  // namespace omogen
