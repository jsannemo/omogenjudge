#include <atomic>
#include <chrono>
#include <memory>
#include <unistd.h>
#include <string>
#include <thread>
#include <vector>

#include "chroot.h"
#include "container.h"
#include "init.h"
#include "util/error.h"
#include "util/log.h"

using std::atomic;
using std::chrono::duration_cast;
using std::chrono::steady_clock;
using std::chrono::time_point;
using std::string;
using std::thread;
using std::vector;

Container::Container() {
    if (pipe(commandPipe) == -1) {
        CRASH_ERROR("pipe");
    }
    containerRoot = CreateTemporaryRoot();
    InitArgs args { commandPipe[0], containerRoot };
    vector<char> stack(100 * 1024);
    initPid = clone(Init, stack.data() + stack.size(), SIGCHLD | CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS | CLONE_NEWPID | CLONE_NEWUSER | CLONE_NEWUTS, &args);
    if (initPid == -1) {
        CRASH_ERROR("clone");
    }
    LOG(TRACE) << "Created new container with process ID is " << initPid << endl;
    if (close(commandPipe[0]) == -1) {
        CRASH_ERROR("close");
    }
    cgroup = make_unique<Cgroup>(initPid);
}

Container::~Container() {
    DestroyDirectory(containerRoot);
    if (initPid != 0) {
        killInit();
        waitInit();
    }
}

bool isTermination(int waitStatus) {
    return WIFEXITED(waitStatus) || WIFSIGNALED(waitStatus) || WIFSTOPPED(waitStatus);
}

void setTermination(Termination* termination, int waitStatus) {
    if (WIFEXITED(waitStatus)) {
        termination->mutable_exit()->set_code(WEXITSTATUS(waitStatus));
    } else if (WIFSIGNALED(waitStatus)) {
        termination->mutable_signal()->set_signal(WTERMSIG(waitStatus));
    } else if (WIFSTOPPED(waitStatus)) {
        termination->mutable_signal()->set_signal(WSTOPSIG(waitStatus));
    } else assert(false && "Invalid exit status");
}

void Container::killInit() {
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
            CRASH_ERROR("waitpid");
        }
        break;
    }
    // In case the container is actually used and thus waited for here, we do not want to
    // kill and wait for the process again when we destroy the container, so we mark that
    // init is done by zeroing the pid.
    initPid = 0;
    return waitStatus;
}

ExecuteResponse Container::monitorInit(const ResourceLimits& limits) {
    ExecuteResponse response;
    ExecutionResult *result = response.mutable_result();
    // Note that mutable_terminaton() creates a new pointer on the first call,
    // so we call it here to avoid a race condition between the monitor and the 
    // normal termination setting.
    Termination *termination = result->mutable_termination();
    auto startTime = steady_clock::now();
    atomic<bool> killedByMonitor(false), isDead(false);
    thread resourceMonitor([&](){
        while (!isDead) {
            // Memory does not need to be monitored, since this is the only limit
            // that the control groups can be limit by itself.
            long long cpuUsedMillis = cgroup->CpuUsed();
            if (cpuUsedMillis > (long long)(limits.cputime() * 1000)) {
                LOG(TRACE) << "CPU exceeded" << endl;
                termination->mutable_resourceexceeded()->set_cputime(true);
                killedByMonitor = true;
                killInit();
                break;
            }

            long long elapsed = duration_cast<std::chrono::milliseconds>(steady_clock::now() - startTime).count();
            if (elapsed > (long long)(limits.walltime() * 1000)) {
                LOG(TRACE) << "Wall time exceeded" << endl;
                termination->mutable_resourceexceeded()->set_walltime(true);
                killedByMonitor = true;
                killInit();
                break;
            }
            long long bytesTransferred = cgroup->BytesTransferred();
            if (bytesTransferred > (long long)(limits.diskio() * 1000)) {
                LOG(TRACE) << "Disk usage exceeded" << endl;
                termination->mutable_resourceexceeded()->set_diskio(true);
                killedByMonitor = true;
                killInit();
                break;
            }

            const timespec pollSleep { 0, 5000000 }; // 5 milliseconds
            nanosleep(&pollSleep, nullptr);
        }
    });
    int waitStatus = waitInit();
    long long elapsed = duration_cast<std::chrono::milliseconds>(steady_clock::now() - startTime).count();
    isDead = true;
    resourceMonitor.join();
    assert(isTermination(waitStatus));
    setTermination(result->mutable_termination(), waitStatus);
    // It may happen that the resource monitor writes its termination cause,
    // but init dies before the monitor kills it. In this case, we may
    // write our termination cause before killedByMonitor is set.
    // Therefore, we clear the termination explicitly instead of e.g. 
    // if (!killedByMonitor) setTermination(...)
    if (killedByMonitor) {
        result->mutable_termination()->clear_signal();
        result->mutable_termination()->clear_exit();
    }
    result->mutable_resourceusage()->set_cputime(cgroup->CpuUsed());
    result->mutable_resourceusage()->set_walltime(elapsed);
    result->mutable_resourceusage()->set_memory(cgroup->MemoryUsed());
    result->mutable_resourceusage()->set_diskio(cgroup->BytesTransferred());
    LOG(TRACE) << "Responding with " << response.DebugString() << endl;
    return response;
}

ExecuteResponse Container::Execute(const ExecuteRequest& request) {
    cgroup->SetMemoryLimit(request.limits().memory());
    cgroup->SetProcessLimit(10);
    cgroup->Reset();
    if (!request.SerializeToFileDescriptor(commandPipe[1])) {
        LOG(FATAL) << "Could not send request to init" << endl;
        CRASH();
    }
    if (close(commandPipe[1]) == -1) {
        CRASH_ERROR("close");
    }
    return monitorInit(request.limits());
}
