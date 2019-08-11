#include "sandbox/container/inside/chroot.h"

#include <fcntl.h>
#include <sys/mount.h>
#include <sys/stat.h>
#include <unistd.h>

#include <iostream>

#include "glog/logging.h"
#include "glog/raw_logging.h"
#include "util/cpp/files.h"

using omogen::util::DirectoryExists;
using omogen::util::MakeDirParents;
using std::string;

namespace omogen {
namespace sandbox {

void Chroot::AddDirectoryMount(const DirectoryMount& rule) {
  if (rule.path_inside_container().empty() ||
      rule.path_inside_container()[0] != '/') {
    RAW_LOG(FATAL, "Directory rule target is not an absolute path");
  }
  string target = rootfs + rule.path_inside_container();
  MakeDirParents(target);
  // MS_NODEV: ensure we can't access special devices somehow in our sandbox.
  // MS_NOSUID: don't allow the sandboxed program to do anything as root by
  // calling some setuid root programs.
  int mountFlags = MS_BIND | MS_NOSUID | MS_NODEV;
  if (!rule.writable()) {
    mountFlags |= MS_RDONLY;
  }
  RAW_CHECK(mount(rule.path_outside_container().c_str(), target.c_str(),
                  nullptr, mountFlags, nullptr) != -1,
            "Could not mount rule");
  // When BIND:ing using mount, read-only (and possible other flags) may
  // require a remount to take effect (see e.g.
  // https://lwn.net/Articles/281157/).
  RAW_CHECK(mount(rule.path_outside_container().c_str(), target.c_str(),
                  nullptr, MS_REMOUNT | mountFlags, nullptr) != -1,
            "Could not remount rule");
}

void Chroot::AddDefaultRules() {
  // Since we are in a new PID namespace, the old procfs will have incorrect
  // data, so we mount a new one.
  string proc_path = rootfs + "/proc";
  MakeDirParents(proc_path);
  RAW_CHECK(mount("/proc", proc_path.c_str(), "proc",
                  MS_NODEV | MS_NOEXEC | MS_NOSUID, nullptr) != -1,
            "Could not mount new proc namespace");

  DirectoryMount rule;
  rule.set_path_outside_container("/bin");
  rule.set_path_inside_container("/bin");
  AddDirectoryMount(rule);

  rule.set_path_outside_container("/usr/bin");
  rule.set_path_inside_container("/usr/bin");
  AddDirectoryMount(rule);

  rule.set_path_outside_container("/usr/lib");
  rule.set_path_inside_container("/usr/lib");
  AddDirectoryMount(rule);

  rule.set_path_outside_container("/lib");
  rule.set_path_inside_container("/lib");
  AddDirectoryMount(rule);

  if (DirectoryExists("/usr/lib32")) {
    rule.set_path_outside_container("/usr/lib32");
    rule.set_path_inside_container("/usr/lib32");
    AddDirectoryMount(rule);
  }

  if (DirectoryExists("/lib64")) {
    rule.set_path_outside_container("/lib64");
    rule.set_path_inside_container("/lib64");
    AddDirectoryMount(rule);
  }

  if (DirectoryExists("/lib32")) {
    rule.set_path_outside_container("/lib32");
    rule.set_path_inside_container("/lib32");
    AddDirectoryMount(rule);
  }
}

void Chroot::SetRoot() {
  RAW_CHECK(chroot(rootfs.c_str()) != -1, "Could not chroot to the new root");
  RAW_CHECK(chdir("/") != -1, "Could not chdir to the new root");
}

void Chroot::ApplyContainerSpec(const ContainerSpec& spec) {
  for (const auto& mount : spec.mounts()) {
    AddDirectoryMount(mount);
  }
}

Chroot::Chroot(const string& path) : rootfs(path) {
  RAW_VLOG(2, "Chroot path %s", path.c_str());
  RAW_CHECK(DirectoryExists(path), "Path does not exist");
  RAW_CHECK(mount(nullptr, "/", nullptr, MS_REC | MS_PRIVATE, nullptr) != -1,
            "Could not update mounts to be private");
  AddDefaultRules();
}

Chroot Chroot::ForNewRoot(const string& new_root) { return Chroot(new_root); }

}  // namespace sandbox
}  // namespace omogen
