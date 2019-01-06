#include <fcntl.h>
#include <sched.h>
#include <sys/mount.h>
#include <sys/prctl.h>
#include <unistd.h>
#include <atomic>
#include <chrono>
#include <condition_variable>
#include <csignal>
#include <memory>
#include <mutex>
#include <string>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "api/execspec.pb.h"
#include "container/inside/chroot.h"
#include "container/outside/container.h"
#include "glog/logging.h"
#include "grpc/grpc.h"
#include "proto/container.pb.h"
#include "util/files.h"
#include "util/statusor.h"
#include "util/time.h"

using namespace std::chrono_literals;  // for ms

namespace omogenexec {

static const int CHILD_STACK_SIZE = 100 * 1000;  // 100 KB

struct SandboxArgs {
  int stdinFd;
  int stdoutFd;

  int closeFd;
  int closeFd2;

  std::string rootfs;
  api::ContainerSpec containerSpec;
};

static char* const argv[] = {strdup("/usr/bin/omogenrunner"),
                             strdup("--logtostderr"), strdup("--v=3"), NULL};

static int startSandbox(void* argp) {
  SandboxArgs args = *static_cast<SandboxArgs*>(argp);
  PCHECK(prctl(PR_SET_KEEPCAPS, 1) != -1) << "Could not keep capabilities";
  // Close the other ends of the pipes we use for stdin/stdout to avoid keeping
  // them open on our side.
  PCHECK(close(args.closeFd) != -1)
      << "Could not close other end of stdin pipe";
  PCHECK(close(args.closeFd2) != -1)
      << "Could not close other end of stdout pipe";
  PCHECK(dup2(args.stdinFd, 0) != -1) << "Could not open stdin";
  PCHECK(dup2(args.stdoutFd, 6) != -1) << "Could not open stdout";
  PCHECK(close(args.stdinFd) != -1) << "Could not close old stdin";
  PCHECK(close(args.stdoutFd) != -1) << "Could not close old stdout";
  Chroot chroot = Chroot::ForNewRoot(args.rootfs);
  chroot.ApplyContainerSpec(args.containerSpec);
  chroot.SetRoot();
  PCHECK(execv("/usr/bin/omogenrunner", argv) != -1)
      << "Could not start sandbox";
  assert(false && "unreachable");
}

Container::Container(const api::ContainerSpec& spec) {
  PCHECK(pipe(commandPipe) != -1) << "Failed creating command pipe";
  PCHECK(pipe(returnPipe) != -1) << "Failed creating return pipe";
  // The container will get a new root with chroot; we store this in a temporary
  // directory
  containerRoot = MakeTempDir();
  // Clone requires us to provide a new stack for the child process
  std::vector<char> stack(CHILD_STACK_SIZE);
  // Clone and create new namespaces for the contained process
  SandboxArgs args{commandPipe[0], returnPipe[1], commandPipe[1],
                   returnPipe[0],  containerRoot, spec};
  PCHECK((initPid = clone(startSandbox, stack.data() + stack.size(),
                          SIGCHLD | CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS |
                              CLONE_NEWPID | CLONE_NEWUSER | CLONE_NEWUTS,
                          &args)) != -1)
      << "Failed cloning new contained process";
  LOG(INFO) << "Created new container with process ID is " << initPid;
  PCHECK(close(commandPipe[0]) != -1)
      << "Failed closing read end of command pipe";
  PCHECK(close(returnPipe[1]) != -1)
      << "Failed closing write end of return pipe";
  cgroup = std::make_unique<Cgroup>(initPid);
}

Container::~Container() {
  VLOG(3) << "Destroying container";
  killInit();
  VLOG(3) << "Going to wait";
  while (waitInit() == -1)
    ;
  VLOG(3) << "Removing container root";
  RemoveTree(containerRoot);
  VLOG(3) << "Removed container root!";
}

void Container::killInit() {
  // Since we immediately move the contained process out of our process group,
  // it is fine to do kill(-initPid)
  kill(-initPid, SIGKILL);
  kill(initPid, SIGKILL);
}

int Container::waitInit() {
  int status;
  int pid = waitpid(initPid, &status, 0);
  if (pid == -1) {
    if (errno == EINTR) {
      return -1;
    }
    PLOG(FATAL) << "Failed waitpid for init";
  }
  return status;
}

static void setTermination(
    api::Termination* termination,
    const proto::ContainerTermination& containerTermination) {
  switch (containerTermination.termination_case()) {
    case proto::ContainerTermination::kSignal:
      termination->mutable_signal()->set_signal(
          containerTermination.signal().signal());
      return;
    case proto::ContainerTermination::kExit:
      termination->mutable_exit()->set_code(containerTermination.exit().code());
      return;
    case proto::ContainerTermination::kError:
      CHECK(false) << "Errors should be handled outside set termination";
    default:
      assert(false && "Invalid exit status");
  }
}

struct MonitorState {
  Container* container;
  std::atomic<bool> isDead;
  std::atomic<bool> shouldKill;
  bool waitReady;
  std::mutex lock;
  std::condition_variable waitCv;
  proto::ContainerTermination termination;
  int returnPipe;

