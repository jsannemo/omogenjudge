#include "sandbox/container/outside/container_id.h"

#include <memory>
#include <numeric>

#include "glog/logging.h"
#include "util/cpp/config.h"

namespace omogen {
namespace sandbox {

// TODO: this should be owned by the executor service rather than being global
ContainerIds CONTAINER_IDS(
    omogen::util::ReadConfigInt("/etc/omogen/sandbox/sandbox.conf", "range"));

ContainerIds::ContainerIds(int limit) : container_ids(limit) {
  std::iota(container_ids.begin(), container_ids.end(), 0);
}

static bool HasIds(std::vector<int>* v) { return !v->empty(); }

std::unique_ptr<ContainerId> ContainerIds::Get() {
  mutex.LockWhen(absl::Condition(HasIds, &container_ids));
  std::unique_ptr<ContainerId> ret =
      std::make_unique<ContainerId>(container_ids.back(), this);
  LOG(INFO) << "Claiming container ID " << ret->Get();
  container_ids.pop_back();
  mutex.Unlock();
  return ret;
}

ContainerId& ContainerId::operator=(ContainerId&& other) {
  if (this != &other) {
    if (container_ids != nullptr) {
      container_ids->Put(id);
    }
    id = other.id;
    container_ids = other.container_ids;
    other.container_ids = nullptr;
  }
  return *this;
}

void ContainerIds::Put(int id) {
  absl::MutexLock lock(&mutex);
  container_ids.push_back(id);
}

ContainerId::~ContainerId() {
  if (container_ids != nullptr) {
    LOG(INFO) << "Returning container ID " << id;
    container_ids->Put(id);
  }
}

std::unique_ptr<ContainerId> ContainerIds::GetId() {
  return CONTAINER_IDS.Get();
}

}  // namespace sandbox
}  // namespace omogen
