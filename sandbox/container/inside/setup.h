#ifndef SANDBOX_CONTAINER_INSIDE_SETUP_H
#define SANDBOX_CONTAINER_INSIDE_SETUP_H

#include "sandbox/proto/container.pb.h"

namespace omogen {
namespace sandbox {

// Setup the process for executing the request, and run the executable given by
// it. A writable file descriptor must be provided for writing errors.
// After setup finishes, the program will write a \1 byte to the file to signify
// a successfull setup. Otherwise, an error message will be printed to it.
[[noreturn]] void SetupAndRun(const proto::ContainerExecution& request,
                              int error_fd);

}  // namespace sandbox
}  // namespace omogen
#endif
