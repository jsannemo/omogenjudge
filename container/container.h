#pragma once

#include <sys/types.h>
#include <sys/wait.h>

#include "cgroups.h"
#include "proto/omogenexec.pb.h"

using std::unique_ptr;

namespace omogenexec {

// A process container, implemented using Linux control groups, namespaces and rlimits.
// Note that the constructor creates the new process ahead-of-time and performs some setup.
// The remaining setup is performed once Execute() is called, minimizing the latency a bit.
class Container {
    // Process ID of the child process we are executing the new program in (called init)
    pid_t initPid;
    // Since we may receive the execution request after starting the new process, we use
    // a pipe to send the request to the process.
    int commandPipe[2];
    // To be able to distinguish between e.g. the program itself exiting and setup failing,
    // we keep an additional pipe that the new process can use to send us errors.
    int errorPipe[2];
    unique_ptr<Cgroup> cgroup;
    // The path to the new root with specific paths mounted to it
	std::string containerRoot;

    // Have the contained process been waited for?
    bool waitedFor;

    ExecuteResponse monitorInit(const ResourceLimits& limits);
    void killInit();
    int waitInit();


public:

    ExecuteResponse Execute(const ExecuteRequest& request);

    Container();
    ~Container();

    Container(const Container&) = default;
    Container& operator=(const Container&) = delete;
};

}
