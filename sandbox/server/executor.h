#include <grpc/grpc.h>

#include <thread>

#include "absl/base/thread_annotations.h"
#include "absl/synchronization/mutex.h"
#include "sandbox/api/exec.grpc.pb.h"
#include "sandbox/container/outside/container.h"

using grpc::ServerContext;
using grpc::ServerReaderWriter;
using grpc::Status;

namespace omogen {
namespace sandbox {

class ExecuteServiceImpl final : public ExecuteService::Service {
  // The cleanup thread for this service, continously cleaning up old
  // containers.
  std::thread* cleanupThread;

  absl::Mutex deleteQueueMutex;

  // A vector with containers that should be cleaned up.
  std::vector<unique_ptr<Container>> deleteQueue GUARDED_BY(deleteQueueMutex);

  // Runs the cleanup loop the remove old containers.
  void cleanup();

 public:
  Status Execute(ServerContext* context,
                 ServerReaderWriter<ExecuteResponse, ExecuteRequest>* stream);

  ExecuteServiceImpl();
};

}  // namespace sandbox
}  // namespace omogen
