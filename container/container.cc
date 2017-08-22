#include <atomic>
#include <chrono>
#include <condition_variable>
#include <iostream>
#include <memory>
#include <mutex>
#include <string>
#include <signal.h>
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
using namespace std::chrono_literals;
using std::condition_variable;
using std::endl;
using std::make_unique;
using std::mutex;
using std::string;
using std::thread;
using std::unique_lock;
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
    if (close(errorPipe[1]) == -1) {
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
    int pid = waitpid(initPid, &waitStatus, 0);
    if (pid == -1) {
        if (errno == EINTR) {
            return -1;
        }
        OE_FATAL("waitpid");
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

struct MonitorState {
    Container *container;
    atomic<bool> isDead;
    atomic<bool> shouldKill;
    bool waitReady;
    mutex lock;
    condition_variable waitCv;
    int waitStatus;

    MonitorState(Container* cont) : container(cont), isDead(false), shouldKill(false), waitReady(false), waitStatus(0) {}
};

ExecuteResponse Container::monitorInit(const ResourceLimits& limits) {
    ExecuteResponse response;
    Stopwatch watch;
    MonitorState monitorState(this);

    // We keep one thread that only waits for the process to complete.
    // We also let this thread be responsible for killing the process in case it exceeds
    // its resource limits. This avoids races between killing the process and waiting for it,
    // something that could otherwise result in us killing an unrelated process after the PID
    // has been reused.
    pthread_t waitThread;
    errno = pthread_create(&waitThread, nullptr, [](void* arg) -> void* {
        MonitorState *state = static_cast<MonitorState*>(arg);
        // The resource monitor loop notifies us if we should kill init by giving us SIGALRM to interrupt our
        // wait. We use a lock and flag to tell the monitor when we have set up our own signal handler to avoid
        // getting such a signal before the handler is installed, otherwise we would get killed by the signal.
        {
            unique_lock<std::mutex> waitLock(state->lock);
            struct sigaction action;
            memset(&action, 0, sizeof(action));
            action.sa_handler = [](int){};
            sigaction(SIGALRM, &action, NULL);
            state->waitReady = true;
        }
        state->waitCv.notify_one();
        while (true) {
            if (state->shouldKill) {
                state->container->killInit();
            }
            int waitStatus = state->container->waitInit();
            if (waitStatus != -1) {
                state->waitStatus = waitStatus;
                break;
            }
        }
        state->isDead = true;
        // To avoid some latency, we wake the resource monitor up from its polling sleep
        // whenever the process is dead.
        state->waitCv.notify_one();
        return nullptr;
    }, &monitorState);
    if (errno != 0) {
        OE_FATAL("pthread_create");
    }

    // Wait for the waitThread to set up its signal handler
    {
        unique_lock<std::mutex> waitLock(monitorState.lock);
        monitorState.waitCv.wait(waitLock, [&]{ return monitorState.waitReady; });
    }

    while (!monitorState.isDead) {
#define CHECK_LIM(current, limit, name) \
        if ((current) > (limit)) { \
            OE_LOG(TRACE) << name << " exceeded" << endl; \
            monitorState.shouldKill = true; \
            pthread_kill(waitThread, SIGALRM); \
            break; \
        }
        // Memory does not need to be monitored, since this is the only limit
        // that the control groups can be limit by itself.
        CHECK_LIM(cgroup->CpuUsed(), (long long)(limits.cputime() * 1000), "CPU");
        CHECK_LIM(watch.millis(), (long long)(limits.walltime() * 1000), "Wall time");
        CHECK_LIM(cgroup->DiskIOUsed(), limits.diskio(), "Disk IO");

        unique_lock<std::mutex> timeoutLock(monitorState.lock);
        monitorState.waitCv.wait_for(timeoutLock, 5ms, [&]{ return !monitorState.isDead; });
#undef CHECK_LIM
    }
    if ((errno = pthread_join(waitThread, nullptr)) != 0) {
        OE_FATAL("pthread_join");
    }
    long long elapsed = watch.millis();

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
    // We do not want to do a call to mutable_resrouceexceeded() unless any resource actually was exceeded.
    // This lets client test for the presence of resourceexceeded to determine if a resource was exceeded
    // instead of checking each resource individually.
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
    // We only want to set the normal termination cause in case we did not exceed any resource,
    // since the termination does not have much meaning otherwise (it is essentially random depending on
    // whether the process completed just before getting killed or not.
    if (!termination->has_resourceexceeded()) {
        setTermination(termination, monitorState.waitStatus);
    }
    return response;
}

ExecuteResponse Container::Execute(const ExecuteRequest& request) {
    cgroup->SetMemoryLimit(request.limits().memory());
    cgroup->SetProcessLimit(request.limits().processes());
    if (!request.SerializeToFileDescriptor(commandPipe[1])) {
        OE_LOG(FATAL) << "Could not send request to init" << endl;
        OE_CRASH();
    }
    if (close(commandPipe[1]) == -1) {
        OE_FATAL("close");
    }
    string errMsg;
    char err[1025];
    while (true) {
        int r = read(errorPipe[0], err, sizeof(err) - 1);
        if (r == 0) {
            break;
        }
        if (r == -1) {
            if (errno == EINTR) {
                continue;
            }
            OE_FATAL("read");
        }
        err[r] = 0;
        errMsg += string(err, err + r);
    }
    if (!errMsg.empty()) {
        // To protect aginst errors we don't handle explicitly, we write a 1 byte just before executing
        // the child process. If this is not the first byte, we got a real, handled error.
        if (errMsg[0] != '\1') {
            ExecuteResponse response;
            response.mutable_failure()->set_error(errMsg);
            waitInit();
            return response;
        }
    } else {
        ExecuteResponse response;
        response.mutable_failure()->set_error("Init crashed before execve");
        waitInit();
        return response;
    }
    cgroup->Reset();
    return monitorInit(request.limits());
}

} // namespace omogenexec
