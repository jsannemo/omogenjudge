#include <pwd.h>
#include <sys/quota.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include <iostream>
#include <string>
#include <vector>

#include "absl/strings/match.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "glog/logging.h"
#include "util/cpp/files.h"

using omogen::util::TokenizeFile;

const std::string kContainerPath = "/var/lib/omogen/sandbox";

namespace omogen {
namespace sandbox {

void SetQuota(uint64_t block_quota, uint64_t inode_quota, int uid) {
  struct stat st;
  PCHECK(stat(kContainerPath.c_str(), &st) == 0)
      << "Could not stat submission folder";

  std::vector<std::string> toks = TokenizeFile("/proc/mounts");
  std::string device;
  for (size_t i = 0; i + 5 < toks.size(); i += 6) {
    std::string dev = toks[i];
    std::string mountPoint = toks[i + 1];
    if (absl::StartsWith(kContainerPath, mountPoint) &&
        dev.size() > device.size()) {
      device = dev;
    }
  }
  CHECK(!device.empty())
      << "Could not find device with the correct mount point";
  struct stat dev_st;
  PCHECK(stat(device.c_str(), &dev_st) == 0) << "Could not stat device";
  CHECK(st.st_dev == dev_st.st_rdev)
      << "Device numbers inconsistent with mount point";

  struct dqblk quota;
  quota.dqb_bhardlimit = uint64_t(block_quota);
  quota.dqb_bsoftlimit = uint64_t(block_quota);
  quota.dqb_ihardlimit = uint64_t(inode_quota);
  quota.dqb_isoftlimit = uint64_t(inode_quota);
  quota.dqb_valid = QIF_LIMITS;
  PCHECK(quotactl(QCMD(Q_SETQUOTA, USRQUOTA), device.c_str(), uid,
                  (caddr_t)&quota) == 0)
      << "Could not set quota correctly";
}

}  // namespace sandbox
}  // namespace omogen
