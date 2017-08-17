#include <sys/wait.h>

#include "cgroups.h"
#include "proto/omogenexec.pb.h"

using std::unique_ptr;

// A process container, implemented using Linux control groups, namespaces and rlimits.
// When initialized, a new process is created that performs some initial setup.
// The remaining setup is performed once Execute() is called. This makes it possible to
// create containers before they are to be used as an optimization.
// Upon Execute(), the various resource measurements are reset, so it is fine to let
// the container be long-lived.
class Container {
    // Process ID of the child process we are executing the new program in (called init)
    pid_t initPid;
    // Since we may receive the execution request after starting the new process, we use
    // a pipe to send the request to the process.
    int commandPipe[2];
    unique_ptr<Cgroup> cgroup;
    // The path to the new root with specific paths mounted to it
    string containerRoot;

    ExecuteResponse monitorInit(const ResourceLimits& limits);
    void killInit();
    int waitInit();

public:
    Container();
    ~Container();

    ExecuteResponse Execute(const ExecuteRequest& request);

};
