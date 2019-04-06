#pragma once

#include "exec/proto/container.pb.h"

namespace omogen {
namespace exec {

// Setup the process for executing the request, and run the executable given by
// it.
[[noreturn]] void SetupAndRun(const proto::ContainerExecution& request,
                              int errorFd);

}  // namespace exec
}  // namespace omogen
