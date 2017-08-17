#include <fcntl.h>
#include <sstream>
#include <sys/stat.h>
#include <unistd.h>

#include "cgroups.h"
#include "errors/errors.h"
#include "logger/logger.h"

using std::stringstream;

// TODO(jsannemo): make this configurable
const string CGROUP_ROOT = "/sys/fs/cgroup";
const string CPU_ACCT = "cpuacct";
const string MEMORY = "memory";
const string PIDS = "pids";
const string BLKIO = "blkio";

string getCgroupName(pid_t pid) {
    stringstream ss;
    ss << "omogen_" << pid;
    return ss.str();
}

string Cgroup::getSubsystemPath(const string& subsystem) {
    return CGROUP_ROOT + "/" + subsystem + "/omogencontain/" + name;
}

void writeTo(const string& path, const string& contents) {
    int fd = open(path.c_str(), O_WRONLY | O_TRUNC);
    if (fd == -1) {
        crashSyscall("open");
    }
    int written = write(fd, contents.c_str(), contents.size());
    if (written == -1) {
        crashSyscall("write");
    }
    if (written != (int)contents.size()) {
        LOG(FATAL) << "Could not write cgroup file" << endl;
        crash();
    }
    if (close(fd) == -1) {
        crashSyscall("close");
    }
}

long long readLongLongFrom(const string& path) {
    int fd = open(path.c_str(), O_RDONLY);
    if (fd == -1) {
        crashSyscall("open");
    }
    char buf[30];
    int readd = read(fd, buf, sizeof(buf));
    if (readd == -1) {
        crashSyscall("read");
    }
    if (readd >= (int)sizeof(buf) - 1) {
        LOG(FATAL) << "Could not read cgroup value" << endl;
        crash();
    }
    if (close(fd) == -1) {
        crashSyscall("close");
    }
    if (readd && buf[readd - 1] == '\n') {
        --readd;
    }
    buf[readd] = 0;
    long long value = strtoll(buf, nullptr, 10);
    return value;
}

void Cgroup::enableSubsystem(const string& subsystem) {
    string subsystemPath = getSubsystemPath(subsystem);
    if (mkdir(subsystemPath.c_str(), 0755) == -1) {
        if (errno != EEXIST) {
            crashSyscall("mkdir");
        }
    }

    stringstream ss;
    ss << pid;
    string toWrite = ss.str();

    string taskFile = subsystemPath + "/tasks";
    writeTo(taskFile, toWrite);
}

void Cgroup::disableSubsystem(const string& subsystem) {
    string subsystemPath = getSubsystemPath(subsystem);
    if (rmdir(subsystemPath.c_str()) == -1) {
        if (errno != ENOENT) {
            crashSyscall("mkdir");
        }
    }
}

long long Cgroup::CpuUsed() {
    string cpuAcctUsagePath = getSubsystemPath(CPU_ACCT) + "/cpuacct.usage";
    long long nanoSeconds = readLongLongFrom(cpuAcctUsagePath);
    return nanoSeconds / 1000000;
}

void Cgroup::SetMemoryLimit(long long memLimitKb) {
    string memoryLimitPath = getSubsystemPath(MEMORY) + "/memory.limit_in_bytes";
    stringstream memoryLimit;
    memoryLimit << memLimitKb * 1000;
    writeTo(memoryLimitPath, memoryLimit.str());
}

long long Cgroup::MemoryUsed() {
    string memoryUsagePath = getSubsystemPath(MEMORY) + "/memory.max_usage_in_bytes";
    long long bytes = readLongLongFrom(memoryUsagePath);
    return bytes / 1000;
}

void Cgroup::SetProcessLimit(int maxProcesses) {
    string processPath = getSubsystemPath(PIDS) + "/pids.max";
    stringstream pidLimit;
    pidLimit << maxProcesses;
    writeTo(processPath, pidLimit.str());
}

long long Cgroup::BytesTransferred() {
    string bytesTransferredPath = getSubsystemPath(BLKIO) + "/blkio.throttle.io_service_bytes";
    // The total bytes used is given as the last integer token in the file, which requires
    // the following somewhat cumbersome reading
    int fd = open(bytesTransferredPath.c_str(), O_RDONLY);
    if (fd == -1) {
        crashSyscall("open");
    }
    char buf[3000];
    int readd = read(fd, buf, sizeof(buf));
    if (readd == -1) {
        crashSyscall("read");
    }
    if (readd >= (int)sizeof(buf) - 1) {
        LOG(FATAL) << "Could not read cgroup value" << endl;
        crash();
    }
    if (close(fd) == -1) {
        crashSyscall("close");
    }
    buf[readd] = 0;
    assert(readd > 0);
    int at = readd - 1;
    bool sawDigit = false;
    while (!sawDigit || isdigit(buf[at])) {
        if (isdigit(buf[at])) {
            sawDigit = true;
        } else {
            buf[at] = 0;
        }
        --at;
    }
    ++at;
    long long bytes = strtoll(buf+at, nullptr, 10) / 1000;
    return bytes;
}

void Cgroup::Reset() {
    string cpuAcctUsagePath = getSubsystemPath(CPU_ACCT) + "/cpuacct.usage";
    writeTo(cpuAcctUsagePath, "0");
    string memoryUsagePath = getSubsystemPath(MEMORY) + "/memory.max_usage_in_bytes";
    writeTo(memoryUsagePath, "0");
    string ioResetPath = getSubsystemPath(BLKIO) + "/blkio.reset_stats";
    writeTo(ioResetPath, "0");
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
