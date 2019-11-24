#ifndef SANDBOX_SERVER_EXECUTOR_H
#define SANDBOX_SERVER_EXECUTOR_H
#include <thread>

#include "absl/base/thread_annotations.h"
#include "absl/synchronization/mutex.h"
#include "grpc/grpc.h"
#include "sandbox/api/execute_service.grpc.pb.h"
#include "sandbox/container/outside/container.h"

using grpc::ServerContext;
using grpc::ServerReaderWriter;
using grpc::Status;

namespace omogen {
namespace sandbox {

// Implementation of the ExecuteService. This server is stateful - it keeps a
// list of use containers that have not yet been cleaned up.
class ExecuteServiceImpl final : public ExecuteService::Service {
  // The cleanup thread for this service, continuously cleaning up old
  // containers.
  std::thread* cleanup_thread;

  absl::Mutex delete_queue_mutex;

  // A vector with containers that should be cleaned up.
  std::vector<unique_ptr<Container>> delete_queue
      GUARDED_BY(delete_queue_mutex);

  // Runs the cleanup loop the remove old containers.
  void Cleanup();

 public:
  // Handler for the Execute requests.
  Status Execute(ServerContext* context,
                 ServerReaderWriter<ExecuteResponse, ExecuteRequest>* stream);

  ExecuteServiceImpl();
};

}  // namespace sandbox
}  // namespace omogen
#endif