  explicit MonitorState(Container* cont, int returnPipe)
      : container(cont),
        isDead(false),
        shouldKill(false),
        waitReady(false),
        termination(),
        returnPipe(returnPipe) {}
};

static long long getLimit(const api::ResourceAmounts& limits,
                          api::ResourceType type) {
  for (const auto& limit : limits.amounts()) {
    if (limit.type() == type) {
      return limit.amount();
    }
  }
  LOG(FATAL) << "Could not find resource type " << type;
}

static bool pollStatus(int fd, int* at, char buf[4]) {
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

StatusOr<api::Termination> Container::monitorInit(
    const api::ResourceAmounts& limits) {
  api::Termination response;
  Stopwatch watch;
  MonitorState monitorState(this, returnPipe[0]);
  long long cpuTimeLimit = getLimit(limits, api::ResourceType::CPU_TIME);
  long long wallTimeLimit = getLimit(limits, api::ResourceType::WALL_TIME);
  long long memoryLimit = getLimit(limits, api::ResourceType::MEMORY);

  // We keep one thread that only waits for the process to complete.
  // We also let this thread be responsible for killing the process in case it
  // exceeds its resource limits. This avoids races between killing the process
  // and waiting for it, something that could otherwise result in us killing an
  // unrelated process after the PID has been reused.
  pthread_t waitThread;
  errno = pthread_create(
      &waitThread, nullptr,
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
          std::unique_lock<std::mutex> waitLock(state->lock);
          struct sigaction action;
          memset(&action, 0, sizeof(action));
          action.sa_handler = [](int) {};
          sigaction(SIGALRM, &action, NULL);
          state->waitReady = true;
        }
        state->waitCv.notify_one();
        int lengthRead = 0;
        char lengthBuf[4];
        while (true) {
          if (state->shouldKill) {
            state->container->killInit();
            state->termination.mutable_signal()->set_signal(9);  // SIGKILL
            break;
          } else if (pollStatus(state->returnPipe, &lengthRead, lengthBuf)) {
            CHECK(lengthRead == 4) << "Could not read termination length";
            int length = 0;
            for (int i = 0; i < 4; i++) {
              length = length << 8 | lengthBuf[i];
            }
            CHECK(0 <= length && length <= 5000)
                << "Unreasonable termination length";
            VLOG(3) << "Got length " << length;
            CHECK(state->termination.ParseFromString(
                ReadFromFd(2, state->returnPipe)))
                << "Could not read resulting termination reason";
            break;
          }
          LOG(INFO) << "Reached end of monitor loop " << state->shouldKill;
        }
        state->isDead = true;
        // To avoid some latency, we wake the resource
        // monitor up from its polling sleep whenever the
        // process is dead.
        state->waitCv.notify_one();
        return nullptr;
      },
      &monitorState);
  PCHECK(errno == 0) << "Could not create monitor thread";

