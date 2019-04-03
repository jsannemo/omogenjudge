#pragma once

#include <sys/types.h>
#include <sys/wait.h>

#include "api/execspec.pb.h"
#include "cgroups.h"
#include "container/outside/container_id.h"
#include "util/statusor.h"

using std::unique_ptr;

namespace omogenexec {

// A process container, implemented using Linux control groups, namespaces and
// rlimits. Note that the constructor creates the new process ahead-of-time and
// performs some setup. The remaining setup is performed once Execute() is
// called, minimizing the latency a bit.
class Container {
  // Process ID of the child process we are executing the new program in (called
  // init)
  pid_t initPid;
  // Since we may receive the execution request after starting the new process,
  // we use a pipe to send the request to the process.
  int commandPipe[2];
  // A pipe used by the container to tell us what the return status of the user
  // program was.
  int returnPipe[2];
  // The path to the new root with specific paths mounted to it
  std::string containerRoot;

  unique_ptr<Cgroup> cgroup;

  unique_ptr<ContainerId> containerId;

  StatusOr<api::Termination> monitorInit(const api::ResourceAmounts& limits);
  void killInit();
  int waitInit();

 public:
  StatusOr<api::Termination> Execute(const api::Execution& request);

  Container(unique_ptr<ContainerId> id, const api::ContainerSpec& spec);
  ~Container();
  
  bool IsDead();

  Container(const Container&) = delete;
  Container& operator=(const Container&) = delete;
};

}  // namespace omogenexec
