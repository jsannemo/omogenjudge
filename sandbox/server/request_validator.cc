#include "sandbox/server/request_validator.h"

#include "sandbox/api/execute_service.grpc.pb.h"

using grpc::Status;

namespace omogen {
namespace sandbox {

Status ValidateExecuteRequest(const ExecuteRequest& request,
                              const ContainerSpec& previous_spec) {
  // TODO(jsannemo): implement
  // Verify command:
  // - exists
  // - is executable
  // - has (group-exec and group omogenjudge-clients) or others-exec
  // - is in one of the readable directories.
  // - all directories up the readable directory is (group-exec/read and group
  // omogenjudge-clients) or others-exec
  //
  // Verify that the working directory:
  // - exists
  // - is readable
  // - all directories up to it is (group-exec/read and group
  // omogenjudge-clients) or others-exec
  //
  // Verify that there are at least three stream mappings
  //
  // Verify that the input mapping
  // - exists
  // - is readable
  // - it and all directories up are readable by omogenjudge-clients
  //
  // Verify that the output/error mapping files
  // - is to writable directories
  // - is in a directory writable to by omogenjudge-clients
  // - is in a directory with read/exec access by omogenjudge-clients
  // - do not already exist
  return Status();
}

}  // namespace sandbox
}  // namespace omogen
