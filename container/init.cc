#include <fcntl.h>
#include <iostream>
#include <sys/prctl.h>
#include <sys/resource.h>
#include <type_traits>
#include <unistd.h>
#include <vector>

#include "chroot.h"
#include "init.h"
#include "util/error.h"
#include "util/files.h"
#include "util/log.h"
#include "proto/omogenexec.pb.h"

using std::endl;
using std::string;
using std::vector;

namespace omogenexec {

static void setResourceLimit(int resource, rlim_t limit) {
    rlimit rlim = { .rlim_cur = limit, .rlim_max = limit };
    if (setrlimit(resource, &rlim) == -1) {
        throw InitException("setrlimit: " + StrError());
    }
}

static void setResourceLimits(const ResourceLimits& resourceLimits) {
    setResourceLimit(RLIMIT_AS, (rlim_t)resourceLimits.memory() * 1000);
    setResourceLimit(RLIMIT_STACK, RLIM_INFINITY);
    setResourceLimit(RLIMIT_MEMLOCK, 0);
    setResourceLimit(RLIMIT_CORE, 0);
    setResourceLimit(RLIMIT_NOFILE, 100);
}

static void openFileWithFd(int wantFd, const string& path, bool writable) {
    if (close(wantFd) == -1) {
        throw InitException("close: " + StrError());
    }
    int fd = writable ?
        open(path.c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666)
        : open(path.c_str(), O_RDONLY);
    if (fd == -1) {
        throw InitException("open: " + StrError());
    }
    if (fd != wantFd) {
        throw InitException("Got the wrong fd for stream");
    }
}

static void setupStreams(const StreamRedirections& streams) {
    // When opening a new file, the lowest unused file descriptor is reused. Thus, we can
    // map descriptors 0/1/2 to a particular file by closing the descriptor and then opening
    // the correct file, since they will be open when the process starts.
    openFileWithFd(0, streams.infile().c_str(), false);
    openFileWithFd(1, streams.outfile().c_str(), true);
    openFileWithFd(2, streams.errfile().c_str(), true);
}

static vector<const char*> setupEnvironment() {
    vector<const char*> env;
    // Path is needed for e.g. gcc, which searchs for some binaries in the path
    env.push_back("PATH=/bin:/usr/bin");
    env.push_back(nullptr);
    return env;
}

[[noreturn]] int Init(void* argp) {
    InitArgs args = *static_cast<InitArgs*>(argp);
    try {
        // If the parent dies for some reason, we wish to be SIGKILLed. This is not
        // a race with the parent's death, since a parental death will cause the
        // request.ParseFromFileDescriptor below to fail.
        if (prctl(PR_SET_PDEATHSIG, SIGKILL) < 0) {
            throw InitException("prctl: " + StrError());
        }

        // We close all file descriptors to prevent leaks from the parent
        CloseFdsExcept(vector<int> {0, 1, 2, args.commandPipe, args.errorPipe});

        // We move ourself to a new process group so that we can kill(-pid) without
        // accidentally killing the parent
        if (setpgrp() == -1) {
            throw InitException("setpgrp: " + StrError());
        }

        Chroot chroot(args.containerRoot);
        chroot.SetWD();

        ExecuteRequest request;
        if (!request.ParseFromFileDescriptor(args.commandPipe)) {
            throw InitException("Failed to read execution request from container");
        }
        if (close(args.commandPipe) == -1) {
            throw InitException("close: " + StrError());
        }

        for (const auto& rule : request.directories()) {
            if (!DirectoryExists(rule.oldpath())) {
                throw InitException("Outside directory " + rule.oldpath() + " did not exist!");
            }
            chroot.AddDirectoryRule(rule);
        }
        chroot.SetRoot();
        setResourceLimits(request.limits());

        const Command& command = request.command();
        vector<const char*> argv;
        argv.push_back(request.command().command().c_str());
        for (int i = 0; i < command.flags_size(); ++i) {
            argv.push_back(command.flags(i).c_str());
        }
        argv.push_back(nullptr);
        // We set up the stream redirections last of all so that we can read error logs
        // as long as possible
        setupStreams(request.streams());
        vector<const char*> env = setupEnvironment();
        if (!FileIsExecutable(argv[0])) {
            throw InitException("Command is not an executable file inside the sandbox");
        }
        WriteToFd(args.errorPipe, "\1");
        close(args.errorPipe);
        execve(argv[0], const_cast<char**>(argv.data()), const_cast<char**>(env.data()));
    } catch (InitException e) {
        OE_LOG(INFO) << "Caught exception" << endl;
        WriteToFd(args.errorPipe, e.what());
        close(args.errorPipe);
    }
    abort();
}

}
