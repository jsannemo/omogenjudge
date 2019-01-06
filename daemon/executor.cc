#include <mutex>

#include "container/outside/container.h"
#include "daemon/executor.h"
#include "grpcpp/grpcpp.h"

using grpc::Status;
using grpc::StatusCode;

namespace omogenexec {

ExecuteServiceImpl::ExecuteServiceImpl() {
  cleanupThread = new std::thread(&ExecuteServiceImpl::cleanup, this);
}

#define PUSH_CONTAINER_TO_CLEANUP(container)          \
  if ((container) != nullptr) {                       \
    absl::MutexLock containerLock(&deleteQueueMutex); \
    deleteQueue.push_back((container));               \
  }

Status ExecuteServiceImpl::Execute(
    ServerContext* context,
    ServerReaderWriter<api::ExecuteResponse, api::ExecuteRequest>* stream) {
  LOG(INFO) << "New request";
  Container* container = nullptr;
  api::ExecuteRequest request;
  while (stream->Read(&request)) {
    // Grab a new container if we have a container spec in the request.
    if (request.has_container_spec()) {
      PUSH_CONTAINER_TO_CLEANUP(container);
      container = new Container(request.container_spec());
    }
    if (container == nullptr) {
      return grpc::Status(
          grpc::StatusCode::FAILED_PRECONDITION,
          "Can't send an execution before sending a container spec");
    }
    StatusOr<api::Termination> result = container->Execute(request.execution());
    if (result.ok()) {
      api::ExecuteResponse response;
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
    std::vector<Container*> toDelete_;
    {
      absl::MutexLock queueLock(&deleteQueueMutex);
      toDelete_ = deleteQueue;
      deleteQueue.clear();
    }
    VLOG(2) << "Cleanup " << toDelete_.size();
    for (Container* container : toDelete_) {
      delete container;
    }
    if (toDelete_.empty()) {
      usleep(1000 * 1000);  // 50 milliseconds in microseconds
    }
  }
}

}  // namespace omogenexec
