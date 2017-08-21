#include <algorithm>
#include <cassert>
#include <dirent.h>
#include <fcntl.h>
#include <ftw.h>
#include <fstream>
#include <sys/stat.h>
#include <unistd.h>

#include "error.h"
#include "log.h"
#include "files.h"

using std::endl;
using std::find;
using std::ifstream;
using std::ofstream;
using std::string;
using std::vector;

namespace omogenexec {

bool DirectoryExists(const string& path) {
    struct stat sb;
    if (stat(path.c_str(), &sb) == -1) {
        if (errno == ENOENT || errno == ENOTDIR) {
            return false;
        }
        OE_FATAL("stat");
    }
    return S_ISDIR(sb.st_mode);
}

bool MakeDir(const string& path) {
    if (mkdir(path.c_str(), 0755) == -1) {
        if (errno == EEXIST) {
            return true;
        }
        OE_FATAL("mkdir");
    }
    return false;
}

void MakeDirParents(const string& path) {
    if (path.empty()) {
        return;
    }
    // Create all directories on a path by repeatedly finding the next directory separator
    // and creating the path up to this location.
    size_t at = 0;
    while (at < path.size()) {
        size_t next = path.find('/', at + 1);
        if (next == string::npos) {
            next = path.size();
        }
        MakeDir(path.substr(0, next));
        at = next;
    }
}

string MakeTempDir() {
    char rootPath[] = "/tmp/omogencontainXXXXXX";
    if (mkdtemp(rootPath) == nullptr) {
        OE_FATAL("mkdtemp");
    }
    return string(rootPath);
}

void RemoveDir(const string& path) {
    if (rmdir(path.c_str()) == -1) {
        if (errno != ENOENT) {
            OE_FATAL("mkdir");
        }
    }
}

static int removeTree0(const char *filePath, const struct stat *statData, int typeflag, struct FTW *ftwbuf) {
    if (S_ISDIR(statData->st_mode)) {
        if (rmdir(filePath) == -1) {
            OE_FATAL("rmdir");
        }
    } else {
        if (unlink(filePath) == -1) {
            OE_FATAL("unlink");
        }
    }
    return FTW_CONTINUE;
}

void RemoveTree(const string& directoryPath) {
    if (nftw(directoryPath.c_str(), removeTree0, 32, FTW_MOUNT | FTW_PHYS | FTW_DEPTH) == -1) {
        OE_FATAL("nftw");
    }
}

void WriteToFile(const string& path, const string& contents) {
    ofstream ofs(path);
    if (!(ofs << contents)) {
        OE_LOG(FATAL) << "Failed writing to " << path << endl;
        OE_CRASH();
    }
}

vector<string> TokenizeFile(const string& path) {
    ifstream ifs(path);
    vector<string> tokens;
    string tok;
    while (ifs >> tok) {
        tokens.push_back(tok);
    }
    if (!ifs.eof()) {
        OE_LOG(FATAL) << "Failed reading from " << path << endl;
        OE_CRASH();
    }
    return tokens;
}

void CloseFdsExcept(vector<int> fdsToKeep) {
    DIR *fdDir = opendir("/proc/self/fd");
    if (fdDir == nullptr) {
        OE_FATAL("opendir");
    }
    // Do not accidentally close the fd directory file descriptor.
    fdsToKeep.push_back(dirfd(fdDir));
    while (true) {
        errno = 0;
        struct dirent *entry = readdir(fdDir);
        if (entry == nullptr) {
            if (errno != 0) {
                OE_FATAL("readdir");
            }
            break;
        }

        errno = 0;
        int fd = strtol(entry->d_name, nullptr, 10);
        if (errno != 0) {
            OE_LOG(WARN) << "Ignoring invalid fd entry: " << entry->d_name << endl;
        } else if (find(fdsToKeep.begin(), fdsToKeep.end(), fd) == fdsToKeep.end()) {
            if (close(fd) == -1) {
                OE_FATAL("close");
            }
        }
    }
    if (closedir(fdDir) == -1) {
        OE_FATAL("closedir");
    }
}

}
