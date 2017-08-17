#include <string>

#include "proto/omogenexec.pb.h"

using std::string;

// Destroy an entire directory tree
void DestroyDirectory(const string& path);
string CreateTemporaryRoot();

class Chroot {

    // The path of the new root
    string rootfs;
    void addDefaultRules();

public:
    Chroot(const string& path);
    void AddDirectoryRule(const DirectoryRule& rule);
    // Changes the root of the current file system to the one given as the path
    void SetRoot();
};
