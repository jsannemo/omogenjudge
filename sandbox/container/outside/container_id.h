#ifndef SANDBOX_CONTAINER_OUTSIDE_CONTAINER_ID_H
#define SANDBOX_CONTAINER_OUTSIDE_CONTAINER_ID_H

#include <memory>
#include <vector>

#include "absl/synchronization/mutex.h"

namespace omogen {
namespace sandbox {

class ContainerIds;

// A container ID is a unique integer mapped to one of a limited number of
// containers on the system. Each container ID X has a corresponding system user
// omogenjudge-clientX allocated to it.
//
// Container IDs can not be copied, and should only ever be acquired from a
// ContainerIds instance. They are reclaimed when the instance is destroyed.
class ContainerId {
  int id;
  ContainerIds* container_ids;

 public:
  int Get() const { return id; }

  ContainerId(int id, ContainerIds* ptr) : id(id), container_ids(ptr) {}
  ContainerId(ContainerId&& other) {
    id = other.id;
    container_ids = other.container_ids;
    other.id = 0;
    other.container_ids = nullptr;
  }
  ContainerId& operator=(ContainerId&& other);
  ~ContainerId();

  ContainerId(const ContainerId&) = delete;
  ContainerId& operator=(const ContainerId&) = delete;
};

// A container of available container IDs.
// This is thread-safe.
class ContainerIds {
  absl::Mutex mutex;
  std::vector<int> container_ids GUARDED_BY(mutex);

 public:
  explicit ContainerIds(int limit);
  std::unique_ptr<ContainerId> Get();
  void Put(int id);

  static std::unique_ptr<ContainerId> GetId();

  ContainerIds(const ContainerIds&) = delete;
  ContainerIds& operator=(const ContainerId&) = delete;
};

}  // namespace sandbox
}  // namespace omogen
#endif