  // Wait for the waitThread to set up its signal handler
  {
    std::unique_lock<std::mutex> waitLock(monitorState.lock);
    monitorState.waitCv.wait(waitLock, [&] { return monitorState.waitReady; });
  }

  while (!monitorState.isDead) {
#define CHECK_LIM(current, limit, name) \
  if ((current) > (limit)) {            \
    LOG(INFO) << (name) << " exceeded"; \
    monitorState.shouldKill = true;     \
    pthread_kill(waitThread, SIGALRM);  \
    break;                              \
  }
    // Memory does not need to be monitored, since this is the only limit
    // that the control groups can be limit by itself.
    CHECK_LIM(cgroup->CpuUsed(), cpuTimeLimit, "CPU");
    CHECK_LIM(watch.millis(), wallTimeLimit, "Wall time");

    std::unique_lock<std::mutex> timeoutLock(monitorState.lock);
    monitorState.waitCv.wait_for(timeoutLock, 1ms,
                                 [&] { return !monitorState.isDead; });
#undef CHECK_LIM
  }
  PCHECK((errno = pthread_join(waitThread, nullptr)) == 0)
      << "Could not join with monitor thread";

  long long elapsedMs = watch.millis();
  long long cpuUsedMs = cgroup->CpuUsed();
  long long memoryUsedKb = cgroup->MemoryUsed();
  api::ResourceAmounts* resourceUsage = response.mutable_used_resources();
  api::ResourceAmount* cpuAmount = resourceUsage->add_amounts();
  cpuAmount->set_type(api::ResourceType::CPU_TIME);
  cpuAmount->set_amount(cpuUsedMs);
  api::ResourceAmount* wallTimeAmount = resourceUsage->add_amounts();
  wallTimeAmount->set_type(api::ResourceType::WALL_TIME);
  wallTimeAmount->set_amount(elapsedMs);
  api::ResourceAmount* memoryAmount = resourceUsage->add_amounts();
  memoryAmount->set_type(api::ResourceType::MEMORY);
  memoryAmount->set_amount(memoryUsedKb);

  if (cpuUsedMs > cpuTimeLimit) {
    response.set_resource_exceeded(api::ResourceType::CPU_TIME);
  } else if (memoryUsedKb > memoryLimit) {
    response.set_resource_exceeded(api::ResourceType::MEMORY);
  } else if (elapsedMs > wallTimeLimit) {
    response.set_resource_exceeded(api::ResourceType::WALL_TIME);
  } else if (monitorState.termination.termination_case() ==
             proto::ContainerTermination::kError) {
    return StatusOr<api::Termination>(
        grpc::Status(grpc::StatusCode::INTERNAL,
                     monitorState.termination.error().error_message()));
  } else {
    setTermination(&response, monitorState.termination);
  }
  LOG(INFO) << "Finished with termination " << response.DebugString();
  return StatusOr<api::Termination>(response);
}

StatusOr<api::Termination> Container::Execute(const api::Execution& request) {
  VLOG(3) << "Setting limits for container";
  cgroup->SetMemoryLimit(
      getLimit(request.resource_limits(), api::ResourceType::MEMORY));
  // The sandbox uses one extra process, so increase the limit with this.
  proto::ContainerExecution containerRequest;
  containerRequest.set_allocated_command(new api::Command(request.command()));
  containerRequest.set_allocated_environment(
      new api::Environment(request.environment()));
  containerRequest.set_process_limit(
      getLimit(request.resource_limits(), api::ResourceType::PROCESSES));
  LOG(INFO) << "Sending execution request " << containerRequest.DebugString()
            << " to init";
  std::string requestBytes;
  containerRequest.SerializeToString(&requestBytes);
  WriteIntToFd(requestBytes.size(), commandPipe[1]);
  WriteToFd(commandPipe[1], requestBytes);
  cgroup->Reset();
  VLOG(2) << "Starting monitoring " << request.command().DebugString();
  return monitorInit(request.resource_limits());
}

}  // namespace omogenexec
