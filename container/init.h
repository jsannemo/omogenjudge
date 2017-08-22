#pragma once

#include <string>

namespace omogenexec {

class InitException : public std::runtime_error {
     std::string msg;

public:
     InitException(const std::string& msg) : runtime_error("Container failed to setup execution"), msg(msg) {}
     const char* what() const noexcept override {
         return msg.c_str();
     }
};

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
