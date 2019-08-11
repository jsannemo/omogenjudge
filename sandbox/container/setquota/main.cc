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
#include "util/cpp/files.h"

using omogen::util::TokenizeFile;
using omogen::util::WriteToFile;

using std::cerr;
using std::cout;
using std::endl;

const std::string kContainerPath = "/var/lib/omogen/sandbox";

int main(int argc, char** argv) {
  if (argc != 5) {
    cerr << "Incorrect number of arguments" << endl;
    return 1;
  }
  int pid_id;
  if (!absl::SimpleAtoi(std::string(argv[1]), &pid_id) || pid_id < 0) {
    cerr << "Invalid first argument" << endl;
    return 1;
  }
  int sandbox_id;
  if (!absl::SimpleAtoi(std::string(argv[2]), &sandbox_id) || sandbox_id < 0) {
    cerr << "Invalid second argument" << endl;
    return 1;
  }
  std::string user = absl::StrCat("omogenjudge-client", sandbox_id);
  struct passwd* pw = getpwnam(user.c_str());
  if (pw == NULL) {
    cerr << "Could not fetch uid" << endl;
    return 1;
  }
  cout << "Setting quota for for " << user << " with uid " << pw->pw_uid
       << endl;

  std::string uidpath = absl::StrCat("/proc/", pid_id, "/uid_map");
  WriteToFile(uidpath, absl::StrCat("65123 ", pw->pw_uid, " 1"));
  std::string gidpath = absl::StrCat("/proc/", pid_id, "/gid_map");
  WriteToFile(gidpath, absl::StrCat("65123 ", pw->pw_gid, " 1"));

  struct stat st;
  if (stat(kContainerPath.c_str(), &st) != 0) {
    cerr << "Could not stat submission folder" << endl;
    return 1;
  }

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
  if (device.empty()) {
    cerr << "Could not find device with the correct mount point" << endl;
    return 1;
  }
  struct stat dev_st;
  if (stat(device.c_str(), &dev_st) != 0) {
    cerr << "Could not stat device" << endl;
    return 1;
  }
  if (st.st_dev != dev_st.st_rdev) {
    cerr << "Device numbers inconsistent with mount point" << endl;
    return 1;
  }

  int block_quota;
  if (!absl::SimpleAtoi(std::string(argv[3]), &block_quota) ||
      block_quota < 0) {
    cerr << "Invalid third argument" << endl;
    return 1;
  }

  int inode_quota;
  if (!absl::SimpleAtoi(std::string(argv[4]), &inode_quota) ||
      inode_quota < 0) {
    cerr << "Invalid fourth argument" << endl;
    return 1;
  }

  struct dqblk quota;
  quota.dqb_bhardlimit = uint64_t(block_quota);
  quota.dqb_bsoftlimit = uint64_t(block_quota);
  quota.dqb_ihardlimit = uint64_t(inode_quota);
  quota.dqb_isoftlimit = uint64_t(inode_quota);
  quota.dqb_valid = QIF_LIMITS;
  if (quotactl(QCMD(Q_SETQUOTA, USRQUOTA), device.c_str(), pw->pw_uid,
               (caddr_t)&quota) != 0) {
    cerr << "Could not set quota correctly" << endl;
    return 1;
  }
}
