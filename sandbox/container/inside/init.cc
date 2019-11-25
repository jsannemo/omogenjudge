// Init is the init process (PID 1) of a container. Its purpose is to receive
// execution requests from the outside container and fork of a process to
// execute it.
#include <grp.h>
#include <pwd.h>
#include <sys/fcntl.h>
#include <sys/prctl.h>
#include <sys/resource.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include <cstdio>
#include <iostream>
#include <string>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "glog/logging.h"
#include "glog/raw_logging.h"
#include "sandbox/api/execspec.pb.h"
#include "sandbox/container/inside/chroot.h"
#include "sandbox/container/inside/quota.h"
#include "sandbox/container/inside/setup.h"
#include "sandbox/proto/container.pb.h"
#include "util/cpp/files.h"

using omogen::sandbox::proto::ContainerExecution;
using omogen::sandbox::proto::ContainerTermination;
using omogen::util::MakeDir;
using omogen::util::ReadFromFd;
using omogen::util::ReadIntFromFd;
using omogen::util::WriteIntToFd;
using omogen::util::WriteToFd;
using std::string;

namespace omogen {
namespace sandbox {

static const char* kContainerRoot = "/var/lib/omogen/sandbox";
static const int kChildStackSize = 100 * 1000;  // 100 KB
static const int kInodeLimit = 1000;

static struct passwd* pw;

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
    // De-privilege now. We must do this before execing, to keep rlimits we set.
    // They are cleared to default values if an elevated setuid execs.
    PCHECK(setresgid(omogen::sandbox::pw->pw_gid, omogen::sandbox::pw->pw_gid,
                     omogen::sandbox::pw->pw_gid) != -1)
        << "Could not set gid";
    PCHECK(setresuid(omogen::sandbox::pw->pw_uid, omogen::sandbox::pw->pw_uid,
                     omogen::sandbox::pw->pw_uid) != -1)
        << "Could not set uid";
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

static std::vector<char> stack(omogen::sandbox::kChildStackSize);
static int in_id;
static int out_id;
static int sandbox_id;
static string container_root;
static gid_t omogenclients_gid;

int startSandbox() {
  // Kill us if the main sandbox is killed, to prevent our child from possibly
  // keep running. This is not a race with the parent death, since the read
  // later will crash us in case our parent dies after the prctl call.
  // Furthermore, as a result of our death we will take with us any processes
  // running in the sandbox since we are PID 1 in a PID namespace.
  CHECK(prctl(PR_SET_PDEATHSIG, SIGKILL) != -1)
      << "Could not set PR_SET_PDEATHSIG";

  LOG(INFO) << "S" << sandbox_id << " Started up container";
  // Keep reading execution requests in a loop in case we want to run more
  // commands in the same sandbox. Requests are written in the format
  // <number of bytes><request bytes>.
  while (true) {
    // Read execution request from the parent.
    int length;
    ContainerExecution request;
    if (!ReadIntFromFd(&length, in_id)) {
      LOG(ERROR) << "S" << sandbox_id << " Failed reading length";
      return 1;
    }
    LOG(INFO) << "S" << sandbox_id << " Request length: " << length;
    string request_bytes = ReadFromFd(length, in_id);
    LOG(INFO) << "S" << sandbox_id
              << " Read request size: " << request_bytes.length();
    if (!request.ParseFromString(request_bytes)) {
      LOG(ERROR) << "Could not read complete request: "
                 << request.DebugString();
      return 1;
    }
    LOG(INFO) << "S" << sandbox_id << " Received request "
              << request.DebugString();
    ContainerTermination response = omogen::sandbox::Execute(request);
    LOG(INFO) << "S" << sandbox_id << " Done with termination "
              << response.DebugString();
    string response_bytes;
    response.SerializeToString(&response_bytes);
    WriteIntToFd(response_bytes.size(), out_id);
    WriteToFd(out_id, response_bytes);
  }
  LOG(INFO) << "S" << sandbox_id << " Closing";
  PCHECK(close(out_id) != -1)
      << "S" << sandbox_id << "Could not close output pipe";
  gflags::ShutDownCommandLineFlags();
  return 0;
}

int main(int argc, char** argv) {
  gflags::ParseCommandLineFlags(&argc, &argv, true);
  google::InitGoogleLogging(argv[0]);
  google::InstallFailureSignalHandler();
  setpgid(getpid(), getpid());
  CHECK(prctl(PR_SET_PDEATHSIG, SIGKILL) != -1)
      << "Could not set PR_SET_PDEATHSIG";

  CHECK(argc == 6) << "Incorrect number of arguments";
  CHECK(absl::SimpleAtoi(string(argv[1]), &sandbox_id))
      << "Can not convert sandbox ID to integer";
  CHECK(sandbox_id >= 0) << "Sandbox ID was negative";
  CHECK(absl::SimpleAtoi(string(argv[2]), &in_id))
      << "Can not convert input FD to integer";
  CHECK(absl::SimpleAtoi(string(argv[3]), &out_id))
      << "Can not convert output FD to integer";
  CHECK(in_id >= 0) << "Input FD was negative";
  CHECK(out_id >= 0) << "Ouptut FD was negative";
  int block_quota, inode_quota;
  CHECK(absl::SimpleAtoi(std::string(argv[4]), &block_quota))
      << "Can not convert block quota to integer";
  CHECK(block_quota >= 0) << "Block quota was negative";
  CHECK(absl::SimpleAtoi(std::string(argv[5]), &inode_quota))
      << "Can not convert Inode quota to integer";
  CHECK(inode_quota >= 0) << "Inode quota was negative";

  container_root =
      absl::StrCat(omogen::sandbox::kContainerRoot, "/", sandbox_id);
  MakeDir(container_root);
  struct group* group = getgrnam("omogenjudge-clients");
  omogenclients_gid = group->gr_gid;
  std::string user = absl::StrCat("omogenjudge-client", sandbox_id);
  omogen::sandbox::pw = getpwnam(user.c_str());
  PCHECK(omogen::sandbox::pw != NULL) << "Could not fetch user";
  omogen::sandbox::SetQuota(block_quota, inode_quota, omogen::sandbox::pw->pw_uid);
  PCHECK(chown(container_root.c_str(), omogen::sandbox::pw->pw_uid, group->gr_gid) != -1)
      << "Could not chown container root";
  PCHECK(chmod(container_root.c_str(), 0775) != -1)
      << "Could not chmod container root";
  struct passwd* sandbox_pw = getpwnam("omogenjudge-sandbox");

  unshare(CLONE_NEWNS);
  unshare(CLONE_NEWIPC);
  unshare(CLONE_NEWNET);
  unshare(CLONE_NEWPID);
  unshare(CLONE_NEWUTS);
  omogen::sandbox::ContainerSpec spec;
  int length;
  CHECK(ReadIntFromFd(&length, in_id)) << "Failed reading length";
  string spec_bytes = ReadFromFd(length, in_id);
  CHECK(spec.ParseFromString(spec_bytes))
      << "Could not read complete request: " << spec.DebugString();
  omogen::sandbox::Chroot chroot =
      omogen::sandbox::Chroot::ForNewRoot(container_root);
  chroot.ApplyContainerSpec(spec);
  chroot.SetRoot();

  PCHECK(prctl(PR_SET_PDEATHSIG, SIGKILL) != -1)
      << "Could not set PR_SET_PDEATHSIG";
  startSandbox();
}
