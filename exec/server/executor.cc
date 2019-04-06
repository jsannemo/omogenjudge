#include "exec/server/executor.h"

#include <mutex>

#include "exec/container/outside/container.h"
#include "grpcpp/grpcpp.h"

using grpc::Status;
using grpc::StatusCode;

namespace omogen {
namespace exec {

ExecuteServiceImpl::ExecuteServiceImpl() {
  cleanupThread = new std::thread(&ExecuteServiceImpl::cleanup, this);
}

#define PUSH_CONTAINER_TO_CLEANUP(container)          \
  if ((container) != nullptr) {                       \
    absl::MutexLock containerLock(&deleteQueueMutex); \
    deleteQueue.push_back(std::move(container));      \
    (container) = nullptr;                            \
  }

Status ExecuteServiceImpl::Execute(
    ServerContext* context,
    ServerReaderWriter<ExecuteResponse, ExecuteRequest>* stream) {
  LOG(INFO) << "New request";
  unique_ptr<Container> container;
  ExecuteRequest request;
  ContainerSpec lastSpec;
  bool firstRequest = false;
  while (stream->Read(&request)) {
    if (request.has_container_spec()) {
      lastSpec = request.container_spec();
    } else if (firstRequest) {
      return grpc::Status(
          grpc::StatusCode::FAILED_PRECONDITION,
          "Can't send an execution before sending a container spec");
    }
    firstRequest = false;
    bool needsNewContainer = (container == nullptr) || container->IsDead() ||
                             request.has_container_spec();
    if (needsNewContainer) {
      PUSH_CONTAINER_TO_CLEANUP(container);
      container = std::make_unique<Container>(ContainerIds::GetId(), lastSpec);
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
  LOG(INFO) << "Finished request";
  return Status::OK;
}

void ExecuteServiceImpl::cleanup() {
  while (true) {
    std::vector<unique_ptr<Container>> toDelete_;
    {
      absl::MutexLock queueLock(&deleteQueueMutex);
      toDelete_.swap(deleteQueue);
    }
    VLOG(2) << "Cleanup " << toDelete_.size();
    if (toDelete_.empty()) {
      usleep(1000 * 1000);  // 1 second in microseconds
    }
  }
}

}  // namespace exec
}  // namespace omogen
