#pragma once

#include "sandbox/proto/container.pb.h"

namespace omogen {
namespace sandbox {

// Setup the process for executing the request, and run the executable given by
// it.
[[noreturn]] void SetupAndRun(const proto::ContainerExecution& request,
                              int errorFd);

}  // namespace sandbox
}  // namespace omogen
