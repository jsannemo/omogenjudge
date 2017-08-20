#pragma once

#include <string>

namespace omogenexec {

struct InitArgs {
    // File descriptor to the read end of the command pipe
    int commandPipe;
    // File descriptor to the write end of the error pipe
    int errorPipe;
    // The path where init should build its new rootfs.
    std::string containerRoot;
};

// The entry point of the init process. The args pointer is to an InitArgs struct.
int Init(void* args);

}
