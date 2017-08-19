#include <chrono>
#include <unistd.h>

#include "errors/errors.h"
#include "logger/log.h"
#include "sandbox.h"
#include "supervisor.h"
#include "util.h"

Container::Container() {
    if (pipe(commandPipe) == -1) {
        crashSyscall("pipe");
    }
    if (pipe(responsePipe) == -1) {
        crashSyscall("pipe");
    }
    pid = fork();
    if (pid == -1) {
        crashSyscall("fork");
    }
    if (pid == 0) {
        SupervisorMain(commandPipe, responsePipe);
    } else {
        if (close(commandPipe[0]) == -1) {
            crashSyscall("close");
        }
        if (close(responsePipe[1]) == -1) {
            crashSyscall("close");
        }
    }
}

Container::~Container() {
    killSupervisor();
}

void Container::killSupervisor() {
    if (kill(-pid, SIGKILL) == -1) {
        if (errno != ESRCH) {
            LOG(WARN) << "Could not destroy container: " << strerror(errno) << endl;
        }
    }
}


ExecuteResponse Container::Execute(const ExecuteRequest& request) {
    chrono::steady_clock::time_point start = chrono::steady_clock::now();
    if (!request.SerializeToFileDescriptor(commandPipe[1])) {
        LOG(ERROR) << "Could not send command to container" << endl;
        killSupervisor();
        return MakeExecutionFailure("Could not send command to container");
    }
    if (close(commandPipe[1]) == -1) {
        LOG(WARN) << "Parent could not close write command pipe: " << strerror(errno) << endl;
    }
    ExecuteResponse response;
    if (!response.ParsePartialFromFileDescriptor(responsePipe[0])) {
        LOG(ERROR) << "Could not read execution result from container" << endl;
        killSupervisor();
        return MakeExecutionFailure("Could not read execution result from container");
    }
    if (close(responsePipe[0]) == -1) {
        LOG(WARN) << "Parent could not close read result pipe: " << strerror(errno) << endl;
    }
    chrono::steady_clock::time_point finish = chrono::steady_clock::now();
    long long elapsedMillis = chrono::duration_cast<chrono::milliseconds>(finish - start).count();
    LOG(TRACE) << "Execution took " << elapsedMillis << " millis" << endl;
    return response;
}
