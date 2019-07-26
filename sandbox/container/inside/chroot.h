#pragma once

#include <string>

#include "sandbox/api/containerspec.pb.h"

using std::string;

namespace omogen {
namespace sandbox {

// A chroot jail, where certain directories from outside the jail can be mounted
// at a given path inside the jail.
class Chroot {
  // The path of the new root
  string rootfs;

  // Adds a set of default rules to ensure that the new environment is a
  // somewhat functional system.
  void addDefaultRules();

  // Creates a new mount point inside the chroot from the given rule.
  void addDirectoryMount(const DirectoryMount& rule);

  // Create a new chroot at the specified path, and initalize it with
  // some default mount points.
  explicit Chroot(const string& newRoot);

 public:
  // Applies a container specification to this environment, setting up all
  // the mount rules specified.
  void ApplyContainerSpec(const ContainerSpec& spec);

  // Changes the root of the current file system to the one given as the path.
  // The new working directory will be the root of the new system.
  void SetRoot();

  // Creates a new chroot jail with a given new root directory.
  static Chroot ForNewRoot(const string& newRoot);

  Chroot(const Chroot&) = delete;
  Chroot& operator=(const Chroot&) = delete;
};

}  // namespace sandbox
}  // namespace omogen
