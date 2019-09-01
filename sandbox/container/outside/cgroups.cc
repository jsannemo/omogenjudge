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
static const string kSubSystemName[] = {"cpuacct", "memory", "invalid"};

static bool ValidateCgroupPath(const char* flagname, const string& value) {
  if (value.empty() || value[0] != '/' || !DirectoryExists(value)) {
    LOG(FATAL) << "Invalid cgroup root -" << flagname << ": " << value
               << " does not exist";
    return false;
  }
  return true;
}
DEFINE_string(cgroup_root, "/sys/fs/cgroup",
              "The root of the cgroup file system");
DEFINE_validator(cgroup_root, &ValidateCgroupPath);

static bool ValidateCgroupParent(const char* flagname, const string& value) {
  for (int i = 0; i < static_cast<int>(CgroupSubsystem::INVALID); ++i) {
    string path =
        FLAGS_cgroup_root + "/" + kSubSystemName[i] + "/" + value + "/";
    if (!DirectoryExists(path)) {
      LOG(FATAL) << "Cgroup parent -" << flagname << ": " << value
                 << " does not contain subsystem " << kSubSystemName[i];
      return false;
    }
  }
  return true;
}
DEFINE_string(parent_cgroup, "omogencontain",
              "The name of the parent cgroup that will be used. The user "
              "executing the container must have read-write access");
DEFINE_validator(parent_cgroup, &ValidateCgroupParent);

DEFINE_string(cgroup_prefix, "omogen_",
              "A prefix used to name the cgroups to avoid collisions");

const string kCpuUsage = "cpuacct.usage";
const string kMemLimit = "memory.limit_in_bytes";
const string kMemUsage = "memory.max_usage_in_bytes";
const string kTasks = "tasks";
const string kNotify = "notify_on_release";

static int IndexForSubsystem(CgroupSubsystem subsystem) {
  return static_cast<int>(subsystem);
}

static string GetCgroupName(pid_t pid) {
  return absl::StrCat(FLAGS_cgroup_prefix, pid);
}

string Cgroup::GetSubsystemPath(CgroupSubsystem subsystem) {
  return FLAGS_cgroup_root + "/" +
         kSubSystemName[IndexForSubsystem(subsystem)] + "/" +
         FLAGS_parent_cgroup + "/" + name;
}

string Cgroup::GetSubsystemOp(CgroupSubsystem subsystem, const string& op) {
  return GetSubsystemPath(subsystem) + "/" + op;
}

void Cgroup::EnableSubsystem(CgroupSubsystem subsystem) {
  MakeDir(GetSubsystemPath(subsystem));
  WriteToFile(GetSubsystemOp(subsystem, kTasks), absl::StrCat(pid));
  WriteToFile(GetSubsystemOp(subsystem, kNotify), "1");
}

void Cgroup::DisableSubsystem(CgroupSubsystem subsystem) {
  string path = FLAGS_cgroup_root + "/" +
                kSubSystemName[IndexForSubsystem(subsystem)] + "/" + kTasks;
  TryRemoveDir(GetSubsystemPath(subsystem));
}

long long Cgroup::CpuUsed() {
  vector<string> tokens =
      TokenizeFile(GetSubsystemOp(CgroupSubsystem::CPU_ACCT, kCpuUsage));
  CHECK(!tokens.empty()) << "CPU usage file for cgroup was empty";
  long long cpu_ns;
  CHECK(absl::SimpleAtoi(tokens[0], &cpu_ns));
  return cpu_ns / 1000000;
}

void Cgroup::SetMemoryLimit(long long mem_limit_kb) {
  VLOG(2) << "Setting memory limit to " << mem_limit_kb;
  if (mem_limit_kb == _mem_limit_kb) {
    return;
  }
  _mem_limit_kb = mem_limit_kb;
  CHECK(mem_limit_kb >= 0) << "Memory limit was negative: " << mem_limit_kb;
  WriteToFile(GetSubsystemOp(CgroupSubsystem::MEMORY, kMemLimit),
              absl::StrCat(mem_limit_kb * 1000));
}

long long Cgroup::MemoryUsed() {
  vector<string> tokens =
      TokenizeFile(GetSubsystemOp(CgroupSubsystem::MEMORY, kMemUsage));
  CHECK(!tokens.empty()) << "CPU usage file for cgroup was empty";
  long long bytes;
  CHECK(absl::SimpleAtoi(tokens[0], &bytes));
  return bytes / 1000;
}

void Cgroup::Reset() {
  WriteToFile(GetSubsystemOp(CgroupSubsystem::CPU_ACCT, kCpuUsage), "0");
  WriteToFile(GetSubsystemOp(CgroupSubsystem::MEMORY, kMemUsage), "0");
}

Cgroup::Cgroup(pid_t pid) : name(GetCgroupName(pid)), pid(pid) {
  EnableSubsystem(CgroupSubsystem::CPU_ACCT);
  EnableSubsystem(CgroupSubsystem::MEMORY);
  Reset();
}

Cgroup::~Cgroup() {
  VLOG(3) << "Removing cgroups for " << pid;
  DisableSubsystem(CgroupSubsystem::CPU_ACCT);
  DisableSubsystem(CgroupSubsystem::MEMORY);
}

}  // namespace sandbox
}  // namespace omogen
