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

DEFINE_string(cgroup_root, "/sys/fs/cgroup", "The root of the cgroup file system");
DEFINE_string(parent_cgroup, "omogencontain", "The name of the parent cgroup that will be used. The user executing the container must have read-write access");
DEFINE_string(cgroup_prefix, "omogen_", "A prefix used to name the cgroups to avoid collisions");

const string CPU_USAGE = "cpuacct.usage";
const string MEM_USAGE = "memory.max_usage_in_bytes";
const string MEM_LIMIT = "memory.limit_in_bytes";
const string IO_USAGE = "blkio.throttle.io_service_bytes";
const string IO_RESET = "blkio.reset_stats";
const string PID_LIMIT = "pids.max";

static map<CgroupSubsystem, string> subsystemNames() {
    map<CgroupSubsystem, string> ret;
    ret[CPU_ACCT] = "cpuacct";
    ret[MEMORY] = "memory";
    ret[PIDS] = "pids";
    ret[BLKIO] = "blkio";
    return ret;
}

static map<CgroupSubsystem, string> subsystemName = subsystemNames();

static string getCgroupName(pid_t pid) {
    stringstream ss;
    ss << FLAGS_cgroup_prefix << pid;
    return ss.str();
}

string Cgroup::getSubsystemPath(CgroupSubsystem subsystem) {
    return FLAGS_cgroup_root + "/" + subsystemName[subsystem] + "/omogencontain/" + name;
}

string Cgroup::getSubsystemOp(CgroupSubsystem subsystem, const string& op) {
    return getSubsystemPath(subsystem) + "/" + op;
}

void Cgroup::enableSubsystem(CgroupSubsystem subsystem) {
    MakeDir(getSubsystemPath(subsystem));
    stringstream ss;
    ss << pid;
    WriteToFile(getSubsystemOp(subsystem, "/tasks"), ss.str());
}

void Cgroup::disableSubsystem(CgroupSubsystem subsystem) {
    RemoveDir(getSubsystemPath(subsystem));
}

long long Cgroup::CpuUsed() {
    vector<string> tokens = TokenizeFile(getSubsystemOp(CPU_ACCT, CPU_USAGE));
    assert(!tokens.empty());
    long long nanoSeconds = StringToLL(tokens[0]);
    return nanoSeconds / 1000000;
}

void Cgroup::SetMemoryLimit(long long memLimitKb) {
    stringstream memoryLimit;
    memoryLimit << memLimitKb * 1000;
    WriteToFile(getSubsystemOp(MEMORY, MEM_LIMIT), memoryLimit.str());
}

long long Cgroup::MemoryUsed() {
    vector<string> tokens = TokenizeFile(getSubsystemOp(MEMORY, MEM_USAGE));
    assert(!tokens.empty());
    long long bytes = StringToLL(tokens[0]);
    return bytes / 1000;
}

void Cgroup::SetProcessLimit(int maxProcesses) {
    stringstream pidLimit;
    pidLimit << maxProcesses;
    WriteToFile(getSubsystemOp(PIDS, PID_LIMIT), pidLimit.str());
}

long long Cgroup::DiskIOUsed() {
    vector<string> tokens = TokenizeFile(getSubsystemOp(BLKIO, IO_USAGE));
    assert(!tokens.empty());
    long long bytes = StringToLL(tokens.back());
    return bytes / 1000;
}

void Cgroup::Reset() {
    WriteToFile(getSubsystemOp(CPU_ACCT, CPU_USAGE), "0");
    WriteToFile(getSubsystemOp(MEMORY, MEM_USAGE), "0");
    WriteToFile(getSubsystemOp(BLKIO, IO_RESET), "0");
}

Cgroup::Cgroup(pid_t pid) : name(getCgroupName(pid)), pid(pid) {
    enableSubsystem(CPU_ACCT);
    enableSubsystem(MEMORY);
    enableSubsystem(PIDS);
    enableSubsystem(BLKIO);
}

Cgroup::~Cgroup() {
    disableSubsystem(CPU_ACCT);
    disableSubsystem(MEMORY);
    disableSubsystem(PIDS);
    disableSubsystem(BLKIO);
}

}
