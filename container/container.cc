#include <atomic>
#include <chrono>
#include <iostream>
#include <memory>
#include <string>
#include <thread>
#include <unistd.h>
#include <vector>

#include "chroot.h"
#include "container.h"
#include "init.h"
#include "util/error.h"
#include "util/files.h"
#include "util/log.h"
#include "util/time.h"

using std::atomic;
using std::chrono::steady_clock;
using std::endl;
using std::make_unique;
using std::string;
using std::thread;
using std::vector;

namespace omogenexec {

const int CHILD_STACK_SIZE = 100 * 1000; // 100 KB

Container::Container() {
    if (pipe(commandPipe) == -1) {
        OE_FATAL("pipe");
    }
    if (pipe(errorPipe) == -1) {
        OE_FATAL("pipe");
    }
    // The container will get a new root with chroot; we store this in a temporary directory
    containerRoot = MakeTempDir();
    // Clone requires us to provide a new stack for the child process
    vector<char> stack(CHILD_STACK_SIZE);
    InitArgs args { commandPipe[0], errorPipe[1], containerRoot };
    // Clone and create new namespaces for the contained process
    initPid = clone(Init, stack.data() + stack.size(), SIGCHLD | CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS | CLONE_NEWPID | CLONE_NEWUSER | CLONE_NEWUTS, &args);
    if (initPid == -1) {
        OE_FATAL("clone");
    }
    OE_LOG(TRACE) << "Created new container with process ID is " << initPid << endl;
    if (close(commandPipe[0]) == -1) {
        OE_FATAL("close");
    }
    cgroup = make_unique<Cgroup>(initPid);
}

Container::~Container() {
    RemoveTree(containerRoot);
    if (!waitedFor) {
        killInit();
        waitInit();
    }
}

void Container::killInit() {
    // Since we immediately move the contained process out of our process group,
    // it is fine to do kill(-initPid)
    kill(-initPid, SIGKILL);
    kill(initPid, SIGKILL);
}

int Container::waitInit() {
    int waitStatus = 0;
    while (true) {
        int pid = waitpid(initPid, &waitStatus, 0);
        if (pid == -1) {
            if (errno == EINTR) {
                continue;
            }
            OE_FATAL("waitpid");
        }
        break;
    }
    // In case the container is actually used and thus waited for here, we do not want to
    // kill and wait for the process again when we destroy the container, so we mark that
    // we have already waited for the process.
    waitedFor = true;
    return waitStatus;
}

static void setTermination(Termination* termination, int waitStatus) {
    if (WIFEXITED(waitStatus)) {
        termination->mutable_exit()->set_code(WEXITSTATUS(waitStatus));
    } else if (WIFSIGNALED(waitStatus)) {
        termination->mutable_signal()->set_signal(WTERMSIG(waitStatus));
    } else if (WIFSTOPPED(waitStatus)) {
        termination->mutable_signal()->set_signal(WSTOPSIG(waitStatus));
    } else assert(false && "Invalid exit status");
}


ExecuteResponse Container::monitorInit(const ResourceLimits& limits) {
    ExecuteResponse response;
    Stopwatch watch;
    atomic<bool> isDead(false);
    // We monitor whether init has exceeded any of the resources we cannot limit explicitly
    // in a separate thread, since we need also need to wait for the process in case it 
    // exits normally.
    thread resourceMonitor([&](){
#define CHECK_LIM(current, limit, name) \
    if ((current) > (limit)) { \
        OE_LOG(TRACE) << name << " exceeded" << endl; \
        killInit(); \
        break; \
    }
        while (!isDead) {
            // Memory does not need to be monitored, since this is the only limit
            // that the control groups can be limit by itself.
            CHECK_LIM(cgroup->CpuUsed(), (long long)(limits.cputime() * 1000), "CPU");
            CHECK_LIM(watch.millis(), (long long)(limits.walltime() * 1000), "Wall time");
            CHECK_LIM(cgroup->DiskIOUsed(), limits.diskio(), "Disk IO");
            const timespec pollSleep { 0, 5000000 }; // 5 milliseconds
            nanosleep(&pollSleep, nullptr);
        }
#undef CHECK_LIM
    });
    int waitStatus = waitInit();
    long long elapsed = watch.millis();
    isDead = true;
    resourceMonitor.join();

    ExecutionResult *result = response.mutable_result();
    ResourceUsage *resourceUsage = result->mutable_resourceusage();
    long long cpuUsedMs = cgroup->CpuUsed();
    long long memoryUsedKb = cgroup->MemoryUsed();
    long long diskIoKb = cgroup->DiskIOUsed();
    resourceUsage->set_cputime(cpuUsedMs);
    resourceUsage->set_walltime(elapsed);
    resourceUsage->set_memory(memoryUsedKb);
    resourceUsage->set_diskio(diskIoKb);

    Termination *termination = result->mutable_termination();
    if (cpuUsedMs > (long long)(limits.cputime() * 1000)) {
        termination->mutable_resourceexceeded()->set_cputime(true);
    }
    if (elapsed > (long long)(limits.walltime() * 1000)) {
        termination->mutable_resourceexceeded()->set_walltime(true);
    }
    if (memoryUsedKb > (long long)(limits.memory())) {
        termination->mutable_resourceexceeded()->set_memory(true);
    }
    if (diskIoKb > (long long)(limits.diskio())) {
        termination->mutable_resourceexceeded()->set_diskio(true);
    }
    // We only set the normal termination cause in case we did not exceed any resource,
    // since 
    if (!termination->has_resourceexceeded()) {
        setTermination(result->mutable_termination(), waitStatus);
    }
    return response;
}

ExecuteResponse Container::Execute(const ExecuteRequest& request) {
    cgroup->SetMemoryLimit(request.limits().memory());
    cgroup->SetProcessLimit(request.limits().processes());
    cgroup->Reset();
    if (!request.SerializeToFileDescriptor(commandPipe[1])) {
        OE_LOG(FATAL) << "Could not send request to init" << endl;
        OE_CRASH();
    }
    if (close(commandPipe[1]) == -1) {
        OE_FATAL("close");
    }
    return monitorInit(request.limits());
}

} // namespace omogenexec
