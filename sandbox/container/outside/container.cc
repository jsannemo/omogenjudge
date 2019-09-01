#include "sandbox/container/outside/container.h"

#include <fcntl.h>
#include <sched.h>
#include <sys/mount.h>
#include <sys/prctl.h>
#include <unistd.h>

#include <atomic>
#include <chrono>
#include <condition_variable>
#include <csignal>
#include <cstdlib>
#include <memory>
#include <mutex>
#include <string>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "glog/logging.h"
#include "glog/raw_logging.h"
#include "grpc/grpc.h"
#include "sandbox/api/execspec.pb.h"
#include "sandbox/proto/container.pb.h"
#include "util/cpp/files.h"
#include "util/cpp/statusor.h"
#include "util/cpp/time.h"

using omogen::util::ReadFromFd;
using omogen::util::ReadIntFromFd;
using omogen::util::RemoveTree;
using omogen::util::Stopwatch;
using omogen::util::WriteIntToFd;
using omogen::util::WriteToFd;
using std::string;

using namespace std::chrono_literals;  // for ms

namespace omogen {
namespace sandbox {

static const char* kContainerRoot = "/var/lib/omogen/sandbox";

static const int kInodeLimit = 1000;

void Container::startSandbox(const ContainerSpec& spec, int sandboxId) {
  init_pid = fork();
  PCHECK(init_pid != -1) << "Failed forking";
  if (init_pid != 0) {
    cgroup = std::make_unique<Cgroup>(init_pid);
    return;
  }
  close(command_pipe[1]);
  close(return_pipe[0]);

  unsigned long long maxDiskBlocks =
      (unsigned long long)spec.max_disk_kb() * 1000 / 1024;
  char* const argv[] = {strdup("/usr/bin/omogenjudge-sandboxr"),
                        strdup("-logtostderr"),
                        strdup("-v=4"),
                        strdup(absl::StrCat(sandboxId).c_str()),
                        strdup(absl::StrCat(command_pipe[0]).c_str()),
                        strdup(absl::StrCat(return_pipe[1]).c_str()),
                        strdup(absl::StrCat(maxDiskBlocks).c_str()),
                        strdup(absl::StrCat(kInodeLimit).c_str()),
                        NULL};
  PCHECK(execv("/usr/bin/omogenjudge-sandboxr", argv) != -1)
      << "Could not start sandbox";
}

Container::Container(std::unique_ptr<ContainerId> container_id_,
                     const ContainerSpec& spec)
    : container_id(std::move(container_id_)) {
  PCHECK(pipe(command_pipe) != -1) << "Failed creating command pipe";
  PCHECK(pipe(return_pipe) != -1) << "Failed creating return pipe";
  // The container will get a new root with chroot; we store this in a temporary
  // directory.
  container_root = absl::StrCat(kContainerRoot, "/", container_id->Get());
  startSandbox(spec, container_id->Get());
  close(command_pipe[0]);
  close(return_pipe[1]);

  string spec_bytes;
  spec.SerializeToString(&spec_bytes);
  WriteIntToFd(spec_bytes.size(), command_pipe[1]);
  WriteToFd(command_pipe[1], spec_bytes);
}

Container::~Container() {
  VLOG(3) << "Destroying container";
  KillInit();
  VLOG(3) << "Going to wait";
  while (WaitInit() == -1)
    ;
  init_pid = 0;
  VLOG(3) << "Removing container root";
  std::string cmd =
      absl::StrCat("/usr/bin/omogenjudge-sandboxc ", container_id->Get());
  CHECK(system(cmd.c_str()) == 0) << "Could not clear submission";
  VLOG(3) << "Removed container root!";
}

void Container::KillInit() {
  // Since we immediately move the contained process out of our process group,
  // it is fine to do kill(-init_pid)
  kill(-init_pid, SIGKILL);
  kill(init_pid, SIGKILL);
}

int Container::WaitInit() {
  int status;
  int pid = waitpid(init_pid, &status, 0);
  if (pid == -1) {
    if (errno == EINTR) {
      return -1;
    }
    PLOG(FATAL) << "Failed waitpid for init";
  }
  return status;
}

bool Container::IsDead() { return init_pid == 0; }

static void SetTermination(
    Termination* termination,
    const proto::ContainerTermination& container_termination) {
  switch (container_termination.termination_case()) {
    case proto::ContainerTermination::kSignal:
      termination->mutable_signal()->set_signal(
          container_termination.signal().signal());
      return;
    case proto::ContainerTermination::kExit:
      termination->mutable_exit()->set_code(
          container_termination.exit().code());
      return;
    case proto::ContainerTermination::kError:
      CHECK(false) << "Errors should be handled outside set termination";
    default:
      assert(false && "Invalid exit status");
  }
}

struct MonitorState {
  Container* container;
  std::atomic<bool> is_dead;
  std::atomic<bool> should_kill;
  bool wait_ready;
  std::mutex lock;
  std::condition_variable wait_cv;
  proto::ContainerTermination termination;
  int return_pipe;

