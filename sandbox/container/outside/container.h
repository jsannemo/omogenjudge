#ifndef SANDBOX_CONTAINER_OUTSIDE_CONTAINER_H
#define SANDBOX_CONTAINER_OUTSIDE_CONTAINER_H

#include <sys/types.h>
#include <sys/wait.h>

#include "sandbox/api/execspec.pb.h"
#include "sandbox/container/outside/cgroups.h"
#include "sandbox/container/outside/container_id.h"
#include "util/cpp/statusor.h"

using omogen::util::StatusOr;
using std::unique_ptr;

namespace omogen {
namespace sandbox {

struct SandboxArgs {
  int id;
  string s_id;

  int stdin_fd;
  string s_stdin_fd;
  int stdout_fd;
  string s_stdout_fd;

  int close_fd;
  int close_fd2;

  string rootfs;
  ContainerSpec container_spec;
};

// A process container, implemented using Linux control groups, namespaces and
// rlimits. Note that the constructor creates the new process ahead-of-time and
// performs some setup. The remaining setup is performed once Execute() is
// called, minimizing the latency a bit.
class Container {
  // Process ID of the child process we are executing the new program in (called
  // init)
  pid_t init_pid;
  // Since we may receive the execution request after starting the new process,
  // we use a pipe to send the request to the process.
  int command_pipe[2];
  // A pipe used by the container to tell us what the return status of the user
  // program was.
  int return_pipe[2];
  // The path to the new root with specific paths mounted to it
  std::string container_root;

  unique_ptr<Cgroup> cgroup;

  unique_ptr<ContainerId> container_id;

  std::vector<char> stack;

  SandboxArgs args;

  StatusOr<Termination> MonitorInit(const ResourceAmounts& limits);

  // Forcibly kill()'s the init process.
  void KillInit();

  // wait()'s the init process after it has been killed.
  int WaitInit();

  void startSandbox(const ContainerSpec& spec, int sandboxId);

 public:
  // Performs an execution in the sandbox. Assumes the given execution
  // has been validated.
  StatusOr<Termination> Execute(const Execution& request);

  // Whether the container process has died. A dead container
  // may not be used for anything.
  bool IsDead();

  void Reset();

  Container(unique_ptr<ContainerId> id, const ContainerSpec& spec);
  ~Container();

  Container(const Container&) = delete;
  Container& operator=(const Container&) = delete;
};

}  // namespace sandbox
}  // namespace omogen
#endif
