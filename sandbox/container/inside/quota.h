#ifndef SANDBOX_CONTAINER_INSIDE_QUOTA_H
#define SANDBOX_CONTAINER_INSIDE_QUOTA_H

namespace omogen {
namespace sandbox {

void SetQuota(uint64_t block_quota, uint64_t inode_quota, int uid);

}  // namespace sandbox
}  // namespace omogen
#endif
