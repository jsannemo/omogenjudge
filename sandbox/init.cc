#include <fcntl.h>
#include <dirent.h>
#include <sys/prctl.h>
#include <sys/resource.h>
#include <type_traits>
#include <unistd.h>

#include "chroot.h"
#include "errors/errors.h"
#include "init.h"
#include "logger/logger.h"
#include "proto/omogenexec.pb.h"

using std::remove_const;

void closeFdsExcept(vector<int> fdsToKeep) {
    DIR *fdDir = opendir("/proc/self/fd");
    if (fdDir == nullptr) {
        crashSyscall("opendir");
    }
    // Do not accidentally close the fd directory fd
    fdsToKeep.push_back(dirfd(fdDir));
    while (true) {
        struct dirent *entry = readdir(fdDir);
        if (entry == nullptr) {
            if (errno != 0) {
                crashSyscall("readdir");
            }
            break;
        }

        errno = 0;
        int fd = strtol(entry->d_name, nullptr, 10);
        if (errno != 0) {
            LOG(WARN) << "Ignoring invalid fd entry: " << entry->d_name << endl;
        } else {
            for (int fdToKeep : fdsToKeep) {
                if (fd == fdToKeep) {
                    goto skip;
                }
            }
            if (close(fd) == -1) {
                crashSyscall("close");
            }
skip:;
        }
    }
    if (closedir(fdDir) == -1) {
        crashSyscall("closedir");
    }
}

void setResourceLimit(int resource, rlim_t limit) {
    rlimit rlim = { .rlim_cur = limit, .rlim_max = limit };
    if (setrlimit(resource, &rlim) == -1) {
        crashSyscall("setrlimit");
    }
}

void setResourceLimits(const ResourceLimits& resourceLimits) {
    setResourceLimit(RLIMIT_AS, (rlim_t)resourceLimits.memory() * 1000);
    setResourceLimit(RLIMIT_STACK, RLIM_INFINITY);
    setResourceLimit(RLIMIT_MEMLOCK, 0);
    setResourceLimit(RLIMIT_CORE, 0);
    setResourceLimit(RLIMIT_NOFILE, 100);
}

void setupStreams(const StreamRedirections& streams) {
    // File descriptors are reused sequentially when we close our streams,
    // so the newly opened files will be mapped to the correct stream file descriptor
    close(0);
    if (open(streams.infile().c_str(), O_RDONLY) == -1) {
        crashSyscall("open");
    }
    close(1);
    if (open(streams.outfile().c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666) == -1) {
        crashSyscall("open");
    }
    close(2);
    if (open(streams.errfile().c_str(), O_WRONLY | O_CREAT | O_TRUNC, 0666) == -1) {
        crashSyscall("open");
    }
}

char* strdupOrDie(const char* str) {
    char *ret = strdup(str);
    if (ret == nullptr) {
        crashSyscall("strdup");
    }
    return ret;
}

// TODO(jsannemo): possible expand this/make it configurable if needed?
char** setupEnvironment() {
    char **env = static_cast<char**>(malloc(2 * sizeof(char*)));
    if (env == nullptr) {
        crashSyscall("malloc");
    }
    env[0] = strdup("PATH=/bin:/usr/bin");
    env[1] = nullptr;
    return env;
}

int Init(void* argp) {
    // If the parent dies for some reason, we wish to be SIGKILLed
    if (prctl(PR_SET_PDEATHSIG, SIGKILL) < 0) {
        crashSyscall("prctl");
    }
    InitArgs args = *static_cast<InitArgs*>(argp);

    // We close all file descriptors to prevent leaks from the parent
    closeFdsExcept(vector<int> {0, 1, 2, args.commandPipe});

    // We move ourself to a new process group so that we can kill(-pid) without
    // accidentally killing the parent
    if (setpgrp() == -1) {
        crashSyscall("setpgrp");
    }

    Chroot chroot(args.containerRoot);

    ExecuteRequest request;
    if (!request.ParseFromFileDescriptor(args.commandPipe)) {
        LOG(FATAL) << "Could not read request from container" << endl;
        crash();
    }
    LOG(TRACE) << "Init got request " << request.DebugString() << endl;
    if (close(args.commandPipe) == -1) {
        crashSyscall("close");
    }

    for (const auto& rule : request.directories()) {
        chroot.AddDirectoryRule(rule);
    }
    chroot.SetRoot();
    setResourceLimits(request.limits());

    const Command& command = request.command();
    char *argv[command.flags_size() + 2];
    argv[0] = strdupOrDie(request.command().command().c_str());
    for (int i = 0; i < command.flags_size(); ++i) {
        argv[i + 1] = strdupOrDie(command.flags(i).c_str());
    }
    argv[command.flags_size() + 1] = nullptr;
    // We set up the stream redirections last of all so that we can read error logs
    // as long as possible
    setupStreams(request.streams());
    if (execve(argv[0], argv, setupEnvironment()) == -1) {
        crashSyscall("execve");
    }
    return 1;
}
