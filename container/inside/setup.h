#pragma once

#include "proto/container.pb.h"

namespace omogenexec {

// Setup the process for executing the request, and run the executable given by
// it.
[[noreturn]] void SetupAndRun(const proto::ContainerExecution& request,
                              int errorFd);

}  // namespace omogenexec
