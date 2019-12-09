// Setup performs the main setup of a process that is to run an execution
// request, such as rlimits and stream redirections.
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
#include "sandbox/api/execute_service.pb.h"
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

class InitException : public std::runtime_error {
  string msg;

 public:
  InitException(const string& msg)
      : runtime_error("Container failed to setup execution"), msg(msg) {}
  const char* what() const noexcept override { return msg.c_str(); }
};

static void SetResourceLimit(int resource, rlim_t limit) {
  rlimit rlim = {.rlim_cur = limit, .rlim_max = limit};
  if (setrlimit(resource, &rlim) == -1) {
    throw InitException("setrlimit failed");
  }
}

static void SetResourceLimits() {
  SetResourceLimit(RLIMIT_AS, RLIM_INFINITY);
  SetResourceLimit(RLIMIT_STACK, RLIM_INFINITY);
  SetResourceLimit(RLIMIT_MEMLOCK, 0);
  SetResourceLimit(RLIMIT_CORE, 0);
}

static int OpenFileWithFd(const string& path, bool writable) {
  RAW_VLOG(2, "Opening path %s", path.c_str());
  int fd = writable ? open(path.c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666)
                    : open(path.c_str(), O_RDONLY);
  if (fd == -1) {
    throw InitException(absl::StrCat("opening ", path, " failed"));
  }
  return fd;
}

static map<int, int> OpenStreams(const Streams& streams) {
  map<int, int> new_fds;
  for (int fd = 0; fd < streams.mappings_size(); fd++) {
    const Streams::Mapping& mapping = streams.mappings(fd);
    switch (mapping.mapping_case()) {
      case Streams::Mapping::kFile:
        new_fds[fd] = OpenFileWithFd(mapping.file().path_inside_container(),
                                     mapping.write());
        break;
      case Streams::Mapping::MAPPING_NOT_SET:
        assert(false && "Unsupported mapping");
    }
  }
  return new_fds;
}

static void ReplaceStreams(const map<int, int>& new_fds) {
  // TODO(jsannemo): this is a bit shaky if we are mapping more
  // fds than 0, 1, 2 since we don't know which file descriptors
  // we got when we opened the above, or what file descriptor
  // the error had. They should be chosen in a manner that won't conflict
  // with the low fds we get when opening.
  for (const auto& to_replace : new_fds) {
    if (dup2(to_replace.second, to_replace.first) == -1) {
      throw InitException("Could not open new file descriptor");
    }
    if (close(to_replace.second) == -1) {
      throw InitException("Could not close old file descriptor");
    }
  }
}

static vector<const char*> GetEnv(
    const google::protobuf::Map<string, string>& env_to_add) {
  vector<const char*> env;
  // Path is needed for e.g. gcc, which searchs for some binaries in the path
  env.push_back("PATH=/bin:/usr/bin");
  for (const auto& variable : env_to_add) {
    env.push_back(
        strdup(absl::StrCat(variable.first, "=", variable.second).c_str()));
  }
  env.push_back(nullptr);
  return env;
}

[[noreturn]] void SetupAndRun(const proto::ContainerExecution& request,
                              int error_pipe) {
  try {
    // We close all file descriptors to prevent leaks from the parent
    CloseFdsExcept(vector<int>{0, 1, 2, error_pipe});
    SetResourceLimits();

    const Command& command = request.command();
    const char** argv = new const char*[2 + command.flags_size()];
    argv[0] = command.command().c_str();
    argv[1] = nullptr;
    for (int i = 0; i < command.flags_size(); ++i) {
      argv[i + 1] = strdup(command.flags(i).c_str());
    }
    argv[command.flags_size() + 1] = nullptr;
    vector<const char*> env = GetEnv(request.environment().env());

    if (!FileIsExecutable(argv[0])) {
      throw InitException(
          "Command is not an executable file inside the sandbox");
    }
    map<int, int> new_fds =
        OpenStreams(request.environment().stream_redirections());
    if (!request.environment().working_directory().empty()) {
      RAW_CHECK(chdir(request.environment().working_directory().c_str()) == 0,
                "could not set working directory");
    }
    // Write a \1 that the parent will read to make sure we didn't crash
    // before we decided to close the error pipe.
    WriteToFd(error_pipe, "\1");
    // TODO(jsannemo): make sure we can wait with writing \1 after fixing file
    // descriptors by keeping the error stream at a high fd.
    ReplaceStreams(new_fds);
    execve(argv[0], const_cast<char**>(argv), const_cast<char**>(env.data()));
    exit(255);
  } catch (InitException* e) {
    RAW_LOG(ERROR, "Caught exception: %s", e->what());
    WriteToFd(error_pipe, e->what());
    close(error_pipe);
    abort();
  }
}

}  // namespace sandbox
}  // namespace omogen
