#include <unistd.h>

#include <algorithm>
#include <cstdlib>
#include <iostream>
#include <string>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"

using std::cerr;
using std::cout;
using std::endl;

const std::string kSubmissionPath = "/var/lib/omogen/submissions";

int main(int argc, char** argv) {
  if (argc != 2) {
    cerr << "Incorrect number of arguments" << endl;
    return 1;
  }
  std::string submission_id = std::string(argv[1]);
  if (std::find(submission_id.begin(), submission_id.end(), '/') !=
      submission_id.end()) {
    cerr << "Invalid first argument" << endl;
    return 1;
  }

  std::string path = absl::StrCat(kSubmissionPath, "/", submission_id);

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

  std::string chown_cmd = absl::StrCat(
      "/bin/chown -R omogenjudge-local:omogenjudge-clients ", path);
  if (system(chown_cmd.c_str())) {
    cerr << "Could not chown" << endl;
  }

  // We must have a mode that allows us to remove all the files afterwards as
  // well.
  std::string chmod_cmd = absl::StrCat("/bin/chmod -R gu+wrx ", path);
  return system(chmod_cmd.c_str());
}
