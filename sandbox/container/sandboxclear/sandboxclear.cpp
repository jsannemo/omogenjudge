#include <unistd.h>

#include <cstdlib>
#include <iostream>
#include <string>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"

using std::cerr;
using std::cout;
using std::endl;

const std::string kSandboxPath = "/var/lib/omogen/sandbox";

int main(int argc, char** argv) {
  if (argc != 2) {
    cerr << "Incorrect number of arguments" << endl;
    return 1;
  }
  int sandbox_id;
  if (!absl::SimpleAtoi(std::string(argv[1]), &sandbox_id) || sandbox_id < 0) {
    cerr << "Invalid first argument" << endl;
    return 1;
  }

  std::string path = absl::StrCat(kSandboxPath, "/", sandbox_id);

  if (setuid(0) != 0) {
    cerr << "Could not setuid" << endl;
  }

  // In case the files were marked immutable, we must unmark them, or else
  // the chown call will fail.
  std::string chattr_cmd = absl::StrCat("/usr/bin/chattr -i -R ", path);
  if (system(chattr_cmd.c_str())) {
    cerr << "Could not chattr" << endl;
    return 1;
  }

  // Remove the old submission
  std::string rm_cmd = absl::StrCat("/bin/rm -rf ", path);
  return system(rm_cmd.c_str());
}
