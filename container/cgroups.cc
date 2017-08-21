#include <gflags/gflags.h>
#include <map>
#include <sstream>
#include <vector>

#include "cgroups.h"
#include "util/error.h"
#include "util/files.h"
#include "util/format.h"
#include "util/log.h"

namespace omogenexec {

using std::map;
using std::string;
using std::stringstream;
using std::vector;

// TODO(#2): write flag validators
DEFINE_string(cgroup_root, "/sys/fs/cgroup", "The root of the cgroup file system");
DEFINE_string(parent_cgroup, "omogencontain", "The name of the parent cgroup that will be used. The user executing the container must have read-write access");
DEFINE_string(cgroup_prefix, "omogen_", "A prefix used to name the cgroups to avoid collisions");

const string IO_RESET = "blkio.reset_stats";
const string IO_USAGE = "blkio.throttle.io_service_bytes";
const string CPU_USAGE = "cpuacct.usage";
const string MEM_LIMIT = "memory.limit_in_bytes";
const string MEM_USAGE = "memory.max_usage_in_bytes";
const string PID_LIMIT = "pids.max";
const string TASKS = "tasks";

static int sub2idx(CgroupSubsystem subsystem) {
    return static_cast<int>(subsystem);
}

static const string subsystemName[] = {
    "blkio",
    "cpuacct",
    "memory",
    "pids",
};

static string getCgroupName(pid_t pid) {
    stringstream ss;
    ss << FLAGS_cgroup_prefix << pid;
    return ss.str();
}

string Cgroup::getSubsystemPath(CgroupSubsystem subsystem) {
    return FLAGS_cgroup_root + "/" + subsystemName[sub2idx(subsystem)] + "/omogencontain/" + name;
}

string Cgroup::getSubsystemOp(CgroupSubsystem subsystem, const string& op) {
    return getSubsystemPath(subsystem) + "/" + op;
}

void Cgroup::enableSubsystem(CgroupSubsystem subsystem) {
    MakeDir(getSubsystemPath(subsystem));
    stringstream ss;
    ss << pid;
    WriteToFile(getSubsystemOp(subsystem, TASKS), ss.str());
}

void Cgroup::disableSubsystem(CgroupSubsystem subsystem) {
    RemoveDir(getSubsystemPath(subsystem));
}

long long Cgroup::CpuUsed() {
    vector<string> tokens = TokenizeFile(getSubsystemOp(CgroupSubsystem::CPU_ACCT, CPU_USAGE));
    assert(!tokens.empty());
    long long nanoSeconds = StringToLL(tokens[0]);
    return nanoSeconds / 1000000;
}

void Cgroup::SetMemoryLimit(long long memLimitKb) {
    stringstream memoryLimit;
    memoryLimit << memLimitKb * 1000;
    WriteToFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_LIMIT), memoryLimit.str());
}

long long Cgroup::MemoryUsed() {
    vector<string> tokens = TokenizeFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_USAGE));
    assert(!tokens.empty());
    long long bytes = StringToLL(tokens[0]);
    return bytes / 1000;
}

void Cgroup::SetProcessLimit(int maxProcesses) {
    stringstream pidLimit;
    pidLimit << maxProcesses;
    WriteToFile(getSubsystemOp(CgroupSubsystem::PIDS, PID_LIMIT), pidLimit.str());
}

long long Cgroup::DiskIOUsed() {
    vector<string> tokens = TokenizeFile(getSubsystemOp(CgroupSubsystem::BLKIO, IO_USAGE));
    assert(!tokens.empty());
    long long bytes = StringToLL(tokens.back());
    return bytes / 1000;
}

void Cgroup::Reset() {
    WriteToFile(getSubsystemOp(CgroupSubsystem::CPU_ACCT, CPU_USAGE), "0");
    WriteToFile(getSubsystemOp(CgroupSubsystem::MEMORY, MEM_USAGE), "0");
    WriteToFile(getSubsystemOp(CgroupSubsystem::BLKIO, IO_RESET), "0");
}

Cgroup::Cgroup(pid_t pid) : name(getCgroupName(pid)), pid(pid) {
    enableSubsystem(CgroupSubsystem::CPU_ACCT);
    enableSubsystem(CgroupSubsystem::MEMORY);
    enableSubsystem(CgroupSubsystem::PIDS);
    enableSubsystem(CgroupSubsystem::BLKIO);
}

Cgroup::~Cgroup() {
    disableSubsystem(CgroupSubsystem::CPU_ACCT);
    disableSubsystem(CgroupSubsystem::MEMORY);
    disableSubsystem(CgroupSubsystem::PIDS);
    disableSubsystem(CgroupSubsystem::BLKIO);
}

}
