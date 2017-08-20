#pragma once

#include <string>
#include <sys/types.h>

#include "proto/omogenexec.pb.h"

namespace omogenexec {

// Note: if adding a new subsystem, its name needs to be added to the list of
// names subsystemName in cgroups.cc
enum class CgroupSubsystem {
    BLKIO = 0,
    CPU_ACCT,
    MEMORY,
    PIDS,
};


// A cgroup enables tracking of a process' resource usage and limiting the memory it can use.
class Cgroup {

    std::string name;
    pid_t pid;

    std::string getSubsystemPath(CgroupSubsystem subsystem);
    std::string getSubsystemOp(CgroupSubsystem subsystem, const std::string& op);
    void enableSubsystem(CgroupSubsystem subsystem);
    void disableSubsystem(CgroupSubsystem subsystem);

public:

    // Creates a new cgroup for a given process and enables the subsystems used for the process.
    Cgroup(pid_t pid);
    ~Cgroup();

    // The total CPU usage of the process and its children, in milliseconds.
    long long CpuUsed();
    // The total memory usage of the process and its children, in kilobytes.
    long long MemoryUsed();
    void SetMemoryLimit(long long memLimitKb);
    // The total disk I/O usage, in kilobytes.
    long long DiskIOUsed();
    long long ProcessesUsed();
    void SetProcessLimit(int maxProcesses);
    // Reset the resource usage statistics.
    void Reset();

    static Cgroup MakeCgroupFor(pid_t pid);

};

}
