#ifndef SANDBOX_SERVER_REQUEST_VALIDATOR_H
#define SANDBOX_SERVER_REQUEST_VALIDATOR_H

#include "grpc++/grpc++.h"
#include "sandbox/api/execute_service.grpc.pb.h"

namespace omogen {
namespace sandbox {

grpc::Status ValidateExecuteRequest(const ExecuteRequest& request,
                                    const ContainerSpec& previous_spec);

}  // namespace sandbox
}  // namespace omogen
#endif
