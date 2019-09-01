#ifndef SANDBOX_CONTAINER_INSIDE_CHROOT_H
#define SANDBOX_CONTAINER_INSIDE_CHROOT_H

#include <string>

#include "sandbox/api/containerspec.pb.h"

using std::string;

namespace omogen {
namespace sandbox {

// A chroot jail, where certain directories from outside the jail can be mounted
// at a given path inside the jail.
class Chroot {
  // The path of the new root.
  string rootfs;

  // Adds a set of default rules to ensure that the new environment is a
  // somewhat functional system.
  void AddDefaultRules();

  // Creates a new mount point inside the chroot from the given rule.
  void AddDirectoryMount(const DirectoryMount& rule);

  // Create a new chroot at the specified path, and initalize it with
  // some default mount points.
  explicit Chroot(const string& new_root);

 public:
  // Applies a container specification to this environment, setting up all
  // the mount rules specified.
  void ApplyContainerSpec(const ContainerSpec& spec);

  // Changes the root of the current file system to the one given as the path.
  // The new working directory will be the root of the new system.
  void SetRoot();

  // Creates a new chroot jail with a given new root directory.
  static Chroot ForNewRoot(const string& new_root);

  Chroot(const Chroot&) = delete;
  Chroot& operator=(const Chroot&) = delete;
};

}  // namespace sandbox
}  // namespace omogen
#endif
