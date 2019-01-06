#include <map>
#include <sstream>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "container/outside/cgroups.h"
#include "gflags/gflags.h"
#include "glog/logging.h"
#include "util/files.h"

namespace omogenexec {

using std::string;

// This names of the cgroup sybsystems needs to be kept in sync with the
// enum in cgroups.h
static const std::string subsystemName[] = {"cpuacct", "memory", "pids",
                                            "invalid"};

static bool validateCgroupPath(const char* flagname, const std::string& value) {
  if (value.empty() || value[0] != '/' || !DirectoryExists(value)) {
    LOG(FATAL) << "Invalid cgroup root -" << flagname << ": " << value
               << " does not exist";
    return false;
  }
  return true;
}
DEFINE_string(cgroup_root, "/sys/fs/cgroup",
              "The root of the cgroup file system");
DEFINE_validator(cgroup_root, &validateCgroupPath);

static bool validateCgroupParent(const char* flagname,
                                 const std::string& value) {
  for (int i = 0; i < static_cast<int>(CgroupSubsystem::INVALID); ++i) {
    std::string path =
        FLAGS_cgroup_root + "/" + subsystemName[i] + "/" + value + "/";
    if (!DirectoryExists(path)) {
      LOG(FATAL) << "Cgroup parent -" << flagname << ": " << value
                 << " does not contain subsystem " << subsystemName[i];
      return false;
    }
  }
  return true;
}
DEFINE_string(parent_cgroup, "omogencontain",
              "The name of the parent cgroup that will be used. The user "
              "executing the container must have read-write access");
DEFINE_validator(parent_cgroup, &validateCgroupParent);

DEFINE_string(cgroup_prefix, "omogen_",
              "A prefix used to name the cgroups to avoid collisions");

const std::string CPU_USAGE = "cpuacct.usage";
const std::string MEM_LIMIT = "memory.limit_in_bytes";
const std::string MEM_USAGE = "memory.max_usage_in_bytes";
const std::string PID_LIMIT = "pids.max";
const std::string TASKS = "tasks";

static int indexForSubsystem(CgroupSubsystem subsystem) {
  return static_cast<int>(subsystem);
}

static std::string getCgroupName(pid_t pid) {
  return absl::StrCat(FLAGS_cgroup_prefix, pid);
}

std::string Cgroup::getSubsystemPath(CgroupSubsystem subsystem) {
  return FLAGS_cgroup_root + "/" + subsystemName[indexForSubsystem(subsystem)] +
         "/" + FLAGS_parent_cgroup + "/" + name;
}

std::string Cgroup::getSubsystemOp(CgroupSubsystem subsystem,
                                   const std::string& op) {
  return getSubsystemPath(subsystem) + "/" + op;
}

void Cgroup::enableSubsystem(CgroupSubsystem subsystem) {
  MakeDir(getSubsystemPath(subsystem));
  WriteToFile(getSubsystemOp(subsystem, TASKS), absl::StrCat(pid));
}

void Cgroup::disableSubsystem(CgroupSubsystem subsystem) {
  RemoveDir(getSubsystemPath(subsystem));
}

long long Cgroup::CpuUsed() {
  std::vector<std::string> tokens =
      TokenizeFile(getSubsystemOp(CgroupSubsystem::CPU_ACCT, CPU_USAGE));
  CHECK(!tokens.empty()) << "CPU usage file for cgroup was empty";
  long long nanoSeconds;
  CHECK(absl::SimpleAtoi(tokens[0], &nanoSeconds));
  return nanoSeconds / 1000000;
}

void Cgroup::SetMemoryLimit(long long memLimitKb) {
  VLOG(2) << "Setting memory limit to " << memLimitKb;
  CHECK(memLimitKb >= 0) << "Memory limit was negative: " << memLimitKb;
  WriteToFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_LIMIT),
              absl::StrCat(memLimitKb * 1000));
}

long long Cgroup::MemoryUsed() {
  std::vector<std::string> tokens =
      TokenizeFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_USAGE));
  CHECK(!tokens.empty()) << "CPU usage file for cgroup was empty";
  long long bytes;
  CHECK(absl::SimpleAtoi(tokens[0], &bytes));
  return bytes / 1000;
}

void Cgroup::SetProcessLimit(int maxProcesses) {
  VLOG(2) << "Setting process limit to " << maxProcesses;
  CHECK(maxProcesses >= 0) << "Process limit was negative: " << maxProcesses;
  WriteToFile(getSubsystemOp(CgroupSubsystem::PIDS, PID_LIMIT),
              absl::StrCat(maxProcesses));
}

void Cgroup::Reset() {
  WriteToFile(getSubsystemOp(CgroupSubsystem::CPU_ACCT, CPU_USAGE), "0");
  WriteToFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_USAGE), "0");
}

Cgroup::Cgroup(pid_t pid) : name(getCgroupName(pid)), pid(pid) {
  enableSubsystem(CgroupSubsystem::CPU_ACCT);
  enableSubsystem(CgroupSubsystem::MEMORY);
  enableSubsystem(CgroupSubsystem::PIDS);
}

Cgroup::~Cgroup() {
  VLOG(3) << "Removing cgroups for " << pid;
  disableSubsystem(CgroupSubsystem::CPU_ACCT);
  disableSubsystem(CgroupSubsystem::MEMORY);
  disableSubsystem(CgroupSubsystem::PIDS);
}

}  // namespace omogenexec
