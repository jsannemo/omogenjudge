#include <fcntl.h>
#include <iostream>
#include <sys/mount.h>
#include <sys/stat.h>
#include <unistd.h>

#include "chroot.h"
#include "util/error.h"
#include "util/files.h"
#include "util/log.h"

using std::endl;
using std::string;

namespace omogenexec {

void Chroot::AddDirectoryRule(const DirectoryRule& rule) {
    if (rule.newpath().empty() || rule.newpath()[0] != '/') {
        OE_LOG(FATAL) << "Directory rule target is not an absolute path" << endl;
        OE_CRASH();
    }
    string target = rootfs + rule.newpath();
    MakeDirParents(target);
    int mountFlags = MS_BIND | MS_NOSUID | MS_NODEV;
    if (!rule.writable()) {
        mountFlags |= MS_RDONLY;
    }
    if (mount(rule.oldpath().c_str(), target.c_str(), nullptr, mountFlags, nullptr) == -1) {
        OE_FATAL("mount");
    }
    // When BIND:ing using mount, read-only (and possible other flags) may 
    // require a remount to take effect (see e.g. https://lwn.net/Articles/281157/)
    if (mount(rule.oldpath().c_str(), target.c_str(), nullptr, MS_REMOUNT | mountFlags, nullptr) == -1) {
        OE_FATAL("mount");
    }
}

void Chroot::addDefaultRules() {
    // Since we are in a new PID namespace, the old procfs will have incorrect data,
    // so we mount a new one.
    string procPath = rootfs + "/proc";
    MakeDirParents(procPath);
    if (mount("/proc", procPath.c_str(), "proc", MS_NODEV | MS_NOEXEC | MS_NOSUID, nullptr) == -1) {
        OE_FATAL("mount");
    }

    DirectoryRule rule;
    rule.set_oldpath("/bin");
    rule.set_newpath("/bin");
    AddDirectoryRule(rule);

    rule.set_oldpath("/usr");
    rule.set_newpath("/usr");
    AddDirectoryRule(rule);

    rule.set_oldpath("/lib");
    rule.set_newpath("/lib");
    AddDirectoryRule(rule);

    if (DirectoryExists("/lib64")) {
        rule.set_oldpath("/lib64");
        rule.set_newpath("/lib64");
        AddDirectoryRule(rule);
    }

    if (DirectoryExists("/lib32")) {
        rule.set_oldpath("/lib32");
        rule.set_newpath("/lib32");
        AddDirectoryRule(rule);
    }

}

void Chroot::SetRoot() {
    if (chroot(rootfs.c_str()) == -1) {
        OE_FATAL("chroot");
    }
    if (chdir("/") == -1) {
        OE_FATAL("chdir");
    }
}

void Chroot::SetWD() {
    if (chdir(rootfs.c_str()) == -1) {
        OE_FATAL("chdir");
    }
}

Chroot::Chroot(const string& path) : rootfs(path) {
    MakeDir(path);
    if (mount(nullptr, "/", nullptr, MS_REC | MS_PRIVATE, nullptr) == -1) {
        OE_FATAL("mount");
    }
    addDefaultRules();
}

}
