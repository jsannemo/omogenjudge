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
        OE_FATAL("setrlimit");
    }
}

static void setResourceLimits(const ResourceLimits& resourceLimits) {
    setResourceLimit(RLIMIT_AS, (rlim_t)resourceLimits.memory() * 1000);
    setResourceLimit(RLIMIT_STACK, RLIM_INFINITY);
    setResourceLimit(RLIMIT_MEMLOCK, 0);
    setResourceLimit(RLIMIT_CORE, 0);
    setResourceLimit(RLIMIT_NOFILE, 100);
}

static void setupStreams(const StreamRedirections& streams) {
    // When opening a new file, the lowest unused file descriptor is reused. Thus, we can
    // map descriptors 0/1/2 to a particular file by closing the descriptor and then opening
    // the correct file.
    int fd;
    close(0);
    fd = open(streams.infile().c_str(), O_RDONLY);
    if (fd == -1) {
        OE_FATAL("open");
    }
    assert(fd == 0);

    close(1);
    fd = open(streams.outfile().c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666);
    if (fd == -1) {
        OE_FATAL("open");
    }
    assert(fd == 1);

    close(2);
    fd = open(streams.errfile().c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666);
    if (fd == -1) {
        OE_FATAL("open");
    }
    assert(fd == 2);
}

// TODO(jsannemo): possibly expand this/make it configurable if needed?
static vector<const char*> setupEnvironment() {
    vector<const char*> env;
    // Path is needed for e.g. gcc, which searchs for some binaries in the path
    env.push_back("PATH=/bin:/usr/bin");
    env.push_back(nullptr);
    return env;
}

[[noreturn]] int Init(void* argp) {
    // If the parent dies for some reason, we wish to be SIGKILLed.
    if (prctl(PR_SET_PDEATHSIG, SIGKILL) < 0) {
        OE_FATAL("prctl");
    }
    InitArgs args = *static_cast<InitArgs*>(argp);

    // We close all file descriptors to prevent leaks from the parent
    CloseFdsExcept(vector<int> {0, 1, 2, args.commandPipe, args.errorPipe});

    // We move ourself to a new process group so that we can kill(-pid) without
    // accidentally killing the parent
    if (setpgrp() == -1) {
        OE_FATAL("setpgrp");
    }

    Chroot chroot(args.containerRoot);
    chroot.SetWD();

    ExecuteRequest request;
    if (!request.ParseFromFileDescriptor(args.commandPipe)) {
        OE_LOG(FATAL) << "Could not read request from container" << endl;
        OE_CRASH();
    }
    if (close(args.commandPipe) == -1) {
        OE_FATAL("close");
    }

    for (const auto& rule : request.directories()) {
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
    execve(argv[0], const_cast<char**>(argv.data()), const_cast<char**>(env.data()));
    abort();
}

}
