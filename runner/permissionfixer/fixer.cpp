#include <unistd.h>

#include <cstdlib>
#include <iostream>
#include <string>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"

using std::cerr;
using std::cout;
using std::endl;

const std::string SUBMISSION_PATH = "/var/lib/omogen/submissions";

int main(int argc, char** argv) {
  if (argc != 2) {
    cerr << "Incorrect number of arguments" << endl;
    return 1;
  }
  int submissionID;
  if (!absl::SimpleAtoi(std::string(argv[1]), &submissionID) ||
      submissionID < 0) {
    cerr << "Invalid first argument" << endl;
    return 1;
  }

  std::string path = absl::StrCat(SUBMISSION_PATH, "/", submissionID);

  if (setuid(0) != 0) {
    cerr << "Could not setuid" << endl;
  }

  std::string chattrCmd = absl::StrCat("/usr/bin/chattr -i -R ", path);
  if (system(chattrCmd.c_str())) {
    cerr << "Could not chattr" << endl;
    return 1;
  }

  std::string chownCmd =
      absl::StrCat("/bin/chown -R omogenjudge:omogenjudge ", path);
  return system(chownCmd.c_str());
}
