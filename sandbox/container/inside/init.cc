// Init is the init process (PID 1) of a container. Its purpose is to recieve
// execution requests from the outside container and fork of a process to
// execute it.
#include <sys/fcntl.h>
#include <sys/prctl.h>
#include <sys/resource.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include <cstdio>
#include <iostream>
#include <string>

#include "absl/strings/numbers.h"
#include "glog/logging.h"
#include "glog/raw_logging.h"
#include "sandbox/container/inside/chroot.h"
#include "sandbox/container/inside/setup.h"
#include "sandbox/proto/container.pb.h"
#include "util/cpp/files.h"

using omogen::sandbox::proto::ContainerExecution;
using omogen::sandbox::proto::ContainerTermination;
using omogen::util::ReadFromFd;
using omogen::util::ReadIntFromFd;
using omogen::util::WriteIntToFd;
using omogen::util::WriteToFd;
using std::string;

namespace omogen {
namespace sandbox {

// Returns a termination cause based on the given status, which should
// be given in the waitpid format.
static ContainerTermination TerminationForStatus(int wait_status) {
  ContainerTermination termination;
  if (WIFEXITED(wait_status)) {
    termination.mutable_exit()->set_code(WEXITSTATUS(wait_status));
  } else if (WIFSIGNALED(wait_status)) {
    termination.mutable_signal()->set_signal(WTERMSIG(wait_status));
  } else {
    assert(false && "Invalid exit status");
  }
  return termination;
}

// Returns a termination for a sandbox-level error that occured.
static ContainerTermination TerminationForError(const string& error_message) {
  ContainerTermination termination;
  termination.mutable_error()->set_error_message(error_message);
  return termination;
}

static ContainerTermination Execute(const ContainerExecution& request) {
  // We need to set this here rather than in setup since we lose privilege to
  // change this to a potentially higher number after the fork.
  rlim_t process_limit = request.process_limit() + 2;  // +1 for this process
  rlimit rlim = {.rlim_cur = process_limit, .rlim_max = process_limit};
  PCHECK(setrlimit(RLIMIT_NPROC, &rlim) != -1)
      << "Could not set the process limit";

  // Start a fork to set up the execution environment for the request.
  // In the parent, we will wait for the request to finish. Since an error
  // may occur during setup of the contained process, we keep a pipe open
  // that the child can send error messages to us with. We use O_CLOEXEC
  // so that we don't need to close the stream ourselves to keep it available
  // for as long as possible.
  int error_pipe[2];
  PCHECK(pipe2(error_pipe, O_CLOEXEC) != -1) << "Could not create error pipe";

  // Note that the child process will not survive the parent death since the
  // parent is the init process of a new PID namespace. That means the child is
  // automatically killed if the parent is killed.
  pid_t which = fork();
  if (which == 0) {
    VLOG(2) << "Container child started";
    close(error_pipe[0]);                 // We only write to the error pipe
    SetupAndRun(request, error_pipe[1]);  // Will either execve or crash
  } else {
    PCHECK(which != -1) << "Could not fork execution process";
    close(error_pipe[1]);  // We only read from the error pipe
    // Try reading an error message until the pipe closes. The child writes
    // a \1 byte just before it decides to close the pipe. This means that we
    // can also be aware of errors that crashes the child (without first being
    // handled by sending us an error message) by checking if we get a \1 byte.
    // This will never we part of a normal error message, since they only
    // contain printable ASCII.
    string error_message;
    char buf[1024];
    while (true) {
      int ret = read(error_pipe[0], buf, sizeof(buf));
      PCHECK(ret != -1 || errno == EINTR) << "Could not read error message";
      // We were interrupted.
      if (ret == -1) {
        continue;
      }
      // Pipe was closed.
      if (ret == 0) {
        break;
      }
      error_message += string(buf, buf + ret);
    }
    // The process crashed before it closed its error stream.
    if (error_message.empty()) {
      return TerminationForError("Unhandled error during setup");
    }
    if (error_message[0] != '\1') {
      return TerminationForError("Error during execution setup: " +
                                 error_message);
    }
    while (true) {
      int status;
      int ret = waitpid(which, &status, 0);
      PCHECK(ret != -1 || errno == EINTR) << "Could not wait";
      // We were interrupted
      if (ret == -1) {
        continue;
      }
      // If our child process decided to start children of their own, they will
      // now become a child of us since we are init. Therefore, we make sure to
      // SIGKILL all of them before we return, in case we want to reuse our
      // sandbox.
      // TODO: is kill(-pid) a race against a fork bomb?
      PCHECK(kill(-1, SIGKILL) != -1 || errno == ESRCH)
          << "Did not manage to kill all remaining processes in the container";
      while (true) {
        int ret = waitpid(-1, nullptr, WNOHANG);
        PCHECK(ret != -1 || errno == ECHILD)
            << "Could not wait for remaining daemons";
        if (ret == 0 || errno == ECHILD) {
          break;
        }
      }
      return TerminationForStatus(status);
    }
  }
}

}  // namespace sandbox
}  // namespace omogen

int main(int argc, char** argv) {
  gflags::ParseCommandLineFlags(&argc, &argv, true);
  google::InitGoogleLogging(argv[0]);
  google::InstallFailureSignalHandler();

  CHECK(argc == 4) << "Incorrect number of arguments";
  int sandbox_id;
  CHECK(absl::SimpleAtoi(string(argv[1]), &sandbox_id))
      << "Can not convert sandbox ID to integer";
  CHECK(sandbox_id >= 0) << "Sandbox ID was negative";
  int in_id, out_id;
  CHECK(absl::SimpleAtoi(string(argv[2]), &in_id))
      << "Can not convert input FD to integer";
  CHECK(absl::SimpleAtoi(string(argv[3]), &out_id))
      << "Can not convert output FD to integer";
  CHECK(in_id >= 0) << "Input FD was negative";
  CHECK(out_id >= 0) << "Ouptut FD was negative";

  // Kill us if the main sandbox is killed, to prevent our child from possibly
  // keep running. This is not a race with the parent death, since the read
  // later will crash us in case our parent dies after the prctl call.
  // Furthermore, as a result of our death we will take with us any processes
  // running in the sandbox since we are PID 1 in a PID namespace.
  CHECK(prctl(PR_SET_PDEATHSIG, SIGKILL) != -1)
      << "Could not set PR_SET_PDEATHSIG";
  LOG(INFO) << "Started up container";
  // Keep reading execution requests in a loop in case we want to run more
  // commands in the same sandbox. Requests are written in the format
  // <number of bytes><request bytes>.
  while (true) {
    // Read execution request from the parent.
    ContainerExecution request;
    int length;
    if (!ReadIntFromFd(&length, in_id)) {
      break;
    }
    LOG(INFO) << "Request length: " << length;
    string request_bytes = ReadFromFd(length, in_id);
    LOG(INFO) << "Read request size: " << request_bytes.length();
    if (!request.ParseFromString(request_bytes)) {
      LOG(ERROR) << "Could not read complete request";
      break;
    }
    LOG(INFO) << "Received request " << request.DebugString();
    ContainerTermination response = omogen::sandbox::Execute(request);
    LOG(INFO) << "Done with termination " << response.DebugString();
    string response_bytes;
    response.SerializeToString(&response_bytes);
    WriteIntToFd(response_bytes.size(), out_id);
    WriteToFd(out_id, response_bytes);
  }
  PCHECK(close(out_id) != -1) << "Could not close output pipe";
  gflags::ShutDownCommandLineFlags();
}
