#include "sandbox/container/outside/cgroups.h"

#include <map>
#include <sstream>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "gflags/gflags.h"
#include "glog/logging.h"
#include "util/cpp/files.h"

namespace omogen {
namespace sandbox {

using omogen::util::DirectoryExists;
using omogen::util::MakeDir;
using omogen::util::TokenizeFile;
using omogen::util::TryRemoveDir;
using omogen::util::WriteToFile;
using std::string;
using std::vector;

// This names of the cgroup sybsystems needs to be kept in sync with the
// enum in cgroups.h
static const string subsystemName[] = {"cpuacct", "memory", "invalid"};

static bool validateCgroupPath(const char* flagname, const string& value) {
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

static bool validateCgroupParent(const char* flagname, const string& value) {
  for (int i = 0; i < static_cast<int>(CgroupSubsystem::INVALID); ++i) {
    string path =
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

const string CPU_USAGE = "cpuacct.usage";
const string MEM_LIMIT = "memory.limit_in_bytes";
const string MEM_USAGE = "memory.max_usage_in_bytes";
const string TASKS = "tasks";

static int indexForSubsystem(CgroupSubsystem subsystem) {
  return static_cast<int>(subsystem);
}

static string getCgroupName(pid_t pid) {
  return absl::StrCat(FLAGS_cgroup_prefix, pid);
}

string Cgroup::getSubsystemPath(CgroupSubsystem subsystem) {
  return FLAGS_cgroup_root + "/" + subsystemName[indexForSubsystem(subsystem)] +
         "/" + FLAGS_parent_cgroup + "/" + name;
}

string Cgroup::getSubsystemOp(CgroupSubsystem subsystem, const string& op) {
  return getSubsystemPath(subsystem) + "/" + op;
}

void Cgroup::enableSubsystem(CgroupSubsystem subsystem) {
  MakeDir(getSubsystemPath(subsystem));
  WriteToFile(getSubsystemOp(subsystem, TASKS), absl::StrCat(pid));
}

void Cgroup::disableSubsystem(CgroupSubsystem subsystem) {
  string path = FLAGS_cgroup_root + "/" +
                subsystemName[indexForSubsystem(subsystem)] + "/" + TASKS;
  TryRemoveDir(getSubsystemPath(subsystem));
}

long long Cgroup::CpuUsed() {
  vector<string> tokens =
      TokenizeFile(getSubsystemOp(CgroupSubsystem::CPU_ACCT, CPU_USAGE));
  CHECK(!tokens.empty()) << "CPU usage file for cgroup was empty";
  long long nanoSeconds;
  CHECK(absl::SimpleAtoi(tokens[0], &nanoSeconds));
  return nanoSeconds / 1000000;
}

void Cgroup::SetMemoryLimit(long long memLimitKb) {
  VLOG(2) << "Setting memory limit to " << memLimitKb;
  if (memLimitKb == _memLimitKb) {
    return;
  }
  _memLimitKb = memLimitKb;
  CHECK(memLimitKb >= 0) << "Memory limit was negative: " << memLimitKb;
  WriteToFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_LIMIT),
              absl::StrCat(memLimitKb * 1000));
}

long long Cgroup::MemoryUsed() {
  vector<string> tokens =
      TokenizeFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_USAGE));
  CHECK(!tokens.empty()) << "CPU usage file for cgroup was empty";
  long long bytes;
  CHECK(absl::SimpleAtoi(tokens[0], &bytes));
  return bytes / 1000;
}

void Cgroup::Reset() {
  WriteToFile(getSubsystemOp(CgroupSubsystem::CPU_ACCT, CPU_USAGE), "0");
  WriteToFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_USAGE), "0");
}

Cgroup::Cgroup(pid_t pid) : name(getCgroupName(pid)), pid(pid) {
  enableSubsystem(CgroupSubsystem::CPU_ACCT);
  enableSubsystem(CgroupSubsystem::MEMORY);
  Reset();
}

Cgroup::~Cgroup() {
  VLOG(3) << "Removing cgroups for " << pid;
  disableSubsystem(CgroupSubsystem::CPU_ACCT);
  disableSubsystem(CgroupSubsystem::MEMORY);
}

}  // namespace sandbox
}  // namespace omogen
