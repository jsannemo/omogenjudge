#include "sandbox/server/executor.h"

#include <mutex>

#include "sandbox/container/outside/container.h"
#include "sandbox/server/request_validator.h"

using grpc::Status;
using grpc::StatusCode;

namespace omogen {
namespace sandbox {

ExecuteServiceImpl::ExecuteServiceImpl() {
  cleanup_thread = new std::thread(&ExecuteServiceImpl::Cleanup, this);
}

#define PUSH_CONTAINER_TO_CLEANUP(container)            \
  if ((container) != nullptr) {                         \
    absl::MutexLock containerLock(&delete_queue_mutex); \
    delete_queue.push_back(std::move(container));       \
    (container) = nullptr;                              \
  }

Status ExecuteServiceImpl::Execute(
    ServerContext* context,
    ServerReaderWriter<ExecuteResponse, ExecuteRequest>* stream) {
  LOG(INFO) << "/ExecService.Execute: start";
  unique_ptr<Container> container;
  ExecuteRequest request;
  ContainerSpec last_spec;
  bool first_request = false;
  while (stream->Read(&request)) {
    LOG(INFO) << "/ExecService.Execute: new request";
    Status status = ValidateExecuteRequest(request, last_spec);
    if (!status.ok()) {
      return status;
    }
    // We always need a container spec on the first request.
    if (request.has_container_spec()) {
      last_spec = request.container_spec();
    } else if (first_request) {
      return grpc::Status(
          grpc::StatusCode::INVALID_ARGUMENT,
          "Can't send an execution before sending a container spec");
    }
    first_request = false;
    // Reuse the old container unless it is considered invalid or has been
    // killed. It is only ever killed by calls to Execute, so it is not a race
    // to check this here before the Execute call.
    bool needs_new_container = (container == nullptr) || container->IsDead() ||
                               request.has_container_spec();
    if (needs_new_container) {
      PUSH_CONTAINER_TO_CLEANUP(container);
      container = std::make_unique<Container>(ContainerIds::GetId(), last_spec);
    }
    StatusOr<Termination> result = container->Execute(request.execution());
    if (result.ok()) {
      ExecuteResponse response;
      *response.mutable_termination() = result.value();
      stream->Write(response);
    } else {
      PUSH_CONTAINER_TO_CLEANUP(container);
      return result.status();
    }
  }
  PUSH_CONTAINER_TO_CLEANUP(container);
  LOG(INFO) << "/ExecService.Execute: finish";
  return Status::OK;
}

void ExecuteServiceImpl::Cleanup() {
  while (true) {
    // Swap to avoid blocking the cleanup queue.
    // This cleans up the containers by destroying the unique pointers
    // when this containing vector is destroyed at the end of the scope.
    std::vector<unique_ptr<Container>> to_delete;
    {
      absl::MutexLock queueLock(&delete_queue_mutex);
      to_delete.swap(delete_queue);
    }
    if (to_delete.empty()) {
      usleep(100 * 1000);  // 100 milliseconds in microseconds
    } else {
      VLOG(2) << "Cleanup " << to_delete.size();
    }
  }
}

}  // namespace sandbox
}  // namespace omogen
