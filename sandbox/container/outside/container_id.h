#pragma once

#include <memory>
#include <vector>

#include "absl/synchronization/mutex.h"

namespace omogen {
namespace sandbox {

class ContainerIds;

class ContainerId {
  int id;
  ContainerIds* containerIds;

 public:
  int Get() const { return id; }

  ContainerId(int id, ContainerIds* ptr) : id(id), containerIds(ptr) {}
  ContainerId(const ContainerId&) = delete;
  ContainerId& operator=(const ContainerId&) = delete;
  ContainerId(ContainerId&& other) {
    id = other.id;
    containerIds = other.containerIds;
    other.id = 0;
    other.containerIds = nullptr;
  }
  ContainerId& operator=(ContainerId&& other);
  ~ContainerId();
};

class ContainerIds {
  absl::Mutex mutex;
  std::vector<int> containerIds GUARDED_BY(mutex);

 public:
  explicit ContainerIds(int limit);
  std::unique_ptr<ContainerId> Get();
  void Put(int id);

  static std::unique_ptr<ContainerId> GetId();
};

}  // namespace sandbox
}  // namespace omogen
