#pragma once

#include <string>

#include "proto/omogenexec.pb.h"

namespace omogenexec {

std::string CreateTemporaryRoot();

// A chroot jail, where certain directories from outside the jail can be mounted
// at a given path inside the jail.
class Chroot {

    // The path of the new root
    std::string rootfs;

    void addDefaultRules();

public:
    // Creates a new mount point inside the chroot from the given rule.
    void AddDirectoryRule(const DirectoryRule& rule);

    // Changes the root of the current file system to the one given as the path.
    void SetRoot();

    // Changes the working directory to the root of the new file system.
    void SetWD();

    // Create a new chroot at the specified path, and initalize it with
    // some default mount points.
    Chroot(const std::string& path);

    Chroot(const Chroot&) = delete;
    Chroot& operator=(const Chroot&) = delete;
};

}
