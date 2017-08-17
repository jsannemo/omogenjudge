#include <fcntl.h>
#include <ftw.h>
#include <sys/mount.h>
#include <sys/stat.h>
#include <unistd.h>

#include "chroot.h"
#include "errors/errors.h"
#include "logger/logger.h"

void makeDir(const string& path) {
    errno = 0;
    if (mkdir(path.c_str(), 0755) == -1 && errno != EEXIST) {
        crashSyscall("mkdir");
    }
}

void makeDirectoryPath(const string& path) {
    if (path.empty()) {
        return;
    }
    // Create all directories on a path by repeatedly finding the next directory separator
    // and creating the path up to this location
    size_t at = 0;
    while (at < path.size()) {
        size_t next = path.find('/', at + 1);
        if (next == string::npos) {
            next = path.size();
        }
        makeDir(path.substr(0, next));
        at = next;
    }
    struct stat st;
    if (stat(path.c_str(), &st) == -1) {
        crashSyscall("stat");
    }
    // Double-check that the path we should have created is indeed a directory now,
    // and not e.g. a file
    if (!S_ISDIR(st.st_mode)) {
        LOG(FATAL) << "Failed to create path" << endl;
        crash();
    }
}

void Chroot::AddDirectoryRule(const DirectoryRule& rule) {
    string target = rootfs + rule.destination();
    makeDirectoryPath(target);
    int mountFlags = MS_BIND | MS_NOSUID;
    if (!rule.writable()) {
        mountFlags |= MS_RDONLY;
    }
    if (!rule.executable()) {
        mountFlags |= MS_NOEXEC;
    }
    LOG(TRACE) << "Mounting " << rule.source() << " on " << target << endl;
    if (mount(rule.source().c_str(), target.c_str(), nullptr, mountFlags, nullptr) == -1) {
        crashSyscall("mount");
    }
    // Remount to make all the flags take effect
    if (mount(rule.source().c_str(), target.c_str(), nullptr, MS_REMOUNT | mountFlags, nullptr) == -1) {
        crashSyscall("mount");
    }
}

void Chroot::addDefaultRules() {
    string procPath = rootfs + "/proc";
    makeDirectoryPath(procPath);
    if (mount("/proc", "proc", "proc", MS_NODEV | MS_NOEXEC | MS_NOSUID, nullptr) == -1) {
        crashSyscall("mount");
    }

    DirectoryRule rule;
    rule.set_source("/bin");
    rule.set_destination("/bin");
    rule.set_executable(true);
    AddDirectoryRule(rule);

    rule.set_source("/usr");
    rule.set_destination("/usr");
    AddDirectoryRule(rule);

    rule.set_source("/lib");
    rule.set_destination("/lib");
    AddDirectoryRule(rule);

    // TODO(jsannemo): verify that this directory actually exists before mounting it
    // since this is not necessarily the case on e.g. 32 bit machines
    rule.set_source("/lib64");
    rule.set_destination("/lib64");
    AddDirectoryRule(rule);
}

void Chroot::SetRoot() {
    if (chroot(rootfs.c_str()) == -1) {
        crashSyscall("chroot");
    }
}

string CreateTemporaryRoot() {
    char *rootPath = strdup("/tmp/omogencontainXXXXXX");
    if (rootPath == nullptr) {
        crashSyscall("strdup");
    }
    if (mkdtemp(rootPath) == nullptr) {
        crashSyscall("mkdtemp");
    }
    string ret(rootPath);
    free(rootPath);
    return ret;
}

Chroot::Chroot(const string& path) : rootfs(path) {
    if (chdir(rootfs.c_str()) == -1) {
        crashSyscall("chdir");
    }
    LOG(TRACE) << "Setting up container root in " << rootfs << endl;
    if (mount(nullptr, "/", nullptr, MS_REC | MS_PRIVATE, nullptr) == -1) {
        crashSyscall("mount");
    }
    addDefaultRules();
}

int removeTree0(const char *filePath, const struct stat *statData, int typeflag, struct FTW *ftwbuf) {
    if (S_ISDIR(statData->st_mode)) {
        if (rmdir(filePath) == -1) {
            crashSyscall("rmdir");
        }
    } else {
        if (unlink(filePath) == -1) {
            crashSyscall("unlink");
        }
    }
    return FTW_CONTINUE;
}

void DestroyDirectory(const string& directoryPath) {
    if (nftw(directoryPath.c_str(), removeTree0, 32, FTW_MOUNT | FTW_PHYS | FTW_DEPTH) == -1) {
        crashSyscall("nftw");
    }
}