  explicit MonitorState(Container* cont, int return_pipe)
      : container(cont),
        is_dead(false),
        should_kill(false),
        wait_ready(false),
        termination(),
        return_pipe(return_pipe) {}
};

static long long GetLimit(const ResourceAmounts& limits, ResourceType type) {
  for (const auto& limit : limits.amounts()) {
    if (limit.type() == type) {
      return limit.amount();
    }
  }
  LOG(FATAL) << "Could not find resource type " << type;
}

static bool PollStatus(int fd, int* at, char buf[4]) {
  VLOG(3) << "Polling exit status";
  while (true) {
    VLOG(3) << "Read " << *at << " so far";
    int r = read(fd, buf, 4 - *at);
    VLOG(3) << "Read returned with " << r;
    PCHECK(r != -1 || errno == EINTR) << "Failed reading return value";
    if (errno == EINTR) {
      return false;
    }
    if (r == 0) {
      return true;
    }
    *at += r;
    if (*at == 4) {
      break;
    }
  }
  return true;
}

StatusOr<Termination> Container::MonitorInit(const ResourceAmounts& limits) {
  Termination response;
  Stopwatch watch;
  MonitorState monitor_state(this, return_pipe[0]);
  long long cpu_time_limit = GetLimit(limits, ResourceType::CPU_TIME);
  long long wall_time_limit = GetLimit(limits, ResourceType::WALL_TIME);
  long long memory_limit = GetLimit(limits, ResourceType::MEMORY);

  // We keep one thread that only waits for the process to complete.
  // We also let this thread be responsible for killing the process in case it
  // exceeds its resource limits. This avoids races between killing the process
  // and waiting for it, something that could otherwise result in us killing an
  // unrelated process after the PID has been reused.
  pthread_t wait_thread;
  errno = pthread_create(
      &wait_thread, nullptr,
      [](void* arg) -> void* {
        MonitorState* state = static_cast<MonitorState*>(arg);
        // The resource monitor loop notifies us if we should
        // kill init by giving us SIGALRM to interrupt our
        // wait. We use a lock and flag to tell the monitor
        // when we have set up our own signal handler to
        // avoid getting such a signal before the handler is
        // installed, otherwise we would get killed by the
        // signal.
        {
          std::unique_lock<std::mutex> wait_lock(state->lock);
          struct sigaction action;
          memset(&action, 0, sizeof(action));
          action.sa_handler = [](int) {};
          sigaction(SIGALRM, &action, NULL);
          state->wait_ready = true;
        }
        state->wait_cv.notify_one();
        int length_read = 0;
        char length_buf[4];
        while (true) {
          if (state->should_kill) {
            state->container->KillInit();
            state->termination.mutable_signal()->set_signal(9);  // SIGKILL
            break;
          } else if (PollStatus(state->return_pipe, &length_read, length_buf)) {
            CHECK(length_read == 4) << "Could not read termination length";
            int length = 0;
            for (int i = 0; i < 4; i++) {
              length = length << 8 | length_buf[i];
            }
            CHECK(0 <= length && length <= 5000)
                << "Unreasonable termination length";
            VLOG(3) << "Got length " << length;
            CHECK(state->termination.ParseFromString(
                ReadFromFd(length, state->return_pipe)))
                << "Could not read resulting termination reason";
            break;
          }
          LOG(INFO) << "Reached end of monitor loop - should kill"
                    << state->should_kill;
        }
        state->is_dead = true;
        // To avoid some latency, we wake the resource
        // monitor up from its polling sleep whenever the
        // process is dead.
        state->wait_cv.notify_one();
        return nullptr;
      },
      &monitor_state);
  PCHECK(errno == 0) << "Could not create monitor thread";

  // Wait for the wait_thread to set up its signal handler
  {
    std::unique_lock<std::mutex> wait_lock(monitor_state.lock);
    monitor_state.wait_cv.wait(wait_lock,
                               [&] { return monitor_state.wait_ready; });
  }

  while (!monitor_state.is_dead) {
#define CHECK_LIM(current, limit, name) \
  if ((current) > (limit)) {            \
    LOG(INFO) << (name) << " exceeded"; \
    monitor_state.should_kill = true;   \
    pthread_kill(wait_thread, SIGALRM); \
    break;                              \
  }
    // Memory does not need to be monitored, since this is the only limit
    // that the control groups can be limit by itself.
    CHECK_LIM(cgroup->CpuUsed(), cpu_time_limit, "CPU");
    CHECK_LIM(watch.millis(), wall_time_limit, "Wall time");

    std::unique_lock<std::mutex> timeoutLock(monitor_state.lock);
    monitor_state.wait_cv.wait_for(timeoutLock, 10ms,
                                   [&] { return !monitor_state.is_dead; });
#undef CHECK_LIM
  }
  PCHECK((errno = pthread_join(wait_thread, nullptr)) == 0)
      << "Could not join with monitor thread";

  long long elapsed_ms = watch.millis();
  long long cpu_used_ms = cgroup->CpuUsed();
  long long memory_used_kb = cgroup->MemoryUsed();
  ResourceAmounts* resource_usage = response.mutable_used_resources();
  ResourceAmount* cpu_amount = resource_usage->add_amounts();
  cpu_amount->set_type(ResourceType::CPU_TIME);
  cpu_amount->set_amount(cpu_used_ms);
  ResourceAmount* wall_time_amount = resource_usage->add_amounts();
  wall_time_amount->set_type(ResourceType::WALL_TIME);
  wall_time_amount->set_amount(elapsed_ms);
  ResourceAmount* memory_amonut = resource_usage->add_amounts();
  memory_amonut->set_type(ResourceType::MEMORY);
  memory_amonut->set_amount(memory_used_kb);

  if (cpu_used_ms > cpu_time_limit) {
    response.set_resource_exceeded(ResourceType::CPU_TIME);
  } else if (memory_used_kb > memory_limit) {
    // TODO: Check if cgroup OOM-killed the program to see if it should have
    // gotten memory limit exceeded
    response.set_resource_exceeded(ResourceType::MEMORY);
  } else if (elapsed_ms > wall_time_limit) {
    response.set_resource_exceeded(ResourceType::WALL_TIME);
  } else if (monitor_state.termination.termination_case() ==
             proto::ContainerTermination::kError) {
    return StatusOr<Termination>(
        grpc::Status(grpc::StatusCode::INTERNAL,
                     monitor_state.termination.error().error_message()));
  } else {
    SetTermination(&response, monitor_state.termination);
  }
  LOG(INFO) << "Finished with termination " << response.DebugString();
  return StatusOr<Termination>(response);
}

StatusOr<Termination> Container::Execute(const Execution& request) {
  if (IsDead()) {
    return grpc::Status(grpc::StatusCode::INTERNAL,
                        "Tried to execute on a dead container");
  }

  long long memory_limit =
      GetLimit(request.resource_limits(), ResourceType::MEMORY);
  cgroup->SetMemoryLimit(memory_limit);
  proto::ContainerExecution container_request;
  *container_request.mutable_command() = request.command();
  *container_request.mutable_environment() = request.environment();
  container_request.set_process_limit(
      GetLimit(request.resource_limits(), ResourceType::PROCESSES));
  LOG(INFO) << "Sending execution request " << container_request.DebugString()
            << " to init";
  string request_bytes;
  container_request.SerializeToString(&request_bytes);
  WriteIntToFd(request_bytes.size(), command_pipe[1]);
  Reset();
  WriteToFd(command_pipe[1], request_bytes);
  VLOG(2) << "Starting monitoring " << request.command().DebugString();
  return MonitorInit(request.resource_limits());
}

void Container::Reset() { cgroup->Reset(); }

}  // namespace sandbox
}  // namespace omogen
