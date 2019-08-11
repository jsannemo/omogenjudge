#ifndef SANDBOX_CONTAINER_OUTSIDE_CGROUPS_H
#define SANDBOX_CONTAINER_OUTSIDE_CGROUPS_H

#include <sys/types.h>

#include <string>

#include "sandbox/api/execute_service.pb.h"

using std::string;

namespace omogen {
namespace sandbox {

// The cgroup subsystems supported by the sandbox.
//
// Note: if adding a new subsystem, its name needs to be added to the list of
// names subsystemName in cgroups.cc
enum class CgroupSubsystem { CPU_ACCT = 0, MEMORY, INVALID };

// A cgroup enables tracking of a process' resource usage.
// It also allows limiting the amount of memory the process can use.
class Cgroup {
  string name;
  pid_t pid;
  long long _mem_limit_kb;

  string GetSubsystemPath(CgroupSubsystem subsystem);
  string GetSubsystemOp(CgroupSubsystem subsystem, const string& op);
  void EnableSubsystem(CgroupSubsystem subsystem);
  void DisableSubsystem(CgroupSubsystem subsystem);

 public:
  // Creates a new cgroup for a given process and enables the subsystems used
  // for the process.
  Cgroup(pid_t pid);
  ~Cgroup();

  // The total CPU usage of the process and its children, in milliseconds.
  long long CpuUsed();
  // The total memory usage of the process and its children, in kilobytes.
  long long MemoryUsed();
  // Sets the maximum amount of memory that the process can use.
  // Note that cgroup swap accounting must be enabled for this to work
  // properly (or swap disabled).
  void SetMemoryLimit(long long memLimitKb);
  // Reset the resource usage statistics.
  void Reset();

  static Cgroup MakeCgroupFor(pid_t pid);
};

}  // namespace sandbox
}  // namespace omogen
#endif
