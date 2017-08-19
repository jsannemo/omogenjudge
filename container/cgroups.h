#include <string>

#include "proto/omogenexec.pb.h"

using std::string;

// A Cgroup puts the given process into a new control group, tracking its resource usage
// and limits the memory used.
class Cgroup {

    string name;
    pid_t pid;

    string getSubsystemPath(const string& subsystem);
    void enableSubsystem(const string& subsystemPath);
    void disableSubsystem(const string& subsystemPath);

public:
    Cgroup(pid_t pid);
    ~Cgroup();

    // The total CPU usage of the process and its children, in milliseconds.
    long long CpuUsed();
    // The total memory usage of the process and its children, in kilobytes.
    long long MemoryUsed();
    void SetMemoryLimit(long long memLimitKb);
    void SetProcessLimit(int maxProcesses);
    // The total disk I/O usage, in kilobytes.
    long long BytesTransferred();
    // Reset the resource usage statistics
    void Reset();
};
