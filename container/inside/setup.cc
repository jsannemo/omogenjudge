#include <fcntl.h>
#include <sys/resource.h>
#include <sys/wait.h>
#include <unistd.h>
#include <iostream>
#include <vector>

#include "absl/strings/str_cat.h"
#include "api/omogenexec.pb.h"
#include "glog/logging.h"
#include "glog/raw_logging.h"
#include "proto/container.pb.h"
#include "util/files.h"

namespace omogenexec {

// Setup performs the main setup of a process that is to run an execution
// request, such as rlimits and stream redirections.

class InitException : public std::runtime_error {
  std::string msg;

 public:
  InitException(const std::string& msg)
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

static int openFileWithFd(const std::string& path, bool writable) {
  VLOG(2) << "Opening path " << path;
  int fd = writable ? open(path.c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666)
                    : open(path.c_str(), O_RDONLY);
  if (fd == -1) {
    throw InitException("open failed");
  }
  return fd;
}

static std::map<int, int> openStreams(const api::Streams& streams) {
  std::map<int, int> newFds;
  for (int fd = 0; fd < streams.mappings_size(); fd++) {
    const api::Streams::Mapping& mapping = streams.mappings(fd);
    switch (mapping.mapping_case()) {
      case api::Streams::Mapping::kEmpty:
        newFds[fd] = openFileWithFd("/dev/null", mapping.write());
        break;
      case api::Streams::Mapping::kFile:
        newFds[fd] = openFileWithFd(mapping.file().path_inside_container(),
                                    mapping.write());
        break;
      case api::Streams::Mapping::MAPPING_NOT_SET:
        assert(false && "Unsupported mapping");
    }
  }
  return newFds;
}

static void replaceStreams(const std::map<int, int>& newFds) {
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

static std::vector<const char*> getEnv(
    const google::protobuf::Map<std::string, std::string>& envToAdd) {
  std::vector<const char*> env;
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
    CloseFdsExcept(std::vector<int>{0, 1, 2, errorPipe});
    setResourceLimits();

    const api::Command& command = request.command();
    std::vector<const char*> argv;
    argv.push_back(request.command().command().c_str());
    for (int i = 0; i < command.flags_size(); ++i) {
      argv.push_back(command.flags(i).c_str());
    }
    argv.push_back(nullptr);
    std::vector<const char*> env = getEnv(request.environment().env());

    if (!FileIsExecutable(argv[0])) {
      throw InitException(
          "Command is not an executable file inside the sandbox");
    }
    std::map<int, int> newFds =
        openStreams(request.environment().stream_redirections());
    // Write a \1 that the parent will read to make sure we didn't crash
    // before we decided to close the error pipe.
    WriteToFd(errorPipe, "\1");
    // TODO(jsannemo): make sure we can wait with writing \1 after fixing file
    // descriptors by keeping the error stream at a high fd.
    replaceStreams(newFds);
    execve(argv[0], const_cast<char**>(argv.data()),
           const_cast<char**>(env.data()));
    exit(255);
  } catch (InitException e) {
    LOG(ERROR) << "Caught exception: " << e.what();
    WriteToFd(errorPipe, e.what());
    close(errorPipe);
    abort();
  }
}

}  // namespace omogenexec
