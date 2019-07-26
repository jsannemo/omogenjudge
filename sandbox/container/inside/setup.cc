#include "sandbox/container/inside/setup.h"

#include <fcntl.h>
#include <sys/resource.h>
#include <sys/wait.h>
#include <unistd.h>

#include <iostream>
#include <vector>

#include "absl/strings/str_cat.h"
#include "glog/logging.h"
#include "glog/raw_logging.h"
#include "sandbox/api/exec.pb.h"
#include "sandbox/proto/container.pb.h"
#include "util/cpp/files.h"

using omogen::util::CloseFdsExcept;
using omogen::util::FileIsExecutable;
using omogen::util::WriteToFd;
using std::map;
using std::string;
using std::vector;

namespace omogen {
namespace sandbox {

// Setup performs the main setup of a process that is to run an execution
// request, such as rlimits and stream redirections.

class InitException : public std::runtime_error {
  string msg;

 public:
  InitException(const string& msg)
      : runtime_error("Container failed to setup execution"), msg(msg) {}
  const char* what() const noexcept override { return msg.c_str(); }
};

static void setResourceLimit(int resource, rlim_t limit) {
  rlimit rlim = {.rlim_cur = limit, .rlim_max = limit};
  if (setrlimit(resource, &rlim) == -1) {
    throw InitException("setrlimit failed");
  }
}

static void setResourceLimits() {
  setResourceLimit(RLIMIT_STACK, RLIM_INFINITY);
  setResourceLimit(RLIMIT_MEMLOCK, 0);
  setResourceLimit(RLIMIT_CORE, 0);
  setResourceLimit(RLIMIT_NOFILE, 100);
}

static int openFileWithFd(const string& path, bool writable) {
  VLOG(2) << "Opening path " << path;
  int fd = writable ? open(path.c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666)
                    : open(path.c_str(), O_RDONLY);
  if (fd == -1) {
    throw InitException("open failed");
  }
  return fd;
}

static map<int, int> openStreams(const Streams& streams) {
  map<int, int> newFds;
  for (int fd = 0; fd < streams.mappings_size(); fd++) {
    const Streams::Mapping& mapping = streams.mappings(fd);
    switch (mapping.mapping_case()) {
      case Streams::Mapping::kEmpty:
        newFds[fd] = openFileWithFd("/dev/null", mapping.write());
        break;
      case Streams::Mapping::kFile:
        newFds[fd] = openFileWithFd(mapping.file().path_inside_container(),
                                    mapping.write());
        break;
      case Streams::Mapping::MAPPING_NOT_SET:
        assert(false && "Unsupported mapping");
    }
  }
  return newFds;
}

static void replaceStreams(const map<int, int>& newFds) {
  // TODO(jsannemo): this is a bit shaky if we are mapping more
  // fds than 0, 1, 2 since we don't know which file descriptors
  // we got when we opened the above, or what file descriptor
  // the error had. They should be chosen in a manner that won't conflict
  // with the low fds we get when opening.
  for (const auto& toReplace : newFds) {
    if (dup2(toReplace.second, toReplace.first) == -1) {
      throw InitException("Could not open new file descriptor");
    }
    if (close(toReplace.second) == -1) {
      throw InitException("Could not close old file descriptor");
    }
  }
}

static vector<const char*> getEnv(
    const google::protobuf::Map<string, string>& envToAdd) {
  vector<const char*> env;
  // Path is needed for e.g. gcc, which searchs for some binaries in the path
  env.push_back("PATH=/bin:/usr/bin");
  for (const auto& variable : envToAdd) {
    env.push_back(
        strdup(absl::StrCat(variable.first, "=", variable.second).c_str()));
  }
  env.push_back(nullptr);
  return env;
}

[[noreturn]] void SetupAndRun(const proto::ContainerExecution& request,
                              int errorPipe) {
  try {
    // We close all file descriptors to prevent leaks from the parent
    CloseFdsExcept(vector<int>{0, 1, 2, errorPipe});
    setResourceLimits();

    const Command& command = request.command();
    const char** argv = new const char*[2 + command.flags_size()];
    argv[0] = command.command().c_str();
    argv[1] = nullptr;
    for (int i = 0; i < command.flags_size(); ++i) {
      argv[i + 1] = strdup(command.flags(i).c_str());
    }
    argv[command.flags_size() + 1] = nullptr;
    vector<const char*> env = getEnv(request.environment().env());

    if (!FileIsExecutable(argv[0])) {
      throw InitException(
          "Command is not an executable file inside the sandbox");
    }
    map<int, int> newFds =
        openStreams(request.environment().stream_redirections());
    if (!request.environment().working_directory().empty()) {
      PCHECK(chdir(request.environment().working_directory().c_str()) == 0);
    }
    // Write a \1 that the parent will read to make sure we didn't crash
    // before we decided to close the error pipe.
    WriteToFd(errorPipe, "\1");
    // TODO(jsannemo): make sure we can wait with writing \1 after fixing file
    // descriptors by keeping the error stream at a high fd.
    replaceStreams(newFds);
    execve(argv[0], const_cast<char**>(argv), const_cast<char**>(env.data()));
    exit(255);
  } catch (InitException e) {
    LOG(ERROR) << "Caught exception: " << e.what();
    WriteToFd(errorPipe, e.what());
    close(errorPipe);
    abort();
  }
}

}  // namespace sandbox
}  // namespace omogen
